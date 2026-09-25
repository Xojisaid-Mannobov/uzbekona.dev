<script setup lang="ts">
import AppIcon from '@/components/ui/AppIcon.vue'
import type { Metric } from '@/types/api'

// Raqam + izoh juftliklari ("75K+" — "Foydalanuvchilar")
const model = defineModel<Metric[]>({ required: true })
withDefaults(defineProps<{ max?: number; errors?: Record<string, string>; prefix?: string }>(), { max: 8, prefix: 'metrics' })

function add() {
  model.value = [...model.value, { value: '', label: '' }]
}

function remove(i: number) {
  model.value = model.value.filter((_, idx) => idx !== i)
}
</script>

<template>
  <div class="metrics">
    <div v-for="(m, i) in model" :key="i" class="metrics__row">
      <input
        v-model="m.value"
        class="a-input metrics__value"
        placeholder="75K+"
        aria-label="Qiymat"
        :class="{ 'is-error': errors?.[`${prefix}.${i}.value`] }"
      />
      <input
        v-model="m.label"
        class="a-input"
        placeholder="Foydalanuvchilar"
        aria-label="Izoh"
        :class="{ 'is-error': errors?.[`${prefix}.${i}.label`] }"
      />
      <button type="button" class="a-icon-btn a-icon-btn--danger" aria-label="O‘chirish" @click="remove(i)">
        <AppIcon name="trash" :size="18" />
      </button>
    </div>
    <button v-if="model.length < max" type="button" class="a-btn a-btn--outline a-btn--sm metrics__add" @click="add">
      <AppIcon name="plus" :size="16" /> Qo‘shish
    </button>
  </div>
</template>

<style scoped>
.metrics {
  display: grid;
  gap: 8px;
}

.metrics__row {
  display: grid;
  grid-template-columns: 120px 1fr auto;
  gap: 8px;
  align-items: center;
}

.metrics__value {
  font-weight: 700;
}

.is-error {
  border-color: var(--danger);
}

.metrics__add {
  justify-self: start;
}
</style>
