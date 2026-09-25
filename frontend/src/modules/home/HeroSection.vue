<script setup lang="ts">
import UiButton from '@/components/ui/UiButton.vue'
import { hero } from '@/content/site'
import { useSettingsStore } from '@/stores/settings'
import { scrollToTarget } from '@/composables/useSmoothScroll'

// Hero animatsiyasi sof CSS'da — birinchi ekran JavaScript kutubxonalarini kutmaydi (LCP)
const settings = useSettingsStore()
</script>

<template>
  <section class="hero" aria-labelledby="hero-title">
    <div class="container hero__inner">
      <p class="hero__eyebrow hero__fade" style="--d: 0ms">
        <span class="hero__dot" :class="{ 'is-on': settings.site.available }" aria-hidden="true" />
        {{ settings.site.available ? 'Yangi loyihalar uchun ochiqmiz' : 'Digital Product Studio' }}
      </p>

      <h1 id="hero-title" class="t-hero hero__title">
        <span v-for="(line, i) in hero.title" :key="line" class="hero__line"
          ><span :style="{ '--i': i }">{{ line }}</span></span
        >
      </h1>

      <div class="hero__bottom">
        <p class="hero__subtitle hero__fade" style="--d: 520ms">{{ hero.subtitle }}</p>
        <div class="hero__actions hero__fade" style="--d: 600ms">
          <UiButton icon="arrow-down" @click="scrollToTarget('#projects')">Loyihalarni ko‘rish</UiButton>
          <UiButton to="/contact" variant="secondary" icon="arrow-up-right">Biz bilan ishlash</UiButton>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.hero {
  min-height: calc(100vh - var(--nav-h));
  min-height: calc(100svh - var(--nav-h));
  display: flex;
  align-items: center;
  padding-block: 64px 96px;
}

.hero__inner {
  display: grid;
  gap: 40px;
}

.hero__eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  width: fit-content;
  height: 44px;
  padding: 0 20px 0 16px;
  border-radius: var(--r-pill);
  border: 1px solid var(--line);
  background: var(--surface);
  font-size: var(--fs-small);
  font-weight: 600;
  color: var(--ink-2);
}

.hero__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--ink-3);
}

.hero__dot.is-on {
  background: var(--success);
  box-shadow: 0 0 0 4px color-mix(in srgb, var(--success) 18%, transparent);
  animation: pulse 2.4s var(--ease) infinite;
}

@keyframes pulse {
  50% {
    box-shadow: 0 0 0 7px color-mix(in srgb, var(--success) 6%, transparent);
  }
}

.hero__title {
  max-width: 1200px;
}

.hero__line {
  display: block;
  overflow: hidden;
  /* descender'lar (g, y) kesilmasligi uchun */
  padding-bottom: 0.14em;
  margin-bottom: -0.14em;
}

.hero__line > span {
  display: inline-block;
  /* Qatorlar pastdan ko'tariladi (text reveal) */
  animation: hero-rise 950ms cubic-bezier(0.16, 1, 0.3, 1) both;
  animation-delay: calc(var(--i) * 90ms + 80ms);
}

.hero__fade {
  animation: hero-fade 800ms var(--ease) both;
  animation-delay: var(--d, 0ms);
}

@keyframes hero-rise {
  from {
    transform: translateY(105%);
  }
}

@keyframes hero-fade {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
}

.hero__bottom {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 48px;
  margin-top: 16px;
}

.hero__subtitle {
  max-width: 680px;
  font-size: clamp(18px, 1.5vw, 22px);
  line-height: 1.55;
  color: var(--ink-2);
}

.hero__actions {
  display: flex;
  gap: 12px;
  flex-shrink: 0;
}

@media (max-width: 1024px) {
  .hero__bottom {
    flex-direction: column;
    align-items: flex-start;
    gap: 32px;
  }
}

@media (max-width: 560px) {
  .hero {
    padding-block: 40px 64px;
  }

  .hero__inner {
    gap: 28px;
  }

  .hero__actions {
    flex-direction: column;
    width: 100%;
  }

  .hero__actions > * {
    width: 100%;
  }
}
</style>
