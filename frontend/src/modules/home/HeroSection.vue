<script setup lang="ts">
import UiButton from '@/components/ui/UiButton.vue'
import GirihPattern from '@/components/ornament/GirihPattern.vue'
import SuzaniRosette from '@/components/ornament/SuzaniRosette.vue'
import { hero } from '@/content/site'
import { useSettingsStore } from '@/stores/settings'
import { scrollToTarget } from '@/composables/useSmoothScroll'

// Hero animatsiyasi sof CSS'da — birinchi ekran JavaScript kutubxonalarini kutmaydi (LCP).
// Fonda girih to'ri va sekin aylanuvchi suzani rozetkasi — o'zbekona kayfiyat.
const settings = useSettingsStore()
</script>

<template>
  <section class="hero" aria-labelledby="hero-title">
    <div class="hero__decor" aria-hidden="true">
      <GirihPattern :size="88" :opacity="0.16" fade="right" />
      <div class="hero__rosette">
        <SuzaniRosette :opacity="0.45" spin />
      </div>
    </div>
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
  position: relative;
  min-height: clamp(560px, calc(92svh - var(--nav-h)), 820px);
  display: flex;
  align-items: center;
  padding-block: 56px 88px;
  overflow: hidden;
  isolation: isolate;
}

.hero__decor {
  position: absolute;
  inset: 0;
  z-index: -1;
}

/* Rozetka o'ng tomonda, yarmi ekrandan chiqib turadi — naqsh matnga xalaqit bermaydi */
.hero__rosette {
  position: absolute;
  top: 50%;
  right: max(-120px, calc((100vw - var(--container)) / 2 - 200px));
  width: min(46vw, 620px);
  aspect-ratio: 1;
  translate: 0 -50%;
}

.hero__rosette::before {
  content: '';
  position: absolute;
  inset: 12%;
  border-radius: 50%;
  background: radial-gradient(circle, color-mix(in srgb, var(--ornament) 10%, transparent), transparent 70%);
}

.hero__inner {
  display: grid;
  gap: 32px;
}

.hero__eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  width: fit-content;
  height: 40px;
  padding: 0 18px 0 14px;
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
  max-width: 880px;
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
  flex-direction: column;
  align-items: flex-start;
  gap: 32px;
  margin-top: 4px;
}

.hero__subtitle {
  max-width: 560px;
  font-size: var(--fs-body-lg);
  line-height: 1.55;
  color: var(--ink-2);
}

.hero__actions {
  display: flex;
  gap: 12px;
  flex-shrink: 0;
}

@media (max-width: 1024px) {
  .hero__rosette {
    width: 60vw;
    right: -22vw;
    top: 30%;
  }
}

@media (max-width: 560px) {
  .hero {
    padding-block: 32px 56px;
  }

  .hero__rosette {
    width: 88vw;
    right: -40vw;
    top: 18%;
    opacity: 0.6;
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
