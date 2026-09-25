<script setup lang="ts">
import { computed, ref } from 'vue'
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
import { useAuthStore } from '@/admin/stores/auth'
import { adminApi, type ArticleInput } from '@/admin/services/adminApi'
import type { Article, ArticleCategory, Block, MediaRef } from '@/types/api'

type ArticleForm = Omit<ArticleInput, 'cover_id' | 'published_at'> & { cover: MediaRef | null; published_at: string }

const route = useRoute()
const auth = useAuthStore()
const id = computed(() => (route.params.id ? Number(route.params.id) : null))
const categories = ref<ArticleCategory[]>([])
adminApi.categories
  .list()
  .then((c) => (categories.value = c))
  .catch(() => {})

// <input type="datetime-local"> formati ↔ ISO
function toLocal(iso: string | null): string {
  if (!iso) return ''
  const d = new Date(iso)
  const off = d.getTimezoneOffset() * 60000
  return new Date(d.getTime() - off).toISOString().slice(0, 16)
}

const { form, entity, loading, loadError, saving, errors, dirty, isNew, save, load } = useEditor<Article, ArticleForm, ArticleInput>({
  id,
  empty: () => ({
    title: '',
    slug: '',
    excerpt: '',
    cover: null,
    category_id: null,
    content: [{ type: 'text', data: { text: '' } }],
    author_name: auth.admin?.name ?? '',
    status: 'draft',
    featured: false,
    seo: { title: '', description: '' },
    published_at: '',
  }),
  fromEntity: (a) => ({
    title: a.title,
    slug: a.slug,
    excerpt: a.excerpt,
    cover: a.cover,
    category_id: a.category_id,
    content: JSON.parse(JSON.stringify(a.content ?? [])),
    author_name: a.author_name,
    status: a.status,
    featured: a.featured,
    seo: { ...a.seo },
    published_at: toLocal(a.published_at),
  }),
  toInput: ({ cover, published_at, content, ...rest }) => ({
    ...rest,
    cover_id: cover?.id ?? null,
    published_at: published_at ? new Date(published_at).toISOString() : null,
    content: content.map((b) => ({ type: b.type, data: b.data }) as Block),
  }),
  load: adminApi.articles.get,
  create: adminApi.articles.create,
  update: adminApi.articles.update,
  editRoute: (a) => `/admin/articles/${a.id}`,
})

const publicUrl = computed(() => (entity.value?.status === 'published' ? `/journal/${entity.value.slug}` : null))
</script>

<template>
  <div>
    <PageHeader :title="isNew ? 'Yangi maqola' : form.title || 'Maqola'" back="/admin/articles" />

    <AdminEmpty v-if="loadError" error title="Yuklab bo‘lmadi" :text="loadError.message" @retry="load" />
    <div v-else-if="loading" class="skeleton" style="height: 520px; border-radius: 20px" />

    <form v-else class="a-edit" novalidate @submit.prevent="save">
      <div class="a-stack">
        <ErrorSummary :errors="errors" :labels="{ content: 'Blok' }" />
        <section class="a-card a-stack">
          <FormField label="Sarlavha" for="title" :error="errors.title">
            <input id="title" v-model="form.title" class="a-input title-input" placeholder="Nega backend uchun Go tanladik" />
          </FormField>
          <FormField label="Slug" for="slug" :error="errors.slug" optional>
            <input id="slug" v-model="form.slug" class="a-input" placeholder="avtomatik" />
          </FormField>
          <FormField label="Qisqa mazmun (excerpt)" for="excerpt" :error="errors.excerpt" hint="Kartada va sarlavha ostida ko‘rinadi">
            <textarea id="excerpt" v-model="form.excerpt" class="a-textarea" />
          </FormField>
        </section>

        <section>
          <h2 class="a-card__title">Maqola matni</h2>
          <BlockBuilder v-model="form.content" />
        </section>

        <section class="a-card">
          <h2 class="a-card__title">SEO</h2>
          <SeoFields
            v-model="form.seo"
            :fallback-title="form.title"
            :fallback-description="form.excerpt"
            :path="`/journal/${form.slug || '…'}`"
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
            <input v-model="form.featured" type="checkbox" />
            <span class="a-switch__track" />
            Tanlangan maqola
          </label>
        </EditorSidebar>

        <section class="a-card a-stack">
          <FormField label="Kategoriya" for="cat">
            <select id="cat" v-model="form.category_id" class="a-select">
              <option :value="null">Kategoriyasiz</option>
              <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
          </FormField>
          <FormField label="Muallif" for="author" :error="errors.author_name">
            <input id="author" v-model="form.author_name" class="a-input" />
          </FormField>
        </section>

        <section class="a-card">
          <h2 class="a-card__title">Cover</h2>
          <MediaField v-model="form.cover" />
          <p class="a-hint" style="margin-top: 10px">Bo‘lmasa, kartada kategoriya asosidagi tipografik muqova ko‘rsatiladi.</p>
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
</style>
