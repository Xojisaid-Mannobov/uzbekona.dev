<script setup lang="ts">
import SortableResourceList from '@/admin/components/SortableResourceList.vue'
import StatusBadge from '@/admin/components/StatusBadge.vue'
import { adminApi } from '@/admin/services/adminApi'
import type { TeamMember } from '@/types/api'
</script>

<template>
  <SortableResourceList
    title="Jamoa"
    subtitle="Saytdagi jamoa a’zolari (3 ustunli grid, rasm 4:5)."
    base-path="/admin/team"
    new-label="Yangi a’zo"
    empty-title="Jamoa ro‘yxati bo‘sh"
    empty-text="Mahsulotni quradigan odamlarni qo‘shing."
    icon="users"
    :item-name="(m: TeamMember) => m.name"
    :api="adminApi.team"
  >
    <template #row="{ item }">
      <span class="avatar">
        <img v-if="item.photo" :src="item.photo.variants[0]?.url ?? item.photo.url" alt="" />
        <span v-else>{{ item.name.charAt(0) }}</span>
      </span>
      <span>
        <span class="a-cell-title">{{ item.name }}</span>
        <span class="a-cell-sub" style="display: block">{{ item.role }}</span>
      </span>
      <StatusBadge :status="item.is_published ? 'on' : 'off'" style="margin-left: auto" />
    </template>
  </SortableResourceList>
</template>

<style scoped>
.avatar {
  display: grid;
  place-items: center;
  width: 48px;
  height: 60px;
  border-radius: 10px;
  overflow: hidden;
  background: var(--surface-2);
  font-weight: 700;
  color: var(--ink-3);
  flex-shrink: 0;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
