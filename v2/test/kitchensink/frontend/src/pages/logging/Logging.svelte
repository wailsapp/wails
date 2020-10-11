<script>
    import Log from './Log/Log.svelte';
    import SetLogLevel from './SetLogLevel/SetLogLevel.svelte';
    
    const loglevels = ["Trace", "Debug", "Info", "Warning", "Error", "Fatal", "Print"];
    
</script>
<div>
    <h4>Logging</h4>

    Logging is part of the Wails Runtime and is accessed through the <code>runtime.Log</code> object. There are {loglevels.length} methods available:
    
    <ul class="list">
        {#each loglevels as option}
        <li>{option}</li>
        {/each}
    </ul>
    <br/>

    The default logger will log all messages to the console EG:<br/>
    <pre>
INFO  | I am an Info message
ERROR | I am an Error message
WARN  | I am a Warning message
    </pre>

    <code>Fatal</code> will print the message and then immediately exit the program.
    
    <div style="padding: 15px"></div>
    
    <Log></Log>
    <br/>
    <SetLogLevel></SetLogLevel>


</div>
