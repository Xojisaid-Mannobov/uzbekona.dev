<script setup lang="ts">
import { computed, useId } from 'vue'
import { GIRIH_TILE } from './geometry'

/**
 * Girih naqshi — butun maydonni qoplovchi dekorativ fon (koshin/NBU uslubida).
 * `fade` — naqsh qaysi tomondan ko'rinib, qaysi tomonga so'nishi (matn o'qilishi uchun).
 */
const props = withDefaults(
  defineProps<{
    size?: number
    tone?: 'accent' | 'gold' | 'ink' | 'light'
    opacity?: number
    fade?: 'right' | 'left' | 'center' | 'top' | 'bottom' | 'none'
    /** Ixtiyoriy aniq rang (masalan loyiha accent rangi) — tone'dan ustun */
    color?: string
  }>(),
  { size: 96, tone: 'accent', opacity: 0.14, fade: 'right' },
)

const id = `girih-${useId()}`
const color = computed(
  () => props.color ?? { accent: 'var(--ornament)', gold: 'var(--gold)', ink: 'var(--ink)', light: 'var(--on-dark)' }[props.tone],
)
</script>

<template>
  <svg class="girih" :class="`girih--${fade}`" :style="{ color, opacity }" aria-hidden="true" focusable="false">
    <defs>
      <pattern :id="id" :width="size" :height="size" patternUnits="userSpaceOnUse" viewBox="0 0 100 100">
        <g fill="none" stroke="currentColor" stroke-width="1.25" stroke-linejoin="round">
          <path v-for="(d, i) in GIRIH_TILE" :key="i" :d="d" />
        </g>
      </pattern>
    </defs>
    <rect width="100%" height="100%" :fill="`url(#${id})`" />
  </svg>
</template>

<style scoped>
.girih {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  user-select: none;
}

.girih--right {
  mask-image: radial-gradient(ellipse 60% 95% at 100% 45%, #000 0%, rgb(0 0 0 / 0.5) 45%, transparent 78%);
}

.girih--left {
  mask-image: radial-gradient(ellipse 60% 95% at 0% 45%, #000 0%, rgb(0 0 0 / 0.5) 45%, transparent 78%);
}

.girih--center {
  mask-image: radial-gradient(ellipse 55% 70% at 50% 50%, #000 0%, transparent 80%);
}

.girih--top {
  mask-image: linear-gradient(180deg, #000 0%, transparent 85%);
}

.girih--bottom {
  mask-image: linear-gradient(0deg, #000 0%, transparent 85%);
}
</style>
