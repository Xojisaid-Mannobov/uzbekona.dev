<script setup lang="ts">
import doppiUrl from '@/assets/ornament/doppi.webp'

// Haqiqiy Chust do'ppisi (3/4 ko'rinish, shaffof fon): qora paxta, oq ipak "qalampir" naqshlari va jiyak.
// Rasm SVG <image> ichida — hero'dagi "o" harfiga ham, footer'dagi ichma-ich SVG'ga ham (x/y/width/height) bir xil joylashadi.
// Kadr 120x80: jiyakning pastki cheti y≈77 da (eski chizma bilan bir xil tayanch nuqta).
</script>

<template>
  <svg class="doppi" viewBox="0 0 120 80" aria-hidden="true" focusable="false">
    <image :href="doppiUrl" x="0" y="0" width="120" height="80" preserveAspectRatio="xMidYMax meet" />
  </svg>
</template>

<style scoped>
.doppi {
  display: block;
  overflow: visible;
  /* qora do'ppi qorong'i fonda yo'qolib qolmasligi uchun nozik yorug' kontur (--doppi-outline) */
  filter: drop-shadow(0 0 0.6px var(--doppi-outline, transparent)) drop-shadow(0 3px 4px rgb(0 0 0 / 0.22));
}
</style>
