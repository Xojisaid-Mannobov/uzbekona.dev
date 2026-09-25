import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { adminApi } from '@/admin/services/adminApi'
import type { Admin } from '@/types/api'

/**
 * Admin sessiyasi. Token HttpOnly cookie'da — JavaScript uni ko'rmaydi,
 * shuning uchun sessiya holati /admin/auth/me orqali aniqlanadi.
 */
export const useAuthStore = defineStore('admin-auth', () => {
  const admin = ref<Admin | null>(null)
  const checked = ref(false)
  let pending: Promise<boolean> | null = null

  const isAuthenticated = computed(() => !!admin.value)

  function check(): Promise<boolean> {
    if (checked.value) return Promise.resolve(isAuthenticated.value)
    pending ??= adminApi.auth
      .me()
      .then((a) => {
        admin.value = a
        return true
      })
      .catch(() => {
        admin.value = null
        return false
      })
      .finally(() => {
        checked.value = true
        pending = null
      })
    return pending
  }

  async function login(email: string, password: string) {
    admin.value = await adminApi.auth.login(email, password)
    checked.value = true
  }

  async function logout() {
    try {
      await adminApi.auth.logout()
    } finally {
      admin.value = null
    }
  }

  /** 401 kelganda (sessiya tugagan) chaqiriladi. */
  function expire() {
    admin.value = null
  }

  return { admin, checked, isAuthenticated, check, login, logout, expire }
})
