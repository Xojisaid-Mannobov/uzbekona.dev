<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import FormField from '@/admin/components/FormField.vue'
import MediaField from '@/admin/components/MediaField.vue'
import SeoFields from '@/admin/components/SeoFields.vue'
import EditorSidebar from '@/admin/components/EditorSidebar.vue'
import ErrorSummary from '@/admin/components/ErrorSummary.vue'
import BlockBuilder from '@/admin/modules/builder/BlockBuilder.vue'
import { useEditor } from '@/admin/composables/useEditor'
import { adminApi, type NewsInput } from '@/admin/services/adminApi'
import { toneStyle } from '@/modules/news/tone'
import { formatNumber } from '@/utils/format'
import type { Block, MediaRef, News } from '@/types/api'

type NewsForm = Omit<NewsInput, 'cover_id' | 'published_at'> & { cover: MediaRef | null; published_at: string }

// Tez tanlanadigan teglar — har biri saytda o'z rangida ko'rinadi
const TAGS = ['Studiya', 'Loyiha', 'Labs', 'Jamoa', 'Tadbir', 'Hamkorlik']

const route = useRoute()
const id = computed(() => (route.params.id ? Number(route.params.id) : null))

// <input type="datetime-local"> formati ↔ ISO
function toLocal(iso: string | null): string {
  if (!iso) return ''
  const d = new Date(iso)
  const off = d.getTimezoneOffset() * 60000
  return new Date(d.getTime() - off).toISOString().slice(0, 16)
}

const { form, entity, loading, loadError, saving, errors, dirty, isNew, save, load } = useEditor<News, NewsForm, NewsInput>({
  id,
  empty: () => ({
    title: '',
    slug: '',
    excerpt: '',
    tag: 'Studiya',
    cover: null,
    content: [{ type: 'text', data: { text: '' } }],
    status: 'draft',
    pinned: false,
    seo: { title: '', description: '' },
    published_at: '',
  }),
  fromEntity: (n) => ({
    title: n.title,
    slug: n.slug,
    excerpt: n.excerpt,
    tag: n.tag,
    cover: n.cover,
    content: JSON.parse(JSON.stringify(n.content ?? [])),
    status: n.status,
    pinned: n.pinned,
    seo: { ...n.seo },
    published_at: toLocal(n.published_at),
  }),
  toInput: ({ cover, published_at, content, ...rest }) => ({
    ...rest,
    cover_id: cover?.id ?? null,
    published_at: published_at ? new Date(published_at).toISOString() : null,
    content: content.map((b) => ({ type: b.type, data: b.data }) as Block),
  }),
  load: adminApi.news.get,
  create: adminApi.news.create,
  update: adminApi.news.update,
  editRoute: (n) => `/admin/news/${n.id}`,
})

const publicUrl = computed(() => (entity.value?.status === 'published' ? `/news/${entity.value.slug}` : null))
</script>

<template>
  <div>
    <PageHeader :title="isNew ? 'Yangi yangilik' : form.title || 'Yangilik'" back="/admin/news" />

    <AdminEmpty v-if="loadError" error title="Yuklab bo‘lmadi" :text="loadError.message" @retry="load" />
    <div v-else-if="loading" class="skeleton" style="height: 520px; border-radius: 20px" />

    <form v-else class="a-edit" novalidate @submit.prevent="save">
      <div class="a-stack">
        <ErrorSummary :errors="errors" :labels="{ content: 'Blok' }" />
        <section class="a-card a-stack">
          <FormField label="Sarlavha" for="title" :error="errors.title">
            <input id="title" v-model="form.title" class="a-input title-input" placeholder="Uzbekona.dev yangi loyihani ishga tushirdi" />
          </FormField>
          <FormField label="Slug" for="slug" :error="errors.slug" optional>
            <input id="slug" v-model="form.slug" class="a-input" placeholder="avtomatik" />
          </FormField>
          <FormField
            label="Qisqa mazmun"
            for="excerpt"
            :error="errors.excerpt"
            hint="Kartada, yangiliklar lentasida va sarlavha ostida ko‘rinadi"
          >
            <textarea id="excerpt" v-model="form.excerpt" class="a-textarea" />
          </FormField>
        </section>

        <section>
          <h2 class="a-card__title">Yangilik matni</h2>
          <BlockBuilder v-model="form.content" />
        </section>

        <section class="a-card">
          <h2 class="a-card__title">SEO</h2>
          <SeoFields
            v-model="form.seo"
            :fallback-title="form.title"
            :fallback-description="form.excerpt"
            :path="`/news/${form.slug || '…'}`"
          />
        </section>
      </div>

      <aside class="a-edit__side">
        <EditorSidebar :saving="saving" :dirty="dirty" :is-new="isNew" :public-url="publicUrl" @save="save">
          <FormField label="Holat" for="status">
            <select id="status" v-model="form.status" class="a-select">
              <option value="draft">Qoralama</option>
              <option value="published">E’lon qilingan</option>
              <option value="archived">Arxiv</option>
            </select>
          </FormField>
          <FormField label="E’lon qilish sanasi" for="pub" hint="Kelajak sana — rejalashtirilgan e’lon" optional>
            <input id="pub" v-model="form.published_at" class="a-input" type="datetime-local" />
          </FormField>
          <label class="a-switch">
            <input v-model="form.pinned" type="checkbox" />
            <span class="a-switch__track" />
            Mahkamlash (asosiy yangilik)
          </label>
          <p class="a-hint">Mahkamlangan yangilik ro‘yxat boshida katta kartada chiqadi. Bir vaqtda faqat bittasi.</p>
        </EditorSidebar>

        <section class="a-card a-stack">
          <FormField label="Teg" for="tag" :error="errors.tag" hint="Rang teg bo‘yicha avtomatik tanlanadi">
            <input id="tag" v-model="form.tag" class="a-input" maxlength="40" list="news-tags" />
            <datalist id="news-tags">
              <option v-for="t in TAGS" :key="t" :value="t" />
            </datalist>
          </FormField>
          <div class="tags">
            <button
              v-for="t in TAGS"
              :key="t"
              type="button"
              class="tags__chip"
              :class="{ 'is-active': form.tag === t }"
              :style="toneStyle(t)"
              @click="form.tag = t"
            >
              {{ t }}
            </button>
          </div>
          <p v-if="entity" class="a-hint">Ko‘rishlar: {{ formatNumber(entity.views) }} (noyob tashrif buyuruvchilar)</p>
        </section>

        <section class="a-card">
          <h2 class="a-card__title">Muqova</h2>
          <MediaField v-model="form.cover" />
          <p class="a-hint" style="margin-top: 10px">Bo‘lmasa, teg rangidagi naqshli muqova avtomatik chiziladi.</p>
        </section>
      </aside>
    </form>
  </div>
</template>

<style scoped>
.title-input {
  height: 56px;
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.tags__chip {
  padding: 5px 11px;
  border-radius: 999px;
  border: 1px solid color-mix(in srgb, var(--tone) 30%, transparent);
  color: var(--tone);
  font-size: 12px;
  font-weight: 700;
  transition: background-color 0.15s ease;
}

.tags__chip:hover,
.tags__chip.is-active {
  background: color-mix(in srgb, var(--tone) 14%, transparent);
}
</style>
