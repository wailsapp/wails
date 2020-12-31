<script>
	import {Dialog} from '@wails/runtime';
	import CodeBlock from '../../../components/CodeBlock.svelte';
	import CodeSnippet from '../../../components/CodeSnippet.svelte';
	import jsCode from './code.jsx';
	import goCode from './code.go';

	let isJs = false;
    $: lang = isJs ? 'Javascript' : 'Go';
    let id = "MessageDialog";

    let options = {
        "Type":             "info",
        "Title":            "",
        "Message":          "",
        "Buttons":          [],
		"DefaultButton":    "",
		"CancelButton":     "",
		"Icon":             "",
    }

    function processMessage() {
        if( isJs ) {
            console.log(options);
            Dialog.Message(options);
        } else {
            backend.main.Dialog.Message(options).then( (result) => {
                console.log(result);
            });            
        }
    }

	function prettyPrintArray(json) {
		if (typeof json === 'string') {
			json = JSON.parse(json);
		}
		return JSON.stringify(json, function (k, v) {
			if (v instanceof Array)
				return JSON.stringify(v);
			return v;
		}, 2).replace(/\\/g, '')
			.replace(/"\[/g, '[')
			.replace(/]"/g, ']')
			.replace(/"{/g, '{')
			.replace(/}"/g, '}');
	}

	let dialogTypes = ["Info", "Warning", "Error", "Question"];
    let dialogTypeSelected = dialogTypes[0];
	let buttonInputs = ["","","",""];

	// Keep buttons in sync
	$: {
		options.Buttons = [];
		buttonInputs.forEach( (button) => {
			if ( button.length > 0 ) {
				options.Buttons.push(button);
            }
        })
    }

    // Keep options in sync with dialog type selected
    $: options.Type = dialogTypeSelected.toLowerCase();

	// Inspired by: https://stackoverflow.com/a/54931396
    $: encodedJSOptions = JSON.stringify(options, function (k, v) {
		if (v instanceof Array)
			return JSON.stringify(v);
		return v;
	}, 4)
        .replace(/\\/g, '')
		.replace(/"\[/g, '[')
		.replace(/]"/g, ']')
		.replace(/"{/g, '{')
		.replace(/}"/g, '}');

    $: encodedGoOptions = encodedJSOptions
        .replace(/ {2}"(.*)":/mg, "  $1:")
		.replace(/Type: "(.*)"/mg, "Type: options." + dialogTypeSelected + "Dialog")
        .replace(/Buttons: \[(.*)],/mg, "Buttons: []string{$1},")
        .replace(/\n}/, ",\n}");

    $: testcodeJs = "import { Dialog } from '@wails/runtime';\n\nDialog.Message(" + encodedJSOptions + ");";
    $: testcodeGo = '// runtime is given through WailsInit()\nimport "github.com/wailsapp/wails/v2/pkg/options"\n\nselectedFiles := runtime.Dialog.Message( &options.MessageDialog' + encodedGoOptions + ')';

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} title="Message" {id} showRun=true>
    <div class="browser-form">
        <form data-wails-no-drag class="mw-full"> 
            <div class="form-row row-eq-spacing-sm">
                <div class="form-group">
                    <div>Dialog Type</div>
                    {#each dialogTypes as option}
                        <div class="custom-radio">
                            <input type="radio" name="dialogType" bind:group="{dialogTypeSelected}" id="{id}-{option}" value="{option}">
                            <label for="{id}-{option}">{option}</label>
                        </div>
                    {/each}
                </div>
            </div>
            <div class="form-row row-eq-spacing-sm">
                <div class="col-sm">
                    <label for="{id}-Title">Title</label>
                    <input type="text" class="form-control" id="{id}-Title" bind:value="{options.Title}">
                    <div class="form-text"> The title for the dialog </div>
                </div>
                <div class="col-sm">
                    <label for="{id}-Message">Message</label>
                    <input type="text" class="form-control" id="{id}-Message" bind:value="{options.Message}">
                    <div class="form-text"> The dialog message </div>
                </div>
            </div>
            <div class="form-row row-eq-spacing-sm">
                <div class="col-sm">
                    <label for="{id}-Button1">Button 1</label>
                    <input type="text" class="form-control" id="{id}-Button1" bind:value="{buttonInputs[0]}">
                </div>
                <div class="col-sm">
                    <label for="{id}-Button2">Button 2</label>
                    <input type="text" class="form-control" id="{id}-Button2" bind:value="{buttonInputs[1]}">
                </div>
                <div class="col-sm">
                    <label for="{id}-Button3">Button 3</label>
                    <input type="text" class="form-control" id="{id}-Button3" bind:value="{buttonInputs[2]}">
                </div>
                <div class="col-sm">
                    <label for="{id}-Button4">Button 4</label>
                    <input type="text" class="form-control" id="{id}-Button4" bind:value="{buttonInputs[3]}">
                </div>
            </div>
            <div class="form-row row-eq-spacing-sm">
                <div class="col-sm">
                    <label for="{id}-DefaultButton">Default Button</label>
                    <input type="text" class="form-control" id="{id}-DefaultButton" bind:value="{options.DefaultButton}">
                    <div class="form-text"> The button that is the default option</div>
                </div>
                <div class="col-sm">
                    <label for="{id}-CancelButton">Cancel Button</label>
                    <input type="text" class="form-control" id="{id}-CancelButton" bind:value="{options.CancelButton}">
                    <div class="form-text"> The button that is the cancel option </div>
                </div>
            </div>           
            <div class="form-row row-eq-spacing-sm">
                <div class="col-sm">
                    <label for="{id}-Icon">Icon</label>
                    <input type="text" class="form-control" id="{id}-Icon" bind:value="{options.Icon}">
                    <div class="form-text"> The icon to use in the dialog </div>
                </div>
            </div>
            
            <input class="btn btn-primary" type="button" on:click="{processMessage}" value="Show message dialog using {lang} runtime">

            <CodeSnippet bind:isJs={isJs} jsCode={testcodeJs} goCode={testcodeGo}/>

        </form>
    </div>
</CodeBlock>

