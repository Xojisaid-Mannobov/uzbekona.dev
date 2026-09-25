<script setup lang="ts">
import { computed, defineAsyncComponent } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import PublicLayout from '@/layouts/PublicLayout.vue'

// Admin layout alohida chunk — public sayt foydalanuvchilari uni yuklamaydi
const AdminLayout = defineAsyncComponent(() => import('@/admin/layouts/AdminLayout.vue'))

const route = useRoute()
const layout = computed(() => {
  switch (route.meta.layout) {
    case 'admin':
      return AdminLayout
    case 'blank':
      return null
    default:
      return PublicLayout
  }
})
</script>

<template>
  <component :is="layout" v-if="layout" />
  <RouterView v-else />
</template>
