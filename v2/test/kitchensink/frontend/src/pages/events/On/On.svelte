<script>
    import { Events } from '@wails/runtime';
    import { writable } from 'svelte/store';
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import description from './description.txt';
    import { UniqueID } from '../../../utils/utils';
    import FakeTerm from '../../../components/FakeTerm.svelte';
    import jsCode from './code.jsx';
    import goCode from './code.go';

    let isJs = false;
    $: lang = isJs ? 'Javascript' : 'Go';

    // Listeners
    let listeners = writable([]);
    let id = UniqueID('events');

    let eventName = "";
    let loggingOutput = writable("");


    function updateLog(eventName, data, source) {
        loggingOutput.update( (log) => {
            let datatext = (data ? JSON.stringify(data) : "(No data given)");
            return log + "[" + eventName + " (" + source + ")] data: " + datatext + "\n";
        });
    }

    // Subscribe to the Go event calls
    Events.On("event fired by go subscriber", (input) => {
        // Format the data for printing
        updateLog(input.Name, input.Data, "Go");
    });

    function subscribe() {
        if (eventName.length == 0) {
            return
        }

        let name = eventName + " (" + (isJs ? 'JS' : 'Go') + ")"
        if( $listeners.includes(name) ) {
            return
        }

        // Add eventName to listeners list
        listeners.update( (current) => {
            return current.concat(name);
        });

        if( isJs ) {
            Events.On(eventName, (...data) => {
                updateLog(eventName, data, "JS");
            })
        } else {
            // We call a function in Go to register a subscriber
            // for us
            backend.main.Events.Subscribe(eventName);
        }
    }

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} {id} title="Events.On(eventName, callback)" {description}>
    <div class="logging-form">
        <form data-wails-no-drag class="mw-full"> 
            {#if $listeners.length > 0 }
            Subscribed to:
            <div class="form-group">
                <ul class="list">
                {#each $listeners as listener}
                    <li>"{listener}"</li>
                {/each}
            </div>
            Now use <code>Events.Emit</code> to trigger the subscribers!<br/>
            {/if}
            <div class="form-group">
                <label for="{id}-eventName" class="required">Event Name to subscribe to</label>
                <input type="text" class="form-control" id="{id}-eventName" placeholder="MyEventName" bind:value="{eventName}" required="required">
            </div>

            <input class="btn btn-primary" type="button" on:click="{subscribe}" value="Subscribe using {lang} runtime">

            <FakeTerm text={$loggingOutput} style="height: 300px; overflow: scroll"></FakeTerm>

        </form>
    </div>
</CodeBlock>
