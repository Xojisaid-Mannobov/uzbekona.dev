<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'
import { adminApi } from '@/admin/services/adminApi'
import { useUiStore } from '@/admin/stores/ui'
import { toApiError } from '@/services/http'
import type { Media } from '@/types/api'

/**
 * Media kutubxonasidan tanlash oynasi. Shu yerning o'zida yangi fayl yuklash mumkin.
 * `multiple` — galereya uchun bir nechta tanlash.
 */
const props = withDefaults(defineProps<{ multiple?: boolean; kind?: 'image' | 'video' | '' }>(), { kind: 'image' })
const open = defineModel<boolean>('open', { required: true })
const emit = defineEmits<{ select: [items: Media[]] }>()

const ui = useUiStore()
const items = ref<Media[]>([])
const page = ref(1)
const total = ref(0)
const loading = ref(false)
const error = ref('')
const q = ref('')
const selected = ref<Media[]>([])
const uploading = ref(false)
const progress = ref(0)
const dragOver = ref(false)
const fileInput = ref<HTMLInputElement>()

const hasMore = computed(() => items.value.length < total.value)
const accept = computed(() =>
  props.kind === 'video'
    ? 'video/mp4'
    : props.kind === 'image'
      ? 'image/png,image/jpeg,image/webp,image/avif,image/svg+xml'
      : 'image/*,video/mp4',
)

async function load(reset = false) {
  if (reset) {
    page.value = 1
    items.value = []
  }
  loading.value = true
  error.value = ''
  try {
    // Rasm tanlashda vektor (SVG) fayllar ham ko'rsatiladi
    const res = await adminApi.media.list({ q: q.value, page: page.value, limit: 30, kind: props.kind === 'video' ? 'video' : undefined })
    const list = props.kind === 'image' ? res.items.filter((m) => m.kind !== 'video') : res.items
    items.value = reset ? list : [...items.value, ...list]
    total.value = res.meta.total
  } catch (e) {
    error.value = toApiError(e).message
  } finally {
    loading.value = false
  }
}

watch(open, (v) => {
  if (v) {
    selected.value = []
    load(true)
  }
})

let searchTimer: ReturnType<typeof setTimeout>
watch(q, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => load(true), 300)
})

function toggle(m: Media) {
  const i = selected.value.findIndex((s) => s.id === m.id)
  if (i >= 0) selected.value.splice(i, 1)
  else if (props.multiple) selected.value.push(m)
  else selected.value = [m]
}

function isSelected(m: Media) {
  return selected.value.some((s) => s.id === m.id)
}

