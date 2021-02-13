<script>
    import Menu from "./Menu.svelte";
    import { selectedMenu } from "./store";

    export let tray = null;

    $: hidden = $selectedMenu !== tray;

    function closeMenu() {
        selectedMenu.set(null);
    }

    function trayClicked() {
        if ( $selectedMenu !== tray ) {
            selectedMenu.set(tray);
        } else {
            selectedMenu.set(null);
        }
    }
    // Source: https://svelte.dev/repl/0ace7a508bd843b798ae599940a91783?version=3.16.7
    /** Dispatch event on click outside of node */
    function clickOutside(node) {

        const handleClick = event => {
            if (node && !node.contains(event.target) && !event.defaultPrevented) {
                node.dispatchEvent(
                    new CustomEvent('click_outside', node)
                )
            }
        }

        document.addEventListener('click', handleClick, true);

        return {
            destroy() {
                document.removeEventListener('click', handleClick, true);
            }
        }
    }
</script>

<span class="tray-menu" use:clickOutside on:click_outside={closeMenu}>
    <!--{#if tray.Image && tray.Image.length > 0}-->
    <!--    <img alt="" src="data:image/png;base64,{tray.Image}"/>-->
    <!--{/if}-->
    <span class="label" on:click={trayClicked}>{tray.Label}</span>
    {#if tray.ProcessedMenu }
        <Menu menu="{tray.ProcessedMenu}" {hidden}/>
    {/if}
</span>

<style>

    .tray-menu {
        padding-left: 0.5rem;
        padding-right: 0.5rem;
        overflow: visible;
        font-size: 14px;
    }

    .label {
        text-align: right;
        padding-right: 10px;
    }
</style>