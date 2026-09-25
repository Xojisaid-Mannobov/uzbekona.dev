<script setup lang="ts">
import { ref } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'
import MediaField from '@/admin/components/MediaField.vue'
import MediaPicker from '@/admin/components/MediaPicker.vue'
import MetricsEditor from '@/admin/components/MetricsEditor.vue'
import TagInput from '@/admin/components/TagInput.vue'
import type { Block, Media } from '@/types/api'
import { COVER_RATIOS } from '@/utils/cover'

// Har bir blok turi uchun tahrirlash maydonlari. Blok obyekti joyida (in-place) o'zgartiriladi.
const props = defineProps<{ block: Block }>()

const galleryOpen = ref(false)

// Matn ichidagi rasm: kichik (markazda), o'rta (matn kengligida), keng (matndan chiqib turadi)
const IMAGE_SIZES = [
  { value: 'small', label: 'Kichik' },
  { value: 'medium', label: 'O‘rta' },
  { value: 'wide', label: 'Keng' },
] as const

function addGallery(items: Media[]) {
  if (props.block.type !== 'gallery') return
  const refs = items.map((m) => ({
    media: { id: m.id, url: m.url, kind: m.kind, mime: m.mime, width: m.width, height: m.height, alt: m.alt, variants: m.variants },
    caption: '',
  }))
  props.block.data.items = [...props.block.data.items, ...refs]
}

function moveItem<T>(arr: T[], i: number, dir: number) {
  const j = i + dir
  if (j < 0 || j >= arr.length) return
  ;[arr[i], arr[j]] = [arr[j], arr[i]]
}
</script>

