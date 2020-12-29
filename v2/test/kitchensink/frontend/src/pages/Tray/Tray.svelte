<script>
    import SetIcon from "./SetIcon/SetIcon.svelte";
</script>

<div>
    Wails includes support for adding icons/menus to the system tray. Tray icons are read from the "trayicons" directory in the project directory and bundled during the compilation phase.
    The filenames follow a convention and only PNG icons are supported. The file `&lt;projectroot&gt;/trayicons/light.png` is given a tray icon ID of 'light'. This is used
    when setting the icon.

    <div style="padding: 15px"></div>

    <SetIcon></SetIcon>
    <br/>     
</div>

<style>
</style>