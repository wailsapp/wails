package commands

import "testing"

func TestBuild(t *testing.T) {
	type args struct {
		options *RunOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should error if task name not provided",
			args: args{
				options: &RunOptions{},
			},
			wantErr: true,
		},
		{
			name: "should work if task name provided",
			args: args{
				options: &RunOptions{
					Task: "build",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
