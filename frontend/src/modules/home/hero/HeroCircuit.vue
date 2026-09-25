<script setup lang="ts">
import { computed } from 'vue'

// Texnologiya ramzi — "sxema" chiziqlari, ular bo'ylab yorug' impuls yuguradi
const props = defineProps<{ variant: 'top' | 'bottom' }>()

type Circuit = { box: string; lines: { d: string; soft?: boolean }[]; nodes: [number, number, 'solid' | 'ring'][] }

const CIRCUITS: Record<'top' | 'bottom', Circuit> = {
  top: {
    box: '0 0 640 240',
    lines: [{ d: 'M0 200 H220 L300 120 H470 L560 30 H640' }, { d: 'M60 236 H300 L360 176 H520 L600 96 H640', soft: true }],
    nodes: [
      [300, 120, 'ring'],
      [470, 120, 'solid'],
      [360, 176, 'ring'],
    ],
  },
  bottom: {
    box: '0 0 620 220',
    lines: [{ d: 'M0 20 C60 24 100 50 130 80 L170 120 H400 L500 220' }, { d: 'M0 76 C40 80 70 96 90 116 L124 150 H300', soft: true }],
    nodes: [
      [170, 120, 'solid'],
      [400, 120, 'ring'],
    ],
  },
}

const c = computed(() => CIRCUITS[props.variant])
</script>

<template>
  <svg class="circuit" :viewBox="c.box" aria-hidden="true" focusable="false">
    <g v-for="(l, i) in c.lines" :key="i">
      <path :d="l.d" class="circuit__line" :class="{ 'circuit__line--soft': l.soft }" />
      <path :d="l.d" class="circuit__pulse" pathLength="100" :style="{ animationDelay: `${i * -2.4}s` }" />
    </g>
    <g v-for="([x, y, kind], i) in c.nodes" :key="`n${i}`">
      <circle :cx="x" :cy="y" r="11" class="circuit__halo" />
      <circle :cx="x" :cy="y" :r="kind === 'solid' ? 6 : 5" :class="kind === 'solid' ? 'circuit__node' : 'circuit__ring'" />
    </g>
  </svg>
</template>

<style scoped>
.circuit {
  display: block;
  width: 100%;
  height: auto;
  overflow: visible;
}

.circuit__line {
  fill: none;
  stroke: color-mix(in srgb, var(--ornament) 55%, transparent);
  stroke-width: 1.5;
}

.circuit__line--soft {
  stroke: color-mix(in srgb, #20b9a4 45%, transparent);
}

.circuit__pulse {
  fill: none;
  stroke: var(--ornament);
  stroke-width: 2.5;
  stroke-linecap: round;
  stroke-dasharray: 6 94;
  animation: circuit-pulse 6s linear infinite;
}

.circuit__halo {
  fill: color-mix(in srgb, var(--ornament) 14%, transparent);
}

.circuit__node {
  fill: var(--ornament);
}

.circuit__ring {
  fill: var(--bg);
  stroke: color-mix(in srgb, var(--ornament) 60%, transparent);
  stroke-width: 2;
}

@keyframes circuit-pulse {
  from {
    stroke-dashoffset: 100;
  }
  to {
    stroke-dashoffset: 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  .circuit__pulse {
    display: none;
  }
}
</style>
