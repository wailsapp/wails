import { createRouter, createMemoryHistory } from 'vue-router'
import Home from '../views/Home.vue'
import About from '../views/About.vue'


const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    // You can only use pre-loading to add routes, not the on-demand loading method.
    component: About
  }
]

const router = createRouter({
  history: createMemoryHistory(),
  routes
})

export default router
