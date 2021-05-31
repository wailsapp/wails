import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

// import wails runtime
import * as Wails from '@wailsapp/runtime';

Wails.Init(() => {
  createApp(App)
    .use(store)
    .use(router)
    .mount('#app')
})
