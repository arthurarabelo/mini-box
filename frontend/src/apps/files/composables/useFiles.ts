import { ref, watch } from 'vue'
import type { Ref } from 'vue'
import { filesApi, type FileItem } from '../../../pkg/api/client'

export { type FileItem }

export function fileIcon(mimeType?: string): string {
  if (!mimeType) return '📄'
  if (mimeType.startsWith('image/')) return '🖼️'
  if (mimeType === 'text/markdown') return '📝'
  if (mimeType === 'text/csv') return '📊'
  if (mimeType.includes('openxmlformats')) return '📃'
  return '📄'
}

export function formatSize(bytes?: number): string {
  if (!bytes) return ''
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

export function useFiles(path: Ref<string>) {
  const items = ref<FileItem[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function load() {
    loading.value = true
    error.value = null
    try {
      items.value = await filesApi.list(path.value)
    } catch (e: any) {
      error.value = e.message
    } finally {
      loading.value = false
    }
  }

  // Recarrega quando o path muda (ex: navegação entre pastas)
  watch(path, load, { immediate: true })

  return { items, loading, error, reload: load }
}
