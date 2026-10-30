package changelog

import "testing"

var constructedProcessor *Processor

func BenchmarkNewProcessor(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		constructedProcessor = NewProcessor()
	}
	result, err := constructedProcessor.ProcessString("## Added\n- Add reusable changelog processors\n")
	if err != nil || !result.HasContent || !result.ValidationResult.Valid {
		b.Fatalf("constructed processor failed validation: %+v, %v", result, err)
	}
}
