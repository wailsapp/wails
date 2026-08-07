package commands

// The CLI MCP server deliberately lives beside the CLI commands, rather than
// inside the application MCP implementation. The application server controls
// a running WebView; this server controls the project lifecycle.

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	mcpTokenEnv  = "WAILS_MCP_TOKEN"
	maxMCPOutput = 2 << 20
	maxMCPJobs   = 32
	maxMCPArgs   = 64
)

// MCPOptions configures the local project MCP server. Transport is selected
// automatically: stdio when launched by an agent (no terminal) and loopback
// Streamable HTTP when run interactively. Both modes can be forced explicitly.
type MCPOptions struct {
	Token string `name:"token" description:"Bearer token for mutating and process-control tools"`
	Root  string `name:"root" description:"Allowed project root (defaults to the current directory)"`
	HTTP  bool   `name:"http" description:"Use Streamable HTTP instead of automatic transport selection"`
	Stdio bool   `name:"stdio" description:"Use stdio instead of automatic transport selection"`
	Port  int    `name:"port" description:"HTTP port (0 chooses a free loopback port)" default:"0"`
}

type mcpServer struct {
	root  string
	token string
	jobs  *mcpJobs
	ctx   context.Context
}

// MCP starts a local MCP server. Transport is selected automatically: stdio
// when no terminal is attached and loopback Streamable HTTP when run
// interactively. The token is returned in the MCP initialize instructions so
// an agent can use it without contaminating the JSON-RPC stdout stream. It is
// also written to stderr for humans debugging a manually launched server.
func MCP(options *MCPOptions) error {
	if options == nil {
		options = &MCPOptions{}
	}
	// MCP owns stdout for JSON-RPC. Never let the normal CLI footer or a
	// terminal renderer corrupt the protocol stream.
	DisableFooter = true
	root, err := resolveMCPRoot(options.Root)
	if err != nil {
		return err
	}
	token := strings.TrimSpace(options.Token)
	if token == "" {
		token = strings.TrimSpace(os.Getenv(mcpTokenEnv))
	}
	if token == "" {
		token, err = newMCPToken()
		if err != nil {
			return fmt.Errorf("generate MCP session token: %w", err)
		}
	}

	serverCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	server := &mcpServer{root: root, token: token, jobs: newMCPJobs(), ctx: serverCtx}
	if options.HTTP && options.Stdio {
		return errors.New("choose at most one of --http and --stdio")
	}
	useHTTP := options.HTTP || (!options.Stdio && terminalMCPMode())
	if useHTTP {
		return server.runHTTP(serverCtx, token, options.Port)
	}
	protocolServer := server.protocolServer()
	fmt.Fprintf(os.Stderr, "wails3 MCP server ready over stdio; root=%s token=%s\n", root, token)
	err = protocolServer.Run(serverCtx, &mcp.StdioTransport{})
	if err == nil {
		return nil
	}
	if errors.Is(err, io.EOF) || strings.HasSuffix(err.Error(), ": EOF") {
		return nil
	}
	return err
}

func terminalMCPMode() bool {
	stdin, err := os.Stdin.Stat()
	if err != nil || stdin.Mode()&os.ModeCharDevice == 0 {
		return false
	}
	stdout, err := os.Stdout.Stat()
	return err == nil && stdout.Mode()&os.ModeCharDevice != 0
}

func (s *mcpServer) protocolServer() *mcp.Server {
	impl := &mcp.Implementation{Name: "wails3-project", Version: BuildSettings["vcs.revision"]}
	if impl.Version == "" {
		impl.Version = "development"
	}
	instructions := fmt.Sprintf("Wails project MCP session. Allowed root: %s. Session token: %s. Include the token in mutating and process-control tool arguments as `token`. The token is scoped to this server process and is not a substitute for user approval.", s.root, s.token)
	server := mcp.NewServer(impl, &mcp.ServerOptions{Instructions: instructions})
	s.registerTools(server)
	return server
}

