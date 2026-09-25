<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import FormField from '@/admin/components/FormField.vue'
import TagInput from '@/admin/components/TagInput.vue'
import MediaField from '@/admin/components/MediaField.vue'
import MediaPicker from '@/admin/components/MediaPicker.vue'
import MetricsEditor from '@/admin/components/MetricsEditor.vue'
import SeoFields from '@/admin/components/SeoFields.vue'
import EditorSidebar from '@/admin/components/EditorSidebar.vue'
import ErrorSummary from '@/admin/components/ErrorSummary.vue'
import BlockBuilder from '@/admin/modules/builder/BlockBuilder.vue'
import { useEditor } from '@/admin/composables/useEditor'
import { adminApi, type ProjectInput } from '@/admin/services/adminApi'
import { useDragSort } from '@/admin/composables/useDragSort'
import type { Block, ContentStatus, Media, MediaRef, Metric, Project, Seo } from '@/types/api'

interface ProjectForm {
  title: string
  slug: string
  tagline: string
  short_description: string
  full_description: string
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
  seo: Seo
  gallery: { media: MediaRef; caption: string }[]
  blocks: Block[]
}

const route = useRoute()
const id = computed(() => (route.params.id ? Number(route.params.id) : null))

const { form, entity, loading, loadError, saving, errors, dirty, isNew, save, load } = useEditor<Project, ProjectForm, ProjectInput>({
  id,
  empty: () => ({
    title: '',
    slug: '',
    tagline: '',
    short_description: '',
    full_description: '',
    cover: null,
    year: new Date().getFullYear(),
    client: '',
    industry: '',
    platforms: ['Web'],
    services: [],
    stack: [],
    metrics: [],
    live_url: '',
    accent: '#2350F5',
    status: 'draft',
    featured: false,
    seo: { title: '', description: '' },
    gallery: [],
    blocks: [],
  }),
  fromEntity: (p) => ({
    title: p.title,
    slug: p.slug,
    tagline: p.tagline,
    short_description: p.short_description,
    full_description: p.full_description,
    cover: p.cover,
    year: p.year,
    client: p.client,
    industry: p.industry,
    platforms: [...p.platforms],
    services: [...p.services],
    stack: [...p.stack],
    metrics: p.metrics.map((m) => ({ ...m })),
    live_url: p.live_url,
    accent: p.accent || '#2350F5',
    status: p.status,
    featured: p.featured,
    seo: { ...p.seo },
    gallery: (p.gallery ?? []).filter((g) => g.media).map((g) => ({ media: g.media!, caption: g.caption })),
    blocks: JSON.parse(JSON.stringify(p.blocks ?? [])),
  }),
  toInput: (f) => ({
    ...f,
    cover_id: f.cover?.id ?? null,
    year: f.year ? Number(f.year) : null,
    metrics: f.metrics.filter((m) => m.value || m.label),
    gallery: f.gallery.map((g) => ({ media_id: g.media.id, caption: g.caption })),
    // id'lar serverda qayta yaratiladi
    blocks: f.blocks.map((b) => ({ type: b.type, data: b.data }) as Block),
  }),
  load: adminApi.projects.get,
  create: adminApi.projects.create,
  update: adminApi.projects.update,
  editRoute: (p) => `/admin/projects/${p.id}`,
})

const tab = ref<'main' | 'builder' | 'gallery' | 'seo'>('main')
const tabs = computed(
  () =>
    [
      { value: 'main', label: 'Asosiy' },
      { value: 'builder', label: `Case study · ${form.value.blocks.length}` },
      { value: 'gallery', label: `Galereya · ${form.value.gallery.length}` },
      { value: 'seo', label: 'SEO' },
    ] as const,
)

const publicUrl = computed(() => (entity.value?.status === 'published' ? `/projects/${entity.value.slug}` : null))

// Galereya
const galleryOpen = ref(false)
const galleryList = computed({
  get: () => form.value.gallery,
  set: (v) => (form.value.gallery = v),
})
const { handlers: gHandlers, dragIndex: gDrag } = useDragSort(galleryList)

function addToGallery(items: Media[]) {
  form.value.gallery = [
    ...form.value.gallery,
    ...items
      .filter((m) => !form.value.gallery.some((g) => g.media.id === m.id))
      .map((m) => ({
        media: { id: m.id, url: m.url, kind: m.kind, mime: m.mime, width: m.width, height: m.height, alt: m.alt, variants: m.variants },
        caption: '',
      })),
  ]
}

const errorLabels = {
  blocks: 'Blok',
  gallery: 'Galereya',
  metrics: 'Metrika',
  title: 'Nomi',
  slug: 'Slug',
  seo: 'SEO',
  live_url: 'Sayt havolasi',
  accent: 'Rang',
}
</script>

