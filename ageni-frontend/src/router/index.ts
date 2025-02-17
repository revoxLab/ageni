import {createRouter, createWebHistory} from 'vue-router'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/:w*',
      component: () => import('../main.vue'),
    },
  ],
})
