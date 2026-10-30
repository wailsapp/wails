//go:build linux && cgo && !android

package application

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestLinuxChildFrameGuardBehaviour(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required for executable JavaScript shim tests")
	}
	script := `
const vm = require("vm");
const guard = ` + strconv.Quote(linuxChildFrameGuardJS) + `;
function frame(isTop) {
  const sent = [];
  const win = { webkit: { messageHandlers: { external: { postMessage(m) { sent.push(m); } } } } };
  win.top = isTop ? win : {};
  win.window = win;
  vm.runInNewContext(guard, win);
  return { win, sent };
}
const child = frame(false);
child.win.webkit.messageHandlers.external.postMessage("wails:event:emit:x");
child.win.webkit.messageHandlers.external.postMessage({ toString() { return "wails:drag"; } });
child.win.webkit.messageHandlers.external.postMessage("raw");
let flip = 0;
child.win.webkit.messageHandlers.external.postMessage({ toString() { return flip++ ? "wails:drag" : "ok"; } });
const tamper = (() => { try { child.win.webkit = {}; child.win.webkit.messageHandlers.external = null; } catch (e) {} return child.win.webkit.messageHandlers.external !== null; })();
const top = frame(true);
top.win.webkit.messageHandlers.external.postMessage("wails:drag");
console.log(JSON.stringify({ child: child.sent, top: top.sent, tamper }));
`
	out, err := exec.Command(node, "-e", script).CombinedOutput()
	if err != nil {
		t.Fatalf("node failed: %v\n%s", err, out)
	}
	got := strings.TrimSpace(string(out))
	want := `{"child":["raw","ok"],"top":["wails:drag"],"tamper":true}`
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}
