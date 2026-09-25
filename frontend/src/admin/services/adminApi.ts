import { http } from '@/services/http'
import type {
  Admin,
  Analytics,
  ApiResponse,
  Article,
  ArticleCategory,
  Block,
  ContactRequest,
  ContactStatus,
  ContentStatus,
  DashboardStats,
  GalleryItem,
  Lab,
  LabStage,
  Media,
  Metric,
  News,
  PageMeta,
  Project,
  Seo,
  Service,
  Settings,
  SiteSettings,
  SocialLink,
  SocialRef,
  TeamMember,
} from '@/types/api'

// ─── Kiruvchi (form) turlari ────────────────────────────

export interface ProjectInput {
  title: string
  slug: string
  tagline: string
  short_description: string
  full_description: string
  cover_id: number | null
  year: number | null
  client: string
  industry: string
  platforms: string[]
  services: string[]
  stack: string[]
  metrics: Metric[]
  live_url: string
  accent: string
  status: ContentStatus
  featured: boolean
  seo: Seo
  gallery: Pick<GalleryItem, 'media_id' | 'caption'>[]
  blocks: Block[]
}

export interface ServiceInput {
  title: string
  slug: string
  summary: string
  description: string
  features: string[]
  stack: string[]
  preview_id: number | null
  status: 'draft' | 'published'
  seo: Seo
}

export interface TeamMemberInput {
  name: string
  role: string
  bio: string
  photo_id: number | null
  socials: SocialRef[]
  is_published: boolean
}

export interface LabInput {
  title: string
  slug: string
  description: string
  stage: LabStage
  url: string
  repo_url: string
  cover_id: number | null
  stack: string[]
  is_published: boolean
}

export interface ArticleInput {
  title: string
  slug: string
  excerpt: string
  cover_id: number | null
  category_id: number | null
  content: Block[]
  author_name: string
  status: ContentStatus
  featured: boolean
  seo: Seo
  published_at: string | null
}

export interface NewsInput {
  title: string
  slug: string
  excerpt: string
  tag: string
  cover_id: number | null
  content: Block[]
  status: ContentStatus
  pinned: boolean
  seo: Seo
  published_at: string | null
}

export interface SettingsInput {
  site: SiteSettings
  metrics: Metric[]
  seo: Seo
  socials: SocialLink[]
}

export interface ListQuery {
  q?: string
  status?: string
  page?: number
  limit?: number
  kind?: string
}

export interface Paged<T> {
  items: T[]
  meta: PageMeta
}

async function get<T>(url: string, params?: object) {
  const { data } = await http.get<ApiResponse<T>>(url, { params })
  return data.data
}

async function paged<T>(url: string, params?: object): Promise<Paged<T>> {
  const { data } = await http.get<ApiResponse<T[]>>(url, { params })
  return { items: data.data, meta: data.meta as PageMeta }
}

async function send<T>(method: 'post' | 'put' | 'patch', url: string, body?: unknown) {
  const { data } = await http.request<ApiResponse<T>>({ method, url, data: body })
  return data?.data
}

async function remove(url: string) {
  await http.delete(url)
}

/** Oddiy CRUD resurs uchun umumiy API (services, team, labs). */
function resource<T, In>(base: string) {
  return {
    list: () => get<T[]>(base),
    get: (id: number) => get<T>(`${base}/${id}`),
    create: (input: In) => send<T>('post', base, input) as Promise<T>,
    update: (id: number, input: In) => send<T>('put', `${base}/${id}`, input) as Promise<T>,
    remove: (id: number) => remove(`${base}/${id}`),
    reorder: (ids: number[]) => send<void>('put', `${base}/reorder`, { ids }),
  }
}

