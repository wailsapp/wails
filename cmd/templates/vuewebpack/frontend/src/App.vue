<template>
  <div id="app">
    <img alt="Wails logo" src="./assets/images/logo.png" class="logo zoomIn">
    <HelloWorld msg="Welcome to Your Wails App!"/>
    <Quote/>
  </div>
</template>

<script>
import HelloWorld from "./components/HelloWorld.vue";
import Quote from "./components/Quote.vue";
import "./assets/css/main.css";

export default {
  name: "app",
  components: {
    HelloWorld,
    Quote
  }
};
</script>

