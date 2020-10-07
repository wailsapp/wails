
<script>
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

    // Handle hiding example
    let hidden = false;

    function toggleExample() {
        hidden = !hidden;
    }
</script>

<div data-wails-no-drag class="codeblock">
    <div class="header">
        <span on:click="{toggleExample}">Title</span>
        <span class="toggle">
            <span>Go</span>
            <span class="custom-switch">
                <input type="checkbox" id="languageToggle" value="" bind:checked={isJs}>
                <label for="languageToggle">Javascript</label>
            </span>
        </span>
    </div>
    {#if !hidden}
    <Highlight language="{lang}" {code} />
    {/if}
</div>

<style>

    .header {
        display: flex;
        justify-content: space-between;
        padding: 15px 15px 0 15px;
    }

    .toggle {
        float: right;
    }

    .custom-switch {
        display: inline-block;
        margin-left: 5px;
    }

    .codeblock {
        background-color: #3F3F4B;
        border-radius: 5px;
        padding-bottom: 15px;
    }

</style>