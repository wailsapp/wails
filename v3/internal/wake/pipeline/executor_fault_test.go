package pipeline

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/wake/cache"
)

type faultExecutorCache struct {
	snapshotDigest string
	snapshotErr    error
	snapshotCalls  int
	snapshotErrors map[int]error
	filesDigest    string
	filesErr       error
	goDigest       string
	goErr          error
	hasReceipt     bool
	receiptErr     error
	lookupStatus   cache.LookupStatus
	lookupArtifact string
	lookupErr      error
	peekStatus     cache.LookupStatus
	peekArtifact   string
	peekErr        error
	recordArtifact string
	recordErr      error
	saveErr        error
	begins         int
	invalidations  int
}

func (store *faultExecutorCache) BeginObservationSession() { store.begins++ }
func (store *faultExecutorCache) InvalidateObservations()  { store.invalidations++ }
func (store *faultExecutorCache) Snapshot(cache.SnapshotOptions) (string, error) {
	store.snapshotCalls++
	if err := store.snapshotErrors[store.snapshotCalls]; err != nil {
		return "", err
	}
	return store.snapshotDigest, store.snapshotErr
}
func (store *faultExecutorCache) SnapshotFiles(string, ...string) (string, error) {
	return store.filesDigest, store.filesErr
}
func (store *faultExecutorCache) SnapshotGoAPI(cache.SnapshotOptions) (string, error) {
	return store.goDigest, store.goErr
}
func (store *faultExecutorCache) HasReceipt(string, string) bool { return store.hasReceipt }
func (store *faultExecutorCache) RecordReceipt(string) error     { return store.receiptErr }
func (store *faultExecutorCache) Lookup(string, string) (cache.LookupStatus, string, error) {
	return store.lookupStatus, store.lookupArtifact, store.lookupErr
}
func (store *faultExecutorCache) Peek(string, string) (cache.LookupStatus, string, error) {
	return store.peekStatus, store.peekArtifact, store.peekErr
}
func (store *faultExecutorCache) RecordAction(string, string) (string, error) {
	return store.recordArtifact, store.recordErr
}
func (store *faultExecutorCache) Save() error { return store.saveErr }

type faultExecutorHandler struct {
	identity     string
	identityErr  error
	identityHook func()
	run          RunResult
	runErr       error
}

type delayedCancelHandler struct {
	started chan struct{}
	once    sync.Once
}

func (*delayedCancelHandler) Identity(context.Context, Node) (string, error) { return "cancel-v1", nil }
func (handler *delayedCancelHandler) Run(ctx context.Context, _ Node) (RunResult, error) {
	handler.once.Do(func() { close(handler.started) })
	<-ctx.Done()
	// Keep the acquired resource long enough for waiters to observe cancellation
	// before the first worker releases it.
	time.Sleep(20 * time.Millisecond)
	return RunResult{}, ctx.Err()
}

func (handler *faultExecutorHandler) Identity(context.Context, Node) (string, error) {
	if handler.identityHook != nil {
		handler.identityHook()
	}
	return handler.identity, handler.identityErr
}
func (handler *faultExecutorHandler) Run(context.Context, Node) (RunResult, error) {
	return handler.run, handler.runErr
}

func faultExecutorPlan(node Node) Plan {
	return Plan{Name: "fault", Roots: []NodeKey{node.Key}, Nodes: map[NodeKey]Node{node.Key: node}}
}

func fixedActionKey(string, any, []string, []string) (string, error) { return "action", nil }

