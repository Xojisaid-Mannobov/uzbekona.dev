<script setup lang="ts">
import { computed } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'
import type { PageMeta } from '@/types/api'

const props = defineProps<{ meta: PageMeta }>()
const page = defineModel<number>({ required: true })

const pages = computed(() => Math.max(1, Math.ceil(props.meta.total / props.meta.limit)))
const from = computed(() => (props.meta.total ? (props.meta.page - 1) * props.meta.limit + 1 : 0))
const to = computed(() => Math.min(props.meta.page * props.meta.limit, props.meta.total))
</script>

<template>
  <div v-if="meta.total > meta.limit" class="a-pagination">
    <span>{{ from }}–{{ to }} / {{ meta.total }}</span>
    <div class="a-pagination__btns">
      <button type="button" class="a-btn a-btn--outline a-btn--sm" :disabled="page <= 1" @click="page--">
        <AppIcon name="arrow-left" :size="16" /> Oldingi
      </button>
      <button type="button" class="a-btn a-btn--outline a-btn--sm" :disabled="page >= pages" @click="page++">
        Keyingi <AppIcon name="arrow-right" :size="16" />
      </button>
    </div>
  </div>
</template>
