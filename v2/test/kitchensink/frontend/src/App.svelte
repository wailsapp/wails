
<script>

  import runtime from '@wailsapp/runtime2';
  import { selectedPage } from './Store';
  import MainPage from './MainPage.svelte';
 
  // Handle Dark/Light themes automatically
  var darkMode = runtime.System.DarkModeEnabled();
  runtime.System.OnThemeChange( (isDarkMode) => {
    darkMode = isDarkMode;
  });

  // Hightlight CSS
  import { atomOneDark, atomOneLight } from "svelte-highlight/styles";
  $: css = darkMode ? atomOneDark : atomOneLight;

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

<div data-wails-drag class="page-wrapper with-sidebar" class:dark-mode="{darkMode}" data-sidebar-type="full-height" >
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
  <div class="content-wrapper noselect" class:dark-content-wrapper="{darkMode}">
    <MainPage></MainPage>
  </div>
</div>


<style global>
  @import 'halfmoon/css/halfmoon-variables.min.css';

  :root {
    --lm-base-body-bg-color: #0000;
    --dm-base-body-bg-color: #0000;
    --lm-sidebar-bg-color: #0000;
    --dm-sidebar-bg-color: #0000;
    --dm-sidebar-link-text-color: white;
    --dm-sidebar-link-text-color-hover: rgb(255, 214, 0);
    --lm-sidebar-link-text-color: black;
    --lm-sidebar-link-text-color-hover: rgb(158, 158, 255);

    --dm-sidebar-link-text-color-active: rgb(255, 214, 0);
    --dm-sidebar-link-text-color-active-hover: rgb(255, 214, 0);

    --sidebar-title-font-size: 1.75rem;
    --sidebar-brand-font-size: 2.3rem;

    /* Switch */
    --dm-switch-bg-color: rgb(28,173,213);
    --lm-switch-bg-color: rgb(28,173,213);
    --dm-switch-bg-color-checked: rgb(239,218,91);
    --lm-switch-bg-color-checked: rgb(239,218,91);
    --lm-switch-slider-bg-color: #FFF;
    --dm-switch-slider-bg-color: #FFF;
  }

  .sidebar-link {
    font-weight: bold;
    cursor: pointer;
  }

  .content-wrapper {
    background-color: #eee;
  }

  .dark-content-wrapper {
    background-color: #25282c;
  }
/* 
  .sidebar {
    background-color: #0000;
  } */

  .sidebar-brand {
    padding-top: 35px;
    padding-bottom: 25px;
    cursor: pointer;
  }

  .sidebar-menu {
    background-color: #0000;
  }

  /* Credit: https://stackoverflow.com/a/4407335 */
  .noselect {
                   cursor: default;
    -webkit-touch-callout: none; /* iOS Safari */
      -webkit-user-select: none; /* Safari */
       -khtml-user-select: none; /* Konqueror HTML */
         -moz-user-select: none; /* Old versions of Firefox */
          -ms-user-select: none; /* Internet Explorer/Edge */
              user-select: none; /* Non-prefixed version, currently
                                    supported by Chrome, Edge, Opera and Firefox */
  }
</style>
