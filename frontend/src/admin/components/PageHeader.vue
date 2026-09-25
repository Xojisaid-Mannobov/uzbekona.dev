<script setup lang="ts">
import { RouterLink, type RouteLocationRaw } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'

defineProps<{ title: string; subtitle?: string; back?: RouteLocationRaw }>()
</script>

<template>
  <header class="ph">
    <div class="ph__text">
      <RouterLink v-if="back" :to="back" class="ph__back"><AppIcon name="arrow-left" :size="16" /> Orqaga</RouterLink>
      <h1>{{ title }}</h1>
      <p v-if="subtitle">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.default" class="ph__actions"><slot /></div>
  </header>
</template>

<style scoped>
.ph {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 20px;
  flex-wrap: wrap;
  margin-bottom: 24px;
}

.ph__back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink-3);
  margin-bottom: 8px;
}

.ph__back:hover {
  color: var(--ink);
}

h1 {
  font-size: clamp(26px, 2.4vw, 34px);
  letter-spacing: -0.035em;
}

p {
  margin-top: 4px;
  color: var(--ink-2);
}

.ph__actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}
</style>
