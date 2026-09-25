<script setup lang="ts">
import { computed, useId } from 'vue'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import UiButton from '@/components/ui/UiButton.vue'
import GirihPattern from '@/components/ornament/GirihPattern.vue'
import SuzaniRosette from '@/components/ornament/SuzaniRosette.vue'
import OrnamentBorder from '@/components/ornament/OrnamentBorder.vue'
import { whyUs } from '@/content/site'
import { useSettingsStore } from '@/stores/settings'
import { vReveal } from '@/composables/reveal'
import { pad2 } from '@/utils/format'

// "Nega aynan biz" — bento kartalar: har birida o'z vizuali (suzani, o'sish, sprint, kalit, monitoring, ishonch)
withDefaults(defineProps<{ index?: string }>(), { index: '05' })

const settings = useSettingsStore()
const card = (key: (typeof whyUs.cards)[number]['key']) => whyUs.cards.find((c) => c.key === key)!
const cards = {
  care: card('care'),
  scale: card('scale'),
  sprint: card('sprint'),
  ownership: card('ownership'),
  support: card('support'),
  trust: card('trust'),
}
const num = (key: string) => pad2(whyUs.cards.findIndex((c) => c.key === key) + 1)

// "Katta miqyos" kartasidagi raqam — admin sozlamalaridagi birinchi metrika (bo'lmasa ko'rsatilmaydi)
const metric = computed(() => settings.metrics[0] ?? null)

// 8 qirrali yulduz (kalit boshi uchun): tashqi va ichki radiuslar almashadi
function star8(cx: number, cy: number, r: number) {
  const pts: string[] = []
  for (let i = 0; i < 16; i++) {
    const rad = i % 2 ? r * 0.765 : r
    const a = (Math.PI / 8) * i - Math.PI / 2
    pts.push(`${(cx + rad * Math.cos(a)).toFixed(1)} ${(cy + rad * Math.sin(a)).toFixed(1)}`)
  }
  return `M${pts.join(' L')} Z`
}
const KEY_BOW = star8(56, 60, 42)
const KEY_INNER = star8(56, 60, 22)

const fillId = `why-fill-${useId()}`

const BARS = [22, 30, 28, 40, 46, 44, 58, 70, 92]
const SPRINTS = [16, 83, 150, 217, 284]
</script>