func (s *mcpServer) runHTTP(ctx context.Context, token string, port int) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("invalid HTTP port %d", port)
	}
	listener, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", fmt.Sprint(port)))
	if err != nil {
		return fmt.Errorf("listen for MCP HTTP: %w", err)
	}
	protocolServer := s.protocolServer()
	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return protocolServer }, nil)
	verifier := func(_ context.Context, candidate string, _ *http.Request) (*auth.TokenInfo, error) {
		if subtle.ConstantTimeCompare([]byte(candidate), []byte(token)) != 1 {
			return nil, auth.ErrInvalidToken
		}
		return &auth.TokenInfo{Expiration: time.Now().Add(24 * time.Hour)}, nil
	}
	httpServer := &http.Server{Handler: auth.RequireBearerToken(verifier, nil)(handler), ReadHeaderTimeout: 10 * time.Second}
	go func() {
		<-ctx.Done()
		_ = httpServer.Shutdown(context.Background())
	}()
	fmt.Fprintf(os.Stderr, "wails3 MCP server ready over Streamable HTTP\nendpoint=http://%s/mcp\nroot=%s\nbearer-token=%s\n", listener.Addr().String(), s.root, token)
	err = httpServer.Serve(listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func resolveMCPRoot(input string) (string, error) {
	root := input
	if root == "" {
		var err error
		root, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("get current directory: %w", err)
		}
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve MCP root: %w", err)
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "", fmt.Errorf("resolve MCP root %q: %w", root, err)
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		if err == nil {
			err = errors.New("not a directory")
		}
		return "", fmt.Errorf("MCP root %q: %w", root, err)
	}
	return root, nil
}

func newMCPToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

type mcpAuthInput struct {
	Token string `json:"token" jsonschema:"session token from server instructions"`
}

type mcpReadInput struct{}

type mcpReadPathInput struct {
	Path string `json:"path,omitempty" jsonschema:"path relative to the allowed project root; defaults to the root"`
}

type mcpInspectOutput struct {
	Root      string   `json:"root"`
	Path      string   `json:"path"`
	Exists    bool     `json:"exists"`
	IsProject bool     `json:"isWailsProject"`
	Files     []string `json:"files,omitempty"`
	Taskfile  bool     `json:"taskfile"`
	Config    bool     `json:"config"`
	GoModule  bool     `json:"goModule"`
	Frontend  bool     `json:"frontend"`
	FileCount int      `json:"fileCount"`
	Truncated bool     `json:"truncated"`
}

type mcpInitInput struct {
	mcpAuthInput
	Directory     string `json:"directory,omitempty" jsonschema:"new project directory relative to the allowed root"`
	Name          string `json:"name" jsonschema:"project name"`
	Template      string `json:"template,omitempty" jsonschema:"built-in or trusted template name; defaults to vanilla"`
	Module        string `json:"module,omitempty" jsonschema:"Go module path"`
	Git           string `json:"git,omitempty" jsonschema:"optional Git remote URL"`
	AllowExternal bool   `json:"allowExternal,omitempty" jsonschema:"explicitly allow remote templates or Git remotes"`
}

type mcpCommandInput struct {
	mcpAuthInput
	Path string   `json:"path,omitempty" jsonschema:"project path relative to the allowed root"`
	Args []string `json:"args,omitempty" jsonschema:"additional validated Wails CLI arguments"`
}

type mcpTaskInput struct {
	mcpCommandInput
	Task string `json:"task" jsonschema:"Taskfile task name"`
}

type mcpJobOutput struct {
	JobID  string `json:"jobId"`
	State  string `json:"state"`
	Path   string `json:"path"`
	Output string `json:"output,omitempty"`
}

type mcpJobInput struct {
	mcpAuthInput
	JobID string `json:"jobId" jsonschema:"job ID returned by a process-control tool"`
}

type mcpJobReadInput struct {
	JobID string `json:"jobId" jsonschema:"job ID returned by a process-control tool"`
}

