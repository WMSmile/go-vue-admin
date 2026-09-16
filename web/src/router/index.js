import { createRouter, createWebHistory } from 'vue-router'
import Layout from '@/layout/index.vue'
import { useUserStore } from '@/store/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { public: true }
  },
  {
    path: '/',
    name: 'Layout',
    component: Layout,
    children: [
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/system/profile/index.vue'),
        meta: { title: '个人中心' }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/system/settings/index.vue'),
        meta: { title: '设置' }
      }
    ]
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { public: true }
  },
  // Render the 404 view INLINE for unknown paths instead of redirecting to
  // /404. A redirect rewrites the address bar to /404, which poisons the URL so
  // that a later refresh keeps loading /404 even though the real page (e.g.
  // /system/config) resolves fine once its route is registered.
  { path: '/:pathMatch(.*)*', name: 'CatchAll', component: () => import('@/views/error/404.vue'), meta: { public: true } }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Dynamic routes are registered in main.js BEFORE the router is installed, so
// the very first navigation — including refreshes / deep links like
// /system/role — already sees the full route table. This guard enforces auth
// and, as a safety net, re-resolves the URL once if it still fell through to
// the catch-all (e.g. the cached menu lacked the route).
let catchAllRetried = false
router.beforeEach(async (to) => {
  const token = localStorage.getItem('token')
  if (!to.meta.public && !token) return '/login'

  // If we have a token but the target did not match any dynamic route (it hit
  // the catch-all), the routes weren't ready. Wait for the fresh menu, then
  // re-navigate to the SAME url so it resolves to the real page.
  if (token && !catchAllRetried && to.matched.some((r) => r.name === 'CatchAll')) {
    catchAllRetried = true
    await useUserStore().rehydrate()
    return { path: to.path, query: to.query, hash: to.hash, replace: true }
  }
  return true
})

export default router
