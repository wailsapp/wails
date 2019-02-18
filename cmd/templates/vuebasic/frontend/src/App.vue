<template>
  <div id="app">
    <img alt="Wails logo" src="./assets/images/logo.png" class="logo zoomIn">
    <HelloWorld/>
  </div>
</template>

<script>
import HelloWorld from "./components/HelloWorld.vue";
import "./assets/css/main.css";

export default {
  name: "app",
  components: {
    HelloWorld
  }
};
</script>
