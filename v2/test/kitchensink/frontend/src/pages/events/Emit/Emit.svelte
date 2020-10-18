<script>
    import { Events } from '@wails/runtime';
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import CodeSnippet from '../../../components/CodeSnippet.svelte';
    import FakeTerm from '../../../components/FakeTerm.svelte';
    import description from './description.txt';
    import { UniqueID } from '../../../utils/utils';
    import jsCode from './code.jsx';
    import goCode from './code.go';
    import { loggingOutput } from '../On/store';


    let isJs = false;
    $: lang = isJs ? 'Javascript' : 'Go';

    let id = UniqueID('events');

    let eventName = "";

    let dataText = "";
    let placeholder = `123, "hello", true`;
    let formattedData = "";

    function emitEvent() {
        let data = JSON.parse("["+formattedData+"]");
        if( isJs ) {
            Events.Emit(eventName.trim(), ...data);
        } else {
            backend.main.Events.Emit(eventName, data);
        }
    }

    let dataValid = true;

    $: {
        console.log("REactive function triggered!");

        if ( dataText.length === 0 ) {
            dataValid = true;
            formattedData = "";
        } else {
            try {
                formattedData = JSON.stringify(JSON.parse("["+dataText+"]")).slice(1,-1);
                dataValid = true;
            } catch(e) {
                formattedData = "";
                dataValid = false;
            }
        }
    }

    $: emitCodeJs = `import { Events } from '@wails/runtime';\n\nEvents.Emit("` + eventName + `"` + (formattedData.length > 0 ? ',' + formattedData : '') + `);`;
    $: emitCodeGo = `Events.Emit("` + eventName + `"` + (formattedData.length > 0 ? ',' + formattedData : '') + `);`;
    

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} {id} title="Events.Emit(eventName [, data])" {description}>
    <div class="logging-form">
        <form data-wails-no-drag class="mw-full"> 

            <div class="form-group">
                <label for="{id}-eventName" class="required">Event Name to Emit</label>
                <input type="text" class="form-control" id="{id}-eventName" placeholder="MyEventName" bind:value="{eventName}" required="required">
            </div>

            <div class="form-group">
                <label for="{id}-data">Optional data:</label>
                <input type="text" class="form-control" style="{ dataValid ? '' : 'border: 1px solid red;' }" id="{id}-data" {placeholder} bind:value="{dataText}">
            </div>

            <input class="btn btn-primary" type="button" on:click="{emitEvent}" disabled="{!dataValid || eventName.trim().length == 0}" value="Emit using {lang} runtime">

            <CodeSnippet bind:isJs={isJs} jsCode={emitCodeJs} goCode={emitCodeGo}></CodeSnippet>

            Listener output:
            <FakeTerm text={$loggingOutput} style="height: 300px; overflow: scroll"></FakeTerm>

        </form>
    </div>
</CodeBlock>
