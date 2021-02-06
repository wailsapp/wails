/*
 _       __      _ __
| |     / /___ _(_) /____
| | /| / / __ `/ / / ___/
| |/ |/ / /_/ / / (__  )
|__/|__/\__,_/_/_/____/
The lightweight framework for web-like apps
(c) Lea Anthony 2019-present
*/
/* jshint esversion: 6 */

function init() {
	// Bridge object
	window.wailsbridge = {
		reconnectOverlay: null,
		reconnectTimer: 300,
		wsURL: 'ws://' + window.location.hostname + ':34115/bridge',
		connectionState: null,
		config: {},
		websocket: null,
		callback: null,
		overlayHTML:
			'<div class="wails-reconnect-overlay"><div class="wails-reconnect-overlay-content"><div class="wails-reconnect-overlay-loadingspinner"></div></div></div>',
		overlayCSS:
			'.wails-reconnect-overlay{position:fixed;top:0;left:0;width:100%;height:100%;background:linear-gradient(45deg,rgb(0 0 0 / 60%) 0%,rgb(4 0 255 / 60%) 100%);font-family:sans-serif;display:none;z-index:999999}.wails-reconnect-overlay-content{text-align:center;width:10em;position:relative;margin:5% auto 0;background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAC8AAAAuCAMAAACPpbA7AAAAflBMVEUAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBAQAAAAAAAAAAAABAQEEBAQAAAAAAAAEBAQAAAADAwMAAAABAQEAAAAAAAAAAAAAAAAAAAACAgICAgIBAQEAAAAAAAAAAAAAAAAAAAACAgIAAAAAAAAAAAAAAAAAAAAAAAAAAAAFBQWCC3waAAAAKXRSTlMALgUMIBk0+xEqJs70Xhb3lu3EjX2EZTlv5eHXvbarQj3cdmpXSqOeUDwaqNAAAAKCSURBVEjHjZTntqsgEIUPVVCwtxg1vfD+L3hHRe8K6snZf+KKn8OewvzsSSeXLruLnz+KHs0gr6DkT3xsRkU6VVn4Ha/UxLe1Z4y64i847sykPBh/AvQ7ry3eFN70oKrfcBJYvm/tQ1qxP4T3emXPeXAkvodPUvtdjbhk+Ft4c0hslTiXVOzxOJ15NWUblQhRsdu3E1AfCjj3Gdm18zSOsiH8Lk4TB480ksy62fiqNo4OpyU8O21l6+hyRtS6z8r1pHlmle5sR1/WXS6Mq2Nl+YeKt3vr+vdH/q4O68tzXuwkiZmngYb4R8Co1jh0+Ww2UTyWxBvtyxLO7QVjO3YOD/lWZpbXDGellFG2Mws58mMnjVZSn7p+XvZ6IF4nn02OJZV0aTO22arp/DgLPtrgpVoi6TPbZm4XQBjY159w02uO0BDdYsfrOEi0M2ulRXlCIPAOuN1NOVhi+riBR3dgwQplYsZRZJLXq23Mlo5njkbY0rZFu3oiNIYG2kqsbVz67OlNuZZIOlfxHDl0UpyRX86z/OYC/3qf1A1xTrMp/PWWM4ePzf8DDp1nesQRpcFk7BlwdzN08ZIALJpCaciQXO0f6k4dnuT/Ewg4l7qSTNzm2SykdHn6GJ12mWc6aCNj/g1cTXpB8YFfr0uVc96aFkkqiIiX4nO+salKwGtIkvfB+Ja8DxMeD3hIXP5mTOYPB4eVT0+32I5ykvPZjesnkGgIREgYnmLrPb0PdV3hoLup2TjcGBPM4mgsfF5BrawZR4/GpzYQzQfrUZCf0TCWYo2DqhdhTJBQ6j4xqmmLN5LjdRIY8LWExiFUsSrza/nmFBqw3I9tEZB9h0lIQSO9if8DkISDAj8CDawAAAAASUVORK5CYII=);background-repeat:no-repeat;background-position:center}.wails-reconnect-overlay-loadingspinner{pointer-events:none;width:2.5em;height:2.5em;border:.4em solid transparent;border-color:#fff #eee0 #fff #eee0;border-radius:50%;animation:loadingspin 1s linear infinite;margin:auto;padding:2.5em}@keyframes loadingspin{100%{transform:rotate(360deg); opacity}}',
		log: function (message) {
			// eslint-disable-next-line
			console.log(
				'%c wails bridge %c ' + message + ' ',
				'background: #aa0000; color: #fff; border-radius: 3px 0px 0px 3px; padding: 1px; font-size: 0.7rem',
				'background: #009900; color: #fff; border-radius: 0px 3px 3px 0px; padding: 1px; font-size: 0.7rem'
			);
		}
	};
}

