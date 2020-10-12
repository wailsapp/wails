
<script>
    import { darkMode } from '../Store';

    import { Highlight } from "svelte-highlight";
    import { go, javascript } from "svelte-highlight/languages";
    
    // Default to Go
    export let isJs = false;

    // Calculate CSS to use
    $: lang = isJs ? javascript : go;

    // Calculate Code for code block
    export let jsCode = "Hi from JS!";
    export let goCode = "Hi from Go!";
    $: code = isJs ? jsCode : goCode;

</script>

<div data-wails-no-drag class={$darkMode ? "codeblock" : "codeblock-light"}>
    <Highlight language="{lang}" {code} style="margin: -5px;"/>
</div>

<style>

    .codeblock {
        background-color: #3F3F4B;
        border: 1px solid #555;
        padding: 5px;
        margin-top: 20px;
        margin-bottom: 10px;
    }

    .codeblock-light {
        background-color: #e5e5e5;
        border: 1px solid #ccc;
        padding: 5px;
        margin-top: 20px;
        margin-bottom: 10px;
    }

</style>