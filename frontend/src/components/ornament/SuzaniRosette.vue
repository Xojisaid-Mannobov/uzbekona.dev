<script setup lang="ts">
import { computed } from 'vue'
import { ROSETTE_DOTS, ROSETTE_PATHS } from './geometry'

/** Suzani rozetkasi — kashtachilikdagi markaziy gul (chiziqli, katta dekorativ element). */
const props = withDefaults(
  defineProps<{ tone?: 'accent' | 'gold' | 'ink' | 'light'; opacity?: number; spin?: boolean; strokeWidth?: number }>(),
  { tone: 'accent', opacity: 0.18, spin: false, strokeWidth: 1.2 },
)

const color = computed(() => ({ accent: 'var(--ornament)', gold: 'var(--gold)', ink: 'var(--ink)', light: 'var(--on-dark)' })[props.tone])
</script>

<template>
  <svg
    class="rosette"
    :class="{ 'rosette--spin': spin }"
    viewBox="-100 -100 200 200"
    :style="{ color, opacity }"
    aria-hidden="true"
    focusable="false"
  >
    <g fill="none" stroke="currentColor" :stroke-width="strokeWidth" stroke-linejoin="round">
      <path v-for="(d, i) in ROSETTE_PATHS" :key="i" :d="d" vector-effect="non-scaling-stroke" />
    </g>
    <circle v-for="([x, y], i) in ROSETTE_DOTS" :key="`d${i}`" :cx="x" :cy="y" r="1.8" fill="currentColor" />
  </svg>
</template>

<style scoped>
.rosette {
  display: block;
  pointer-events: none;
  user-select: none;
}

/* Juda sekin aylanish — tirik, lekin e'tiborni tortmaydigan harakat */
.rosette--spin {
  animation: rosette-spin 140s linear infinite;
}

@keyframes rosette-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