func TestExecutorEntryPointsRejectInvalidStateAndCacheOpenFailures(t *testing.T) {
	want := errors.New("injected cache open failure")
	validNode := Node{Key: "work", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever}
	validPlan := faultExecutorPlan(validNode)

	_, err := (Executor{}).Inspect(t.Context(), validPlan, t.TempDir())
	require.ErrorContains(t, err, "requires a handler")
	_, err = (Executor{}).Execute(t.Context(), validPlan, ExecuteOptions{Root: t.TempDir()})
	require.ErrorContains(t, err, "requires a handler")

	handler := &faultExecutorHandler{}
	_, err = (Executor{Handler: handler}).Inspect(t.Context(), Plan{}, t.TempDir())
	require.ErrorContains(t, err, "no nodes")
	_, err = (Executor{Handler: handler}).Execute(t.Context(), Plan{}, ExecuteOptions{Root: t.TempDir()})
	require.ErrorContains(t, err, "no nodes")

	operations := &executorOperations{
		openReadOnly: func(string) (executorCache, error) { return nil, want },
		open:         func(string) (executorCache, error) { return nil, want },
		actionKey:    fixedActionKey,
	}
	executor := Executor{Handler: handler, operations: operations}
	_, err = executor.Inspect(t.Context(), validPlan, t.TempDir())
	require.ErrorIs(t, err, want)
	_, err = executor.Execute(t.Context(), validPlan, ExecuteOptions{Root: t.TempDir()})
	require.ErrorIs(t, err, want)
}

func TestInspectNodeCoversEveryCacheDecisionAndFailure(t *testing.T) {
	want := errors.New("injected inspection failure")
	baseNode := Node{Key: "work", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheArtifact, Output: "bin/app"}
	handler := &faultExecutorHandler{identity: "tool-v1"}
	executor := Executor{Handler: handler}
	operations := map[NodeKey]OperationInspection{}
	results := map[NodeKey]Result{}

	store := &faultExecutorCache{snapshotErr: want}
	node := baseNode
	node.Inputs = []InputSpec{{Label: "source", Root: ".", IncludeAll: true}}
	_, _, err := executor.inspectNode(t.Context(), store, node, operations, results, fixedActionKey)
	require.ErrorIs(t, err, want)

	store = &faultExecutorCache{}
	node = baseNode
	node.Cache = CacheNever
	operation, result, err := executor.inspectNode(t.Context(), store, node, operations, results, fixedActionKey)
	require.NoError(t, err)
	assert.Equal(t, "run", operation.Decision)
	assert.Equal(t, node.Output, result.Output)

	node = baseNode
	node.Dependencies = []NodeKey{"dependency"}
	operation, _, err = executor.inspectNode(t.Context(), store, node, map[NodeKey]OperationInspection{"dependency": {Decision: "run"}}, results, fixedActionKey)
	require.NoError(t, err)
	assert.Equal(t, "run", operation.Decision)

	node.Dependencies = nil
	handler.identityErr = want
	_, _, err = executor.inspectNode(t.Context(), store, node, operations, results, fixedActionKey)
	require.ErrorIs(t, err, want)
	handler.identityErr = nil

	actionFailure := func(string, any, []string, []string) (string, error) { return "", want }
	_, _, err = executor.inspectNode(t.Context(), store, node, operations, results, actionFailure)
	require.ErrorIs(t, err, want)

	for name, state := range map[string]struct {
		status   cache.LookupStatus
		decision string
	}{
		"hit":      {cache.LookupHit, "cached"},
		"restored": {cache.LookupRestored, "restore"},
		"dirty":    {cache.LookupDirty, "run"},
		"miss":     {cache.LookupMiss, "run"},
	} {
		t.Run(name, func(t *testing.T) {
			store := &faultExecutorCache{peekStatus: state.status, peekArtifact: "artifact"}
			operation, result, err := executor.inspectNode(t.Context(), store, baseNode, operations, results, fixedActionKey)
			require.NoError(t, err)
			assert.Equal(t, state.decision, operation.Decision)
			assert.Equal(t, state.status, result.Status)
		})
	}
	store = &faultExecutorCache{peekErr: want}
	_, _, err = executor.inspectNode(t.Context(), store, baseNode, operations, results, fixedActionKey)
	require.ErrorIs(t, err, want)

	receipt := baseNode
	receipt.Cache = CacheReceipt
	receipt.Marker = "node_modules/.receipt"
	receipt.Inputs = []InputSpec{{Root: ".", IncludeAll: true}}
	for _, hit := range []bool{false, true} {
		store = &faultExecutorCache{hasReceipt: hit}
		operation, result, err = executor.inspectNode(t.Context(), store, receipt, operations, results, fixedActionKey)
		require.NoError(t, err)
		if hit {
			assert.Equal(t, "cached", operation.Decision)
			assert.Equal(t, cache.LookupHit, result.Status)
		} else {
			assert.Equal(t, "run", operation.Decision)
		}
	}

	uncachedOutput := baseNode
	uncachedOutput.Output = ""
	operation, _, err = executor.inspectNode(t.Context(), store, uncachedOutput, operations, results, fixedActionKey)
	require.NoError(t, err)
	assert.Equal(t, "run", operation.Decision)

	dependency := baseNode
	dependency.Dependencies = []NodeKey{"dependency"}
	store = &faultExecutorCache{peekStatus: cache.LookupHit}
	_, _, err = executor.inspectNode(t.Context(), store, dependency,
		map[NodeKey]OperationInspection{"dependency": {Decision: "cached"}},
		map[NodeKey]Result{"dependency": {Artifact: "dependency-artifact"}}, fixedActionKey)
	require.NoError(t, err)
}

