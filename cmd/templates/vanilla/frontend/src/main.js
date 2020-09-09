import 'core-js/stable';
const runtime = require('@wailsapp/runtime');

// Main entry point
function start() {

	var mystore = runtime.Store.New('Counter');

	// Ensure the default app div is 100% wide/high
	var app = document.getElementById('app');
	app.style.width = '100%';
	app.style.height = '100%';

	// Inject html
	app.innerHTML = `
	<div class='logo'></div>
	<div class='container'>
	<button onClick='window.backend.Counter.Increment()'>
		Increment Counter
	</button>
	<button onClick='window.backend.Counter.Decrement()'>
		Decrement Counter
		</button>
	</div>
	<div class='result'>Counter: <span id='counter'></span></div>
	<div class='container'>
		<input id='newCounter' type="number" value="0"/>
		<button id='setvalue'>Set Counter Value</button>
		<button onclick='window.backend.Counter.RandomValue()'>Set to Random Value</button>
	</div>
	`;

	// Connect counter value button to Go method
	document.getElementById('setvalue').onclick = function() {
		let newValue = parseInt(document.getElementById('newCounter').value,10);
		mystore.set(newValue);
	};

	mystore.subscribe( function(state) {
		document.getElementById('counter').innerText = state;
	});
	
	mystore.set(0);
};

// We provide our entrypoint as a callback for runtime.Init
runtime.Init(start);