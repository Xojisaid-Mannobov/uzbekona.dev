<script setup lang="ts">
import { computed } from 'vue'
import type { MediaRef } from '@/types/api'
import MediaImage from './MediaImage.vue'

/**
 * Loyiha vizuali. Cover yuklangan bo'lsa — real screenshot.
 * Hali screenshot bo'lmasa — loyiha accent rangidagi toza interfeys eskizi
 * (soxta ma'lumotsiz: faqat struktura va loyiha nomi).
 */
const props = withDefaults(
  defineProps<{
    title: string
    subtitle?: string
    accent?: string
    cover?: MediaRef | null
    variant?: 'auto' | 'dashboard' | 'landing' | 'mobile'
    seed?: string
    sizes?: string
    priority?: boolean
  }>(),
  { variant: 'auto', sizes: '(min-width: 1440px) 1360px, 100vw' },
)

const accent = computed(() => (props.accent && /^#[0-9a-f]{6}$/i.test(props.accent) ? props.accent : '#2350F5'))

const kind = computed(() => {
  if (props.variant !== 'auto') return props.variant
  const key = props.seed ?? props.title
  let hash = 0
  for (const ch of key) hash = (hash * 31 + ch.charCodeAt(0)) >>> 0
  return (['dashboard', 'landing', 'mobile'] as const)[hash % 3]
})
</script>

<template>
  <div class="visual" :style="{ '--pa': accent }">
    <MediaImage v-if="cover" :media="cover" :sizes="sizes" :priority="priority" fill :alt="cover.alt || title" />

    <div v-else class="mock" :class="`mock--${kind}`" role="img" :aria-label="`${title} — interfeys eskizi`">
      <!-- Dashboard: sidebar + KPI + grafik -->
      <div v-if="kind === 'dashboard'" class="win">
        <div class="win__bar">
          <i /><i /><i /><span class="win__url">{{ title.toLowerCase() }}</span>
        </div>
        <div class="dash">
          <aside class="dash__side">
            <b class="logo">{{ title.charAt(0) }}</b>
            <span v-for="n in 6" :key="n" class="line" :class="{ 'line--active': n === 2 }" />
          </aside>
          <div class="dash__main">
            <div class="dash__head">
              <strong>{{ title }}</strong>
              <span class="pill" />
            </div>
            <div class="dash__kpis">
              <div v-for="n in 3" :key="n" class="card">
                <span class="line line--xs" />
                <span class="bar" :style="{ width: `${48 + n * 12}%` }" />
              </div>
            </div>
            <div class="dash__row">
              <div class="card card--chart">
                <span class="line line--xs" />
                <svg viewBox="0 0 300 100" preserveAspectRatio="none" aria-hidden="true">
                  <path class="area" d="M0 80 C 30 70, 50 76, 75 60 S 120 40, 150 48 S 200 22, 225 30 S 270 12, 300 8 V100 H0 Z" />
                  <path class="stroke" d="M0 80 C 30 70, 50 76, 75 60 S 120 40, 150 48 S 200 22, 225 30 S 270 12, 300 8" />
                </svg>
              </div>
              <div class="card card--list">
                <div v-for="n in 5" :key="n" class="row"><i /><span class="line" /></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Landing: katta sarlavha + CTA + kontent bloklari -->
      <div v-else-if="kind === 'landing'" class="win">
        <div class="win__bar">
          <i /><i /><i /><span class="win__url">{{ title.toLowerCase() }}</span>
        </div>
        <div class="land">
          <div class="land__nav">
            <b>{{ title }}</b>
            <span class="line" /><span class="line" /><span class="line" />
            <span class="pill pill--solid" />
          </div>
          <p class="land__title">{{ subtitle || title }}</p>
          <div class="land__cta"><span class="pill pill--solid" /><span class="pill" /></div>
          <div class="land__grid">
            <div class="tile tile--a" />
            <div class="tile" />
            <div class="tile" />
          </div>
        </div>
      </div>

      <!-- Mobile: ikki telefon ekrani -->
      <div v-else class="phones">
        <div v-for="n in 2" :key="n" class="phone" :class="`phone--${n}`">
          <div class="phone__screen">
            <div class="phone__notch" />
            <strong v-if="n === 1">{{ title }}</strong>
            <span v-else class="line line--xs" />
            <div class="phone__hero" />
            <div v-for="r in 4" :key="r" class="row"><i /><span class="line" /></div>
            <div class="phone__tab"><i v-for="t in 4" :key="t" :class="{ on: t === 1 }" /></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.visual {
  --tint: color-mix(in srgb, var(--pa) 14%, var(--surface));
  --tint-2: color-mix(in srgb, var(--pa) 26%, var(--surface));
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background:
    radial-gradient(120% 90% at 85% 0%, color-mix(in srgb, var(--pa) 30%, transparent) 0%, transparent 60%),
    linear-gradient(160deg, var(--tint) 0%, var(--tint-2) 100%);
  container-type: size;
}

.mock {
  position: absolute;
  inset: 0;
  color: #111312;
}

/* ─── Brauzer oynasi ─────────────────────────────── */
.win {
  position: absolute;
  left: 7%;
  right: 7%;
  top: 11%;
  bottom: -6%;
  background: #fbfbf9;
  border-radius: clamp(10px, 1.6cqw, 22px);
  box-shadow:
    0 2px 6px rgb(17 19 18 / 0.06),
    0 30px 80px color-mix(in srgb, var(--pa) 30%, rgb(17 19 18 / 0.25));
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.win__bar {
  display: flex;
  align-items: center;
  gap: 0.6cqw;
  padding: 1.1cqw 1.4cqw;
  border-bottom: 1px solid rgb(17 19 18 / 0.07);
  background: #f3f3f0;
}

.win__bar i {
  width: max(6px, 0.8cqw);
  aspect-ratio: 1;
  border-radius: 50%;
  background: rgb(17 19 18 / 0.14);
}

.win__url {
  margin-inline: auto;
  padding: 0.3cqw 2.4cqw;
  border-radius: 99px;
  background: rgb(17 19 18 / 0.05);
  font-size: max(9px, 0.95cqw);
  color: rgb(17 19 18 / 0.68);
  font-weight: 600;
}

.line {
  display: block;
  height: max(5px, 0.7cqw);
  border-radius: 99px;
  background: rgb(17 19 18 / 0.09);
}

.line--xs {
  width: 40%;
  height: max(4px, 0.55cqw);
}

.line--active {
  background: var(--pa);
  opacity: 0.85;
}

.pill {
  display: inline-block;
  width: 9cqw;
  height: 2.6cqw;
  border-radius: 99px;
  border: 1px solid rgb(17 19 18 / 0.14);
}

.pill--solid {
  background: var(--pa);
  border-color: transparent;
}

.card {
  background: #fff;
  border: 1px solid rgb(17 19 18 / 0.06);
  border-radius: clamp(6px, 1.1cqw, 16px);
  padding: 1.4cqw;
  display: grid;
  gap: 1cqw;
  align-content: start;
}

.bar {
  height: 2.2cqw;
  border-radius: 0.6cqw;
  background: #111312;
}

.row {
  display: flex;
  align-items: center;
  gap: 1cqw;
}

.row i {
  width: 2.2cqw;
  aspect-ratio: 1;
  border-radius: 30%;
  background: color-mix(in srgb, var(--pa) 22%, #fff);
  flex-shrink: 0;
}

.row .line {
  flex: 1;
}

/* ─── Dashboard ──────────────────────────────────── */
.dash {
  flex: 1;
  display: grid;
  grid-template-columns: 17% 1fr;
  min-height: 0;
}

.dash__side {
  border-right: 1px solid rgb(17 19 18 / 0.07);
  padding: 1.8cqw 1.4cqw;
  display: grid;
  gap: 1.4cqw;
  align-content: start;
}

.logo {
  width: 3.4cqw;
  aspect-ratio: 1;
  border-radius: 0.9cqw;
  background: var(--pa);
  color: #fff;
  display: grid;
  place-items: center;
  font-size: 1.6cqw;
  margin-bottom: 1.2cqw;
}

.dash__main {
  padding: 2cqw;
  display: grid;
  gap: 1.6cqw;
  grid-template-rows: auto auto 1fr;
  min-height: 0;
}

.dash__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.dash__head strong {
  font-size: 2.2cqw;
  letter-spacing: -0.03em;
}

.dash__kpis {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.4cqw;
}

.dash__row {
  display: grid;
  grid-template-columns: 1.6fr 1fr;
  gap: 1.4cqw;
  min-height: 0;
}

.card--chart svg {
  width: 100%;
  height: 100%;
  min-height: 10cqw;
}

.area {
  fill: color-mix(in srgb, var(--pa) 16%, transparent);
}

.stroke {
  fill: none;
  stroke: var(--pa);
  stroke-width: 2.5;
  vector-effect: non-scaling-stroke;
}

.card--list {
  gap: 1.4cqw;
}

/* ─── Landing ────────────────────────────────────── */
.land {
  flex: 1;
  padding: 2.4cqw 4cqw;
  display: flex;
  flex-direction: column;
  gap: 2.6cqw;
}

.land__nav {
  display: flex;
  align-items: center;
  gap: 2.4cqw;
}

.land__nav b {
  font-size: 1.6cqw;
  margin-right: auto;
  letter-spacing: -0.02em;
}

.land__nav .line {
  width: 5cqw;
}

.land__nav .pill {
  width: 8cqw;
}

.land__title {
  font-size: 5.2cqw;
  line-height: 0.98;
  letter-spacing: -0.045em;
  font-weight: 600;
  max-width: 70%;
}

.land__cta {
  display: flex;
  gap: 1.2cqw;
}

.land__cta .pill {
  width: 12cqw;
  height: 3.4cqw;
}

.land__grid {
  flex: 1;
  display: grid;
  grid-template-columns: 1.4fr 1fr 1fr;
  gap: 1.4cqw;
}

.tile {
  border-radius: clamp(8px, 1.4cqw, 18px);
  background: color-mix(in srgb, var(--pa) 12%, #f1f1ee);
}

.tile--a {
  background: linear-gradient(140deg, var(--pa), color-mix(in srgb, var(--pa) 55%, #111312));
}

/* ─── Mobile ─────────────────────────────────────── */
.phones {
  position: absolute;
  inset: 0;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  gap: 5cqw;
  padding-top: 9cqh;
}

.phone {
  width: min(24cqw, 44cqh);
  aspect-ratio: 9 / 19;
  border-radius: min(4.2cqw, 7cqh);
  background: #111312;
  padding: min(0.9cqw, 1.5cqh);
  box-shadow: 0 40px 80px color-mix(in srgb, var(--pa) 35%, rgb(17 19 18 / 0.3));
}

.phone--2 {
  margin-top: 8cqh;
}

.phone__screen {
  position: relative;
  height: 100%;
  border-radius: min(3.4cqw, 5.8cqh);
  background: #fbfbf9;
  padding: min(5cqw, 8cqh) min(1.8cqw, 3cqh) 0;
  display: flex;
  flex-direction: column;
  gap: min(1.4cqw, 2.4cqh);
  overflow: hidden;
}

.phone__screen strong {
  font-size: min(2cqw, 3.4cqh);
  letter-spacing: -0.03em;
}

.phone__notch {
  position: absolute;
  top: min(1cqw, 1.6cqh);
  left: 50%;
  translate: -50% 0;
  width: 30%;
  height: min(1.6cqw, 2.6cqh);
  border-radius: 99px;
  background: #111312;
}

.phone__hero {
  height: 28%;
  border-radius: min(1.6cqw, 2.6cqh);
  background: linear-gradient(140deg, var(--pa), color-mix(in srgb, var(--pa) 50%, #111312));
  flex-shrink: 0;
}

.phone .row i {
  width: min(2.6cqw, 4.2cqh);
}

.phone .line {
  height: min(0.8cqw, 1.3cqh);
}

.phone__tab {
  margin-top: auto;
  display: flex;
  justify-content: space-around;
  padding: min(1.4cqw, 2.4cqh) 0 min(2cqw, 3.2cqh);
  border-top: 1px solid rgb(17 19 18 / 0.07);
}

.phone__tab i {
  width: min(1.8cqw, 3cqh);
  aspect-ratio: 1;
  border-radius: 30%;
  background: rgb(17 19 18 / 0.14);
}

.phone__tab i.on {
  background: var(--pa);
}
</style>
