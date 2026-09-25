<script setup lang="ts">
// Chust do'ppisi (yon tomondan 3/4 ko'rinish): qora gumbaz, har yuzida oq "qalampir" naqshi,
// pastida jiyak (hoshiya bandi). Wordmark'dagi "O" harfiga kiydirish uchun.
</script>

<template>
  <svg class="doppi" viewBox="0 0 120 80" aria-hidden="true" focusable="false">
    <!-- gumbaz -->
    <path d="M7 71 L13 32 C14 22 30 12 60 11 C90 12 106 22 107 32 L113 71 Z" class="doppi__body" />
    <!-- old tomondagi qirra (to'rt burchakli do'ppining burchagi) -->
    <path d="M60 12 Q58.5 42 60 60" class="doppi__edge" />
    <!-- tepa chok -->
    <path d="M17 31 C19 23 33 16.5 60 16 C87 16.5 101 23 103 31" class="doppi__stitch" />

    <!-- qalampir naqshlari: chap va o'ng yuzlarda, bir-biriga ko'zgu; ichida qora kontur -->
    <g class="doppi__pepper">
      <path d="M32 55 C22 55 20 42 27 37 C31 34 36 33 39 29 C41 26 40 22 36 22 C39 19 45 22 45 29 C45 36 43 55 32 55 Z" />
      <path d="M88 55 C98 55 100 42 93 37 C89 34 84 33 81 29 C79 26 80 22 84 22 C81 19 75 22 75 29 C75 36 77 55 88 55 Z" />
    </g>
    <g class="doppi__seed">
      <path
        d="M32.7 49.6 C27.2 49.6 26.1 42.5 29.9 39.7 C32.1 38.0 34.9 37.5 36.5 35.3 C37.6 33.6 37.1 31.4 34.9 31.4 C36.5 29.8 39.8 31.4 39.8 35.3 C39.8 39.1 38.7 49.6 32.7 49.6 Z"
      />
      <path
        d="M87.3 49.6 C92.8 49.6 93.9 42.5 90.1 39.7 C87.9 38 85.1 37.5 83.5 35.3 C82.4 33.6 82.9 31.4 85.1 31.4 C83.5 29.8 80.2 31.4 80.2 35.3 C80.2 39.1 81.3 49.6 87.3 49.6 Z"
      />
    </g>

    <!-- jiyak: oq chiziqlar va mayda "tishlar" -->
    <path d="M9.3 60 L110.7 60 L113 71 L7 71 Z" class="doppi__band" />
    <path d="M9.3 60 L110.7 60" class="doppi__line" />
    <path
      d="M12 66 l4 -3 l4 3 l4 -3 l4 3 l4 -3 l4 3 l4 -3 l4 3 l4 -3 l4 3 l4 -3 l4 3 l4 -3 l4 3 l4 -3 l4 3 l4 -3 l4 3 l4 -3 l4 3 l4 -3 l4 3 l4 -3 l4 3"
      class="doppi__zigzag"
    />
  </svg>
</template>

<style scoped>
.doppi {
  display: block;
  overflow: visible;
}

.doppi__body,
.doppi__band {
  fill: var(--doppi-body, #0a0b0b);
  stroke: var(--doppi-outline, transparent);
  stroke-width: 1.5;
  stroke-linejoin: round;
}

.doppi__edge {
  fill: none;
  stroke: rgb(255 255 255 / 0.14);
  stroke-width: 1.4;
}

.doppi__stitch,
.doppi__line {
  fill: none;
  stroke: var(--doppi-ink, #f3f3ef);
  stroke-width: 1.3;
  opacity: 0.75;
}

.doppi__stitch {
  stroke-dasharray: 3 3;
}

.doppi__pepper path {
  fill: var(--doppi-ink, #f3f3ef);
}

/* Qalampir ichidagi qora ichki kontur */
.doppi__seed path {
  fill: none;
  stroke: var(--doppi-body, #0a0b0b);
  stroke-width: 1.8;
}

.doppi__zigzag {
  fill: none;
  stroke: var(--doppi-ink, #f3f3ef);
  stroke-width: 1.3;
  stroke-linejoin: round;
  opacity: 0.8;
}
</style>
