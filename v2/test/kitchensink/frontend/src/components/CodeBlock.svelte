
<script>
    import { Highlight } from "svelte-highlight";
    import { go, javascript } from "svelte-highlight/languages";
    let isJs = false;
    $: lang = isJs ? javascript : go;
    export let jsCode = "Hi form JS!";
    export let goCode = "Hi form Go!";

    $: code = isJs ? jsCode : goCode;
</script>

<div data-wails-no-drag class="code">
    <div class="header">
        <span>Title</span>
        <span class="toggle">
            <span>Go</span>
            <span class="custom-switch">
                <input type="checkbox" id="languageToggle" value="" bind:checked={isJs}>
                <label for="languageToggle">Javascript</label>
            </span>
        </span>
    </div>
    <Highlight language="{lang}" {code} />
</div>

<style>

    .header {
        margin: 5px;
    }

    .toggle {
        float: right;
    }

    .custom-switch {
        display: inline-block;
    }

</style>