async function upload(files: FileList | File[] | null) {
  if (!files?.length) return
  uploading.value = true
  progress.value = 0
  try {
    const res = await adminApi.media.upload([...files], (p) => (progress.value = p))
    items.value = [...res.data, ...items.value]
    total.value += res.data.length
    // Yangi yuklanganlar avtomatik tanlanadi
    selected.value = props.multiple ? [...res.data, ...selected.value] : res.data.slice(0, 1)
    res.failed.forEach((f) => ui.error(`${f.name}: ${f.message}`))
    if (res.data.length) ui.success(`${res.data.length} ta fayl yuklandi`)
  } catch (e) {
    ui.error(toApiError(e).message)
  } finally {
    uploading.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

function onDrop(e: DragEvent) {
  dragOver.value = false
  upload(e.dataTransfer?.files ?? null)
}

function confirm() {
  emit('select', selected.value)
  open.value = false
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="open" class="overlay" @click.self="open = false" @keydown.esc="open = false">
        <div
          class="picker"
          role="dialog"
          aria-modal="true"
          aria-label="Media tanlash"
          :class="{ 'is-drag': dragOver }"
          @dragover.prevent="dragOver = true"
          @dragleave.self="dragOver = false"
          @drop.prevent="onDrop"
        >
          <header class="picker__head">
            <h2>Media kutubxonasi</h2>
            <div class="a-search">
              <AppIcon name="search" :size="18" />
              <input v-model="q" class="a-input" type="search" placeholder="Fayl nomi bo‘yicha qidirish" />
            </div>
            <button type="button" class="a-btn a-btn--primary" :disabled="uploading" @click="fileInput?.click()">
              <AppIcon name="upload" :size="18" /> {{ uploading ? `Yuklanmoqda ${progress}%` : 'Yuklash' }}
            </button>
            <input
              ref="fileInput"
              type="file"
              :accept="accept"
              multiple
              hidden
              @change="upload(($event.target as HTMLInputElement).files)"
            />
            <button type="button" class="a-icon-btn" aria-label="Yopish" @click="open = false"><AppIcon name="close" /></button>
          </header>

          <div class="picker__body">
            <p v-if="error" class="a-error">{{ error }}</p>
            <div v-if="!loading && !items.length && !error" class="a-empty">
              <span class="a-empty__icon"><AppIcon name="image" :size="24" /></span>
              <strong>Hozircha fayl yo‘q</strong>
              <p>Fayllarni shu oynaga sudrab tashlang yoki “Yuklash” tugmasini bosing.</p>
            </div>
            <div class="grid">
              <button
                v-for="m in items"
                :key="m.id"
                type="button"
                class="tile"
                :class="{ 'is-selected': isSelected(m) }"
                :aria-pressed="isSelected(m)"
                :title="m.original_name"
                @click="toggle(m)"
                @dblclick="!multiple && ((selected = [m]), confirm())"
              >
                <video v-if="m.kind === 'video'" :src="m.url" muted preload="metadata" />
                <img v-else :src="m.variants[0]?.url ?? m.url" :alt="m.alt" loading="lazy" />
                <span class="tile__check"><AppIcon name="check" :size="16" /></span>
                <span v-if="m.kind === 'video'" class="tile__kind"><AppIcon name="play" :size="14" /></span>
              </button>
              <div v-for="n in loading ? 6 : 0" :key="`s${n}`" class="skeleton tile" />
            </div>
            <div v-if="hasMore && !loading" class="picker__more">
              <button type="button" class="a-btn a-btn--outline" @click="(page++, load())">Ko‘proq ko‘rsatish</button>
            </div>
          </div>

          <footer class="picker__foot">
            <span class="a-hint">{{ selected.length ? `${selected.length} ta tanlandi` : 'Faylni tanlang' }}</span>
            <div class="picker__btns">
              <button type="button" class="a-btn a-btn--ghost" @click="open = false">Bekor qilish</button>
              <button type="button" class="a-btn a-btn--accent" :disabled="!selected.length" @click="confirm">Tanlash</button>
            </div>
          </footer>

          <div v-if="dragOver" class="picker__drop"><AppIcon name="upload" :size="32" /> Fayllarni qo‘yib yuboring</div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  z-index: 80;
  display: grid;
  place-items: center;
  padding: 24px;
  background: rgb(17 19 18 / 0.45);
  backdrop-filter: blur(4px);
}

.picker {
  position: relative;
  width: min(1080px, 100%);
  height: min(760px, 100%);
  display: grid;
  grid-template-rows: auto 1fr auto;
  border-radius: 24px;
  background: var(--surface);
  box-shadow: var(--shadow-lg);
  overflow: hidden;
  font-size: 15px;
}

.picker__head {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 18px 20px;
  border-bottom: 1px solid var(--line);
}

.picker__head h2 {
  font-size: 18px;
  letter-spacing: -0.02em;
  margin-right: auto;
}

.picker__head .a-search {
  max-width: 300px;
}

.picker__body {
  overflow-y: auto;
  padding: 20px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 12px;
}

.tile {
  position: relative;
  aspect-ratio: 1;
  border-radius: 14px;
  overflow: hidden;
  background: var(--surface-2);
  outline: 2px solid transparent;
  outline-offset: 2px;
  transition: outline-color var(--dur-fast) var(--ease);
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

.tile__check {
  position: absolute;
  top: 8px;
  right: 8px;
  display: grid;
  place-items: center;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: var(--accent);
  color: #fff;
  opacity: 0;
  transform: scale(0.7);
  transition:
    opacity var(--dur-fast) var(--ease),
    transform var(--dur-fast) var(--ease);
}

.is-selected .tile__check {
  opacity: 1;
  transform: none;
}

.tile__kind {
  position: absolute;
  left: 8px;
  bottom: 8px;
  display: grid;
  place-items: center;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: rgb(0 0 0 / 0.6);
  color: #fff;
}

.picker__more {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.picker__foot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-top: 1px solid var(--line);
}

.picker__btns {
  display: flex;
  gap: 10px;
}

.picker__drop {
  position: absolute;
  inset: 12px;
  display: grid;
  place-content: center;
  justify-items: center;
  gap: 10px;
  border: 2px dashed var(--accent);
  border-radius: 20px;
  background: color-mix(in srgb, var(--surface) 88%, transparent);
  color: var(--accent);
  font-weight: 700;
  pointer-events: none;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 200ms var(--ease);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

@media (max-width: 720px) {
  .picker__head {
    flex-wrap: wrap;
  }

  .picker__head .a-search {
    order: 5;
    max-width: none;
    flex-basis: 100%;
  }
}
</style>
