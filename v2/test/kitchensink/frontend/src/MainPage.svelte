<script>

    import {selectedPage} from './Store';
    import TitlePage from './TitlePage.svelte';
    import Logging from './Logging.svelte';
    import Events from './Events.svelte';

</script>

<div class="mainpage">
    {#if $selectedPage == undefined} <TitlePage ></TitlePage> {/if}
    {#if $selectedPage == "Logging"} <Logging></Logging> {/if}
    {#if $selectedPage == "Events"} <Events></Events> {/if}
</div>

<style>
    .mainpage {
        margin-top: 55px;
        margin-left: 25px;
    }
    h3 {
        margin-left: 20px;
    }
</style>