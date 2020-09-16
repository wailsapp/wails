import App from './App.svelte';

const app = new App({
	target: document.body,
	props: {
		name: 'Wails User'
	}
});

export default app;
