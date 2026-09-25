<script setup lang="ts">
import { ref } from 'vue'
import MediaImage from '@/components/media/MediaImage.vue'
import type { MediaRef } from '@/types/api'
import { aspect } from '@/utils/media'

// Oldin/keyin taqqoslash: slayder (klaviatura bilan ham boshqariladi)
defineProps<{ before: MediaRef; after: MediaRef }>()
const pos = ref(50)
</script>

<template>
  <div class="ba" :style="{ aspectRatio: aspect(after) }">
    <MediaImage :media="before" fill sizes="(min-width: 1440px) 1360px, 100vw" />
    <div class="ba__after" :style="{ clipPath: `inset(0 0 0 ${pos}%)` }">
      <MediaImage :media="after" fill sizes="(min-width: 1440px) 1360px, 100vw" />
    </div>
    <div class="ba__handle" :style="{ left: `${pos}%` }" aria-hidden="true"><span>⟷</span></div>
    <span class="ba__tag ba__tag--l">Oldin</span>
    <span class="ba__tag ba__tag--r">Keyin</span>
    <input v-model.number="pos" class="ba__range" type="range" min="0" max="100" step="0.5" aria-label="Oldin va keyin taqqoslash" />
  </div>
</template>

<style scoped>
.ba {
  position: relative;
  border-radius: var(--r-lg);
  overflow: hidden;
  background: var(--surface-2);
}

.ba__after {
  position: absolute;
  inset: 0;
}

.ba__handle {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #fff;
  translate: -1px 0;
  pointer-events: none;
}

.ba__handle span {
  position: absolute;
  top: 50%;
  left: 50%;
  translate: -50% -50%;
  display: grid;
  place-items: center;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: #fff;
  color: #111312;
  font-size: 22px;
  box-shadow: var(--shadow-md);
}

.ba__tag {
  position: absolute;
  top: 20px;
  padding: 8px 16px;
  border-radius: var(--r-pill);
  background: rgb(17 19 18 / 0.7);
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  pointer-events: none;
}

.ba__tag--l {
  left: 20px;
}

.ba__tag--r {
  right: 20px;
}

.ba__range {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: ew-resize;
  margin: 0;
}

.ba__range:focus-visible + .ba__handle,
.ba:focus-within .ba__handle span {
  outline: 2px solid var(--accent);
}
</style>
