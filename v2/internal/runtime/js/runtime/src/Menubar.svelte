<script>

    import {menuVisible} from './store'
    import {fade} from 'svelte/transition';

    import {trays} from './store'
    import TrayMenu from "./TrayMenu.svelte";

</script>

{#if $menuVisible }
    <div class="wails-menubar" transition:fade>
    <span class="tray-menus">
    {#each $trays as tray}
        <TrayMenu {tray}></TrayMenu>
    {/each}
    </span>
    </div>
{/if}

<style>

    .tray-menus {
        display: flex;
        flex-direction: row;
        justify-content: flex-end;
    }
    .wails-menubar { position: relative;
        display: block;
        top: 0;
        height: 2rem;
        width: 100%;
        border-bottom: 1px solid #b3b3b3;
        box-shadow: antiquewhite;
        box-shadow: 0px 0px 10px 0px #33333360;
    }
</style>