export const adminApi = {
  auth: {
    login: (email: string, password: string) => send<Admin>('post', '/admin/auth/login', { email, password }) as Promise<Admin>,
    logout: () => send<void>('post', '/admin/auth/logout'),
    me: () => get<Admin>('/admin/auth/me'),
    changePassword: (current_password: string, new_password: string) =>
      send<void>('put', '/admin/auth/password', { current_password, new_password }),
  },

  dashboard: () => get<DashboardStats>('/admin/dashboard'),

  projects: {
    list: (q: ListQuery = {}) => paged<Project>('/admin/projects', q),
    get: (id: number) => get<Project>(`/admin/projects/${id}`),
    create: (input: ProjectInput) => send<Project>('post', '/admin/projects', input) as Promise<Project>,
    update: (id: number, input: ProjectInput) => send<Project>('put', `/admin/projects/${id}`, input) as Promise<Project>,
    setStatus: (id: number, status: ContentStatus) => send<void>('patch', `/admin/projects/${id}/status`, { status }),
    setFeatured: (id: number, featured: boolean) => send<void>('patch', `/admin/projects/${id}/featured`, { featured }),
    reorder: (ids: number[]) => send<void>('put', '/admin/projects/reorder', { ids }),
    remove: (id: number) => remove(`/admin/projects/${id}`),
  },

  services: resource<Service, ServiceInput>('/admin/services'),
  team: resource<TeamMember, TeamMemberInput>('/admin/team'),
  labs: resource<Lab, LabInput>('/admin/labs'),

  articles: {
    list: (q: ListQuery = {}) => paged<Article>('/admin/articles', q),
    get: (id: number) => get<Article>(`/admin/articles/${id}`),
    create: (input: ArticleInput) => send<Article>('post', '/admin/articles', input) as Promise<Article>,
    update: (id: number, input: ArticleInput) => send<Article>('put', `/admin/articles/${id}`, input) as Promise<Article>,
    remove: (id: number) => remove(`/admin/articles/${id}`),
  },

  news: {
    list: (q: ListQuery = {}) => paged<News>('/admin/news', q),
    get: (id: number) => get<News>(`/admin/news/${id}`),
    create: (input: NewsInput) => send<News>('post', '/admin/news', input) as Promise<News>,
    update: (id: number, input: NewsInput) => send<News>('put', `/admin/news/${id}`, input) as Promise<News>,
    remove: (id: number) => remove(`/admin/news/${id}`),
  },

  analytics: (days: number) => get<Analytics>('/admin/analytics', { days }),

  categories: {
    list: () => get<ArticleCategory[]>('/admin/article-categories'),
    create: (input: { name: string; slug: string; position: number }) =>
      send<ArticleCategory[]>('post', '/admin/article-categories', input),
    update: (id: number, input: { name: string; slug: string; position: number }) =>
      send<ArticleCategory[]>('put', `/admin/article-categories/${id}`, input),
    remove: (id: number) => remove(`/admin/article-categories/${id}`),
  },

  media: {
    list: (q: ListQuery = {}) => paged<Media>('/admin/media', q),
    async upload(files: File[], onProgress?: (percent: number) => void) {
      const form = new FormData()
      files.forEach((f) => form.append('files', f))
      const { data } = await http.post<{ data: Media[]; failed: { name: string; message: string }[] }>('/admin/media', form, {
        timeout: 0,
        onUploadProgress: (e) => onProgress?.(e.total ? Math.round((e.loaded / e.total) * 100) : 0),
      })
      return data
    },
    update: (id: number, alt: string) => send<Media>('put', `/admin/media/${id}`, { alt }) as Promise<Media>,
    usage: (id: number) => get<{ count: number }>(`/admin/media/${id}/usage`),
    remove: (id: number) => remove(`/admin/media/${id}`),
  },

  requests: {
    list: (q: ListQuery = {}) => paged<ContactRequest>('/admin/requests', q),
    update: (id: number, status: ContactStatus, note: string) =>
      send<ContactRequest>('patch', `/admin/requests/${id}`, { status, note }) as Promise<ContactRequest>,
    remove: (id: number) => remove(`/admin/requests/${id}`),
  },

  settings: {
    get: () => get<Settings>('/admin/settings'),
    update: (input: SettingsInput) => send<Settings>('put', '/admin/settings', input) as Promise<Settings>,
  },

  users: {
    list: () => get<Admin[]>('/admin/users'),
    create: (input: { name: string; email: string; password: string }) => send<Admin>('post', '/admin/users', input) as Promise<Admin>,
    remove: (id: number) => remove(`/admin/users/${id}`),
  },
}
