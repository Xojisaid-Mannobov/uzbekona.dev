<script setup lang="ts">
import { computed } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'

// Server qaytargan barcha maydon xatolari (bloklar ichidagilar ham) bitta ro'yxatda
const props = defineProps<{ errors: Record<string, string>; labels?: Record<string, string> }>()

const list = computed(() =>
  Object.entries(props.errors).map(([key, msg]) => {
    const [root, idx, ...rest] = key.split('.')
    const label = props.labels?.[root] ?? root
    const where =
      idx !== undefined && /^\d+$/.test(idx) ? `${label} #${Number(idx) + 1}${rest.length ? ` → ${rest.join('.')}` : ''}` : label
    return { key, where, msg }
  }),
)
</script>

<template>
  <div v-if="list.length" class="summary" role="alert">
    <p class="summary__title"><AppIcon name="alert" :size="18" /> Saqlab bo‘lmadi — quyidagilarni tekshiring:</p>
    <ul>
      <li v-for="e in list" :key="e.key">
        <strong>{{ e.where }}:</strong> {{ e.msg }}
      </li>
    </ul>
  </div>
</template>

<style scoped>
.summary {
  padding: 16px 18px;
  border-radius: 16px;
  background: color-mix(in srgb, var(--danger) 8%, var(--surface));
  border: 1px solid color-mix(in srgb, var(--danger) 30%, transparent);
  margin-bottom: 16px;
  font-size: 14px;
}

.summary__title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
  color: var(--danger);
}

ul {
  margin-top: 8px;
  padding-left: 26px;
  color: var(--ink-2);
}
</style>
