<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ status: string }>()

const labels: Record<string, string> = {
  draft: 'Qoralama',
  published: 'E’lon qilingan',
  archived: 'Arxiv',
  new: 'Yangi',
  in_progress: 'Jarayonda',
  done: 'Yakunlangan',
  spam: 'Spam',
  on: 'Ko‘rinadi',
  off: 'Yashirin',
}

const label = computed(() => labels[props.status] ?? props.status)
</script>

<template>
  <span class="a-badge" :class="`a-badge--${status}`">{{ label }}</span>
</template>
