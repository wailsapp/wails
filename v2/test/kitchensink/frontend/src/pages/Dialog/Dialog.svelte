<script>
	import Open from './Open/Open.svelte';
	import Save from './Save/Save.svelte';
	import Message from "./Message/Message.svelte";
</script>
<div>
    Sometimes there is a need to prompt the user for input. Dialogs are a good approach to this and Wails provides access to native dialogs. The dialog types supported are:

    <ul class="list">
        <li>Open - Prompts the user to select one or more files or folders</li>
        <li>Save - Prompts the user to select a file or input a filename</li>
        <li>Message - Prompts the user with information and optionally provides action buttons</li>
    </ul>

    <br/><br/>

    <Open/>
    <br/>
    <Save/>
    <br/>
    <Message/>
</div>

