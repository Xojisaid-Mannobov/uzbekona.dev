<script setup lang="ts">
import { vReveal } from '@/composables/reveal'
import GirihPattern from '@/components/ornament/GirihPattern.vue'
import SuzaniRosette from '@/components/ornament/SuzaniRosette.vue'
import OrnamentStar from '@/components/ornament/OrnamentStar.vue'

// Ichki sahifalar uchun katta sarlavha bloki
defineProps<{ label: string; title: string; lead?: string }>()
</script>

<template>
  <header class="page-hero">
    <GirihPattern :size="80" :opacity="0.13" fade="right" />
    <div class="page-hero__rosette" aria-hidden="true"><SuzaniRosette :opacity="0.4" spin /></div>
    <div class="container page-hero__inner">
      <p class="t-label page-hero__label" v-reveal><OrnamentStar :size="14" /> {{ label }}</p>
      <h1 class="t-display page-hero__title" v-reveal="60">{{ title }}</h1>
      <div v-if="lead || $slots.default" class="page-hero__bottom" v-reveal="120">
        <p v-if="lead" class="t-lead page-hero__lead">{{ lead }}</p>
        <slot />
      </div>
    </div>
  </header>
</template>

<style scoped>
.page-hero {
  position: relative;
  overflow: hidden;
  isolation: isolate;
  padding-top: clamp(56px, 6vw, 96px);
  padding-bottom: clamp(48px, 5.4vw, 88px);
}

.page-hero__inner {
  position: relative;
  z-index: 1;
}

.page-hero__rosette {
  position: absolute;
  top: 50%;
  right: max(-160px, calc((100vw - var(--container)) / 2 - 220px));
  width: min(38vw, 480px);
  aspect-ratio: 1;
  translate: 0 -50%;
  z-index: 0;
}

.page-hero__label {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 26px;
}

@media (max-width: 768px) {
  .page-hero__rosette {
    width: 70vw;
    right: -34vw;
    top: 28%;
    opacity: 0.6;
  }
}

.page-hero__title {
  max-width: 14ch;
  white-space: pre-line;
}

.page-hero__bottom {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 40px;
  flex-wrap: wrap;
  margin-top: clamp(32px, 4vw, 56px);
}

.page-hero__lead {
  max-width: 640px;
}
</style>
