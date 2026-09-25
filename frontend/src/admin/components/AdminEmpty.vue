<script setup lang="ts">
import AppIcon from '@/components/ui/AppIcon.vue'
import type { IconName } from '@/components/ui/icons'

withDefaults(defineProps<{ title: string; text?: string; icon?: IconName; error?: boolean }>(), { icon: 'layers' })
defineEmits<{ retry: [] }>()
</script>

<template>
  <div class="a-empty" :role="error ? 'alert' : undefined">
    <span class="a-empty__icon" :style="error ? { color: 'var(--danger)' } : undefined">
      <AppIcon :name="error ? 'alert' : icon" :size="24" />
    </span>
    <strong>{{ title }}</strong>
    <p v-if="text">{{ text }}</p>
    <button v-if="error" type="button" class="a-btn a-btn--outline a-btn--sm" @click="$emit('retry')">
      <AppIcon name="refresh" :size="16" /> Qayta urinish
    </button>
    <slot />
  </div>
</template>
