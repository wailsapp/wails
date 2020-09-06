

<script>
	import Modal from 'svelte-simple-modal';
	import HelloWorld from './components/HelloWorld.svelte'
	import logo from './logo.png';
</script>

<main>
	<div class="App">
		<header class="App-header">
			<Modal>
				<img src={logo} class="App-logo" alt="logo" />
				<p>Welcome to your new <code>wails/svelte</code> project.</p>
				<HelloWorld/>
			</Modal>
		</header>
	</div>
</main>

<style>
	:global(body) {
		margin: 0;
		font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Roboto", "Oxygen",
			"Ubuntu", "Cantarell", "Fira Sans", "Droid Sans", "Helvetica Neue",
			sans-serif;
		-webkit-font-smoothing: antialiased;
		-moz-osx-font-smoothing: grayscale;
	}

	code {
		font-family: source-code-pro, Menlo, Monaco, Consolas, "Courier New", monospace;
	}

	.App {
		text-align: center;
	}

	.App-logo {
		height: 40vmin;
		pointer-events: none;
	}

	@media (prefers-reduced-motion: no-preference) {
		.App-logo {
			animation: App-logo-spin infinite 20s linear;
		}
	}

	.App-header {
		background-color: #282c34;
		min-height: 100vh;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		font-size: calc(10px + 2vmin);
		color: white;
	}

	@keyframes App-logo-spin {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}
</style>