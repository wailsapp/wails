<script>
    import { Browser } from '@wails/runtime';
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import CodeSnippet from '../../../components/CodeSnippet.svelte';
    import jsCode from './code.jsx';
    import goCode from './code.go';

    var isJs = false;
    $: lang = isJs ? 'Javascript' : 'Go';
    var id = "Browser";

    var userInput = "";

    function processOpen() {
        if( userInput.length > 0 ) {
            if( isJs ) {
                Browser.Open(userInput)
            } else {
                backend.main.Browser.Open(userInput)            
            }
        }
    }

    $: encodedMessage = userInput.replace(`"`, `\"`);
    $: testcodeJs = "import { runtime } from '@wails/runtime';\nruntime.Browser.Open(`" + encodedMessage + "`);";
    $: testcodeGo = '// runtime is given through WailsInit()\nruntime.Browser.Open("' + encodedMessage + '")'; 

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} title="Open" {id} showRun=true>
    <div class="browser-form">
        <form data-wails-no-drag class="mw-full"> 
            <div class="form-group">
                <label for="{id}-userInput" class="required">Enter Filename or URL</label>
                <input type="text" class="form-control" id="{id}-userInput" placeholder="https://www.duckduckgo.com" bind:value="{userInput}" required="required">
            </div>

            <input class="btn btn-primary" type="button" on:click="{processOpen}" value="Open using {lang} runtime">

            <CodeSnippet bind:isJs={isJs} jsCode={testcodeJs} goCode={testcodeGo}></CodeSnippet>

        </form>
    </div>
</CodeBlock>

