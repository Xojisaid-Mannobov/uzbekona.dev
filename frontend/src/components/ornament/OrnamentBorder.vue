<script setup lang="ts">
import { computed, useId } from 'vue'
import { BORDER_TILE } from './geometry'

/** Hoshiya — an'anaviy chegaraviy naqsh bandi (bo'limlar orasidagi ajratkich). */
const props = withDefaults(defineProps<{ tone?: 'accent' | 'gold' | 'ink' | 'light'; opacity?: number }>(), {
  tone: 'gold',
  opacity: 0.55,
})

const id = `hoshiya-${useId()}`
const color = computed(() => ({ accent: 'var(--ornament)', gold: 'var(--gold)', ink: 'var(--ink)', light: 'var(--on-dark)' })[props.tone])
</script>

<template>
  <svg class="hoshiya" :style="{ color, opacity }" height="28" aria-hidden="true" focusable="false">
    <defs>
      <pattern :id="id" width="56" height="28" patternUnits="userSpaceOnUse">
        <g fill="none" stroke="currentColor" stroke-width="1.2" stroke-linejoin="round">
          <path v-for="(d, i) in BORDER_TILE" :key="i" :d="d" />
        </g>
      </pattern>
    </defs>
    <rect width="100%" height="28" :fill="`url(#${id})`" />
  </svg>
</template>

<style scoped>
.hoshiya {
  display: block;
  width: 100%;
  height: 28px;
  pointer-events: none;
}
</style>
