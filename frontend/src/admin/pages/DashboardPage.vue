<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import type { IconName } from '@/components/ui/icons'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import StatusBadge from '@/admin/components/StatusBadge.vue'
import { useAsync } from '@/composables/useAsync'
import { useAuthStore } from '@/admin/stores/auth'
import { adminApi } from '@/admin/services/adminApi'
import { timeAgo } from '@/utils/format'

const auth = useAuthStore()
const { data: stats, loading, error, reload } = useAsync(adminApi.dashboard)

const cards = computed<{ label: string; value: number; icon: IconName; to: string; accent?: boolean }[]>(() => {
  const s = stats.value
  if (!s) return []
  return [
    { label: 'Bugun saytda (kishi)', value: s.visitors_today, icon: 'pulse', to: '/admin/analytics' },
    { label: 'Faol loyihalar', value: s.active_projects, icon: 'folder', to: '/admin/projects' },
    { label: 'E’lon qilingan', value: s.published_projects, icon: 'eye', to: '/admin/projects?status=published' },
    { label: 'Yangiliklar', value: s.news, icon: 'news', to: '/admin/news' },
    { label: 'Maqolalar', value: s.articles, icon: 'file-text', to: '/admin/articles' },
    { label: 'Yangi so‘rovlar', value: s.incoming_requests, icon: 'inbox', to: '/admin/requests', accent: s.incoming_requests > 0 },
  ]
})

const maxDay = computed(() => Math.max(1, ...(stats.value?.requests_by_day.map((d) => d.count) ?? [1])))
const greeting = computed(() => {
  const h = new Date().getHours()
  return h < 12 ? 'Xayrli tong' : h < 18 ? 'Xayrli kun' : 'Xayrli kech'
})
</script>

<template>
  <div>
    <PageHeader :title="`${greeting}, ${auth.admin?.name?.split(' ')[0] ?? ''}`" subtitle="Saytdagi so‘nggi holat bir qarashda.">
      <RouterLink to="/admin/projects/new" class="a-btn a-btn--primary"><AppIcon name="plus" :size="18" /> Yangi loyiha</RouterLink>
    </PageHeader>

    <AdminEmpty v-if="error" error title="Ma’lumotlarni yuklab bo‘lmadi" :text="error.message" @retry="reload" />

    <template v-else>
      <div class="stats">
        <template v-if="loading && !stats">
          <div v-for="n in 6" :key="n" class="skeleton stat-sk" />
        </template>
        <RouterLink v-for="c in cards" :key="c.label" :to="c.to" class="stat" :class="{ 'stat--accent': c.accent }">
          <span class="stat__icon"><AppIcon :name="c.icon" :size="20" /></span>
          <span class="stat__value">{{ c.value }}</span>
          <span class="stat__label">{{ c.label }}</span>
        </RouterLink>
      </div>

      <div v-if="stats" class="grid">
        <section class="a-card chart-card">
          <div class="card-head">
            <h2 class="a-card__title">So‘rovlar — oxirgi 14 kun</h2>
            <span class="a-hint">Jami: {{ stats.total_requests }}</span>
          </div>
          <div class="chart" role="img" :aria-label="`Oxirgi 14 kunda ${stats.requests_by_day.reduce((a, d) => a + d.count, 0)} ta so‘rov`">
            <div v-for="d in stats.requests_by_day" :key="d.day" class="chart__col" :title="`${d.day}: ${d.count}`">
              <span class="chart__bar" :style="{ height: `${(d.count / maxDay) * 100}%` }" :class="{ 'is-zero': !d.count }" />
              <span class="chart__day">{{ d.day.slice(8) }}</span>
            </div>
          </div>
        </section>

        <section class="a-card">
          <h2 class="a-card__title">Tizim</h2>
          <ul class="sys" role="list">
            <li>
              <RouterLink to="/admin/media"
                ><AppIcon name="image" :size="18" /> Media fayllar <strong>{{ stats.media }}</strong></RouterLink
              >
            </li>
            <li>
              <RouterLink to="/admin/users"
                ><AppIcon name="user" :size="18" /> Adminlar <strong>{{ stats.admins }}</strong></RouterLink
              >
            </li>
            <li>
              <RouterLink to="/admin/articles"
                ><AppIcon name="file-text" :size="18" /> Maqolalar <strong>{{ stats.articles }}</strong></RouterLink
              >
            </li>
            <li>
              <a href="/" target="_blank"><AppIcon name="external" :size="18" /> Public sayt <strong>↗</strong></a>
            </li>
          </ul>
        </section>

        <section class="a-card">
          <div class="card-head">
            <h2 class="a-card__title">So‘nggi so‘rovlar</h2>
            <RouterLink to="/admin/requests" class="more">Barchasi</RouterLink>
          </div>
          <p v-if="!stats.recent_requests.length" class="a-hint">
            Hozircha so‘rov yo‘q. Contact forma orqali kelgan xabarlar shu yerda ko‘rinadi.
          </p>
          <ul v-else class="list" role="list">
            <li v-for="r in stats.recent_requests" :key="r.id">
              <RouterLink :to="{ path: '/admin/requests', query: { id: r.id } }">
                <span class="list__main">
                  <strong>{{ r.name }}</strong>
                  <span>{{ r.project_type || r.message.slice(0, 60) }}</span>
                </span>
                <span class="list__meta">
                  <StatusBadge :status="r.status" />
                  <small>{{ timeAgo(r.created_at) }}</small>
                </span>
              </RouterLink>
            </li>
          </ul>
        </section>

        <section class="a-card">
          <div class="card-head">
            <h2 class="a-card__title">Oxirgi o‘zgarishlar</h2>
            <RouterLink to="/admin/projects" class="more">Loyihalar</RouterLink>
          </div>
          <ul class="list" role="list">
            <li v-for="p in stats.recent_projects" :key="p.id">
              <RouterLink :to="`/admin/projects/${p.id}`">
                <span class="list__main">
                  <strong>{{ p.title }}</strong>
                  <span>{{ p.short_description || p.tagline }}</span>
                </span>
                <span class="list__meta">
                  <StatusBadge :status="p.status" />
                  <small>{{ timeAgo(p.updated_at) }}</small>
                </span>
              </RouterLink>
            </li>
          </ul>
        </section>
      </div>
    </template>
  </div>
