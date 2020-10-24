<script>
    import { Dialog } from '@wails/runtime';
    import CodeBlock from '../../../components/CodeBlock.svelte';
    import CodeSnippet from '../../../components/CodeSnippet.svelte';
    import jsCode from './code.jsx';
    import goCode from './code.go';

    var isJs = false;
    $: lang = isJs ? 'Javascript' : 'Go';
    var id = "Dialog";

    let options = {
        "DefaultDirectory":           "",
        "DefaultFilename":            "",
        "Title":                      "",
        "Filters":                    "",
        "AllowFiles":                 false,
        "AllowDirectories":           false,
        "AllowMultiple":              false,
        "ShowHiddenFiles":            false,
        "CanCreateDirectories":       false,
        "ResolveAliases":             false,
        "TreatPackagesAsDirectories": false
    }

    function processOpen() {
        if( isJs ) {
            console.log(options);
            Dialog.Open(options);
        } else {
            backend.main.Dialog.Open(options).then( (result) => {
                console.log(result);
            });            
        }
    }

    $: encodedJSOptions = JSON.stringify(options, null, "  ");
    $: encodedGoOptions = encodedJSOptions
        .replace(/\ {2}"(.*)":/mg, "  $1:")
        .replace(/\n}/, ",\n}");
    
    $: testcodeJs = "import { runtime } from '@wails/runtime';\nruntime.Dialog.Open(" + encodedJSOptions + ");";
    $: testcodeGo = '// runtime is given through WailsInit()\nimport "github.com/wailsapp/wails/v2/pkg/options"\n\nselectedFiles := runtime.Dialog.Open( &options.OpenDialog' + encodedGoOptions + ')'; 

</script>

<CodeBlock bind:isJs={isJs} {jsCode} {goCode} title="Open" {id} showRun=true>
    <div class="browser-form">
        <form data-wails-no-drag class="mw-full"> 
            <div class="form-row row-eq-spacing-sm">
                <div class="col-sm">
                    <label for="{id}-Title">Title</label>
                    <input type="text" class="form-control" id="{id}-Title" bind:value="{options.Title}">
                    <div class="form-text"> The title for the dialog </div>
                </div>        
                <div class="col-sm">
                    <label for="{id}-defaultDirectory">Default Directory</label>
                    <input type="text" class="form-control" id="{id}-defaultDirectory" bind:value="{options.DefaultDirectory}">
                    <div class="form-text"> The directory the dialog will default to </div>
                </div>
            </div>
            <div class="form-row row-eq-spacing-sm">
                <div class="col-sm">
                    <label for="{id}-defaultFilename">Default Filename</label>
                    <input type="text" class="form-control" id="{id}-defaultFilename" bind:value="{options.DefaultFilename}">
                    <div class="form-text"> The filename the dialog will suggest to use </div>
                </div>                
                <div class="col-sm">
                    <label for="{id}-Filters">Filters</label>
                    <input type="text" class="form-control" id="{id}-Filters" bind:value="{options.Filters}">
                    <div class="form-text"> A list of extensions eg <code>*.jpg,*.jpeg</code> </div>
                </div>       
            </div>     
            <div class="form-row row-eq-spacing-sm">
                <div class="col-sm">
                    <input type="checkbox" id="{id}-AllowFiles" bind:checked="{options.AllowFiles}">
                    <label for="{id}-AllowFiles">Allow files to be selected</label>
                </div>            
                <div class="col-sm">
                    <input type="checkbox" id="{id}-AllowDirectories" bind:checked="{options.AllowDirectories}">
                    <label for="{id}-AllowDirectories">Allow directories to be selected</label>
                </div>    
            </div>        
            <div class="form-row row-eq-spacing-sm">
                <div class="col-sm">
                    <input type="checkbox" id="{id}-AllowMultiple" bind:checked="{options.AllowMultiple}">
                    <label for="{id}-AllowMultiple">Allow multiple selection</label>
                </div>            
                <div class="col-sm">
                    <input type="checkbox" id="{id}-CanCreateDirectories" bind:checked="{options.CanCreateDirectories}">
                    <label for="{id}-CanCreateDirectories">Can create directories</label>
                </div>
            </div>
            
            <input class="btn btn-primary" type="button" on:click="{processOpen}" value="Open using {lang} runtime">

            <CodeSnippet bind:isJs={isJs} jsCode={testcodeJs} goCode={testcodeGo}></CodeSnippet>

        </form>
    </div>
</CodeBlock>