func TestInspectWrapsNodeFailuresAndBeginsOneObservation(t *testing.T) {
	want := errors.New("injected snapshot failure")
	store := &faultExecutorCache{snapshotErr: want}
	operations := &executorOperations{
		openReadOnly: func(string) (executorCache, error) { return store, nil },
		open:         func(string) (executorCache, error) { return store, nil },
		actionKey:    fixedActionKey,
	}
	node := Node{Key: "work", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheArtifact, Output: "bin/app", Inputs: []InputSpec{{Root: ".", IncludeAll: true}}}
	_, err := (Executor{Handler: &faultExecutorHandler{}, operations: operations}).Inspect(t.Context(), faultExecutorPlan(node), t.TempDir())
	require.ErrorContains(t, err, "work")
	require.ErrorIs(t, err, want)
	assert.Equal(t, 1, store.begins)
}

func TestInspectOrdersIndependentNodesAndDetectsPostValidationCycles(t *testing.T) {
	store := &faultExecutorCache{peekStatus: cache.LookupMiss}
	plan := Plan{Name: "ordered", Roots: []NodeKey{"child", "z-independent"}, Nodes: map[NodeKey]Node{
		"a-parent":      {Key: "a-parent", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever},
		"child":         {Key: "child", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever, Dependencies: []NodeKey{"a-parent"}},
		"z-independent": {Key: "z-independent", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever},
	}}
	operations := &executorOperations{
		openReadOnly: func(string) (executorCache, error) { return store, nil },
		open:         func(string) (executorCache, error) { return store, nil },
		actionKey:    fixedActionKey,
	}
	inspection, err := (Executor{Handler: &faultExecutorHandler{}, operations: operations}).Inspect(t.Context(), plan, t.TempDir())
	require.NoError(t, err)
	assert.Len(t, inspection.Operations, 3)

	cycle := Plan{Name: "cycle-after-validation", Roots: []NodeKey{"work"}, Nodes: map[NodeKey]Node{
		"work": {Key: "work", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever},
	}}
	operations.openReadOnly = func(string) (executorCache, error) {
		node := cycle.Nodes["work"]
		node.Dependencies = []NodeKey{"work"}
		cycle.Nodes["work"] = node
		return store, nil
	}
	_, err = (Executor{Handler: &faultExecutorHandler{}, operations: operations}).Inspect(t.Context(), cycle, t.TempDir())
	require.ErrorContains(t, err, "could not traverse")
}