func (s *mcpServer) registerTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{Name: "wails_project_inspect", Description: "Inspect a Wails project within the allowed root without modifying it.", Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: boolPtr(false)}}, s.inspect)
	mcp.AddTool(server, &mcp.Tool{Name: "wails_project_init", Description: "Initialize a new Wails project within the allowed root. Remote templates and Git URLs are external inputs and should be explicitly approved.", Annotations: &mcp.ToolAnnotations{OpenWorldHint: boolPtr(true), DestructiveHint: boolPtr(false)}}, s.init)
	mcp.AddTool(server, &mcp.Tool{Name: "wails_project_doctor", Description: "Run Wails environment diagnostics and return machine-readable JSON.", Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: boolPtr(false)}}, s.doctor)
	mcp.AddTool(server, &mcp.Tool{Name: "wails_project_task_list", Description: "List the Taskfile tasks available in a project without running them.", Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: boolPtr(false)}}, s.taskList)
	mcp.AddTool(server, &mcp.Tool{Name: "wails_project_build", Description: "Start a Wails build job inside the allowed root.", Annotations: &mcp.ToolAnnotations{OpenWorldHint: boolPtr(false), DestructiveHint: boolPtr(false)}}, s.build)
	mcp.AddTool(server, &mcp.Tool{Name: "wails_project_dev_start", Description: "Start Wails development mode inside the allowed root.", Annotations: &mcp.ToolAnnotations{OpenWorldHint: boolPtr(false), DestructiveHint: boolPtr(false)}}, s.dev)
	mcp.AddTool(server, &mcp.Tool{Name: "wails_project_generate_bindings", Description: "Start Wails Go-to-TypeScript binding generation inside the allowed root.", Annotations: &mcp.ToolAnnotations{OpenWorldHint: boolPtr(false), DestructiveHint: boolPtr(false)}}, s.bindings)
	mcp.AddTool(server, &mcp.Tool{Name: "wails_project_task_run", Description: "Start a named Taskfile task inside the allowed root. Arbitrary shell commands are not accepted as tool input.", Annotations: &mcp.ToolAnnotations{OpenWorldHint: boolPtr(false), DestructiveHint: boolPtr(false)}}, s.task)
	mcp.AddTool(server, &mcp.Tool{Name: "wails_job_status", Description: "Read status and bounded initial output for a Wails MCP job.", Annotations: &mcp.ToolAnnotations{ReadOnlyHint: true, OpenWorldHint: boolPtr(false)}}, s.status)
	mcp.AddTool(server, &mcp.Tool{Name: "wails_job_stop", Description: "Stop a Wails MCP job started by this server.", Annotations: &mcp.ToolAnnotations{OpenWorldHint: boolPtr(false), DestructiveHint: boolPtr(false)}}, s.stop)
}

func boolPtr(v bool) *bool { return &v }

func (s *mcpServer) authorize(token string) error {
	if subtle.ConstantTimeCompare([]byte(token), []byte(s.token)) != 1 {
		return errors.New("invalid or missing MCP session token")
	}
	return nil
}

func (s *mcpServer) projectPath(input string) (string, error) {
	if input == "" {
		return s.root, nil
	}
	p, err := filepath.Abs(filepath.Join(s.root, input))
	if err != nil {
		return "", err
	}
	cleanRoot := filepath.Clean(s.root)
	cleanPath := filepath.Clean(p)
	if cleanPath != cleanRoot && !strings.HasPrefix(cleanPath, cleanRoot+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q is outside the allowed MCP root", input)
	}
	// Resolve symlinks on the longest existing prefix to catch cases where an
	// intermediate directory component is a symlink that escapes the root even
	// when the final path component does not yet exist.
	check := cleanPath
	for check != filepath.Dir(check) {
		if _, err := os.Lstat(check); err == nil {
			realPath, err := filepath.EvalSymlinks(check)
			if err != nil || (realPath != cleanRoot && !strings.HasPrefix(realPath, cleanRoot+string(filepath.Separator))) {
				return "", fmt.Errorf("path %q resolves outside the allowed MCP root", input)
			}
			break
		}
		check = filepath.Dir(check)
	}
	return cleanPath, nil
}

