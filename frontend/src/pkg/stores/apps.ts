import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { AppDefinition, NavItem } from '../types'

export const useAppsStore = defineStore('apps', () => {
  const registeredApps = ref<AppDefinition[]>([])
  const navItems = computed<NavItem[]>(() =>
    registeredApps.value.flatMap((app) => app.navItems)
  )

  function register(app: AppDefinition) {
    registeredApps.value.push(app)
  }

  return { registeredApps, navItems, register }
})