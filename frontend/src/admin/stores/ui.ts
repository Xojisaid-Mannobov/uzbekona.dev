import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Toast {
  id: number
  tone: 'success' | 'error' | 'info'
  message: string
}

interface ConfirmState {
  title: string
  text?: string
  confirmLabel: string
  danger: boolean
  resolve: (ok: boolean) => void
}

/** Admin UI: bildirishnomalar (toast) va tasdiqlash oynasi. */
export const useUiStore = defineStore('admin-ui', () => {
  const toasts = ref<Toast[]>([])
  const confirmState = ref<ConfirmState | null>(null)
  let seq = 0

  function toast(message: string, tone: Toast['tone'] = 'success', ttl = 3500) {
    const id = ++seq
    toasts.value.push({ id, tone, message })
    setTimeout(() => dismiss(id), ttl)
  }

  function dismiss(id: number) {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }

  function confirm(opts: { title: string; text?: string; confirmLabel?: string; danger?: boolean }): Promise<boolean> {
    confirmState.value?.resolve(false)
    return new Promise((resolve) => {
      confirmState.value = {
        title: opts.title,
        text: opts.text,
        confirmLabel: opts.confirmLabel ?? 'Ha, o‘chirish',
        danger: opts.danger ?? true,
        resolve,
      }
    })
  }

  function closeConfirm(ok: boolean) {
    confirmState.value?.resolve(ok)
    confirmState.value = null
  }

  return {
    toasts,
    confirmState,
    toast,
    success: (m: string) => toast(m, 'success'),
    error: (m: string) => toast(m, 'error', 5000),
    dismiss,
    confirm,
    closeConfirm,
  }
})
