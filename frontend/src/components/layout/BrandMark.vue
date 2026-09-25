<script setup lang="ts">
import { useId } from 'vue'
import { GIRIH_TILE } from '@/components/ornament/geometry'

/**
 * Uzbekona.dev belgisi — lola shaklidagi "U":
 * chap bargi koshin naqshli ko'k, o'ng bargi O'zbekiston bayrog'i ranglarida, tepasida "</>".
 */
withDefaults(defineProps<{ code?: boolean }>(), { code: true })

const uid = useId()
const ids = {
  blue: `bm-blue-${uid}`,
  sky: `bm-sky-${uid}`,
  green: `bm-green-${uid}`,
  tile: `bm-tile-${uid}`,
  right: `bm-right-${uid}`,
}

// Ikki barg: chap (naqshli) va o'ng (bayroq), pastda S-chiziq bilan tutashadi
const LEFT = 'M6 20 C24 19 40 27 46 43 L46 76 C46 90 52 99 60 100 C64 110 64 120 60 127 C28 127 6 107 6 79 Z'
const RIGHT = 'M114 20 C96 19 80 27 74 43 L74 76 C74 90 68 99 60 100 C64 110 64 120 60 127 C92 127 114 107 114 79 Z'
</script>

<template>
  <svg class="mark" :viewBox="code ? '0 -26 120 153' : '0 16 120 111'" aria-hidden="true" focusable="false">
    <defs>
      <linearGradient :id="ids.blue" x1="0" y1="0" x2="0.4" y2="1">
        <stop offset="0" stop-color="#2f86ff" />
        <stop offset="1" stop-color="#1552d6" />
      </linearGradient>
      <linearGradient :id="ids.sky" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0" stop-color="#1f6ff0" />
        <stop offset="1" stop-color="#39b6ff" />
      </linearGradient>
      <linearGradient :id="ids.green" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0" stop-color="#22b04b" />
        <stop offset="1" stop-color="#138a3a" />
      </linearGradient>
      <pattern :id="ids.tile" width="30" height="30" patternUnits="userSpaceOnUse" viewBox="0 0 100 100">
        <g fill="none" stroke="#a9d6ff" stroke-width="6" stroke-linejoin="round" opacity="0.7">
          <path v-for="(d, i) in GIRIH_TILE" :key="i" :d="d" />
        </g>
      </pattern>
      <clipPath :id="ids.right"><path :d="RIGHT" /></clipPath>
    </defs>

    <g v-if="code" fill="none" stroke-linecap="round" stroke-linejoin="round" stroke-width="7">
      <path d="M40 -16 L28 -4 L40 8" stroke="#1f6ff0" />
      <path d="M80 -16 L92 -4 L80 8" stroke="#1f6ff0" />
      <path d="M66 -20 L54 12" stroke="#e5262c" />
    </g>

    <path :d="LEFT" :fill="`url(#${ids.blue})`" />
    <path :d="LEFT" :fill="`url(#${ids.tile})`" />

    <g :clip-path="`url(#${ids.right})`">
      <rect x="40" y="10" width="90" height="130" :fill="`url(#${ids.sky})`" />
      <!-- oq band va yashil band, orasida qizil chiziqlar (bayroq tartibi) -->
      <path d="M50 84 C72 76 94 62 124 36 L124 140 L50 140 Z" fill="#fff" />
      <path d="M50 102 C74 94 96 80 124 56 L124 140 L50 140 Z" :fill="`url(#${ids.green})`" />
      <path d="M50 84 C72 76 94 62 124 36" fill="none" stroke="#e5262c" stroke-width="3" />
      <path d="M50 101 C74 93 96 79 124 55" fill="none" stroke="#e5262c" stroke-width="1.6" />
    </g>
  </svg>
</template>

<style scoped>
.mark {
  display: block;
  height: 100%;
  width: auto;
  overflow: visible;
}
</style>
