import type { RouteRecordRaw } from 'vue-router'

export type AuthContext = 'anonymous' | 'user'

export interface NavItem {
  id: string
  label: string
  icon: string
  path: string
}

export interface AppDefinition {
  id: string
  name: string
  routes: RouteRecordRaw[]
  navItems: NavItem[]
}

export function defineApp(def: AppDefinition): AppDefinition {
  return def
}