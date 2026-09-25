<script setup lang="ts">
import { computed } from 'vue'
import { leafDate } from './tone'

// Taqvim varag'i: tepada ikki "halqa" teshigi, katta kun raqami va oy — yangilik sanasi uchun
const props = withDefaults(defineProps<{ date: string | null; size?: 'sm' | 'md' | 'lg' }>(), { size: 'md' })
const d = computed(() => leafDate(props.date))
</script>

<template>
  <time class="leaf" :class="`leaf--${size}`" :datetime="date ?? undefined">
    <span class="leaf__top" aria-hidden="true"><i /><i /></span>
    <span class="leaf__day">{{ d.day }}</span>
    <span class="leaf__month">{{ d.month }} {{ size === 'sm' ? '' : d.year }}</span>
  </time>
</template>

<style scoped>
.leaf {
  --w: 64px;
  display: inline-grid;
  justify-items: center;
  align-content: start;
  width: var(--w);
  padding-bottom: 8px;
  border-radius: 12px;
  background: var(--surface);
  box-shadow:
    0 1px 0 var(--line),
    var(--shadow-sm);
  border: 1px solid var(--line);
  overflow: hidden;
  line-height: 1;
  flex-shrink: 0;
}

.leaf__top {
  display: flex;
  justify-content: center;
  gap: 18px;
  width: 100%;
  height: 12px;
  background: var(--leaf-accent, var(--accent));
  margin-bottom: 8px;
}

.leaf__top i {
  width: 5px;
  height: 5px;
  margin-top: 4px;
  border-radius: 50%;
  background: var(--surface);
}

.leaf__day {
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.05em;
  color: var(--ink);
}

.leaf__month {
  margin-top: 4px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--ink-3);
}

.leaf--sm {
  --w: 54px;
  border-radius: 10px;
}

.leaf--sm .leaf__day {
  font-size: 22px;
}

.leaf--lg {
  --w: 84px;
  border-radius: 14px;
  padding-bottom: 12px;
}

.leaf--lg .leaf__top {
  height: 16px;
  gap: 26px;
  margin-bottom: 10px;
}

.leaf--lg .leaf__top i {
  width: 6px;
  height: 6px;
  margin-top: 5px;
}

.leaf--lg .leaf__day {
  font-size: 36px;
}

.leaf--lg .leaf__month {
  font-size: 12px;
}
</style>