<template>
  <section class="section why" aria-labelledby="why-title">
    <div class="container">
      <SectionHeader title-id="why-title" :index="index" label="Nega biz" :title="whyUs.title" :lead="whyUs.lead" />

      <div class="bento">
        <!-- 01 · Mehr -->
        <article class="tile tile--care" v-reveal>
          <GirihPattern tone="light" :size="72" :opacity="0.1" fade="right" />
          <div class="tile__rosette" aria-hidden="true"><SuzaniRosette tone="light" :opacity="0.5" spin /></div>
          <span class="tile__num">{{ num('care') }}</span>
          <div class="tile__body">
            <h3 class="tile__title tile__title--xl">{{ cards.care.title }}</h3>
            <p class="tile__text">{{ cards.care.text }}</p>
          </div>
        </article>

        <!-- 02 · Miqyos -->
        <article class="tile tile--scale" v-reveal="100">
          <span class="tile__num">{{ num('scale') }}</span>
          <div class="tile__body">
            <h3 class="tile__title">{{ cards.scale.title }}</h3>
            <p class="tile__text">{{ cards.scale.text }}</p>
          </div>
          <div class="scale">
            <p v-if="metric" class="scale__metric">
              <strong>{{ metric.value }}</strong>
              <span>{{ metric.label }}</span>
            </p>
            <svg class="scale__bars" viewBox="0 0 220 100" aria-hidden="true" focusable="false">
              <rect
                v-for="(h, i) in BARS"
                :key="i"
                :x="i * 24 + 4"
                :y="100 - h"
                width="16"
                :height="h"
                rx="4"
                :class="i === BARS.length - 1 ? 'bar bar--top' : 'bar'"
                :style="{ '--i': i }"
              />
            </svg>
          </div>
        </article>

        <!-- 03 · Sprint -->
        <article class="tile" v-reveal>
          <span class="tile__num">{{ num('sprint') }}</span>
          <svg class="sprint" viewBox="0 0 300 64" aria-hidden="true" focusable="false">
            <line x1="16" y1="22" x2="284" y2="22" class="sprint__track" />
            <line x1="16" y1="22" x2="217" y2="22" class="sprint__done" />
            <g v-for="(x, i) in SPRINTS" :key="x">
              <template v-if="i < SPRINTS.length - 1">
                <circle :cx="x" cy="22" r="11" class="sprint__node" />
                <path :d="`M${x - 4.5} 22 l3 3 l6 -6.5`" class="sprint__check" />
              </template>
              <template v-else>
                <circle :cx="x" cy="22" r="16" class="sprint__pulse" />
                <circle :cx="x" cy="22" r="10" class="sprint__now" />
              </template>
              <text :x="x" y="58" text-anchor="middle" class="sprint__label">{{ (i + 1) * 2 }}-hafta</text>
            </g>
          </svg>
          <div class="tile__body">
            <h3 class="tile__title">{{ cards.sprint.title }}</h3>
            <p class="tile__text">{{ cards.sprint.text }}</p>
          </div>
        </article>

        <!-- 04 · Egalik -->
        <article class="tile tile--ownership" v-reveal="100">
          <span class="tile__num">{{ num('ownership') }}</span>
          <svg class="key" viewBox="0 0 220 120" aria-hidden="true" focusable="false">
            <g transform="rotate(-14 110 60)">
              <path d="M92 53 H196 a7 7 0 0 1 0 14 H186 V86 H172 V67 H160 V80 H146 V67 H92 Z" class="key__metal" />
              <path :d="KEY_BOW" class="key__metal" />
              <path :d="KEY_INNER" class="key__inner" />
              <circle cx="56" cy="60" r="9" class="key__hole" />
            </g>
          </svg>
          <div class="tile__body">
            <h3 class="tile__title">{{ cards.ownership.title }}</h3>
            <p class="tile__text">{{ cards.ownership.text }}</p>
          </div>
        </article>

        <!-- 05 · Qo'llab-quvvatlash -->
        <article class="tile" v-reveal="200">
          <span class="tile__num">{{ num('support') }}</span>
          <div class="monitor" aria-hidden="true">
            <p class="monitor__head"><span class="monitor__dot" /> Monitoring 24/7</p>
            <svg class="monitor__chart" viewBox="0 0 300 80" preserveAspectRatio="none" focusable="false">
              <defs>
                <linearGradient :id="fillId" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0" class="monitor__stop" stop-opacity="0.28" />
                  <stop offset="1" class="monitor__stop" stop-opacity="0" />
                </linearGradient>
              </defs>
              <path
                d="M0 34 L20 30 L40 33 L60 27 L80 31 L100 26 L120 29 L140 24 L160 28 L180 22 L200 25 L220 20 L240 23 L260 18 L280 21 L300 16 L300 80 L0 80 Z"
                :fill="`url(#${fillId})`"
              />
              <path
                d="M0 34 L20 30 L40 33 L60 27 L80 31 L100 26 L120 29 L140 24 L160 28 L180 22 L200 25 L220 20 L240 23 L260 18 L280 21 L300 16"
                class="monitor__line"
              />
            </svg>
          </div>
          <div class="tile__body">
            <h3 class="tile__title">{{ cards.support.title }}</h3>
            <p class="tile__text">{{ cards.support.text }}</p>
          </div>
        </article>

        <!-- 06 · Ishonch -->
        <article class="tile tile--trust" v-reveal>
          <OrnamentBorder class="tile__hoshiya" :opacity="0.5" />
          <div class="tile__rosette tile__rosette--gold" aria-hidden="true"><SuzaniRosette tone="gold" :opacity="0.75" /></div>
          <span class="tile__num">{{ num('trust') }}</span>
          <div class="tile__body trust">
            <h3 class="tile__title tile__title--xl">{{ cards.trust.title }}</h3>
            <p class="tile__text">{{ cards.trust.text }}</p>
            <UiButton to="/contact" icon="arrow-up-right">Loyihani muhokama qilish</UiButton>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped>
