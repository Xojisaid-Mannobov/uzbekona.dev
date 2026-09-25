<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import TrendChart from '@/admin/components/TrendChart.vue'
import { useAsync } from '@/composables/useAsync'
import { adminApi } from '@/admin/services/adminApi'
import { formatNumber } from '@/utils/format'

// Sayt statistikasi: davr tanlanadi (7/30/90 kun) — barcha kartalar va ro'yxatlar shu davr bo'yicha.
const PERIODS = [
  { days: 7, label: '7 kun' },
  { days: 30, label: '30 kun' },
  { days: 90, label: '90 kun' },
]
const days = ref(30)
const { data, loading, error, reload } = useAsync(() => adminApi.analytics(days.value), { watch: [days] })

const fmt = formatNumber
const compact = (n: number) => (n >= 10_000 ? `${(n / 1000).toFixed(n >= 100_000 ? 0 : 1)}K` : fmt(n))

/** Oldingi shu uzunlikdagi davrga nisbatan o'zgarish */
function delta(cur: number, prev: number) {
  if (!prev) return cur ? { text: 'yangi', dir: 'up' as const } : null
  const p = Math.round(((cur - prev) / prev) * 100)
  return { text: `${p > 0 ? '+' : ''}${p}%`, dir: p > 0 ? ('up' as const) : p < 0 ? ('down' as const) : ('flat' as const) }
}

const tiles = computed(() => {
  const s = data.value?.summary
  if (!s) return []
  return [
    { label: 'Ko‘rishlar', hint: 'ochilgan sahifalar', value: s.views, delta: delta(s.views, s.prev_views) },
    { label: 'Tashriflar', hint: 'saytga kirishlar (sessiyalar)', value: s.visits, delta: delta(s.visits, s.prev_visits) },
    { label: 'Tashrif buyuruvchilar', hint: 'noyob odamlar', value: s.visitors, delta: delta(s.visitors, s.prev_visitors) },
  ]
})

// Sahifa yo'llarini tushunarli nomlarga aylantirish
const PAGE_NAMES: Record<string, string> = {
  '/': 'Bosh sahifa',
  '/projects': 'Loyihalar',
  '/services': 'Xizmatlar',
  '/about': 'Biz haqimizda',
  '/team': 'Jamoa',
  '/news': 'Yangiliklar',
  '/journal': 'Journal',
  '/contact': 'Bog‘lanish',
}
function pageName(path: string) {
  if (PAGE_NAMES[path]) return PAGE_NAMES[path]
  const [, section, slug] = path.split('/')
  return `${PAGE_NAMES[`/${section}`] ?? section} · ${slug}`
}

const SOURCES: Record<string, string> = {
  direct: 'To‘g‘ridan-to‘g‘ri',
  't.me': 'Telegram',
  'web.telegram.org': 'Telegram Web',
  'google.com': 'Google',
  'yandex.ru': 'Yandex',
  'yandex.uz': 'Yandex',
  'instagram.com': 'Instagram',
  'l.instagram.com': 'Instagram',
  'facebook.com': 'Facebook',
  'l.facebook.com': 'Facebook',
  'linkedin.com': 'LinkedIn',
  'github.com': 'GitHub',
}
const DEVICES: Record<string, string> = { desktop: 'Kompyuter', mobile: 'Telefon', tablet: 'Planshet' }
const DEVICE_COLORS: Record<string, string> = { desktop: 'var(--chart-1)', mobile: 'var(--chart-2)', tablet: 'var(--chart-3)' }

const maxPage = computed(() => Math.max(1, ...(data.value?.top_pages.map((p) => p.views) ?? [1])))
const maxRef = computed(() => Math.max(1, ...(data.value?.referrers.map((p) => p.views) ?? [1])))
const deviceTotal = computed(() => data.value?.devices.reduce((a, d) => a + d.visitors, 0) ?? 0)
const pct = (n: number, total: number) => (total ? Math.round((n / total) * 100) : 0)
</script>

