<script>
    import { Tray } from '@wails/runtime';
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import CodeSnippet from '../../../components/CodeSnippet.svelte';
    import description from './description.txt';
    import { UniqueID } from '../../../utils/utils';
    import jsCode from './code.jsx';
    import goCode from './code.go';
    import { darkMode } from '../../../Store';

    let isJs = false;
    $: lang = isJs ? 'Javascript' : 'Go';

    let id = UniqueID('tray');

    var icons = ["light", "dark", "svelte"];
    let darkmode = $darkmode;
    let iconName = darkMode ? 'light' : 'dark';

    function setIcon() {
    	console.log(iconName);
        if( isJs ) {
            Tray.SetIcon(iconName);
        } else {
            backend.main.Tray.SetIcon(iconName);
        }
    }

    $: exampleCodeJS = `import { Tray } from '@wails/runtime';\n\nTray.SetIcon('` + iconName + `');`;
    $: exampleCodeGo = `// runtime is given through WailsInit()\nruntime.Tray.SetIcon("` + iconName + `");`;
    

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} {id} title="Tray.SetIcon(trayIconID)" {description}>
    <div class="logging-form">
        <form data-wails-no-drag class="mw-full">
            <!-- Radio -->
            <div class="form-group">
                <div>Select Tray Icon</div>
                {#each icons as option}
                    <div class="custom-radio">
                        <input type="radio" name="trayicon" bind:group="{iconName}" id="{id}-{option}" value="{option}">
                        <label for="{id}-{option}">{option}</label>
                    </div>
                {/each}
            </div>

            <input class="btn btn-primary" type="button" on:click="{setIcon}" value="Set Icon using {lang} runtime">

            <CodeSnippet bind:isJs={isJs} jsCode={exampleCodeJS} goCode={exampleCodeGo}></CodeSnippet>

        </form>
    </div>
</CodeBlock>
