<script setup lang="ts">
import { STAR_8 } from './geometry'

// Kichik 8 qirrali yulduz — bo'lim label'lari oldidagi milliy belgi
withDefaults(defineProps<{ size?: number; filled?: boolean }>(), { size: 14, filled: false })
</script>

<template>
  <svg class="star8" :width="size" :height="size" viewBox="0 0 24 24" aria-hidden="true" focusable="false">
    <path :d="STAR_8" :fill="filled ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="1.6" stroke-linejoin="round" />
    <circle v-if="!filled" cx="12" cy="12" r="2.4" fill="currentColor" />
  </svg>
</template>

<style scoped>
.star8 {
  display: inline-block;
  flex-shrink: 0;
  color: var(--star-color, var(--gold));
}
</style>
