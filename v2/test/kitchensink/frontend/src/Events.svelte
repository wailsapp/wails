<script>

  import runtime from '@wailsapp/runtime2';
  import { Highlight } from "svelte-highlight";
  import { typescript } from "svelte-highlight/languages";
  import { atomOneDark, atomOneLight } from "svelte-highlight/styles";

  let css = runtime.System.DarkModeEnabled() ? atomOneDark : atomOneLight;
  $: code = `const add = (a: number, b: number) => a + b;`;

  runtime.System.OnThemeChange( (isDarkMode) => {
    css = isDarkMode ? atomOneDark : atomOneLight;
  });
</script>

<svelte:head>
  {@html css}
</svelte:head>

<div class="page">
    <h4>Events</h4>

<Highlight language="{typescript}" {code} />
    Logging is part of the Wails Runtime and is accessed through the <code>runtime.Log</code> object. There are 5 methods available:
    
    <ul class="list">
        <li>Debug</li>
        <li>Info</li>
        <li>Warning</li>
        <li>Error</li>
        <li>Fatal</li>
    </ul>
    All methods will log to the console and <code>Fatal</code> will also exit the program.

    <h5>Try Me</h5>
    <hr>
    <div class="logging-form">
    </div>
</div>

<style>
    .page {
        margin-left: 50px;
        margin-top: 10px;
    }

    h5 {
        margin-top: 3rem;
    }
</style>