func (s *mcpServer) inspect(ctx context.Context, _ *mcp.CallToolRequest, in mcpReadPathInput) (*mcp.CallToolResult, mcpInspectOutput, error) {
	path, err := s.projectPath(in.Path)
	if err != nil {
		return nil, mcpInspectOutput{}, err
	}
	out := mcpInspectOutput{Root: s.root, Path: path}
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, out, nil
	}
	if err != nil {
		return nil, out, err
	}
	out.Exists = true
	if !info.IsDir() {
		return nil, out, errors.New("project path is not a directory")
	}
	err = filepath.WalkDir(path, func(p string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() && (d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == "bin") {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(path, p)
		if err != nil {
			return err
		}
		out.FileCount++
		if len(out.Files) < 1000 {
			out.Files = append(out.Files, filepath.ToSlash(rel))
		} else {
			out.Truncated = true
		}
		if rel == "Taskfile.yml" || rel == "Taskfile.yaml" {
			out.Taskfile = true
		}
		if rel == "build/config.yml" {
			out.Config = true
		}
		if rel == "go.mod" {
			out.GoModule = true
		}
		if rel == "frontend/package.json" {
			out.Frontend = true
		}
		return nil
	})
	if err != nil {
		return nil, out, err
	}
	out.IsProject = out.Taskfile && out.Config && out.GoModule
	return nil, out, nil
}

func (s *mcpServer) init(ctx context.Context, _ *mcp.CallToolRequest, in mcpInitInput) (*mcp.CallToolResult, mcpInspectOutput, error) {
	if err := s.authorize(in.Token); err != nil {
		return nil, mcpInspectOutput{}, err
	}
	if strings.TrimSpace(in.Name) == "" {
		return nil, mcpInspectOutput{}, errors.New("name is required")
	}
	if (in.Git != "" || strings.Contains(in.Template, "://")) && !in.AllowExternal {
		return nil, mcpInspectOutput{}, errors.New("remote templates and Git remotes require allowExternal=true")
	}
	dir, err := s.projectPath(in.Directory)
	if err != nil {
		return nil, mcpInspectOutput{}, err
	}
	args := []string{"init", "-n", in.Name, "-d", dir, "-q"}
	if in.Template != "" {
		args = append(args, "-t", in.Template)
	}
	if in.Module != "" {
		args = append(args, "-mod", in.Module)
	}
	if in.Git != "" {
		args = append(args, "-git", in.Git)
	}
	// Run with the allowed root as the working directory; the target project
	// directory may not exist yet and exec would fail if used as cmd.Dir.
	result, err := s.runSync(ctx, s.root, args...)
	if err != nil {
		return nil, mcpInspectOutput{}, fmt.Errorf("initialize project: %w: %s", err, result)
	}
	return nil, mcpInspectOutput{Root: s.root, Path: dir, Exists: true}, nil
}

func (s *mcpServer) doctor(ctx context.Context, _ *mcp.CallToolRequest, in mcpReadInput) (*mcp.CallToolResult, map[string]any, error) {
	result, err := s.runSync(ctx, s.root, "doctor", "--json")
	if err != nil {
		return nil, nil, fmt.Errorf("doctor: %w: %s", err, result)
	}
	return nil, map[string]any{"root": s.root, "report": result}, nil
}

func (s *mcpServer) taskList(ctx context.Context, _ *mcp.CallToolRequest, in mcpReadPathInput) (*mcp.CallToolResult, map[string]any, error) {
	path, err := s.projectPath(in.Path)
	if err != nil {
		return nil, nil, err
	}
	result, err := s.runSync(ctx, path, "task", "--list", "--json")
	if err != nil {
		return nil, nil, fmt.Errorf("list tasks: %w: %s", err, result)
	}
	return nil, map[string]any{"path": path, "tasks": result}, nil
}

