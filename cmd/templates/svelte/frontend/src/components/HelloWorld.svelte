<script>
	import { getContext } from 'svelte';
	import ModalContent from './ModalContent.svelte'
	
	const { open } = getContext('simple-modal');

	const handleOpenModal = () => {
		window.backend.basic().then((result) => {
			open(ModalContent, { message: result });
		});
	};
</script>

<main>
	<p><button on:click={handleOpenModal}>Hello</button></p>
</main>

<style></style>