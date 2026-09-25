<script setup lang="ts">
import { RouterLink } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import { aboutStatement } from '@/content/site'
import { vReveal } from '@/composables/reveal'
import OrnamentStar from '@/components/ornament/OrnamentStar.vue'
import SuzaniRosette from '@/components/ornament/SuzaniRosette.vue'

// Katta bayonot (48–64px): ikki bosqichli yumshoq reveal
</script>

<template>
  <section class="section statement" aria-labelledby="about-title">
    <div class="container statement__grid">
      <div class="statement__side">
        <p class="t-label statement__label"><OrnamentStar :size="14" /> {{ aboutStatement.label }}</p>
        <SuzaniRosette class="statement__rosette" tone="gold" :opacity="0.7" />
      </div>
      <div>
        <h2 id="about-title" class="t-statement statement__lead" v-reveal>{{ aboutStatement.lead }}</h2>
        <p class="t-statement statement__text" v-reveal="120">{{ aboutStatement.text }}</p>
        <RouterLink to="/about" class="statement__link"> Biz haqimizda batafsil <AppIcon name="arrow-up-right" :size="20" /> </RouterLink>
      </div>
    </div>
  </section>
</template>

<style scoped>
.statement__grid {
  display: grid;
  grid-template-columns: minmax(180px, 1fr) minmax(0, 3fr);
  gap: 48px;
}

.statement__side {
  display: grid;
  align-content: start;
  gap: 40px;
}

.statement__rosette {
  width: min(100%, 200px);
}

.statement__label {
  display: flex;
  align-items: center;
  gap: 10px;
  padding-top: 18px;
}

.statement__lead {
  color: var(--ink);
}

.statement__text {
  margin-top: 0.1em;
  max-width: 20ch;
  color: var(--ink-2);
}

.statement__link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-top: 56px;
  font-size: 18px;
  font-weight: 600;
  border-bottom: 1px solid var(--line-strong);
  padding-bottom: 4px;
  transition:
    border-color var(--dur-fast) var(--ease),
    color var(--dur-fast) var(--ease);
}

.statement__link:hover {
  color: var(--accent);
  border-color: var(--accent);
}

@media (max-width: 768px) {
  .statement__grid {
    grid-template-columns: 1fr;
    gap: 24px;
  }

  .statement__side {
    grid-template-columns: 1fr auto;
    align-items: center;
    gap: 16px;
  }

  .statement__rosette {
    width: 88px;
  }

  .statement__label {
    padding-top: 0;
  }
}
</style>
