<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import PaginationBar from '@/admin/components/PaginationBar.vue'
import { useAsync } from '@/composables/useAsync'
import { useUiStore } from '@/admin/stores/ui'
import { adminApi } from '@/admin/services/adminApi'
import { toApiError } from '@/services/http'
import { formatBytes, formatDateTime } from '@/utils/format'
import type { Media } from '@/types/api'

const ui = useUiStore()
const kind = ref('')
const q = ref('')
const page = ref(1)
const { data, loading, error, reload } = useAsync(
  () => adminApi.media.list({ kind: kind.value, q: q.value, page: page.value, limit: 40 }),
  {
    watch: [kind, page],
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

// ─── Yuklash ────────────────────────────────────────
const uploading = ref(false)
const progress = ref(0)
const dragOver = ref(false)
const fileInput = ref<HTMLInputElement>()

async function upload(files: FileList | File[] | null | undefined) {
  if (!files?.length) return
  uploading.value = true
  progress.value = 0
  try {
    const res = await adminApi.media.upload([...files], (p) => (progress.value = p))
    if (res.data.length) ui.success(`${res.data.length} ta fayl yuklandi va optimallashtirildi`)
    res.failed.forEach((f) => ui.error(`${f.name}: ${f.message}`))
    page.value = 1
    reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  } finally {
    uploading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

function onDrop(e: DragEvent) {
  dragOver.value = false
  upload(e.dataTransfer?.files)
}

// ─── Tafsilotlar paneli ─────────────────────────────
const selected = ref<Media | null>(null)
const alt = ref('')
const usage = ref<number | null>(null)
const savingAlt = ref(false)

watch(selected, async (m) => {
  alt.value = m?.alt ?? ''
  usage.value = null
  if (m) {
    try {
      usage.value = (await adminApi.media.usage(m.id)).count
    } catch {
      usage.value = null
    }
  }
})

async function saveAlt() {
  if (!selected.value) return
  savingAlt.value = true
  try {
    const updated = await adminApi.media.update(selected.value.id, alt.value)
    Object.assign(selected.value, updated)
    ui.success('Alt matn saqlandi')
  } catch (e) {
    ui.error(toApiError(e).message)
  } finally {
    savingAlt.value = false
  }
}

async function copyUrl(m: Media) {
  try {
    await navigator.clipboard.writeText(new URL(m.url, window.location.origin).href)
    ui.success('Havola nusxalandi')
  } catch {
    ui.error('Nusxalab bo‘lmadi')
  }
}

async function remove(m: Media) {
  const ok = await ui.confirm({
    title: 'Fayl o‘chirilsinmi?',
    text: usage.value
      ? `Bu fayl ${usage.value} joyda ishlatilmoqda. O‘chirilsa, u yerlarda rasm ko‘rinmay qoladi.`
      : 'Fayl va uning barcha o‘lchamdagi nusxalari o‘chiriladi.',
  })
  if (!ok) return
  try {
    await adminApi.media.remove(m.id)
    selected.value = null
    ui.success('O‘chirildi')
    reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}

const kinds = [
  { v: '', l: 'Barchasi' },
  { v: 'image', l: 'Rasmlar' },
  { v: 'vector', l: 'SVG' },
  { v: 'video', l: 'Video' },
]
const items = computed(() => data.value?.items ?? [])
</script>

<template>
  <div class="media-page" @dragover.prevent="dragOver = true" @dragleave.self="dragOver = false" @drop.prevent="onDrop">
    <PageHeader title="Media" subtitle="PNG, JPG, WEBP, AVIF, SVG va MP4. Rasmlar uchun responsive o‘lchamlar avtomatik yaratiladi.">
      <button type="button" class="a-btn a-btn--primary" :disabled="uploading" @click="fileInput?.click()">
        <AppIcon name="upload" :size="18" /> {{ uploading ? `Yuklanmoqda ${progress}%` : 'Fayl yuklash' }}
      </button>
      <input
        ref="fileInput"
        type="file"
        multiple
        accept="image/png,image/jpeg,image/webp,image/avif,image/svg+xml,video/mp4"
        hidden
        @change="upload(($event.target as HTMLInputElement).files)"
      />
    </PageHeader>

    <div v-if="uploading" class="progress"><span :style="{ width: `${progress}%` }" /></div>

    <div class="a-toolbar">
      <div class="a-tabs">
        <button
          v-for="k in kinds"
          :key="k.v"
          type="button"
          class="a-tab"
          :class="{ 'is-active': kind === k.v }"
          @click="((kind = k.v), (page = 1))"
        >
          {{ k.l }}
        </button>
      </div>
      <div class="a-search">
        <AppIcon name="search" :size="18" />
        <input v-model="q" class="a-input" type="search" placeholder="Fayl nomi yoki alt matn" aria-label="Qidirish" />
      </div>
    </div>

    <div class="layout" :class="{ 'has-panel': selected }">
      <div>
        <AdminEmpty v-if="error" error title="Yuklab bo‘lmadi" :text="error.message" @retry="reload" />
        <div v-else-if="loading && !data" class="grid">
          <div v-for="n in 12" :key="n" class="skeleton tile-sk" />
        </div>
        <div v-else-if="!items.length" class="drop" :class="{ 'is-over': dragOver }" @click="fileInput?.click()">
          <AppIcon name="upload" :size="32" />
          <strong>{{ q || kind ? 'Hech narsa topilmadi' : 'Fayllarni shu yerga sudrab tashlang' }}</strong>
          <span>yoki bosib tanlang · bir martada 20 tagacha fayl</span>
        </div>
        <template v-else>
          <div class="grid">
            <button
              v-for="m in items"
              :key="m.id"
              type="button"
              class="tile"
              :class="{ 'is-selected': selected?.id === m.id }"
              @click="selected = selected?.id === m.id ? null : m"
            >
              <video v-if="m.kind === 'video'" :src="m.url" muted preload="metadata" />
              <img v-else :src="m.variants[0]?.url ?? m.url" :alt="m.alt" loading="lazy" />
              <span class="tile__name">{{ m.original_name }}</span>
            </button>
          </div>
          <PaginationBar v-if="data" v-model="page" :meta="data.meta" />
        </template>
      </div>

      <aside v-if="selected" class="a-card panel">
        <div class="panel__head">
          <h2 class="a-card__title" style="margin: 0">Fayl</h2>
          <button type="button" class="a-icon-btn" aria-label="Yopish" @click="selected = null"><AppIcon name="close" /></button>
        </div>
        <div class="panel__preview">
          <video v-if="selected.kind === 'video'" :src="selected.url" controls />
          <img v-else :src="selected.variants.at(-1)?.url ?? selected.url" :alt="selected.alt" />
        </div>
        <dl class="panel__meta">
          <div>
            <dt>Nomi</dt>
            <dd>{{ selected.original_name }}</dd>
          </div>
          <div>
            <dt>Turi</dt>
            <dd>{{ selected.mime }}</dd>
          </div>
          <div>
            <dt>Hajmi</dt>
            <dd>{{ formatBytes(selected.size) }}</dd>
          </div>
          <div v-if="selected.width">
            <dt>O‘lcham</dt>
            <dd>{{ selected.width }} × {{ selected.height }}</dd>
          </div>
          <div v-if="selected.variants.length">
            <dt>Variantlar</dt>
            <dd>{{ selected.variants.map((v) => v.width).join(', ') }}px</dd>
          </div>
          <div>
            <dt>Yuklangan</dt>
            <dd>{{ formatDateTime(selected.created_at) }}</dd>
          </div>
          <div>
            <dt>Ishlatilgan</dt>
            <dd>{{ usage === null ? '…' : usage ? `${usage} joyda` : 'Ishlatilmagan' }}</dd>
          </div>
        </dl>
        <div class="a-field">
          <label class="a-label" for="alt">Alt matn <small>SEO va accessibility</small></label>
          <textarea id="alt" v-model="alt" class="a-textarea" style="min-height: 80px" />
          <button type="button" class="a-btn a-btn--outline a-btn--sm" :disabled="savingAlt || alt === selected.alt" @click="saveAlt">
            Saqlash
          </button>
        </div>
        <div class="panel__actions">
          <button type="button" class="a-btn a-btn--ghost a-btn--sm" @click="copyUrl(selected)">
            <AppIcon name="copy" :size="16" /> Havola
          </button>
          <a :href="selected.url" target="_blank" class="a-btn a-btn--ghost a-btn--sm"><AppIcon name="external" :size="16" /> Ochish</a>
          <button type="button" class="a-btn a-btn--danger a-btn--sm" @click="remove(selected)">
            <AppIcon name="trash" :size="16" /> O‘chirish
          </button>
        </div>
      </aside>
    </div>

    <div v-if="dragOver" class="overlay-drop"><AppIcon name="upload" :size="36" /> Fayllarni qo‘yib yuboring</div>
  </div>
</template>

<style scoped>
.progress {
  height: 4px;
  border-radius: 4px;
  background: var(--surface-2);
  margin-bottom: 16px;
  overflow: hidden;
}

.progress span {
  display: block;
  height: 100%;
  background: var(--accent);
  transition: width 200ms linear;
}

.layout {
  display: grid;
  gap: 20px;
  align-items: start;
}

.layout.has-panel {
  grid-template-columns: minmax(0, 1fr) 340px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(170px, 1fr));
  gap: 12px;
}

.tile-sk {
  aspect-ratio: 1;
  border-radius: 14px;
}

.tile {
  position: relative;
  aspect-ratio: 1;
  border-radius: 14px;
  overflow: hidden;
  background: var(--surface-2);
  outline: 2px solid transparent;
  outline-offset: 2px;
}

.tile img,
.tile video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.tile:hover {
  outline-color: var(--line-strong);
}

.tile.is-selected {
  outline-color: var(--accent);
}

.tile__name {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 18px 10px 8px;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
  text-align: left;
  background: linear-gradient(transparent, rgb(0 0 0 / 0.6));
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  opacity: 0;
  transition: opacity var(--dur-fast) var(--ease);
}

.tile:hover .tile__name,
.tile.is-selected .tile__name {
  opacity: 1;
}

.drop {
  display: grid;
  justify-items: center;
  gap: 8px;
  padding: 80px 24px;
  border-radius: var(--a-radius-lg);
  border: 2px dashed var(--line-strong);
  color: var(--ink-2);
  cursor: pointer;
  text-align: center;
}

.drop strong {
  font-size: 17px;
  color: var(--ink);
}

.drop:hover,
.drop.is-over {
  border-color: var(--accent);
  background: var(--accent-soft);
}

.panel {
  position: sticky;
  top: calc(var(--a-top-h) + 20px);
  display: grid;
  gap: 16px;
}

.panel__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel__preview {
  border-radius: 14px;
  overflow: hidden;
  background: repeating-conic-gradient(var(--surface-2) 0% 25%, var(--surface) 0% 50%) 50% / 20px 20px;
}

.panel__preview img,
.panel__preview video {
  width: 100%;
  max-height: 280px;
  object-fit: contain;
}

.panel__meta {
  display: grid;
  gap: 6px;
  font-size: 13px;
}

.panel__meta div {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.panel__meta dt {
  color: var(--ink-3);
}

.panel__meta dd {
  font-weight: 600;
  text-align: right;
  word-break: break-all;
}

.panel__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.overlay-drop {
  position: fixed;
  inset: 16px;
  z-index: 70;
  display: grid;
  place-content: center;
  justify-items: center;
  gap: 12px;
  border: 2px dashed var(--accent);
  border-radius: 24px;
  background: color-mix(in srgb, var(--bg) 85%, transparent);
  color: var(--accent);
  font-size: 18px;
  font-weight: 700;
  pointer-events: none;
}

@media (max-width: 1100px) {
  .layout.has-panel {
    grid-template-columns: 1fr;
  }

  .panel {
    position: static;
  }
}
</style>
