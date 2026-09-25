<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { vReveal } from '@/composables/reveal'
import { loadMotion, reducedMotion } from '@/composables/useSmoothScroll'

// Metrikalar admin paneldagi Sozlamalar'dan keladi (faqat real raqamlar)
const settings = useSettingsStore()
const root = ref<HTMLElement>()
let triggers: { kill(): void }[] = []

/** "75K+" → { prefix: '', num: 75, decimals: 0, suffix: 'K+' } */
function parse(value: string) {
  const m = value.match(/^(\D*)(\d+(?:[.,]\d+)?)(.*)$/)
  if (!m) return null
  const num = parseFloat(m[2].replace(',', '.'))
  const decimals = (m[2].split(/[.,]/)[1] ?? '').length
  return { prefix: m[1], num, decimals, suffix: m[3] }
}

// Raqamlar ekranga kirganda 0 dan sanab chiqadi
async function animate() {
  if (reducedMotion() || !root.value) return
  const { gsap, ScrollTrigger } = await loadMotion()
  root.value?.querySelectorAll<HTMLElement>('[data-value]').forEach((el) => {
    const p = parse(el.dataset.value ?? '')
    if (!p) return
    const counter = { v: 0 }
    const render = () => (el.textContent = `${p.prefix}${counter.v.toFixed(p.decimals)}${p.suffix}`)
    render()
    const t = gsap.to(counter, {
      v: p.num,
      duration: 1.6,
      ease: 'expo.out',
      onUpdate: render,
      paused: true,
    })
    triggers.push(ScrollTrigger.create({ trigger: el, start: 'top 90%', once: true, onEnter: () => t.play() }))
  })
}

watch(
  () => settings.metrics,
  async (m) => {
    if (!m.length) return
    await nextTick()
    animate()
  },
  { immediate: true },
)

onBeforeUnmount(() => triggers.forEach((t) => t.kill()))
</script>

<template>
  <section v-if="settings.metrics.length" ref="root" class="section metrics" aria-label="Raqamlarda">
    <div class="container">
      <p class="t-label metrics__label" v-reveal>Raqamlarda</p>
      <dl class="metrics__grid">
        <div v-for="(m, i) in settings.metrics" :key="m.label" class="metrics__item" v-reveal="i * 90">
          <dt class="metrics__desc">{{ m.label }}</dt>
          <dd class="metrics__value" :data-value="m.value">{{ m.value }}</dd>
        </div>
      </dl>
    </div>
  </section>
</template>

<style scoped>
.metrics__label {
  margin-bottom: 48px;
}

.metrics__grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 32px;
}

.metrics__item {
  display: flex;
  flex-direction: column-reverse;
  justify-content: flex-end;
  gap: 16px;
  padding-top: 32px;
  border-top: 1px solid var(--line-strong);
}

.metrics__value {
  font-size: var(--fs-metric);
  line-height: 0.9;
  font-weight: 600;
  letter-spacing: -0.055em;
  font-variant-numeric: tabular-nums;
}

.metrics__desc {
  font-size: clamp(16px, 1.25vw, 18px);
  color: var(--ink-2);
  font-weight: 500;
}

@media (max-width: 1024px) {
  .metrics__grid {
    grid-template-columns: repeat(2, 1fr);
    row-gap: 56px;
  }
}
</style>
