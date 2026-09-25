<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import StatusBadge from '@/admin/components/StatusBadge.vue'
import PaginationBar from '@/admin/components/PaginationBar.vue'
import { useAsync } from '@/composables/useAsync'
import { useUiStore } from '@/admin/stores/ui'
import { adminApi } from '@/admin/services/adminApi'
import { toApiError } from '@/services/http'
import { formatDateTime, timeAgo } from '@/utils/format'
import type { ContactRequest, ContactStatus } from '@/types/api'

const route = useRoute()
const ui = useUiStore()
const status = ref<string>('new')
const q = ref('')
const page = ref(1)

const { data, loading, error, reload } = useAsync(
  () => adminApi.requests.list({ status: status.value, q: q.value, page: page.value, limit: 30 }),
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

const selected = ref<ContactRequest | null>(null)
const note = ref('')
const saving = ref(false)

// Dashboard'dan ?id= bilan kelinsa — shu so'rov ochiladi
watch(data, (d) => {
  const id = Number(route.query.id)
  if (id && !selected.value) {
    const found = d?.items.find((r) => r.id === id)
    if (found) open(found)
    else if (status.value) status.value = ''
  }
})

function open(r: ContactRequest) {
  selected.value = r
  note.value = r.note
}

async function update(s: ContactStatus) {
  if (!selected.value) return
  saving.value = true
  try {
    const updated = await adminApi.requests.update(selected.value.id, s, note.value)
    Object.assign(selected.value, updated)
    ui.success('Yangilandi')
    if (status.value && status.value !== s) reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  } finally {
    saving.value = false
  }
}

async function remove(r: ContactRequest) {
  const ok = await ui.confirm({ title: 'So‘rov o‘chirilsinmi?', text: `${r.name} dan kelgan so‘rov butunlay o‘chiriladi.` })
  if (!ok) return
  try {
    await adminApi.requests.remove(r.id)
    selected.value = null
    ui.success('O‘chirildi')
    reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}

const tabs = [
  { v: 'new', l: 'Yangi' },
  { v: 'in_progress', l: 'Jarayonda' },
  { v: 'done', l: 'Yakunlangan' },
  { v: 'spam', l: 'Spam' },
  { v: '', l: 'Barchasi' },
]

const telegramLink = computed(() => {
  const c = selected.value?.contact.trim() ?? ''
  if (c.startsWith('@')) return `https://t.me/${c.slice(1)}`
  if (/^t\.me\//.test(c)) return `https://${c}`
  return null
})
const phoneLink = computed(() => {
  const c = selected.value?.contact.replace(/[^\d+]/g, '') ?? ''
  return c.length >= 9 ? `tel:${c}` : null
})
</script>

<template>
  <div>
    <PageHeader title="So‘rovlar" subtitle="Contact forma orqali kelgan murojaatlar." />

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
        <input v-model="q" class="a-input" type="search" placeholder="Ism, email, telefon yoki matn" aria-label="Qidirish" />
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
          icon="inbox"
          :title="status === 'new' ? 'Yangi so‘rov yo‘q' : 'So‘rov topilmadi'"
          text="Saytdagi contact forma orqali kelgan murojaatlar shu yerda ko‘rinadi."
        />
        <template v-else>
          <ul class="list" role="list">
            <li v-for="r in data.items" :key="r.id">
              <button
                type="button"
                class="req"
                :class="{ 'is-active': selected?.id === r.id, 'is-new': r.status === 'new' }"
                @click="open(r)"
              >
                <span class="req__top">
                  <strong>{{ r.name }}</strong>
                  <small>{{ timeAgo(r.created_at) }}</small>
                </span>
                <span class="req__msg">{{ r.message }}</span>
                <span class="req__tags">
                  <StatusBadge :status="r.status" />
                  <span v-if="r.project_type" class="req__tag">{{ r.project_type }}</span>
                  <span v-if="r.budget" class="req__tag">{{ r.budget }}</span>
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
              <h2>{{ selected.name }}</h2>
              <p class="a-hint">{{ formatDateTime(selected.created_at) }} · #{{ selected.id }}</p>
            </div>
            <StatusBadge :status="selected.status" />
          </div>

          <dl class="detail__meta">
            <div v-if="selected.contact">
              <dt>Aloqa</dt>
              <dd>
                {{ selected.contact }}
                <a v-if="telegramLink" :href="telegramLink" target="_blank" class="link">Telegram ↗</a>
                <a v-else-if="phoneLink" :href="phoneLink" class="link">Qo‘ng‘iroq</a>
              </dd>
            </div>
            <div v-if="selected.email">
              <dt>Email</dt>
              <dd>
                <a :href="`mailto:${selected.email}`" class="link">{{ selected.email }}</a>
              </dd>
            </div>
            <div v-if="selected.project_type">
              <dt>Loyiha turi</dt>
              <dd>{{ selected.project_type }}</dd>
            </div>
            <div v-if="selected.budget">
              <dt>Budjet</dt>
              <dd>{{ selected.budget }}</dd>
            </div>
          </dl>

          <p class="detail__msg">{{ selected.message }}</p>

          <div class="a-field">
            <label class="a-label" for="note">Ichki izoh <small>faqat adminlar ko‘radi</small></label>
            <textarea id="note" v-model="note" class="a-textarea" placeholder="Qo‘ng‘iroq qilindi, taklif yuborildi…" />
          </div>

          <div class="detail__status">
            <button
              type="button"
              class="a-btn a-btn--sm"
              :class="selected.status === 'in_progress' ? 'a-btn--primary' : 'a-btn--outline'"
              :disabled="saving"
              @click="update('in_progress')"
            >
              Jarayonda
            </button>
            <button
              type="button"
              class="a-btn a-btn--sm"
              :class="selected.status === 'done' ? 'a-btn--primary' : 'a-btn--outline'"
              :disabled="saving"
              @click="update('done')"
            >
              Yakunlandi
            </button>
            <button
              type="button"
              class="a-btn a-btn--sm"
              :class="selected.status === 'new' ? 'a-btn--primary' : 'a-btn--outline'"
              :disabled="saving"
              @click="update('new')"
            >
              Yangi
            </button>
            <button type="button" class="a-btn a-btn--sm a-btn--ghost" :disabled="saving" @click="update('spam')">Spam</button>
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
            <button type="button" class="a-btn a-btn--danger a-btn--sm" @click="remove(selected)">
              <AppIcon name="trash" :size="16" /> O‘chirish
            </button>
          </div>
          <p class="a-hint detail__ua">IP: {{ selected.ip || '—' }}</p>
        </template>
        <div v-else class="detail__empty">
          <AppIcon name="inbox" :size="28" />
          <p>So‘rovni tanlang — to‘liq matn, aloqa ma’lumotlari va holatni boshqarish shu yerda.</p>
        </div>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 420px;
  gap: 20px;
  align-items: start;
}

.list {
  display: grid;
  gap: 8px;
}

.req {
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

.req:hover {
  border-color: var(--line-strong);
}

.req.is-active {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-soft);
}

.req.is-new strong::before {
  content: '';
  display: inline-block;
  width: 8px;
  height: 8px;
  margin-right: 8px;
  border-radius: 50%;
  background: var(--accent);
  vertical-align: middle;
}

.req__top {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.req__top small {
  color: var(--ink-3);
  font-size: 12px;
  white-space: nowrap;
}

.req__msg {
  color: var(--ink-2);
  font-size: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.req__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 4px;
}

.req__tag {
  font-size: 12px;
  font-weight: 600;
  padding: 3px 10px;
  border-radius: 99px;
  background: var(--surface-2);
  color: var(--ink-2);
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
  grid-template-columns: 110px 1fr;
  gap: 8px;
}

.detail__meta dt {
  color: var(--ink-3);
}

.detail__meta dd {
  font-weight: 600;
}

.link {
  color: var(--accent);
  margin-left: 8px;
  font-weight: 700;
}

.detail__msg {
  padding: 16px;
  border-radius: 14px;
  background: var(--surface-2);
  white-space: pre-wrap;
  line-height: 1.6;
}

.detail__status,
.detail__foot {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.detail__foot {
  justify-content: space-between;
  padding-top: 16px;
  border-top: 1px solid var(--line);
}

.detail__ua {
  font-size: 12px;
}

.detail__empty {
  display: grid;
  justify-items: center;
  gap: 12px;
  padding: 48px 12px;
  text-align: center;
  color: var(--ink-3);
}

@media (max-width: 1100px) {
  .layout {
    grid-template-columns: 1fr;
  }

  .detail {
    position: static;
  }
}
</style>