// Adapted from webview - thanks zserge!
function injectCSS(css) {
	const elem = document.createElement('style');
	elem.setAttribute('type', 'text/css');
	elem.appendChild(document.createTextNode(css));
	const head = document.head || document.getElementsByTagName('head')[0];
	head.appendChild(elem);
}

// Creates a node in the Dom
/**
 * @param {HTMLElement} parent
 * @param {string} elementType
 * @param {string} id
 * @param {string|null} className
 * @param {string|null} content
 */
function createNode(parent, elementType, id, className, content) {
	const d = document.createElement(elementType);
	if (id) {
		d.id = id;
	}
	if (className) {
		d.className = className;
	}
	if (content) {
		d.innerHTML = content;
	}
	parent.appendChild(d);
	return d;
}

// Sets up the overlay
function setupOverlay() {
	const body = document.body;
	const wailsBridgeNode = createNode(body, 'div', 'wails-bridge', null, null);
	wailsBridgeNode.innerHTML = window.wailsbridge.overlayHTML;

	// Inject the overlay CSS
	injectCSS(window.wailsbridge.overlayCSS);
}

// Start the Wails Bridge
function startBridge() {
	// Setup the overlay
	setupOverlay();

	window.wailsbridge.websocket = null;
	window.wailsbridge.connectTimer = null;
	window.wailsbridge.reconnectOverlay = document.querySelector(
		'.wails-reconnect-overlay'
	);
	window.wailsbridge.connectionState = 'disconnected';

	// Shows the overlay
	function showReconnectOverlay() {
		window.wailsbridge.reconnectOverlay.style.display = 'block';
	}

	// Hides the overlay
	const hideReconnectOverlay = function () {
		window.wailsbridge.reconnectOverlay.style.display = 'none';
	}
	window.wailsbridge.hideReconnectOverlay = hideReconnectOverlay;

	// Adds a script to the Dom.
	// Removes it if second parameter is true.
	function addScript(script, remove) {
		const s = document.createElement('script');
		s.setAttribute('type', 'text/javascript');
		s.textContent = script;
		document.head.appendChild(s);

		// Remove internal messages from the DOM
		if (remove) {
			s.parentNode.removeChild(s);
		}
	}

	// Handles incoming websocket connections
	function handleConnect() {
		window.wailsbridge.log('Connected to backend');
		hideReconnectOverlay();
		clearInterval(window.wailsbridge.connectTimer);
		window.wailsbridge.websocket.onclose = handleDisconnect;
		window.wailsbridge.websocket.onmessage = handleMessage;
		window.wailsbridge.connectionState = 'connected';
	}

	// Handles websocket disconnects
	function handleDisconnect() {
		window.wailsbridge.log('Disconnected from backend');
		window.wailsbridge.websocket = null;
		window.wailsbridge.connectionState = 'disconnected';
		showReconnectOverlay();
		connect();
	}

	// Try to connect to the backend every 300ms (default value).
	// Change this value in the main wailsbridge object.
	function connect() {
		window.wailsbridge.connectTimer = setInterval(function () {
			if (window.wailsbridge.websocket == null) {
				window.wailsbridge.websocket = new WebSocket(window.wailsbridge.wsURL);
				window.wailsbridge.websocket.onopen = handleConnect;
				window.wailsbridge.websocket.onerror = function (e) {
					e.stopImmediatePropagation();
					e.stopPropagation();
					e.preventDefault();
					window.wailsbridge.websocket = null;
					return false;
				};
			}
		}, window.wailsbridge.reconnectTimer);
	}

	function handleMessage(message) {
		// As a bridge we ignore js and css injections
		switch (message.data[0]) {
		// Wails library - inject!
		case 'w':
			addScript(message.data.slice(1));

			// Now wails runtime is loaded, wails for the ready event
			// and callback to the main app
			window.wails.Events.On('wails:loaded', function () {
				window.wailsbridge.log('Wails Ready');
				if (window.wailsbridge.callback) {
					window.wailsbridge.log('Notifying application');
					window.wailsbridge.callback(window.wails);
				}
			});
			window.wailsbridge.log('Loaded Wails Runtime');
			break;
			// Notifications
		case 'n':
			addScript(message.data.slice(1), true);
			break;
			// Binding
		case 'b':
			const binding = message.data.slice(1);
			//log("Binding: " + binding)
			window.wails._.NewBinding(binding);
			break;
			// Call back
		case 'c':
			const callbackData = message.data.slice(1);
			window.wails._.Callback(callbackData);
			break;
		default:
			window.wails.Log.Error('Unknown message type received: ' + message.data[0]);
		}
	}

	// Start by showing the overlay...
	showReconnectOverlay();

	// ...and attempt to connect
	connect();
}

function InitBridge(callback) {
	// Set up the bridge
	init();

	// Save the callback
	window.wailsbridge.callback = callback;

	// Start Bridge
	startBridge();
}

module.exports = {
	InitBridge: InitBridge,
}
