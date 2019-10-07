import Vue from 'vue';
import App from './App.vue';

Vue.config.productionTip = false;
Vue.config.devtools = true;

import Wails from '@wailsapp/runtime';

Wails.Init(() => {
	new Vue({
		render: h => h(App)
	}).$mount('#app');
});
