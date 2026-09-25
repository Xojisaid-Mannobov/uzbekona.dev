<script setup lang="ts">
import { vReveal } from '@/composables/reveal'

// Ichki sahifalar uchun katta sarlavha bloki
defineProps<{ label: string; title: string; lead?: string }>()
</script>

<template>
  <header class="page-hero">
    <div class="container">
      <p class="t-label page-hero__label" v-reveal>{{ label }}</p>
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
  padding-top: clamp(64px, 8vw, 128px);
  padding-bottom: clamp(56px, 7vw, 112px);
}

.page-hero__label {
  margin-bottom: 32px;
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
