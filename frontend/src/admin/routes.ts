import type { RouteComponent, RouteRecordRaw } from 'vue-router'

type Loader = () => Promise<RouteComponent | { default: RouteComponent }>

// Admin panel marshrutlari — public sayt bilan bitta Vue ilova ichida, /admin prefiksi bilan
const page = (path: string, name: string, component: Loader, title: string): RouteRecordRaw => ({
  path,
  name,
  component,
  meta: { layout: 'admin', requiresAuth: true, title },
})

export const adminRoutes: RouteRecordRaw[] = [
  {
    path: '/admin/login',
    name: 'admin-login',
    component: () => import('@/admin/pages/LoginPage.vue'),
    meta: { layout: 'blank', title: 'Kirish' },
  },
  { path: '/admin', redirect: '/admin/dashboard' },
  page('/admin/dashboard', 'admin-dashboard', () => import('@/admin/pages/DashboardPage.vue'), 'Dashboard'),

  page('/admin/projects', 'admin-projects', () => import('@/admin/pages/ProjectsListPage.vue'), 'Loyihalar'),
  page('/admin/projects/new', 'admin-project-new', () => import('@/admin/pages/ProjectEditPage.vue'), 'Yangi loyiha'),
  page('/admin/projects/:id(\\d+)', 'admin-project-edit', () => import('@/admin/pages/ProjectEditPage.vue'), 'Loyihani tahrirlash'),

  page('/admin/services', 'admin-services', () => import('@/admin/pages/ServicesListPage.vue'), 'Xizmatlar'),
  page('/admin/services/new', 'admin-service-new', () => import('@/admin/pages/ServiceEditPage.vue'), 'Yangi xizmat'),
  page('/admin/services/:id(\\d+)', 'admin-service-edit', () => import('@/admin/pages/ServiceEditPage.vue'), 'Xizmatni tahrirlash'),

  page('/admin/team', 'admin-team', () => import('@/admin/pages/TeamListPage.vue'), 'Jamoa'),
  page('/admin/team/new', 'admin-member-new', () => import('@/admin/pages/TeamEditPage.vue'), 'Yangi a’zo'),
  page('/admin/team/:id(\\d+)', 'admin-member-edit', () => import('@/admin/pages/TeamEditPage.vue'), 'A’zoni tahrirlash'),

  page('/admin/labs', 'admin-labs', () => import('@/admin/pages/LabsListPage.vue'), 'Uzbekona Labs'),
  page('/admin/labs/new', 'admin-lab-new', () => import('@/admin/pages/LabEditPage.vue'), 'Yangi Labs loyihasi'),
  page('/admin/labs/:id(\\d+)', 'admin-lab-edit', () => import('@/admin/pages/LabEditPage.vue'), 'Labs loyihasini tahrirlash'),

  page('/admin/analytics', 'admin-analytics', () => import('@/admin/pages/AnalyticsPage.vue'), 'Statistika'),

  page('/admin/news', 'admin-news', () => import('@/admin/pages/NewsListPage.vue'), 'Yangiliklar'),
  page('/admin/news/new', 'admin-news-new', () => import('@/admin/pages/NewsEditPage.vue'), 'Yangi yangilik'),
  page('/admin/news/:id(\\d+)', 'admin-news-edit', () => import('@/admin/pages/NewsEditPage.vue'), 'Yangilikni tahrirlash'),

  page('/admin/articles', 'admin-articles', () => import('@/admin/pages/ArticlesListPage.vue'), 'Maqolalar'),
  page('/admin/articles/new', 'admin-article-new', () => import('@/admin/pages/ArticleEditPage.vue'), 'Yangi maqola'),
  page('/admin/articles/:id(\\d+)', 'admin-article-edit', () => import('@/admin/pages/ArticleEditPage.vue'), 'Maqolani tahrirlash'),

  page('/admin/media', 'admin-media', () => import('@/admin/pages/MediaPage.vue'), 'Media'),
  page('/admin/requests', 'admin-requests', () => import('@/admin/pages/RequestsPage.vue'), 'So‘rovlar'),
  { path: '/admin/messages', redirect: '/admin/requests' },
  page('/admin/settings', 'admin-settings', () => import('@/admin/pages/SettingsPage.vue'), 'Sozlamalar'),
  page('/admin/users', 'admin-users', () => import('@/admin/pages/UsersPage.vue'), 'Adminlar'),
]
