package pipeline

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
	"golang.org/x/sync/semaphore"
)

type Executor struct{ Handler Handler }

type nodeOutcome struct {
	key    NodeKey
	result Result
	err    error
}

func (e Executor) Execute(ctx context.Context, plan Plan, options ExecuteOptions) (map[NodeKey]Result, error) {
	if e.Handler == nil {
		return nil, fmt.Errorf("wake: executor requires a handler")
	}
	if err := plan.Validate(options.Root); err != nil {
		return nil, err
	}
	workers := options.Workers
	if workers <= 0 {
		workers = max(1, runtime.GOMAXPROCS(0))
	}
	reporter := options.Reporter
	if reporter == nil {
		reporter = report.Nop{}
	}
	cacheStore, err := cache.OpenCache(options.Root)
	if err != nil {
		return nil, err
	}

	indegree := map[NodeKey]int{}
	dependents := map[NodeKey][]NodeKey{}
	for key, node := range plan.Nodes {
		indegree[key] = len(node.Dependencies)
		for _, dep := range node.Dependencies {
			dependents[dep] = append(dependents[dep], key)
		}
	}
	critical := criticalPaths(plan, dependents)
	ready := make([]NodeKey, 0)
	for key, n := range indegree {
		if n == 0 {
			ready = append(ready, key)
		}
	}
	sortReady(ready, critical)
	results := map[NodeKey]Result{}
	var resultsMu sync.RWMutex
	var cacheMu sync.Mutex
	cpu := semaphore.NewWeighted(int64(workers))
	memoryLimit := options.MemoryLimitMB
	if memoryLimit <= 0 {
		// One logical GiB per worker is a portable default capacity. Explicit
		// callers and constrained runners can provide a tighter real budget.
		memoryLimit = int64(workers) * 1024
	}
	memory := semaphore.NewWeighted(memoryLimit)
	var exclusiveMu sync.Mutex
	exclusive := map[string]*sync.Mutex{}
	failed := map[NodeKey]error{}
	blocked := map[NodeKey]bool{}
	jobs := make(chan NodeKey)
	outcomes := make(chan nodeOutcome, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for key := range jobs {
				node := plan.Nodes[key]
				claim := node.Claims.CPU
				if claim < 1 {
					claim = 1
				}
				if claim > workers {
					claim = workers
				}
				if err := cpu.Acquire(ctx, int64(claim)); err != nil {
					outcomes <- nodeOutcome{key: key, err: err}
					continue
				}
				memoryClaim := node.Claims.MemoryMB
				if memoryClaim < 1 {
					memoryClaim = 1
				}
				if memoryClaim > memoryLimit {
					memoryClaim = memoryLimit
				}
				if err := memory.Acquire(ctx, memoryClaim); err != nil {
					cpu.Release(int64(claim))
					outcomes <- nodeOutcome{key: key, err: err}
					continue
				}
				var exclusiveLock *sync.Mutex
				if node.Claims.Exclusive != "" {
					exclusiveMu.Lock()
					exclusiveLock = exclusive[node.Claims.Exclusive]
					if exclusiveLock == nil {
						exclusiveLock = &sync.Mutex{}
						exclusive[node.Claims.Exclusive] = exclusiveLock
					}
					exclusiveMu.Unlock()
					exclusiveLock.Lock()
				}
				dependencyResults := make(map[NodeKey]Result, len(node.Dependencies))
				resultsMu.RLock()
				for _, dependency := range node.Dependencies {
					dependencyResults[dependency] = results[dependency]
				}
				resultsMu.RUnlock()
				result, err := e.runNode(ctx, cacheStore, &cacheMu, node, dependencyResults, options.Force, reporter)
				if exclusiveLock != nil {
					exclusiveLock.Unlock()
				}
				memory.Release(memoryClaim)
				cpu.Release(int64(claim))
				outcomes <- nodeOutcome{key: key, result: result, err: err}
			}
		}()
	}
	running, finished := 0, 0
	for finished < len(plan.Nodes) {
		for running < workers && len(ready) > 0 {
			key := ready[0]
			ready = ready[1:]
			jobs <- key
			running++
		}
		if running == 0 {
			break
		}
		outcome := <-outcomes
		running--
		finished++
		if outcome.err != nil {
			failed[outcome.key] = outcome.err
		} else {
			resultsMu.Lock()
			results[outcome.key] = outcome.result
			resultsMu.Unlock()
		}
		for _, child := range dependents[outcome.key] {
			indegree[child]--
			if outcome.err != nil || blocked[outcome.key] {
				blocked[child] = true
			}
			if indegree[child] == 0 {
				if blocked[child] {
					failed[child] = fmt.Errorf("blocked by failed dependency")
					finished++
					markBlocked(child, dependents, indegree, blocked, failed, &finished)
				} else {
					ready = append(ready, child)
					sortReady(ready, critical)
				}
			}
		}
	}
	close(jobs)
	wg.Wait()
	if err := cacheStore.Save(); err != nil && len(failed) == 0 {
		return results, err
	}
	if len(failed) > 0 {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return results, ctxErr
		}
		keys := make([]string, 0, len(failed))
		for key := range failed {
			keys = append(keys, string(key))
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, key := range keys {
			parts = append(parts, key+": "+failed[NodeKey(key)].Error())
		}
		return results, errors.New(strings.Join(parts, "; "))
	}
	return results, nil
}

