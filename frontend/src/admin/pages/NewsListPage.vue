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
import { formatDate, formatNumber } from '@/utils/format'
import { toneStyle } from '@/modules/news/tone'
import type { News } from '@/types/api'

const ui = useUiStore()
const status = ref('')
const q = ref('')
const page = ref(1)

const { data, loading, error, reload } = useAsync(
  () => adminApi.news.list({ status: status.value, q: q.value, page: page.value, limit: 20 }),
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

async function remove(n: News) {
  const ok = await ui.confirm({ title: `“${n.title}” o‘chirilsinmi?`, text: 'Yangilik butunlay o‘chiriladi.' })
  if (!ok) return
  try {
    await adminApi.news.remove(n.id)
    ui.success('O‘chirildi')
    reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}
</script>

<template>
  <div>
    <PageHeader title="Yangiliklar" subtitle="Studiya xabarlari: yangi loyihalar, Labs, jamoa va hamkorliklar.">
      <RouterLink to="/admin/news/new" class="a-btn a-btn--primary"><AppIcon name="plus" :size="18" /> Yangi yangilik</RouterLink>
    </PageHeader>

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
        <input v-model="q" class="a-input" type="search" placeholder="Sarlavha yoki teg bo‘yicha" aria-label="Qidirish" />
      </div>
    </div>

    <AdminEmpty v-if="error" error title="Yuklab bo‘lmadi" :text="error.message" @retry="reload" />
    <div v-else-if="loading && !data" class="a-stack">
      <div v-for="n in 4" :key="n" class="skeleton" style="height: 72px; border-radius: 14px" />
    </div>
    <AdminEmpty
      v-else-if="!data?.items.length"
      icon="news"
      title="Yangilik topilmadi"
      text="Birinchi xabarni e’lon qiling — sayt tirik ko‘rinadi."
    >
      <RouterLink to="/admin/news/new" class="a-btn a-btn--primary a-btn--sm">Yangilik yozish</RouterLink>
    </AdminEmpty>

    <template v-else>
      <table class="a-table">
        <thead>
          <tr>
            <th>Sarlavha</th>
            <th>Teg</th>
            <th>Holat</th>
            <th style="text-align: right">Ko‘rishlar</th>
            <th>Sana</th>
            <th style="text-align: right">Amallar</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="n in data.items" :key="n.id">
            <td>
              <RouterLink :to="`/admin/news/${n.id}`">
                <span class="a-cell-title">
                  <AppIcon v-if="n.pinned" name="star" :size="15" class="pin" />
                  {{ n.title }}
                </span>
                <span class="a-cell-sub" style="display: block">/news/{{ n.slug }}</span>
              </RouterLink>
            </td>
            <td>
              <span v-if="n.tag" class="tag" :style="toneStyle(n.tag)">{{ n.tag }}</span>
              <span v-else class="a-cell-sub">—</span>
            </td>
            <td><StatusBadge :status="n.status" /></td>
            <td class="views"><AppIcon name="eye" :size="16" /> {{ formatNumber(n.views) }}</td>
            <td class="a-cell-sub">{{ formatDate(n.published_at ?? n.created_at) }}</td>
            <td>
              <div class="a-actions">
                <a v-if="n.status === 'published'" :href="`/news/${n.slug}`" target="_blank" class="a-icon-btn" title="Saytda ko‘rish"
                  ><AppIcon name="external" :size="18"
                /></a>
                <RouterLink :to="`/admin/news/${n.id}`" class="a-icon-btn" title="Tahrirlash"
                  ><AppIcon name="edit" :size="18"
                /></RouterLink>
                <button type="button" class="a-icon-btn a-icon-btn--danger" title="O‘chirish" @click="remove(n)">
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
</template>

<style scoped>
.pin {
  display: inline-block;
  vertical-align: -2px;
  margin-right: 4px;
  color: var(--gold);
}

.tag {
  display: inline-flex;
  padding: 4px 10px;
  border-radius: var(--r-pill);
  background: color-mix(in srgb, var(--tone) 14%, transparent);
  color: var(--tone);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.views {
  text-align: right;
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
  color: var(--ink-2);
}

.views :deep(svg) {
  vertical-align: -3px;
  margin-right: 4px;
  color: var(--ink-3);
}
</style>
