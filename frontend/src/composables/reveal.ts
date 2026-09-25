import type { Directive } from 'vue'

/**
 * v-reveal — element ekranga kirganda yumshoq paydo bo'ladi.
 * Faqat bloklar (sarlavha, rasm, karta) uchun — har bir paragraf alohida uchib chiqmaydi.
 *
 *   <div v-reveal>…</div>
 *   <div v-reveal="{ delay: 120, variant: 'scale' }">…</div>
 */
interface RevealOptions {
  delay?: number
  variant?: 'up' | 'fade' | 'scale'
}

let observer: IntersectionObserver | null = null

function getObserver() {
  if (observer) return observer
  observer = new IntersectionObserver(
    (entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          entry.target.classList.add('is-revealed')
          observer?.unobserve(entry.target)
        }
      }
    },
    { rootMargin: '0px 0px -8% 0px', threshold: 0.08 },
  )
  return observer
}

export const vReveal: Directive<HTMLElement, RevealOptions | number | undefined> = {
  mounted(el, binding) {
    const opts: RevealOptions = typeof binding.value === 'number' ? { delay: binding.value } : (binding.value ?? {})
    el.dataset.reveal = opts.variant ?? 'up'
    if (opts.delay) el.style.setProperty('--reveal-delay', `${opts.delay}ms`)

    if (typeof IntersectionObserver === 'undefined') {
      el.classList.add('is-revealed')
      return
    }
    getObserver().observe(el)
  },
  unmounted(el) {
    observer?.unobserve(el)
  },
}
