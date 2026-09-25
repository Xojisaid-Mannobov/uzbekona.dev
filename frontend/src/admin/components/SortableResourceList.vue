<script setup lang="ts" generic="T extends { id: number }">
import { ref, watch, type Ref } from 'vue'
import { RouterLink } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import type { IconName } from '@/components/ui/icons'
import PageHeader from './PageHeader.vue'
import AdminEmpty from './AdminEmpty.vue'
import { useAsync } from '@/composables/useAsync'
import { useDragSort } from '@/admin/composables/useDragSort'
import { useUiStore } from '@/admin/stores/ui'
import { toApiError } from '@/services/http'

/**
 * Xizmatlar, jamoa va Labs uchun umumiy ro'yxat: drag & drop tartiblash,
 * tahrirlash va o'chirish. Qator mazmuni slot orqali beriladi.
 */
const props = defineProps<{
  title: string
  subtitle: string
  basePath: string
  newLabel: string
  emptyTitle: string
  emptyText: string
  icon: IconName
  itemName: (item: T) => string
  publicUrl?: (item: T) => string | null
  api: {
    list: () => Promise<T[]>
    remove: (id: number) => Promise<void>
    reorder: (ids: number[]) => Promise<unknown>
  }
}>()

defineSlots<{ row(props: { item: T; index: number }): unknown }>()

const ui = useUiStore()
const { data, loading, error, reload } = useAsync(props.api.list)
const items = ref([]) as Ref<T[]>
watch(data, (d) => (items.value = d ? [...d] : []))

const { handlers, dragIndex, overIndex, move } = useDragSort(items, saveOrder)

async function saveOrder() {
  try {
    await props.api.reorder(items.value.map((i) => i.id))
    ui.success('Tartib saqlandi')
  } catch (e) {
    ui.error(toApiError(e).message)
    reload()
  }
}

async function remove(item: T) {
  const ok = await ui.confirm({ title: `“${props.itemName(item)}” o‘chirilsinmi?`, text: 'Bu amalni qaytarib bo‘lmaydi.' })
  if (!ok) return
  try {
    await props.api.remove(item.id)
    ui.success('O‘chirildi')
    reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}
</script>

<template>
  <div>
    <PageHeader :title="title" :subtitle="subtitle">
      <RouterLink :to="`${basePath}/new`" class="a-btn a-btn--primary"><AppIcon name="plus" :size="18" /> {{ newLabel }}</RouterLink>
    </PageHeader>

    <AdminEmpty v-if="error" error title="Yuklab bo‘lmadi" :text="error.message" @retry="reload" />

    <div v-else-if="loading && !data" class="a-stack">
      <div v-for="n in 4" :key="n" class="skeleton" style="height: 72px; border-radius: 14px" />
    </div>

    <AdminEmpty v-else-if="!items.length" :icon="icon" :title="emptyTitle" :text="emptyText">
      <RouterLink :to="`${basePath}/new`" class="a-btn a-btn--primary a-btn--sm">Birinchisini qo‘shish</RouterLink>
    </AdminEmpty>

    <template v-else>
      <ul class="rl" role="list">
        <li
          v-for="(item, i) in items"
          :key="item.id"
          class="rl__item"
          :class="{ 'is-dragging': dragIndex === i, 'is-over': overIndex === i && dragIndex !== i }"
          v-bind="handlers(i)"
        >
          <span class="a-handle" title="Sudrab tartiblang"><AppIcon name="grip" :size="18" /></span>
          <RouterLink :to="`${basePath}/${item.id}`" class="rl__main">
            <slot name="row" :item="item" :index="i" />
          </RouterLink>
          <div class="a-actions">
            <button type="button" class="a-icon-btn" aria-label="Yuqoriga" :disabled="i === 0" @click="move(i, i - 1)">
              <AppIcon name="arrow-up" :size="17" />
            </button>
            <button type="button" class="a-icon-btn" aria-label="Pastga" :disabled="i === items.length - 1" @click="move(i, i + 1)">
              <AppIcon name="arrow-down" :size="17" />
            </button>
            <a v-if="publicUrl?.(item)" :href="publicUrl(item)!" target="_blank" class="a-icon-btn" title="Saytda ko‘rish"
              ><AppIcon name="external" :size="18"
            /></a>
            <RouterLink :to="`${basePath}/${item.id}`" class="a-icon-btn" title="Tahrirlash"><AppIcon name="edit" :size="18" /></RouterLink>
            <button type="button" class="a-icon-btn a-icon-btn--danger" title="O‘chirish" @click="remove(item)">
              <AppIcon name="trash" :size="18" />
            </button>
          </div>
        </li>
      </ul>
      <p class="a-hint" style="margin-top: 12px">
        Tartib saytda xuddi shunday ko‘rinadi. Sudrab yoki ↑/↓ bilan o‘zgartiring — avtomatik saqlanadi.
      </p>
    </template>
  </div>
</template>

<style scoped>
.rl {
  display: grid;
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: var(--a-radius-lg);
  overflow: hidden;
}

.rl__item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 12px 12px 12px 16px;
  border-bottom: 1px solid var(--line);
  transition: background-color var(--dur-fast) var(--ease);
}

.rl__item:last-child {
  border-bottom: 0;
}

.rl__item:hover {
  background: color-mix(in srgb, var(--surface-2) 40%, var(--surface));
}

.rl__item.is-dragging {
  opacity: 0.4;
}

.rl__item.is-over {
  box-shadow: inset 0 2px 0 var(--accent);
}

.rl__main {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 14px;
}
</style>
