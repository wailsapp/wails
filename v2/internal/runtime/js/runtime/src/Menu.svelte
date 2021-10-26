s
<script>

    export let menu;

    export let hidden = true;
</script>

{#if !hidden}
    <div class="menu">
        <span class="menu-title">{menu.Label}</span>
        <!--{#each menu.Menu.Items as menuItem}-->
        <!--    <div class="menuitem">-->
        <!--    {#if menuItem.Image }-->
        <!--        <div><img alt="" src="data:image/png;base64,{menuItem.Image}"/></div>-->
        <!--    {/if}-->
        <!--    <div class="menulabel">{menuItem.Label}</div>-->
        <!--    </div>-->
        <!--{/each}-->
    </div>
{/if}

<style>

    .menu {
        padding: 3px;
        background-color: #0008;
        color: #EEF;
        border-radius: 5px;
        margin-top: 5px;
        position: absolute;
        backdrop-filter: blur(3px) saturate(160%) contrast(45%) brightness(140%);
        border: 1px solid rgb(88,88,88);
        box-shadow: 0 0 1px rgb(146,146,148) inset;

    }

    .menuitem {
        display: flex;
        align-items: center;
        padding: 1px 5px;
    }

    .menuitem:hover {
        display: flex;
        align-items: center;
        background-color: rgb(57,131,223);
        padding: 1px 5px;
        border-radius: 5px;
    }

    .menuitem img {
        padding-right: 5px;
    }

</style>