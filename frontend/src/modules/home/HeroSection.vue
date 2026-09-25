<script setup lang="ts">
import UiButton from '@/components/ui/UiButton.vue'
import BrandMark from '@/components/layout/BrandMark.vue'
import Doppi from '@/components/ornament/Doppi.vue'
import HeroSkyline from './hero/HeroSkyline.vue'
import HeroFlagWave from './hero/HeroFlagWave.vue'
import HeroCircuit from './hero/HeroCircuit.vue'
import { hero } from '@/content/site'
import { useSettingsStore } from '@/stores/settings'
import { scrollToTarget } from '@/composables/useSmoothScroll'

// Hero: brend belgisi + Registon silueti + hilpirayotgan bayroq — "Raqamli O'zbekiston" kayfiyati.
// Illyustratsiya to'liq SVG (rasm yuklanmaydi), animatsiya sof CSS'da — LCP JavaScript'ni kutmaydi.
const settings = useSettingsStore()
</script>

<template>
  <section class="hero" aria-labelledby="hero-title">
    <div class="hero__art" aria-hidden="true">
      <HeroCircuit variant="top" class="hero__circuit hero__circuit--top" />
      <HeroCircuit variant="bottom" class="hero__circuit hero__circuit--bottom" />
      <HeroSkyline class="hero__skyline" />
      <HeroFlagWave class="hero__wave" />
    </div>

    <div class="container hero__inner">
      <div class="hero__brand">
        <BrandMark class="hero__mark" />
        <div class="hero__brand-text">
          <p class="hero__tagline hero__fade" style="--d: 80ms">{{ hero.tagline }}</p>
          <h1 id="hero-title" class="hero__title">
            <span class="hero__rise"
              >Uzbek<span class="hero__o">o<Doppi class="hero__doppi" /></span>na<span class="hero__tld">.dev</span></span
            >
          </h1>
          <p class="hero__services hero__fade" style="--d: 300ms">
            <template v-for="(s, i) in hero.services" :key="s">
              <span v-if="i" class="hero__sep" aria-hidden="true" />
              <span>{{ s }}</span>
            </template>
          </p>
        </div>
      </div>

      <p class="hero__lead hero__fade" style="--d: 420ms">{{ hero.lead }}</p>

      <div class="hero__actions hero__fade" style="--d: 520ms">
        <UiButton icon="arrow-down" @click="scrollToTarget('#projects')">Loyihalarni ko‘rish</UiButton>
        <UiButton to="/contact" variant="secondary" icon="arrow-up-right">Biz bilan ishlash</UiButton>
      </div>

      <p v-if="settings.site.available" class="hero__status hero__fade" style="--d: 620ms">
        <span class="hero__dot" aria-hidden="true" />
        Yangi loyihalar uchun ochiqmiz
      </p>
    </div>
  </section>
</template>

<style scoped>
.hero {
  --wave-h: clamp(190px, 27vw, 400px);
  --edge: max(var(--gutter), (100vw - var(--container)) / 2);
  position: relative;
  min-height: clamp(640px, calc(100svh - var(--nav-h)), 860px);
  display: flex;
  align-items: center;
  padding-top: 48px;
  padding-bottom: calc(var(--wave-h) * 0.72);
  overflow: hidden;
  isolation: isolate;
  background: radial-gradient(60% 70% at 78% 46%, color-mix(in srgb, var(--ornament) 11%, transparent), transparent 70%);
}

/* ─── Illyustratsiya ─────────────────────────────── */
.hero__art {
  position: absolute;
  inset: 0;
  z-index: -1;
  pointer-events: none;
}

.hero__circuit {
  position: absolute;
  opacity: 0.8;
}

.hero__circuit--top {
  top: -72px;
  left: 50%;
  width: min(38vw, 560px);
}

.hero__circuit--bottom {
  left: 0;
  bottom: 12px;
  width: min(32vw, 460px);
}

