package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/wailsapp/wails/v3/internal/flags"
	"github.com/wailsapp/wails/v3/internal/report"
	"github.com/wailsapp/wails/v3/internal/report/pulse"
	"github.com/wailsapp/wails/v3/internal/term"
	"github.com/wailsapp/wails/v3/internal/version"
	"github.com/wailsapp/wails/v3/internal/wake"
	"github.com/wailsapp/wails/v3/internal/wake/manifest"
	"github.com/wailsapp/wails/v3/internal/wake/packagetemplate"
	"github.com/wailsapp/wails/v3/internal/wake/pipeline"
	"gopkg.in/yaml.v3"
)

type manifestRunOptions struct {
	Context                             context.Context
	Verb, Profile, TargetOS, TargetArch string
	Loaded                              *manifest.Loaded
	Targets                             []pipeline.Target
	Formats                             []string
	Environment                         []string
	Development, Force, Obfuscated      bool
	Tags                                []string
	GarbleArgs                          []string
}

type manifestPipelineRun struct {
	Plan    pipeline.Plan
	Results map[pipeline.NodeKey]pipeline.Result
}

type manifestPlanOutput struct {
	SchemaVersion int                     `json:"schema_version"`
	Request       manifestPlanRequest     `json:"request"`
	Operations    []manifestPlanOperation `json:"operations"`
	Artifacts     []manifestPlanArtifact  `json:"artifacts"`
}

type manifestPlanRequest struct {
	Command   string                 `json:"command"`
	Profile   string                 `json:"profile,omitempty"`
	Targets   []string               `json:"targets"`
	Formats   []string               `json:"formats,omitempty"`
	Compilers []manifestPlanCompiler `json:"compilers"`
}

type manifestPlanCompiler struct {
	Target     string   `json:"target"`
	Toolchain  string   `json:"toolchain"`
	Tags       []string `json:"tags"`
	Obfuscated bool     `json:"obfuscated"`
	GarbleArgs []string `json:"garble_args,omitempty"`
}

type manifestPlanOperation struct {
	ID        string               `json:"id"`
	Kind      string               `json:"kind"`
	Stage     string               `json:"stage"`
	Scope     string               `json:"scope"`
	Decision  string               `json:"decision"`
	Cache     string               `json:"cache"`
	DependsOn []string             `json:"depends_on,omitempty"`
	Output    string               `json:"output,omitempty"`
	Origins   []manifestPlanOrigin `json:"origins,omitempty"`
	Inputs    []manifestPlanInput  `json:"inputs,omitempty"`
}

type manifestPlanInput struct {
	Label    string `json:"label"`
	Snapshot string `json:"snapshot"`
}

type manifestPlanOrigin struct {
	Field  string `json:"field"`
	Source string `json:"source"`
}

type manifestPlanArtifact struct {
	Target    string `json:"target"`
	Format    string `json:"format"`
	Kind      string `json:"kind"`
	Path      string `json:"path"`
	Signed    bool   `json:"signed"`
	Notarized bool   `json:"notarized"`
}

func runManifestPipeline(options manifestRunOptions) error {
	_, err := runManifestPipelineResult(options)
	return err
}

func runManifestPipelineResult(options manifestRunOptions) (manifestPipelineRun, error) {
	started := time.Now()
	root, loaded, plan, err := resolveManifestPlan(options)
	if err != nil {
		return manifestPipelineRun{}, err
	}
	reporter := pulse.New(os.Stdout, report.Normal)
	term.Header(strings.ToUpper(options.Verb[:1]) + options.Verb[1:])
	report.SetActive(reporter)
	defer report.SetActive(nil)
	reporter.BuildStart(options.Verb, plan.Target, len(plan.Nodes))
	ctx := options.Context
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()
	}
	results, err := (pipeline.Executor{Handler: &manifestHandler{root: root, config: loaded.Config, environment: options.Environment}}).Execute(ctx, plan, pipeline.ExecuteOptions{Root: root, Force: options.Force, Reporter: reporter})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			reporter.BuildCanceled(time.Since(started))
			return manifestPipelineRun{Plan: plan, Results: results}, err
		}
		reporter.BuildEnd(time.Since(started), false)
		return manifestPipelineRun{Plan: plan, Results: results}, wake.MarkReported(err)
	}
	resultKeys := make([]string, 0, len(results))
	for key := range results {
		resultKeys = append(resultKeys, string(key))
	}
	sort.Strings(resultKeys)
	for _, rawKey := range resultKeys {
		key := pipeline.NodeKey(rawKey)
		node := plan.Nodes[key]
		result := results[key]
		if result.Output != "" && !node.Artifact.Empty() {
			reporter.Artifact(report.Artifact{Path: filepath.Join(root, filepath.FromSlash(result.Output)), Kind: node.Artifact.DisplayKind()})
		}
	}
	reporter.BuildEnd(time.Since(started), true)
	return manifestPipelineRun{Plan: plan, Results: results}, nil
}

func resolveManifestPlan(options manifestRunOptions) (string, *manifest.Loaded, pipeline.Plan, error) {
	return resolveManifestPlanWithOperations(options, manifestPlanOperations{
		getwd: os.Getwd,
		load:  manifest.Load,
		plan: func(config manifest.Config, request pipeline.Request) (pipeline.Plan, error) {
			return pipeline.PlanBuildForHost(config, request, pipeline.CurrentHostCapabilities(manifestCredentialNames(config)...))
		},
	})
}

func manifestCredentialNames(config manifest.Config) []string {
	platforms := []manifest.SigningPlatform{config.Signing.Windows, config.Signing.Darwin, config.Signing.Linux, config.Signing.IOS, config.Signing.Android}
	result := make([]string, 0, len(platforms)*2)
	for _, platform := range platforms {
		if platform.Credential != "" {
			result = append(result, platform.Credential)
		}
		if platform.NotarizationCredential != "" {
			result = append(result, platform.NotarizationCredential)
		}
	}
	return result
}

type manifestPlanOperations struct {
	getwd func() (string, error)
	load  func(string, string) (*manifest.Loaded, error)
	plan  func(manifest.Config, pipeline.Request) (pipeline.Plan, error)
}

func resolveManifestPlanWithOperations(options manifestRunOptions, operations manifestPlanOperations) (string, *manifest.Loaded, pipeline.Plan, error) {
	root, err := operations.getwd()
	if err != nil {
		return "", nil, pipeline.Plan{}, err
	}
	loaded := options.Loaded
	if loaded == nil {
		loaded, err = operations.load(root, options.Profile)
		if err != nil {
			return "", nil, pipeline.Plan{}, err
		}
	}
	if loaded.Config.Root != "" {
		root = loaded.Config.Root
	}
	plan, err := operations.plan(loaded.Config, pipeline.Request{Verb: options.Verb, TargetOS: options.TargetOS, TargetArch: options.TargetArch, Targets: options.Targets, Formats: options.Formats, Development: options.Development, ExtraTags: options.Tags, GarbleArgs: options.GarbleArgs, Obfuscated: options.Obfuscated})
	if err != nil {
		return "", nil, pipeline.Plan{}, err
	}
	return root, loaded, plan, nil
}

func printManifestPlan(options manifestRunOptions, asJSON bool) error {
	root, loaded, plan, err := resolveManifestPlan(options)
	if err != nil {
		return err
	}
	ctx := options.Context
	if ctx == nil {
		ctx = context.Background()
	}
	inspection, err := (pipeline.Executor{Handler: &manifestHandler{root: root, config: loaded.Config, environment: options.Environment}}).Inspect(ctx, plan, root)
	if err != nil {
		return err
	}
	if asJSON {
		return json.NewEncoder(os.Stdout).Encode(planOutput(plan, inspection))
	}
	document := planOutput(plan, inspection)
	fmt.Printf("Plan: %s\n", chooseString(document.Request.Profile != "", document.Request.Profile, "anonymous"))
	fmt.Printf("Manifest: %s\n", filepath.ToSlash(filepath.Join(".", manifest.Filename)))
	fmt.Printf("Targets: %s\n", strings.Join(document.Request.Targets, " · "))
	if len(document.Request.Formats) > 0 {
		fmt.Printf("Formats: %s\n", strings.Join(document.Request.Formats, " · "))
	}
	for _, compiler := range document.Request.Compilers {
		fmt.Printf("Compiler %s: toolchain=%s; tags=%s; obfuscated=%s", compiler.Target, compiler.Toolchain, strings.Join(compiler.Tags, ", "), chooseString(compiler.Obfuscated, "yes", "no"))
		if len(compiler.GarbleArgs) > 0 {
			fmt.Printf("; garble args=%s", strings.Join(compiler.GarbleArgs, " "))
		}
		fmt.Println()
	}
	fmt.Println("STAGE\tSCOPE\tDECISION\tOUTPUT")
	for _, operation := range document.Operations {
		output := operation.Output
		if output == "" {
			output = operation.ID
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", operation.Stage, operation.Scope, strings.ToUpper(operation.Decision), output)
	}
	fmt.Println("No files will be changed because --plan was used.")
	return nil
}

func planOutput(plan pipeline.Plan, inspections ...pipeline.Inspection) manifestPlanOutput {
	targets := make([]string, len(plan.Intent.Targets))
	var formats []string
	for index, intent := range plan.Intent.Targets {
		targets[index] = intent.Target.OS + "/" + intent.Target.Arch
		formats = appendUniqueStrings(formats, intent.Formats...)
	}
	sort.Strings(formats)
	result := manifestPlanOutput{SchemaVersion: 1, Request: manifestPlanRequest{Command: plan.Intent.Command, Profile: plan.Intent.Profile, Targets: targets, Formats: formats, Compilers: resolvedPlanCompilers(plan)}}
	keys := make([]string, 0, len(plan.Nodes))
	for key := range plan.Nodes {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)
	for _, key := range keys {
		node := plan.Nodes[pipeline.NodeKey(key)]
		dependencies := make([]string, len(node.Dependencies))
		for index, dependency := range node.Dependencies {
			dependencies[index] = string(dependency)
		}
		decision := "run"
		if node.Cache != pipeline.CacheNever {
			decision = "cache-check"
		}
		var inputs []manifestPlanInput
		if len(inspections) != 0 {
			if inspected, ok := inspections[0].Operations[node.Key]; ok {
				decision = inspected.Decision
				inputs = make([]manifestPlanInput, len(inspected.Inputs))
				for index, input := range inspected.Inputs {
					inputs[index] = manifestPlanInput{Label: input.Label, Snapshot: input.Digest}
				}
			}
		}
		origins := make([]manifestPlanOrigin, len(node.Origins))
		for index, origin := range node.Origins {
			origins[index] = manifestPlanOrigin{Field: origin.Field, Source: formatPlanOrigin(origin.Origin)}
		}
		result.Operations = append(result.Operations, manifestPlanOperation{ID: key, Kind: string(node.Kind), Stage: planStage(node.Kind), Scope: string(node.Scope), Decision: decision, Cache: string(node.Cache), DependsOn: dependencies, Output: node.Output, Origins: origins, Inputs: inputs})
	}
	for _, root := range plan.Artifacts {
		if artifact, ok := planArtifact(plan.Nodes[root]); ok {
			result.Artifacts = append(result.Artifacts, artifact)
		}
	}
	return result
}

func formatPlanOrigin(origin manifest.Origin) string {
	if origin.Kind == manifest.OriginDefault {
		return "default"
	}
	if origin.Range.Filename == "" {
		return string(origin.Kind)
	}
	return fmt.Sprintf("%s:%d:%d", filepath.ToSlash(origin.Range.Filename), origin.Range.StartLine, origin.Range.StartColumn)
}

func planArtifact(node pipeline.Node) (manifestPlanArtifact, bool) {
	if node.Output == "" {
		return manifestPlanArtifact{}, false
	}
	switch spec := node.Spec.(type) {
	case pipeline.PublishSpec:
		identity := node.Artifact
		format := identity.Format
		if format == "" {
			format = string(identity.Kind)
		}
		return manifestPlanArtifact{Target: identity.Target.OS + "/" + identity.Target.Arch, Format: format, Kind: identity.DisplayKind(), Path: spec.Destination, Signed: identity.Signed, Notarized: identity.Notarized}, true
	case pipeline.SignSpec:
		return manifestPlanArtifact{Target: spec.TargetOS + "/" + spec.TargetArch, Format: node.Artifact.Format, Kind: node.Artifact.DisplayKind(), Path: node.Output, Signed: node.Artifact.Signed, Notarized: node.Artifact.Notarized}, true
	case pipeline.PackageSpec:
		return manifestPlanArtifact{Target: spec.TargetOS + "/" + spec.TargetArch, Format: spec.Format, Kind: node.Artifact.DisplayKind(), Path: node.Output}, true
	case pipeline.CompileSpec:
		return manifestPlanArtifact{Target: spec.TargetOS + "/" + spec.TargetArch, Format: "binary", Kind: node.Artifact.DisplayKind(), Path: node.Output}, true
	default:
		return manifestPlanArtifact{}, false
	}
}

func resolvedPlanCompilers(plan pipeline.Plan) []manifestPlanCompiler {
	result := make([]manifestPlanCompiler, 0, len(plan.Nodes))
	for _, node := range plan.Nodes {
		if node.Kind != pipeline.CompileApplication {
			continue
		}
		spec, ok := node.Spec.(pipeline.CompileSpec)
		if !ok {
			continue
		}
		result = append(result, manifestPlanCompiler{Target: spec.TargetOS + "/" + spec.TargetArch, Toolchain: spec.Toolchain, Tags: append([]string(nil), spec.Tags...), Obfuscated: spec.Obfuscated, GarbleArgs: append([]string(nil), spec.GarbleArgs...)})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Target < result[j].Target })
	return result
}

