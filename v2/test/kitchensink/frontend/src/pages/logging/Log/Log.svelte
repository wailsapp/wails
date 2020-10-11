<script>
    import { Log } from '@wails/runtime';
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import jsCode from './code.jsx';
    import goCode from './code.go';

    import { logLevel } from '../../../Store';

    var message = '';
    var isJs = false;

    const loglevels = ["Trace", "Debug", "Info", "Warning", "Error", "Fatal", "Print"];
    var loglevel = loglevels[0];

    $: lang = isJs ? 'Javascript' : 'Go';

    var id = "Logging";

    function sendLogMessage() {
        if( message.length > 0 ) {
            if( isJs ) {
                // Call JS runtime
                Log[loglevel](message);
            } else {
                // Call Go method which calls Go Runtime
                backend.main.Logger[loglevel](message);              
            }
        }
    }

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} title="Logging" {id}>
    <div class="logging-form">
        <form data-wails-no-drag class="w-500 mw-full"> 
            <div class="form-group">
                <label for="Debug">Select Logging Level</label>
                {#each loglevels as option, index}
                {#if index === $logLevel}
                <span style="margin-top: 5px; height: 20px; display: inline-block;"><hr style="width: 270px;display: inline-block; vertical-align: middle; margin-right: 10px"/> Current Log Level </span>
                {/if}
                <div class="custom-radio">
                    <input type="radio" name="logging" bind:group="{loglevel}" id="{id}-{option}" value="{option}">
                    <label for="{id}-{option}">{option}</label>
                </div>   
                {/each}
            </div>

            <div class="form-group">
                <label for="{id}-message" class="required">Message</label>
                <input type="text" class="form-control" id="{id}-message" placeholder="Hello World!" bind:value="{message}" required="required">
            </div>

            <input class="btn btn-primary" type="button" on:click="{sendLogMessage}" value="Log using {lang} runtime">
        </form>
    </div>
</CodeBlock>

