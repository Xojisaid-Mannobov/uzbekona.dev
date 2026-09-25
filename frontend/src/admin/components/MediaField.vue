<script setup lang="ts">
import { ref } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'
import MediaPicker from './MediaPicker.vue'
import type { Media, MediaRef } from '@/types/api'

/** Bitta media tanlash maydoni (cover, rasm, video). Qiymat — media snapshot. */
const model = defineModel<MediaRef | null>({ required: true })
withDefaults(defineProps<{ kind?: 'image' | 'video'; ratio?: string; label?: string }>(), { kind: 'image', ratio: '16 / 10' })

const open = ref(false)

// Media obyektidan faqat kerakli maydonlar olinadi (bloklar JSON'ida saqlanadi)
function toRef(m: Media): MediaRef {
  return { id: m.id, url: m.url, kind: m.kind, mime: m.mime, width: m.width, height: m.height, alt: m.alt, variants: m.variants }
}

function onSelect(items: Media[]) {
  if (items[0]) model.value = toRef(items[0])
}
</script>

<template>
  <div class="mf">
    <div v-if="model" class="mf__preview" :style="{ aspectRatio: ratio }">
      <video v-if="model.kind === 'video'" :src="model.url" muted controls preload="metadata" />
      <img v-else :src="model.variants[1]?.url ?? model.variants[0]?.url ?? model.url" :alt="model.alt" />
      <div class="mf__actions">
        <button type="button" class="a-btn a-btn--ghost a-btn--sm" @click="open = true">
          <AppIcon name="refresh" :size="16" /> Almashtirish
        </button>
        <button type="button" class="a-icon-btn a-icon-btn--danger mf__remove" aria-label="Olib tashlash" @click="model = null">
          <AppIcon name="trash" :size="18" />
        </button>
      </div>
    </div>
    <button v-else type="button" class="mf__empty" :style="{ aspectRatio: ratio }" @click="open = true">
      <AppIcon :name="kind === 'video' ? 'video' : 'image'" :size="28" />
      <span>{{ label ?? (kind === 'video' ? 'Video tanlash' : 'Rasm tanlash') }}</span>
      <small>Kutubxonadan yoki yangi yuklash</small>
    </button>
    <MediaPicker v-model:open="open" :kind="kind" @select="onSelect" />
  </div>
</template>

<style scoped>
.mf__preview {
  position: relative;
  border-radius: var(--a-radius);
  overflow: hidden;
  background: var(--surface-2);
  border: 1px solid var(--line);
}

.mf__preview img,
.mf__preview video {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.mf__actions {
  position: absolute;
  left: 10px;
  right: 10px;
  bottom: 10px;
  display: flex;
  justify-content: space-between;
  opacity: 0;
  transition: opacity var(--dur-fast) var(--ease);
}

.mf__preview:hover .mf__actions,
.mf__preview:focus-within .mf__actions {
  opacity: 1;
}

.mf__remove {
  background: var(--surface);
}

.mf__empty {
  width: 100%;
  display: grid;
  place-content: center;
  justify-items: center;
  gap: 6px;
  border-radius: var(--a-radius);
  border: 1.5px dashed var(--line-strong);
  color: var(--ink-2);
  font-weight: 600;
  transition:
    border-color var(--dur-fast) var(--ease),
    background-color var(--dur-fast) var(--ease);
}

.mf__empty small {
  font-weight: 500;
  color: var(--ink-3);
}

.mf__empty:hover {
  border-color: var(--accent);
  background: var(--accent-soft);
  color: var(--accent);
}

@media (hover: none) {
  .mf__actions {
    opacity: 1;
  }
}
</style>
