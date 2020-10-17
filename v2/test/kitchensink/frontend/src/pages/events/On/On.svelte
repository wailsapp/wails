<script>
    import { Events } from '@wails/runtime';
    import { writable } from 'svelte/store';
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import description from './description.txt';
    import { UniqueID } from '../../../utils/utils';
    import FakeTerm from '../../../components/FakeTerm.svelte';

    let isJs = false;
    let jsCode = "js";
    let goCode = "go";
    $: lang = isJs ? 'Javascript' : 'Go';

    // Listeners
    let listeners = writable([]);
    let id = UniqueID('events');

    let eventName = "";
    let loggingOutput = writable("");

    function subscribe() {
        if (eventName.length == 0) {
            return
        }

        // Add eventName to listeners list
        listeners.update( (current) => {
            // Don't add twice
            if( current.includes(eventName) ) {
                return current;
            }
            return current.concat(eventName);
        });

        if( isJs ) {
            console.log("Adding listener for " + eventName);
            Events.On(eventName, (data) => {
                console.log("CALLED! " + eventName);
                loggingOutput.update( (log) => {
                    let datatext = (data ? JSON.stringify(data) : "(No data given)");
                    return log + "[" + eventName + "] " + datatext + "\n";
                });
            })
        } else {
            console.log("go!");
        }
    }

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} {id} title="Events.On(eventName, callback)" {description}>
    <div class="logging-form">
        <form data-wails-no-drag class="mw-full"> 
            Subscribed to:
            <div class="form-group">
                <ul class="list">
                {#each $listeners as listener}
                    <li>"{listener}"</li>
                {/each}
            </div>

            <div class="form-group">
                <label for="{id}-eventName" class="required">Event Name to subscribe to</label>
                <input type="text" class="form-control" id="{id}-eventName" placeholder="MyEventName" bind:value="{eventName}" required="required">
            </div>

            <input class="btn btn-primary" type="button" on:click="{subscribe}" value="Subscribe using {lang} runtime">

            <FakeTerm text={$loggingOutput} style="height: 300px; overflow: scroll"></FakeTerm>

        </form>
    </div>
</CodeBlock>
