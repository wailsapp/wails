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
    let counts = {};
    let id = UniqueID('events');

    let eventName = "";
    let maxCallbacks = 1;

    function removeSubscriber(eventName, source) {
        listeners.update( (current) => {
            let name = '"' + eventName + '" (' + source + ')';
            const index = current.indexOf(name);
            if (index > -1) { 
                current.splice(index, 1); 
            }
            return current;
        });
    }

    function updateLog(eventName, data, source) {
        let name = '"' + eventName + '" (' + source + ")";
        console.log(counts);
        counts[name].current = counts[name].current + 1;
        console.log(counts);
        let countText = counts[name].current + "/" + counts[name].max;
        if( counts[name].current === counts[name].max) {
            removeSubscriber(name);
        }
        loggingOutput.update( (log) => {
            let datatext = (data ? JSON.stringify(data) : "(No data given)");
            let destroyText = "";
            if( counts[name].current === counts[name].max) {
                destroyText = " (Listener Now Destroyed) ";
            }
            return log + "[" + eventName + " (" + source + ") " + countText + "] data: " + datatext + destroyText + "\n";
        });
    }

    // Subscribe to the Go event calls
    Events.On("onmultiple event fired by go subscriber", (input) => {
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

        counts[name] = { max: maxCallbacks, current: 0 }
        if( isJs ) {
            Events.OnMultiple(eventName, (...data) => {
                updateLog(eventName, data, "JS");
            }, maxCallbacks);
        } else {
            // We call a function in Go to register a subscriber
            // for us
            backend.main.Events.OnMultiple(eventName, maxCallbacks);
        }
    }

    $: testcodeJs = "import { Events } from '@wails/runtime';\nEvents.OnMultiple('" + eventName + "', callback, " + maxCallbacks + ");";
    $: testcodeGo = '// runtime is given through WailsInit()\nruntime.Events.OnMultiple("' + eventName + '", func(optionalData ...interface{} {\n // Process data\n}), " + maxCallbacks + ")'; 

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} {id} title="Events.OnMultiple(eventName, callback)" {description}>
    <div class="logging-form">
        <form data-wails-no-drag class="mw-full"> 
            <div class="form-group">
                <label for="{id}-eventName" class="required">Event Name to subscribe to</label>
                <input type="text" class="form-control" id="{id}-eventName" placeholder="MyEventName" bind:value="{eventName}" required="required">
            </div>
            <div class="form-group">
                <label for="{id}-maxTimes" class="required">Number of times the callback should fire</label>
                <input type="number" class="form-control" id="{id}-maxTimes" placeholder="1" bind:value="{maxCallbacks}">
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
