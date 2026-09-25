<script setup lang="ts">
import MediaImage from '@/components/media/MediaImage.vue'
import { COVER_FOCUS, COVER_RATIOS, focusCss, ratioCss } from '@/utils/cover'
import type { CoverFocus, CoverRatio, MediaRef } from '@/types/api'

// Muqova rasmi kartada qanday o'lchamda chiqishi: nisbat + kesish fokusi, yonida jonli ko'rinish
const props = defineProps<{ media: MediaRef | null }>()
const ratio = defineModel<CoverRatio>('ratio', { required: true })
const focus = defineModel<CoverFocus>('focus', { required: true })

// Tugmadagi kichik to'rtburchak — nisbatning vizual belgisi
function swatch(r: CoverRatio) {
  if (r === 'auto') return { width: '18px', height: '14px' }
  const [w, h] = r.split(':').map(Number)
  const k = 18 / Math.max(w, h)
  return { width: `${Math.round(w * k)}px`, height: `${Math.round(h * k)}px` }
}
</script>

<template>
  <div class="size">
    <p class="a-label">Kartadagi o‘lchami</p>
    <div class="size__ratios" role="radiogroup" aria-label="Muqova nisbati">
      <button
        v-for="r in COVER_RATIOS"
        :key="r.value"
        type="button"
        role="radio"
        class="size__ratio"
        :class="{ 'is-active': ratio === r.value }"
        :aria-checked="ratio === r.value"
        @click="ratio = r.value"
      >
        <i :style="swatch(r.value)" :class="{ 'is-auto': r.value === 'auto' }" />
        {{ r.label }}
      </button>
    </div>

    <p class="a-label">Kesishda saqlanadigan qism</p>
    <div class="size__focus" role="radiogroup" aria-label="Fokus">
      <button
        v-for="f in COVER_FOCUS"
        :key="f.value"
        type="button"
        role="radio"
        class="size__focus-btn"
        :class="{ 'is-active': focus === f.value }"
        :aria-checked="focus === f.value"
        :disabled="ratio === 'auto'"
        @click="focus = f.value"
      >
        {{ f.label }}
      </button>
    </div>

    <div v-if="props.media" class="size__preview">
      <p class="a-hint">Kartadagi ko‘rinishi</p>
      <div class="size__frame" :style="{ aspectRatio: ratioCss(ratio, props.media, '16 / 11') }">
        <MediaImage :media="props.media" sizes="320px" :position="focusCss(focus)" fill />
      </div>
    </div>
    <p v-else class="a-hint">Rasm tanlang — shu yerda kartadagi ko‘rinishi chiqadi.</p>
  </div>
</template>

<style scoped>
.size {
  display: grid;
  gap: 10px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--line);
}

.size__ratios {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 6px;
  margin-bottom: 6px;
}

.size__ratio {
  display: flex;
  align-items: center;
  gap: 8px;
  height: 40px;
  padding: 0 10px;
  border-radius: 10px;
  border: 1px solid var(--line-strong);
  font-size: 13px;
  font-weight: 700;
  transition:
    border-color 0.15s ease,
    background-color 0.15s ease;
}

.size__ratio i {
  flex-shrink: 0;
  border-radius: 3px;
  border: 1.5px solid currentColor;
  opacity: 0.7;
}

.size__ratio i.is-auto {
  border-style: dashed;
}

.size__ratio:hover {
  border-color: var(--ink-3);
}

.size__ratio.is-active {
  border-color: var(--accent);
  background: var(--accent-soft);
  color: var(--accent);
}

.size__focus {
  display: flex;
  gap: 4px;
  padding: 4px;
  border-radius: 12px;
  background: var(--surface-2);
}

.size__focus-btn {
  flex: 1;
  height: 34px;
  border-radius: 9px;
  font-size: 13px;
  font-weight: 700;
  color: var(--ink-2);
}

.size__focus-btn.is-active {
  background: var(--surface);
  color: var(--ink);
  box-shadow: var(--shadow-sm);
}

.size__focus-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.size__preview {
  display: grid;
  gap: 8px;
  margin-top: 6px;
}

.size__frame {
  position: relative;
  max-height: 360px;
  border-radius: 14px;
  overflow: hidden;
  background: var(--surface-2);
}
</style>
