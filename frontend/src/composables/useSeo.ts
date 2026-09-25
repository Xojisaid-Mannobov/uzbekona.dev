import { onBeforeUnmount, toValue, watchEffect, type MaybeRefOrGetter } from 'vue'
import { useSettingsStore } from '@/stores/settings'

interface SeoInput {
  title?: string | null
  description?: string | null
  image?: string | null
  type?: 'website' | 'article'
  noindex?: boolean
}

function setMeta(attr: 'name' | 'property', key: string, content: string | null | undefined) {
  let el = document.head.querySelector<HTMLMetaElement>(`meta[${attr}="${key}"]`)
  if (!content) {
    el?.remove()
    return
  }
  if (!el) {
    el = document.createElement('meta')
    el.setAttribute(attr, key)
    document.head.appendChild(el)
  }
  el.content = content
}

function setCanonical(href: string) {
  let el = document.head.querySelector<HTMLLinkElement>('link[rel="canonical"]')
  if (!el) {
    el = document.createElement('link')
    el.rel = 'canonical'
    document.head.appendChild(el)
  }
  el.href = href
}

/** Sahifa sarlavhasi, description va Open Graph teglarini yangilaydi. */
export function useSeo(input: MaybeRefOrGetter<SeoInput>) {
  const settings = useSettingsStore()

  const stop = watchEffect(() => {
    const seo = toValue(input)
    const siteName = settings.site.name || 'Uzbekona.dev'
    const defaultTitle = settings.seo.title || `${siteName} — Digital Product Studio`
    const title = seo.title ? `${seo.title} — ${siteName}` : defaultTitle
    const description = seo.description || settings.seo.description
    const url = window.location.origin + window.location.pathname

    document.title = title
    setMeta('name', 'description', description)
    setMeta('property', 'og:title', title)
    setMeta('property', 'og:description', description)
    setMeta('property', 'og:type', seo.type ?? 'website')
    setMeta('property', 'og:url', url)
    setMeta('property', 'og:image', seo.image ? new URL(seo.image, window.location.origin).href : `${window.location.origin}/og-image.png`)
    setMeta('name', 'robots', seo.noindex ? 'noindex, nofollow' : null)
    setCanonical(url)
  })

  onBeforeUnmount(stop)
}