<template>
  <div>
    <PageHeader title="Statistika" subtitle="Saytga nechta tashrif bo‘ldi va necha kishi ko‘rdi — cookie’siz, o‘z serverimizda.">
      <div class="periods" role="group" aria-label="Davr">
        <button
          v-for="p in PERIODS"
          :key="p.days"
          type="button"
          class="a-tab"
          :class="{ 'is-active': days === p.days }"
          :aria-pressed="days === p.days"
          @click="days = p.days"
        >
          {{ p.label }}
        </button>
      </div>
    </PageHeader>

    <AdminEmpty v-if="error && !data" error title="Statistikani yuklab bo‘lmadi" :text="error.message" @retry="reload" />
    <div v-else-if="!data" class="a-stack">
      <div class="skeleton" style="height: 120px; border-radius: 20px" />
      <div class="skeleton" style="height: 360px; border-radius: 20px" />
    </div>

    <!-- Qayta yuklanayotganda eski ko'rinish xiralashadi, sakrash bo'lmaydi -->
    <div v-else class="analytics" :class="{ 'is-loading': loading }">
      <div class="tiles">
        <article v-for="t in tiles" :key="t.label" class="a-card tile">
          <p class="tile__label">{{ t.label }}</p>
          <p class="tile__value">{{ compact(t.value) }}</p>
          <p class="tile__foot">
            <span v-if="t.delta" class="tile__delta" :class="`is-${t.delta.dir}`">
              <AppIcon v-if="t.delta.dir !== 'flat'" :name="t.delta.dir === 'up' ? 'arrow-up' : 'arrow-down'" :size="14" />
              {{ t.delta.text }}
            </span>
            <span>{{ t.hint }}</span>
          </p>
        </article>

        <article class="a-card tile tile--live">
          <p class="tile__label"><span class="live-dot" aria-hidden="true" /> Hozir saytda</p>
          <p class="tile__value">{{ fmt(data.summary.online_now) }}</p>
          <p class="tile__foot">
            <span>Bugun: {{ fmt(data.summary.today_visitors) }} kishi · {{ fmt(data.summary.today_views) }} ko‘rish</span>
          </p>
        </article>
      </div>

      <section class="a-card">
        <div class="card-head">
          <h2 class="a-card__title">Kunlik dinamika — oxirgi {{ data.days }} kun</h2>
          <span class="a-hint">
            Hamma vaqt: {{ fmt(data.summary.all_time_visitors) }} kishi · {{ fmt(data.summary.all_time_views) }} ko‘rish
          </span>
        </div>
        <TrendChart :days="data.by_day" />
      </section>

      <div class="grid">
        <section class="a-card">
          <h2 class="a-card__title">Eng ko‘p ko‘rilgan sahifalar</h2>
          <p v-if="!data.top_pages.length" class="a-hint">Bu davrda hali ko‘rishlar yo‘q.</p>
          <ol v-else role="list" class="bars">
            <li v-for="p in data.top_pages" :key="p.key">
              <div class="bars__text">
                <a :href="p.key" target="_blank" class="bars__name">{{ pageName(p.key) }}</a>
                <span class="bars__value"
                  >{{ fmt(p.views) }} <small>· {{ fmt(p.visitors) }} kishi</small></span
                >
              </div>
              <span class="bars__track"><span class="bars__fill" :style="{ width: `${(p.views / maxPage) * 100}%` }" /></span>
            </li>
          </ol>
        </section>

        <div class="a-stack">
          <section class="a-card">
            <h2 class="a-card__title">Qayerdan kelishdi</h2>
            <p v-if="!data.referrers.length" class="a-hint">Manbalar ma’lumoti hali yo‘q.</p>
            <ol v-else role="list" class="bars">
              <li v-for="r in data.referrers" :key="r.key">
                <div class="bars__text">
                  <span class="bars__name">{{ SOURCES[r.key] ?? r.key }}</span>
                  <span class="bars__value">{{ fmt(r.views) }} <small>tashrif</small></span>
                </div>
                <span class="bars__track"><span class="bars__fill" :style="{ width: `${(r.views / maxRef) * 100}%` }" /></span>
              </li>
            </ol>
          </section>

          <section class="a-card">
            <h2 class="a-card__title">Qurilmalar</h2>
            <p v-if="!deviceTotal" class="a-hint">Ma’lumot yo‘q.</p>
            <template v-else>
              <div
                class="stack"
                role="img"
                :aria-label="data.devices.map((d) => `${DEVICES[d.key]} ${pct(d.visitors, deviceTotal)}%`).join(', ')"
              >
                <span
                  v-for="d in data.devices"
                  :key="d.key"
                  class="stack__seg"
                  :style="{ flexGrow: d.visitors, background: DEVICE_COLORS[d.key] }"
                  :title="`${DEVICES[d.key]}: ${fmt(d.visitors)} kishi`"
                />
              </div>
              <ul role="list" class="stack__legend">
                <li v-for="d in data.devices" :key="d.key">
                  <i :style="{ background: DEVICE_COLORS[d.key] }" />
                  {{ DEVICES[d.key] }}
                  <strong>{{ pct(d.visitors, deviceTotal) }}%</strong>
                  <small>{{ fmt(d.visitors) }} kishi</small>
                </li>
              </ul>
            </template>
          </section>

          <section class="a-card">
            <div class="card-head">
              <h2 class="a-card__title">Eng ko‘p o‘qilgan yangiliklar</h2>
              <RouterLink to="/admin/news" class="a-hint">Barchasi →</RouterLink>
            </div>
            <p v-if="!data.top_news.length" class="a-hint">E’lon qilingan yangilik yo‘q.</p>
            <ol v-else role="list" class="top-news">
              <li v-for="(n, i) in data.top_news" :key="n.id">
                <span class="top-news__rank">{{ i + 1 }}</span>
                <RouterLink :to="`/admin/news/${n.id}`" class="top-news__title">{{ n.title }}</RouterLink>
                <span class="top-news__views"><AppIcon name="eye" :size="15" /> {{ fmt(n.views) }}</span>
              </li>
            </ol>
          </section>
        </div>
      </div>

      <p class="a-hint privacy">
        <AppIcon name="eye-off" :size="16" /> IP manzil va brauzer ma’lumotlari saqlanmaydi: tashrif buyuruvchi anonim xesh bilan
        hisoblanadi. Robotlar, admin sahifalari va “Do Not Track” yoqilgan brauzerlar hisobga olinmaydi.
      </p>
    </div>
  </div>
