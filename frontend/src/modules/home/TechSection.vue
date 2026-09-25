<script setup lang="ts">
import SectionHeader from '@/components/ui/SectionHeader.vue'
import { technologies } from '@/content/site'
import { vReveal } from '@/composables/reveal'

// Text-based texnologiyalar taqdimoti: kategoriya → yirik yozuvlar
withDefaults(defineProps<{ index?: string }>(), { index: '03' })
</script>

<template>
  <section class="section tech" aria-labelledby="tech-title">
    <div class="container">
      <SectionHeader
        title-id="tech-title"
        :index="index"
        label="Texnologiyalar"
        :title="'Chiroyli interfeysning\norqasida kuchli engineering.'"
        lead="Har bir loyiha uchun vositalarni vazifaga qarab tanlaymiz — moda uchun emas, yillar davomida ishlashi uchun."
      />

      <div class="tech__groups">
        <div v-for="(g, i) in technologies" :key="g.group" class="tech__group" v-reveal="i * 70">
          <p class="tech__name">{{ g.group }}</p>
          <ul role="list" class="tech__items">
            <li v-for="item in g.items" :key="item" class="tech__item">{{ item }}</li>
          </ul>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.tech__groups {
  border-top: 1px solid var(--line-strong);
}

.tech__group {
  display: grid;
  grid-template-columns: minmax(200px, 1fr) 3fr;
  gap: 32px;
  align-items: baseline;
  padding-block: 36px;
  border-bottom: 1px solid var(--line);
}

.tech__name {
  font-size: var(--fs-small);
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: var(--ink-3);
}

.tech__items {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 14px;
}

.tech__item {
  display: inline-flex;
  align-items: center;
  height: 64px;
  padding: 0 28px;
  border-radius: var(--r-pill);
  background: var(--surface);
  border: 1px solid var(--line);
  font-size: clamp(20px, 1.9vw, 30px);
  font-weight: 600;
  letter-spacing: -0.03em;
  transition:
    background-color var(--dur) var(--ease),
    color var(--dur) var(--ease),
    border-color var(--dur) var(--ease);
}

.tech__item:hover {
  background: var(--ink);
  color: var(--bg);
  border-color: var(--ink);
}

@media (max-width: 768px) {
  .tech__group {
    grid-template-columns: 1fr;
    gap: 16px;
    padding-block: 28px;
  }

  .tech__item {
    height: 52px;
    padding: 0 22px;
  }
}
</style>
