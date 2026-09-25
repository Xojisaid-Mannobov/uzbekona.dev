<script setup lang="ts">
import { nextTick, ref, watch } from 'vue'
import { useUiStore } from '@/admin/stores/ui'

const ui = useUiStore()
const confirmBtn = ref<HTMLButtonElement>()

// Oyna ochilganda fokus tasdiqlash tugmasiga o'tadi (klaviatura bilan ishlash uchun)
watch(
  () => ui.confirmState,
  async (s) => {
    if (s) {
      await nextTick()
      confirmBtn.value?.focus()
    }
  },
)
</script>

<template>
  <Transition name="modal">
    <div v-if="ui.confirmState" class="overlay" @click.self="ui.closeConfirm(false)" @keydown.esc="ui.closeConfirm(false)">
      <div class="dialog" role="alertdialog" aria-modal="true" aria-labelledby="confirm-title">
        <h2 id="confirm-title">{{ ui.confirmState.title }}</h2>
        <p v-if="ui.confirmState.text">{{ ui.confirmState.text }}</p>
        <div class="dialog__actions">
          <button type="button" class="a-btn a-btn--ghost" @click="ui.closeConfirm(false)">Bekor qilish</button>
          <button
            ref="confirmBtn"
            type="button"
            class="a-btn"
            :class="ui.confirmState.danger ? 'a-btn--danger' : 'a-btn--primary'"
            @click="ui.closeConfirm(true)"
          >
            {{ ui.confirmState.confirmLabel }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  z-index: 90;
  display: grid;
  place-items: center;
  padding: 16px;
  background: rgb(17 19 18 / 0.45);
  backdrop-filter: blur(4px);
}

.dialog {
  width: min(440px, 100%);
  padding: 28px;
  border-radius: 24px;
  background: var(--surface);
  box-shadow: var(--shadow-lg);
}

.dialog h2 {
  font-size: 20px;
  letter-spacing: -0.02em;
}

.dialog p {
  margin-top: 10px;
  color: var(--ink-2);
}

.dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 28px;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 200ms var(--ease);
}

.modal-enter-active .dialog,
.modal-leave-active .dialog {
  transition: transform 260ms var(--ease);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .dialog {
  transform: scale(0.96) translateY(8px);
}
</style>
