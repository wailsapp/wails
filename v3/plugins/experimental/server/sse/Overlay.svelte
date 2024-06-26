<script>

    import {overlayVisible} from './store'
    import {fade,} from 'svelte/transition';
</script>

{#if $overlayVisible }
    <div class="wails-reconnect-overlay" transition:fade="{{ duration: 300 }}">
        <div class="wails-reconnect-overlay-content">
            <div class="wails-reconnect-overlay-loadingspinner"></div>
        </div>
    </div>
{/if}

<style>
    .wails-reconnect-overlay {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        backdrop-filter: blur(2px) saturate(0%) contrast(50%) brightness(25%);
        z-index: 999999
    }

    .wails-reconnect-overlay-content {
        position: relative;
        top: 50%;
        transform: translateY(-50%);
        margin: 0;
        background-image: url(data:image/png;base64,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);
        background-repeat: no-repeat;
        background-position: center
    }

    .wails-reconnect-overlay-loadingspinner {
        pointer-events: none;
        width: 2.5em;
        height: 2.5em;
        border: .4em solid transparent;
        border-color: #f00 #eee0 #f00 #eee0;
        border-radius: 50%;
        animation: loadingspin 1s linear infinite;
        margin: auto;
        padding: 2.5em
    }

    @keyframes loadingspin {
        100% {
            transform: rotate(360deg);
        }
    }

</style>