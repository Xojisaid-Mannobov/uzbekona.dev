import type { Router } from 'vue-router'

// Cookie'siz tashriflar hisobi (o'z serverimizda, uchinchi tomon skriptlarisiz).
// Brauzerda faqat tasodifiy ID'lar saqlanadi: tashrif buyuruvchi — localStorage, sessiya — sessionStorage.
// Server ularning xeshini yozadi; IP va User-Agent bazaga tushmaydi. "Do Not Track" yoqilgan bo'lsa — hech narsa yuborilmaydi.

const endpoint = `${import.meta.env.VITE_API_URL ?? '/api/v1'}/track`

function randomId(): string {
  return crypto.randomUUID?.() ?? `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 14)}`
}

/** Saqlangan ID'ni oladi yoki yangisini yaratadi; `fresh` — hozirgina yaratilganmi. */
function storedId(getStorage: () => Storage, key: string): { id: string; fresh: boolean } {
  try {
    const storage = getStorage()
    const existing = storage.getItem(key)
    if (existing) return { id: existing, fresh: false }
    const id = randomId()
    storage.setItem(key, id)
    return { id, fresh: true }
  } catch {
    // Maxfiy rejim yoki bloklangan storage — server kunlik anonim xeshdan foydalanadi
    return { id: '', fresh: true }
  }
}

function send(body: string) {
  try {
    if (navigator.sendBeacon?.(endpoint, new Blob([body], { type: 'application/json' }))) return
  } catch {
    // sendBeacon ishlamasa — oddiy so'rov
  }
  fetch(endpoint, { method: 'POST', body, keepalive: true, credentials: 'omit', headers: { 'Content-Type': 'application/json' } }).catch(
    () => {},
  )
}

export function installAnalytics(router: Router) {
  if (typeof navigator === 'undefined' || navigator.doNotTrack === '1') return

  let lastPath = ''
  router.afterEach((to, _from, failure) => {
    // Faqat public sahifalar; admin, login va 404 hisobga olinmaydi
    if (failure || to.meta.layout !== 'public' || to.name === 'not-found') return
    // Faqat query (?category=…) o'zgarsa — yangi ko'rish emas
    if (to.path === lastPath) return
    lastPath = to.path

    const visitor = storedId(() => localStorage, 'uz_vid')
    const session = storedId(() => sessionStorage, 'uz_sid')
    send(
      JSON.stringify({
        path: to.path,
        referrer: session.fresh ? document.referrer : '',
        vid: visitor.id,
        sid: session.id,
        landing: session.fresh,
      }),
    )
  })
}
