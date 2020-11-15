<script>

    import {selectedPage} from './Store';
    import TitlePage from './pages/TitlePage.svelte';
    import Logging from './pages/Logging/Logging.svelte';
    import Events from './pages/Events/Events.svelte';
    import Browser from './pages/Browser/Browser.svelte';
    import Dialog from './pages/Dialog/Dialog.svelte';
    import System from './pages/System/System.svelte';
    import Window from './pages/Window/Window.svelte';
    
</script>

<h4 class="title">{$selectedPage || "" }</h4>
<div class="mainpage">
    {#if $selectedPage == undefined} <TitlePage></TitlePage> {/if}
    {#if $selectedPage == "Logging"} <Logging></Logging> {/if}
    {#if $selectedPage == "Events"} <Events></Events> {/if}
    {#if $selectedPage == "Browser"} <Browser></Browser> {/if}
    {#if $selectedPage == "Dialog"} <Dialog></Dialog> {/if}
    {#if $selectedPage == "System"} <System></System> {/if}
    {#if $selectedPage == "Window"} <Window></Window> {/if}
</div>

<style>
    .title {
        position: fixed;
        top: 50px;
        margin-left: 20px;
    }
    .mainpage {
        position: absolute;
        top: 70px;
        margin-left: 20px;
        padding-right: 10px;
        margin-top: 20px;
        overflow-y: auto;
        height: calc(100% - 100px);
        width: calc(100% - 20px);
    }
</style>