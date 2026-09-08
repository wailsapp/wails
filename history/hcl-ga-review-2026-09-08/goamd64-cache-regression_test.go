// Copy into v3/internal/commands/review_goamd64_cache_test.go; run from v3:
// go test ./internal/commands -run TestReviewGOAMD64 -v
// This test PASSES when the reviewed bug reproduces.
package commands
import("context";"testing";"github.com/wailsapp/wails/v3/internal/wake/pipeline";"github.com/wailsapp/wails/v3/internal/wake/cache")
func TestReviewGOAMD64CacheIdentity(t *testing.T) {
 node:=pipeline.Node{Kind:pipeline.CompileApplication,Spec:pipeline.CompileSpec{TargetOS:"linux",TargetArch:"amd64",Production:true,Toolchain:"native"}}
 handler:=&manifestHandler{}
 t.Setenv("GOAMD64","v3")
 before,err:=handler.Identity(context.Background(),node);if err!=nil{t.Fatal(err)}
 beforeKey,err:=cache.ActionKey(string(node.Kind),map[string]any{"spec":node.Spec,"tool":before},nil,nil);if err!=nil{t.Fatal(err)}
 t.Setenv("GOAMD64","v1")
 after,err:=handler.Identity(context.Background(),node);if err!=nil{t.Fatal(err)}
 afterKey,err:=cache.ActionKey(string(node.Kind),map[string]any{"spec":node.Spec,"tool":after},nil,nil);if err!=nil{t.Fatal(err)}
 if before!=after || beforeKey!=afterKey {t.Fatal("bug did not reproduce: GOAMD64 invalidated identity")}
 t.Logf("CONFIRMED: GOAMD64=v3 and GOAMD64=v1 share handler identity and action key %s",beforeKey)
}
