<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import BrandIcon from '@/components/ui/BrandIcon.vue'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import StatusBadge from '@/admin/components/StatusBadge.vue'
import PaginationBar from '@/admin/components/PaginationBar.vue'
import { useAsync } from '@/composables/useAsync'
import { useUiStore } from '@/admin/stores/ui'
import { adminApi } from '@/admin/services/adminApi'
import { toApiError } from '@/services/http'
import { formatDateTime, timeAgo } from '@/utils/format'
import type { ApplicationStatus, JobApplication } from '@/types/api'

// "Jamoaga qo'shilish" arizalari: nomzodni o'qib chiqish, holatini yuritish va bog'lanish
const route = useRoute()
const ui = useUiStore()
const status = ref<string>('new')
const q = ref('')
const page = ref(1)

const { data, loading, error, reload } = useAsync(
  () => adminApi.applications.list({ status: status.value, q: q.value, page: page.value, limit: 30 }),
  { watch: [status, page] },
)

let timer: ReturnType<typeof setTimeout>
watch(q, () => {
  clearTimeout(timer)
  timer = setTimeout(() => {
    page.value = 1
    reload()
  }, 300)
})

const selected = ref<JobApplication | null>(null)
const note = ref('')
const saving = ref(false)

// Dashboard'dan ?id= bilan kelinsa — shu ariza ochiladi
watch(data, (d) => {
  const id = Number(route.query.id)
  if (id && !selected.value) {
    const found = d?.items.find((a) => a.id === id)
    if (found) open(found)
    else if (status.value) status.value = ''
  }
})

function open(a: JobApplication) {
  selected.value = a
  note.value = a.note
}

async function update(s: ApplicationStatus) {
  if (!selected.value) return
  saving.value = true
  try {
    const updated = await adminApi.applications.update(selected.value.id, s, note.value)
    Object.assign(selected.value, updated)
    ui.success('Yangilandi')
    if (status.value && status.value !== s) reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  } finally {
    saving.value = false
  }
}

