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

type Executor struct {
	Handler    Handler
	operations *executorOperations
}

type executorCache interface {
	BeginObservationSession()
	InvalidateObservations()
	Snapshot(cache.SnapshotOptions) (string, error)
	SnapshotFiles(string, ...string) (string, error)
	SnapshotGoAPI(cache.SnapshotOptions) (string, error)
	HasReceipt(string, string) bool
	RecordReceipt(string) error
	Lookup(string, string) (cache.LookupStatus, string, error)
	Peek(string, string) (cache.LookupStatus, string, error)
	RecordAction(string, string) (string, error)
	Save() error
}

type executorOperations struct {
	openReadOnly func(string) (executorCache, error)
	open         func(string) (executorCache, error)
	actionKey    func(string, any, []string, []string) (string, error)
}

func (e Executor) cacheOperations() executorOperations {
	if e.operations != nil {
		return *e.operations
	}
	return executorOperations{
		openReadOnly: func(root string) (executorCache, error) { return cache.OpenCacheReadOnly(root) },
		open:         func(root string) (executorCache, error) { return cache.OpenCache(root) },
		actionKey:    cache.ActionKey,
	}
}

type nodeOutcome struct {
	key    NodeKey
	result Result
	err    error
}

func (e Executor) Inspect(ctx context.Context, plan Plan, root string) (Inspection, error) {
	if e.Handler == nil {
		return Inspection{}, fmt.Errorf("wake: inspector requires a handler")
	}
	if err := plan.Validate(root); err != nil {
		return Inspection{}, err
	}
	operations := e.cacheOperations()
	store, err := operations.openReadOnly(root)
	if err != nil {
		return Inspection{}, err
	}
	store.BeginObservationSession()
	indegree := make(map[NodeKey]int, len(plan.Nodes))
	dependents := make(map[NodeKey][]NodeKey, len(plan.Nodes))
	for key, node := range plan.Nodes {
		indegree[key] = len(node.Dependencies)
		for _, dependency := range node.Dependencies {
			dependents[dependency] = append(dependents[dependency], key)
		}
	}
	ready := make([]NodeKey, 0, len(plan.Nodes))
	for key, degree := range indegree {
		if degree == 0 {
			ready = append(ready, key)
		}
	}
	sort.Slice(ready, func(i, j int) bool { return ready[i] < ready[j] })
	inspection := Inspection{Operations: make(map[NodeKey]OperationInspection, len(plan.Nodes))}
	results := make(map[NodeKey]Result, len(plan.Nodes))
	for len(ready) != 0 {
		key := ready[0]
		ready = ready[1:]
		node := plan.Nodes[key]
		operation, result, inspectErr := e.inspectNode(ctx, store, node, inspection.Operations, results, operations.actionKey)
		if inspectErr != nil {
			return Inspection{}, fmt.Errorf("%s: %w", key, inspectErr)
		}
		inspection.Operations[key], results[key] = operation, result
		for _, dependent := range dependents[key] {
			indegree[dependent]--
			if indegree[dependent] == 0 {
				ready = append(ready, dependent)
				sort.Slice(ready, func(i, j int) bool { return ready[i] < ready[j] })
			}
		}
	}
	if len(inspection.Operations) != len(plan.Nodes) {
		return Inspection{}, fmt.Errorf("wake: inspection could not traverse the complete plan")
	}
	return inspection, nil
}

