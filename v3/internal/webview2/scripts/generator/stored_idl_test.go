package generator

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"updater/generator/idlreader"
)

func TestStoredIDLRegression(t *testing.T) {
	for _, tt := range []struct{ name, model, output string }{
		{"WebView2.1.0.1823.32.idl", "6bd6b0cb996f5794b3dacd367e07d5ae9fb0832d9e0a96df129f9b5b6993b670", "6e5c69f7a1ecbd6925deb4aae2ae9219e925ae3c1f0f531d29a9db692c1cba4d"},
		{"WebView2.1.0.1901.177.idl", "b02844a66b8f7f2fb883467ecb9cf309891e07af1aae25d716444bb39f6fdf11", "cdbc9dd90676eb8d852e663e72495b3b47f1e659057dfb4c62b0b81430dd6a7e"},
		{"WebView2.1.0.2045.28.idl", "aa7827f4746fb66b4cf4eb7c30682bf9bb3f15aa08ec42e0664526f3483651d4", "1e0e2cd9329c4d7adaa92953d6f95ea5de925ecd341921dfe3617d2b6f995390"},
		{"WebView2.1.0.2592.51.idl", "6672cfe82773978f52b6807b069e2df02efadd9acc9e0cfdc291aadf04086cea", "0d0ed6cfabef4cb374a8b57f50b78fb3e7c2bd9f1362b896f2772e927ceb5ac1"},
		{"WebView2.1.0.2739.15.idl", "f304076fd92b8e8e83cd82d5c1870bc5f91d1a1ac5e2716fa0b4b99ba5e0f7a4", "19cb0fbf7a4ba92e451668d009dad0e1bf28a6f38ce840401bd9234c608a91bb"},
		{"WebView2.1.0.2903.40.idl", "56ee0e044bf47f6fcd0ecca0cd4ef28511a85009151f15dc2ad2a48f43a77443", "b826f8ecc84210ebf561a353e39ee4aa65b4a4250ad5215acfe1e38238363476"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			data, err := os.ReadFile("../" + tt.name)
			if err != nil {
				t.Fatal(err)
			}
			tree, err := idlreader.Parse(data)
			if err != nil {
				t.Fatal(err)
			}
			model, err := json.Marshal(tree)
			if err != nil {
				t.Fatal(err)
			}
			if got := fmt.Sprintf("%x", sha256.Sum256(model)); got != tt.model {
				t.Fatalf("parsed model changed: %s", got)
			}
			files, err := ParseIDL(data)
			if err != nil {
				t.Fatal(err)
			}
			outputs := make(map[string]string, len(files))
			for _, f := range files {
				outputs[f.FileName] = f.Package + "\n" + f.Content.String()
			}
			output, err := json.Marshal(outputs)
			if err != nil {
				t.Fatal(err)
			}
			if got := fmt.Sprintf("%x", sha256.Sum256(output)); got != tt.output {
				t.Fatalf("generated output changed: %s", got)
			}
		})
	}
}
