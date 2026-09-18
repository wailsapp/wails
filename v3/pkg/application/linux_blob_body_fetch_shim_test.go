//go:build linux && cgo && !android

package application

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestLinuxBlobBodyFetchShimIsScopedAndConvertsUnsupportedBodies(t *testing.T) {
	required := []string{
		`protocol === "wails:"`,
		`body instanceof Blob`,
		`body instanceof FormData`,
		`await body.arrayBuffer()`,
		`await encoded.arrayBuffer()`,
		`encoded.headers.get("Content-Type")`,
		`request.clone().arrayBuffer()`,
		`Object.assign({}, init,`,
		`new Request(request,`,
		`originalFetch.apply(this, arguments)`,
	}
	for _, fragment := range required {
		if !strings.Contains(linuxBlobBodyFetchShimJS, fragment) {
			t.Errorf("Linux Blob body fetch shim is missing %q", fragment)
		}
	}
}

func TestLinuxBlobBodyFetchShimHasIdempotencyGuard(t *testing.T) {
	const marker = "__wailsLinuxBlobBodyFetchShim"
	if count := strings.Count(linuxBlobBodyFetchShimJS, marker); count < 2 {
		t.Fatalf("expected the shim marker to be checked and set, got %d occurrences", count)
	}
}

func TestLinuxBlobBodyFetchShimBehaviour(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required for executable JavaScript shim tests")
	}

	harness := fmt.Sprintf(`
	import assert from "node:assert/strict";
globalThis.location = {href: "wails://wails/"};
let calls = [];
globalThis.fetch = async function () {
  calls.push(Array.from(arguments));
  return new Response("ok");
};
const originalFetch = globalThis.fetch;
const install = () => eval(%s);
install();
const wrappedFetch = globalThis.fetch;

await wrappedFetch("wails://wails/blob", {
  method: "POST",
  body: new Blob(["blob bytes"], {type: "text/plain"})
});
assert(calls[0][1].body instanceof ArrayBuffer);
assert.equal(new TextDecoder().decode(calls[0][1].body), "blob bytes");
assert.equal(calls[0][1].headers.get("Content-Type"), "text/plain");

const form = new FormData();
form.append("message", "form bytes");
await wrappedFetch("wails://wails/form", {method: "POST", body: form});
assert(calls[1][1].body instanceof ArrayBuffer);
assert.match(new TextDecoder().decode(calls[1][1].body), /form bytes/);
assert.match(calls[1][1].headers.get("Content-Type"), /^multipart\/form-data; boundary=/);

const externalBlob = new Blob(["external"]);
await wrappedFetch("https://example.com/", {method: "POST", body: externalBlob});
assert.strictEqual(calls[2][1].body, externalBlob);

await wrappedFetch("wails://wails/string", {method: "POST", body: "unchanged"});
assert.equal(calls[3][1].body, "unchanged");

const request = new Request("wails://wails/request", {
  method: "POST",
  headers: {Authorization: "Bearer secret", "X-Original": "preserved"},
  body: new Blob(["request bytes"], {type: "text/plain"})
});
await wrappedFetch(request);
assert(calls[4][0] instanceof Request);
assert.notStrictEqual(calls[4][0], request);
assert.equal(await calls[4][0].text(), "request bytes");
assert.equal(calls[4][0].headers.get("Authorization"), "Bearer secret");
assert.equal(request.bodyUsed, false);

await wrappedFetch(request, {
  body: new Blob(["replacement"], {type: "application/octet-stream"})
});
assert.equal(calls[5][1].headers.get("Authorization"), "Bearer secret");
assert.equal(calls[5][1].headers.get("X-Original"), "preserved");
assert.equal(calls[5][1].headers.get("Content-Type"), "text/plain");

await wrappedFetch(request, {
  body: new Blob(["replacement"], {type: "application/octet-stream"}),
  headers: {"X-Override": "authoritative"}
});
assert.equal(calls[6][1].headers.get("Authorization"), null);
assert.equal(calls[6][1].headers.get("X-Original"), null);
assert.equal(calls[6][1].headers.get("X-Override"), "authoritative");
assert.equal(calls[6][1].headers.get("Content-Type"), "application/octet-stream");

install();
assert.strictEqual(globalThis.fetch, wrappedFetch);
assert.notStrictEqual(globalThis.fetch, originalFetch);
`, strconv.Quote(linuxBlobBodyFetchShimJS))

	command := exec.Command(node, "--input-type=module", "--eval", harness)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("executable JavaScript shim test failed: %v\n%s", err, output)
	}
}