</template>

<style scoped>
.stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(170px, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.stat-sk {
  height: 150px;
  border-radius: var(--a-radius-lg);
}

.stat {
  display: grid;
  gap: 4px;
  padding: 22px;
  border-radius: var(--a-radius-lg);
  background: var(--surface);
  border: 1px solid var(--line);
  transition:
    border-color var(--dur-fast) var(--ease),
    transform var(--dur-fast) var(--ease);
}

.stat:hover {
  border-color: var(--line-strong);
  transform: translateY(-2px);
}

.stat__icon {
  display: grid;
  place-items: center;
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: var(--surface-2);
  margin-bottom: 14px;
}

.stat__value {
  font-size: 40px;
  font-weight: 700;
  letter-spacing: -0.05em;
  line-height: 1;
}

.stat__label {
  color: var(--ink-2);
  font-weight: 600;
  font-size: 14px;
}

.stat--accent {
  background: var(--accent);
  border-color: var(--accent);
  color: #fff;
}

.stat--accent .stat__icon {
  background: rgb(255 255 255 / 0.18);
}

.stat--accent .stat__label {
  color: rgb(255 255 255 / 0.85);
}

.grid {
  display: grid;
  grid-template-columns: minmax(0, 2fr) minmax(0, 1fr);
  gap: 16px;
}

.grid .a-card + .a-card {
  margin-top: 0;
}

.card-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
}

.more {
  font-size: 13px;
  font-weight: 700;
  color: var(--accent);
}

.chart {
  display: grid;
  grid-template-columns: repeat(14, 1fr);
  gap: 8px;
  height: 180px;
  align-items: end;
}

.chart__col {
  display: grid;
  grid-template-rows: 1fr auto;
  height: 100%;
  gap: 8px;
  justify-items: center;
}

.chart__bar {
  align-self: end;
  width: 100%;
  max-width: 28px;
  min-height: 4px;
  border-radius: 8px 8px 4px 4px;
  background: var(--accent);
}

.chart__bar.is-zero {
  background: var(--surface-3);
}

.chart__day {
  font-size: 11px;
  color: var(--ink-3);
  font-variant-numeric: tabular-nums;
}

.sys {
  display: grid;
  gap: 4px;
}

.sys a {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  font-weight: 600;
}

.sys a:hover {
  background: var(--surface-2);
}

.sys strong {
  margin-left: auto;
}

.list {
  display: grid;
}

.list a {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
  border-bottom: 1px solid var(--line);
}

.list li:last-child a {
  border-bottom: 0;
}

.list a:hover strong {
  color: var(--accent);
}

.list__main {
  display: grid;
  min-width: 0;
}

.list__main span {
  font-size: 13px;
  color: var(--ink-3);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.list__meta {
  display: grid;
  justify-items: end;
  gap: 4px;
  flex-shrink: 0;
}

.list__meta small {
  font-size: 12px;
  color: var(--ink-3);
}

@media (max-width: 1100px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
