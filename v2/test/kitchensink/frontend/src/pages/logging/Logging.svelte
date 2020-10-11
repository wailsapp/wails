<script>
    import FakeTerm from '../../components/FakeTerm.svelte';
    import Link from '../../components/Link.svelte';
    import Log from './Log/Log.svelte';
    import SetLogLevel from './SetLogLevel/SetLogLevel.svelte';
    
    const loglevels = ["Trace", "Debug", "Info", "Warning", "Error", "Fatal", "Print"];
    
</script>
<div>
    <h4>Logging</h4>

    Logging is part of the Wails Runtime and is accessed through the <code>runtime.Log</code> object. 
    
    There are {loglevels.length} methods available:
    
    <ol class="list">
        {#each loglevels as option}
        <li>{option}</li>
        {/each}
    </ol>
    <br/>

    Logs are only output if they are above or equal to the current log level. Example: If the current log level is 
    <code>Info</code>, then <code>Trace</code> and <code>Debug</code> messages will be surpressed. <br/>
    <code>Fatal</code> will log the message and then immediately exit the program.<br/>
    <code>Print</code> will send a raw message to the log. You can use <code>Print</code> at any log level.<br/>
    
    The default logger will log messages to the console in the following format:<br/>
    <FakeTerm>
INFO  | I am an Info message
ERROR | I am an Error message
WARN  | I am a Warning message
    </FakeTerm>
    Custom loggers may be given to your Wails application. More details <Link href="https://www.google.com">here</Link>.
    
    <div style="padding: 15px"></div>
    
    <Log></Log>
    <br/>
    <SetLogLevel></SetLogLevel>


</div>
