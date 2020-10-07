<script>

    import { Log } from '@wailsapp/runtime2';

    var loglevel = 'Debug';
    var message = '';

    var options = ["Debug", "Info", "Warning", "Error", "Fatal"];

    function sendLogMessage() {
        if( message.length > 0 ) {
            Log[loglevel](message);
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

    <h5>Try Me</h5>
    <hr>
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

            <input class="btn btn-primary" type="submit" on:click="{sendLogMessage}" value="Run!">
        </form>
    </div>
</div>

<style>
    .logging-form {
        margin-top: 2rem;
    }
    .list {
        margin-top: 2rem;
        margin-left: 2rem;
    }

    .list li {
        margin-bottom: 0.5rem;
    }

    h5 {
        margin-top: 3rem;
    }
</style>