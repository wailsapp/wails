// Copy into v3/internal/wake/pipeline/review_cache_inputs_test.go; run from v3:
// go test ./internal/wake/pipeline -run TestReview -v
// These tests PASS when the reviewed bug reproduces.
package pipeline

import (
 "os"
 "path/filepath"
 "testing"
 "github.com/wailsapp/wails/v3/internal/wake/cache"
)
func TestReviewEmbeddedDataCache(t *testing.T) {
 c := testConfig(t)
 if err := os.WriteFile(filepath.Join(c.Root,"main.go"), []byte("package main\nimport _ \"embed\"\n//go:embed data.json\nvar data string\nfunc main() {}\n"),0600); err != nil {t.Fatal(err)}
 p:=filepath.Join(c.Root,"data.json")
 if err:=os.WriteFile(p,[]byte("before"),0600);err!=nil {t.Fatal(err)}
 plan,err:=PlanBuild(c,Request{Verb:"build",TargetOS:"linux",TargetArch:"amd64"});if err!=nil {t.Fatal(err)}
 node:=plan.Nodes["target:linux/amd64:compile"]
 store,err:=cache.OpenCache(c.Root);if err!=nil {t.Fatal(err)}
 before,err:=snapshotNodeInputs(store,node);if err!=nil {t.Fatal(err)}
 if err:=os.WriteFile(p,[]byte("after"),0600);err!=nil {t.Fatal(err)}
 store.InvalidateObservations()
 after,err:=snapshotNodeInputs(store,node);if err!=nil {t.Fatal(err)}
 if len(before)!=len(after) {t.Fatal("changed")}
 for i:= range before { if before[i]!=after[i] {t.Fatal("snapshot did invalidate")}}
 t.Log("CONFIRMED: compile snapshots unchanged after embedded JSON edit")
}
func TestReviewCustomFrontendCommand(t *testing.T) {
 c:=testConfig(t)
 c.Frontend.Install=[]string{"python3","install.py"}
 host:=CurrentHostCapabilities()
 if _,err:=os.Stat("/usr/bin/python3");err!=nil {t.Skip(err)}
 _,err:=PlanBuildForHost(c,Request{Verb:"build",TargetOS:"linux",TargetArch:"amd64"},host)
 if err==nil {t.Fatal("unexpected acceptance")}
 t.Logf("CONFIRMED: /usr/bin/python3 exists but plan rejected: %v",err)
}
