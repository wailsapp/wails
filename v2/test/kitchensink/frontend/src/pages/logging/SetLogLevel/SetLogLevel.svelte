<script>
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import CodeSnippet from '../../../components/CodeSnippet.svelte';
    import { logLevel } from '../../../Store';

    import { Log } from '@wails/runtime';

    import jsCode from './code.jsx';
    import goCode from './code.go';

    var options = ["Trace", "Debug", "Info", "Warning", "Error"];
    let isJs = false;
    var id = "SetLogLevel";
    let loglevelText = options[$logLevel];

    $: setLogLevelMethod = isJs ? Log.SetLogLevel : backend.main.Logger.SetLogLevel;

    function setLogLevel() {
        let logLevelUpper = loglevelText.toUpperCase();
        let logLevelNumber = Log.Level[logLevelUpper];
        setLogLevelMethod(logLevelNumber);
    };
    $: lang = isJs ? 'Javascript' : 'Go';

    let description = `You can set the log level using Log.SetLogLevel(). It accepts a log level (number) but the log levels supported have been added to Log: Log.TRACE
`;    
    $: testcodeJs = "import { Log } from '@wails/runtime';\nLog.SetLogLevel(Log.Level." + loglevelText.toUpperCase() + ");";
    $: testcodeGo = 'import "github.com/wailsapp/wails/v2/pkg/logger"\n\n// runtime is given through WailsInit()\nruntime.Log.SetLogLevel(logger.' + loglevelText.toUpperCase() + ')'; 
</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} title="SetLogLevel" {id} {description}>
    <div class="logging-form">
        <form data-wails-no-drag class="mw-full">
            <!-- Radio -->
            <div class="form-group">
                <label for="Debug">Select Logging Level</label>
                {#each options as option}
                <div class="custom-radio">
                    <input type="radio" name="logging" bind:group="{loglevelText}" id="{id}-{option}" value="{option}">
                    <label for="{id}-{option}">{option}</label>
                </div>   
                {/each}
            </div>
            <input class="btn btn-primary" type="button" on:click="{setLogLevel}" value="SetLogLevel using {lang} runtime">
            <CodeSnippet bind:isJs={isJs} jsCode={testcodeJs} goCode={testcodeGo}></CodeSnippet>
        </form>
    </div>
</CodeBlock>