func (e Executor) runNode(ctx context.Context, store *cache.Cache, cacheMu *sync.Mutex, node Node, deps map[NodeKey]Result, force bool, reporter report.Reporter) (Result, error) {
	id := reporter.StepStart(string(node.Key), node.Label)
	started := time.Now()
	fail := func(err error, output string) (Result, error) {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return Result{}, ctxErr
		}
		reporter.StepFailed(id, report.Failure{Task: string(node.Key), Err: err, ExitCode: exitCode(err), Output: output})
		return Result{}, err
	}
	identity, err := e.Handler.Identity(ctx, node)
	if err != nil {
		return fail(err, "")
	}
	cacheMu.Lock()
	inputs, err := snapshotNodeInputs(store, node)
	cacheMu.Unlock()
	if err != nil {
		return fail(err, "")
	}
	depArtifacts := make([]string, 0, len(node.Dependencies))
	for _, key := range node.Dependencies {
		if artifact := deps[key].Artifact; artifact != "" {
			depArtifacts = append(depArtifacts, artifact)
		}
	}
	action, err := cache.ActionKey(string(node.Kind), map[string]any{"spec": node.Spec, "tool": identity}, inputs, depArtifacts)
	if err != nil {
		return fail(err, "")
	}
	if !force && node.Cache == CacheReceipt {
		cacheMu.Lock()
		hit := store.HasReceipt(action, node.Marker)
		cacheMu.Unlock()
		if !hit {
			goto execute
		}
		reporter.StepEnd(id, report.StatusCached, time.Since(started))
		return Result{Key: node.Key, Status: cache.LookupHit}, nil
	}
	if !force && node.Cache == CacheArtifact && node.Output != "" {
		cacheMu.Lock()
		status, artifact, lookupErr := store.Lookup(action, node.Output)
		cacheMu.Unlock()
		if lookupErr != nil {
			return fail(lookupErr, "")
		}
		if status == cache.LookupHit || status == cache.LookupRestored {
			reporter.StepInfo(id, string(status))
			reporter.StepEnd(id, report.StatusCached, time.Since(started))
			return Result{Key: node.Key, Status: status, Artifact: artifact, Output: node.Output}, nil
		}
		if status == cache.LookupDirty {
			reporter.StepInfo(id, "generated output modified; rebuilding")
		}
	}

execute:
	run, err := e.Handler.Run(ctx, node)
	if err != nil {
		return fail(err, run.Detail)
	}
	if run.Detail != "" {
		reporter.StepInfo(id, run.Detail)
	}
	result := Result{Key: node.Key, Status: cache.LookupMiss, Output: node.Output}
	switch node.Cache {
	case CacheReceipt:
		cacheMu.Lock()
		// Installers may create or rewrite a lockfile. Store the Receipt under
		// the post-run identity so the immediately following build is a hit.
		postInputs, snapshotErr := snapshotNodeInputs(store, node)
		if snapshotErr == nil {
			action, snapshotErr = cache.ActionKey(string(node.Kind), map[string]any{"spec": node.Spec, "tool": identity}, postInputs, depArtifacts)
		}
		if snapshotErr != nil {
			cacheMu.Unlock()
			return fail(snapshotErr, "")
		}
		if err := store.RecordReceipt(action); err != nil {
			cacheMu.Unlock()
			return fail(err, "")
		}
		cacheMu.Unlock()
	case CacheArtifact:
		cacheMu.Lock()
		artifact, err := store.RecordAction(action, node.Output)
		cacheMu.Unlock()
		if err != nil {
			return fail(err, "")
		}
		result.Artifact = artifact
	}
	reporter.StepEnd(id, report.StatusOK, time.Since(started))
	return result, nil
}

func snapshotNodeInputs(store *cache.Cache, node Node) ([]string, error) {
	inputs := make([]string, 0, len(node.Inputs))
	for _, input := range node.Inputs {
		var (
			digest string
			err    error
		)
		if input.SemanticGo {
			digest, err = store.SnapshotGoAPI(input.Label, input.Root, input.ExcludeDirs)
		} else if len(input.Files) > 0 {
			digest, err = store.SnapshotFiles(input.Label, input.Files...)
		} else {
			digest, err = store.Snapshot(cache.SnapshotOptions{Label: input.Label, Root: input.Root, IncludeAll: input.IncludeAll, IncludeNames: input.IncludeNames, IncludeExtensions: input.IncludeExtensions, ExcludeDirs: input.ExcludeDirs})
		}
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, digest)
	}
	return inputs, nil
}

func criticalPaths(plan Plan, dependents map[NodeKey][]NodeKey) map[NodeKey]int64 {
	memo := map[NodeKey]int64{}
	var visit func(NodeKey) int64
	visit = func(key NodeKey) int64 {
		if v, ok := memo[key]; ok {
			return v
		}
		best := int64(0)
		for _, child := range dependents[key] {
			if v := visit(child); v > best {
				best = v
			}
		}
		memo[key] = plan.Nodes[key].EstimateMS + best
		return memo[key]
	}
	for key := range plan.Nodes {
		visit(key)
	}
	return memo
}
func sortReady(ready []NodeKey, critical map[NodeKey]int64) {
	sort.SliceStable(ready, func(i, j int) bool {
		if critical[ready[i]] == critical[ready[j]] {
			return ready[i] < ready[j]
		}
		return critical[ready[i]] > critical[ready[j]]
	})
}
func markBlocked(key NodeKey, dependents map[NodeKey][]NodeKey, indegree map[NodeKey]int, blocked map[NodeKey]bool, failed map[NodeKey]error, finished *int) {
	for _, child := range dependents[key] {
		indegree[child]--
		blocked[child] = true
		if indegree[child] == 0 {
			failed[child] = fmt.Errorf("blocked by failed dependency")
			*finished++
			markBlocked(child, dependents, indegree, blocked, failed, finished)
		}
	}
}

type exitCoder interface{ ExitCode() int }

func exitCode(err error) int {
	var value exitCoder
	if errors.As(err, &value) {
		return value.ExitCode()
	}
	return -1
}
