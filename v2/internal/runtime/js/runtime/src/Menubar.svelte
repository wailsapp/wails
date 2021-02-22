<script>

    import {menuVisible} from './store'
    import {fade} from 'svelte/transition';

    import {trays} from './store'
    import TrayMenu from "./TrayMenu.svelte";
    import {onMount} from "svelte";

    let time = new Date();
    $: day = time.toLocaleString("default", { weekday: "short" })
    $: dom = time.getDate()
    $: mon = time.toLocaleString("default", { month: "short" })
    $: currentTime = time.toLocaleString('en-US', { hour: 'numeric', minute: 'numeric', hour12: true }).toLowerCase()
    $: dateTimeString = `${day} ${dom} ${mon} ${currentTime}`

    onMount(() => {
        const interval = setInterval(() => {
            time = new Date();
        }, 1000);

        return () => {
            clearInterval(interval);
        };
    });

    function handleKeydown(e) {
        // Backtick toggle
        if( e.keyCode == 192 ) {
            menuVisible.update( (current) => {
                return !current;
            });
        }
    }

</script>

{#if $menuVisible }
    <div class="wails-menubar" transition:fade>
    <span class="tray-menus">
    {#each $trays as tray}
        <TrayMenu {tray}/>
    {/each}
    <span class="time">{dateTimeString}</span>
    </span>
    </div>
{/if}

<svelte:window on:keydown={handleKeydown}/>

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
        box-shadow: 0 0 10px 0 #33333360;
    }
    .time {
        padding-left: 0.5rem;
        padding-right: 1.5rem;
        overflow: visible;
        font-size: 14px;
    }
</style>