package commands
import (
 "context"
 "errors"
 "os"
 "path/filepath"
 "strings"
 "testing"
 "time"
 "syscall"
 "strconv"
 "github.com/wailsapp/wails/v3/internal/wake/manifest"
)
func TestReviewDevUsesDiscoveredRoot(t *testing.T) {
 ops, _ := newManifestDevTestOps(t)
 loaded, _ := ops.load("", "")
 ops.getwd = func()(string,error){ return filepath.Join(loaded.Config.Root,"subdir"),nil }
 ops.startFrontend = func(root string, _ manifest.Config, _ string, _ int, _ string)(*manifestProcess,error){
  if root != loaded.Config.Root { t.Errorf("frontend root = %q; want discovered root %q",root,loaded.Config.Root) }
  return nil,errors.New("stop after observation")
 }
 _ = runManifestDevContextWithOps(context.Background(), &DevOptions{},ops)
}
func TestReviewMigrationDoesNotLoseTaskfileEnvironment(t *testing.T) {
 root:=t.TempDir()
 if err:=copyMigrationFixture(filepath.Join("..","..","examples","badge"),root); err!=nil {t.Fatal(err)}
 file:=filepath.Join(root,"Taskfile.yml")
 data,err:=os.ReadFile(file); if err!=nil {t.Fatal(err)}
 data=append(data,[]byte("\nenv:\n  CGO_CFLAGS: '-DREVIEW_REQUIRED_FEATURE=1'\n")...)
 if err=os.WriteFile(file,data,0644);err!=nil{t.Fatal(err)}
 report,doc,err:=analyseMigration(root);if err!=nil{t.Fatal(err)}
 encoded,err:=manifest.EncodeDocument(doc);if err!=nil{t.Fatal(err)}
 if report.Complete && !strings.Contains(string(encoded),"REVIEW_REQUIRED_FEATURE") {t.Errorf("migration says complete but silently drops Taskfile CGO_CFLAGS; diagnostics=%+v",report.Diagnostics)}
}
func TestReviewMigrationDoesNotLosePort(t *testing.T) {
 root:=t.TempDir()
 if err:=copyMigrationFixture(filepath.Join("..","..","examples","badge"),root); err!=nil {t.Fatal(err)}
 file:=filepath.Join(root,"Taskfile.yml")
 data,err:=os.ReadFile(file); if err!=nil {t.Fatal(err)}
 data=[]byte(strings.Replace(string(data),"'{{.WAILS_VITE_PORT | default 9245}}'","9357",1))
 if err=os.WriteFile(file,data,0644);err!=nil{t.Fatal(err)}
 report,doc,err:=analyseMigration(root);if err!=nil{t.Fatal(err)}
 encoded,err:=manifest.EncodeDocument(doc);if err!=nil{t.Fatal(err)}
 if report.Complete && !strings.Contains(string(encoded),"9357") {t.Errorf("migration says complete but loses VITE_PORT=9357; internal Port=%d",doc.Dev.Port)}
}

func TestReviewFrontendDevReceivesEnvironment(t *testing.T) {
 root:=t.TempDir()
 config:=manifest.Config{Frontend:manifest.Frontend{Directory:".",Dev:[]string{"sh","-c","printf '%s' \"$WAILS_REVIEW_FRONTEND_VALUE\" > env-result"},Environment:map[string]string{"WAILS_REVIEW_FRONTEND_VALUE":"configured-value"}}}
 p,err:=startFrontendDev(root,config,"127.0.0.1",9245,"http://127.0.0.1:9245"); if err!=nil{t.Fatal(err)}
 <-p.done
 data,err:=os.ReadFile(filepath.Join(root,"env-result"));if err!=nil{t.Fatal(err)}
 if string(data)!="configured-value"{t.Errorf("frontend.dev environment = %q; want configured-value",data)}
}
func TestReviewStopCleansChildrenAfterParentExit(t *testing.T) {
 root:=t.TempDir()
 p,err:=startManifestProcess(root,"sh",nil,"-c","sleep 60 & echo $! > child-pid")
 if err!=nil{t.Fatal(err)}
 <-p.done
 data,err:=os.ReadFile(filepath.Join(root,"child-pid"));if err!=nil{t.Fatal(err)}
 pid,err:=strconv.Atoi(strings.TrimSpace(string(data)));if err!=nil{t.Fatal(err)}
 defer syscall.Kill(pid,syscall.SIGKILL)
 p.stop(10*time.Millisecond)
 if syscall.Kill(pid,0)==nil{t.Errorf("child PID %d remains alive after stop on exited parent",pid)}
}
