<script>

    export let menu;

</script>

<div class="menu">
    {#each menu.Menu.Items as menuItem}
        {#if menuItem.Image.length > 0}
            <img alt="" src="data:image/png;base64,{menuItem.Image}"/>
        {/if}
    <span class="menuitem">{menuItem.Label}</span>
    {/each}
</div>

<style>

</style>