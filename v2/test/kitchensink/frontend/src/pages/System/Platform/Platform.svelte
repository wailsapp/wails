<script>
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import CodeSnippet from '../../../components/CodeSnippet.svelte';
    import FakeTerm from '../../../components/FakeTerm.svelte';
    import jsCode from './code.jsx';
    import goCode from './code.go';
    import { writable } from 'svelte/store';
    import { System } from '@wails/runtime';
    import description from './description.txt';

    var isJs = false;
    $: lang = isJs ? 'Javascript' : 'Go';
    var id = "Platform";

    let output = writable("");

    function log(message) {
        output.update( (current) => {
            current += message;
            return current;
        });
    }

    function getPlatform() {
        if( isJs ) {
            log("Platform from JS runtime: " + System.Platform() + "\n");
        } else {
            backend.main.System.Platform().then( (platformFromGo) => {
                log("Platform from Go runtime: " + platformFromGo + "\n");
            })           
        }
    }

    $: testcodeJs = "import { runtime } from '@wails/runtime';\nruntime.Log.Info('Platform from JS runtime: ' + runtime.System.Platform);";
    $: testcodeGo = '// runtime is given through WailsInit()\nruntime.Log.Info("Platform from Go runtime: " + runtime.System.Platform)'; 

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} {description} title="Platform()" {id} showRun=true>
    <div class="browser-form">
        <form data-wails-no-drag class="mw-full"> 
            <input class="btn btn-primary" type="button" on:click="{getPlatform}" value="Fetch platform using {lang} runtime">

            <CodeSnippet bind:isJs={isJs} jsCode={testcodeJs} goCode={testcodeGo}></CodeSnippet>
            <FakeTerm text={$output} style="height: 100px; overflow: scroll"></FakeTerm>

        </form>
    </div>
</CodeBlock>

