
<script>

  import { darkMode, selectedPage } from './Store';
  import MainPage from './MainPage.svelte';
 
  // Hightlight CSS
  import { atomOneDark, atomOneLight } from "svelte-highlight/styles";
  $: css = $darkMode ? atomOneDark : atomOneLight;

  function linkClicked(event) {
    let linkText = event.target.innerText;
    selectedPage.set(linkText); 
    console.log(event.target.innerText);
  }

  function homepageClicked() {
    selectedPage.set(null); 
  }

  let runtimePages = [
    'Logging',
    'Events',
    'Calls',
    'Dialog',
    'Browser',
    'File System',
    'Window',
  ];

</script>

<svelte:head>
  {@html css}
</svelte:head>

<div data-wails-drag class="page-wrapper with-sidebar" class:dark-mode="{$darkMode}" data-sidebar-type="full-height" >
  <!-- Sticky alerts (toasts), empty container -->
  <div class="sticky-alerts"></div>
  <!-- Sidebar -->
  <div class="sidebar noselect">
    <div data-wails-no-drag class="sidebar-menu">
      <!-- Sidebar brand -->
      <div on:click="{ homepageClicked }" class="sidebar-brand">        
      Wails Kitchen Sink
      </div>
      <!-- Sidebar links and titles -->
      <h5 class="sidebar-title">Runtime</h5>
      <div class="sidebar-divider"></div>
      {#each runtimePages as link}
        <span on:click="{linkClicked}" class="sidebar-link" class:active="{$selectedPage == link}">{link}</span> 
      {/each}
      <br />
      <h5 class="sidebar-title">Links</h5>
      <div class="sidebar-divider"></div>
      <span on:click="{linkClicked}" class="sidebar-link">Github</span>
      <span on:click="{linkClicked}" class="sidebar-link">Website</span>
    </div>
  </div>
  <!-- Content wrapper -->
    <div class="content-wrapper noselect" class:dark-content-wrapper="{$darkMode}">
      <div class="inner-content">
        <MainPage></MainPage>
      </div>
    </div>
</div>


<style global>
  @import 'halfmoon/css/halfmoon-variables.min.css';
  /* @import './assets/fonts/roboto.css'; */
  @import './App.css';
</style>