<template>
  <div>
    <PageHeader
      :title="isNew ? 'Yangi loyiha' : form.title || 'Loyiha'"
      :subtitle="isNew ? 'Case-study yaratish' : `/projects/${entity?.slug ?? ''}`"
      back="/admin/projects"
    />

    <AdminEmpty
      v-if="loadError"
      error
      :title="loadError.status === 404 ? 'Loyiha topilmadi' : 'Yuklab bo‘lmadi'"
      :text="loadError.message"
      @retry="load"
    />

    <div v-else-if="loading" class="a-edit">
      <div class="skeleton" style="height: 520px; border-radius: 20px" />
      <div class="skeleton" style="height: 320px; border-radius: 20px" />
    </div>

    <form v-else class="a-edit" novalidate @submit.prevent="save">
      <div>
        <ErrorSummary :errors="errors" :labels="errorLabels" />

        <div class="a-tabs tabs" role="tablist">
          <button
            v-for="t in tabs"
            :key="t.value"
            type="button"
            role="tab"
            class="a-tab"
            :class="{ 'is-active': tab === t.value }"
            :aria-selected="tab === t.value"
            @click="tab = t.value"
          >
            {{ t.label }}
          </button>
        </div>

        <!-- Asosiy -->
        <section v-show="tab === 'main'" class="a-card a-stack">
          <FormField label="Nomi" for="title" :error="errors.title">
            <input id="title" v-model="form.title" class="a-input title-input" placeholder="KUAF" />
          </FormField>
          <FormField label="Slug" for="slug" :error="errors.slug" hint="Bo‘sh qoldirilsa nomidan avtomatik yasaladi" optional>
            <div class="slug">
              <span>/projects/</span>
              <input id="slug" v-model="form.slug" class="a-input" placeholder="kuaf" />
            </div>
          </FormField>
          <FormField label="Tagline (hero sarlavha)" for="tagline" :error="errors.tagline" hint="Case-study sahifasidagi katta jumla">
            <input id="tagline" v-model="form.tagline" class="a-input" placeholder="Universitet uchun yagona raqamli ekotizim." />
          </FormField>
          <FormField label="Qisqa tavsif" for="short" :error="errors.short_description" hint="Kartalarda nom ostida ko‘rinadi">
            <input id="short" v-model="form.short_description" class="a-input" placeholder="University Digital Ecosystem" />
          </FormField>
          <FormField
            label="To‘liq tavsif"
            for="full"
            :error="errors.full_description"
            hint="“Loyiha haqida” bo‘limi. Bo‘sh qator — yangi paragraf"
          >
            <textarea id="full" v-model="form.full_description" class="a-textarea a-textarea--lg" />
          </FormField>
        </section>

        <!-- Content builder -->
        <section v-show="tab === 'builder'">
          <BlockBuilder v-model="form.blocks" />
        </section>

        <!-- Galereya -->
        <section v-show="tab === 'gallery'" class="a-card">
          <div class="gallery">
            <div
              v-for="(g, i) in form.gallery"
              :key="g.media.id"
              class="gallery__item"
              :class="{ 'is-dragging': gDrag === i }"
              v-bind="gHandlers(i)"
            >
              <img :src="g.media.variants[0]?.url ?? g.media.url" :alt="g.media.alt" draggable="false" />
              <input v-model="g.caption" class="a-input" placeholder="Izoh" />
              <button
                type="button"
                class="a-icon-btn a-icon-btn--danger gallery__rm"
                aria-label="Olib tashlash"
                @click="form.gallery.splice(i, 1)"
              >
                <AppIcon name="trash" :size="16" />
              </button>
            </div>
            <button type="button" class="gallery__add" @click="galleryOpen = true"><AppIcon name="plus" :size="24" /> Rasm qo‘shish</button>
          </div>
          <p class="a-hint" style="margin-top: 14px">Rasmlarni sudrab tartiblang. Har 3-rasm to‘liq kenglikda ko‘rsatiladi.</p>
          <MediaPicker v-model:open="galleryOpen" multiple @select="addToGallery" />
        </section>

        <!-- SEO -->
        <section v-show="tab === 'seo'" class="a-card">
          <SeoFields
            v-model="form.seo"
            :fallback-title="form.title"
            :fallback-description="form.tagline"
            :path="`/projects/${form.slug || '…'}`"
          />
        </section>
      </div>

      <aside class="a-edit__side">
        <EditorSidebar :saving="saving" :dirty="dirty" :is-new="isNew" :public-url="publicUrl" @save="save">
          <FormField label="Holat" for="status" :error="errors.status">
            <select id="status" v-model="form.status" class="a-select">
              <option value="draft">Qoralama</option>
              <option value="published">E’lon qilingan</option>
              <option value="archived">Arxiv</option>
            </select>
          </FormField>
          <label class="a-switch">
            <input v-model="form.featured" type="checkbox" />
            <span class="a-switch__track" />
            Bosh sahifada (featured)
          </label>
        </EditorSidebar>

        <section class="a-card">
          <h2 class="a-card__title">Cover</h2>
          <MediaField v-model="form.cover" label="Cover rasm tanlash" />
          <p class="a-hint" style="margin-top: 10px">Katta screenshot (kamida 2400px kenglik tavsiya etiladi).</p>
        </section>

        <section class="a-card a-stack">
          <h2 class="a-card__title" style="margin: 0">Ma’lumotlar</h2>
          <div class="a-form-grid">
            <FormField label="Yil" for="year" :error="errors.year">
              <input id="year" v-model.number="form.year" class="a-input" type="number" min="1990" max="2100" />
            </FormField>
            <FormField label="Accent rang" for="accent" :error="errors.accent">
              <div class="color">
                <input v-model="form.accent" type="color" aria-label="Rang tanlash" />
                <input id="accent" v-model="form.accent" class="a-input" maxlength="7" />
              </div>
            </FormField>
          </div>
          <FormField label="Mijoz (Client)" for="client" :error="errors.client">
            <input id="client" v-model="form.client" class="a-input" />
          </FormField>
          <FormField label="Soha (Industry)" for="industry" :error="errors.industry">
            <input id="industry" v-model="form.industry" class="a-input" placeholder="Education" />
          </FormField>
          <FormField label="Platformalar" :error="errors.platforms">
            <TagInput v-model="form.platforms" :suggestions="['Web', 'Mobile', 'Telegram', 'Desktop']" />
          </FormField>
          <FormField label="Xizmatlar" :error="errors.services">
            <TagInput v-model="form.services" placeholder="Web platforma, Product Design…" />
          </FormField>
          <FormField label="Texnologiyalar (Stack)" :error="errors.stack">
            <TagInput v-model="form.stack" :suggestions="['Vue.js', 'TypeScript', 'Go', 'PostgreSQL', 'Redis', 'Docker']" />
          </FormField>
          <FormField label="Sayt havolasi" for="live" :error="errors.live_url" optional>
            <input id="live" v-model="form.live_url" class="a-input" placeholder="https://kuaf.uz" />
          </FormField>
        </section>

        <section class="a-card">
          <h2 class="a-card__title">Natija metrikalari</h2>
          <MetricsEditor v-model="form.metrics" :errors="errors" />
          <p class="a-hint" style="margin-top: 10px">Faqat real, tasdiqlangan raqamlar.</p>
        </section>
      </aside>
    </form>
  </div>
