<script setup lang="ts">
import { vReveal } from '@/composables/reveal'
import OrnamentStar from '@/components/ornament/OrnamentStar.vue'

// "01 / TANLANGAN LOYIHALAR" uslubidagi raqamli label + katta sarlavha
defineProps<{
  index?: string
  label: string
  title?: string
  lead?: string
  as?: 'h1' | 'h2'
  titleId?: string
}>()
</script>

<template>
  <header class="section-header">
    <p class="section-header__label t-label" v-reveal>
      <OrnamentStar :size="14" />
      <span v-if="index" class="section-header__index">{{ index }} /</span>
      {{ label }}
    </p>
    <div class="section-header__main">
      <component :is="as ?? 'h2'" v-if="title" :id="titleId" class="t-h2 section-header__title" v-reveal="60">
        <slot name="title">{{ title }}</slot>
      </component>
      <div v-if="lead || $slots.aside" class="section-header__aside" v-reveal="120">
        <p v-if="lead" class="t-lead">{{ lead }}</p>
        <slot name="aside" />
      </div>
    </div>
  </header>
</template>

<style scoped>
.section-header {
  display: grid;
  gap: 22px;
  margin-bottom: clamp(44px, 4.6vw, 72px);
}

.section-header__label {
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-header__index {
  color: var(--ink);
}

.section-header__main {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(0, 1fr);
  gap: 48px;
  align-items: end;
}

.section-header__title {
  max-width: 20ch;
  white-space: pre-line;
}

.section-header__aside {
  display: grid;
  gap: 28px;
  justify-items: start;
  max-width: 480px;
  justify-self: end;
}

@media (max-width: 1024px) {
  .section-header__main {
    grid-template-columns: 1fr;
    gap: 28px;
  }

  .section-header__aside {
    justify-self: start;
  }
}
</style>