func TestRunNodePropagatesEveryHandlerAndCacheFailure(t *testing.T) {
	want := errors.New("injected node failure")
	baseNode := Node{Key: "work", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheArtifact, Output: "bin/app"}
	cacheMu := &sync.Mutex{}
	reporter := &failureReporter{}

	run := func(handler *faultExecutorHandler, store *faultExecutorCache, node Node, actionKey func(string, any, []string, []string) (string, error), force bool) (Result, error) {
		return (Executor{Handler: handler}).runNode(t.Context(), store, cacheMu, node, nil, force, reporter, actionKey)
	}
	handler := &faultExecutorHandler{identity: "tool-v1"}
	store := &faultExecutorCache{}
	handler.identityErr = want
	_, err := run(handler, store, baseNode, fixedActionKey, false)
	require.ErrorIs(t, err, want)
	handler.identityErr = nil
	node := baseNode
	node.Inputs = []InputSpec{{Root: ".", IncludeAll: true}}
	store.snapshotErr = want
	_, err = run(handler, store, node, fixedActionKey, false)
	require.ErrorIs(t, err, want)
	store.snapshotErr = nil
	_, err = run(handler, store, node, func(string, any, []string, []string) (string, error) { return "", want }, false)
	require.ErrorIs(t, err, want)

	store = &faultExecutorCache{lookupErr: want}
	_, err = run(handler, store, baseNode, fixedActionKey, false)
	require.ErrorIs(t, err, want)
	store = &faultExecutorCache{lookupStatus: cache.LookupDirty, recordArtifact: "artifact"}
	result, err := run(handler, store, baseNode, fixedActionKey, false)
	require.NoError(t, err)
	assert.Equal(t, "artifact", result.Artifact)
	assert.Equal(t, 1, store.invalidations)
	store = &faultExecutorCache{lookupStatus: cache.LookupHit, lookupArtifact: "artifact"}
	result, err = run(handler, store, baseNode, fixedActionKey, false)
	require.NoError(t, err)
	assert.Equal(t, cache.LookupHit, result.Status)

	handler.run = RunResult{Detail: "compiler detail"}
	handler.runErr = want
	_, err = run(handler, &faultExecutorCache{}, Node{Key: "never", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever}, fixedActionKey, false)
	require.ErrorIs(t, err, want)
	assert.Equal(t, "compiler detail", reporter.failure.Output)
	handler.runErr = nil
	result, err = run(handler, &faultExecutorCache{}, Node{Key: "never", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever}, fixedActionKey, false)
	require.NoError(t, err)
	assert.Equal(t, NodeKey("never"), result.Key)
	handler.run = RunResult{}

	receipt := baseNode
	receipt.Cache = CacheReceipt
	receipt.Marker = "node_modules/.receipt"
	receipt.Inputs = []InputSpec{{Root: ".", IncludeAll: true}}
	store = &faultExecutorCache{hasReceipt: true}
	result, err = run(handler, store, receipt, fixedActionKey, false)
	require.NoError(t, err)
	assert.Equal(t, cache.LookupHit, result.Status)
	store = &faultExecutorCache{snapshotErrors: map[int]error{2: want}}
	_, err = run(handler, store, receipt, fixedActionKey, true)
	require.ErrorIs(t, err, want)
	store = &faultExecutorCache{}
	actionCalls := 0
	_, err = run(handler, store, receipt, func(string, any, []string, []string) (string, error) {
		actionCalls++
		if actionCalls == 2 {
			return "", want
		}
		return "action", nil
	}, false)
	require.ErrorIs(t, err, want)
	store = &faultExecutorCache{receiptErr: want}
	_, err = run(handler, store, receipt, fixedActionKey, false)
	require.ErrorIs(t, err, want)
	store = &faultExecutorCache{}
	_, err = run(handler, store, receipt, fixedActionKey, false)
	require.NoError(t, err)

	store = &faultExecutorCache{recordErr: want}
	_, err = run(handler, store, baseNode, fixedActionKey, true)
	require.ErrorIs(t, err, want)
}

func TestSnapshotNodeInputsCoversAllKindsAndStopsOnFailure(t *testing.T) {
	store := &faultExecutorCache{snapshotDigest: "tree", filesDigest: "files", goDigest: "go"}
	node := Node{Inputs: []InputSpec{
		{Label: "go", Root: ".", SemanticGo: true},
		{Label: "files", Files: []string{"go.mod"}},
		{Label: "tree", Root: ".", IncludeAll: true},
	}}
	digests, err := snapshotNodeInputs(store, node)
	require.NoError(t, err)
	assert.Equal(t, []string{"go", "files", "tree"}, digests)
	store.filesErr = errors.New("files failed")
	_, err = snapshotNodeInputs(store, node)
	require.ErrorContains(t, err, "files failed")
}

