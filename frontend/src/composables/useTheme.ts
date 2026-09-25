import { ref } from 'vue'

export type ThemeMode = 'system' | 'light' | 'dark'

const STORAGE_KEY = 'uzb-theme'

function read(): ThemeMode {
  try {
    const v = localStorage.getItem(STORAGE_KEY)
    return v === 'light' || v === 'dark' ? v : 'system'
  } catch {
    return 'system'
  }
}

const mode = ref<ThemeMode>(read())

function apply(value: ThemeMode) {
  const root = document.documentElement
  if (value === 'system') delete root.dataset.theme
  else root.dataset.theme = value
}

/** Light / dark / system mavzusini boshqaradi (tanlov brauzerda saqlanadi). */
export function useTheme() {
  function set(value: ThemeMode) {
    mode.value = value
    apply(value)
    try {
      if (value === 'system') localStorage.removeItem(STORAGE_KEY)
      else localStorage.setItem(STORAGE_KEY, value)
    } catch {
      /* private rejim — tanlov faqat shu sessiyada */
    }
  }

  function isDark() {
    if (mode.value === 'dark') return true
    if (mode.value === 'light') return false
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  }

  function toggle() {
    set(isDark() ? 'light' : 'dark')
  }

  return { mode, set, toggle, isDark }
}
