
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

    export let id;
    
    // Handle hiding example
    let showRun = true;

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
        <div class="description">{description}</div>
        {/if}
    </div>
    <div class="example"> 
        <div class="{showCode ? 'code-title-open' : 'code-title-closed'}" on:click="{toggleExample}" >
            <span style="display: inline-block; width: 15px;">{showCode?'▼':'▶'}</span>
            Example Code
        </div>
        {#if showCode}
        <Highlight class="allow-select" language="{lang}" {code} style="max-height: 500px;"/>
        {/if}
    </div>

    <div class="run"> 
        <div class="{showRun ? 'run-title-open' : 'run-title-closed'}" on:click="{toggleRun}">
            <span style="display: inline-block; width: 15px;">{showRun?'▼':'▶'}</span>
            Try Me!
        </div>
        {#if showRun}
        <div class={$darkMode ? "run-content-dark" : "run-content-light"}>
            <slot></slot>
        </div>
        {/if}
    </div>
</div>

<style>

    .header {
        display: flex;
        justify-content: space-between;
        border-bottom: 1px solid #5555;
        flex-wrap: wrap;
        padding-bottom: 10px;
    }

    .title {
        font-size: 1.75rem;
    }

    .code-title-open {
        padding-top: 15px;
        padding-bottom: 0;
        cursor: pointer;
    }

    .code-title-closed {
        padding-top: 15px;
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
        margin-top: 15px;
        margin-bottom: 15px;
        padding-bottom: 0;
        cursor: pointer;
    }

    .run-title-closed {
        margin-top: 15px;
        margin-bottom: 15px;
        cursor: pointer;
    }

    .toggle {
        float: right;
        margin-top: 2px;
    }

    .example {
        padding-right: 5px;
        border-bottom: 1px solid #5555;
        padding-bottom: 10px;
    }
    
    .custom-switch {
        display: inline-block;
        margin-left: 5px;
    }

    .codeblock {
        background-color: #3F3F4B;
        border-radius: 5px;
        border: 1px solid #555;
        padding: 15px;
    }

    .codeblock-light {
        background-color: #e5e5e5;
        border-radius: 5px;
        border: 1px solid #ccc;
        padding: 15px;
    }

    .description {
        border-top: 1px solid #5555;
        margin-top: 10px;
        padding-top: 5px;
        width: 100%;
    }

</style>