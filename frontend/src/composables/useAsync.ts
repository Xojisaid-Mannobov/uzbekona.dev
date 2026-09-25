import { ref, shallowRef, watch, type Ref, type WatchSource } from 'vue'
import { toApiError, type ApiError } from '@/services/http'

/**
 * Ma'lumot yuklashning barcha holatlari bitta joyda: loading / error / data.
 * Har bir sahifa skeleton, xato va bo'sh holatlarni shu orqali ko'rsatadi.
 */
export function useAsync<T>(fetcher: () => Promise<T>, options: { watch?: WatchSource[]; immediate?: boolean } = {}) {
  const data = shallowRef<T | null>(null) as Ref<T | null>
  const error = ref<ApiError | null>(null)
  const loading = ref(false)
  let requestId = 0

  async function run() {
    const id = ++requestId
    loading.value = true
    error.value = null
    try {
      const result = await fetcher()
      // Faqat eng oxirgi so'rov natijasi qabul qilinadi (race condition'dan himoya)
      if (id === requestId) data.value = result
    } catch (e) {
      if (id === requestId) error.value = toApiError(e)
    } finally {
      if (id === requestId) loading.value = false
    }
  }

  if (options.watch?.length) watch(options.watch, run)
  if (options.immediate !== false) run()

  return { data, error, loading, reload: run }
}
