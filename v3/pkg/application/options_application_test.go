package application

import (
	"reflect"
	"testing"
)

func TestOptions_getOptions(t *testing.T) {

	tests := []struct {
		name      string
		input     Options
		debugMode bool
		want      Options
	}{
		{
			name: "Override Icon in Production",
			input: Options{
				Icon: []byte("debug-icon"),
				ProductionOverrides: &Options{
					Icon: []byte("prod-icon"),
				},
			},
			want: Options{
				Icon: []byte("prod-icon"),
			},
		},
		{
			name: "Don't override Icon in debug",
			input: Options{
				Icon: []byte("debug-icon"),
				ProductionOverrides: &Options{
					Icon: []byte("prod-icon"),
				},
			},
			debugMode: true,
			want: Options{
				Icon: []byte("debug-icon"),
			},
		},
		{
			name: "Override Mac in Production",
			input: Options{
				Mac: MacOptions{
					ApplicationShouldTerminateAfterLastWindowClosed: false,
				},
				ProductionOverrides: &Options{
					Mac: MacOptions{
						ApplicationShouldTerminateAfterLastWindowClosed: true,
					},
				},
			},
			want: Options{
				Mac: MacOptions{
					ApplicationShouldTerminateAfterLastWindowClosed: true,
				},
			},
		},
		{
			name: "Don't override Mac in debug",
			input: Options{
				Mac: MacOptions{
					ApplicationShouldTerminateAfterLastWindowClosed: false,
				},
				ProductionOverrides: &Options{
					Mac: MacOptions{
						ApplicationShouldTerminateAfterLastWindowClosed: true,
					},
				},
			},
			debugMode: true,
			want: Options{
				Mac: MacOptions{
					ApplicationShouldTerminateAfterLastWindowClosed: false,
				},
			},
		},
		{
			name: "Override Flags in Production",
			input: Options{
				Flags: map[string]interface{}{
					"environment": "debug",
				},
				ProductionOverrides: &Options{
					Flags: map[string]interface{}{
						"environment": "prod",
					},
				},
			},
			want: Options{
				Flags: map[string]interface{}{
					"environment": "prod",
				},
			},
		},
		{
			name: "Do not override Flags in debug",
			input: Options{
				Flags: map[string]interface{}{
					"environment": "debug",
				},
				ProductionOverrides: &Options{
					Flags: map[string]interface{}{
						"environment": "prod",
					},
				},
			},
			debugMode: true,
			want: Options{
				Flags: map[string]interface{}{
					"environment": "debug",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.getOptions(tt.debugMode); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
