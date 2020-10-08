<script>

    import CodeBlock from '../components/CodeBlock.svelte';

</script>

<div>
    <h4>Events</h4>

    Wails includes a unified events system which allows you to do the following in either Go or Javascript:

    <ul>
        <li>On</li>
        <li>Once</li>
        <li>OnMultiple</li>
        <li>Emit</li>
    </ul>

    <CodeBlock title="Events.On(eventName, callback)" description="etc etc etc "></CodeBlock>
</div>

<style>
    h5 {
        margin-top: 3rem;
    }
</style>