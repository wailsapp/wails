<script>

    import {selectedPage} from './Store';
    import TitlePage from './pages/TitlePage.svelte';
    import Logging from './pages/logging/Logging.svelte';
    import Events from './pages/events/Events.svelte';

</script>

<div class="mainpage">
    {#if $selectedPage == undefined} <TitlePage></TitlePage> {/if}
    {#if $selectedPage == "Logging"} <Logging></Logging> {/if}
    {#if $selectedPage == "Events"} <Events></Events> {/if}
</div>

<style>
    .mainpage {
        margin-left: 45px;
        margin-right: 45px;
        margin-top: 20px;
        margin-bottom: 20px;
    }
</style>