import { useUserStore } from '@/store/user'

// v-permission="'user:add'" hides the element when the user lacks the permission.
export const permissionDirective = {
  mounted(el, binding) {
    const store = useUserStore()
    const need = binding.value
    const ok = Array.isArray(need)
      ? need.some((p) => store.permissions.includes(p))
      : store.permissions.includes(need)
    if (!ok && el.parentNode) el.parentNode.removeChild(el)
  }
}
