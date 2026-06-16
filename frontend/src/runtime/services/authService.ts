import type { Router } from 'vue-router'
import { authApi, graphApi } from '../../pkg/api/client'
import { useAuthStore } from '../../pkg/stores/auth'
import { useUserStore } from '../../pkg/stores/user'

class AuthService {
  private router: Router | null = null

  init(router: Router) {
    this.router = router
    this.restoreSession()
  }

  // Se já há token salvo, busca o perfil do usuário no backend
  private async restoreSession() {
    const authStore = useAuthStore()
    if (!authStore.accessToken) return

    try {
      const me = await graphApi.me()
      const userStore = useUserStore()
      userStore.setUser({
        username: me.onPremisesSamAccountName,
        name: me.displayName,
        email: me.mail,
      })
    } catch {
      // Token inválido — limpa sessão
      authStore.setToken(null)
    }
  }

  async login(username: string, password: string) {
    const data = await authApi.login(username, password)

    const authStore = useAuthStore()
    const userStore = useUserStore()
    authStore.setToken(data.token)
    userStore.setUser({ username, name: data.name, email: data.email })
  }

  logout() {
    const authStore = useAuthStore()
    const userStore = useUserStore()
    authStore.setToken(null)
    userStore.reset()
    this.router?.push({ name: 'login' })
  }
}

export const authService = new AuthService()