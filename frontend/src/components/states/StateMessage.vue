<script setup lang="ts">
import AppIcon from '@/components/ui/AppIcon.vue'
import type { IconName } from '@/components/ui/icons'
import UiButton from '@/components/ui/UiButton.vue'

/** Xato yoki bo'sh holat uchun yagona, sokin ko'rinish. */
withDefaults(
  defineProps<{
    tone?: 'error' | 'empty'
    title: string
    text?: string
    icon?: IconName
    actionLabel?: string
  }>(),
  { tone: 'empty' },
)

defineEmits<{ action: [] }>()
</script>

<template>
  <div class="state" :class="`state--${tone}`" :role="tone === 'error' ? 'alert' : undefined">
    <span class="state__icon">
      <AppIcon :name="icon ?? (tone === 'error' ? 'alert' : 'layers')" :size="28" />
    </span>
    <h3 class="state__title">{{ title }}</h3>
    <p v-if="text" class="state__text">{{ text }}</p>
    <UiButton
      v-if="actionLabel"
      variant="secondary"
      size="md"
      :icon-left="tone === 'error' ? 'refresh' : undefined"
      @click="$emit('action')"
    >
      {{ actionLabel }}
    </UiButton>
    <slot />
  </div>
</template>

<style scoped>
.state {
  display: grid;
  justify-items: center;
  text-align: center;
  gap: 16px;
  padding: clamp(56px, 8vw, 112px) 24px;
  border: 1px dashed var(--line-strong);
  border-radius: var(--r-lg);
}

.state__icon {
  display: grid;
  place-items: center;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: var(--surface-2);
  color: var(--ink-2);
  margin-bottom: 8px;
}

.state--error .state__icon {
  background: color-mix(in srgb, var(--danger) 12%, transparent);
  color: var(--danger);
}

.state__title {
  font-size: clamp(24px, 2vw, 30px);
  letter-spacing: -0.025em;
}

.state__text {
  max-width: 46ch;
  color: var(--ink-2);
  margin-bottom: 8px;
}
</style>
