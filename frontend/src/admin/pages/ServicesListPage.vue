<script setup lang="ts">
import SortableResourceList from '@/admin/components/SortableResourceList.vue'
import StatusBadge from '@/admin/components/StatusBadge.vue'
import { adminApi } from '@/admin/services/adminApi'
import { pad2 } from '@/utils/format'
import type { Service } from '@/types/api'
</script>

<template>
  <SortableResourceList
    title="Xizmatlar"
    subtitle="Saytdagi “Xizmatlar” bo‘limi va /services/:slug sahifalari."
    base-path="/admin/services"
    new-label="Yangi xizmat"
    empty-title="Hozircha xizmat yo‘q"
    empty-text="Studiya taklif qiladigan yo‘nalishlarni qo‘shing."
    icon="layers"
    :item-name="(s: Service) => s.title"
    :public-url="(s: Service) => (s.status === 'published' ? `/services/${s.slug}` : null)"
    :api="adminApi.services"
  >
    <template #row="{ item, index }">
      <span class="num">{{ pad2(index + 1) }}</span>
      <span style="min-width: 0">
        <span class="a-cell-title">{{ item.title }}</span>
        <span class="a-cell-sub" style="display: block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap">{{
          item.summary
        }}</span>
      </span>
      <StatusBadge :status="item.status" style="margin-left: auto" />
    </template>
  </SortableResourceList>
</template>

<style scoped>
.num {
  font-weight: 700;
  color: var(--ink-3);
  font-variant-numeric: tabular-nums;
  width: 28px;
}
</style>
