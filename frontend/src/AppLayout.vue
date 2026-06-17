<template>
  <div class="app-layout">
    <aside class="sidebar">
      <div class="sidebar-logo">📦 MiniBox</div>
      <nav class="sidebar-nav">
        <RouterLink
          v-for="item in navItems"
          :key="item.id"
          :to="item.path"
          class="nav-item"
          active-class="nav-item--active"
        >
          <span>{{ item.icon }}</span>
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>
      <div class="sidebar-footer">
        <div class="user-info">
          <span class="user-name">{{ user?.name }}</span>
          <span class="user-email">{{ user?.email }}</span>
        </div>
        <button class="logout-btn" @click="authService.logout()">Sair</button>
      </div>
    </aside>
    <main class="main-content">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useAppsStore } from './pkg/stores/apps'
import { useUserStore } from './pkg/stores/user'
import { authService } from './runtime/services/authService'

const { navItems } = storeToRefs(useAppsStore())
const { user } = storeToRefs(useUserStore())
</script>
