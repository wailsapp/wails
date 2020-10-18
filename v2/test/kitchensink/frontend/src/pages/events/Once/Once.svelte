<script>
    import { Events } from '@wails/runtime';
    import { writable } from 'svelte/store';
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import CodeSnippet from '../../../components/CodeSnippet.svelte';
    import description from './description.txt';
    import { UniqueID } from '../../../utils/utils';
    import FakeTerm from '../../../components/FakeTerm.svelte';
    import jsCode from './code.jsx';
    import goCode from './code.go';
    import { loggingOutput } from '../On/store';

    let isJs = false;
    $: lang = isJs ? 'Javascript' : 'Go';

    // Listeners
    let listeners = writable([]);
    let id = UniqueID('events');

    let eventName = "";

    function removeSubscriber(eventName, source) {
        listeners.update( (current) => {
            let name = '"' + eventName + '" (' + source + ')';
            console.log(name);
            const index = current.indexOf(name);
            console.log("index = ", index);
            if (index > -1) { 
                current.splice(index, 1); 
            }
            console.log(current);
            return current;
        });
    }

    function updateLog(eventName, data, source) {
        removeSubscriber(eventName, source);
        loggingOutput.update( (log) => {
            let datatext = (data ? JSON.stringify(data) : "(No data given)");
            return log + "[" + eventName + " (" + source + ")] data: " + datatext + " (Listener now destroyed)\n";
        });
    }

    // Subscribe to the Go event calls
    Events.On("once event fired by go subscriber", (input) => {
        // Format the data for printing
        updateLog(input.Name, input.Data, "Go");
    });

    function subscribe() {
        if (eventName.length == 0) {
            return
        }

        let name = '"' + eventName + '" (' + (isJs ? 'JS' : 'Go') + ")"
        if( $listeners.includes(name) ) {
            return
        }

        // Add eventName to listeners list
        listeners.update( (current) => {
            return current.concat(name);
        });

        if( isJs ) {
            Events.Once(eventName, (...data) => {
                updateLog(eventName, data, "JS");
            })
        } else {
            // We call a function in Go to register a subscriber
            // for us
            backend.main.Events.Once(eventName);
        }
    }

    $: testcodeJs = "import { Events } from '@wails/runtime';\nEvents.Once('" + eventName + "', callback);";
    $: testcodeGo = '// runtime is given through WailsInit()\nruntime.Events.Once("' + eventName + '", func(optionalData ...interface{} {\n // Process data\n}))'; 

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} {id} title="Events.Once(eventName, callback)" {description}>
    <div class="logging-form">
        <form data-wails-no-drag class="mw-full"> 
            <div class="form-group">
                <label for="{id}-eventName" class="required">Event Name to subscribe to</label>
                <input type="text" class="form-control" id="{id}-eventName" placeholder="MyEventName" bind:value="{eventName}" required="required">
            </div>

            <input class="btn btn-primary" type="button" on:click="{subscribe}" value="Subscribe using {lang} runtime">
            <CodeSnippet bind:isJs={isJs} jsCode={testcodeJs} goCode={testcodeGo}></CodeSnippet>
            {#if $listeners.length > 0 }
            <div class="form-group" style="margin-top:10px">
                Subscribed to:
                <ul class="list">
                {#each $listeners as listener}
                    <li>{listener}</li>
                {/each}
            </div>
            Now use <code>Events.Emit</code> to trigger the subscribers!<br/>
            <div style="margin-top: 10px">Subscriber output will be printed below:</div>
            <FakeTerm text={$loggingOutput} style="height: 300px; overflow: scroll"></FakeTerm>
            {/if}

        </form>
    </div>
</CodeBlock>
