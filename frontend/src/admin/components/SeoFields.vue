<script setup lang="ts">
import FormField from './FormField.vue'
import type { Seo } from '@/types/api'

// SEO sarlavha va description + Google natijasi ko'rinishi
const model = defineModel<Seo>({ required: true })
defineProps<{ fallbackTitle?: string; fallbackDescription?: string; path?: string }>()
</script>

<template>
  <div class="a-stack">
    <FormField label="SEO sarlavha" :hint="`${model.title.length}/60 — bo‘sh bo‘lsa asosiy sarlavha ishlatiladi`" for="seo-title">
      <input id="seo-title" v-model="model.title" class="a-input" maxlength="160" :placeholder="fallbackTitle" />
    </FormField>
    <FormField label="Meta description" :hint="`${model.description.length}/160`" for="seo-desc">
      <textarea id="seo-desc" v-model="model.description" class="a-textarea" maxlength="320" :placeholder="fallbackDescription" />
    </FormField>
    <div class="serp" aria-label="Qidiruv natijasi ko‘rinishi">
      <p class="serp__url">uzbekona.dev{{ path }}</p>
      <p class="serp__title">{{ model.title || fallbackTitle || 'Sarlavha' }} — Uzbekona.dev</p>
      <p class="serp__desc">{{ model.description || fallbackDescription || 'Sahifa tavsifi shu yerda ko‘rinadi.' }}</p>
    </div>
  </div>
</template>

<style scoped>
.serp {
  padding: 16px;
  border-radius: var(--a-radius);
  background: var(--surface-2);
  font-family: Arial, sans-serif;
}

.serp__url {
  font-size: 12px;
  color: var(--ink-3);
}

.serp__title {
  margin-top: 4px;
  font-size: 17px;
  color: #1a0dab;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:root[data-theme='dark'] .serp__title {
  color: #8ab4f8;
}

.serp__desc {
  margin-top: 2px;
  font-size: 13px;
  color: var(--ink-2);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
