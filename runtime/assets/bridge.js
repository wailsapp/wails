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
			'<div class="wails-reconnect-overlay"><div class="wails-reconnect-overlay-content"><div class="wails-reconnect-overlay-title">Wails Bridge</div><br><div class="wails-reconnect-overlay-loadingspinner"></div><br><div id="wails-reconnect-overlay-message">Waiting for backend</div></div></div>',
		overlayCSS:
			'.wails-reconnect-overlay{position:fixed;top:0;left:0;width:100%;height:100%;background:rgba(0,0,0,.6);font-family:sans-serif;display:none;z-index:999999}.wails-reconnect-overlay-content{padding:20px 30px;text-align:center;width:20em;position:relative;height:14em;border-radius:1em;margin:5% auto 0;background-color:#fff;box-shadow:1px 1px 20px 3px;background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAC8AAAAuCAMAAACPpbA7AAAAqFBMVEUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBAQAAAAAAAAEBAQAAAAAAAAAAAAEBAQEBAQDAwMBAQEAAAABAQEAAAAAAAAAAAABAQEAAAAAAAACAgICAgIBAQEAAAAAAAAAAAAAAAAAAAAAAAABAQEAAAACAgIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFBQWKCj6oAAAAN3RSTlMALiIqDhkGBAswJjP0GxP6NR4W9/ztjRDMhWU50G9g5eHXvbZ9XEI9xZTcqZl2aldKo55QwoCvZUgzhAAAAs9JREFUSMeNleeWqjAUhU0BCaH3Itiw9zKT93+zG02QK1hm/5HF+jzZJ6fQe6cyXE+jg9X7o9wxuylIIf4Tv2V3+bOrEXnf8dwQ/KQIGDN2/S+4OmVCVXL/ScBnfibxURqIByP/hONE8r8T+bDMlQ98KSl7Y8hzjpS8v1qtDh8u5f8KQpGpfnPPhqG8JeogN37Hq9eaN2xRhIwAaGnvws8F1ShxqK5ob2twYi1FAMD4rXsYtnC/JEiRbl4cUrCWhnMCLRFemXezXbb59QK4WASOsm6n2W1+4CBT2JmtzQ6fsrbGubR/NFbd2g5Y179+5w/GEHaKsHjYCet7CgrXU3txarNC7YxOVJtIj4/ERzMdZfzc31hp+8cD6eGILgarZY9uZ12hAs03vfBD9C171gS5Omz7OcvxALQIn4u8RRBBBcsi9WW2woO9ipLgfzpYlggg3ZRdROUC8KT7QLqq3W9KB5BbdFVg4929kdwp6+qaZnMCCNBdj+NyN1W885Ry/AL3D4AQbsVV4noCiM/C83kyYq80XlDAYQtralOiDzoRAHlotWl8q2tjvYlOgcg1A8jEApZa+C06TBdAz2Qv0wu11I/zZOyJQ6EwGez2P2b8PIQr1hwwnAZsAxwA4UAYOyXUxM/xp6tHAn4GUmPGM9R28oVxgC0e/zQJJI6DyhyZ1r7uzRQhpcW7x7vTaWSzKSG6aep77kroTEl3U81uSVaUTtgEINfC8epx+Q4F9SpplHG84Ek6m4RAq9/TLkOBrxyeuddZhHvGIp1XXfFy3Z3vtwNblKGiDn+J+92vwwABHghj7HnzlS1H5kB49AZvdGCFgiBPq69qfXPr3y++yilF0ON4R8eR7spAsLpZ95NqAW5tab1c4vkZm6aleajchMwYTdILQQTwE2OV411ZM9WztDjPql12caBi6gDpUKmDd4U1XNdQxZ4LIXQ5/Tr4P7I9tYcFrDK3AAAAAElFTkSuQmCC);background-repeat:no-repeat;background-position:center}.wails-reconnect-overlay-title{font-size:2em}.wails-reconnect-overlay-message{font-size:1.3em}.wails-reconnect-overlay-loadingspinner{pointer-events:none;width:2.5em;height:2.5em;border:.4em solid transparent;border-color:#3E67EC #eee #eee;border-radius:50%;animation:loadingspin 1s linear infinite;margin:auto;padding:2.5em}@keyframes loadingspin{100%{transform:rotate(360deg)}}',
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
	var elem = document.createElement('style');
	elem.setAttribute('type', 'text/css');
	if (elem.styleSheet) {
		elem.styleSheet.cssText = css;
	} else {
		elem.appendChild(document.createTextNode(css));
	}
	var head = document.head || document.getElementsByTagName('head')[0];
	head.appendChild(elem);
}

// Creates a node in the Dom
function createNode(parent, elementType, id, className, content) {
	var d = document.createElement(elementType);
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
	var body = document.body;
	var wailsBridgeNode = createNode(body, 'div', 'wails-bridge');
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
	function hideReconnectOverlay() {
		window.wailsbridge.reconnectOverlay.style.display = 'none';
	}

	// Adds a script to the Dom.
	// Removes it if second parameter is true.
	function addScript(script, remove) {
		var s = document.createElement('script');
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
			var binding = message.data.slice(1);
			//log("Binding: " + binding)
			window.wails._.NewBinding(binding);
			break;
			// Call back
		case 'c':
			var callbackData = message.data.slice(1);
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

function start(callback) {

	// Set up the bridge
	init();

	// Save the callback
	window.wailsbridge.callback = callback;

	// Start Bridge
	startBridge();
}

function Init(callback) {
	start(callback);
}

module.exports = Init;
