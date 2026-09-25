<script setup lang="ts">
import SectionHeader from '@/components/ui/SectionHeader.vue'
import { processSteps } from '@/content/site'
import { vReveal } from '@/composables/reveal'
import { pad2 } from '@/utils/format'

withDefaults(defineProps<{ index?: string }>(), { index: '04' })
</script>

<template>
  <section class="section process" aria-labelledby="process-title">
    <div class="container">
      <SectionHeader
        title-id="process-title"
        :index="index"
        label="Ish jarayoni"
        title="Oltita aniq bosqich."
        lead="Har bosqich oxirida ko‘rish, sinash va fikr bildirish mumkin bo‘lgan natija bo‘ladi."
      />

      <ol role="list" class="process__grid">
        <li v-for="(step, i) in processSteps" :key="step.title" class="step" v-reveal="(i % 3) * 90">
          <span class="step__num">{{ pad2(i + 1) }}</span>
          <h3 class="step__title">{{ step.title }}</h3>
          <p class="step__text">{{ step.text }}</p>
        </li>
      </ol>
    </div>
  </section>
</template>

<style scoped>
.process__grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0 32px;
}

.step {
  display: grid;
  align-content: start;
  gap: 16px;
  padding: 40px 0 64px;
  border-top: 1px solid var(--line-strong);
}

.step__num {
  font-size: var(--fs-step);
  line-height: 0.9;
  font-weight: 600;
  letter-spacing: -0.06em;
  color: var(--accent);
  margin-bottom: 24px;
  font-variant-numeric: tabular-nums;
}

.step__title {
  font-size: clamp(21px, 1.7vw, 26px);
  letter-spacing: -0.03em;
}

.step__text {
  max-width: 36ch;
  font-size: 16px;
  color: var(--ink-2);
}

@media (max-width: 1024px) {
  .process__grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .process__grid {
    grid-template-columns: 1fr;
  }

  .step {
    padding: 32px 0 48px;
  }
}
</style>