async function remove(a: JobApplication) {
  const ok = await ui.confirm({ title: 'Ariza o‘chirilsinmi?', text: `${a.full_name} arizasi butunlay o‘chiriladi.` })
  if (!ok) return
  try {
    await adminApi.applications.remove(a.id)
    selected.value = null
    ui.success('O‘chirildi')
    reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}

const tabs = [
  { v: 'new', l: 'Yangi' },
  { v: 'reviewing', l: 'Ko‘rib chiqilmoqda' },
  { v: 'interview', l: 'Suhbat' },
  { v: 'accepted', l: 'Qabul qilingan' },
  { v: 'rejected', l: 'Rad etilgan' },
  { v: '', l: 'Barchasi' },
]

// Holat bosqichlari — tugmalar tartibi nomzod yo'lini takrorlaydi
const flow: { v: ApplicationStatus; l: string }[] = [
  { v: 'reviewing', l: 'Ko‘rib chiqish' },
  { v: 'interview', l: 'Suhbatga chaqirish' },
  { v: 'accepted', l: 'Qabul qilish' },
]

const telegramLink = computed(() => {
  const t = selected.value?.telegram.trim().replace(/^@/, '') ?? ''
  return /^[A-Za-z0-9_]{4,}$/.test(t) ? `https://t.me/${t}` : null
})
const phoneLink = computed(() => {
  const p = selected.value?.phone.replace(/[^\d+]/g, '') ?? ''
  return p.length >= 9 ? `tel:${p}` : null
})
</script>

<template>
  <div>
    <PageHeader title="Nomzodlar" subtitle="“Jamoaga qo‘shilish” sahifasi orqali kelgan arizalar." />

    <div class="a-toolbar">
      <div class="a-tabs">
        <button
          v-for="t in tabs"
          :key="t.v"
          type="button"
          class="a-tab"
          :class="{ 'is-active': status === t.v }"
          @click="((status = t.v), (page = 1), (selected = null))"
        >
          {{ t.l }}
        </button>
      </div>
      <div class="a-search">
        <AppIcon name="search" :size="18" />
        <input v-model="q" class="a-input" type="search" placeholder="Ism, email, yo‘nalish yoki matn" aria-label="Qidirish" />
      </div>
    </div>

    <div class="layout">
      <div>
        <AdminEmpty v-if="error" error title="Yuklab bo‘lmadi" :text="error.message" @retry="reload" />
        <div v-else-if="loading && !data" class="a-stack">
          <div v-for="n in 5" :key="n" class="skeleton" style="height: 84px; border-radius: 14px" />
        </div>
        <AdminEmpty
          v-else-if="!data?.items.length"
          icon="briefcase"
          :title="status === 'new' ? 'Yangi ariza yo‘q' : 'Ariza topilmadi'"
          text="Saytdagi “Jamoaga qo‘shilish” formasidan kelgan arizalar shu yerda ko‘rinadi."
        />
        <template v-else>
          <ul class="list" role="list">
            <li v-for="a in data.items" :key="a.id">
              <button
                type="button"
                class="app"
                :class="{ 'is-active': selected?.id === a.id, 'is-new': a.status === 'new' }"
                @click="open(a)"
              >
                <span class="app__top">
                  <strong>{{ a.full_name }}</strong>
                  <small>{{ timeAgo(a.created_at) }}</small>
                </span>
                <span class="app__msg">{{ a.about }}</span>
                <span class="app__tags">
                  <StatusBadge :status="a.status" />
                  <span class="app__tag app__tag--pos">{{ a.position }}</span>
                  <span v-if="a.experience" class="app__tag">{{ a.experience }}</span>
                </span>
              </button>
            </li>
          </ul>
          <PaginationBar v-model="page" :meta="data.meta" />
        </template>
      </div>

      <aside class="a-card detail">
        <template v-if="selected">
          <div class="detail__head">
            <div>
              <h2>{{ selected.full_name }}</h2>
              <p class="a-hint">{{ selected.position }} · {{ formatDateTime(selected.created_at) }} · #{{ selected.id }}</p>
            </div>
            <StatusBadge :status="selected.status" />
          </div>

          <dl class="detail__meta">
            <div>
              <dt>Email</dt>
              <dd>
                <a :href="`mailto:${selected.email}`" class="link">{{ selected.email }}</a>
              </dd>
            </div>
            <div v-if="selected.phone">
              <dt>Telefon</dt>
              <dd>{{ selected.phone }} <a v-if="phoneLink" :href="phoneLink" class="link">Qo‘ng‘iroq</a></dd>
            </div>
            <div v-if="selected.telegram">
              <dt>Telegram</dt>
              <dd>
                <a v-if="telegramLink" :href="telegramLink" target="_blank" rel="noopener noreferrer" class="tg">
                  <BrandIcon name="telegram" :size="16" /> {{ selected.telegram }}
                </a>
                <template v-else>{{ selected.telegram }}</template>
              </dd>
            </div>
            <div v-if="selected.experience">
              <dt>Tajriba</dt>
              <dd>{{ selected.experience }}</dd>
            </div>
            <div v-if="selected.portfolio_url">
              <dt>Portfolio</dt>
              <dd>
                <a :href="selected.portfolio_url" target="_blank" rel="noopener noreferrer nofollow" class="link url">{{
                  selected.portfolio_url
                }}</a>
              </dd>
            </div>
            <div v-if="selected.resume_url">
              <dt>Rezyume</dt>
              <dd>
                <a :href="selected.resume_url" target="_blank" rel="noopener noreferrer nofollow" class="link url">{{
                  selected.resume_url
                }}</a>
              </dd>
            </div>
          </dl>

          <div>
            <p class="a-label">O‘zi haqida</p>
            <p class="detail__msg">{{ selected.about }}</p>
          </div>

          <div class="a-field">
            <label class="a-label" for="note">Ichki izoh <small>faqat adminlar ko‘radi</small></label>
            <textarea id="note" v-model="note" class="a-textarea" placeholder="Suhbat sanasi, topshiriq natijasi, taassurot…" />
          </div>

          <div class="detail__status">
            <button
              v-for="f in flow"
              :key="f.v"
              type="button"
              class="a-btn a-btn--sm"
              :class="selected.status === f.v ? 'a-btn--primary' : 'a-btn--outline'"
              :disabled="saving"
              @click="update(f.v)"
            >
              {{ f.l }}
            </button>
            <button
              type="button"
              class="a-btn a-btn--sm"
              :class="selected.status === 'rejected' ? 'a-btn--primary' : 'a-btn--ghost'"
              :disabled="saving"
              @click="update('rejected')"
            >
              Rad etish
            </button>
          </div>
          <div class="detail__foot">
            <button
              type="button"
              class="a-btn a-btn--outline a-btn--sm"
              :disabled="saving || note === selected.note"
              @click="update(selected.status)"
            >
              Izohni saqlash
            </button>
            <a
              :href="`mailto:${selected.email}?subject=${encodeURIComponent('Uzbekona.dev — arizangiz bo‘yicha')}`"
              class="a-btn a-btn--outline a-btn--sm"
            >
              <AppIcon name="send" :size="16" /> Javob yozish
            </a>
            <button type="button" class="a-btn a-btn--danger a-btn--sm" @click="remove(selected)">
              <AppIcon name="trash" :size="16" /> O‘chirish
            </button>
          </div>
        </template>
        <div v-else class="detail__empty">
          <AppIcon name="briefcase" :size="28" />
          <p>Nomzodni tanlang — arizasi, havolalari va holatini boshqarish shu yerda.</p>
        </div>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 440px;
  gap: 20px;
  align-items: start;
}

