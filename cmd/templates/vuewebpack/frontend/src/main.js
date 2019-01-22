import Vue from "vue";
export const eventBus = new Vue();

import App from "./App.vue";
new Vue({
  render: h => h(App)
}).$mount("#app");

import Bridge from "./wailsbridge";
Bridge.Start(startApp);

function startApp() {
  eventBus.$emit("ready");
}