</template>

<style scoped>
.tabs {
  margin-bottom: 16px;
}

.title-input {
  height: 56px;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.slug {
  display: flex;
  align-items: center;
  border: 1px solid var(--line-strong);
  border-radius: var(--a-radius);
  background: var(--surface);
  overflow: hidden;
}

.slug span {
  padding: 0 4px 0 16px;
  color: var(--ink-3);
  font-size: 14px;
}

.slug .a-input {
  border: 0;
  padding-left: 0;
  box-shadow: none;
}

.slug:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 0 4px var(--accent-soft);
}

.color {
  display: flex;
  gap: 8px;
}

.color input[type='color'] {
  width: 48px;
  height: 48px;
  padding: 4px;
  border-radius: var(--a-radius);
  border: 1px solid var(--line-strong);
  background: var(--surface);
  cursor: pointer;
  flex-shrink: 0;
}

.gallery {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 14px;
}

.gallery__item {
  position: relative;
  display: grid;
  gap: 8px;
  cursor: grab;
}

.gallery__item.is-dragging {
  opacity: 0.4;
}

.gallery__item img {
  aspect-ratio: 4 / 3;
  width: 100%;
  object-fit: cover;
  border-radius: 14px;
  background: var(--surface-2);
}

.gallery__item .a-input {
  height: 40px;
  font-size: 13px;
}

.gallery__rm {
  position: absolute;
  top: 8px;
  right: 8px;
  background: var(--surface);
}

.gallery__add {
  aspect-ratio: 4 / 3;
  display: grid;
  place-content: center;
  justify-items: center;
  gap: 8px;
  border-radius: 14px;
  border: 1.5px dashed var(--line-strong);
  color: var(--ink-2);
  font-weight: 600;
}

.gallery__add:hover {
  border-color: var(--accent);
  color: var(--accent);
}
</style>
