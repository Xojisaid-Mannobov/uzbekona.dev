import { http } from './http'
import type {
  ApiResponse,
  Article,
  ArticleCategory,
  ArticleDetail,
  ContactPayload,
  Lab,
  PageMeta,
  Project,
  ProjectDetail,
  Service,
  Settings,
  TeamMember,
} from '@/types/api'

// Public sayt uchun API funksiyalari

async function get<T>(url: string, params?: Record<string, unknown>): Promise<T> {
  const { data } = await http.get<ApiResponse<T>>(url, { params })
  return data.data
}

export const publicApi = {
  settings: () => get<Settings>('/settings'),

  projects: (params: { featured?: boolean; limit?: number } = {}) => get<Project[]>('/projects', params),
  project: (slug: string) => get<ProjectDetail>(`/projects/${encodeURIComponent(slug)}`),

  services: () => get<Service[]>('/services'),
  service: (slug: string) => get<Service>(`/services/${encodeURIComponent(slug)}`),

  team: () => get<TeamMember[]>('/team'),
  labs: () => get<Lab[]>('/labs'),

  categories: () => get<ArticleCategory[]>('/article-categories'),
  async articles(params: { category?: string; page?: number; limit?: number } = {}) {
    const { data } = await http.get<ApiResponse<Article[]>>('/articles', { params })
    return { items: data.data, meta: data.meta as PageMeta }
  },
  article: (slug: string) => get<ArticleDetail>(`/articles/${encodeURIComponent(slug)}`),

  async contact(payload: ContactPayload) {
    const { data } = await http.post<ApiResponse<{ message: string }>>('/contact', payload)
    return data.data
  },
}
