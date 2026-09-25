<script setup lang="ts">
import AppIcon from '@/components/ui/AppIcon.vue'

// Tahrirlash sahifasining "Saqlash" kartasi: holat, saqlanmagan o'zgarishlar, havolalar
defineProps<{ saving: boolean; dirty: boolean; isNew: boolean; publicUrl?: string | null }>()
defineEmits<{ save: [] }>()
</script>

<template>
  <section class="a-card publish">
    <slot />
    <div class="publish__foot">
      <span class="publish__state" :class="{ 'is-dirty': dirty }">
        <span class="publish__dot" />
        {{ dirty ? 'Saqlanmagan o‘zgarishlar' : isNew ? 'Yangi yozuv' : 'Barcha o‘zgarishlar saqlangan' }}
      </span>
      <button type="button" class="a-btn a-btn--primary publish__save" :disabled="saving || (!dirty && !isNew)" @click="$emit('save')">
        <AppIcon name="check" :size="18" /> {{ saving ? 'Saqlanmoqda…' : isNew ? 'Yaratish' : 'Saqlash' }}
      </button>
      <a v-if="publicUrl" :href="publicUrl" target="_blank" class="a-btn a-btn--ghost publish__view"
        ><AppIcon name="external" :size="16" /> Saytda ko‘rish</a
      >
      <p class="a-hint publish__kbd">Tezkor saqlash: <kbd>Ctrl</kbd> + <kbd>S</kbd></p>
    </div>
  </section>
</template>

<style scoped>
.publish {
  display: grid;
  gap: 18px;
}

.publish__foot {
  display: grid;
  gap: 10px;
  padding-top: 18px;
  border-top: 1px solid var(--line);
}

.publish__state {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink-3);
}

.publish__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--success);
}

.is-dirty {
  color: var(--warning);
}

.is-dirty .publish__dot {
  background: var(--warning);
}

.publish__save,
.publish__view {
  width: 100%;
}

.publish__kbd {
  text-align: center;
}

kbd {
  font-family: var(--font-mono);
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 6px;
  border: 1px solid var(--line-strong);
}
</style>
