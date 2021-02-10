<script>

    import { overlayVisible } from './store'
    import { fade,  } from 'svelte/transition';

</script>

{#if $overlayVisible }
<div class="wails-reconnect-overlay" transition:fade="{{ duration: 200 }}">
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
        backdrop-filter: blur(20px) saturate(160%) contrast(45%) brightness(140%);
        z-index: 999999
    }

    .wails-reconnect-overlay-content {
        position: relative;
        top: 50%;
        transform: translateY(-50%);
        margin: 0;
        background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAC8AAAAuCAMAAACPpbA7AAAAflBMVEUAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBAQAAAAAAAAAAAABAQEEBAQAAAAAAAAEBAQAAAADAwMAAAABAQEAAAAAAAAAAAAAAAAAAAACAgICAgIBAQEAAAAAAAAAAAAAAAAAAAACAgIAAAAAAAAAAAAAAAAAAAAAAAAAAAAFBQWCC3waAAAAKXRSTlMALgUMIBk0+xEqJs70Xhb3lu3EjX2EZTlv5eHXvbarQj3cdmpXSqOeUDwaqNAAAAKCSURBVEjHjZTntqsgEIUPVVCwtxg1vfD+L3hHRe8K6snZf+KKn8OewvzsSSeXLruLnz+KHs0gr6DkT3xsRkU6VVn4Ha/UxLe1Z4y64i847sykPBh/AvQ7ry3eFN70oKrfcBJYvm/tQ1qxP4T3emXPeXAkvodPUvtdjbhk+Ft4c0hslTiXVOzxOJ15NWUblQhRsdu3E1AfCjj3Gdm18zSOsiH8Lk4TB480ksy62fiqNo4OpyU8O21l6+hyRtS6z8r1pHlmle5sR1/WXS6Mq2Nl+YeKt3vr+vdH/q4O68tzXuwkiZmngYb4R8Co1jh0+Ww2UTyWxBvtyxLO7QVjO3YOD/lWZpbXDGellFG2Mws58mMnjVZSn7p+XvZ6IF4nn02OJZV0aTO22arp/DgLPtrgpVoi6TPbZm4XQBjY159w02uO0BDdYsfrOEi0M2ulRXlCIPAOuN1NOVhi+riBR3dgwQplYsZRZJLXq23Mlo5njkbY0rZFu3oiNIYG2kqsbVz67OlNuZZIOlfxHDl0UpyRX86z/OYC/3qf1A1xTrMp/PWWM4ePzf8DDp1nesQRpcFk7BlwdzN08ZIALJpCaciQXO0f6k4dnuT/Ewg4l7qSTNzm2SykdHn6GJ12mWc6aCNj/g1cTXpB8YFfr0uVc96aFkkqiIiX4nO+salKwGtIkvfB+Ja8DxMeD3hIXP5mTOYPB4eVT0+32I5ykvPZjesnkGgIREgYnmLrPb0PdV3hoLup2TjcGBPM4mgsfF5BrawZR4/GpzYQzQfrUZCf0TCWYo2DqhdhTJBQ6j4xqmmLN5LjdRIY8LWExiFUsSrza/nmFBqw3I9tEZB9h0lIQSO9if8DkISDAj8CDawAAAAASUVORK5CYII=);
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