func planStage(kind pipeline.NodeKind) string {
	switch kind {
	case pipeline.InstallFrontendDependencies:
		return "prepare"
	case pipeline.GenerateBindings, pipeline.GeneratePlatformAssets:
		return "generate"
	case pipeline.BuildFrontend:
		return "frontend"
	case pipeline.CompileApplication:
		return "compile"
	case pipeline.MergeUniversalBinaries:
		return "assemble"
	case pipeline.AssembleApplication:
		return "assemble"
	case pipeline.PackageArtifact:
		return "package"
	case pipeline.SignArtifact:
		return "sign"
	case pipeline.PublishArtifact:
		return "publish"
	case pipeline.CollectArtifacts:
		return "collect"
	default:
		return "unknown"
	}
}

type manifestHandler struct {
	root        string
	config      manifest.Config
	environment []string
}

// manifestExecutable is injectable so package adapters can be tested without
// invoking the test binary as a second CLI process.
var manifestExecutable = os.Executable
var manifestHostOS = runtime.GOOS
var manifestLipo = ToolLipo
var manifestToolPackage = toolPackage
var manifestSign = Sign
var manifestDockerImageIdentity = func(ctx context.Context) (string, error) {
	output, err := exec.CommandContext(ctx, "docker", "image", "inspect", "--format", "{{.Id}}", "wails-cross").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (h *manifestHandler) Identity(ctx context.Context, node pipeline.Node) (string, error) {
	var tools []string
	var externalIdentity string
	if executable, err := os.Executable(); err == nil {
		tools = append(tools, executable)
	}
	switch node.Kind {
	case pipeline.InstallFrontendDependencies, pipeline.BuildFrontend:
		manager := h.config.Frontend.PackageManager
		path, err := exec.LookPath(manager)
		if err != nil {
			return "", fmt.Errorf("%s not found: %w", manager, err)
		}
		tools = append(tools, path)
		if manager == "npm" {
			if node, err := exec.LookPath("node"); err == nil {
				tools = append(tools, node)
			}
		}
	case pipeline.GenerateBindings, pipeline.CompileApplication:
		names := []string{"go"}
		if node.Kind == pipeline.CompileApplication {
			spec, err := manifestNodeSpec[pipeline.CompileSpec](node)
			if err != nil {
				return "", err
			}
			if spec.Obfuscated {
				names = append(names, "garble")
			}
			switch spec.Toolchain {
			case "zig":
				names = append(names, "zig")
			case "docker":
				names = []string{"docker"}
				externalIdentity, err = manifestDockerImageIdentity(ctx)
				if err != nil {
					return "", fmt.Errorf("inspect wails-cross image identity: %w", err)
				}
			}
		}
		for _, name := range names {
			path, err := exec.LookPath(name)
			if err != nil {
				return "", fmt.Errorf("%s not found: %w", name, err)
			}
			tools = append(tools, path)
		}
	case pipeline.PackageArtifact:
		spec, err := manifestNodeSpec[pipeline.PackageSpec](node)
		if err != nil {
			return "", err
		}
		name := map[string]string{"nsis": "makensis", "msix": "MakeAppx.exe"}[spec.Format]
		if name != "" {
			path, err := exec.LookPath(name)
			if err != nil {
				return "", fmt.Errorf("%s packaging requires %s: %w", spec.Format, name, err)
			}
			tools = append(tools, path)
		}
	case pipeline.SignArtifact:
		// The signer is implemented by this CLI; its executable identity below
		// is the complete tool identity.
	default:
		// Built-in handlers are identified by the running CLI below.
	}
	identity, err := toolIdentity(tools...)
	if err != nil {
		return "", err
	}
	return "wails-" + version.String() + ":" + string(node.Kind) + "|" + identity + "|external:" + externalIdentity + "|env:" + relevantEnvironment(node, h.environment), nil
}

func relevantEnvironment(node pipeline.Node, overrides []string) string {
	keys := []string{"PATH"}
	switch node.Kind {
	case pipeline.InstallFrontendDependencies:
		keys = append(keys, "CI", "NODE_ENV", "NPM_CONFIG_USERCONFIG", "NPM_CONFIG_REGISTRY", "PNPM_HOME", "YARN_CACHE_FOLDER")
	case pipeline.BuildFrontend:
		keys = append(keys, "CI", "NODE_ENV", "NPM_CONFIG_USERCONFIG", "NPM_CONFIG_REGISTRY", "PNPM_HOME", "YARN_CACHE_FOLDER")
		if spec, ok := node.Spec.(pipeline.FrontendSpec); ok && !spec.Production {
			keys = append(keys, "FRONTEND_DEVSERVER_URL", wailsVitePort)
		}
	case pipeline.GenerateBindings, pipeline.CompileApplication:
		keys = append(keys, "GOENV", "GOFLAGS", "GOTOOLCHAIN", "GOWORK", "CGO_ENABLED", "CC", "CXX", "CGO_CFLAGS", "CGO_CPPFLAGS", "CGO_CXXFLAGS", "CGO_LDFLAGS")
		if spec, ok := node.Spec.(pipeline.CompileSpec); ok && !spec.Production {
			keys = append(keys, "FRONTEND_DEVSERVER_URL", wailsVitePort)
		}
	}
	values := make([]string, 0, len(keys))
	resolved := environmentValues(mergeEnvironment(os.Environ(), overrides))
	for _, key := range keys {
		if value, ok := resolved[key]; ok {
			values = append(values, key+"="+value)
		}
	}
	sort.Strings(values)
	return strings.Join(values, "\x00")
}

func environmentValues(environment []string) map[string]string {
	result := make(map[string]string, len(environment))
	for _, entry := range environment {
		if key, value, ok := strings.Cut(entry, "="); ok {
			result[key] = value
		}
	}
	return result
}

func (h *manifestHandler) Run(ctx context.Context, node pipeline.Node) (pipeline.RunResult, error) {
	switch node.Kind {
	case pipeline.InstallFrontendDependencies:
		spec, err := manifestNodeSpec[pipeline.InstallSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.install(ctx, spec)
	case pipeline.GenerateBindings:
		spec, err := manifestNodeSpec[pipeline.BindingsSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.bindings(spec, node.Output)
	case pipeline.BuildFrontend:
		spec, err := manifestNodeSpec[pipeline.FrontendSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.frontend(ctx, spec)
	case pipeline.CompileApplication:
		spec, err := manifestNodeSpec[pipeline.CompileSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.compile(ctx, spec)
	case pipeline.MergeUniversalBinaries:
		spec, err := manifestNodeSpec[pipeline.MergeSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.mergeUniversal(spec)
	case pipeline.GeneratePlatformAssets:
		spec, err := manifestNodeSpec[pipeline.AssetsSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.assets(spec)
	case pipeline.AssembleApplication:
		spec, err := manifestNodeSpec[pipeline.PackageSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.packageArtifact(ctx, spec)
	case pipeline.PackageArtifact:
		spec, err := manifestNodeSpec[pipeline.PackageSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.packageArtifact(ctx, spec)
	case pipeline.SignArtifact:
		spec, err := manifestNodeSpec[pipeline.SignSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.sign(ctx, spec)
	case pipeline.PublishArtifact:
		spec, err := manifestNodeSpec[pipeline.PublishSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.publish(spec)
	case pipeline.CollectArtifacts:
		spec, err := manifestNodeSpec[pipeline.CollectSpec](node)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		return h.collect(spec)
	default:
		return pipeline.RunResult{}, fmt.Errorf("wake: no handler for %s", node.Kind)
	}
}

func (h *manifestHandler) collect(spec pipeline.CollectSpec) (pipeline.RunResult, error) {
	for _, artifact := range spec.Artifacts {
		path, err := existingArtifactPathInsideProject(h.root, artifact.Path)
		if err != nil {
			return pipeline.RunResult{}, fmt.Errorf("collect artifact %s: %w", artifact.Key, err)
		}
		info, err := os.Stat(path)
		if err != nil {
			return pipeline.RunResult{}, fmt.Errorf("collect artifact %s: %w", artifact.Key, err)
		}
		if !info.Mode().IsRegular() && !info.IsDir() {
			return pipeline.RunResult{}, fmt.Errorf("collect artifact %s: unsupported file type", artifact.Key)
		}
	}
	if _, err := pipeline.WriteArtifactReceipt(h.root, spec.Receipt, spec.Artifacts); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{Output: spec.Receipt, Detail: fmt.Sprintf("collected %d artifact(s)", len(spec.Artifacts))}, nil
}

func (h *manifestHandler) publish(spec pipeline.PublishSpec) (pipeline.RunResult, error) {
	cleanSource := filepath.ToSlash(filepath.Clean(spec.Source))
	if cleanSource != ".wails" && !strings.HasPrefix(cleanSource, ".wails/") {
		return pipeline.RunResult{}, fmt.Errorf("publish source %q is not in the generated .wails workspace", spec.Source)
	}
	cleanDestination := filepath.ToSlash(filepath.Clean(spec.Destination))
	if cleanDestination == ".wails" || strings.HasPrefix(cleanDestination, ".wails/") {
		return pipeline.RunResult{}, fmt.Errorf("publish destination %q is inside the generated .wails workspace", spec.Destination)
	}
	source, err := existingArtifactPathInsideProject(h.root, spec.Source)
	if err != nil {
		return pipeline.RunResult{}, fmt.Errorf("publish source: %w", err)
	}
	destination, err := pathInsideProject(h.root, spec.Destination)
	if err != nil {
		return pipeline.RunResult{}, fmt.Errorf("publish destination: %w", err)
	}
	parent := filepath.Dir(destination)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	transaction, err := os.MkdirTemp(parent, ".wails-publish-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer os.RemoveAll(transaction)
	staged := filepath.Join(transaction, "artifact")
	if err := copyManifestPath(source, staged); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := replacePathTransactional(staged, destination); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{Output: spec.Destination}, nil
}

func existingArtifactPathInsideProject(root, value string) (string, error) {
	if value == "" || !strings.HasPrefix(filepath.ToSlash(filepath.Clean(value)), ".wails/") {
		return existingPathInsideProject(root, value)
	}
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(filepath.Join(root, filepath.FromSlash(value)))
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(resolvedRoot, resolved)
	if err != nil {
		return "", err
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q resolves outside the project", value)
	}
	return resolved, nil
}

func manifestNodeSpec[T pipeline.NodeSpec](node pipeline.Node) (T, error) {
	spec, ok := node.Spec.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("wake: %s node carries %T, expected %T", node.Kind, node.Spec, zero)
	}
	return spec, nil
}

func (h *manifestHandler) install(ctx context.Context, s pipeline.InstallSpec) (pipeline.RunResult, error) {
	directory, err := existingPathInsideProject(h.root, s.Directory)
	if err != nil {
		return pipeline.RunResult{}, fmt.Errorf("frontend directory: %w", err)
	}
	if len(s.Arguments) > 0 {
		output, err := runManifestCommand(ctx, directory, declaredEnvironment(h.environment, s.Environment), s.Arguments[0], s.Arguments[1:]...)
		return pipeline.RunResult{Detail: output}, err
	}
	args := []string{s.Command}
	if s.Manager == "npm" {
		args = append(args, "--no-audit", "--no-fund")
	}
	output, err := runPackageManager(ctx, directory, declaredEnvironment(h.environment, s.Environment), s.Manager, args...)
	return pipeline.RunResult{Detail: output}, err
}

func (h *manifestHandler) bindings(s pipeline.BindingsSpec, output string) (pipeline.RunResult, error) {
	// GenerateBindings sets the process-wide footer flag; restore the caller's
	// value after using the legacy command adapter.
	previous := DisableFooter
	defer func() { DisableFooter = previous }()
	flagsString := ""
	if len(s.Tags) > 0 {
		flagsString = "-tags " + strings.Join(s.Tags, ",")
	}
	outputDirectory, err := pathInsideProject(h.root, output)
	if err != nil {
		return pipeline.RunResult{}, fmt.Errorf("bindings output: %w", err)
	}
	err = GenerateBindings(&flags.GenerateBindingsOptions{BuildFlagsString: flagsString, OutputDirectory: outputDirectory, ModelsFilename: s.Config.ModelsFilename, IndexFilename: s.Config.IndexFilename, TimeType: s.Config.TimeType, TS: s.Config.TypeScript, UseInterfaces: s.Config.Interfaces, Clean: true, Silent: true, Obfuscated: s.Obfuscated}, nil)
	return pipeline.RunResult{}, err
}

func (h *manifestHandler) frontend(ctx context.Context, s pipeline.FrontendSpec) (pipeline.RunResult, error) {
	directory, err := existingPathInsideProject(h.root, s.Directory)
	if err != nil {
		return pipeline.RunResult{}, fmt.Errorf("frontend directory: %w", err)
	}
	if len(s.Arguments) > 0 {
		env := declaredEnvironment(h.environment, s.Environment)
		env = append(env, "PRODUCTION="+fmt.Sprint(s.Production))
		output, err := runManifestCommand(ctx, directory, env, s.Arguments[0], s.Arguments[1:]...)
		return pipeline.RunResult{Detail: output}, err
	}
	var args []string
	switch s.Manager {
	case "npm", "pnpm", "bun":
		args = []string{"run", s.Command}
	case "yarn":
		args = []string{s.Command}
	}
	if s.Manager == "npm" {
		args = append(args, "--silent")
	}
	env := declaredEnvironment(h.environment, s.Environment)
	env = append(env, "PRODUCTION="+fmt.Sprint(s.Production))
	output, err := runPackageManager(ctx, directory, env, s.Manager, args...)
	return pipeline.RunResult{Detail: output}, err
}

func (h *manifestHandler) compile(ctx context.Context, s pipeline.CompileSpec) (pipeline.RunResult, error) {
	if s.Toolchain == "docker" && s.TargetOS != "android" && s.TargetOS != "ios" {
		return h.compileDocker(ctx, s)
	}
	switch s.TargetOS {
	case "android":
		return h.compileAndroid(ctx, s)
	case "ios":
		return h.compileIOS(ctx, s)
	}
	output := filepath.Join(h.root, filepath.FromSlash(s.Output))
	tool, args := compileGoArgs(s)
	if s.TargetOS == "windows" {
		overlay, err := h.windowsResourceOverlay(s)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		args = append(args, "-overlay", overlay)
	}
	env := declaredEnvironment(h.environment, s.Environment)
	env = append(env, "GOOS="+s.TargetOS, "GOARCH="+s.TargetArch)
	if s.Toolchain == "zig" {
		triple, err := zigTarget(s.TargetOS, s.TargetArch)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		env = append(env, "CGO_ENABLED=1", "CC=zig cc -target "+triple, "CXX=zig c++ -target "+triple)
	}
	if s.TargetOS == "darwin" && s.MinimumVersion != "" {
		minimumFlag := "-mmacosx-version-min=" + s.MinimumVersion
		env = append(env, "MACOSX_DEPLOYMENT_TARGET="+s.MinimumVersion, "CGO_CFLAGS="+minimumFlag, "CGO_LDFLAGS="+minimumFlag)
	}
	var result string
	err := producePathTransactional(output, ".compile-output-*", func(staged string) error {
		var runErr error
		result, runErr = runManifestCommand(ctx, h.root, env, tool, append(args, "-o", staged, ".")...)
		return runErr
	})
	if err != nil && strings.Contains(result, "updates to go.mod needed") {
		return pipeline.RunResult{}, fmt.Errorf("Go module files need maintenance; run `go mod tidy` manually: %w", err)
	}
	return pipeline.RunResult{Detail: result}, err
}

func zigTarget(targetOS, targetArch string) (string, error) {
	targets := map[string]string{
		"windows/amd64": "x86_64-windows-gnu",
		"windows/arm64": "aarch64-windows-gnu",
		"linux/amd64":   "x86_64-linux-gnu",
		"linux/arm64":   "aarch64-linux-gnu",
	}
	if target := targets[targetOS+"/"+targetArch]; target != "" {
		return target, nil
	}
	return "", fmt.Errorf("zig toolchain does not support %s/%s", targetOS, targetArch)
}

func (h *manifestHandler) compileDocker(ctx context.Context, s pipeline.CompileSpec) (pipeline.RunResult, error) {
	if manifestHostOS == "windows" {
		return pipeline.RunResult{}, fmt.Errorf("Docker cross-compilation from Windows is not yet supported by the generated-workspace mount adapter")
	}
	output := filepath.Join(h.root, filepath.FromSlash(s.Output))
	tool, buildArgs := compileGoArgs(s)
	if s.TargetOS == "windows" {
		overlay, err := h.windowsResourceOverlay(s)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		buildArgs = append(buildArgs, "-overlay", overlay)
	}
	cc, cxx, err := dockerCompilers(s.TargetOS, s.TargetArch)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	arguments := []string{"run", "--rm", "-w", h.root, "-v", h.root + ":" + h.root}
	for _, localRoot := range s.LocalRoots {
		arguments = append(arguments, "-v", localRoot+":"+localRoot+":ro")
	}
	for _, value := range []string{"GOOS=" + s.TargetOS, "GOARCH=" + s.TargetArch, "CGO_ENABLED=1", "CC=" + cc, "CXX=" + cxx} {
		arguments = append(arguments, "-e", value)
	}
	environment := environmentValues(declaredEnvironment(h.environment, s.Environment))
	keys := make([]string, 0, len(environment))
	for key := range environment {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		arguments = append(arguments, "-e", key+"="+environment[key])
	}
	if s.TargetOS == "darwin" && s.MinimumVersion != "" {
		arguments = append(arguments, "-e", "MACOSX_DEPLOYMENT_TARGET="+s.MinimumVersion)
	}
	arguments = append(arguments, "--entrypoint", tool, "wails-cross")
	var result string
	err = producePathTransactional(output, ".docker-compile-output-*", func(staged string) error {
		commandArguments := append(append([]string(nil), arguments...), buildArgs...)
		commandArguments = append(commandArguments, "-o", staged, ".")
		var runErr error
		result, runErr = runManifestCommand(ctx, h.root, nil, "docker", commandArguments...)
		return runErr
	})
	if err != nil && strings.Contains(result, "updates to go.mod needed") {
		return pipeline.RunResult{}, fmt.Errorf("Go module files need maintenance; run `go mod tidy` manually: %w", err)
	}
	return pipeline.RunResult{Detail: result}, err
}

func dockerCompilers(targetOS, targetArch string) (string, string, error) {
	if targetArch != "amd64" && targetArch != "arm64" {
		return "", "", fmt.Errorf("Docker toolchain does not support %s/%s", targetOS, targetArch)
	}
	if targetOS == "linux" {
		arch := map[string]string{"amd64": "x86_64", "arm64": "aarch64"}[targetArch]
		return "zig cc -target " + arch + "-linux-gnu", "zig c++ -target " + arch + "-linux-gnu", nil
	}
	if targetOS == "windows" || targetOS == "darwin" {
		compiler := "zcc-" + targetOS + "-" + targetArch
		return compiler, compiler, nil
	}
	return "", "", fmt.Errorf("Docker toolchain does not support %s/%s", targetOS, targetArch)
}

func (h *manifestHandler) mergeUniversal(s pipeline.MergeSpec) (pipeline.RunResult, error) {
	inputs := make([]string, len(s.Inputs))
	for index, input := range s.Inputs {
		inputs[index] = filepath.Join(h.root, filepath.FromSlash(input))
	}
	output := filepath.Join(h.root, filepath.FromSlash(s.Output))
	previous := DisableFooter
	defer func() { DisableFooter = previous }()
	err := producePathTransactional(output, ".universal-output-*", func(staged string) error {
		return manifestLipo(&flags.Lipo{Inputs: inputs, Output: staged})
	})
	return pipeline.RunResult{}, err
}

func compileGoArgs(s pipeline.CompileSpec) (string, []string) {
	tool := "go"
	args := []string{"build"}
	if s.Obfuscated {
		tool = "garble"
		args = append(append([]string{}, s.GarbleArgs...), "build")
	}
	if len(s.Tags) > 0 {
		args = append(args, "-tags", strings.Join(s.Tags, ","))
	}
	if s.TrimPath {
		args = append(args, "-trimpath")
	}
	if !s.VCSInfo {
		args = append(args, "-buildvcs=false")
	}
	if len(s.CompilerFlags) > 0 {
		args = append(args, "-gcflags", strings.Join(s.CompilerFlags, " "))
	}
	ld := append([]string{}, s.LinkerFlags...)
	if s.Strip {
		ld = appendUniqueStrings(ld, "-w", "-s")
	}
	if len(ld) > 0 {
		args = append(args, "-ldflags", strings.Join(ld, " "))
	}
	return tool, args
}

func (h *manifestHandler) windowsResourceOverlay(s pipeline.CompileSpec) (string, error) {
	assets := filepath.Join(h.root, filepath.FromSlash(s.Assets), "windows")
	generated := filepath.Join(h.root, ".wails", "generated", "windows", s.TargetArch)
	syso := filepath.Join(generated, "wails_windows_"+s.TargetArch+".syso")
	virtual := filepath.Join(h.root, "wails_windows_"+s.TargetArch+".syso")
	data, err := json.Marshal(struct {
		Replace map[string]string `json:"Replace"`
	}{map[string]string{virtual: syso}})
	if err != nil {
		return "", err
	}
	overlay := filepath.Join(generated, "overlay.json")
	err = prepareGeneratedWorkspace(generated, ".windows-resource-stage-*", func(staged string) error {
		if err := GenerateSyso(&SysoOptions{Manifest: filepath.Join(assets, "wails.exe.manifest"), Info: filepath.Join(assets, "info.json"), Icon: filepath.Join(assets, "icon.ico"), Out: filepath.Join(staged, filepath.Base(syso)), Arch: s.TargetArch}); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(staged, "overlay.json"), data, 0o644)
	})
	return overlay, err
}

func (h *manifestHandler) compileAndroid(ctx context.Context, s pipeline.CompileSpec) (pipeline.RunResult, error) {
	sdk := firstNonempty(os.Getenv("ANDROID_HOME"), os.Getenv("ANDROID_SDK_ROOT"))
	ndk := os.Getenv("ANDROID_NDK_HOME")
	if ndk == "" && sdk != "" {
		matches, _ := filepath.Glob(filepath.Join(sdk, "ndk", "*"))
		sort.Strings(matches)
		if len(matches) > 0 {
			ndk = matches[len(matches)-1]
		}
	}
	if ndk == "" {
		return pipeline.RunResult{}, fmt.Errorf("Android NDK not found; set ANDROID_NDK_HOME or install an NDK under ANDROID_HOME")
	}
	hostTag := "linux-x86_64"
	switch manifestHostOS {
	case "darwin":
		hostTag = "darwin-x86_64"
	case "windows":
		hostTag = "windows-x86_64"
	}
	triple := "aarch64-linux-android"
	jni := "arm64-v8a"
	if s.TargetArch == "amd64" {
		triple = "x86_64-linux-android"
		jni = "x86_64"
	}
	if s.TargetArch != "arm64" && s.TargetArch != "amd64" {
		return pipeline.RunResult{}, fmt.Errorf("unsupported Android architecture %s", s.TargetArch)
	}
	minSDK := s.MinimumVersion
	if minSDK == "" {
		minSDK = "21"
	}
	toolchain := filepath.Join(ndk, "toolchains", "llvm", "prebuilt", hostTag, "bin")
	cc := filepath.Join(toolchain, triple+minSDK+"-clang")
	cxx := filepath.Join(toolchain, triple+minSDK+"-clang++")
	if manifestHostOS == "windows" {
		cc += ".cmd"
		cxx += ".cmd"
	}
	if _, err := os.Stat(cc); err != nil {
		return pipeline.RunResult{}, fmt.Errorf("Android compiler not found: %s", cc)
	}
	output := filepath.Join(h.root, s.Output)
	s.Tags = appendUniqueStrings(s.Tags, "android")
	tool, args := compileGoArgs(s)
	args = append(args, "-buildmode=c-shared", "-overlay", filepath.Join(h.root, s.Assets, "android", "overlay.json"))
	env := []string{"GOOS=android", "GOARCH=" + s.TargetArch, "CGO_ENABLED=1", "CC=" + cc, "CXX=" + cxx, "WAILS_ANDROID_JNI=" + jni}
	var result string
	err := producePathTransactional(output, ".android-compile-output-*", func(staged string) error {
		var runErr error
		result, runErr = runManifestCommand(ctx, h.root, env, tool, append(args, "-o", staged, ".")...)
		return runErr
	})
	return pipeline.RunResult{Detail: result}, err
}

func (h *manifestHandler) compileIOS(ctx context.Context, s pipeline.CompileSpec) (pipeline.RunResult, error) {
	if manifestHostOS != "darwin" {
		return pipeline.RunResult{}, fmt.Errorf("iOS builds require macOS and Xcode")
	}
	destination := s.Destination
	if destination == "" {
		destination = "simulator"
	}
	sdk := "iphonesimulator"
	suffix := "-simulator"
	minFlag := "-mios-simulator-version-min="
	if destination == "device" {
		sdk = "iphoneos"
		suffix = ""
		minFlag = "-miphoneos-version-min="
	}
	sdkPath, err := commandOutput(ctx, h.root, "xcrun", "--sdk", sdk, "--show-sdk-path")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	clang, err := commandOutput(ctx, h.root, "xcrun", "--sdk", sdk, "--find", "clang")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	minimum := s.MinimumVersion
	if minimum == "" {
		minimum = "15.0"
	}
	target := s.TargetArch + "-apple-ios" + minimum + suffix
	output := filepath.Join(h.root, s.Output)
	s.Tags = appendUniqueStrings(s.Tags, "ios")
	tool, args := compileGoArgs(s)
	args = append(args, "-buildmode=c-archive", "-overlay", filepath.Join(h.root, s.Assets, "ios", "xcode", "overlay.json"))
	flags := "-isysroot " + sdkPath + " -target " + target + " " + minFlag + minimum
	env := []string{"GOOS=ios", "GOARCH=" + s.TargetArch, "CGO_ENABLED=1", "CC=" + clang, "CGO_CFLAGS=" + flags, "CGO_LDFLAGS=-isysroot " + sdkPath + " -target " + target}
	var result string
	err = producePathTransactional(output, ".ios-compile-output-*", func(staged string) error {
		var runErr error
		result, runErr = runManifestCommand(ctx, h.root, env, tool, append(args, "-o", staged, ".")...)
		return runErr
	})
	return pipeline.RunResult{Detail: result}, err
}

func commandOutput(ctx context.Context, dir, name string, args ...string) (string, error) {
	output, err := runManifestCommand(ctx, dir, nil, name, args...)
	if err != nil {
		return "", fmt.Errorf("%s: %s: %w", name, output, err)
	}
	return strings.TrimSpace(output), nil
}
func firstNonempty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func (h *manifestHandler) assets(s pipeline.AssetsSpec) (pipeline.RunResult, error) {
	output := filepath.Join(h.root, filepath.FromSlash(s.Directory))
	parent := filepath.Dir(output)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	tmp, err := os.MkdirTemp(parent, ".assets-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer os.RemoveAll(tmp)
	options := &BuildAssetsOptions{Dir: tmp, Name: s.Project.Name, BinaryName: s.Project.BinaryName, ProductName: s.Project.ProductName, ProductDescription: s.Project.Description, ProductVersion: s.Project.Version, ProductCompany: s.Project.CompanyName, ProductIdentifier: s.Project.Identifier, ProductCopyright: s.Project.Copyright, ProductComments: s.Project.Comments, Silent: true, Typescript: h.config.Frontend.Bindings.TypeScript, UseInterfaces: h.config.Frontend.Bindings.Interfaces}
	if err := GenerateBuildAssets(options); err != nil {
		return pipeline.RunResult{}, err
	}
	if s.Project.Icon != "" {
		source, err := existingPathInsideProject(h.root, s.Project.Icon)
		if err != nil {
			return pipeline.RunResult{}, fmt.Errorf("project.icon: %w", err)
		}
		icon := filepath.Join(tmp, "appicon.png")
		if err := copyManifestPath(source, icon); err != nil {
			return pipeline.RunResult{}, err
		}
		if err := GenerateIcons(&IconsOptions{Input: icon, WindowsFilename: filepath.Join(tmp, "windows", "icon.ico"), MacFilename: filepath.Join(tmp, "darwin", "icons.icns")}); err != nil {
			return pipeline.RunResult{}, err
		}
	}
	associations, err := h.stageAssociationIcons(tmp, s.TargetOS, s.Associations)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	configPath, err := writeGeneratedConfig(tmp, s.Project, associations, s.Protocols)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	if err := UpdateBuildAssets(&UpdateBuildAssetsOptions{Dir: tmp, Name: s.Project.Name, BinaryName: s.Project.BinaryName, Config: configPath, Silent: true}); err != nil {
		return pipeline.RunResult{}, err
	}
	staged, err := os.MkdirTemp(parent, ".assets-stage-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer os.RemoveAll(staged)
	source := filepath.Join(tmp, s.TargetOS)
	if err := copyManifestPath(source, filepath.Join(staged, s.TargetOS)); err != nil {
		return pipeline.RunResult{}, err
	}
	_ = os.Remove(filepath.Join(staged, s.TargetOS, "Taskfile.yml"))
	if err := copyManifestPath(filepath.Join(tmp, "appicon.png"), filepath.Join(staged, "appicon.png")); err != nil {
		return pipeline.RunResult{}, err
	}
	resolvedConfig := filepath.Join(staged, "config.yml")
	if err := copyManifestPath(configPath, resolvedConfig); err != nil {
		return pipeline.RunResult{}, err
	}
	switch s.TargetOS {
	case "ios":
		xcode := filepath.Join(staged, "ios", "xcode")
		if err := IOSOverlayGen(&IOSOverlayGenOptions{Out: filepath.Join(xcode, "overlay.json"), Config: resolvedConfig}); err != nil {
			return pipeline.RunResult{}, err
		}
		if err := IOSXcodeGen(&IOSXcodeGenOptions{OutDir: xcode, Config: resolvedConfig}); err != nil {
			return pipeline.RunResult{}, err
		}
	case "android":
		if err := AndroidOverlayGen(&AndroidOverlayGenOptions{Out: filepath.Join(staged, "android", "overlay.json"), Config: resolvedConfig}); err != nil {
			return pipeline.RunResult{}, err
		}
	}
	if err := h.applyUserPlatformInputs(staged, s.TargetOS, s.Project.BinaryName); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := h.applyUserSigningInputs(staged, s.TargetOS); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := applyGeneratedTargetSettings(staged, s); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{}, replacePathTransactional(staged, output)
}

func (h *manifestHandler) stageAssociationIcons(generatedRoot, targetOS string, associations []manifest.Association) ([]manifest.Association, error) {
	result := append([]manifest.Association(nil), associations...)
	for index := range result {
		if result[index].Icon == "" {
			continue
		}
		source, err := existingPathInsideProject(h.root, result[index].Icon)
		if err != nil {
			return nil, fmt.Errorf("file_association[%q].icon: %w", result[index].Name, err)
		}
		stem := fmt.Sprintf("association-%03d", index)
		switch targetOS {
		case "windows":
			destination := filepath.Join(generatedRoot, "windows", stem+".ico")
			if strings.EqualFold(filepath.Ext(source), ".ico") {
				err = copyManifestPath(source, destination)
			} else {
				err = GenerateIcons(&IconsOptions{Input: source, WindowsFilename: destination})
			}
			result[index].Icon = stem
		case "darwin":
			destination := filepath.Join(generatedRoot, "darwin", stem+".icns")
			if strings.EqualFold(filepath.Ext(source), ".icns") {
				err = copyManifestPath(source, destination)
			} else {
				err = GenerateIcons(&IconsOptions{Input: source, MacFilename: destination})
			}
			result[index].Icon = stem
		case "ios":
			name := stem + filepath.Ext(source)
			destination := filepath.Join(generatedRoot, "ios", "xcode", "main", name)
			err = copyManifestPath(source, destination)
			result[index].Icon = name
		case "linux":
			name := stem + filepath.Ext(source)
			destination := filepath.Join(generatedRoot, "linux", name)
			err = copyManifestPath(source, destination)
			result[index].Icon = name
		case "android":
			name := stem + filepath.Ext(source)
			destination := filepath.Join(generatedRoot, "android", "app", "src", "main", "assets", name)
			err = copyManifestPath(source, destination)
			result[index].Icon = name
		default:
			return nil, fmt.Errorf("file_association[%q].icon: unsupported target %s", result[index].Name, targetOS)
		}
		if err != nil {
			return nil, fmt.Errorf("file_association[%q].icon: %w", result[index].Name, err)
		}
	}
	return result, nil
}

func (h *manifestHandler) applyUserSigningInputs(assetsRoot, targetOS string) error {
	signing := manifestSigningPlatform(h.config.Signing, targetOS)
	if signing == nil {
		return fmt.Errorf("unsupported signing assets %q", targetOS)
	}
	for _, input := range []struct {
		field, source, name string
	}{
		{"entitlements", signing.Entitlements, "entitlements"},
		{"provisioning_profile", signing.ProvisioningProfile, "provisioning-profile"},
	} {
		if input.source == "" {
			continue
		}
		source, err := existingPathInsideProject(h.root, input.source)
		if err != nil {
			return fmt.Errorf("signing.%s.%s: %w", targetOS, input.field, err)
		}
		destination := filepath.Join(assetsRoot, "signing", input.name+filepath.Ext(source))
		if err := copyManifestPath(source, destination); err != nil {
			return err
		}
	}
	return nil
}

func manifestSigningPlatform(signing manifest.Signing, targetOS string) *manifest.SigningPlatform {
	switch targetOS {
	case "windows":
		return &signing.Windows
	case "darwin":
		return &signing.Darwin
	case "linux":
		return &signing.Linux
	case "ios":
		return &signing.IOS
	case "android":
		return &signing.Android
	default:
		return nil
	}
}

func replacePathTransactional(staged, destination string) error {
	return replacePathTransactionalWithOperations(staged, destination, replacePathOperations{
		mkdirTemp: os.MkdirTemp,
		lstat:     os.Lstat,
		rename:    os.Rename,
		removeAll: os.RemoveAll,
	})
}

type replacePathOperations struct {
	mkdirTemp func(string, string) (string, error)
	lstat     func(string) (fs.FileInfo, error)
	rename    func(string, string) error
	removeAll func(string) error
}

func replacePathTransactionalWithOperations(staged, destination string, operations replacePathOperations) error {
	transaction, err := operations.mkdirTemp(filepath.Dir(destination), ".wails-replace-*")
	if err != nil {
		return err
	}
	defer operations.removeAll(transaction)
	previous := filepath.Join(transaction, "previous")
	hadPrevious := false
	if _, err := operations.lstat(destination); err == nil {
		if err := operations.rename(destination, previous); err != nil {
			return err
		}
		hadPrevious = true
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	if err := operations.rename(staged, destination); err != nil {
		if hadPrevious {
			if rollbackErr := operations.rename(previous, destination); rollbackErr != nil {
				return errors.Join(err, fmt.Errorf("restore previous generated workspace: %w", rollbackErr))
			}
		}
		return err
	}
	return nil
}

func (h *manifestHandler) applyUserPlatformInputs(assetsRoot, targetOS, binaryName string) error {
	platform := manifestPlatform(h.config.Targets, targetOS)
	if platform == nil {
		return fmt.Errorf("unsupported platform assets %q", targetOS)
	}
	type replacement struct {
		field, source, destination string
	}
	var replacements []replacement
	switch targetOS {
	case "windows":
		replacements = []replacement{{"icon", platform.Icon, "windows/icon.ico"}, {"manifest", platform.Manifest, "windows/wails.exe.manifest"}}
	case "darwin":
		replacements = []replacement{{"icon", platform.Icon, "darwin/icons.icns"}, {"assets_car", platform.AssetsCar, "darwin/Assets.car"}, {"info_plist", platform.InfoPlist, "darwin/Info.plist"}}
	case "linux":
		replacements = []replacement{{"icon", platform.Icon, "appicon.png"}, {"desktop_entry", platform.DesktopEntry, filepath.Join("linux", binaryName+".desktop")}}
	case "ios":
		replacements = []replacement{{"icon", platform.Icon, "ios/icon.png"}, {"assets_car", platform.AssetsCar, "ios/xcode/main/Assets.car"}, {"info_plist", platform.InfoPlist, "ios/xcode/main/Info.plist"}}
	case "android":
		replacements = []replacement{{"manifest", platform.Manifest, "android/app/src/main/AndroidManifest.xml"}}
	}
	for _, item := range replacements {
		if item.source == "" {
			continue
		}
		source, err := existingPathInsideProject(h.root, item.source)
		if err != nil {
			return fmt.Errorf("%s.%s: %w", targetOS, item.field, err)
		}
		if err := copyManifestPath(source, filepath.Join(assetsRoot, item.destination)); err != nil {
			return err
		}
	}
	return nil
}

func manifestPlatform(targets manifest.Targets, targetOS string) *manifest.Platform {
	switch targetOS {
	case "windows":
		return &targets.Windows
	case "darwin":
		return &targets.Darwin
	case "linux":
		return &targets.Linux
	case "ios":
		return &targets.IOS
	case "android":
		return &targets.Android
	default:
		return nil
	}
}

func applyGeneratedTargetSettings(output string, spec pipeline.AssetsSpec) error {
	if spec.Project.BuildNumber > 0 {
		value := strconv.Itoa(spec.Project.BuildNumber)
		for _, path := range []string{
			filepath.Join(output, "darwin", "Info.plist"),
			filepath.Join(output, "ios", "xcode", "main", "Info.plist"),
		} {
			if err := replacePlistString(path, "CFBundleVersion", value); err != nil {
				return err
			}
		}
		gradle := filepath.Join(output, "android", "app", "build.gradle")
		if err := replaceGeneratedPattern(gradle, regexp.MustCompile(`(?m)^(\s*versionCode\s+)\d+`), "${1}"+value); err != nil {
			return err
		}
		nsis := filepath.Join(output, "windows", "nsis", "project.nsi")
		if err := replaceGeneratedLiteral(nsis, `${INFO_PRODUCTVERSION}.0`, `${INFO_PRODUCTVERSION}.`+value); err != nil {
			return err
		}
	}
	if spec.MinimumVersion != "" {
		keys := map[string]string{"darwin": "LSMinimumSystemVersion", "ios": "MinimumOSVersion"}
		if key := keys[spec.TargetOS]; key != "" {
			path := filepath.Join(output, spec.TargetOS, "Info.plist")
			if spec.TargetOS == "ios" {
				path = filepath.Join(output, "ios", "xcode", "main", "Info.plist")
			}
			if err := replacePlistString(path, key, spec.MinimumVersion); err != nil {
				return err
			}
		}
	}
	return nil
}

func replacePlistString(path, key, value string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	marker := "<key>" + key + "</key>"
	start := bytes.Index(data, []byte(marker))
	if start < 0 {
		return nil
	}
	afterMarker := data[start+len(marker):]
	stringStart := bytes.Index(afterMarker, []byte("<string>"))
	selfClosing := bytes.Index(afterMarker, []byte("<string/>"))
	if stringStart >= 0 && (selfClosing < 0 || stringStart < selfClosing) {
		stringStart += start + len(marker) + len("<string>")
		stringEnd := bytes.Index(data[stringStart:], []byte("</string>"))
		if stringEnd < 0 {
			return nil
		}
		stringEnd += stringStart
		result := append(append(append([]byte(nil), data[:stringStart]...), value...), data[stringEnd:]...)
		return os.WriteFile(path, result, 0o644)
	}
	if selfClosing < 0 {
		return nil
	}
	selfClosing += start + len(marker)
	result := append(append(append([]byte(nil), data[:selfClosing]...), []byte("<string>")...), value...)
	result = append(result, []byte("</string>")...)
	result = append(result, data[selfClosing+len("<string/>"):]...)
	return os.WriteFile(path, result, 0o644)
}

func replaceGeneratedPattern(path string, pattern *regexp.Regexp, replacement string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return os.WriteFile(path, pattern.ReplaceAll(data, []byte(replacement)), 0o644)
}

func replaceGeneratedLiteral(path, old, replacement string) error {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes.ReplaceAll(data, []byte(old), []byte(replacement)), 0o644)
}

func (h *manifestHandler) copyPackageReplacement(s pipeline.PackageSpec, destination string) error {
	source, err := existingPathInsideProject(h.root, s.Config.Template)
	if err != nil {
		return fmt.Errorf("package template: %w", err)
	}
	return packagetemplate.Copy(source, destination)
}

func packageStringOption(options map[string]any, name, fallback string) string {
	if value, ok := options[name].(string); ok && value != "" {
		return value
	}
	return fallback
}

func packageIntOption(options map[string]any, name string, fallback int) int {
	switch value := options[name].(type) {
	case int:
		return value
	case int64:
		return int(value)
	case float64:
		return int(value)
	default:
		return fallback
	}
}

func resolveDMGFiles(root, value string) (string, error) {
	if value == "" {
		return "", nil
	}
	var resolved []string
	for _, item := range strings.Split(value, ",") {
		name, path, ok := strings.Cut(item, "=")
		if !ok || strings.TrimSpace(path) == "" {
			return "", fmt.Errorf("DMG files entry %q must be name=path", item)
		}
		path = strings.TrimSpace(path)
		path, err := existingPathInsideProject(root, path)
		if err != nil {
			return "", fmt.Errorf("DMG files path: %w", err)
		}
		resolved = append(resolved, strings.TrimSpace(name)+"="+path)
	}
	return strings.Join(resolved, ","), nil
}

func encodeDMGFiles(files map[string]string) string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	values := make([]string, 0, len(names))
	for _, name := range names {
		values = append(values, name+"="+files[name])
	}
	return strings.Join(values, ",")
}

func (h *manifestHandler) packageWorkspace(s pipeline.PackageSpec, elements ...string) string {
	profile := s.Profile
	if profile == "" {
		profile = "default"
	}
	target := strings.ReplaceAll(s.TargetOS+"/"+s.TargetArch, "/", "-")
	parts := []string{h.root, ".wails", "build", profile, target, "package", s.Format}
	return filepath.Join(append(parts, elements...)...)
}

func prepareGeneratedWorkspace(destination, pattern string, prepare func(string) error) error {
	parent := filepath.Dir(destination)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	staged, err := os.MkdirTemp(parent, pattern)
	if err != nil {
		return err
	}
	defer os.RemoveAll(staged)
	if err := prepare(staged); err != nil {
		return err
	}
	return replacePathTransactional(staged, destination)
}

func copyPathTransactional(source, destination, pattern string) error {
	return producePathTransactional(destination, pattern, func(staged string) error {
		return copyManifestPath(source, staged)
	})
}

func producePathTransactional(destination, pattern string, produce func(string) error) error {
	return produceNamedPathTransactional(destination, pattern, "artifact", produce)
}

func produceNamedPathTransactional(destination, pattern, name string, produce func(string) error) error {
	parent := filepath.Dir(destination)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	transaction, err := os.MkdirTemp(parent, pattern)
	if err != nil {
		return err
	}
	defer os.RemoveAll(transaction)
	staged := filepath.Join(transaction, name)
	if err := produce(staged); err != nil {
		return err
	}
	return replacePathTransactional(staged, destination)
}

func (h *manifestHandler) packageArtifact(ctx context.Context, s pipeline.PackageSpec) (pipeline.RunResult, error) {
	switch s.Format {
	case "app":
		if s.TargetOS == "ios" {
			return h.packageIOS(ctx, s, false)
		}
		return h.packageApp(s)
	case "dmg":
		options, err := h.dmgOptions(s)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		outputDirectory := filepath.Join(h.packageWorkspace(s), "output")
		if err := removeGenerated(outputDirectory); err != nil {
			return pipeline.RunResult{}, err
		}
		if err := os.MkdirAll(outputDirectory, 0o755); err != nil {
			return pipeline.RunResult{}, err
		}
		options.Out = outputDirectory
		if err := manifestToolPackage(options); err != nil {
			return pipeline.RunResult{}, err
		}
		return pipeline.RunResult{}, findAndMovePackage(outputDirectory, s.Project.BinaryName, "dmg", filepath.Join(h.root, s.Output))
	case "appimage":
		return h.packageAppImage(ctx, s)
	case "deb", "rpm", "archlinux":
		return h.packageLinux(s)
	case "nsis":
		return h.packageNSIS(ctx, s)
	case "msix":
		return h.packageMSIX(s)
	case "ipa":
		return h.packageIOS(ctx, s, true)
	case "apk", "aab":
		return h.packageAndroid(ctx, s)
	default:
		return pipeline.RunResult{}, fmt.Errorf("unsupported package format %q", s.Format)
	}
}

func (h *manifestHandler) dmgOptions(s pipeline.PackageSpec) (*flags.ToolPackage, error) {
	return h.dmgOptionsWithRead(s, os.ReadFile)
}

func (h *manifestHandler) dmgOptionsWithRead(s pipeline.PackageSpec, readFile func(string) ([]byte, error)) (*flags.ToolPackage, error) {
	workspace := h.packageWorkspace(s)
	if err := os.MkdirAll(filepath.Dir(workspace), 0o755); err != nil {
		return nil, err
	}
	staged, err := os.MkdirTemp(filepath.Dir(workspace), ".dmg-stage-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(staged)

	values := map[string]any{}
	if s.Config.Template != "" {
		path := filepath.Join(staged, "dmg.json")
		if err := h.copyPackageReplacement(s, path); err != nil {
			return nil, err
		}
		data, err := readFile(path)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &values); err != nil {
			return nil, fmt.Errorf("parse rendered DMG template: %w", err)
		}
	}
	for key, value := range map[string]string{
		"background": s.Config.Background, "volume_icon": s.Config.VolumeIcon, "file_icon": s.Config.FileIcon,
	} {
		if value != "" {
			values[key] = value
		}
	}
	if len(s.Config.Files) > 0 {
		values["files"] = encodeDMGFiles(s.Config.Files)
	}
	if s.Config.WindowWidth != 0 {
		values["window_width"] = s.Config.WindowWidth
	}
	if s.Config.WindowHeight != 0 {
		values["window_height"] = s.Config.WindowHeight
	}
	resolve := func(field, name, value string) (string, error) {
		if value == "" {
			return "", nil
		}
		source, err := existingPathInsideProject(h.root, value)
		if err != nil {
			return "", fmt.Errorf("DMG %s: %w", field, err)
		}
		relative := filepath.Join("resources", name+filepath.Ext(source))
		if err := copyManifestPath(source, filepath.Join(staged, relative)); err != nil {
			return "", fmt.Errorf("DMG %s: %w", field, err)
		}
		return relative, nil
	}
	background, err := resolve("background", "background", packageStringOption(values, "background", ""))
	if err != nil {
		return nil, err
	}
	volumeIcon, err := resolve("volume_icon", "volume-icon", packageStringOption(values, "volume_icon", ""))
	if err != nil {
		return nil, err
	}
	fileIcon, err := resolve("file_icon", "file-icon", packageStringOption(values, "file_icon", ""))
	if err != nil {
		return nil, err
	}
	var files []string
	if value := packageStringOption(values, "files", ""); value != "" {
		for index, item := range strings.Split(value, ",") {
			name, sourceValue, ok := strings.Cut(item, "=")
			if !ok || strings.TrimSpace(sourceValue) == "" {
				return nil, fmt.Errorf("DMG files entry %q must be name=path", item)
			}
			source, resolveErr := existingPathInsideProject(h.root, strings.TrimSpace(sourceValue))
			if resolveErr != nil {
				return nil, fmt.Errorf("DMG files path: %w", resolveErr)
			}
			relative := filepath.Join("resources", fmt.Sprintf("file-%03d%s", index, filepath.Ext(source)))
			if copyErr := copyManifestPath(source, filepath.Join(staged, relative)); copyErr != nil {
				return nil, fmt.Errorf("DMG files path: %w", copyErr)
			}
			files = append(files, strings.TrimSpace(name)+"="+relative)
		}
	}
	if err := replacePathTransactional(staged, workspace); err != nil {
		return nil, err
	}
	canonicalWorkspace, err := filepath.EvalSymlinks(workspace)
	if err != nil {
		return nil, err
	}
	ownedPath := func(relative string) string {
		if relative == "" {
			return ""
		}
		return filepath.Join(canonicalWorkspace, relative)
	}
	for index, file := range files {
		name, relative, _ := strings.Cut(file, "=")
		files[index] = name + "=" + ownedPath(relative)
	}
	return &flags.ToolPackage{
		Format: "dmg", ExecutableName: s.Project.BinaryName,
		Out:             filepath.Dir(filepath.Join(h.root, s.Output)),
		BackgroundImage: ownedPath(background),
		DmgVolumeIcon:   ownedPath(volumeIcon),
		DmgFileIcon:     ownedPath(fileIcon),
		DmgFiles:        strings.Join(files, ","),
		DmgWindowWidth:  packageIntOption(values, "window_width", 540),
		DmgWindowHeight: packageIntOption(values, "window_height", 380),
	}, nil
}

func (h *manifestHandler) packageApp(s pipeline.PackageSpec) (pipeline.RunResult, error) {
	output := filepath.Join(h.root, filepath.FromSlash(s.Output))
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	staged, err := os.MkdirTemp(filepath.Dir(output), ".darwin-app-stage-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer os.RemoveAll(staged)
	macos := filepath.Join(staged, "Contents", "MacOS")
	resources := filepath.Join(staged, "Contents", "Resources")
	if err := os.MkdirAll(macos, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := os.MkdirAll(resources, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := copyManifestPath(filepath.Join(h.root, s.Binary), filepath.Join(macos, s.Project.BinaryName)); err != nil {
		return pipeline.RunResult{}, err
	}
	assets := filepath.Join(h.root, s.Assets, "darwin")
	for _, name := range []string{"Info.plist", "icons.icns", "Assets.car"} {
		source := filepath.Join(assets, name)
		if _, err := os.Stat(source); err == nil {
			dest := filepath.Join(staged, "Contents", name)
			if name == "icons.icns" || name == "Assets.car" {
				dest = filepath.Join(resources, name)
			}
			if err := copyManifestPath(source, dest); err != nil {
				return pipeline.RunResult{}, err
			}
		}
	}
	associationIcons, err := filepath.Glob(filepath.Join(assets, "association-*.icns"))
	if err != nil {
		return pipeline.RunResult{}, err
	}
	for _, source := range associationIcons {
		if err := copyManifestPath(source, filepath.Join(resources, filepath.Base(source))); err != nil {
			return pipeline.RunResult{}, err
		}
	}
	return pipeline.RunResult{}, replacePathTransactional(staged, output)
}

func (h *manifestHandler) packageAppImage(ctx context.Context, s pipeline.PackageSpec) (pipeline.RunResult, error) {
	desktop, icon, buildDir, err := h.prepareAppImageInputs(s)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	outDir := filepath.Join(buildDir, "output")
	if err := removeGenerated(outDir); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	executable, err := manifestExecutable()
	if err != nil {
		return pipeline.RunResult{}, err
	}
	detail, err := runManifestCommand(ctx, h.root, nil, executable,
		"generate", "appimage",
		"-binary", filepath.Join(h.root, s.Binary),
		"-icon", icon,
		"-desktopfile", desktop,
		"-outputdir", outDir,
		"-builddir", buildDir,
	)
	if err != nil {
		return pipeline.RunResult{Detail: detail}, err
	}
	return pipeline.RunResult{}, ensureExpectedAppImage(outDir, filepath.Join(h.root, s.Output))
}

func (h *manifestHandler) prepareAppImageInputs(s pipeline.PackageSpec) (desktop, icon, buildDir string, err error) {
	buildDir = h.packageWorkspace(s)
	icon = filepath.Join(buildDir, s.Project.BinaryName+".png")
	desktop = filepath.Join(buildDir, s.Project.BinaryName+".desktop")
	iconSource := filepath.Join(h.root, s.Assets, "appicon.png")
	if s.Config.Icon != "" {
		iconSource, err = existingPathInsideProject(h.root, s.Config.Icon)
		if err != nil {
			return "", "", "", fmt.Errorf("package appimage icon: %w", err)
		}
	}
	err = prepareGeneratedWorkspace(buildDir, ".appimage-stage-*", func(staged string) error {
		if err := copyManifestPath(iconSource, filepath.Join(staged, filepath.Base(icon))); err != nil {
			return err
		}
		stagedDesktop := filepath.Join(staged, filepath.Base(desktop))
		if s.Config.DesktopEntry != "" {
			desktopSource, resolveErr := existingPathInsideProject(h.root, s.Config.DesktopEntry)
			if resolveErr != nil {
				return fmt.Errorf("package appimage desktop_entry: %w", resolveErr)
			}
			return copyManifestPath(desktopSource, stagedDesktop)
		}
		categories := "Utility;"
		if len(s.Config.Categories) > 0 {
			categories = strings.Join(s.Config.Categories, ";") + ";"
		}
		return generateDotDesktop(&DotDesktopOptions{OutputFile: stagedDesktop, Type: "Application", Name: s.Project.ProductName, Exec: s.Project.BinaryName, Icon: s.Project.BinaryName, Comment: s.Project.Description, Categories: categories, Version: "1.0"})
	})
	return desktop, icon, buildDir, err
}

func (h *manifestHandler) packageLinux(s pipeline.PackageSpec) (pipeline.RunResult, error) {
	config, err := h.prepareLinuxPackageConfig(s)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	packageRoot := h.packageWorkspace(s)
	if err := os.MkdirAll(packageRoot, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	outDir, err := os.MkdirTemp(packageRoot, "linux-output-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer os.RemoveAll(outDir)
	if err := toolPackage(&flags.ToolPackage{Format: s.Format, ExecutableName: s.Project.BinaryName, ConfigPath: config, Out: outDir}); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{}, findAndMovePackage(outDir, s.Project.BinaryName, s.Format, filepath.Join(h.root, s.Output))
}

func (h *manifestHandler) prepareLinuxPackageConfig(s pipeline.PackageSpec) (string, error) {
	workspace := h.packageWorkspace(s)
	configPath := filepath.Join(workspace, "nfpm.yaml")
	if s.Config.Template != "" {
		err := prepareGeneratedWorkspace(workspace, ".linux-package-stage-*", func(staged string) error {
			return h.copyPackageReplacement(s, filepath.Join(staged, "nfpm.yaml"))
		})
		return configPath, err
	}
	assets := filepath.Join(h.root, s.Assets)
	desktop := filepath.Join(workspace, s.Project.BinaryName+".desktop")
	source := filepath.Join(assets, "linux", "nfpm", "nfpm.yaml")
	data, err := os.ReadFile(source)
	if err != nil {
		return "", err
	}
	var config map[string]any
	if err := yaml.Unmarshal(data, &config); err != nil {
		return "", err
	}
	config["name"] = s.Project.BinaryName
	config["arch"] = s.TargetArch
	config["platform"] = "linux"
	config["version"] = s.Project.Version
	if s.Config.Maintainer != "" {
		config["maintainer"] = s.Config.Maintainer
	}
	if s.Config.Section != "" {
		config["section"] = s.Config.Section
	}
	if len(s.Config.Dependencies) > 0 {
		config["depends"] = append([]string(nil), s.Config.Dependencies...)
	}
	if contents, ok := config["contents"].([]any); ok {
		for _, raw := range contents {
			entry, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			source, _ := entry["src"].(string)
			switch {
			case strings.Contains(source, "/bin/"):
				entry["src"] = filepath.Join(h.root, s.Binary)
			case strings.HasSuffix(source, "appicon.png"):
				entry["src"] = filepath.Join(assets, "appicon.png")
			case strings.HasSuffix(source, ".desktop"):
				entry["src"] = desktop
			}
		}
	}
	if scripts, ok := config["scripts"].(map[string]any); ok {
		for name, raw := range scripts {
			path, ok := raw.(string)
			if ok {
				scripts[name] = filepath.Join(assets, "linux", "nfpm", "scripts", filepath.Base(path))
			}
		}
	}
	configuredScripts := map[string]string{
		"preinstall": s.Config.PreInstall, "postinstall": s.Config.PostInstall,
		"preremove": s.Config.PreRemove, "postremove": s.Config.PostRemove,
	}
	err = prepareGeneratedWorkspace(workspace, ".linux-package-stage-*", func(staged string) error {
		stagedDesktop := filepath.Join(staged, filepath.Base(desktop))
		generatedDesktop := filepath.Join(assets, "linux", s.Project.BinaryName+".desktop")
		if _, statErr := os.Stat(generatedDesktop); statErr == nil {
			if copyErr := copyManifestPath(generatedDesktop, stagedDesktop); copyErr != nil {
				return copyErr
			}
		} else if generateErr := generateDotDesktop(&DotDesktopOptions{OutputFile: stagedDesktop, Type: "Application", Name: s.Project.ProductName, Exec: s.Project.BinaryName, Icon: s.Project.BinaryName, Comment: s.Project.Description, Categories: "Utility;", Version: "1.0"}); generateErr != nil {
			return generateErr
		}
		if len(configuredScripts) > 0 {
			scripts, _ := config["scripts"].(map[string]any)
			if scripts == nil {
				scripts = map[string]any{}
				config["scripts"] = scripts
			}
			names := make([]string, 0, len(configuredScripts))
			for name := range configuredScripts {
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				configured := configuredScripts[name]
				if configured == "" {
					continue
				}
				source, resolveErr := existingPathInsideProject(h.root, configured)
				if resolveErr != nil {
					return fmt.Errorf("package %s %s: %w", s.Format, name, resolveErr)
				}
				finalDestination := filepath.Join(workspace, "scripts", name+filepath.Ext(source))
				stagedDestination := filepath.Join(staged, "scripts", name+filepath.Ext(source))
				if copyErr := copyManifestPath(source, stagedDestination); copyErr != nil {
					return copyErr
				}
				scripts[name] = finalDestination
			}
		}
		resolved, marshalErr := yaml.Marshal(config)
		if marshalErr != nil {
			return marshalErr
		}
		return os.WriteFile(filepath.Join(staged, "nfpm.yaml"), resolved, 0o644)
	})
	return configPath, err
}

func (h *manifestHandler) packageNSIS(ctx context.Context, s pipeline.PackageSpec) (pipeline.RunResult, error) {
	dir, packageRoot, err := h.prepareNSISWorkspace(s)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	flag := "AMD64"
	if s.TargetArch == "arm64" {
		flag = "ARM64"
	}
	binary, err := filepath.Abs(filepath.Join(h.root, s.Binary))
	if err != nil {
		return pipeline.RunResult{}, err
	}
	args := []string{"-DARG_WAILS_" + flag + "_BINARY=" + binary}
	if s.Config.InstallScope != "" {
		args = append(args, "-DWAILS_INSTALL_SCOPE="+s.Config.InstallScope)
	}
	args = append(args, "project.nsi")
	result, err := runManifestCommand(ctx, dir, nil, "makensis", args...)
	if err != nil {
		return pipeline.RunResult{Detail: result}, err
	}
	generated := filepath.Join(packageRoot, "bin", s.Project.Name+"-"+s.TargetArch+"-installer.exe")
	if err := copyPathTransactional(generated, filepath.Join(h.root, s.Output), ".nsis-output-*"); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{Detail: result}, nil
}

func (h *manifestHandler) prepareNSISWorkspace(s pipeline.PackageSpec) (dir, packageRoot string, err error) {
	packageRoot = h.packageWorkspace(s)
	dir = filepath.Join(packageRoot, "assets", "windows", "nsis")
	err = prepareGeneratedWorkspace(packageRoot, ".nsis-stage-*", func(staged string) error {
		stagedDir := filepath.Join(staged, "assets", "windows", "nsis")
		if err := copyManifestPath(filepath.Join(h.root, s.Assets, "windows", "nsis"), stagedDir); err != nil {
			return err
		}
		if s.Config.Template != "" {
			if err := h.copyPackageReplacement(s, filepath.Join(stagedDir, "project.nsi")); err != nil {
				return err
			}
		}
		return GenerateWebView2Bootstrapper(&GenerateWebView2Options{Directory: stagedDir})
	})
	return dir, packageRoot, err
}

func (h *manifestHandler) packageMSIX(s pipeline.PackageSpec) (pipeline.RunResult, error) {
	if manifestHostOS != "windows" {
		return pipeline.RunResult{}, fmt.Errorf("MSIX packaging requires Windows")
	}
	workspace := h.packageWorkspace(s)
	var associations []map[string]any
	for _, association := range s.Associations {
		for _, extension := range association.Extensions {
			associations = append(associations, map[string]any{"ext": extension, "name": association.Name, "description": association.Description, "iconName": association.Icon, "role": association.Role, "mimeType": association.MIMEType})
		}
	}
	var protocols []map[string]any
	for _, protocol := range s.Protocols {
		protocols = append(protocols, map[string]any{"scheme": protocol.Scheme, "description": protocol.Description})
	}
	config := map[string]any{"info": map[string]any{"companyName": s.Project.CompanyName, "productName": s.Project.ProductName, "productIdentifier": s.Project.Identifier, "description": s.Project.Description, "copyright": s.Project.Copyright, "comments": s.Project.Comments, "version": s.Project.Version}, "fileAssociations": associations, "protocols": protocols}
	data, err := json.Marshal(config)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	configPath := filepath.Join(workspace, "config.json")
	arch := map[string]string{"amd64": "x64", "arm64": "arm64", "386": "x86"}[s.TargetArch]
	if arch == "" {
		return pipeline.RunResult{}, fmt.Errorf("unsupported MSIX architecture %s", s.TargetArch)
	}
	signing := h.config.Signing.Windows
	publisher := "CN=" + s.Project.CompanyName
	if s.Config.Publisher != "" {
		publisher = s.Config.Publisher
	}
	if signing.Identity != "" {
		publisher = signing.Identity
	}
	appxManifest := ""
	if s.Config.Manifest != "" {
		appxManifest = filepath.Join(workspace, "AppxManifest.xml")
	}
	if err := prepareGeneratedWorkspace(workspace, ".msix-stage-*", func(staged string) error {
		if err := os.WriteFile(filepath.Join(staged, "config.json"), data, 0o644); err != nil {
			return err
		}
		if s.Config.Manifest == "" {
			return nil
		}
		source, err := existingPathInsideProject(h.root, s.Config.Manifest)
		if err != nil {
			return fmt.Errorf("package msix manifest: %w", err)
		}
		return copyManifestPath(source, filepath.Join(staged, "AppxManifest.xml"))
	}); err != nil {
		return pipeline.RunResult{}, err
	}
	generatedOutput := filepath.Join(workspace, "output", filepath.Base(s.Output))
	if err := os.MkdirAll(filepath.Dir(generatedOutput), 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := ToolMSIX(&flags.ToolMSIX{ConfigPath: configPath, Publisher: publisher, CertificatePath: signing.Certificate, Arch: arch, ExecutableName: s.Project.BinaryName + ".exe", ExecutablePath: filepath.Join(h.root, s.Binary), OutputPath: generatedOutput, AppxManifest: appxManifest, UseMakeAppx: true}); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{}, copyPathTransactional(generatedOutput, filepath.Join(h.root, s.Output), ".msix-output-*")
}

func (h *manifestHandler) packageAndroid(ctx context.Context, s pipeline.PackageSpec) (pipeline.RunResult, error) {
	workspace := h.packageWorkspace(s)
	if err := os.MkdirAll(filepath.Dir(workspace), 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	staged, err := os.MkdirTemp(filepath.Dir(workspace), ".android-package-stage-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer os.RemoveAll(staged)
	source := filepath.Join(h.root, s.Assets, "android")
	if err := copyManifestPath(source, staged); err != nil {
		return pipeline.RunResult{}, err
	}
	binaries := s.Binaries
	if len(binaries) == 0 {
		binaries = []pipeline.ComponentBinary{{Arch: s.TargetArch, Path: s.Binary}}
	}
	for _, binary := range binaries {
		abi, err := androidABI(binary.Arch)
		if err != nil {
			return pipeline.RunResult{}, err
		}
		jni := filepath.Join(staged, "app", "src", "main", "jniLibs", abi, "libwails.so")
		if err := copyManifestPath(filepath.Join(h.root, binary.Path), jni); err != nil {
			return pipeline.RunResult{}, err
		}
	}
	gradlew := filepath.Join(staged, "gradlew")
	if err := os.Chmod(gradlew, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	task := "assembleRelease"
	producedRelative := filepath.Join("app", "build", "outputs", "apk", "release", "app-release.apk")
	if s.Format == "aab" {
		task = "bundleRelease"
		producedRelative = filepath.Join("app", "build", "outputs", "bundle", "release", "app-release.aab")
	}
	output, err := runManifestCommand(ctx, staged, nil, gradlew, task)
	if err != nil {
		return pipeline.RunResult{Detail: output}, err
	}
	if err := replacePathTransactional(staged, workspace); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := copyPathTransactional(filepath.Join(workspace, producedRelative), filepath.Join(h.root, s.Output), ".android-output-*"); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{Detail: output}, nil
}

func androidABI(arch string) (string, error) {
	switch arch {
	case "arm64":
		return "arm64-v8a", nil
	case "amd64":
		return "x86_64", nil
	default:
		return "", fmt.Errorf("unsupported Android architecture %s", arch)
	}
}

func (h *manifestHandler) packageIOS(ctx context.Context, s pipeline.PackageSpec, ipa bool) (pipeline.RunResult, error) {
	if manifestHostOS != "darwin" {
		return pipeline.RunResult{}, fmt.Errorf("iOS packaging requires macOS and Xcode")
	}
	destination := s.Destination
	if destination == "" {
		destination = "simulator"
	}
	if ipa && destination != "device" {
		return pipeline.RunResult{}, fmt.Errorf("IPA packaging requires profile destination = \"device\"")
	}
	if ipa {
		return h.packageIPA(ctx, s)
	}
	sdk := "iphonesimulator"
	suffix := "-simulator"
	if destination == "device" {
		sdk = "iphoneos"
		suffix = ""
	}
	minimum := s.MinimumVersion
	if minimum == "" {
		minimum = "15.0"
	}
	target := s.TargetArch + "-apple-ios" + minimum + suffix
	sdkPath, err := commandOutput(ctx, h.root, "xcrun", "--sdk", sdk, "--show-sdk-path")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	workspace := h.packageWorkspace(s)
	if err := os.MkdirAll(filepath.Dir(workspace), 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	stagedWorkspace, err := os.MkdirTemp(filepath.Dir(workspace), ".ios-package-stage-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer os.RemoveAll(stagedWorkspace)
	executable := filepath.Join(stagedWorkspace, s.Project.BinaryName)
	mainM := filepath.Join(h.root, s.Assets, "ios", "xcode", "main", "main.m")
	frameworks := []string{"Foundation", "UIKit", "WebKit", "Security", "CoreFoundation", "UniformTypeIdentifiers", "LocalAuthentication", "UserNotifications", "AVFoundation", "CoreLocation", "CoreMotion", "SystemConfiguration"}
	args := []string{"-sdk", sdk, "clang", "-target", target, "-isysroot", sdkPath}
	for _, framework := range frameworks {
		args = append(args, "-framework", framework)
	}
	args = append(args, "-lresolv", "-o", executable, mainM, "-Wl,-force_load,"+filepath.Join(h.root, s.Binary))
	linkOutput, err := runManifestCommand(ctx, h.root, nil, "xcrun", args...)
	if err != nil {
		return pipeline.RunResult{Detail: linkOutput}, err
	}
	appOutput := filepath.Join(h.root, s.Output)
	if err := os.MkdirAll(filepath.Dir(appOutput), 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	stagedApp, err := os.MkdirTemp(filepath.Dir(appOutput), ".ios-app-stage-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer os.RemoveAll(stagedApp)
	if err := copyManifestPath(executable, filepath.Join(stagedApp, s.Project.BinaryName)); err != nil {
		return pipeline.RunResult{}, err
	}
	info := filepath.Join(h.root, s.Assets, "ios", "xcode", "main", "Info.plist")
	if err := copyManifestPath(info, filepath.Join(stagedApp, "Info.plist")); err != nil {
		return pipeline.RunResult{}, err
	}
	assetInput := filepath.Join(h.root, s.Assets, "ios", "xcode", "main", "Assets.xcassets")
	assetTemp := filepath.Join(stagedWorkspace, "compiled-assets")
	precompiledAssets := filepath.Join(h.root, s.Assets, "ios", "xcode", "main", "Assets.car")
	if _, err := os.Stat(precompiledAssets); err == nil {
		if err := copyManifestPath(precompiledAssets, filepath.Join(stagedApp, "Assets.car")); err != nil {
			return pipeline.RunResult{}, err
		}
	} else if _, err := os.Stat(assetInput); err == nil {
		if err := os.MkdirAll(assetTemp, 0o755); err != nil {
			return pipeline.RunResult{}, err
		}
		actool := []string{"actool", "--compile", assetTemp, "--app-icon", "AppIcon", "--platform", sdk, "--minimum-deployment-target", minimum, "--product-type", "com.apple.product-type.application", "--target-device", "iphone", "--target-device", "ipad", "--output-partial-info-plist", filepath.Join(stagedApp, "assetcatalog_generated_info.plist"), assetInput}
		if output, err := runManifestCommand(ctx, h.root, nil, "xcrun", actool...); err != nil {
			return pipeline.RunResult{Detail: output}, err
		}
		if _, err := os.Stat(filepath.Join(assetTemp, "Assets.car")); err == nil {
			if err := copyManifestPath(filepath.Join(assetTemp, "Assets.car"), filepath.Join(stagedApp, "Assets.car")); err != nil {
				return pipeline.RunResult{}, err
			}
		}
	}
	associationIcons, err := filepath.Glob(filepath.Join(h.root, s.Assets, "ios", "xcode", "main", "association-*"))
	if err != nil {
		return pipeline.RunResult{}, err
	}
	for _, source := range associationIcons {
		if err := copyManifestPath(source, filepath.Join(stagedApp, filepath.Base(source))); err != nil {
			return pipeline.RunResult{}, err
		}
	}
	identity := "-"
	signing := h.config.Signing.IOS
	if destination == "device" && signing.ProvisioningProfile != "" {
		profile, err := existingArtifactPathInsideProject(h.root, filepath.ToSlash(filepath.Join(s.Assets, "signing", "provisioning-profile"+filepath.Ext(signing.ProvisioningProfile))))
		if err != nil {
			return pipeline.RunResult{}, fmt.Errorf("signing.ios.provisioning_profile: %w", err)
		}
		if err := copyManifestPath(profile, filepath.Join(stagedApp, "embedded.mobileprovision")); err != nil {
			return pipeline.RunResult{}, err
		}
	}
	if signing.Enabled && signing.Identity != "" {
		identity = signing.Identity
	}
	signArgs := []string{"--force", "--sign", identity}
	if destination == "device" && signing.Entitlements != "" {
		entitlements, err := existingArtifactPathInsideProject(h.root, filepath.ToSlash(filepath.Join(s.Assets, "signing", "entitlements"+filepath.Ext(signing.Entitlements))))
		if err != nil {
			return pipeline.RunResult{}, fmt.Errorf("signing.ios.entitlements: %w", err)
		}
		signArgs = append(signArgs, "--entitlements", entitlements)
	}
	signArgs = append(signArgs, stagedApp)
	if output, err := runManifestCommand(ctx, h.root, nil, "codesign", signArgs...); err != nil {
		return pipeline.RunResult{Detail: output}, err
	}
	if err := replacePathTransactional(stagedWorkspace, workspace); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := replacePathTransactional(stagedApp, appOutput); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{Detail: linkOutput}, nil
}

func (h *manifestHandler) packageIPA(ctx context.Context, s pipeline.PackageSpec) (pipeline.RunResult, error) {
	appInput, err := existingArtifactPathInsideProject(h.root, s.Binary)
	if err != nil {
		return pipeline.RunResult{}, fmt.Errorf("IPA app bundle: %w", err)
	}
	info, err := os.Stat(appInput)
	if err != nil {
		return pipeline.RunResult{}, err
	}
	if !info.IsDir() || !strings.HasSuffix(strings.ToLower(appInput), ".app") {
		return pipeline.RunResult{}, fmt.Errorf("IPA input must be an assembled .app bundle")
	}
	workspace := h.packageWorkspace(s)
	if err := os.MkdirAll(filepath.Dir(workspace), 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	staged, err := os.MkdirTemp(filepath.Dir(workspace), ".ipa-package-stage-*")
	if err != nil {
		return pipeline.RunResult{}, err
	}
	defer os.RemoveAll(staged)
	payload := filepath.Join(staged, "Payload")
	if err := os.MkdirAll(payload, 0o755); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := copyManifestPath(appInput, filepath.Join(payload, filepath.Base(appInput))); err != nil {
		return pipeline.RunResult{}, err
	}
	generatedOutput := filepath.Join(staged, "artifact.ipa")
	zipOutput, err := runManifestCommand(ctx, staged, nil, "zip", "-qry", generatedOutput, "Payload")
	if err != nil {
		return pipeline.RunResult{Detail: zipOutput}, err
	}
	if err := replacePathTransactional(staged, workspace); err != nil {
		return pipeline.RunResult{}, err
	}
	if err := copyPathTransactional(filepath.Join(workspace, "artifact.ipa"), filepath.Join(h.root, s.Output), ".ipa-output-*"); err != nil {
		return pipeline.RunResult{}, err
	}
	return pipeline.RunResult{Detail: zipOutput}, nil
}

func (h *manifestHandler) sign(ctx context.Context, s pipeline.SignSpec) (pipeline.RunResult, error) {
	if !s.Config.Enabled {
		return pipeline.RunResult{}, fmt.Errorf("signing is not enabled for %s", s.TargetOS)
	}
	input := filepath.Join(h.root, s.Input)
	output := filepath.Join(h.root, s.Input+".signed")
	if s.TargetOS == "ios" {
		if s.Config.Identity == "" {
			return pipeline.RunResult{}, fmt.Errorf("iOS signing requires signing.ios.identity")
		}
		// iOS bundles must be signed while being assembled. packageIOS used the
		// configured identity; retain the normal signed-artifact output contract.
		return pipeline.RunResult{}, copyPathTransactional(input, output, ".ios-sign-output-*")
	}
	if s.TargetOS == "android" {
		keyAlias := chooseString(s.Config.KeyAlias != "", s.Config.KeyAlias, s.Config.Identity)
		if s.Config.Certificate == "" || keyAlias == "" || s.Config.Credential == "" {
			return pipeline.RunResult{}, fmt.Errorf("Android signing requires certificate (keystore), identity (key alias), and credential (password environment variable name)")
		}
		if os.Getenv(s.Config.Credential) == "" {
			return pipeline.RunResult{}, fmt.Errorf("Android signing password environment variable %s is empty", s.Config.Credential)
		}
		keystore, err := existingPathInsideProject(h.root, s.Config.Certificate)
		if err != nil {
			return pipeline.RunResult{}, fmt.Errorf("signing.android.certificate: %w", err)
		}
		var detail string
		err = produceNamedPathTransactional(output, ".android-sign-output-*", filepath.Base(input), func(staged string) error {
			var signErr error
			if s.Format == "apk" {
				detail, signErr = runManifestCommand(ctx, h.root, nil, "apksigner", "sign", "--ks", keystore, "--ks-key-alias", keyAlias, "--ks-pass", "env:"+s.Config.Credential, "--out", staged, input)
			} else {
				detail, signErr = runManifestCommand(ctx, h.root, nil, "jarsigner", "-keystore", keystore, "-storepass:env", s.Config.Credential, "-signedjar", staged, input, keyAlias)
			}
			return signErr
		})
		return pipeline.RunResult{Detail: detail}, err
	}
	certificate := s.Config.Certificate
	if s.TargetOS != "linux" {
		var err error
		certificate, err = existingOptionalPathInsideProject(h.root, s.Config.Certificate)
		if err != nil {
			return pipeline.RunResult{}, fmt.Errorf("signing.%s.certificate: %w", s.TargetOS, err)
		}
	}
	entitlements, err := existingOptionalArtifactPathInsideProject(h.root, s.Config.Entitlements)
	if err != nil {
		return pipeline.RunResult{}, fmt.Errorf("signing.%s.entitlements: %w", s.TargetOS, err)
	}
	err = produceNamedPathTransactional(output, ".sign-output-*", filepath.Base(input), func(staged string) error {
		if err := copyManifestPath(input, staged); err != nil {
			return err
		}
		return manifestSign(&flags.Sign{Input: staged, Certificate: certificate, Thumbprint: s.Config.Thumbprint, Timestamp: s.Config.TimestampServer, Identity: s.Config.Identity, Entitlements: entitlements, Notarize: s.Config.Notarize, KeychainProfile: s.Config.Credential, PGPKey: chooseString(s.TargetOS == "linux", certificate, ""), Role: chooseString(s.TargetOS == "linux", s.Config.Identity, "")})
	})
	return pipeline.RunResult{}, err
}

func runManifestCommand(ctx context.Context, dir string, extraEnv []string, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Env = mergeEnvironment(os.Environ(), extraEnv)
	configureManifestProcess(cmd)
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return os.ErrProcessDone
		}
		return killManifestProcess(cmd.Process)
	}
	cmd.WaitDelay = 5 * time.Second
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	return strings.TrimSpace(output.String()), err
}

func mergeEnvironment(inherited, overrides []string) []string {
	type entry struct {
		key, value string
	}
	values := make(map[string]entry, len(inherited)+len(overrides))
	add := func(raw string) {
		key, _, ok := strings.Cut(raw, "=")
		if !ok || key == "" {
			return
		}
		identity := key
		if runtime.GOOS == "windows" {
			identity = strings.ToUpper(key)
		}
		values[identity] = entry{key: key, value: raw}
	}
	for _, raw := range inherited {
		add(raw)
	}
	for _, raw := range overrides {
		add(raw)
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, values[key].value)
	}
	return result
}

func declaredEnvironment(overrides []string, declared map[string]string) []string {
	result := append([]string(nil), overrides...)
	keys := make([]string, 0, len(declared))
	for key := range declared {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		result = append(result, key+"="+declared[key])
	}
	return result
}

func runPackageManager(ctx context.Context, dir string, extraEnv []string, manager string, args ...string) (string, error) {
	if manager != "npm" {
		return runManifestCommand(ctx, dir, extraEnv, manager, args...)
	}
	npmPath, err := exec.LookPath("npm")
	if err != nil {
		return "", err
	}
	script, err := filepath.EvalSymlinks(npmPath)
	if err != nil {
		script = npmPath
	}
	nodePath, err := exec.LookPath("node")
	if err != nil {
		return "", err
	}
	return runManifestCommand(ctx, dir, extraEnv, nodePath, append([]string{script}, args...)...)
}
func toolIdentity(paths ...string) (string, error) {
	parts := make([]string, 0, len(paths))
	for _, path := range paths {
		resolved, err := filepath.EvalSymlinks(path)
		if err == nil {
			path = resolved
		}
		info, err := os.Stat(path)
		if err != nil {
			return "", err
		}
		parts = append(parts, fmt.Sprintf("%s:%d:%d:%o", path, info.Size(), info.ModTime().UnixNano(), info.Mode()))
	}
	sort.Strings(parts)
	return strings.Join(parts, "|"), nil
}
func appendUniqueStrings(values []string, extra ...string) []string {
	seen := map[string]bool{}
	for _, v := range values {
		seen[v] = true
	}
	for _, v := range extra {
		if !seen[v] {
			values = append(values, v)
			seen[v] = true
		}
	}
	return values
}
func chooseString(condition bool, value, fallback string) string {
	if condition {
		return value
	}
	return fallback
}
func existingOptionalPathInsideProject(root, value string) (string, error) {
	if value == "" {
		return "", nil
	}
	return existingPathInsideProject(root, value)
}

func existingOptionalArtifactPathInsideProject(root, value string) (string, error) {
	if value == "" {
		return "", nil
	}
	return existingArtifactPathInsideProject(root, value)
}
func pathInsideProject(root, value string) (string, error) {
	return manifest.ResolveProjectPath(root, "path", value, false)
}

func existingPathInsideProject(root, value string) (string, error) {
	return manifest.ResolveProjectPath(root, "path", value, true)
}

func writeGeneratedConfig(dir string, project manifest.Project, associations []manifest.Association, protocols []manifest.Protocol) (string, error) {
	type info struct {
		CompanyName       string `yaml:"companyName"`
		ProductName       string `yaml:"productName"`
		ProductIdentifier string `yaml:"productIdentifier"`
		Description       string `yaml:"description"`
		Copyright         string `yaml:"copyright"`
		Comments          string `yaml:"comments"`
		Version           string `yaml:"version"`
	}
	var files []FileAssociation
	for _, a := range associations {
		for _, ext := range a.Extensions {
			files = append(files, FileAssociation{Ext: ext, Name: a.Name, Description: a.Description, IconName: a.Icon, Role: a.Role, MimeType: a.MIMEType})
		}
	}
	var schemes []ProtocolConfig
	for _, p := range protocols {
		schemes = append(schemes, ProtocolConfig{Scheme: p.Scheme, Description: p.Description})
	}
	data, err := yaml.Marshal(struct {
		Info             info              `yaml:"info"`
		FileAssociations []FileAssociation `yaml:"fileAssociations,omitempty"`
		Protocols        []ProtocolConfig  `yaml:"protocols,omitempty"`
	}{Info: info{project.CompanyName, project.ProductName, project.Identifier, project.Description, project.Copyright, project.Comments, project.Version}, FileAssociations: files, Protocols: schemes})
	if err != nil {
		return "", err
	}
	path := filepath.Join(dir, "config.yml")
	return path, os.WriteFile(path, data, 0o644)
}

func removeGenerated(path string) error {
	clean := filepath.Clean(path)
	if !strings.Contains(filepath.ToSlash(clean), "/.wails/") && !strings.HasSuffix(filepath.ToSlash(clean), ".app") {
		if _, err := os.Stat(clean); errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("refusing to replace non-generated path %s", path)
	}
	return os.RemoveAll(clean)
}
func copyManifestPath(source, destination string) error {
	info, err := os.Lstat(source)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(source)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
			return err
		}
		if err := os.Remove(destination); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return os.Symlink(target, destination)
	}
	if info.IsDir() {
		if err := os.MkdirAll(destination, info.Mode().Perm()); err != nil {
			return err
		}
		entries, err := os.ReadDir(source)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := copyManifestPath(filepath.Join(source, entry.Name()), filepath.Join(destination, entry.Name())); err != nil {
				return err
			}
		}
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode().Perm())
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
func ensureExpectedAppImage(dir, expected string) error {
	matches, err := filepath.Glob(filepath.Join(dir, "*.AppImage"))
	if err != nil {
		return err
	}
	if len(matches) != 1 {
		return fmt.Errorf("AppImage packager produced %d candidates, expected %s", len(matches), expected)
	}
	return copyPathTransactional(matches[0], expected, ".appimage-output-*")
}
func findAndMovePackage(dir, name, format, expected string) error {
	extension := "." + format
	if format == "archlinux" {
		extension = ".pkg.tar.zst"
	}
	matches, err := filepath.Glob(filepath.Join(dir, name+"*"+extension))
	if err != nil {
		return err
	}
	if len(matches) != 1 {
		return fmt.Errorf("packager produced %d %s candidates", len(matches), format)
	}
	return copyPathTransactional(matches[0], expected, ".linux-package-output-*")
}

var _ = runtime.GOOS
