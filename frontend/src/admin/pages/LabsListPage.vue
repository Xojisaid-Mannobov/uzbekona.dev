<script setup lang="ts">
import SortableResourceList from '@/admin/components/SortableResourceList.vue'
import StatusBadge from '@/admin/components/StatusBadge.vue'
import { adminApi } from '@/admin/services/adminApi'
import { labStageLabels } from '@/content/site'
import type { Lab } from '@/types/api'
</script>

<template>
  <SortableResourceList
    title="Uzbekona Labs"
    subtitle="Ichki va open-source loyihalar: Telekit, Save Flow, AI tools …"
    base-path="/admin/labs"
    new-label="Yangi loyiha"
    empty-title="Labs bo‘limi bo‘sh"
    empty-text="Open-source yoki eksperimental loyihalarni qo‘shing."
    icon="flask"
    :item-name="(l: Lab) => l.title"
    :api="adminApi.labs"
  >
    <template #row="{ item }">
      <span>
        <span class="a-cell-title">{{ item.title }}</span>
        <span class="a-cell-sub" style="display: block">{{ labStageLabels[item.stage] }} · {{ item.stack.join(', ') }}</span>
      </span>
      <StatusBadge :status="item.is_published ? 'on' : 'off'" style="margin-left: auto" />
    </template>
  </SortableResourceList>
</template>
