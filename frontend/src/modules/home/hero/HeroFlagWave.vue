<script setup lang="ts">
import { useId } from 'vue'
import { GIRIH_TILE } from '@/components/ornament/geometry'

// Hilpirayotgan bayroq to'lqini: ko'k — qizil chiziq — oq — qizil chiziq — yashil (O'zbekiston bayrog'i tartibi).
// preserveAspectRatio="none" — har qanday kenglikka cho'ziladi, chiziqlar qalinligi o'zgarmaydi.
const uid = useId()
const ids = { teal: `fw-teal-${uid}`, blue: `fw-blue-${uid}`, green: `fw-green-${uid}`, deep: `fw-deep-${uid}`, tile: `fw-tile-${uid}` }

const CURVE = {
  teal: 'M560 400 C800 330 940 260 1100 150 C1230 60 1330 10 1440 -30',
  blue: 'M0 300 C220 290 400 200 600 214 C800 228 900 250 1060 190 C1220 130 1340 40 1440 0',
  white: 'M0 380 C260 372 470 300 680 310 C880 320 990 322 1140 262 C1280 206 1380 136 1440 104',
  green: 'M0 420 C300 416 520 356 740 360 C920 364 1040 360 1180 306 C1300 260 1390 196 1440 168',
  deep: 'M420 440 C700 430 880 408 1040 404 C1150 400 1240 380 1320 350 C1380 326 1420 300 1440 290',
}
const fill = (d: string) => `${d} L1440 440 L0 440 Z`
</script>

<template>
  <svg class="wave" viewBox="0 0 1440 400" preserveAspectRatio="none" aria-hidden="true" focusable="false">
    <defs>
      <linearGradient :id="ids.teal" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0.35" stop-color="#8fe3d6" stop-opacity="0" />
        <stop offset="1" stop-color="#8fe3d6" stop-opacity="0.75" />
      </linearGradient>
      <linearGradient :id="ids.blue" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#5a9cff" />
        <stop offset="0.55" stop-color="#2a6ff0" />
        <stop offset="1" stop-color="#1b58dc" />
      </linearGradient>
      <linearGradient :id="ids.green" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#3cc768" />
        <stop offset="0.6" stop-color="#1aa74c" />
        <stop offset="1" stop-color="#128a3b" />
      </linearGradient>
      <linearGradient :id="ids.deep" x1="0" y1="0" x2="1" y2="0">
        <stop offset="0" stop-color="#4a8cfb" />
        <stop offset="1" stop-color="#1f62e2" />
      </linearGradient>
      <pattern :id="ids.tile" width="64" height="64" patternUnits="userSpaceOnUse" viewBox="0 0 100 100">
        <g fill="none" stroke="#fff" stroke-width="2" stroke-linejoin="round" opacity="0.22">
          <path v-for="(d, i) in GIRIH_TILE" :key="i" :d="d" />
        </g>
      </pattern>
    </defs>

    <path :d="fill(CURVE.teal)" :fill="`url(#${ids.teal})`" />
    <path :d="fill(CURVE.blue)" :fill="`url(#${ids.blue})`" />
    <path :d="CURVE.blue" class="gloss" />
    <path :d="fill(CURVE.white)" class="white" />
    <path :d="fill(CURVE.green)" :fill="`url(#${ids.green})`" />
    <path :d="fill(CURVE.deep)" :fill="`url(#${ids.deep})`" />
    <path :d="fill(CURVE.deep)" :fill="`url(#${ids.tile})`" />
    <path :d="CURVE.white" class="red red--main" />
    <path :d="CURVE.green" class="red" />
  </svg>
</template>

<style scoped>
.wave {
  display: block;
  width: 100%;
  height: 100%;
}

.white {
  fill: var(--flag-white);
}

.gloss {
  fill: none;
  stroke: rgb(255 255 255 / 0.45);
  stroke-width: 1.5;
  vector-effect: non-scaling-stroke;
}

.red {
  fill: none;
  stroke: #e5262c;
  stroke-width: 1.6;
  vector-effect: non-scaling-stroke;
}

.red--main {
  stroke-width: 3;
}
</style>
