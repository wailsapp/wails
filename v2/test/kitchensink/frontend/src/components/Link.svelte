<script>
    import runtime from '@wails/runtime';
    import { darkMode } from '../Store';

    export let href;

    $: linkclass = darkMode ? 'dark' : 'light';

    function openLink() {
        runtime.Browser.Open(href);
    }
</script>

<span class="link {linkclass}" on:click="{openLink}"><slot></slot></span>

<style>

    .link {
        text-decoration: underline;
        cursor: pointer;
    }

    .link.dark:hover {
        color: #ff8bfb;
    }

    .link.dark {
        color: #ff8bfb88;
    }

    .link.light {
        color: #5c5cff;
        text-decoration: underline;
    }

    
</style>