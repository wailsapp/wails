<script>
	let name = "";
	let greeting = "";

	function greet() {
		window.go.main.App.Greet(name).then((result) => {
			greeting = result;
		});
	}
</script>

<main>
	<div id="logo"></div>
	<div id="input" data-wails-no-drag>
		<input id="name" type="text" bind:value={name}>
		<button class="button" on:click={greet}>Greet</button>
	</div>
	{#if greeting}
		<div id="result">{greeting}</div>
	{/if}
</main>

<style>

	main {
		height: 100%;
		width: 100%;
	}

	#result {
		margin-top: 1rem;
		font-size: 1.5rem;
	}

	button {
		-webkit-appearance: default-button;
		padding: 6px;
	}

	#name {
		border-radius: 3px;
		outline: none;
		-webkit-font-smoothing: antialiased;
	}

	#logo {
		width: 40%;
		height: 40%;
		padding-top: 20%;
		margin: auto;
		display: block;
		background-position: center;
		background-repeat: no-repeat;
		background-image: url("assets/images/logo-dark.svg");
	}

</style>