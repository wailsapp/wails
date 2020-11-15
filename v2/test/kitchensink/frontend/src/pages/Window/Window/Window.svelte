<script>
    import { Window } from '@wails/runtime';
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import CodeSnippet from '../../../components/CodeSnippet.svelte';
    import jsCode from './code.jsx';
    import goCode from './code.go';

    var message = '';
    var isJs = false;

    let windowActions = ["Maximise", "Unmaximise", "Minimise", "Unminimise", "Center", "Show", "Hide", "SetSize", "SetPosition", "Close"]

    var windowAction = windowActions[0];

    $: lang = isJs ? 'Javascript' : 'Go';

    var id = "Window";

    function processAction() {
        if ( lang == 'Javascript' ) {
            switch( windowAction ) {
                case 'SetSize':
                    Window.SetSize(sizeWidth, sizeHeight);
                    break;
                case 'SetPosition':
                    Window.SetPosition(positionX, positionY);
                    break;
                case 'Hide':
                    Window.Hide();
                    setTimeout( Window.Show, 3000 );
                case 'Minimise':
                    Window.Hide();
                    setTimeout( Window.Unminimise, 3000 );
                default:
                    Window[windowAction]();
            }
        } else {
            switch( windowAction ) {
                case 'SetSize':
                    backend.main.Window.SetSize(sizeWidth, sizeHeight);
                    break;
                case 'SetPosition':
                    backend.main.Window.SetPosition(positionX, positionY);
                    break;
                case 'Hide':
                    backend.main.Window.Hide();
                    setTimeout( backend.main.Window.Show, 3000 );
                case 'Minimise':
                    backend.main.Window.Minimise();
                    setTimeout( backend.main.Window.Unminimise, 3000 );
                default:
                    backend.main.Window[windowAction]();
                    break;
            }
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
            default:
                params = '';
                break;
                
        }
        
    }
    $: testcodeJs = "import { Window } from '@wails/runtime';\Window." + windowAction + "(" + params + ");";
    $: testcodeGo = '// runtime is given through WailsInit()\nruntime.Window.' + windowAction + '(' + params + ')'; 

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} title="Window" {id} showRun=true>
    <div class="logging-form">
        <form data-wails-no-drag class="mw-full"> 
            <div class="form-group">
                <div>Select Window Method</div>
                {#each windowActions as option, index}
                <div class="custom-radio">
                    <input type="radio" name="window" bind:group="{windowAction}" id="{id}-{option}" value="{option}" disabled={['Show', 'Unminimise'].includes(option)}>
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