.hero__skyline {
  position: absolute;
  right: calc(var(--edge) - 24px);
  bottom: calc(var(--wave-h) * 0.5);
  width: min(52vw, 780px);
  mask-image: linear-gradient(90deg, transparent 0%, #000 20%), linear-gradient(180deg, #000 58%, transparent 96%);
  mask-composite: intersect;
  animation: skyline-in 1400ms cubic-bezier(0.16, 1, 0.3, 1) 150ms both;
}

.hero__wave {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: var(--wave-h);
  mask-image: linear-gradient(90deg, transparent 0%, rgb(0 0 0 / 0.35) 22%, #000 48%), linear-gradient(180deg, #000 84%, transparent 100%);
  mask-composite: intersect;
  /* bayroq chapdan o'ngga "yoyiladi" */
  animation: wave-in 1600ms cubic-bezier(0.65, 0, 0.35, 1) 250ms both;
}

@keyframes skyline-in {
  from {
    opacity: 0;
    transform: translateY(32px);
  }
}

@keyframes wave-in {
  from {
    clip-path: inset(0 100% 0 0);
  }
  to {
    clip-path: inset(0 0 0 0);
  }
}

/* ─── Matn ───────────────────────────────────────── */
.hero__inner {
  display: grid;
  justify-items: start;
  gap: 32px;
}

.hero__brand {
  display: flex;
  align-items: center;
  gap: clamp(20px, 2vw, 32px);
}

.hero__mark {
  height: clamp(96px, 9.6vw, 148px);
  flex-shrink: 0;
  animation: mark-in 1100ms cubic-bezier(0.16, 1, 0.3, 1) both;
}

@keyframes mark-in {
  from {
    opacity: 0;
    transform: scale(0.82) translateY(12px);
  }
}

.hero__brand-text {
  display: grid;
  gap: 12px;
  min-width: 0;
}

.hero__tagline {
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.26em;
  text-transform: uppercase;
  color: var(--ink-2);
}

.hero__title {
  font-size: clamp(48px, 5.6vw, 84px);
  font-weight: 800;
  line-height: 1;
  letter-spacing: -0.045em;
  overflow: hidden;
  /* tepada do'ppi uchun joy, pastda descender'lar uchun */
  padding-block: 0.3em 0.08em;
  margin-block: -0.3em -0.08em;
}

.hero__rise {
  display: inline-block;
  animation: hero-rise 950ms cubic-bezier(0.16, 1, 0.3, 1) 120ms both;
}

.hero__tld {
  color: var(--accent);
}

/* "o" harfiga kiydirilgan do'ppi — o'lchami harf bilan birga (em) o'zgaradi */
.hero__o {
  position: relative;
  display: inline-block;
}

.hero__doppi {
  position: absolute;
  left: 50%;
  bottom: 0.45em;
  width: 0.84em;
  translate: -52% 0;
  rotate: -10deg;
  --doppi-outline: color-mix(in srgb, var(--ink) 35%, transparent);
  animation: doppi-drop 900ms cubic-bezier(0.34, 1.56, 0.64, 1) 700ms both;
}

@keyframes doppi-drop {
  from {
    opacity: 0;
    translate: -52% -0.4em;
  }
}

.hero__services {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px 14px;
  font-size: clamp(14px, 1.2vw, 18px);
  font-weight: 500;
  letter-spacing: 0.12em;
  color: var(--ink-2);
}

.hero__sep {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent);
}

.hero__lead {
  max-width: 560px;
  font-size: var(--fs-body-lg);
  line-height: 1.55;
  color: var(--ink-2);
}

.hero__actions {
  display: flex;
  gap: 12px;
}

.hero__status {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  margin-top: -12px;
  font-size: var(--fs-small);
  font-weight: 600;
  color: var(--ink-2);
}

.hero__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--success);
  box-shadow: 0 0 0 4px color-mix(in srgb, var(--success) 18%, transparent);
  animation: pulse 2.4s var(--ease) infinite;
}

@keyframes pulse {
  50% {
    box-shadow: 0 0 0 7px color-mix(in srgb, var(--success) 6%, transparent);
  }
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

/* ─── Moslashuvchanlik ───────────────────────────── */
@media (max-width: 1024px) {
  .hero {
    min-height: auto;
    align-items: flex-start;
    padding-top: 56px;
    padding-bottom: calc(var(--wave-h) + min(44vw, 360px));
  }

  .hero__skyline {
    width: min(92vw, 720px);
    right: -4vw;
    bottom: calc(var(--wave-h) * 0.45);
  }

  .hero__circuit--top {
    left: auto;
    right: -80px;
    width: 420px;
    opacity: 0.5;
  }

  .hero__circuit--bottom {
    display: none;
  }
}

@media (max-width: 560px) {
  .hero {
    --wave-h: 150px;
    padding-top: 32px;
    padding-bottom: calc(var(--wave-h) + 50vw);
  }

  .hero__brand {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;
  }

  .hero__mark {
    height: 84px;
  }

  .hero__title {
    font-size: clamp(40px, 12.4vw, 52px);
  }

  .hero__tagline {
    font-size: 11px;
    letter-spacing: 0.2em;
  }

  .hero__services {
    font-size: 12px;
    letter-spacing: 0.04em;
    gap: 4px 8px;
  }

  .hero__sep {
    width: 4px;
    height: 4px;
  }

  .hero__circuit--top {
    display: none;
  }

  .hero__skyline {
    width: 118vw;
    right: -6vw;
    bottom: calc(var(--wave-h) * 0.4);
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