func (s *mcpServer) build(ctx context.Context, _ *mcp.CallToolRequest, in mcpCommandInput) (*mcp.CallToolResult, mcpJobOutput, error) {
	if err := s.authorize(in.Token); err != nil {
		return nil, mcpJobOutput{}, err
	}
	path, err := s.projectPath(in.Path)
	if err != nil {
		return nil, mcpJobOutput{}, err
	}
	if err := validateMCPArgs(in.Args); err != nil {
		return nil, mcpJobOutput{}, err
	}
	return s.startJob(ctx, path, append([]string{"build"}, in.Args...))
}

func (s *mcpServer) dev(ctx context.Context, _ *mcp.CallToolRequest, in mcpCommandInput) (*mcp.CallToolResult, mcpJobOutput, error) {
	if err := s.authorize(in.Token); err != nil {
		return nil, mcpJobOutput{}, err
	}
	path, err := s.projectPath(in.Path)
	if err != nil {
		return nil, mcpJobOutput{}, err
	}
	if err := validateMCPArgs(in.Args); err != nil {
		return nil, mcpJobOutput{}, err
	}
	return s.startJob(ctx, path, append([]string{"dev"}, in.Args...))
}

func (s *mcpServer) bindings(ctx context.Context, _ *mcp.CallToolRequest, in mcpCommandInput) (*mcp.CallToolResult, mcpJobOutput, error) {
	if err := s.authorize(in.Token); err != nil {
		return nil, mcpJobOutput{}, err
	}
	path, err := s.projectPath(in.Path)
	if err != nil {
		return nil, mcpJobOutput{}, err
	}
	if err := validateMCPArgs(in.Args); err != nil {
		return nil, mcpJobOutput{}, err
	}
	return s.startJob(ctx, path, append([]string{"generate", "bindings"}, in.Args...))
}

func (s *mcpServer) task(ctx context.Context, _ *mcp.CallToolRequest, in mcpTaskInput) (*mcp.CallToolResult, mcpJobOutput, error) {
	if err := s.authorize(in.Token); err != nil {
		return nil, mcpJobOutput{}, err
	}
	if strings.TrimSpace(in.Task) == "" || strings.ContainsAny(in.Task, "\r\n") {
		return nil, mcpJobOutput{}, errors.New("task is required and may not contain newlines")
	}
	if err := validateMCPArgs(in.Args); err != nil {
		return nil, mcpJobOutput{}, err
	}
	path, err := s.projectPath(in.Path)
	if err != nil {
		return nil, mcpJobOutput{}, err
	}
	return s.startJob(ctx, path, append([]string{"task", in.Task}, in.Args...))
}

func (s *mcpServer) status(ctx context.Context, _ *mcp.CallToolRequest, in mcpJobReadInput) (*mcp.CallToolResult, mcpJobOutput, error) {
	job, ok := s.jobs.get(in.JobID)
	if !ok {
		return nil, mcpJobOutput{}, errors.New("unknown job")
	}
	return nil, job.snapshot(), nil
}

func (s *mcpServer) stop(ctx context.Context, _ *mcp.CallToolRequest, in mcpJobInput) (*mcp.CallToolResult, mcpJobOutput, error) {
	if err := s.authorize(in.Token); err != nil {
		return nil, mcpJobOutput{}, err
	}
	job, ok := s.jobs.get(in.JobID)
	if !ok {
		return nil, mcpJobOutput{}, errors.New("unknown job")
	}
	job.cancel()
	return nil, job.snapshot(), nil
}

func (s *mcpServer) runSync(ctx context.Context, dir string, args ...string) (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	cmd := exec.CommandContext(ctx, executable, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "NO_COLOR=1", "WAILS_MCP_CHILD=1")
	var output cappedBuffer
	cmd.Stdout, cmd.Stderr = &output, &output
	err = cmd.Run()
	return output.String(), err
}

