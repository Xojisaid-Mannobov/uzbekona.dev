<script setup lang="ts">
import UiButton from '@/components/ui/UiButton.vue'
import { vReveal } from '@/composables/reveal'
import { useSettingsStore } from '@/stores/settings'

const settings = useSettingsStore()
</script>

<template>
  <section class="cta" aria-labelledby="cta-title">
    <div class="container cta__inner">
      <p class="t-label cta__label" v-reveal>Keyingi qadam</p>
      <h2 id="cta-title" class="t-display cta__title" v-reveal="80">Biror narsa<br />qurmoqchimisiz?</h2>
      <div class="cta__bottom" v-reveal="160">
        <p class="cta__text">
          G‘oyangizni, muammoingizni yoki mavjud jarayoningizni aytib bering. Birinchi suhbat bepul va hech narsaga majburlamaydi.
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
  background: var(--dark);
  color: var(--on-dark);
  padding-block: clamp(96px, 11vw, 180px);
  border-radius: var(--r-xl) var(--r-xl) 0 0;
}

.cta__label {
  color: var(--on-dark-2);
  margin-bottom: 32px;
}

.cta__title {
  max-width: 11ch;
}

.cta__bottom {
  margin-top: clamp(48px, 5vw, 80px);
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 48px;
  flex-wrap: wrap;
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
