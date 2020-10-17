<script>
import { afterUpdate } from 'svelte';

    import { darkMode } from '../Store';

    $: termClass = $darkMode ? 'faketerm-dark' : 'faketerm-light';

    let termElement;

    afterUpdate( () => {
        termElement.scrollTop = termElement.scrollHeight;
    });

    export let text = "";
    export let style = null;
</script>

<div bind:this={termElement} class="common {termClass}" {style}>
<pre>
{#if text && text.length > 0}
{text}
{:else}
<slot></slot>
{/if}
</pre>
</div>

<style>
    pre {
        margin: 0;
        padding: 5px;
    }

    .common {
        font-family: 'Courier New', Courier, monospace;
        padding: 5px;
        white-space: pre-line;
        margin-top: 10px;
        margin-bottom: 10px;
        border: 1px solid #5555;
        overflow-y: auto;
        overflow-wrap: break-word;

    }

    .faketerm-dark {
        background-color: black;
        color: white;
    }

    .faketerm-light {
        background-color: #ddd;
        color: black;
    }
</style>