.list {
  display: grid;
  gap: 8px;
}

.app {
  width: 100%;
  display: grid;
  gap: 6px;
  padding: 16px 18px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--line);
  text-align: left;
  transition: border-color var(--dur-fast) var(--ease);
}

.app:hover {
  border-color: var(--line-strong);
}

.app.is-active {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-soft);
}

.app.is-new strong::before {
  content: '';
  display: inline-block;
  width: 8px;
  height: 8px;
  margin-right: 8px;
  border-radius: 50%;
  background: var(--accent);
  vertical-align: middle;
}

.app__top {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.app__top small {
  color: var(--ink-3);
  font-size: 12px;
  white-space: nowrap;
}

.app__msg {
  color: var(--ink-2);
  font-size: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.app__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 4px;
}

.app__tag {
  font-size: 12px;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 99px;
  background: var(--surface-2);
  color: var(--ink-2);
}

.app__tag--pos {
  background: var(--accent-soft);
  color: var(--accent);
}

.detail {
  position: sticky;
  top: calc(var(--a-top-h) + 20px);
  display: grid;
  gap: 18px;
}

.detail__head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 12px;
}

.detail__head h2 {
  font-size: 22px;
  letter-spacing: -0.02em;
}

.detail__meta {
  display: grid;
  gap: 8px;
  font-size: 14px;
}

.detail__meta div {
  display: grid;
  grid-template-columns: 100px minmax(0, 1fr);
  gap: 8px;
}

.detail__meta dt {
  color: var(--ink-3);
}

.detail__meta dd {
  font-weight: 600;
  min-width: 0;
}

.link {
  color: var(--accent);
  font-weight: 700;
}

.url {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tg {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--accent);
  font-weight: 700;
}

.detail__msg {
  margin-top: 8px;
  padding: 16px;
  border-radius: 14px;
  background: var(--surface-2);
  white-space: pre-wrap;
  line-height: 1.6;
  max-height: 320px;
  overflow: auto;
}

.detail__status,
.detail__foot {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.detail__foot {
  padding-top: 14px;
  border-top: 1px solid var(--line);
}

.detail__empty {
  display: grid;
  justify-items: center;
  gap: 12px;
  padding: 48px 12px;
  text-align: center;
  color: var(--ink-3);
}

@media (max-width: 1200px) {
  .layout {
    grid-template-columns: 1fr;
  }

  .detail {
    position: static;
  }
}
</style>
