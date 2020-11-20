
<script>
    import { darkMode } from '../Store';

    import { Highlight } from "svelte-highlight";
    import { go, javascript } from "svelte-highlight/languages";
    
    // Default to Go
    export let isJs = false;

    export let title;
    export let description;
    
    // Calculate CSS to use
    $: lang = isJs ? javascript : go;

    // Calculate Code for code block
    export let jsCode = "Hi from JS!";
    export let goCode = "Hi from Go!";
    $: code = isJs ? jsCode : goCode;

    // Handle hiding example
    let showCode = false;

    function toggleExample() {
        showCode = !showCode;
    }   

    export let id = "toggle-" + Date.now().toString() + Math.random().toString();
    
    // Handle hiding example
    export let showRun = false;

    function toggleRun() {
        showRun = !showRun;
    }
</script>

<div data-wails-no-drag class={$darkMode ? "codeblock" : "codeblock-light"}>
    <div class="header">
        <span class="title">{title}</span>
        <span class="toggle">
            <span>Go</span>
            <span class="custom-switch">
                <input type="checkbox" {id} value="" bind:checked={isJs}>
                <label for={id}>Javascript</label>
            </span>
        </span>
        {#if description}
        <div class="description">{@html description}</div>
        {/if}
    </div>
    <div class="run"> 
        <div class="{showRun ? 'run-title-open' : 'run-title-closed'}" on:click="{toggleRun}">
            <span class="arrow">{showRun?'▼':'▶'}</span>
            Try Me!
        </div>
        {#if showRun}
        <div class={$darkMode ? "run-content-dark" : "run-content-light"}>
            <slot></slot>
        </div>
        {/if}
    </div>
    <div class="example allow-select"> 
        <div class="{showCode ? 'code-title-open' : 'code-title-closed'}" on:click="{toggleExample}" >
            <span class="arrow">{showCode?'▼':'▶'}</span>
            Example Code
        </div>
        {#if showCode}
        <Highlight style="margin-bottom: 0" language="{lang}" {code}/>
        {/if}
    </div>
</div>

<style>

    .arrow {
        display: inline-block; 
        width: calc(var(--base-font-size));
        padding: 2px;
    }

    .header {
        display: flex;
        justify-content: space-between;
        border-bottom: 1px solid #5555;
        flex-wrap: wrap;
        align-items: center;
        padding-bottom: 5px;
        padding-left: 5px;
        padding-right: 5px;
    }

    .title {
        font-size: calc(var(--base-font-size) * 1.1);
    }

    .code-title-open {
        margin-top: 5px;
        margin-bottom: -5px;
        padding-left: 5px;
        cursor: pointer;
    }

    .code-title-closed {
        margin-top: 5px;
        padding-left: 5px;
        padding-bottom: 5px;
        cursor: pointer;
    }

    .run-content-dark {
        padding: 15px;
        background-color: #282c34;
    }

    .run-content-light {
        padding: 15px;
        background-color: #fafafa;
    }

    .run-title-open {
        margin-top: 5px;
        margin-bottom: 5px;
        padding-bottom: 0;
        padding-left: 5px;
        cursor: pointer;
    }

    .run-title-closed {
        margin-top: 5px;
        margin-bottom: 5px;
        padding-left: 5px;
        cursor: pointer;
    }

    .toggle {
        float: right;
        margin-top: 2px;
        font-size: calc(var(--base-font-size) * 0.9);
    }

    .example {
        border-top: 1px solid #5555;
    }
    
    .custom-switch {
        display: inline-block;
        margin-left: 5px;
    }

    .codeblock {
        /* background-color: #3F3F4B; */
        border-radius: 5px;
        border: 1px solid #555;
        padding: 5px;
        margin-top: 20px;
        margin-bottom: 10px;
    }

    .codeblock-light {
        /* background-color: #e5e5e5; */
        border-radius: 5px;
        border: 1px solid #ccc;
        padding: 5px;
        margin-top: 20px;
        margin-bottom: 10px;
    }

    .description {
        border-top: 1px solid #5555;
        margin-top: 10px;
        padding-top: 5px;
        width: 100%;
    }

</style>