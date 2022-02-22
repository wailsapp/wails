package docker

import (
	"context"
	dc "github.com/fsouza/go-dockerclient"
	"io"
	"strings"
)

func CrossCompile(projectPath string, outputfilename string, platform string, verbose bool, loggerWriter io.Writer) error {

	client, err := dc.NewClientFromEnv()
	if err != nil {
		return err
	}

	tag := strings.ReplaceAll(platform, "/", "-")

	outputStream := io.Discard
	if verbose {
		outputStream = loggerWriter
	}

	err = client.PullImage(dc.PullImageOptions{
		All:               false,
		Repository:        "wailsapp/cc",
		Tag:               tag,
		Platform:          "",
		Registry:          "",
		OutputStream:      outputStream,
		RawJSONStream:     false,
		InactivityTimeout: 0,
		Context:           context.Background(),
	}, dc.AuthConfiguration{})
	if err != nil {
		return err
	}

	projectDir := projectPath + ":/usr/src/myapp:rw"
	env := []string{"CGO_ENABLED=1"}
	if platform == "linux/arm64" {
		env = append(env, "CC_FOR_TARGET=aarch64-linux-gnu-gcc", "PKG_CONFIG_PATH=/usr/lib/aarch64-linux-gnu/pkgconfig")
	}

	container, err := client.CreateContainer(dc.CreateContainerOptions{
		//Name: "wails-" + tag + "-builder",
		Config: &dc.Config{
			Env:          env,
			Cmd:          []string{"go", "build", "-x", "-tags", "desktop,production", "-o", "build/bin/" + outputfilename},
			Image:        "wailsapp/cc:" + tag,
			WorkingDir:   "/usr/src/myapp",
			AttachStdin:  false,
			AttachStdout: verbose,
			AttachStderr: verbose,
			ArgsEscaped:  false,
		},
		HostConfig: &dc.HostConfig{
			Binds:      []string{projectDir},
			Privileged: true,
			AutoRemove: false,
		},
		NetworkingConfig: nil,
		Context:          nil,
	})
	if err != nil {
		return err
	}

	err = client.StartContainer(container.ID, &dc.HostConfig{
		Privileged: true,
		AutoRemove: false,
	})
	if err != nil {
		return err
	}

	statusCode, err := client.WaitContainer(container.ID)
	if statusCode != 0 {
		return err
	}

	/*
		client, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return fmt.Errorf("error from Docker1: %v", err)
		}


		ctx := context.Background()
		pulloutput, err := client.ImagePull(ctx, image, types.ImagePullOptions{})
		defer pulloutput.Close()
		output := loggerWriter
		if output == nil {
			output = os.Stdout
		}
		_, _ = io.Copy(output, pulloutput)
		//if err != nil {
		//	if verbose {
		//		println("here")
		//
		//
		//	}
		//	return fmt.Errorf("error from Docker2: %v", err)
		//}
		resp, err := client.ContainerCreate(ctx, &container.Config{
			AttachStdout: verbose,
			AttachStderr: verbose,
			Image:        image,
			Env:          []string{"CGO_ENABLED=1"},
			Cmd:          []string{"go", "build", "-tags", "desktop,production"},
		}, &container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: projectPath,
					Target: "/usr/src/myapp",
				},
			},
		}, nil, nil, "")
		if err != nil {
			return fmt.Errorf("error from Docker3: %v", err)
		}

		err = client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
		if err != nil {
			return fmt.Errorf("error from Docker4: %v", err)
		}

		statusCh, errCh := client.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				return fmt.Errorf("error from Docker5: %v", err)
			}
		case <-statusCh:
		}

		if verbose {
			out, err := client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
			if err != nil {
				return fmt.Errorf("error from Docke6r: %v", err)
			}

			_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
		}
	*/
	return nil
}
