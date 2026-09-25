<script setup lang="ts">
import { useId } from 'vue'
import { GIRIH_TILE } from '@/components/ornament/geometry'
import { SKYLINE } from './art'

// Registon uslubidagi me'moriy silueti — bir rangli (accent) tonlar bilan, pastga qarab tumanga singadi
const uid = useId()
const ids = { wall: `sk-wall-${uid}`, tower: `sk-tower-${uid}`, tile: `sk-tile-${uid}` }
</script>

<template>
  <svg class="skyline" viewBox="0 0 900 520" aria-hidden="true" focusable="false">
    <defs>
      <linearGradient :id="ids.wall" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0" class="stop-wall-top" />
        <stop offset="1" class="stop-wall-bottom" />
      </linearGradient>
      <linearGradient :id="ids.tower" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0" class="stop-tower-top" />
        <stop offset="1" class="stop-tower-bottom" />
      </linearGradient>
      <pattern :id="ids.tile" width="30" height="30" patternUnits="userSpaceOnUse" viewBox="0 0 100 100">
        <g class="tile" fill="none" stroke-width="3" stroke-linejoin="round">
          <path v-for="(d, i) in GIRIH_TILE" :key="i" :d="d" />
        </g>
      </pattern>
    </defs>

    <path class="far" :d="SKYLINE.far" />
    <path :d="SKYLINE.wall" :fill="`url(#${ids.wall})`" />
    <!-- peshtoq koshinlari -->
    <path d="M286 520 V120 H458 V520 Z" :fill="`url(#${ids.tile})`" />
    <path :d="SKYLINE.tower" :fill="`url(#${ids.tower})`" />
    <path class="deep" :d="SKYLINE.deep" />
    <path class="tree" :d="SKYLINE.tree" />
    <path class="line" :d="SKYLINE.line" />
  </svg>
</template>

<style scoped>
.skyline {
  --sk: var(--ornament);
  display: block;
  width: 100%;
  height: auto;
  overflow: visible;
}

.stop-wall-top {
  stop-color: color-mix(in srgb, var(--sk) 30%, var(--bg));
}

.stop-wall-bottom {
  stop-color: color-mix(in srgb, var(--sk) 12%, var(--bg));
}

.stop-tower-top {
  stop-color: color-mix(in srgb, var(--sk) 44%, var(--bg));
}

.stop-tower-bottom {
  stop-color: color-mix(in srgb, var(--sk) 17%, var(--bg));
}

.far {
  fill: color-mix(in srgb, var(--sk) 11%, var(--bg));
}

.tile {
  stroke: color-mix(in srgb, var(--sk) 40%, var(--bg));
}

.deep {
  fill: color-mix(in srgb, var(--sk) 36%, var(--bg));
}

.tree {
  fill: color-mix(in srgb, var(--sk) 28%, var(--bg));
}

.line {
  fill: none;
  stroke: color-mix(in srgb, var(--sk) 56%, var(--bg));
  stroke-width: 1.4;
  stroke-linecap: round;
}
</style>
