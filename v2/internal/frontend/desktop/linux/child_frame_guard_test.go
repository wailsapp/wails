//go:build linux

package linux

import (
	"os/exec"
	"strconv"
	"testing"
)

func TestChildFrameBridgeAdmission(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required for bridge JavaScript tests")
	}
	script := `
const vm = require('vm');
const assert = require('assert');
const guard = ` + strconv.Quote(childFrameGuardJS) + `;
for (const mainFrame of [false, true]) {
 const messages = [];
 const window = {webkit: {messageHandlers: {external: {postMessage(m) {messages.push(m)}}}}};
 window.window = window;
 window.top = mainFrame ? window : {};
 vm.runInNewContext(guard, window);
 for (const message of ['C{"name":"main.App.Witness","args":[]}', 'Q', 'BO:https://example.com', 'EE{}']) {
  window.webkit.messageHandlers.external.postMessage(message);
 }
 assert.equal(messages.length, mainFrame ? 4 : 0);
 if (!mainFrame) {
  assert.equal(Reflect.deleteProperty(window, 'webkit'), false);
  assert.throws(() => Object.defineProperty(window, 'webkit', {value: {}}));
 }
}
`
	if output, err := exec.Command(node, "-e", script).CombinedOutput(); err != nil {
		t.Fatalf("bridge admission: %v\n%s", err, output)
	}
}
