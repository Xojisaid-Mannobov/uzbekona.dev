<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { formatNumber } from '@/utils/format'

// Kunlik ko'rishlar va tashrif buyuruvchilar: bitta o'q (ikkalasi ham "soni"), chiziq + yengil maydon,
// kursor/klaviatura bilan kun tanlanadi — tooltip ikkala qiymatni ko'rsatadi.
type Day = { day: string; views: number; visitors: number }
const props = defineProps<{ days: Day[] }>()

const SERIES = [
  { key: 'views', label: 'Ko‘rishlar', color: 'var(--chart-1)' },
  { key: 'visitors', label: 'Tashrif buyuruvchilar', color: 'var(--chart-2)' },
] as const

const box = ref<HTMLElement>()
const width = ref(720)
const H = 260
const PAD = { top: 16, right: 16, bottom: 32, left: 48 }
let ro: ResizeObserver | undefined

onMounted(() => {
  ro = new ResizeObserver(([entry]) => (width.value = Math.max(280, entry.contentRect.width)))
  if (box.value) ro.observe(box.value)
})
onBeforeUnmount(() => ro?.disconnect())

// "Chiroyli" qadam (1, 1.5, 2, 2.5, 3, 4, 5, 6, 8 × 10ⁿ) — eng katta qiymat ustida ortiqcha bo'shliq qolmaydi
const scale = computed(() => {
  const max = Math.max(4, ...props.days.map((d) => d.views))
  const raw = max / 4
  const mag = 10 ** Math.floor(Math.log10(raw))
  const step = [1, 1.5, 2, 2.5, 3, 4, 5, 6, 8, 10].map((m) => m * mag).find((s) => s >= raw) ?? raw
  const count = Math.ceil(max / step)
  return { top: step * count, ticks: Array.from({ length: count + 1 }, (_, i) => i * step) }
})

const plotW = computed(() => width.value - PAD.left - PAD.right)
const plotH = H - PAD.top - PAD.bottom
const x = (i: number) => PAD.left + (props.days.length < 2 ? plotW.value / 2 : (i * plotW.value) / (props.days.length - 1))
const y = (v: number) => PAD.top + plotH - (v / scale.value.top) * plotH

function linePath(key: 'views' | 'visitors') {
  return props.days.map((d, i) => `${i ? 'L' : 'M'}${x(i).toFixed(1)} ${y(d[key]).toFixed(1)}`).join(' ')
}
const areaPath = computed(() => {
  const n = props.days.length
  if (!n) return ''
  return `${linePath('views')} L${x(n - 1).toFixed(1)} ${y(0)} L${x(0).toFixed(1)} ${y(0)} Z`
})

const fmt = formatNumber
const MONTHS = ['yan', 'fev', 'mar', 'apr', 'may', 'iyun', 'iyul', 'avg', 'sen', 'okt', 'noy', 'dek']
function dayLabel(iso: string, long = false) {
  const [yy, mm, dd] = iso.split('-').map(Number)
  return long ? `${dd}-${MONTHS[mm - 1]}, ${yy}` : `${dd} ${MONTHS[mm - 1]}`
}

// X o'qida ~6 ta yorliq — bir-birini bosmasligi uchun
const xLabels = computed(() => {
  const n = props.days.length
  const every = Math.max(1, Math.ceil(n / Math.max(2, Math.floor(plotW.value / 90))))
  return props.days.map((d, i) => ({ i, text: dayLabel(d.day) })).filter(({ i }) => i % every === 0 || i === n - 1)
})

// ─── Hover / fokus ──────────────────────────────────────
const active = ref<number | null>(null)
function onMove(e: PointerEvent) {
  const rect = (e.currentTarget as SVGElement).getBoundingClientRect()
  const px = e.clientX - rect.left - PAD.left
  const n = props.days.length
  active.value = Math.min(n - 1, Math.max(0, Math.round((px / plotW.value) * (n - 1))))
}
function onKey(e: KeyboardEvent) {
  const n = props.days.length
  if (e.key === 'ArrowRight' || e.key === 'ArrowLeft') {
    e.preventDefault()
    const cur = active.value ?? n - 1
    active.value = Math.min(n - 1, Math.max(0, cur + (e.key === 'ArrowRight' ? 1 : -1)))
  }
}
const tip = computed(() => {
  if (active.value === null) return null
  const d = props.days[active.value]
  const left = x(active.value)
  return { d, left, flip: left > width.value - 200 }
})
</script>

