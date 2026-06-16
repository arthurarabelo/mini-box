import { ref } from 'vue'
import { defineStore } from 'pinia'

export interface UserProfile {
  username: string
  name: string
  email: string
}

export const useUserStore = defineStore('user', () => {
  const user = ref<UserProfile | null>(null)

  function setUser(profile: UserProfile) { user.value = profile }
  function reset() { user.value = null }

  return { user, setUser, reset }
})