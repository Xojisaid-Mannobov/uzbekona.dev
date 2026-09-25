<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import StatusBadge from '@/admin/components/StatusBadge.vue'
import PaginationBar from '@/admin/components/PaginationBar.vue'
import { useAsync } from '@/composables/useAsync'
import { useDragSort } from '@/admin/composables/useDragSort'
import { useUiStore } from '@/admin/stores/ui'
import { adminApi } from '@/admin/services/adminApi'
import { toApiError } from '@/services/http'
import { timeAgo } from '@/utils/format'
import type { ContentStatus, Project } from '@/types/api'

const route = useRoute()
const router = useRouter()
const ui = useUiStore()

// Filtrlar URL'da saqlanadi — sahifani yangilaganda ham holat qoladi
const status = computed(() => (typeof route.query.status === 'string' ? route.query.status : ''))
const q = ref(typeof route.query.q === 'string' ? route.query.q : '')
const page = ref(1)

const { data, loading, error, reload } = useAsync(
  () => adminApi.projects.list({ status: status.value, q: q.value, page: page.value, limit: 50 }),
  {
    watch: [status, page],
  },
)

const items = ref<Project[]>([])
watch(data, (d) => (items.value = d?.items ?? []))

let timer: ReturnType<typeof setTimeout>
watch(q, (v) => {
  clearTimeout(timer)
  timer = setTimeout(() => {
    router.replace({ query: { ...route.query, q: v || undefined } })
    page.value = 1
    reload()
  }, 300)
})

const tabs = [
  { value: '', label: 'Barchasi' },
  { value: 'published', label: 'E’lon qilingan' },
  { value: 'draft', label: 'Qoralama' },
  { value: 'archived', label: 'Arxiv' },
]

function setStatus(v: string) {
  page.value = 1
  router.replace({ query: { ...route.query, status: v || undefined } })
}

// Tartiblash faqat filtrsiz ro'yxatda (aks holda pozitsiyalar chalkashadi)
const canSort = computed(() => !status.value && !q.value && (data.value?.meta.total ?? 0) <= 50)
const { handlers, dragIndex, overIndex } = useDragSort(items, saveOrder)

async function saveOrder() {
  try {
    await adminApi.projects.reorder(items.value.map((p) => p.id))
    ui.success('Tartib saqlandi')
  } catch (e) {
    ui.error(toApiError(e).message)
    reload()
  }
}

