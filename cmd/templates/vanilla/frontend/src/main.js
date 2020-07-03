import 'core-js/stable';
const runtime = require('@wailsapp/runtime');

// Main entry point
function start() {

	// Ensure the default app div is 100% wide/high
	var app = document.getElementById('app');
	app.style.width = '100%';
	app.style.height = '100%';

	// Inject html
	app.innerHTML = `
	<div class='logo'></div>
	<div class='container'>
			<button id='button'>Click Me!</button>
			<div id='result'/>
	</div>
	`;

	// Connect button to Go method
	document.getElementById('button').onclick = function() {
		window.backend.basic().then( function(result) {
			document.getElementById('result').innerText = result;
		});
	};
};

// We provide our entrypoint as a callback for runtime.Init
runtime.Init(start);