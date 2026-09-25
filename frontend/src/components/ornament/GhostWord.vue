<script setup lang="ts">
// Footer'dagi "UZBEKONA" kabi ulkan "arvoh" so'z — konteyner kengligini to'liq egallaydi (SVG textLength),
// pastki qismi keyingi blokka kirib turadi. Faqat bezak.
withDefaults(defineProps<{ text: string; tone?: 'light' | 'dark' }>(), { tone: 'light' })
</script>

<template>
  <svg
    class="ghost-word"
    :class="`ghost-word--${tone}`"
    viewBox="0 0 1000 150"
    preserveAspectRatio="xMidYMax meet"
    aria-hidden="true"
    focusable="false"
  >
    <text x="500" y="148" text-anchor="middle" textLength="980" lengthAdjust="spacingAndGlyphs">{{ text }}</text>
  </svg>
</template>

<style scoped>
.ghost-word {
  display: block;
  width: 100%;
  height: auto;
  margin-bottom: -2.6vw;
  font-family: var(--font);
  font-size: 196px;
  font-weight: 700;
  letter-spacing: -0.05em;
  user-select: none;
  pointer-events: none;
}

.ghost-word--light {
  fill: color-mix(in srgb, var(--ink) 6%, var(--bg));
}

.ghost-word--dark {
  fill: var(--dark-2);
}
</style>