async function changeStatus(p: Project, s: ContentStatus) {
  try {
    await adminApi.projects.setStatus(p.id, s)
    p.status = s
    ui.success(s === 'published' ? 'E’lon qilindi' : s === 'archived' ? 'Arxivga o‘tkazildi' : 'Qoralamaga o‘tkazildi')
    if (status.value) reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}

async function toggleFeatured(p: Project) {
  try {
    await adminApi.projects.setFeatured(p.id, !p.featured)
    p.featured = !p.featured
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}

async function remove(p: Project) {
  const ok = await ui.confirm({
    title: `“${p.title}” o‘chirilsinmi?`,
    text: 'Loyiha, uning bloklari va galereyasi butunlay o‘chiriladi. Bu amalni qaytarib bo‘lmaydi.',
  })
  if (!ok) return
  try {
    await adminApi.projects.remove(p.id)
    ui.success('O‘chirildi')
    reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}
</script>

<template>
  <div>
    <PageHeader title="Loyihalar" subtitle="Case-study’lar, tartib va e’lon qilish holati.">
      <RouterLink to="/admin/projects/new" class="a-btn a-btn--primary"><AppIcon name="plus" :size="18" /> Yangi loyiha</RouterLink>
    </PageHeader>

    <div class="a-toolbar">
      <div class="a-tabs" role="tablist">
        <button
          v-for="t in tabs"
          :key="t.value"
          type="button"
          role="tab"
          class="a-tab"
          :class="{ 'is-active': status === t.value }"
          :aria-selected="status === t.value"
          @click="setStatus(t.value)"
        >
          {{ t.label }}
        </button>
      </div>
      <div class="a-search">
        <AppIcon name="search" :size="18" />
        <input v-model="q" class="a-input" type="search" placeholder="Nomi, mijoz yoki slug" aria-label="Qidirish" />
      </div>
    </div>

    <AdminEmpty v-if="error" error title="Loyihalarni yuklab bo‘lmadi" :text="error.message" @retry="reload" />

    <div v-else-if="loading && !data" class="a-stack">
      <div v-for="n in 4" :key="n" class="skeleton" style="height: 76px; border-radius: 14px" />
    </div>

    <AdminEmpty
      v-else-if="!items.length"
      icon="folder"
      :title="q || status ? 'Hech narsa topilmadi' : 'Hozircha loyiha yo‘q'"
      :text="q || status ? 'Boshqa filtr yoki so‘z bilan urinib ko‘ring.' : 'Birinchi case-study’ni yarating.'"
    >
      <RouterLink v-if="!q && !status" to="/admin/projects/new" class="a-btn a-btn--primary a-btn--sm">Birinchisini qo‘shish</RouterLink>
    </AdminEmpty>

    <template v-else>
      <table class="a-table">
        <thead>
          <tr>
            <th style="width: 44px"><span class="visually-hidden">Tartib</span></th>
            <th>Loyiha</th>
            <th>Holat</th>
            <th>Yil</th>
            <th>Yangilangan</th>
            <th style="text-align: right">Amallar</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(p, i) in items"
            :key="p.id"
            :class="{ 'is-dragging': dragIndex === i, 'is-drop-target': overIndex === i && dragIndex !== i }"
            v-bind="canSort ? handlers(i) : {}"
          >
            <td>
              <span v-if="canSort" class="a-handle" title="Sudrab tartiblang"><AppIcon name="grip" :size="18" /></span>
            </td>
            <td>
              <RouterLink :to="`/admin/projects/${p.id}`" class="a-row">
                <span class="a-thumb" :style="!p.cover ? { background: p.accent || 'var(--surface-2)' } : undefined">
                  <img v-if="p.cover" :src="p.cover.variants[0]?.url ?? p.cover.url" alt="" />
                </span>
                <span>
                  <span class="a-cell-title">{{ p.title }}</span>
                  <span class="a-cell-sub" style="display: block">/{{ p.slug }} · {{ p.short_description }}</span>
                </span>
              </RouterLink>
            </td>
            <td><StatusBadge :status="p.status" /></td>
            <td>{{ p.year ?? '—' }}</td>
            <td class="a-cell-sub">{{ timeAgo(p.updated_at) }}</td>
            <td>
              <div class="a-actions">
                <button
                  type="button"
                  class="a-icon-btn"
                  :class="{ 'is-active': p.featured }"
                  :title="p.featured ? 'Bosh sahifadan olish' : 'Bosh sahifaga (featured)'"
                  :aria-pressed="p.featured"
                  @click="toggleFeatured(p)"
                >
                  <AppIcon name="star" :size="18" />
                </button>
                <button
                  v-if="p.status !== 'published'"
                  type="button"
                  class="a-icon-btn"
                  title="E’lon qilish"
                  @click="changeStatus(p, 'published')"
                >
                  <AppIcon name="eye" :size="18" />
                </button>
                <button v-else type="button" class="a-icon-btn" title="Qoralamaga o‘tkazish" @click="changeStatus(p, 'draft')">
                  <AppIcon name="eye-off" :size="18" />
                </button>
                <button
                  v-if="p.status !== 'archived'"
                  type="button"
                  class="a-icon-btn"
                  title="Arxivlash"
                  @click="changeStatus(p, 'archived')"
                >
                  <AppIcon name="archive" :size="18" />
                </button>
                <a v-if="p.status === 'published'" :href="`/projects/${p.slug}`" target="_blank" class="a-icon-btn" title="Saytda ko‘rish"
                  ><AppIcon name="external" :size="18"
                /></a>
                <RouterLink :to="`/admin/projects/${p.id}`" class="a-icon-btn" title="Tahrirlash"
                  ><AppIcon name="edit" :size="18"
                /></RouterLink>
                <button type="button" class="a-icon-btn a-icon-btn--danger" title="O‘chirish" @click="remove(p)">
                  <AppIcon name="trash" :size="18" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
      <p v-if="canSort" class="a-hint" style="margin-top: 12px">
        Qatorlarni sudrab tartiblang — saytdagi tartib avtomatik saqlanadi. ★ — bosh sahifada ko‘rinadi.
      </p>
      <PaginationBar v-if="data" v-model="page" :meta="data.meta" />
    </template>
  </div>
</template>
