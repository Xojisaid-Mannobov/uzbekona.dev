import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { adminRoutes } from '@/admin/routes'
import { http } from '@/services/http'

declare module 'vue-router' {
  interface RouteMeta {
    /** Qaysi layout ishlatiladi: public sayt, admin panel yoki bo'sh (login) */
    layout?: 'public' | 'admin' | 'blank'
    /** Footer oldidagi katta CTA blokini yashirish */
    hideCta?: boolean
    /** Admin sahifasi — autentifikatsiya talab qiladi */
    requiresAuth?: boolean
    title?: string
  }
}

// Sahifalar lazy-load qilinadi — har biri alohida chunk
const publicRoutes: RouteRecordRaw[] = [
  { path: '/', name: 'home', component: () => import('@/pages/HomePage.vue') },
  { path: '/projects', name: 'projects', component: () => import('@/pages/ProjectsPage.vue') },
  { path: '/projects/:slug', name: 'project', component: () => import('@/pages/ProjectDetailPage.vue') },
  { path: '/services', name: 'services', component: () => import('@/pages/ServicesPage.vue') },
  { path: '/services/:slug', name: 'service', component: () => import('@/pages/ServiceDetailPage.vue') },
  { path: '/about', name: 'about', component: () => import('@/pages/AboutPage.vue') },
  { path: '/team', name: 'team', component: () => import('@/pages/TeamPage.vue') },
  { path: '/journal', name: 'journal', component: () => import('@/pages/JournalPage.vue') },
  { path: '/journal/:slug', name: 'article', component: () => import('@/pages/ArticlePage.vue') },
  { path: '/news', name: 'news', component: () => import('@/pages/NewsPage.vue') },
  { path: '/news/:slug', name: 'news-item', component: () => import('@/pages/NewsDetailPage.vue') },
  { path: '/join', name: 'join', component: () => import('@/pages/JoinPage.vue'), meta: { hideCta: true } },
  { path: '/contact', name: 'contact', component: () => import('@/pages/ContactPage.vue'), meta: { hideCta: true } },
].map((r) => ({ ...r, meta: { layout: 'public' as const, ...r.meta } }))

export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    ...publicRoutes,
    ...adminRoutes,
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/pages/NotFoundPage.vue'),
      meta: { layout: 'public', hideCta: true },
    },
  ],
  scrollBehavior(to, _from, saved) {
    if (saved) return saved
    if (to.hash) return { el: to.hash, top: 96 }
    return { top: 0 }
  },
})

// ─── Admin himoyasi ─────────────────────────────────────
// Auth kodi faqat /admin sahifalariga kirilganda yuklanadi (public bundle'ga tushmaydi)
let adminInterceptor = false

router.beforeEach(async (to) => {
  if (!to.path.startsWith('/admin')) return

  const { useAuthStore } = await import('@/admin/stores/auth')
  const auth = useAuthStore()

  if (!adminInterceptor) {
    adminInterceptor = true
    // Sessiya muddati tugasa (401) — login sahifasiga qaytarish
    http.interceptors.response.use(undefined, (err) => {
      const route = router.currentRoute.value
      if (err?.status === 401 && route.meta.requiresAuth) {
        auth.expire()
        router.push({ name: 'admin-login', query: { redirect: route.fullPath } })
      }
      return Promise.reject(err)
    })
  }

  const ok = await auth.check()
  if (to.meta.requiresAuth && !ok) return { name: 'admin-login', query: { redirect: to.fullPath } }
  if (to.name === 'admin-login' && ok) return { name: 'admin-dashboard' }
})

router.afterEach((to) => {
  if (to.meta.layout === 'admin' || to.meta.layout === 'blank') {
    document.title = `${to.meta.title ?? 'Admin'} — Uzbekona.dev Admin`
  }
})
