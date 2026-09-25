<script setup lang="ts">
import { ref, toRef } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'
import BlockEditor from './BlockEditor.vue'
import { blockDefs, blockSummary, blockTypeList, createBlock } from './blockTypes'
import { useDragSort } from '@/admin/composables/useDragSort'
import { useUiStore } from '@/admin/stores/ui'
import type { Block, BlockType } from '@/types/api'

/**
 * Modular content builder: bloklarni qo'shish, tahrirlash, o'chirish,
 * nusxalash va tartiblash (drag & drop yoki ↑/↓ tugmalari).
 */
const model = defineModel<Block[]>({ required: true })
const ui = useUiStore()

// Har bir blok uchun barqaror kalit (v-for) — ma'lumotga qo'shimcha maydon qo'shmasdan
const keys = new WeakMap<Block, number>()
let keySeq = 0
function keyOf(b: Block) {
  if (!keys.has(b)) keys.set(b, ++keySeq)
  return keys.get(b)!
}

const collapsed = ref(new Set<number>())
const paletteAt = ref<number | null>(null)

const { handlers, dragIndex, overIndex, move } = useDragSort(toRef(model))

function toggle(b: Block) {
  const k = keyOf(b)
  const s = new Set(collapsed.value)
  if (s.has(k)) s.delete(k)
  else s.add(k)
  collapsed.value = s
}

function insert(type: BlockType) {
  const at = paletteAt.value ?? model.value.length
  const list = [...model.value]
  list.splice(at, 0, createBlock(type))
  model.value = list
  paletteAt.value = null
}

function duplicate(i: number) {
  const copy = JSON.parse(JSON.stringify(model.value[i])) as Block
  delete copy.id
  const list = [...model.value]
  list.splice(i + 1, 0, copy)
  model.value = list
}

async function remove(i: number) {
  const ok = await ui.confirm({
    title: 'Blok o‘chirilsinmi?',
    text: `“${blockDefs[model.value[i].type].label}” bloki ro‘yxatdan olib tashlanadi.`,
  })
  if (ok) model.value = model.value.filter((_, idx) => idx !== i)
}

function collapseAll(value: boolean) {
  collapsed.value = value ? new Set(model.value.map(keyOf)) : new Set()
}
</script>

<template>
  <div class="builder">
    <div v-if="model.length" class="builder__bar">
      <span class="a-hint">{{ model.length }} ta blok · sudrab tartiblang</span>
      <div>
        <button type="button" class="a-btn a-btn--ghost a-btn--sm" @click="collapseAll(true)">Hammasini yig‘ish</button>
        <button type="button" class="a-btn a-btn--ghost a-btn--sm" @click="collapseAll(false)">Ochish</button>
      </div>
    </div>

    <TransitionGroup name="blk" tag="div" class="builder__list">
      <div v-for="(block, i) in model" :key="keyOf(block)" class="builder__slot">
        <article
          class="blk"
          :class="{ 'is-dragging': dragIndex === i, 'is-over': overIndex === i && dragIndex !== i }"
          @dragover="handlers(i).onDragover"
          @drop="handlers(i).onDrop"
        >
          <!-- Faqat sarlavha qatori sudraladi — input'lardagi matn tanlashga xalaqit bermaydi -->
          <header class="blk__head" draggable="true" @dragstart="handlers(i).onDragstart" @dragend="handlers(i).onDragend">
            <span class="a-handle" aria-hidden="true"><AppIcon name="grip" :size="18" /></span>
            <span class="blk__icon"><AppIcon :name="blockDefs[block.type].icon" :size="18" /></span>
            <button type="button" class="blk__title" :aria-expanded="!collapsed.has(keyOf(block))" @click="toggle(block)">
              <strong>{{ blockDefs[block.type].label }}</strong>
              <span>{{ blockSummary(block) }}</span>
            </button>
            <div class="blk__actions">
              <button type="button" class="a-icon-btn" aria-label="Yuqoriga" :disabled="i === 0" @click="move(i, i - 1)">
                <AppIcon name="arrow-up" :size="17" />
              </button>
              <button type="button" class="a-icon-btn" aria-label="Pastga" :disabled="i === model.length - 1" @click="move(i, i + 1)">
                <AppIcon name="arrow-down" :size="17" />
              </button>
              <button type="button" class="a-icon-btn" aria-label="Nusxalash" @click="duplicate(i)">
                <AppIcon name="copy" :size="17" />
              </button>
              <button type="button" class="a-icon-btn a-icon-btn--danger" aria-label="O‘chirish" @click="remove(i)">
                <AppIcon name="trash" :size="17" />
              </button>
              <button
                type="button"
                class="a-icon-btn"
                :aria-label="collapsed.has(keyOf(block)) ? 'Ochish' : 'Yig‘ish'"
                @click="toggle(block)"
              >
                <AppIcon name="chevron-down" :size="18" :class="{ 'is-rot': collapsed.has(keyOf(block)) }" />
              </button>
            </div>
          </header>
          <div v-show="!collapsed.has(keyOf(block))" class="blk__body">
            <BlockEditor :block="block" />
          </div>
        </article>
        <button type="button" class="builder__insert" :aria-label="`${i + 1}-blokdan keyin qo‘shish`" @click="paletteAt = i + 1">
          <AppIcon name="plus" :size="16" />
        </button>
      </div>
    </TransitionGroup>

    <div v-if="!model.length" class="a-empty">
      <span class="a-empty__icon"><AppIcon name="layers" :size="24" /></span>
      <strong>Case-study hali bo‘sh</strong>
      <p>Sarlavha, matn, katta rasmlar, raqamlar va boshqa bloklardan sahifani yig‘ing.</p>
    </div>

    <button type="button" class="builder__add" @click="paletteAt = model.length"><AppIcon name="plus" :size="20" /> Blok qo‘shish</button>

    <!-- Blok turlari palitrasi -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="paletteAt !== null" class="palette-overlay admin" @click.self="paletteAt = null" @keydown.esc="paletteAt = null">
          <div class="palette" role="dialog" aria-modal="true" aria-label="Blok turini tanlang">
            <header>
              <h2>Blok qo‘shish</h2>
              <button type="button" class="a-icon-btn" aria-label="Yopish" @click="paletteAt = null"><AppIcon name="close" /></button>
            </header>
            <div class="palette__grid">
              <button v-for="def in blockTypeList" :key="def.type" type="button" class="palette__item" @click="insert(def.type)">
                <span class="palette__icon"><AppIcon :name="def.icon" :size="22" /></span>
                <strong>{{ def.label }}</strong>
                <span>{{ def.description }}</span>
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.builder {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 12px;
}