func (s *mcpServer) startJob(ctx context.Context, dir string, args []string) (*mcp.CallToolResult, mcpJobOutput, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, mcpJobOutput{}, err
	}
	// A request context ends when the tool response is delivered. Long-running
	// jobs must therefore be tied to the server lifetime, not the call lifetime.
	jobParent := s.ctx
	if jobParent == nil {
		jobParent = context.Background()
	}
	jobCtx, cancel := context.WithCancel(jobParent)
	job, err := newMCPJob(dir, cancel)
	if err != nil {
		cancel()
		return nil, mcpJobOutput{}, err
	}
	if !s.jobs.add(job) {
		cancel()
		return nil, mcpJobOutput{}, errors.New("too many active or retained MCP jobs")
	}
	cmd := exec.CommandContext(jobCtx, executable, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "NO_COLOR=1", "WAILS_MCP_CHILD=1")
	cmd.Stdout, cmd.Stderr = job, job
	go func() { err := cmd.Run(); job.finish(err) }()
	return nil, job.snapshot(), nil
}

type cappedBuffer struct {
	mu sync.Mutex
	b  []byte
}

func (b *cappedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.b) < maxMCPOutput {
		n := maxMCPOutput - len(b.b)
		if len(p) < n {
			n = len(p)
		}
		b.b = append(b.b, p[:n]...)
	}
	return len(p), nil
}
func (b *cappedBuffer) String() string { b.mu.Lock(); defer b.mu.Unlock(); return string(b.b) }

type mcpJob struct {
	mu                      sync.Mutex
	id, path, state, output string
	cancel                  context.CancelFunc
}

func newMCPJob(path string, cancel context.CancelFunc) (*mcpJob, error) {
	id, err := newMCPToken()
	if err != nil {
		return nil, fmt.Errorf("generate MCP job ID: %w", err)
	}
	return &mcpJob{id: id, path: path, state: "running", cancel: cancel}, nil
}
func (j *mcpJob) Write(p []byte) (int, error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if len(j.output) < maxMCPOutput {
		n := maxMCPOutput - len(j.output)
		if len(p) < n {
			n = len(p)
		}
		j.output += string(p[:n])
	}
	return len(p), nil
}
func (j *mcpJob) finish(err error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if errors.Is(err, context.Canceled) {
		j.state = "stopped"
	} else if err != nil {
		j.state = "failed"
	} else {
		j.state = "completed"
	}
}
func (j *mcpJob) snapshot() mcpJobOutput {
	j.mu.Lock()
	defer j.mu.Unlock()
	return mcpJobOutput{JobID: j.id, State: j.state, Path: j.path, Output: j.output}
}

type mcpJobs struct {
	mu   sync.Mutex
	jobs map[string]*mcpJob
}

func newMCPJobs() *mcpJobs { return &mcpJobs{jobs: make(map[string]*mcpJob)} }
func (j *mcpJobs) add(job *mcpJob) bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	if len(j.jobs) >= maxMCPJobs {
		return false
	}
	j.jobs[job.id] = job
	return true
}
func (j *mcpJobs) get(id string) (*mcpJob, bool) {
	j.mu.Lock()
	defer j.mu.Unlock()
	job, ok := j.jobs[id]
	return job, ok
}

func validateMCPArgs(args []string) error {
	if len(args) > maxMCPArgs {
		return fmt.Errorf("too many command arguments: maximum is %d", maxMCPArgs)
	}
	// Flags that accept path arguments and could be used to write files outside
	// the allowed root. Reject any use of these flags; the MCP tools already
	// pin the working directory to a validated project path.
	pathFlags := map[string]bool{
		"-d": true, "--dir": true,
		"-o": true, "--output": true,
		"--config": true,
	}
	for i, arg := range args {
		if strings.IndexByte(arg, 0) >= 0 {
			return errors.New("command arguments may not contain NUL bytes")
		}
		if len(arg) > 4096 {
			return errors.New("command arguments may not exceed 4096 bytes")
		}
		// Block path-valued flags that could redirect output outside the root.
		if pathFlags[arg] {
			return fmt.Errorf("flag %q is not permitted in MCP tool arguments; use the path parameter instead", arg)
		}
		// Also handle --flag=value forms.
		for flag := range pathFlags {
			if strings.HasPrefix(arg, flag+"=") {
				return fmt.Errorf("flag %q is not permitted in MCP tool arguments; use the path parameter instead", arg[:len(flag)])
			}
		}
	}
	return nil
}
