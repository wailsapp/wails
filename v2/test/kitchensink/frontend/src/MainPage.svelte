<script>

    import {selectedPage} from './Store';
    import TitlePage from './pages/TitlePage.svelte';
    import Logging from './pages/logging/Logging.svelte';
    import Events from './pages/events/Events.svelte';

</script>

<h4 class="title">{$selectedPage || "" }</h4>
<div class="mainpage">
    {#if $selectedPage == undefined} <TitlePage></TitlePage> {/if}
    {#if $selectedPage == "Logging"} <Logging></Logging> {/if}
    {#if $selectedPage == "Events"} <Events></Events> {/if}
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