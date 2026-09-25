// Backend (Go) javoblari bilan bir xil turlar.

export type ContentStatus = 'draft' | 'published' | 'archived'

export interface MediaVariant {
  width: number
  url: string
}

export interface MediaRef {
  id: number
  url: string
  kind: 'image' | 'vector' | 'video'
  mime: string
  width: number | null
  height: number | null
  alt: string
  variants: MediaVariant[]
}

export interface Media extends MediaRef {
  original_name: string
  size: number
  created_at: string
}

export interface Metric {
  value: string
  label: string
}

export interface Seo {
  title: string
  description: string
}

export interface SocialRef {
  platform: string
  url: string
}

// ─── Content builder bloklari ───────────────────────────

export type BlockType =
  | 'heading'
  | 'text'
  | 'large_text'
  | 'image'
  | 'full_image'
  | 'gallery'
  | 'video'
  | 'stats'
  | 'quote'
  | 'two_columns'
  | 'three_columns'
  | 'technology'
  | 'process'
  | 'before_after'

export interface Column {
  title: string
  text: string
}

export interface BlockDataMap {
  heading: { text: string; label?: string }
  text: { text: string }
  large_text: { text: string }
  image: { media: MediaRef | null; caption?: string }
  full_image: { media: MediaRef | null; caption?: string }
  gallery: { items: { media: MediaRef; caption?: string }[] }
  video: { media: MediaRef | null; poster: MediaRef | null; url?: string; caption?: string }
  stats: { items: Metric[] }
  quote: { text: string; author?: string; role?: string }
  two_columns: { columns: Column[] }
  three_columns: { columns: Column[] }
  technology: { title?: string; items: string[] }
  process: { steps: Column[] }
  before_after: { before: MediaRef | null; after: MediaRef | null; caption?: string }
}

export type Block<T extends BlockType = BlockType> = {
  [K in T]: { id?: number; type: K; data: BlockDataMap[K] }
}[T]

// ─── Kontent ────────────────────────────────────────────

export interface GalleryItem {
  media_id: number
  caption: string
  media?: MediaRef | null
}

export interface Project {
  id: number
  slug: string
  title: string
  tagline: string
  short_description: string
  full_description: string
  cover_id: number | null
  cover: MediaRef | null
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
  position: number
  seo: Seo
  published_at: string | null
  created_at: string
  updated_at: string
  blocks?: Block[]
  gallery?: GalleryItem[]
}

export interface ProjectDetail extends Project {
  next: Project | null
}

export interface Service {
  id: number
  slug: string
  title: string
  summary: string
  description: string
  features: string[]
  stack: string[]
  preview_id: number | null
  preview: MediaRef | null
  status: 'draft' | 'published'
  position: number
  seo: Seo
}

export interface TeamMember {
  id: number
  name: string
  role: string
  bio: string
  photo_id: number | null
  photo: MediaRef | null
  socials: SocialRef[]
  is_published: boolean
  position: number
}

export type LabStage = 'open_source' | 'experimental' | 'in_development'

export interface Lab {
  id: number
  slug: string
  title: string
  description: string
  stage: LabStage
  url: string
  repo_url: string
  cover_id: number | null
  cover: MediaRef | null
  stack: string[]
  is_published: boolean
  position: number
}

export interface ArticleCategory {
  id: number
  name: string
  slug: string
  position: number
  count: number
}

export interface Article {
  id: number
  slug: string
  title: string
  excerpt: string
  cover_id: number | null
  cover: MediaRef | null
  category_id: number | null
  category: ArticleCategory | null
  content?: Block[]
  reading_time: number
  author_name: string
  status: ContentStatus
  featured: boolean
  seo: Seo
  published_at: string | null
  created_at: string
  updated_at: string
}

export interface ArticleDetail extends Article {
  related: Article[]
}

// ─── Yangiliklar ────────────────────────────────────────

export interface News {
  id: number
  slug: string
  title: string
  excerpt: string
  tag: string
  cover_id: number | null
  cover: MediaRef | null
  content?: Block[]
  status: ContentStatus
  pinned: boolean
  views: number
  seo: Seo
  published_at: string | null
  created_at: string
  updated_at: string
}

export interface NewsDetail extends News {
  related: News[]
}

// ─── Tizim ──────────────────────────────────────────────

export interface SiteSettings {
  name: string
  tagline: string
  email: string
  phone: string
  telegram: string
  address: string
  available: boolean
}

export interface SocialLink {
  id?: number
  platform: string
  label: string
  url: string
}

export interface Settings {
  site: SiteSettings
  metrics: Metric[]
  seo: Seo
  socials: SocialLink[]
}

export type ContactStatus = 'new' | 'in_progress' | 'done' | 'spam'

export interface ContactRequest {
  id: number
  name: string
  contact: string
  email: string
  project_type: string
  budget: string
  message: string
  status: ContactStatus
  note: string
  ip: string
  user_agent: string
  created_at: string
  updated_at: string
}

export interface ContactPayload {
  name: string
  contact: string
  email: string
  project_type: string
  budget: string
  message: string
  website: string
}

export interface Admin {
  id: number
  name: string
  email: string
  last_login_at: string | null
  created_at: string
}

export interface DashboardStats {
  active_projects: number
  published_projects: number
  articles: number
  news: number
  views_today: number
  visitors_today: number
  incoming_requests: number
  total_requests: number
  media: number
  admins: number
  recent_requests: ContactRequest[]
  recent_projects: Project[]
  requests_by_day: { day: string; count: number }[]
}

// ─── Statistika ─────────────────────────────────────────

export interface AnalyticsRow {
  key: string
  views: number
  visitors: number
}

export interface Analytics {
  days: number
  summary: {
    views: number
    visits: number
    visitors: number
    prev_views: number
    prev_visits: number
    prev_visitors: number
    today_views: number
    today_visitors: number
    online_now: number
    all_time_views: number
    all_time_visitors: number
  }
  by_day: { day: string; views: number; visitors: number }[]
  top_pages: AnalyticsRow[]
  referrers: AnalyticsRow[]
  devices: AnalyticsRow[]
  top_news: { id: number; slug: string; title: string; views: number }[]
}

// ─── API konvertlari ────────────────────────────────────

export interface PageMeta {
  page: number
  limit: number
  total: number
}

export interface ApiResponse<T> {
  data: T
  meta?: PageMeta
}

export interface ApiErrorBody {
  code: string
  message: string
  fields?: Record<string, string>
}
