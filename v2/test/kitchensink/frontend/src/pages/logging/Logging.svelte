<script>
    import FakeTerm from '../../components/FakeTerm.svelte';
    import Link from '../../components/Link.svelte';
    import Log from './Log/Log.svelte';
    import SetLogLevel from './SetLogLevel/SetLogLevel.svelte';
    
    const loglevels = ["Trace", "Debug", "Info", "Warning", "Error", "Fatal", "Print"];
    
</script>
<div>
    Logging is part of the Wails Runtime and is accessed through the <code>runtime.Log</code> object. 
    
    There are {loglevels.length} methods available:
    
    <ol class="list">
        {#each loglevels as option}
        <li>{option}</li>
        {/each}
    </ol>
    <br/>

    Logs are only output if they are above or equal to the current log level. Example: If the current log level is 
    <code>Info</code>, then <code>Trace</code> and <code>Debug</code> messages will be surpressed. <br/><br/>
    <code>Fatal</code> will log the message and then immediately exit the program.<br/><br/>
    <code>Print</code> will send a raw message to the log regardless of log level.<br/><br/>

    The default logger will log messages to the console in the following format:<br/>
    <FakeTerm>
TRACE | I am a Trace message
DEBUG | I am a Debug message
INFO  | I am an Info message
WARN  | I am a Warning message
ERROR | I am an Error message
FATAL | I am a Fatal message
I am a Print message
    </FakeTerm>
    <br/>
    Custom loggers may be given to your Wails application. More details <Link href="https://wails.app">here</Link>.
    
    <br/><br/>
    
    <Log></Log>
    <br/>
    <SetLogLevel></SetLogLevel>
</div>

