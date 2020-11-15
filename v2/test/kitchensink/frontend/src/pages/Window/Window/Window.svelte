<script>
    import { Window } from '@wails/runtime';
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import CodeSnippet from '../../../components/CodeSnippet.svelte';
    import jsCode from './code.jsx';
    import goCode from './code.go';

    var message = '';
    var isJs = false;
    var title = '';

    let windowActions = ["SetTitle", "Fullscreen", "UnFullscreen", "Maximise", "Unmaximise", "Minimise", "Unminimise", "Center", "Show", "Hide", "SetSize", "SetPosition", "Close"]
    var disabledActions = ['Show', 'Unminimise'];
    var windowAction = windowActions[0];

    $: windowRuntime = lang == 'Javascript' ? Window : backend.main.Window;
    $: lang = isJs ? 'Javascript' : 'Go';

    var id = "Window";

    function processAction() {

        switch( windowAction ) {
            case 'SetSize':
                windowRuntime.SetSize(sizeWidth, sizeHeight);
                break;
            case 'SetPosition':
                windowRuntime.SetPosition(positionX, positionY);
                break;
            case 'SetTitle':
                windowRuntime.SetTitle(title);
                break;
            case 'Hide':
                windowRuntime.Hide();
                setTimeout( windowRuntime.Show, 3000 );
            case 'Minimise':
                windowRuntime.Hide();
                setTimeout( windowRuntime.Unminimise, 3000 );
            default:
                windowRuntime[windowAction]();
        }
    }

    var params = "";
    var sizeWidth = 1024;
    var sizeHeight = 768;
    var positionX = 100;
    var positionY = 100;

    $: {
        switch (windowAction) {
            case 'SetSize':
                params = sizeWidth + ", " + sizeHeight;
                break;
            case 'SetPosition':
                params = positionX + ", " + positionY;
                break;
            case 'SetTitle':
                params = `'` + title.replace(`"`, `\"`) + `'`;
                break;
            default:
                params = '';
                break;
                
        }
        
    }
    $: testcodeJs = "import { Window } from '@wails/runtime';\nWindow." + windowAction + "(" + params + ");";
    $: testcodeGo = '// runtime is given through WailsInit()\nruntime.Window.' + windowAction + '(' + params + ')'; 

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} title="Window" {id} showRun=true>
    <div class="logging-form">
        <form data-wails-no-drag class="mw-full"> 
            <div class="form-group">
                <div>Select Window Method</div>
                {#each windowActions as option, index}
                <div class="custom-radio">
                    <input type="radio" name="window" bind:group="{windowAction}" id="{id}-{option}" value="{option}" disabled={disabledActions.includes(option)}>
                    <label for="{id}-{option}">{option}
                        {#if option == 'Hide' } - Show() will be called after 3 seconds {/if}
                        {#if option == 'Minimise' } - Unminimise() will be called after 3 seconds {/if}
                    </label>

                    {#if option == "SetSize"}
                    {#if windowAction == "SetSize" }
                    <div class="form-inline form-group numberInputGroup">Width: <input type="number" class="form-control numberInput" bind:value={sizeWidth}> Height: <input type="number" class="form-control numberInput" bind:value={sizeHeight}></div>
                    {/if}
                    {/if}

                    {#if option == "SetPosition"}
                    {#if windowAction == "SetPosition" }
                    <div class="form-inline form-group numberInputGroup">X: <input type="number" class="form-control numberInput" bind:value={positionX}> Y: <input type="number" class="form-control numberInput" bind:value={positionY}></div>
                    {/if}
                    {/if}

                    {#if option == "SetTitle"}
                    {#if windowAction == "SetTitle" }
                    <div class="form-inline form-group numberInputGroup">Title: <input type="text" class="form-control" bind:value={title}></div>
                    {/if}
                    {/if}

                </div>   
                {/each}
            </div>

            <input class="btn btn-primary" type="button" on:click="{processAction}" value="Call using {lang} runtime">

            <CodeSnippet bind:isJs={isJs} jsCode={testcodeJs} goCode={testcodeGo}></CodeSnippet>

        </form>
    </div>
</CodeBlock>

<style>
    .numberInputGroup {
        margin: 10px 25px 10px;
        width: 40%;
        min-width: 350px;
    }

    .numberInput {
        margin-left: 10px; 
        width: 80px;
    }
</style>