<template>
  <div class="trend">
    <ul class="trend__legend" role="list">
      <li v-for="s in SERIES" :key="s.key"><i :style="{ background: s.color }" />{{ s.label }}</li>
    </ul>

    <div ref="box" class="trend__plot">
      <svg
        :width="width"
        :height="H"
        role="img"
        :aria-label="`Kunlik ko‘rishlar va tashrif buyuruvchilar, ${days.length} kun. Chap va o‘ng strelkalar bilan kunlarni tanlang.`"
        tabindex="0"
        @pointermove="onMove"
        @pointerleave="active = null"
        @focus="active = active ?? days.length - 1"
        @blur="active = null"
        @keydown="onKey"
      >
        <g class="trend__grid">
          <g v-for="t in scale.ticks" :key="t">
            <line :x1="PAD.left" :x2="width - PAD.right" :y1="y(t)" :y2="y(t)" />
            <text :x="PAD.left - 10" :y="y(t)" dy="0.32em" text-anchor="end">{{ fmt(t) }}</text>
          </g>
          <text v-for="l in xLabels" :key="l.i" :x="x(l.i)" :y="H - 8" text-anchor="middle">{{ l.text }}</text>
        </g>

        <path :d="areaPath" class="trend__area" />
        <path v-for="s in SERIES" :key="s.key" :d="linePath(s.key)" class="trend__line" :style="{ stroke: s.color }" />

        <template v-if="tip">
          <line :x1="tip.left" :x2="tip.left" :y1="PAD.top" :y2="PAD.top + plotH" class="trend__cross" />
          <circle
            v-for="s in SERIES"
            :key="s.key"
            :cx="tip.left"
            :cy="y(tip.d[s.key])"
            r="4.5"
            class="trend__dot"
            :style="{ fill: s.color }"
          />
        </template>
      </svg>

      <div v-if="tip" class="trend__tip" :class="{ 'is-flip': tip.flip }" :style="{ left: `${tip.left}px` }" aria-live="polite">
        <p class="trend__tip-day">{{ dayLabel(tip.d.day, true) }}</p>
        <p v-for="s in SERIES" :key="s.key" class="trend__tip-row">
          <i :style="{ background: s.color }" />
          <strong>{{ fmt(tip.d[s.key]) }}</strong>
          <span>{{ s.label }}</span>
        </p>
      </div>
    </div>

    <details class="trend__table">
      <summary>Jadval ko‘rinishida</summary>
      <table class="a-table">
        <thead>
          <tr>
            <th>Kun</th>
            <th style="text-align: right">Ko‘rishlar</th>
            <th style="text-align: right">Tashrif buyuruvchilar</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="d in [...days].reverse()" :key="d.day">
            <td>{{ dayLabel(d.day, true) }}</td>
            <td class="num">{{ fmt(d.views) }}</td>
            <td class="num">{{ fmt(d.visitors) }}</td>
          </tr>
        </tbody>
      </table>
    </details>
  </div>
</template>

<style scoped>
.trend__legend {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 20px;
  margin-bottom: 12px;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink-2);
}

.trend__legend li {
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

/* Legenda belgisi chiziq shaklida — grafikdagi belgini takrorlaydi */
.trend__legend i,
.trend__tip-row i {
  width: 14px;
  height: 2px;
  border-radius: 2px;
}

.trend__plot {
  position: relative;
}

.trend__plot svg {
  display: block;
  overflow: visible;
  outline: none;
  touch-action: pan-y;
}

.trend__plot svg:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 4px;
  border-radius: 8px;
}

.trend__grid line {
  stroke: var(--line);
  stroke-width: 1;
}

.trend__grid text {
  font-size: 11px;
  fill: var(--ink-3);
  font-variant-numeric: tabular-nums;
}

.trend__area {
  fill: var(--chart-1);
  opacity: 0.1;
}

.trend__line {
  fill: none;
  stroke-width: 2;
  stroke-linejoin: round;
  stroke-linecap: round;
}

.trend__cross {
  stroke: var(--ink-3);
  stroke-width: 1;
}

.trend__dot {
  stroke: var(--surface);
  stroke-width: 2;
}

.trend__tip {
  position: absolute;
  top: 0;
  min-width: 190px;
  margin-left: 14px;
  padding: 10px 12px;
  border-radius: 12px;
  background: var(--surface);
  border: 1px solid var(--line-strong);
  box-shadow: var(--shadow-md);
  pointer-events: none;
}

.trend__tip.is-flip {
  transform: translateX(calc(-100% - 28px));
}

.trend__tip-day {
  margin-bottom: 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--ink-3);
}

.trend__tip-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.trend__tip-row strong {
  min-width: 36px;
  font-size: 15px;
  color: var(--ink);
}

.trend__tip-row span {
  color: var(--ink-2);
}

.trend__table {
  margin-top: 14px;
}

.trend__table summary {
  cursor: pointer;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink-2);
}

.trend__table table {
  margin-top: 10px;
}

.num {
  text-align: right;
  font-variant-numeric: tabular-nums;
}
</style>