func (e Executor) inspectNode(ctx context.Context, store executorCache, node Node, operations map[NodeKey]OperationInspection, results map[NodeKey]Result, actionKey func(string, any, []string, []string) (string, error)) (OperationInspection, Result, error) {
	inputs, err := snapshotNodeInputs(store, node)
	if err != nil {
		return OperationInspection{}, Result{}, err
	}
	operation := OperationInspection{Decision: "run", Status: cache.LookupMiss, Inputs: make([]InputSnapshot, len(inputs))}
	for index, digest := range inputs {
		operation.Inputs[index] = InputSnapshot{Label: node.Inputs[index].Label, Digest: digest}
	}
	if node.Cache == CacheNever {
		return operation, Result{Key: node.Key, Status: cache.LookupMiss, Output: node.Output}, nil
	}
	for _, dependency := range node.Dependencies {
		if operations[dependency].Decision == "run" {
			return operation, Result{Key: node.Key, Status: cache.LookupMiss, Output: node.Output}, nil
		}
	}
	identity, err := e.Handler.Identity(ctx, node)
	if err != nil {
		return OperationInspection{}, Result{}, err
	}
	dependencyArtifacts := make([]string, 0, len(node.Dependencies))
	for _, dependency := range node.Dependencies {
		if artifact := results[dependency].Artifact; artifact != "" {
			dependencyArtifacts = append(dependencyArtifacts, artifact)
		}
	}
	action, err := actionKey(string(node.Kind), map[string]any{"spec": node.Spec, "tool": identity}, inputs, dependencyArtifacts)
	if err != nil {
		return OperationInspection{}, Result{}, err
	}
	if node.Cache == CacheReceipt {
		if store.HasReceipt(action, node.Marker) {
			operation.Decision, operation.Status = "cached", cache.LookupHit
		}
		return operation, Result{Key: node.Key, Status: operation.Status}, nil
	}
	if node.Cache == CacheArtifact && node.Output != "" {
		status, artifact, err := store.Peek(action, node.Output)
		if err != nil {
			return OperationInspection{}, Result{}, err
		}
		operation.Status = status
		switch status {
		case cache.LookupHit:
			operation.Decision = "cached"
		case cache.LookupRestored:
			operation.Decision = "restore"
		}
		return operation, Result{Key: node.Key, Status: status, Artifact: artifact, Output: node.Output}, nil
	}
	return operation, Result{Key: node.Key, Status: cache.LookupMiss, Output: node.Output}, nil
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
	operations := e.cacheOperations()
	cacheStore, err := operations.open(options.Root)
	if err != nil {
		return nil, err
	}
	cacheStore.BeginObservationSession()

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
				result, err := e.runNode(ctx, cacheStore, &cacheMu, node, dependencyResults, options.Force, reporter, operations.actionKey)
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

func (e Executor) runNode(ctx context.Context, store executorCache, cacheMu *sync.Mutex, node Node, deps map[NodeKey]Result, force bool, reporter report.Reporter, actionKey func(string, any, []string, []string) (string, error)) (Result, error) {
	id := reporter.StepStart(string(node.Key), node.Label)
	started := time.Now()
	fail := func(err error, output string) (Result, error) {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return Result{}, ctxErr
		}
		reporter.StepFailed(id, report.Failure{Task: string(node.Key), Err: err, ExitCode: exitCode(err), Output: output})
		return Result{}, err
	}
	var identity, action string
	var depArtifacts []string
	if node.Cache != CacheNever {
		var err error
		identity, err = e.Handler.Identity(ctx, node)
		if err != nil {
			return fail(err, "")
		}
		cacheMu.Lock()
		inputs, snapshotErr := snapshotNodeInputs(store, node)
		cacheMu.Unlock()
		if snapshotErr != nil {
			return fail(snapshotErr, "")
		}
		depArtifacts = make([]string, 0, len(node.Dependencies))
		for _, key := range node.Dependencies {
			if artifact := deps[key].Artifact; artifact != "" {
				depArtifacts = append(depArtifacts, artifact)
			}
		}
		action, err = actionKey(string(node.Kind), map[string]any{"spec": node.Spec, "tool": identity}, inputs, depArtifacts)
		if err != nil {
			return fail(err, "")
		}
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
	cacheMu.Lock()
	store.InvalidateObservations()
	cacheMu.Unlock()
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
			action, snapshotErr = actionKey(string(node.Kind), map[string]any{"spec": node.Spec, "tool": identity}, postInputs, depArtifacts)
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

func snapshotNodeInputs(store executorCache, node Node) ([]string, error) {
	inputs := make([]string, 0, len(node.Inputs))
	for _, input := range node.Inputs {
		var (
			digest string
			err    error
		)
		if input.SemanticGo {
			digest, err = store.SnapshotGoAPI(cache.SnapshotOptions{Label: input.Label, Root: input.Root, ExcludeDirs: input.ExcludeDirs, UseGitIgnore: input.UseGitIgnore})
		} else if len(input.Files) > 0 {
			digest, err = store.SnapshotFiles(input.Label, input.Files...)
		} else {
			digest, err = store.Snapshot(cache.SnapshotOptions{Label: input.Label, Root: input.Root, IncludeAll: input.IncludeAll, IncludeNames: input.IncludeNames, IncludeExtensions: input.IncludeExtensions, ExcludeDirs: input.ExcludeDirs, ExcludeSuffixes: input.ExcludeSuffixes, UseGitIgnore: input.UseGitIgnore})
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
