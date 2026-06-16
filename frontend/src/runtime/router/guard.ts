import type { Router } from 'vue-router'
import { useAuthStore } from '../../pkg/stores/auth'

export function setupAuthGuard(router: Router) {
  router.beforeEach((to) => {
    const authContext = (to.meta.authContext as string) ?? 'user'
    if (authContext === 'anonymous') return true

    const authStore = useAuthStore()
    if (!authStore.userContextReady) {
      return { name: 'login', query: { redirectUrl: to.fullPath } }
    }
    return true
  })
}