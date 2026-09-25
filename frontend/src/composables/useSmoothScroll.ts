import type Lenis from 'lenis'
import type { gsap as GsapType } from 'gsap'
import type { ScrollTrigger as ScrollTriggerType } from 'gsap/ScrollTrigger'

/**
 * Harakat kutubxonalari (GSAP, ScrollTrigger, Lenis) alohida chunk'da va kechiktirib yuklanadi —
 * birinchi ekran (hero) ularsiz, faqat CSS bilan chiziladi. Bu Lighthouse Performance uchun muhim.
 */
export interface Motion {
  gsap: typeof GsapType
  ScrollTrigger: typeof ScrollTriggerType
}

let motionPromise: Promise<Motion> | null = null
let lenis: Lenis | null = null
let tick: ((time: number) => void) | null = null

export const reducedMotion = () => window.matchMedia('(prefers-reduced-motion: reduce)').matches

export function loadMotion(): Promise<Motion> {
  motionPromise ??= Promise.all([import('gsap'), import('gsap/ScrollTrigger')]).then(([g, st]) => {
    g.gsap.registerPlugin(st.ScrollTrigger)
    return { gsap: g.gsap, ScrollTrigger: st.ScrollTrigger }
  })
  return motionPromise
}

/** Lenis smooth scroll — faqat public sayt, sichqoncha bilan va harakat cheklanmagan bo'lsa. */
export async function startSmoothScroll() {
  if (lenis || reducedMotion() || window.matchMedia('(pointer: coarse)').matches) return
  const [{ default: LenisCtor }, { gsap, ScrollTrigger }] = await Promise.all([import('lenis'), loadMotion()])
  if (lenis) return
  lenis = new LenisCtor({ lerp: 0.11, wheelMultiplier: 0.95 })
  lenis.on('scroll', ScrollTrigger.update)
  tick = (time: number) => lenis?.raf(time * 1000)
  gsap.ticker.add(tick)
  gsap.ticker.lagSmoothing(0)
}

export async function stopSmoothScroll() {
  if (!lenis) return
  const { gsap } = await loadMotion()
  if (tick) gsap.ticker.remove(tick)
  lenis.destroy()
  lenis = null
}

export function scrollToTarget(target: string | HTMLElement) {
  if (lenis) lenis.scrollTo(target, { offset: -80, duration: 1.2 })
  else {
    const el = typeof target === 'string' ? document.querySelector(target) : target
    el?.scrollIntoView({ behavior: reducedMotion() ? 'auto' : 'smooth' })
  }
}

/** Brauzer bo'sh paytida bajarish (asosiy kontent chizilgandan keyin). */
export function whenIdle(fn: () => void) {
  if ('requestIdleCallback' in window) window.requestIdleCallback(fn, { timeout: 2000 })
  else setTimeout(fn, 600)
}
