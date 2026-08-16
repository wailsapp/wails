package benchmark

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBuildOutputCapturesWorkAndCriticalPathDuration(t *testing.T) {
	sample := Sample{}
	parseBuildOutput("\x1b[31m✓  build succeeded    5ms    0 ran  •  4 cached\x1b[0m", &sample)
	assert.Equal(t, 0, sample.ExecutedSteps)
	assert.Equal(t, 4, sample.CachedSteps)
	assert.Equal(t, 5.0, sample.ReportedBuildMS)
}

func TestCheckUsesAbsoluteAndRelativeMedianBudgets(t *testing.T) {
	baseline := &Result{Scenario: "badge-noop", MedianWallMS: 70}
	require.NoError(t, Check(Result{Scenario: "badge-noop", MedianWallMS: 80}, baseline, 100, 20))
	assert.ErrorContains(t, Check(Result{Scenario: "badge-noop", MedianWallMS: 101}, baseline, 100, 20), "exceeds 100.00ms")
	assert.ErrorContains(t, Check(Result{Scenario: "badge-noop", MedianWallMS: 85}, baseline, 100, 20), "baseline budget")
	assert.ErrorContains(t, Check(Result{Scenario: "other", MedianWallMS: 80}, baseline, 100, 20), "does not match")
}

func TestCheckBudgetEnforcesSampleVarianceWorkAndOverhead(t *testing.T) {
	zero, four := 0, 4
	result := Result{
		Scenario: "badge-noop", MedianWallMS: 84, MedianReportedBuildMS: 81,
		MedianAbsoluteDeviationMS: 3, MADPercent: 3.57, OrchestrationOverheadPercent: 3.70,
		Samples: []Sample{
			{WallMS: 80, ExecutedSteps: 0, CachedSteps: 4, StepsReported: true},
			{WallMS: 84, ExecutedSteps: 0, CachedSteps: 4, StepsReported: true},
			{WallMS: 87, ExecutedSteps: 0, CachedSteps: 4, StepsReported: true},
		},
	}
	budget := Budget{MinSamples: 3, MaxMADPercent: 15, MaxOrchestrationOverheadPercent: 5, ExpectedExecutedSteps: &zero, ExpectedCachedSteps: &four}
	require.NoError(t, CheckBudget(result, nil, budget))

	bad := result
	bad.MADPercent = 16
	assert.ErrorContains(t, CheckBudget(bad, nil, budget), "variance budget")
	bad = result
	bad.OrchestrationOverheadPercent = 6
	assert.ErrorContains(t, CheckBudget(bad, nil, budget), "orchestration overhead")
	bad = result
	bad.Samples[1].ExecutedSteps = 1
	assert.ErrorContains(t, CheckBudget(bad, nil, budget), "sample 2 executed 1 steps")
	bad = result
	bad.Samples = bad.Samples[:2]
	assert.ErrorContains(t, CheckBudget(bad, nil, budget), "at least 3")
}

func TestCheckBudgetRequiresEvidenceForRequestedChecks(t *testing.T) {
	zero := 0
	result := Result{Scenario: "badge-noop", MedianWallMS: 80, Samples: []Sample{{WallMS: 80}}}
	assert.ErrorContains(t, CheckBudget(result, nil, Budget{MaxRegressionPercent: 20}), "requires a baseline")
	assert.ErrorContains(t, CheckBudget(result, nil, Budget{MaxOrchestrationOverheadPercent: 5}), "no reported build duration")
	assert.ErrorContains(t, CheckBudget(result, nil, Budget{ExpectedExecutedSteps: &zero}), "did not report")
}

func TestRunMeasuresCommandSamples(t *testing.T) {
	command := []string{"sh", "-c", "printf 'build succeeded 2ms  1 ran • 3 cached\\n'"}
	if runtime.GOOS == "windows" {
		command = []string{"cmd.exe", "/c", "echo build succeeded 2ms  1 ran • 3 cached"}
	}
	result, err := Run(context.Background(), Config{Scenario: "test", Command: command, Warmups: 1, Samples: 3})
	require.NoError(t, err)
	assert.Len(t, result.Samples, 3)
	assert.Equal(t, 1, result.Samples[0].ExecutedSteps)
	assert.Equal(t, 3, result.Samples[0].CachedSteps)
	assert.True(t, result.Samples[0].StepsReported)
	assert.Positive(t, result.MedianWallMS)
	assert.Equal(t, 2.0, result.MedianReportedBuildMS)
	assert.GreaterOrEqual(t, result.MedianAbsoluteDeviationMS, 0.0)
	assert.GreaterOrEqual(t, result.MADPercent, 0.0)
}

func TestResultRoundTrip(t *testing.T) {
	path := t.TempDir() + string(os.PathSeparator) + "result.json"
	want := Result{Version: 1, Scenario: "badge", MedianWallMS: 70, Samples: []Sample{{WallMS: 70}}}
	require.NoError(t, Write(path, want))
	got, err := Read(path)
	require.NoError(t, err)
	assert.Equal(t, want, *got)
}
