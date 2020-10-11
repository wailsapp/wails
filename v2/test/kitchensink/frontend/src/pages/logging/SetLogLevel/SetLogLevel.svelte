<script>
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import { logLevel } from '../../../Store';

    import jsCode from './code.jsx';
    import goCode from './code.go';

    var options = ["Trace", "Debug", "Info", "Warning", "Error"];

    var id = "SetLogLevel";
    let loglevelText = options[$logLevel];

    $: logLevel.set(options.indexOf(loglevelText));
    
</script>

<CodeBlock {jsCode} {goCode} title="SetLogLevel" {id}>
    <div class="logging-form">
        <form data-wails-no-drag class="w-500 mw-full">
            <!-- Radio -->
            <div class="form-group">
                <label for="Debug">Select Logging Level</label>
                {#each options as option}
                <div class="custom-radio">
                    <input type="radio" name="logging" bind:group="{loglevelText}" id="{id}-{option}" value="{option}">
                    <label for="{id}-{option}">{option}</label>
                </div>   
                {/each}
            </div>
        </form>
    </div>
</CodeBlock>