<template>
  <div class="be">
    <!-- Sarlavha -->
    <template v-if="block.type === 'heading'">
      <div class="be__row">
        <label class="a-field be__narrow">
          <span class="a-label">Label <small>ixtiyoriy</small></span>
          <input v-model="block.data.label" class="a-input" placeholder="Jarayon" />
        </label>
        <label class="a-field">
          <span class="a-label">Sarlavha</span>
          <input v-model="block.data.text" class="a-input be__heading" placeholder="Qanday ishladik" />
        </label>
      </div>
    </template>

    <!-- Matn / katta matn -->
    <label v-else-if="block.type === 'text' || block.type === 'large_text'" class="a-field">
      <span class="a-label">Matn <small>bo‘sh qator — yangi paragraf</small></span>
      <textarea v-model="block.data.text" class="a-textarea a-textarea--lg" :class="{ be__large: block.type === 'large_text' }" />
    </label>

    <!-- Rasm -->
    <template v-else-if="block.type === 'image' || block.type === 'full_image'">
      <MediaField v-model="block.data.media" ratio="16 / 8" />
      <div v-if="block.type === 'image'" class="be__row be__imgopts">
        <div class="a-field">
          <span class="a-label">O‘lchami</span>
          <div class="be__seg" role="radiogroup" aria-label="Rasm o‘lchami">
            <button
              v-for="o in IMAGE_SIZES"
              :key="o.value"
              type="button"
              role="radio"
              :aria-checked="(block.data.size ?? 'medium') === o.value"
              :class="{ 'is-active': (block.data.size ?? 'medium') === o.value }"
              @click="block.data.size = o.value"
            >
              {{ o.label }}
            </button>
          </div>
        </div>
        <div class="a-field">
          <span class="a-label">Nisbati</span>
          <div class="be__seg" role="radiogroup" aria-label="Rasm nisbati">
            <button
              v-for="r in COVER_RATIOS"
              :key="r.value"
              type="button"
              role="radio"
              :aria-checked="(block.data.ratio ?? 'auto') === r.value"
              :class="{ 'is-active': (block.data.ratio ?? 'auto') === r.value }"
              @click="block.data.ratio = r.value"
            >
              {{ r.label }}
            </button>
          </div>
        </div>
      </div>
      <label class="a-field">
        <span class="a-label">Izoh <small>ixtiyoriy</small></span>
        <input v-model="block.data.caption" class="a-input" placeholder="Talaba kabineti — asosiy sahifa" />
      </label>
    </template>

    <!-- Galereya -->
    <template v-else-if="block.type === 'gallery'">
      <div class="be__gallery">
        <div v-for="(item, i) in block.data.items" :key="item.media.id + '-' + i" class="be__gitem">
          <img :src="item.media.variants[0]?.url ?? item.media.url" :alt="item.media.alt" />
          <input v-model="item.caption" class="a-input" placeholder="Izoh" />
          <div class="be__gactions">
            <button type="button" class="a-icon-btn" aria-label="Chapga" :disabled="i === 0" @click="moveItem(block.data.items, i, -1)">
              <AppIcon name="arrow-left" :size="16" />
            </button>
            <button
              type="button"
              class="a-icon-btn"
              aria-label="O‘ngga"
              :disabled="i === block.data.items.length - 1"
              @click="moveItem(block.data.items, i, 1)"
            >
              <AppIcon name="arrow-right" :size="16" />
            </button>
            <button type="button" class="a-icon-btn a-icon-btn--danger" aria-label="O‘chirish" @click="block.data.items.splice(i, 1)">
              <AppIcon name="trash" :size="16" />
            </button>
          </div>
        </div>
        <button type="button" class="be__gadd" @click="galleryOpen = true"><AppIcon name="plus" :size="24" /> Rasm qo‘shish</button>
      </div>
      <MediaPicker v-model:open="galleryOpen" multiple @select="addGallery" />
    </template>

    <!-- Video -->
    <template v-else-if="block.type === 'video'">
      <label class="a-field">
        <span class="a-label">YouTube havola <small>yoki pastda MP4 tanlang</small></span>
        <input v-model="block.data.url" class="a-input" placeholder="https://youtube.com/watch?v=…" />
      </label>
      <div class="be__row">
        <div class="a-field">
          <span class="a-label">MP4 video</span>
          <MediaField v-model="block.data.media" kind="video" ratio="16 / 9" />
        </div>
        <div class="a-field">
          <span class="a-label">Poster rasm <small>ixtiyoriy</small></span>
          <MediaField v-model="block.data.poster" ratio="16 / 9" />
        </div>
      </div>
      <label class="a-field">
        <span class="a-label">Izoh</span>
        <input v-model="block.data.caption" class="a-input" />
      </label>
    </template>

    <!-- Raqamlar -->
    <MetricsEditor v-else-if="block.type === 'stats'" v-model="block.data.items" :max="6" />

    <!-- Iqtibos -->
    <template v-else-if="block.type === 'quote'">
      <label class="a-field">
        <span class="a-label">Iqtibos matni</span>
        <textarea v-model="block.data.text" class="a-textarea" />
      </label>
      <div class="be__row">
        <label class="a-field">
          <span class="a-label">Muallif</span>
          <input v-model="block.data.author" class="a-input" placeholder="Khuja Mannopov" />
        </label>
        <label class="a-field">
          <span class="a-label">Lavozim</span>
          <input v-model="block.data.role" class="a-input" placeholder="Rektor" />
        </label>
      </div>
    </template>

    <!-- Ustunlar -->
    <div
      v-else-if="block.type === 'two_columns' || block.type === 'three_columns'"
      class="be__cols"
      :class="`be__cols--${block.data.columns.length}`"
    >
      <div v-for="(col, i) in block.data.columns" :key="i" class="be__col">
        <input v-model="col.title" class="a-input" :placeholder="`Ustun ${i + 1} sarlavhasi`" />
        <textarea v-model="col.text" class="a-textarea" placeholder="Matn" />
      </div>
    </div>

    <!-- Texnologiyalar -->
    <template v-else-if="block.type === 'technology'">
      <label class="a-field">
        <span class="a-label">Sarlavha</span>
        <input v-model="block.data.title" class="a-input" />
      </label>
      <div class="a-field">
        <span class="a-label">Texnologiyalar</span>
        <TagInput
          v-model="block.data.items"
          placeholder="Vue.js, Go, PostgreSQL…"
          :suggestions="['Vue.js', 'TypeScript', 'Go', 'PostgreSQL', 'Redis', 'Docker', 'Nginx']"
        />
      </div>
    </template>

    <!-- Jarayon -->
    <div v-else-if="block.type === 'process'" class="be__steps">
      <div v-for="(step, i) in block.data.steps" :key="i" class="be__step">
        <span class="be__num">{{ String(i + 1).padStart(2, '0') }}</span>
        <div class="be__step-fields">
          <input v-model="step.title" class="a-input" placeholder="Bosqich nomi" />
          <textarea v-model="step.text" class="a-textarea" placeholder="Qisqa tavsif" />
        </div>
        <button type="button" class="a-icon-btn a-icon-btn--danger" aria-label="O‘chirish" @click="block.data.steps.splice(i, 1)">
          <AppIcon name="trash" :size="16" />
        </button>
      </div>
      <button
        type="button"
        class="a-btn a-btn--outline a-btn--sm"
        style="justify-self: start"
        @click="block.data.steps.push({ title: '', text: '' })"
      >
        <AppIcon name="plus" :size="16" /> Bosqich qo‘shish
      </button>
    </div>

    <!-- Oldin / keyin -->
    <template v-else-if="block.type === 'before_after'">
      <div class="be__row">
        <div class="a-field">
          <span class="a-label">Oldin</span>
          <MediaField v-model="block.data.before" ratio="16 / 10" />
        </div>
        <div class="a-field">
          <span class="a-label">Keyin</span>
          <MediaField v-model="block.data.after" ratio="16 / 10" />
        </div>
      </div>
      <label class="a-field">
        <span class="a-label">Izoh</span>
        <input v-model="block.data.caption" class="a-input" />
      </label>
    </template>
  </div>
