<template>
  <div class="files-page">
    <header class="page-header">
      <h2>{{ currentPath === '/' ? 'Meus Arquivos' : currentPath }}</h2>
    </header>

    <nav class="breadcrumb" v-if="currentPath !== '/'">
      <RouterLink to="/files">Início</RouterLink>
      <span v-for="(segment, i) in pathSegments" :key="i">
        / <RouterLink :to="`/files/${pathSegments.slice(0, i + 1).join('/')}`">
          {{ segment }}
        </RouterLink>
      </span>
    </nav>

    <!-- Upload -->
    <div class="upload-bar">
      <input type="file" ref="fileInput" @change="handleUpload" style="display:none" />
      <button @click="fileInput?.click()" :disabled="uploading">
        {{ uploading ? 'Enviando...' : '⬆ Upload' }}
      </button>
    </div>

    <p v-if="loading" class="empty-state">Carregando...</p>
    <p v-else-if="error" class="empty-state error-msg">{{ error }}</p>

    <div v-else class="file-list">
      <div
        v-for="item in items"
        :key="item.name"
        class="file-row"
        :class="{ 'file-row--folder': item.type === 'folder' }"
        @click="item.type === 'folder' ? navigate(item.name) : null"
      >
        <span class="file-row__icon">
          {{ item.type === 'folder' ? '📁' : fileIcon(item.mimeType) }}
        </span>
        <span class="file-row__name">{{ item.name }}</span>
        <span class="file-row__size">{{ formatSize(item.size) }}</span>
        <a
          v-if="item.type === 'file'"
          :href="downloadUrl(item.name)"
          class="file-row__download"
          @click.stop
        >⬇</a>
      </div>
      <p v-if="items.length === 0" class="empty-state">Pasta vazia</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { filesApi } from '../../pkg/api/client'
import { useFiles, fileIcon, formatSize } from './composables/useFiles'

const route = useRoute()
const router = useRouter()

const currentPath = computed(() => {
  const p = route.params.folderPath as string
  return p ? `/${p}` : '/'
})

const pathSegments = computed(() =>
  currentPath.value.split('/').filter(Boolean)
)

const { items, loading, error, reload } = useFiles(currentPath)

function navigate(folderName: string) {
  const base = currentPath.value === '/' ? '' : currentPath.value
  router.push(`/files${base}/${folderName}`)
}

function downloadUrl(fileName: string) {
  const base = currentPath.value === '/' ? '' : currentPath.value
  return filesApi.downloadUrl(`${base}/${fileName}`)
}

// Upload
const fileInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

async function handleUpload(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  uploading.value = true
  try {
    await filesApi.upload(currentPath.value, file)
    await reload()
  } finally {
    uploading.value = false
  }
}
</script>
