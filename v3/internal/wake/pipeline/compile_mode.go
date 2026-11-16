package pipeline

import (
	"fmt"
	"slices"
)

// Resolve after primary selection and companion expansion: executables in the
// same application may use different backends and CGo dependencies.
func resolveCompileModes(plan *Plan, request Request) error {
	for key, node := range plan.Nodes {
		if node.Kind != CompileApplication {
			continue
		}
		spec := node.Spec.(CompileSpec)
		mode, err := compileMode(spec, request.CGOEnabled)
		if err != nil {
			return fmt.Errorf("%s: %w", node.Label, err)
		}
		if request.CC != "" && spec.Environment["CC"] == "" {
			spec.Environment = cloneStringMap(spec.Environment)
			if spec.Environment == nil {
				spec.Environment = map[string]string{}
			}
			spec.Environment["CC"] = request.CC
		}
		spec.CGOEnabled = mode
		if mode == "0" && !spec.Console && (spec.TargetOS == "linux" || spec.TargetOS == "darwin") {
			spec.Tags = appendUnique(spec.Tags, "purego")
		}
		node.Spec = spec
		plan.Nodes[key] = node
	}
	return nil
}

// Binding analysis runs on the host. Select its Go-only backend whenever a GUI
// uses pure Go so that package loading cannot introduce a host C compiler.
func resolveBindingsBackend(plan *Plan) {
	binding, ok := plan.Nodes["frontend:bindings"]
	if !ok {
		return
	}
	bindings := binding.Spec.(BindingsSpec)
	for _, node := range plan.Nodes {
		if spec, ok := node.Spec.(CompileSpec); ok && spec.Package == bindings.Package && !spec.Console && spec.CGOEnabled == "0" {
			bindings.CGOEnabled = "0"
			bindings.Tags = appendUnique(bindings.Tags, "purego")
			binding.Spec = bindings
			plan.Nodes[binding.Key] = binding
			return
		}
	}
}

func compileMode(spec CompileSpec, inherited string) (string, error) {
	mode := inherited
	if declared, ok := spec.Environment["CGO_ENABLED"]; ok {
		mode = declared
	}
	if mode == "" {
		mode = "1"
		if spec.Console || spec.TargetOS == "windows" || slices.Contains(spec.Tags, "purego") {
			mode = "0"
		}
	}
	if mode != "0" && mode != "1" {
		return "", fmt.Errorf("CGO_ENABLED must be 0 or 1, got %q", mode)
	}
	if mode == "0" && (spec.TargetOS == "android" || spec.TargetOS == "ios") {
		return "", fmt.Errorf("mobile targets require CGO_ENABLED=1")
	}
	return mode, nil
}

func developmentCompileEnvironment(environment map[string]string) map[string]string {
	result := map[string]string{}
	for _, name := range []string{"CGO_ENABLED", "CC", "CXX"} {
		if value, ok := environment[name]; ok {
			result[name] = value
		}
	}
	return result
}
