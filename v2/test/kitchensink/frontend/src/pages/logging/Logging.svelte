<script>

    import { Log } from '@wailsapp/runtime2';
    import CodeBlock from '../../components/CodeBlock.svelte';
    import jsCode from './code.jsx';
    import goCode from './code.go';

    var loglevel = 'Debug';
    var message = '';
    var isJs = false;

    var options = ["Debug", "Info", "Warning", "Error", "Fatal"];

    $: lang = isJs ? 'Javascript' : 'Go';

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
<div>
    <h4>Logging</h4>

    Logging is part of the Wails Runtime and is accessed through the <code>runtime.Log</code> object. There are 5 methods available:
    
    <ul class="list">
        <li>Debug</li>
        <li>Info</li>
        <li>Warning</li>
        <li>Error</li>
        <li>Fatal</li>
    </ul>
    All methods will log to the console and <code>Fatal</code> will also exit the program.
    
    <div style="padding: 15px"></div>
    
    <CodeBlock bind:isJs={isJs} {jsCode} {goCode} title="Logging" >
        <div class="logging-form">
            <form data-wails-no-drag class="w-400 mw-full"> <!-- w-400 = width: 40rem (400px), mw-full = max-width: 100% -->
                <!-- Radio -->
                <div class="form-group">
                    <label for="Debug">Log Level</label>
                    {#each options as option}
                    <div class="custom-radio">
                        <input type="radio" name="logging" bind:group="{loglevel}" id="{option}" value="{option}">
                        <label for="{option}">{option}</label>
                    </div>   
                    {/each}
                </div>

                <!-- Input -->
                <div class="form-group">
                    <label for="message" class="required">Message</label>
                    <input type="text" class="form-control" id="message" placeholder="Hello World!" bind:value="{message}" required="required">
                </div>

                <input class="btn btn-primary" type="button" on:click="{sendLogMessage}" value="Call {lang} method">
            </form>
        </div>
    </CodeBlock>
</div>
