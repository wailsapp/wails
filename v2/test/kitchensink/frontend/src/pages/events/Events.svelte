<script>
    import On from './On/On.svelte';
</script>

<div>
    <h4>Events</h4>

    Wails includes a unified events system which allows you to send events in either Go or Javascript and receive them in either Go or Javascript. The events methods may be accessed through the <code>runtime.Events</code> object. Available methods are:

    <ul class="list">
        <li>On</li>
        <li>Once</li>
        <li>OnMultiple</li>
        <li>Emit</li>
    </ul>

    Events are intraoperable and data sent in events is translated between Go and Javascript automatically.

    <div style="padding: 15px"></div>
     
    <On></On>

</div>

<style>
</style>