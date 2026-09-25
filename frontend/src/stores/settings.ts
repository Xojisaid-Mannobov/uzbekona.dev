import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { publicApi } from '@/services/public'
import type { Metric, Seo, Settings, SiteSettings, SocialLink } from '@/types/api'

/** Sayt sozlamalari (kontaktlar, metrikalar, ijtimoiy tarmoqlar) — bir marta yuklanadi. */
export const useSettingsStore = defineStore('settings', () => {
  const site = ref<SiteSettings>({
    name: 'Uzbekona.dev',
    tagline: 'Digital Product Studio',
    email: '',
    phone: '',
    telegram: '',
    address: '',
    available: true,
  })
  const metrics = ref<Metric[]>([])
  const seo = ref<Seo>({ title: '', description: '' })
  const socials = ref<SocialLink[]>([])
  const loaded = ref(false)
  let pending: Promise<void> | null = null

  function apply(s: Settings) {
    site.value = { ...site.value, ...s.site }
    metrics.value = s.metrics
    seo.value = s.seo
    socials.value = s.socials
    loaded.value = true
  }

  function load(force = false) {
    if ((loaded.value && !force) || pending) return pending ?? Promise.resolve()
    pending = publicApi
      .settings()
      .then(apply)
      .catch(() => {
        /* sozlamalar yuklanmasa standart qiymatlar qoladi */
      })
      .finally(() => {
        pending = null
      })
    return pending
  }

  const telegramUrl = computed(() => {
    const t = site.value.telegram.trim()
    if (!t) return ''
    return t.startsWith('http') ? t : `https://t.me/${t.replace(/^@/, '')}`
  })

  return { site, metrics, seo, socials, loaded, load, apply, telegramUrl }
})