.builder__bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.builder__list {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
}

.builder__slot {
  position: relative;
  min-width: 0;
  padding-bottom: 14px;
}

.blk {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: 18px;
  transition:
    box-shadow var(--dur-fast) var(--ease),
    border-color var(--dur-fast) var(--ease),
    opacity var(--dur-fast) var(--ease);
}

.blk:hover {
  border-color: var(--line-strong);
}

.blk.is-dragging {
  opacity: 0.4;
}

.blk.is-over {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-soft);
}

.blk__head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 10px 10px 14px;
}

.blk__icon {
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  border-radius: 10px;
  background: var(--accent-soft);
  color: var(--accent);
  flex-shrink: 0;
}

.blk__title {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: baseline;
  gap: 10px;
  text-align: left;
}

.blk__title strong {
  font-size: 14px;
  white-space: nowrap;
}

.blk__title span {
  font-size: 13px;
  color: var(--ink-3);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.blk__actions {
  display: flex;
  gap: 0;
}

.blk__actions .icon.is-rot {
  transform: rotate(-90deg);
}

.blk__body {
  padding: 4px 18px 18px;
  border-top: 1px solid var(--line);
  padding-top: 18px;
}

.builder__insert {
  position: absolute;
  left: 50%;
  bottom: -5px;
  z-index: 2;
  translate: -50% 0;
  display: grid;
  place-items: center;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  background: var(--surface);
  border: 1px solid var(--line-strong);
  color: var(--ink-2);
  opacity: 0;
  transition:
    opacity var(--dur-fast) var(--ease),
    background-color var(--dur-fast) var(--ease);
}

.builder__slot:hover .builder__insert,
.builder__insert:focus-visible {
  opacity: 1;
}

.builder__insert:hover {
  background: var(--accent);
  border-color: var(--accent);
  color: #fff;
}

.builder__add {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 64px;
  border-radius: 18px;
  border: 1.5px dashed var(--line-strong);
  font-weight: 700;
  color: var(--ink-2);
  transition:
    border-color var(--dur-fast) var(--ease),
    color var(--dur-fast) var(--ease),
    background-color var(--dur-fast) var(--ease);
}

.builder__add:hover {
  border-color: var(--accent);
  color: var(--accent);
  background: var(--accent-soft);
}

.palette-overlay {
  position: fixed;
  inset: 0;
  z-index: 85;
  display: grid;
  place-items: center;
  padding: 16px;
  background: rgb(17 19 18 / 0.45);
  backdrop-filter: blur(4px);
}

.palette {
  width: min(760px, 100%);
  max-height: 90vh;
  overflow-y: auto;
  padding: 24px;
  border-radius: 24px;
  background: var(--surface);
  box-shadow: var(--shadow-lg);
}

.palette header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
}

.palette h2 {
  font-size: 20px;
  letter-spacing: -0.02em;
}

.palette__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 10px;
}

.palette__item {
  display: grid;
  justify-items: start;
  gap: 4px;
  padding: 16px;
  border-radius: 16px;
  border: 1px solid var(--line);
  text-align: left;
  transition:
    border-color var(--dur-fast) var(--ease),
    background-color var(--dur-fast) var(--ease);
}

.palette__item:hover,
.palette__item:focus-visible {
  border-color: var(--accent);
  background: var(--accent-soft);
}

.palette__icon {
  display: grid;
  place-items: center;
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: var(--surface-2);
  margin-bottom: 8px;
}

.palette__item strong {
  font-size: 15px;
}

.palette__item span:last-child {
  font-size: 13px;
  color: var(--ink-3);
}

.blk-move,
.blk-enter-active,
.blk-leave-active {
  transition: all 300ms var(--ease);
}

.blk-enter-from,
.blk-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

.blk-leave-active {
  position: absolute;
  width: 100%;
}

.modal-enter-active,
.modal-leave-active {
  transition: opacity 200ms var(--ease);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
