<script setup lang="ts">
import UiButton from '@/components/ui/UiButton.vue'
import GirihPattern from '@/components/ornament/GirihPattern.vue'
import SuzaniRosette from '@/components/ornament/SuzaniRosette.vue'
import OrnamentStar from '@/components/ornament/OrnamentStar.vue'
import { vReveal } from '@/composables/reveal'
import { useSettingsStore } from '@/stores/settings'

const settings = useSettingsStore()
</script>

<template>
  <section class="cta" aria-labelledby="cta-title">
    <!-- Qora fonda oltin girih va katta suzani — zardo'zi kayfiyati -->
    <GirihPattern tone="gold" :size="80" :opacity="0.22" fade="right" />
    <div class="cta__rosette" aria-hidden="true">
      <SuzaniRosette tone="gold" :opacity="0.55" spin />
    </div>
    <div class="container cta__inner">
      <p class="t-label cta__label" v-reveal><OrnamentStar :size="14" /> Keyingi qadam</p>
      <h2 id="cta-title" class="t-display cta__title" v-reveal="80">G‘oyangiz bormi?<br />Keling, birga quramiz.</h2>
      <div class="cta__bottom" v-reveal="160">
        <p class="cta__text">
          Vazifangizni qisqacha yozing — bir ish kuni ichida javob beramiz. Birinchi suhbat bepul va hech narsaga majburlamaydi.
        </p>
        <div class="cta__actions">
          <UiButton to="/contact" variant="light" icon="arrow-up-right">Loyihani boshlash</UiButton>
          <UiButton v-if="settings.telegramUrl" :href="settings.telegramUrl" variant="outline-light">Telegram’da yozish</UiButton>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.cta {
  position: relative;
  overflow: hidden;
  isolation: isolate;
  background: var(--dark);
  color: var(--on-dark);
  padding-block: clamp(80px, 8.5vw, 136px);
  border-radius: var(--r-xl) var(--r-xl) 0 0;
}

.cta__inner {
  position: relative;
  z-index: 1;
}

.cta__rosette {
  position: absolute;
  top: 50%;
  right: max(-140px, calc((100vw - var(--container)) / 2 - 200px));
  width: min(44vw, 580px);
  aspect-ratio: 1;
  translate: 0 -50%;
}

@media (max-width: 768px) {
  .cta__rosette {
    width: 90vw;
    right: -45vw;
    top: 22%;
    opacity: 0.6;
  }
}

.cta__label {
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--on-dark-2);
  margin-bottom: 28px;
}

.cta__title {
  max-width: 11ch;
}

.cta__bottom {
  margin-top: clamp(28px, 3vw, 40px);
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 32px;
}

.cta__text {
  max-width: 520px;
  font-size: var(--fs-body-lg);
  color: var(--on-dark-2);
}

.cta__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}
</style>