.bento {
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  gap: clamp(14px, 1.4vw, 20px);
}

.tile {
  position: relative;
  grid-column: span 4;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 32px;
  min-height: 380px;
  padding: clamp(26px, 2.6vw, 40px);
  border-radius: var(--r-lg);
  border: 1px solid var(--line);
  background: var(--surface);
  overflow: hidden;
  isolation: isolate;
  transition:
    transform var(--dur) var(--ease),
    box-shadow var(--dur) var(--ease);
}

.tile:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.tile__num {
  font-size: var(--fs-small);
  font-weight: 700;
  color: var(--ink-3);
}

.tile__body {
  position: relative;
  display: grid;
  gap: 14px;
}

.tile__title {
  font-size: var(--fs-card);
  font-weight: 600;
  line-height: 1.15;
  letter-spacing: -0.03em;
  max-width: 22ch;
}

.tile__title--xl {
  font-size: clamp(28px, 2.7vw, 42px);
  line-height: 1.08;
  letter-spacing: -0.04em;
  max-width: 17ch;
}

.tile__text {
  max-width: 44ch;
  color: var(--ink-2);
  line-height: 1.55;
}

/* 01 — accent karta, oq suzani */
.tile--care {
  grid-column: span 7;
  min-height: clamp(420px, 36vw, 520px);
  border: 0;
  color: #fff;
  background:
    radial-gradient(90% 90% at 100% 100%, rgb(255 255 255 / 0.14), transparent 60%),
    linear-gradient(140deg, color-mix(in srgb, var(--accent) 90%, #0b1b5c) 0%, color-mix(in srgb, var(--accent) 62%, #0b1b5c) 100%);
}

.tile--care .tile__num,
.tile--care .tile__text {
  color: rgb(255 255 255 / 0.78);
}

.tile--care .tile__text {
  max-width: 48ch;
  font-size: var(--fs-body-lg);
}

.tile__rosette {
  position: absolute;
  z-index: -1;
  right: -14%;
  top: -18%;
  width: min(62%, 420px);
  aspect-ratio: 1;
}

/* 02 — o'sish ustunlari */
.tile--scale {
  grid-column: span 5;
  min-height: clamp(420px, 36vw, 520px);
}

.tile--scale {
  justify-content: flex-start;
}

.scale {
  margin-top: auto;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 20px;
}

.scale__metric {
  display: grid;
  gap: 4px;
}

.scale__metric strong {
  font-size: var(--fs-metric);
  font-weight: 700;
  line-height: 0.95;
  letter-spacing: -0.05em;
  color: var(--accent);
}

.scale__metric span {
  font-size: var(--fs-small);
  font-weight: 600;
  color: var(--ink-3);
}

.scale__bars {
  width: min(52%, 220px);
  flex-shrink: 0;
}

.bar {
  fill: color-mix(in srgb, var(--accent) 18%, transparent);
  transform-box: fill-box;
  transform-origin: bottom;
}

.bar--top {
  fill: var(--accent);
}

.tile--scale:hover .bar {
  animation: bar-grow 700ms var(--ease) both;
  animation-delay: calc(var(--i) * 40ms);
}

@keyframes bar-grow {
  from {
    transform: scaleY(0.2);
  }
}

/* 03 — sprint chizig'i */
.sprint {
  width: 100%;
  overflow: visible;
}

.sprint__track {
  stroke: var(--line-strong);
  stroke-width: 3;
  stroke-linecap: round;
}

.sprint__done {
  stroke: var(--accent);
  stroke-width: 3;
  stroke-linecap: round;
}

.sprint__node {
  fill: var(--accent);
}

.sprint__check {
  fill: none;
  stroke: #fff;
  stroke-width: 2.2;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.sprint__now {
  fill: var(--surface);
  stroke: var(--accent);
  stroke-width: 3;
}

.sprint__pulse {
  fill: color-mix(in srgb, var(--accent) 16%, transparent);
  transform-box: fill-box;
  transform-origin: center;
  animation: sprint-pulse 2.2s var(--ease) infinite;
}

@keyframes sprint-pulse {
  50% {
    transform: scale(1.35);
    opacity: 0.4;
  }
}

.sprint__label {
  font-size: 11px;
  font-weight: 600;
  fill: var(--ink-3);
}

/* 04 — oltin kalit (boshi 8 qirrali yulduz) */
.tile--ownership {
  background: radial-gradient(80% 70% at 80% 10%, var(--gold-soft), transparent 70%), var(--surface);
}

.key {
  width: min(100%, 250px);
  overflow: visible;
}

.key__metal {
  fill: var(--gold);
}

.key__inner {
  fill: none;
  stroke: var(--surface);
  stroke-width: 2;
  opacity: 0.7;
}

.key__hole {
  fill: var(--surface);
}

.tile--ownership:hover .key {
  transform: rotate(8deg);
}

.key {
  transition: transform var(--dur-slow) var(--ease);
}

/* 05 — monitoring grafigi */
.monitor {
  display: grid;
  gap: 14px;
}

.monitor__head {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  font-size: var(--fs-small);
  font-weight: 600;
  color: var(--ink-2);
}

.monitor__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--success);
  box-shadow: 0 0 0 4px color-mix(in srgb, var(--success) 18%, transparent);
}

.monitor__chart {
  width: 100%;
  height: 80px;
}

.monitor__stop {
  stop-color: var(--success);
}

.monitor__line {
  fill: none;
  stroke: var(--success);
  stroke-width: 2;
  vector-effect: non-scaling-stroke;
  stroke-linejoin: round;
}

/* 06 — ishonch: oltin hoshiya va suzani */
.tile--trust {
  grid-column: span 12;
  min-height: 0;
  padding-top: calc(clamp(26px, 2.6vw, 40px) + 28px);
  background: color-mix(in srgb, var(--gold) 9%, var(--surface));
  border-color: color-mix(in srgb, var(--gold) 26%, transparent);
}

.tile__hoshiya {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
}

.tile__rosette--gold {
  top: 50%;
  right: clamp(-120px, -6vw, -40px);
  width: clamp(240px, 26vw, 380px);
  translate: 0 -50%;
}

.trust {
  justify-items: start;
  gap: 18px;
  max-width: 760px;
}

.trust .tile__title--xl {
  max-width: 22ch;
}

.trust .tile__text {
  max-width: 58ch;
  font-size: var(--fs-body-lg);
  margin-bottom: 10px;
}

@media (max-width: 1024px) {
  .tile,
  .tile--care,
  .tile--scale {
    grid-column: span 6;
  }

  .tile--care {
    grid-column: span 12;
  }

  .tile--scale {
    min-height: 380px;
  }

  .tile:nth-child(5) {
    grid-column: span 12;
  }

  .tile--trust {
    grid-column: span 12;
  }
}

@media (max-width: 680px) {
  .tile,
  .tile--care,
  .tile--scale,
  .tile:nth-child(5) {
    grid-column: span 12;
    min-height: 0;
  }

  .tile--care {
    min-height: 420px;
  }

  .tile__rosette {
    width: 78%;
    right: -26%;
    top: -12%;
  }

  .tile__rosette--gold {
    width: 200px;
    right: -80px;
    top: auto;
    bottom: -60px;
    translate: none;
    opacity: 0.6;
  }
}

@media (prefers-reduced-motion: reduce) {
  .sprint__pulse {
    animation: none;
  }
}
</style>
