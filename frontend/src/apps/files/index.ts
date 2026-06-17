import type { RouteRecordRaw } from 'vue-router'
import { defineApp } from '../../pkg/types'
import FilesPage from './FilesPage.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/files/:folderPath(.*)?',
    name: 'files',
    component: FilesPage,
    meta: { authContext: 'user' },
  },
]

export default defineApp({
  id: 'files',
  name: 'Files',
  routes,
  navItems: [{ id: 'files', label: 'Arquivos', icon: '📁', path: '/files' }],
})