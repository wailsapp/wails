import Vue from "vue";
import App from "./App.vue";

Vue.config.productionTip = false;

import Bridge from "./wailsbridge";

Bridge.Start(() => {
  new Vue({
    render: h => h(App)
  }).$mount("#app");
});
