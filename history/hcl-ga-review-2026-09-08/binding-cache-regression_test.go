// Copy into v3/internal/wake/cache/review_method_registration_test.go; run from v3:
// go test ./internal/wake/cache -run TestReview -v
// This test PASSES when the reviewed bug reproduces.
package cache
import("os"; "path/filepath"; "testing")
func TestReviewMethodRegistrationIgnored(t *testing.T) {
 root:=t.TempDir(); source:=filepath.Join(root,"service.go")
 before:=[]byte("package app\nimport \"github.com/wailsapp/wails/v3/pkg/application\"\ntype Service struct{}\nfunc (*Service) Startup() { application.RegisterEvent[string](\"first\") }\n")
 after:=[]byte("package app\nimport \"github.com/wailsapp/wails/v3/pkg/application\"\ntype Service struct{}\nfunc (*Service) Startup() { application.RegisterEvent[int](\"other\") }\n")
 if e:=os.WriteFile(source,before,0600);e!=nil{t.Fatal(e)}
 c,e:=OpenCache(root);if e!=nil{t.Fatal(e)}
 a,e:=c.SnapshotGoAPI(SnapshotOptions{Root:root});if e!=nil{t.Fatal(e)}
 if e:=os.WriteFile(source,after,0600);e!=nil{t.Fatal(e)}
 c.InvalidateObservations();b,e:=c.SnapshotGoAPI(SnapshotOptions{Root:root});if e!=nil{t.Fatal(e)}
 if a!=b {t.Fatal("not reproduced")};t.Log("CONFIRMED: changing registered event name AND payload in method gives identical binding API digest")
}
