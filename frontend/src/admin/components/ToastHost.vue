<script setup lang="ts">
import AppIcon from '@/components/ui/AppIcon.vue'
import { useUiStore } from '@/admin/stores/ui'

const ui = useUiStore()
</script>

<template>
  <div class="toasts" role="status" aria-live="polite">
    <TransitionGroup name="toast">
      <div v-for="t in ui.toasts" :key="t.id" class="toast" :class="`toast--${t.tone}`">
        <AppIcon :name="t.tone === 'error' ? 'alert' : 'check'" :size="18" />
        <span>{{ t.message }}</span>
        <button type="button" class="toast__close" aria-label="Yopish" @click="ui.dismiss(t.id)">
          <AppIcon name="close" :size="16" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toasts {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 100;
  display: grid;
  gap: 10px;
  width: min(380px, calc(100vw - 32px));
}

.toast {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 14px 14px 18px;
  border-radius: 16px;
  background: var(--ink);
  color: var(--bg);
  font-weight: 600;
  font-size: 14px;
  box-shadow: var(--shadow-lg);
}

.toast--success .icon {
  color: #4ade80;
}

.toast--error {
  background: var(--danger);
  color: #fff;
}

.toast__close {
  margin-left: auto;
  opacity: 0.6;
}

.toast__close:hover {
  opacity: 1;
}

.toast-enter-active,
.toast-leave-active {
  transition:
    opacity 280ms var(--ease),
    transform 280ms var(--ease);
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(12px);
}
</style>
