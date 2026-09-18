//go:build server

package application

// streamPrelude is prepended to the runtime bundle so the stream transport is
// chosen before any module body runs.
//
// custom.js cannot do this job. It is injected with loadOptionalScript, which
// does a HEAD request and then appends a <script> tag, so it lands well after
// the runtime's dependents have evaluated. A stream created at module scope —
// which is exactly the shape generated bindings use:
//
//	export const Telemetry = Stream("telemetry");
//
// would therefore connect before the factory existed and silently take the poll
// transport. Prepending to the bundle is synchronous by construction: ES module
// dependencies evaluate before their importers, so this runs before any
// generated module can call Stream.
func streamPrelude() []byte {
	return []byte(`(function(){
	window._wails = window._wails || {};
	window._wails.streamFactory = function(name) {
		var p = location.protocol === 'https:' ? 'wss:' : 'ws:';
		var sock = new WebSocket(p + '//' + location.host + '/wails/stream/ws?name=' + encodeURIComponent(name));
		sock.binaryType = 'arraybuffer';
		return sock;
	};
})();
`)
}
