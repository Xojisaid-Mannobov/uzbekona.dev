<script setup lang="ts">
import { ref, watch } from 'vue'
import { RouterLink } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import StatusBadge from '@/admin/components/StatusBadge.vue'
import PaginationBar from '@/admin/components/PaginationBar.vue'
import { useAsync } from '@/composables/useAsync'
import { useUiStore } from '@/admin/stores/ui'
import { adminApi } from '@/admin/services/adminApi'
import { toApiError } from '@/services/http'
import { formatDate } from '@/utils/format'
import type { Article, ArticleCategory } from '@/types/api'

const ui = useUiStore()
const status = ref('')
const q = ref('')
const page = ref(1)

const { data, loading, error, reload } = useAsync(
  () => adminApi.articles.list({ status: status.value, q: q.value, page: page.value, limit: 20 }),
  {
    watch: [status, page],
  },
)

let timer: ReturnType<typeof setTimeout>
watch(q, () => {
  clearTimeout(timer)
  timer = setTimeout(() => {
    page.value = 1
    reload()
  }, 300)
})

async function remove(a: Article) {
  const ok = await ui.confirm({ title: `“${a.title}” o‘chirilsinmi?`, text: 'Maqola butunlay o‘chiriladi.' })
  if (!ok) return
  try {
    await adminApi.articles.remove(a.id)
    ui.success('O‘chirildi')
    reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}

// ─── Kategoriyalar ──────────────────────────────────
const categories = ref<ArticleCategory[]>([])
const catForm = ref({ id: 0, name: '', slug: '', position: 0 })
const catErrors = ref<Record<string, string>>({})

async function loadCategories() {
  try {
    categories.value = await adminApi.categories.list()
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}
loadCategories()

function editCategory(c: ArticleCategory) {
  catForm.value = { id: c.id, name: c.name, slug: c.slug, position: c.position }
}

async function saveCategory() {
  catErrors.value = {}
  const { id, ...input } = catForm.value
  try {
    categories.value = (id ? await adminApi.categories.update(id, input) : await adminApi.categories.create(input)) ?? []
    catForm.value = { id: 0, name: '', slug: '', position: categories.value.length + 1 }
    ui.success('Kategoriya saqlandi')
  } catch (e) {
    catErrors.value = toApiError(e).fields
    ui.error(toApiError(e).message)
  }
}

async function removeCategory(c: ArticleCategory) {
  const ok = await ui.confirm({
    title: `“${c.name}” kategoriyasi o‘chirilsinmi?`,
    text: 'Maqolalar o‘chirilmaydi — ular kategoriyasiz qoladi.',
  })
  if (!ok) return
  try {
    await adminApi.categories.remove(c.id)
    await loadCategories()
    ui.success('O‘chirildi')
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}
</script>

<template>
  <div>
    <PageHeader title="Maqolalar" subtitle="Journal: engineering, product, dizayn, AI va infratuzilma.">
      <RouterLink to="/admin/articles/new" class="a-btn a-btn--primary"><AppIcon name="plus" :size="18" /> Yangi maqola</RouterLink>
    </PageHeader>

    <div class="layout">
      <div>
        <div class="a-toolbar">
          <div class="a-tabs">
            <button
              v-for="t in [
                { v: '', l: 'Barchasi' },
                { v: 'published', l: 'E’lon qilingan' },
                { v: 'draft', l: 'Qoralama' },
                { v: 'archived', l: 'Arxiv' },
              ]"
              :key="t.v"
              type="button"
              class="a-tab"
              :class="{ 'is-active': status === t.v }"
              @click="((status = t.v), (page = 1))"
            >
              {{ t.l }}
            </button>
          </div>
          <div class="a-search">
            <AppIcon name="search" :size="18" />
            <input v-model="q" class="a-input" type="search" placeholder="Sarlavha bo‘yicha qidirish" aria-label="Qidirish" />
          </div>
        </div>

        <AdminEmpty v-if="error" error title="Yuklab bo‘lmadi" :text="error.message" @retry="reload" />
        <div v-else-if="loading && !data" class="a-stack">
          <div v-for="n in 4" :key="n" class="skeleton" style="height: 72px; border-radius: 14px" />
        </div>
        <AdminEmpty
          v-else-if="!data?.items.length"
          icon="file-text"
          title="Maqola topilmadi"
          text="Birinchi maqolani yozing — tajriba va kuzatuvlaringizni ulashing."
        >
          <RouterLink to="/admin/articles/new" class="a-btn a-btn--primary a-btn--sm">Maqola yozish</RouterLink>
        </AdminEmpty>

        <template v-else>
          <table class="a-table">
            <thead>
              <tr>
                <th>Sarlavha</th>
                <th>Kategoriya</th>
                <th>Holat</th>
                <th>Sana</th>
                <th style="text-align: right">Amallar</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="a in data.items" :key="a.id">
                <td>
                  <RouterLink :to="`/admin/articles/${a.id}`">
                    <span class="a-cell-title">{{ a.title }}</span>
                    <span class="a-cell-sub" style="display: block">{{ a.reading_time }} daqiqa · {{ a.author_name || 'Muallifsiz' }}</span>
                  </RouterLink>
                </td>
                <td>{{ a.category?.name ?? '—' }}</td>
                <td><StatusBadge :status="a.status" /></td>
                <td class="a-cell-sub">{{ formatDate(a.published_at ?? a.created_at) }}</td>
                <td>
                  <div class="a-actions">
                    <a
                      v-if="a.status === 'published'"
                      :href="`/journal/${a.slug}`"
                      target="_blank"
                      class="a-icon-btn"
                      title="Saytda ko‘rish"
                      ><AppIcon name="external" :size="18"
                    /></a>
                    <RouterLink :to="`/admin/articles/${a.id}`" class="a-icon-btn" title="Tahrirlash"
                      ><AppIcon name="edit" :size="18"
                    /></RouterLink>
                    <button type="button" class="a-icon-btn a-icon-btn--danger" title="O‘chirish" @click="remove(a)">
                      <AppIcon name="trash" :size="18" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          <PaginationBar v-model="page" :meta="data.meta" />
        </template>
      </div>

      <aside class="a-card cats">
        <h2 class="a-card__title">Kategoriyalar</h2>
        <ul role="list" class="cats__list">
          <li v-for="c in categories" :key="c.id">
            <button type="button" class="cats__name" @click="editCategory(c)">
              {{ c.name }} <small>{{ c.count }}</small>
            </button>
            <button type="button" class="a-icon-btn a-icon-btn--danger" :aria-label="`${c.name} — o‘chirish`" @click="removeCategory(c)">
              <AppIcon name="trash" :size="16" />
            </button>
          </li>
        </ul>
        <form class="cats__form" @submit.prevent="saveCategory">
          <input v-model="catForm.name" class="a-input" placeholder="Yangi kategoriya" aria-label="Kategoriya nomi" />
          <input v-model="catForm.slug" class="a-input" placeholder="slug (ixtiyoriy)" aria-label="Slug" />
          <p v-if="catErrors.slug || catErrors.name" class="a-error">{{ catErrors.slug || catErrors.name }}</p>
          <div class="cats__btns">
            <button type="submit" class="a-btn a-btn--primary a-btn--sm" :disabled="!catForm.name.trim()">
              {{ catForm.id ? 'Saqlash' : 'Qo‘shish' }}
            </button>
            <button
              v-if="catForm.id"
              type="button"
              class="a-btn a-btn--ghost a-btn--sm"
              @click="catForm = { id: 0, name: '', slug: '', position: 0 }"
            >
              Bekor qilish
            </button>
          </div>
        </form>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 20px;
  align-items: start;
}

.cats__list {
  display: grid;
  gap: 2px;
  margin-bottom: 16px;
}

.cats__list li {
  display: flex;
  align-items: center;
}

.cats__name {
  flex: 1;
  display: flex;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: 10px;
  font-weight: 600;
  text-align: left;
}

.cats__name:hover {
  background: var(--surface-2);
}

.cats__name small {
  color: var(--ink-3);
}

.cats__form {
  display: grid;
  gap: 8px;
  padding-top: 16px;
  border-top: 1px solid var(--line);
}

.cats__form .a-input {
  height: 42px;
}

.cats__btns {
  display: flex;
  gap: 8px;
}

@media (max-width: 1100px) {
  .layout {
    grid-template-columns: 1fr;
  }
}
</style>
