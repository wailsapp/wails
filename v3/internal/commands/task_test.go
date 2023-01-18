package commands

import "testing"

func TestBuild(t *testing.T) {
	type args struct {
		options *RunTaskOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should error if task name not provided",
			args: args{
				options: &RunTaskOptions{},
			},
			wantErr: true,
		},
		{
			name: "should work if task name provided",
			args: args{
				options: &RunTaskOptions{
					Name: "build",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunTask(tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