func TestExecutorReportsSaveFailuresAndRecursiveDependencyBlocking(t *testing.T) {
	want := errors.New("injected save failure")
	store := &faultExecutorCache{saveErr: want}
	operations := &executorOperations{
		openReadOnly: func(string) (executorCache, error) { return store, nil },
		open:         func(string) (executorCache, error) { return store, nil },
		actionKey:    fixedActionKey,
	}
	node := Node{Key: "work", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever}
	_, err := (Executor{Handler: &faultExecutorHandler{}, operations: operations}).Execute(t.Context(), faultExecutorPlan(node), ExecuteOptions{Root: t.TempDir(), Workers: 1})
	require.ErrorIs(t, err, want)

	finished := 0
	failed := map[NodeKey]error{}
	blocked := map[NodeKey]bool{}
	indegree := map[NodeKey]int{"child": 1, "grandchild": 1}
	markBlocked("root", map[NodeKey][]NodeKey{"root": {"child"}, "child": {"grandchild"}}, indegree, blocked, failed, &finished)
	assert.Equal(t, 2, finished)
	assert.True(t, blocked["child"])
	assert.True(t, blocked["grandchild"])
	assert.ErrorContains(t, failed["grandchild"], "blocked")
}

func TestExecutorCancellationReleasesCPUAndMemoryWaiters(t *testing.T) {
	for _, scenario := range []struct {
		name        string
		claims      ResourceClaims
		memoryLimit int64
	}{
		{name: "cpu", claims: ResourceClaims{CPU: 2, MemoryMB: 1}, memoryLimit: 2},
		{name: "memory", claims: ResourceClaims{CPU: 1, MemoryMB: 3}, memoryLimit: 2},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			plan := Plan{Name: scenario.name, Roots: []NodeKey{"one", "two"}, Nodes: map[NodeKey]Node{
				"one": {Key: "one", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever, Claims: scenario.claims},
				"two": {Key: "two", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever, Claims: scenario.claims},
			}}
			handler := &delayedCancelHandler{started: make(chan struct{})}
			ctx, cancel := context.WithCancel(context.Background())
			done := make(chan error, 1)
			go func() {
				_, err := (Executor{Handler: handler}).Execute(ctx, plan, ExecuteOptions{Root: t.TempDir(), Workers: 2, MemoryLimitMB: scenario.memoryLimit})
				done <- err
			}()
			<-handler.started
			time.Sleep(20 * time.Millisecond)
			cancel()
			require.ErrorIs(t, <-done, context.Canceled)
		})
	}
}

func TestExecutorDetectsSchedulerStateMutatedAfterValidation(t *testing.T) {
	store := &faultExecutorCache{}
	plan := Plan{Name: "mutated", Roots: []NodeKey{"work"}, Nodes: map[NodeKey]Node{
		"work": {Key: "work", Kind: CompileApplication, Spec: CompileSpec{}, Cache: CacheNever},
	}}
	operations := &executorOperations{
		openReadOnly: func(string) (executorCache, error) { return store, nil },
		open: func(string) (executorCache, error) {
			node := plan.Nodes["work"]
			node.Dependencies = []NodeKey{"missing-after-validation"}
			plan.Nodes["work"] = node
			return store, nil
		},
		actionKey: fixedActionKey,
	}
	results, err := (Executor{Handler: &faultExecutorHandler{}, operations: operations}).Execute(t.Context(), plan, ExecuteOptions{Root: t.TempDir(), Workers: 1})
	require.NoError(t, err)
	assert.Empty(t, results, "the defensive scheduler guard exits when post-validation mutation leaves no runnable node")
}

type testExitError struct{ code int }

func (errorValue testExitError) Error() string { return fmt.Sprintf("exit %d", errorValue.code) }
func (errorValue testExitError) ExitCode() int { return errorValue.code }

func TestExitCodeUsesTypedProcessFailures(t *testing.T) {
	assert.Equal(t, 42, exitCode(testExitError{code: 42}))
	assert.Equal(t, -1, exitCode(errors.New("plain")))
}

var _ report.Reporter = (*failureReporter)(nil)