</template>

<style scoped>
.periods {
  display: flex;
  gap: 4px;
  padding: 4px;
  border-radius: 999px;
  background: var(--surface);
  border: 1px solid var(--line);
}

.analytics {
  display: grid;
  gap: 16px;
  transition: opacity 0.2s ease;
}

.analytics.is-loading {
  opacity: 0.55;
}

.analytics .a-card + .a-card {
  margin-top: 0;
}

.tiles {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.tile__label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink-2);
}

.tile__value {
  margin: 10px 0 8px;
  font-size: 40px;
  font-weight: 700;
  line-height: 1;
  letter-spacing: -0.04em;
}

.tile__foot {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px 10px;
  font-size: 12px;
  color: var(--ink-3);
}

.tile__delta {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 2px 8px;
  border-radius: 999px;
  font-weight: 700;
  color: var(--ink-2);
  background: var(--surface-2);
}

.tile__delta.is-up {
  color: var(--chart-good);
  background: color-mix(in srgb, var(--chart-good) 12%, transparent);
}

.tile__delta.is-down {
  color: var(--chart-bad);
  background: color-mix(in srgb, var(--chart-bad) 12%, transparent);
}

.live-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--chart-good);
  animation: live 1.8s ease-out infinite;
}

@keyframes live {
  from {
    box-shadow: 0 0 0 0 color-mix(in srgb, var(--chart-good) 55%, transparent);
  }
  to {
    box-shadow: 0 0 0 8px transparent;
  }
}

.card-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  flex-wrap: wrap;
  gap: 8px 16px;
  margin-bottom: 14px;
}

.card-head .a-card__title {
  margin-bottom: 0;
}

.grid {
  display: grid;
  grid-template-columns: minmax(0, 1.3fr) minmax(0, 1fr);
  gap: 16px;
  align-items: start;
}

/* Gorizontal ustunlar: nom va qiymat tepada, chiziq ostida */
.bars {
  display: grid;
  gap: 14px;
}

.bars__text {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 6px;
  font-size: 14px;
}

.bars__name {
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

a.bars__name:hover {
  color: var(--accent);
}

.bars__value {
  flex-shrink: 0;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
}

.bars__value small {
  font-weight: 500;
  color: var(--ink-3);
}

.bars__track {
  display: block;
  height: 8px;
  border-radius: 4px;
  background: var(--surface-2);
}

.bars__fill {
  display: block;
  height: 100%;
  min-width: 4px;
  border-radius: 4px;
  background: var(--chart-1);
}

/* Qurilmalar: 100% ustun, segmentlar orasida 2px bo'shliq */
.stack {
  display: flex;
  gap: 2px;
  height: 14px;
  border-radius: 4px;
  overflow: hidden;
}

.stack__seg {
  flex-basis: 0;
  min-width: 4px;
}

.stack__legend {
  display: grid;
  gap: 8px;
  margin-top: 14px;
  font-size: 14px;
}

.stack__legend li {
  display: flex;
  align-items: center;
  gap: 8px;
}

.stack__legend i {
  width: 10px;
  height: 10px;
  border-radius: 3px;
}

.stack__legend strong {
  margin-left: auto;
}

.stack__legend small {
  min-width: 70px;
  text-align: right;
  color: var(--ink-3);
}

.top-news {
  display: grid;
  gap: 4px;
}

.top-news li {
  display: grid;
  grid-template-columns: 24px minmax(0, 1fr) auto;
  gap: 10px;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid var(--line);
  font-size: 14px;
}

.top-news li:last-child {
  border-bottom: 0;
}

.top-news__rank {
  font-weight: 700;
  color: var(--ink-3);
}

.top-news__title {
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.top-news__title:hover {
  color: var(--accent);
}

.top-news__views {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
}

.privacy {
  display: flex;
  gap: 8px;
  align-items: flex-start;
  max-width: 820px;
}

@media (max-width: 1100px) {
  .tiles {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 560px) {
  .tiles {
    grid-template-columns: 1fr;
  }
}
</style>
