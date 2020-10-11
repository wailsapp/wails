<script>
    import { darkMode } from '../Store';

    $: termClass = $darkMode ? 'faketerm-dark' : 'faketerm-light';
</script>

<div class="common {termClass}">
<pre>
<slot></slot>
</pre>
</div>

<style>
    pre {
        margin: 0;
        padding: 5px;
    }

    .common {
        font-family: 'Courier New', Courier, monospace;
        padding: 5px;
        white-space: pre-line;
        margin-top: 10px;
        margin-bottom: 10px;
        border: 1px solid #5555;
    }

    .faketerm-dark {
        background-color: black;
        color: white;
    }

    .faketerm-light {
        background-color: #ddd;
        color: black;
    }
</style>