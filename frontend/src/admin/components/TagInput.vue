<script setup lang="ts">
import { ref } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'

/** Massiv qiymatlar (stack, xizmatlar, platformalar) uchun chip input. Enter yoki vergul — qo'shish. */
const model = defineModel<string[]>({ required: true })
defineProps<{ placeholder?: string; id?: string; suggestions?: string[] }>()

const draft = ref('')

function add(value = draft.value) {
  const v = value.trim().replace(/,$/, '')
  if (v && !model.value.includes(v)) model.value = [...model.value, v]
  draft.value = ''
}

function removeAt(i: number) {
  model.value = model.value.filter((_, idx) => idx !== i)
}

function onKey(e: KeyboardEvent) {
  if (e.key === 'Enter' || e.key === ',') {
    e.preventDefault()
    add()
  } else if (e.key === 'Backspace' && !draft.value && model.value.length) {
    removeAt(model.value.length - 1)
  }
}
</script>

<template>
  <div class="tags">
    <div class="tags__box">
      <span v-for="(t, i) in model" :key="t" class="tags__chip">
        {{ t }}
        <button type="button" :aria-label="`${t} — olib tashlash`" @click="removeAt(i)"><AppIcon name="close" :size="14" /></button>
      </span>
      <input :id="id" v-model="draft" class="tags__input" :placeholder="model.length ? '' : placeholder" @keydown="onKey" @blur="add()" />
    </div>
    <div v-if="suggestions?.length" class="tags__suggest">
      <button v-for="s in suggestions.filter((x) => !model.includes(x))" :key="s" type="button" @click="add(s)">+ {{ s }}</button>
    </div>
  </div>
</template>

<style scoped>
.tags__box {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  min-height: var(--a-control-h);
  padding: 6px 8px;
  border-radius: var(--a-radius);
  border: 1px solid var(--line-strong);
  background: var(--surface);
  transition:
    border-color var(--dur-fast) var(--ease),
    box-shadow var(--dur-fast) var(--ease);
}

.tags__box:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 0 4px var(--accent-soft);
}

.tags__chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 32px;
  padding: 0 6px 0 12px;
  border-radius: 99px;
  background: var(--surface-2);
  font-size: 13px;
  font-weight: 600;
}

.tags__chip button {
  display: grid;
  place-items: center;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  color: var(--ink-3);
}

.tags__chip button:hover {
  background: var(--surface-3);
  color: var(--ink);
}

.tags__input {
  flex: 1;
  min-width: 120px;
  height: 32px;
  border: 0;
  outline: none;
  background: transparent;
  padding: 0 6px;
  font-size: 14px;
}

.tags__suggest {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
}

.tags__suggest button {
  font-size: 12px;
  font-weight: 600;
  color: var(--ink-3);
  padding: 4px 10px;
  border-radius: 99px;
  border: 1px dashed var(--line-strong);
}

.tags__suggest button:hover {
  color: var(--ink);
  border-color: var(--ink);
}
</style>
