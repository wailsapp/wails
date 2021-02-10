<script>
    import Menu from "./Menu.svelte";

    export let tray;

</script>

<span class="tray-menu">
    <span class="label">{tray.Label}</span>
<!--    <Menu menu="{tray.ProcessedMenu}"/>-->
</span>

<style>

    .tray-menu {
        padding-left: 0.5rem;
        padding-right: 0.5rem;
        overflow: visible;
        font-size: 14px;
    }
</style>