</template>

<style scoped>
.be__imgopts {
  flex-wrap: wrap;
  gap: 16px;
}

.be__seg {
  display: inline-flex;
  flex-wrap: wrap;
  gap: 4px;
  padding: 4px;
  border-radius: 12px;
  background: var(--surface-2);
}

.be__seg button {
  height: 32px;
  padding: 0 12px;
  border-radius: 9px;
  font-size: 13px;
  font-weight: 700;
  color: var(--ink-2);
}

.be__seg button.is-active {
  background: var(--surface);
  color: var(--ink);
  box-shadow: var(--shadow-sm);
}

.be {
  display: grid;
  gap: 16px;
}

.be__row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.be__row:has(.be__narrow) {
  grid-template-columns: 200px 1fr;
}

.be__heading {
  font-size: 18px;
  font-weight: 700;
}

.be__large {
  font-size: 20px;
  font-weight: 500;
  letter-spacing: -0.01em;
}

.be__cols {
  display: grid;
  gap: 12px;
}

.be__cols--2 {
  grid-template-columns: 1fr 1fr;
}

.be__cols--3 {
  grid-template-columns: repeat(3, 1fr);
}

.be__col {
  display: grid;
  gap: 8px;
  align-content: start;
}

.be__gallery {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 12px;
}

.be__gitem {
  display: grid;
  gap: 6px;
}

.be__gitem img {
  aspect-ratio: 4 / 3;
  width: 100%;
  object-fit: cover;
  border-radius: 12px;
  background: var(--surface-2);
}

.be__gitem .a-input {
  height: 38px;
  font-size: 13px;
}

.be__gactions {
  display: flex;
  justify-content: center;
}

.be__gadd {
  aspect-ratio: 4 / 3;
  display: grid;
  place-content: center;
  justify-items: center;
  gap: 6px;
  border-radius: 12px;
  border: 1.5px dashed var(--line-strong);
  color: var(--ink-2);
  font-weight: 600;
  font-size: 13px;
}

.be__gadd:hover {
  border-color: var(--accent);
  color: var(--accent);
}

.be__steps {
  display: grid;
  gap: 12px;
}

.be__step {
  display: grid;
  grid-template-columns: 44px 1fr auto;
  gap: 12px;
  align-items: start;
}

.be__num {
  font-size: 24px;
  font-weight: 700;
  color: var(--accent);
  letter-spacing: -0.04em;
  padding-top: 8px;
}

.be__step-fields {
  display: grid;
  gap: 8px;
}

.be__step-fields .a-textarea {
  min-height: 72px;
}

@media (max-width: 720px) {
  .be__row,
  .be__row:has(.be__narrow),
  .be__cols--2,
  .be__cols--3 {
    grid-template-columns: 1fr;
  }
}
</style>
