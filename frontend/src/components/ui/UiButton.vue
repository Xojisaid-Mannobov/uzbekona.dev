<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, type RouteLocationRaw } from 'vue-router'
import AppIcon from './AppIcon.vue'
import type { IconName } from './icons'

/**
 * Asosiy tugma. TZ bo'yicha balandlik 56px (md — 52px), to'liq yumaloq (999px).
 * `to` berilsa RouterLink, `href` berilsa <a>, aks holda <button>.
 */
const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'secondary' | 'ghost' | 'light' | 'outline-light'
    size?: 'lg' | 'md'
    to?: RouteLocationRaw
    href?: string
    icon?: IconName
    iconLeft?: IconName
    loading?: boolean
    disabled?: boolean
    type?: 'button' | 'submit'
    block?: boolean
  }>(),
  { variant: 'primary', size: 'lg', type: 'button' },
)

const tag = computed(() => (props.to ? RouterLink : props.href ? 'a' : 'button'))
const external = computed(() => !!props.href && /^https?:\/\//.test(props.href))
</script>

<template>
  <component
    :is="tag"
    class="btn"
    :class="[`btn--${variant}`, `btn--${size}`, { 'btn--block': block, 'is-loading': loading }]"
    v-bind="to ? { to } : href ? { href } : {}"
    :type="tag === 'button' ? type : undefined"
    :disabled="tag === 'button' ? disabled || loading : undefined"
    :aria-busy="loading || undefined"
    :target="external ? '_blank' : undefined"
    :rel="external ? 'noopener noreferrer' : undefined"
  >
    <AppIcon v-if="iconLeft" :name="iconLeft" :size="18" />
    <slot name="lead" />
    <span class="btn__label"><slot /></span>
    <span v-if="icon" class="btn__icon"><AppIcon :name="icon" :size="18" /></span>
    <span v-if="loading" class="btn__spinner" aria-hidden="true" />
  </component>
</template>

<style scoped>
.btn {
  --h: var(--btn-h);
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  height: var(--h);
  padding: 0 28px;
  border-radius: var(--r-pill);
  font-size: 16px;
  font-weight: 600;
  letter-spacing: -0.01em;
  white-space: nowrap;
  border: 1px solid transparent;
  transition:
    background-color var(--dur-fast) var(--ease),
    color var(--dur-fast) var(--ease),
    border-color var(--dur-fast) var(--ease),
    transform var(--dur-fast) var(--ease);
  user-select: none;
}

.btn--md {
  --h: var(--btn-h-md);
  padding: 0 24px;
}

.btn--block {
  width: 100%;
}

.btn:active:not(:disabled) {
  transform: scale(0.98);
}

.btn:disabled {
  opacity: 0.55;
  cursor: not-allowed;
}

.btn__icon {
  display: inline-flex;
  transition: transform var(--dur) var(--ease);
}

.btn:hover .btn__icon {
  transform: translate(2px, -2px);
}

.btn--primary {
  background: var(--ink);
  color: var(--bg);
}

.btn--primary:hover:not(:disabled) {
  background: var(--accent);
  color: var(--on-accent);
}

.btn--secondary {
  background: transparent;
  color: var(--ink);
  border-color: var(--line-strong);
}

.btn--secondary:hover:not(:disabled) {
  border-color: var(--ink);
  background: var(--surface);
}

.btn--ghost {
  background: var(--surface-2);
  color: var(--ink);
}

.btn--ghost:hover:not(:disabled) {
  background: var(--surface-3);
}

.btn--light {
  background: var(--on-dark);
  color: var(--dark);
}

.btn--light:hover:not(:disabled) {
  background: var(--accent);
  color: var(--on-accent);
}

.btn--outline-light {
  color: var(--on-dark);
  border-color: var(--dark-line);
}

.btn--outline-light:hover:not(:disabled) {
  border-color: var(--on-dark);
}

.is-loading .btn__label,
.is-loading .btn__icon {
  opacity: 0;
}

.btn__spinner {
  position: absolute;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 2px solid currentColor;
  border-right-color: transparent;
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
