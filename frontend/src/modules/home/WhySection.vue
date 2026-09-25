<script setup lang="ts">
import SectionHeader from '@/components/ui/SectionHeader.vue'
import { principles } from '@/content/site'
import { vReveal } from '@/composables/reveal'
import { pad2 } from '@/utils/format'
import GirihPattern from '@/components/ornament/GirihPattern.vue'

withDefaults(defineProps<{ index?: string }>(), { index: '05' })
</script>

<template>
  <section class="section why" aria-labelledby="why-title">
    <div class="container">
      <SectionHeader title-id="why-title" :index="index" label="Nega Uzbekona" title="Kod — natijaning faqat bir qismi." />

      <div class="why__grid">
        <article v-for="(p, i) in principles" :key="p.title" class="why__block" v-reveal="(i % 2) * 110">
          <GirihPattern :size="64" :opacity="0.14" fade="right" />
          <span class="why__num">{{ pad2(i + 1) }}</span>
          <div>
            <h3 class="t-h3">{{ p.title }}</h3>
            <p class="why__text">{{ p.text }}</p>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped>
.why__grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: clamp(16px, 1.6vw, 24px);
}

.why__block {
  position: relative;
  overflow: hidden;
  isolation: isolate;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 64px;
  min-height: 290px;
  padding: clamp(32px, 3.4vw, 56px);
  background: var(--surface);
  border-radius: var(--r-lg);
  border: 1px solid var(--line);
  transition:
    transform var(--dur) var(--ease),
    box-shadow var(--dur) var(--ease);
}

.why__block:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.why__num {
  font-size: var(--fs-small);
  font-weight: 700;
  color: var(--ink-3);
}

.why__text {
  margin-top: 16px;
  max-width: 30ch;
  font-size: var(--fs-body-lg);
  color: var(--ink-2);
}

@media (max-width: 768px) {
  .why__grid {
    grid-template-columns: 1fr;
  }

  .why__block {
    min-height: 240px;
  }
}
</style>
