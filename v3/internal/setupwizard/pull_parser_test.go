package setupwizard

import (
	"strings"
	"testing"
)

func TestParseSize(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"10B", 10},
		{"10KB", 10 * 1024},
		{"10MB", 10 * 1024 * 1024},
		{"10GB", 10 * 1024 * 1024 * 1024},
		{"1.5MB", 1.5 * 1024 * 1024},
		{"100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parseSize(tt.input)
			if got != tt.expected {
				t.Errorf("parseSize(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestPullParser(t *testing.T) {
	dockerOutput := `latest: Pulling from wailsapp/wails-cross
1074353eec0d: Already exists
665f62578bce: Already exists
61d7d5b4f231: Pulling fs layer
202b93a508bb: Pulling fs layer
61d7d5b4f231: Downloading  1.5MB/10MB
202b93a508bb: Downloading  2MB/20MB
61d7d5b4f231: Downloading  5MB/10MB
202b93a508bb: Downloading  10MB/20MB
61d7d5b4f231: Verifying Checksum
61d7d5b4f231: Download complete
61d7d5b4f231: Extracting  1MB/10MB
61d7d5b4f231: Extracting  10MB/10MB
61d7d5b4f231: Pull complete
202b93a508bb: Verifying Checksum
202b93a508bb: Download complete
202b93a508bb: Extracting  5MB/20MB
202b93a508bb: Extracting  20MB/20MB
202b93a508bb: Pull complete
Digest: sha256:abc123
Status: Downloaded newer image`

	parser := newPullParser()
	lines := strings.Split(dockerOutput, "\n")

	type checkpoint struct {
		lineContains string
		minProgress  int
		maxProgress  int
		stage        string
	}

	checkpoints := []checkpoint{
		{"Pulling from", 0, 0, "Connecting"},
		{"Pulling fs layer", 0, 0, "Downloading"},
		{"Downloading  1.5MB/10MB", 1, 20, "Downloading"},
		{"Downloading  5MB/10MB", 20, 60, "Downloading"},
		{"Download complete", 30, 70, "Extracting"},
		{"Pull complete", 30, 100, "Extracting"},
	}

	checkIdx := 0
	for _, line := range lines {
		progress := parser.ParseLine(line)

		if checkIdx < len(checkpoints) && strings.Contains(line, checkpoints[checkIdx].lineContains) {
			cp := checkpoints[checkIdx]
			if progress.Progress < cp.minProgress || progress.Progress > cp.maxProgress {
				t.Errorf("After %q: progress=%d, want between %d and %d",
					cp.lineContains, progress.Progress, cp.minProgress, cp.maxProgress)
			}
			if progress.Stage != cp.stage {
				t.Errorf("After %q: stage=%q, want %q", cp.lineContains, progress.Stage, cp.stage)
			}
			checkIdx++
		}
	}

	finalProgress := parser.ParseLine("")
	if finalProgress.Progress != 100 {
		t.Errorf("Final progress = %d, want 100", finalProgress.Progress)
	}
}

func TestPullParserRealOutput(t *testing.T) {
	realOutput := `latest: Pulling from wailsapp/wails-cross
1074353eec0d: Already exists
665f62578bce: Already exists
61d7d5b4f231: Pulling fs layer
202b93a508bb: Pulling fs layer
604349a0d76e: Pulling fs layer
61d7d5b4f231: Downloading  1MB/10MB
202b93a508bb: Downloading  2MB/20MB
604349a0d76e: Downloading  1MB/15MB
61d7d5b4f231: Downloading  5MB/10MB
202b93a508bb: Downloading  10MB/20MB
604349a0d76e: Downloading  8MB/15MB
61d7d5b4f231: Verifying Checksum
61d7d5b4f231: Download complete
61d7d5b4f231: Pull complete
202b93a508bb: Verifying Checksum
202b93a508bb: Download complete
202b93a508bb: Pull complete
604349a0d76e: Verifying Checksum
604349a0d76e: Download complete
604349a0d76e: Pull complete
Digest: sha256:abc123
Status: Downloaded newer image`

	parser := newPullParser()
	lines := strings.Split(realOutput, "\n")

	var lastProgress PullProgress
	for _, line := range lines {
		lastProgress = parser.ParseLine(line)
		t.Logf("Line: %q -> Stage: %s, Progress: %d%%", line, lastProgress.Stage, lastProgress.Progress)
	}

	if lastProgress.Stage == "Connecting" {
		t.Errorf("Stage should not be 'Connecting' after parsing layers, got: %s", lastProgress.Stage)
	}
	if lastProgress.Progress != 100 {
		t.Errorf("Final progress should be 100%%, got: %d%%", lastProgress.Progress)
	}
}

func TestPullParserWithANSI(t *testing.T) {
	lines := []string{
		"latest: Pulling from wailsapp/wails-cross",
		"[1A[2K1074353eec0d: Pulling fs layer [1B",
		"[1A[2K665f62578bce: Pulling fs layer [1B",
		"[22A[2K665f62578bce: Downloading [==>                                                ]  16.38kB/296.1kB[22B",
		"[21A[2K5c445a0e108b: Downloading [>                                                  ]  538.1kB/60.15MB[21B",
		"[23A[2K1074353eec0d: Downloading [=================================================> ]  3.811MB/3.86MB[23B",
		"[23A[2K1074353eec0d: Verifying Checksum [23B",
		"[23A[2K1074353eec0d: Download complete [23B",
		"[23A[2K1074353eec0d: Extracting [>                                                  ]  65.54kB/3.86MB[23B",
		"[23A[2K1074353eec0d: Pull complete [23B",
		"[22A[2K665f62578bce: Pull complete [22B",
	}

	parser := newPullParser()
	var lastProgress PullProgress
	for _, line := range lines {
		lastProgress = parser.ParseLine(line)
		t.Logf("Line: %.60q -> Stage: %s, Progress: %d%%", line, lastProgress.Stage, lastProgress.Progress)
	}

	if lastProgress.Stage == "Connecting" {
		t.Errorf("Stage should not be 'Connecting', got: %s", lastProgress.Stage)
	}

	if len(parser.layerSizes) == 0 {
		t.Errorf("Should have parsed layer sizes, got none")
	}

	if lastProgress.Progress == 0 {
		t.Errorf("Progress should not be 0 after Pull complete")
	}
}

func TestPullParserNoSizeInfo(t *testing.T) {
	output := `latest: Pulling from wailsapp/wails-cross
layer1: Pulling fs layer
layer2: Pulling fs layer
layer3: Pulling fs layer
layer4: Pulling fs layer
layer1: Verifying Checksum
layer1: Download complete
layer1: Pull complete
layer2: Verifying Checksum
layer2: Download complete
layer2: Pull complete
layer3: Verifying Checksum
layer3: Download complete
layer3: Pull complete
layer4: Verifying Checksum
layer4: Download complete
layer4: Pull complete`

	parser := newPullParser()
	lines := strings.Split(output, "\n")

	progressHistory := []int{}
	for _, line := range lines {
		progress := parser.ParseLine(line)
		progressHistory = append(progressHistory, progress.Progress)
		t.Logf("Line: %q -> Stage: %s, Progress: %d%%", line, progress.Stage, progress.Progress)
	}

	finalProgress := progressHistory[len(progressHistory)-1]
	if finalProgress != 100 {
		t.Errorf("Final progress should be 100%%, got: %d%%", finalProgress)
	}

	for i := 1; i < len(progressHistory); i++ {
		if progressHistory[i] < progressHistory[i-1] {
			t.Errorf("Progress should not decrease: %d -> %d", progressHistory[i-1], progressHistory[i])
		}
	}
}
