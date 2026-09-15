import { defineStore } from 'pinia'
import { login as loginApi } from '@/api'
import router from '@/router'
import RouteView from '@/layout/RouteView.vue'

// Strip leading slash so the path is relative to its parent route (Layout = '/').
const rel = (p) => (p && p.startsWith('/') ? p.slice(1) : p)

// Statically collect all view modules so Vite can resolve lazy imports.
// A fully-dynamic `import(\`@/views/${x}.vue\`)` cannot be analyzed by Vite and
// would fail at runtime, breaking the post-login redirect.
const viewModules = import.meta.glob('../views/**/*.vue')

// Build Vue Router routes from the backend menu tree (catalog -> nested, menu -> page).
export function genRoutes(menus) {
  return (menus || []).map((m) => {
    if (m.type === 1) {
      return {
        path: rel(m.path),
        component: RouteView,
        meta: { title: m.title, icon: m.icon },
        children: genRoutes(m.children || [])
      }
    }
    const loader = viewModules[`../views/${m.component}.vue`]
    return {
      path: rel(m.path),
      name: m.name,
      component: loader || viewModules['../views/error/404.vue'],
      meta: { title: m.title, icon: m.icon }
    }
  })
}

function firstLeaf(menus, parent = '') {
  for (const m of menus || []) {
    const path = m.path && m.path.startsWith('/') ? m.path : parent ? `${parent}/${m.path}` : m.path
    if (m.type === 2) return path
    if (m.children && m.children.length) {
      const p = firstLeaf(m.children, path)
      if (p) return p
    }
  }
  return null
}

// Prefer the dashboard (top-level page) as the landing route after login.
function findDashboard(menus, parent = '') {
  for (const m of menus || []) {
    const path = m.path && m.path.startsWith('/') ? m.path : parent ? `${parent}/${m.path}` : m.path
    if (m.type === 2 && (m.name === 'Dashboard' || m.path === '/dashboard')) return path
    if (m.children && m.children.length) {
      const p = findDashboard(m.children, path)
      if (p) return p
    }
  }
  return null
}

export const useUserStore = defineStore('user', {
  state: () => ({
    token: '',
    userInfo: null,
    menus: [],
    permissions: [],
    roles: []
  }),
  actions: {
    async login(form) {
      const data = await loginApi(form)
      this.token = data.token
      this.userInfo = data.user
      this.menus = data.menus
      this.permissions = data.permissions
      this.roles = (data.user?.roles || []).map((r) => r.keyword)
      localStorage.setItem('token', data.token)
      localStorage.setItem('menus', JSON.stringify(data.menus))
      localStorage.setItem('permissions', JSON.stringify(data.permissions))
      localStorage.setItem('userInfo', JSON.stringify(data.user))
      genRoutes(data.menus).forEach((r) => router.addRoute('Layout', r))
      const fp = findDashboard(data.menus) || firstLeaf(data.menus)
      try {
        await router.push(fp || '/')
      } catch (e) {
        // navigation may be superseded; ignore
      }
      return data
    },
    rehydrate() {
      const token = localStorage.getItem('token')
      const menus = JSON.parse(localStorage.getItem('menus') || '[]')
      if (token && menus.length) {
        this.token = token
        this.menus = menus
        this.permissions = JSON.parse(localStorage.getItem('permissions') || '[]')
        this.userInfo = JSON.parse(localStorage.getItem('userInfo') || 'null')
        this.roles = (this.userInfo?.roles || []).map((r) => r.keyword)
        genRoutes(menus).forEach((r) => router.addRoute('Layout', r))
      }
    },
    logout() {
      this.token = ''
      this.userInfo = null
      this.menus = []
      this.permissions = []
      this.roles = []
      localStorage.removeItem('token')
      localStorage.removeItem('menus')
      localStorage.removeItem('permissions')
      localStorage.removeItem('userInfo')
    }
  }
})
