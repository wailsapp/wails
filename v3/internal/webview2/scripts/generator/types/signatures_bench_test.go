package types

import (
	"fmt"
	"testing"
)

func BenchmarkMethodSignatures(b *testing.B) {
	for _, count := range []int{0, 4, 16} {
		b.Run(fmt.Sprintf("params_%d", count), func(b *testing.B) {
			var method InterfaceMethod
			for index := range count {
				method.inputParams = append(method.inputParams, &Param{
					Name: fmt.Sprintf("input%d", index), Type: "int", GoType: "int",
				})
				method.outputParams = append(method.outputParams, &Param{
					Name: fmt.Sprintf("output%d", index), Type: "int", GoType: "int",
					Direction: &Direction{Dir: "out"},
				})
			}
			b.ReportAllocs()
			for b.Loop() {
				method.processInputParams()
				method.processOutputParams()
			}
			if method.GoReturnTypes == "" {
				b.Fatal("missing return signature")
			}
		})
	}
}
