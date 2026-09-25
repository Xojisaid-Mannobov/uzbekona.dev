<script setup lang="ts">
import { computed } from 'vue'
import MediaImage from '@/components/media/MediaImage.vue'
import GirihPattern from '@/components/ornament/GirihPattern.vue'
import SuzaniRosette from '@/components/ornament/SuzaniRosette.vue'
import type { News } from '@/types/api'
import { toneStyle } from './tone'

// Yangilik muqovasi: rasm bo'lsa — rasm, bo'lmasa teg rangidagi naqshli "plakat"
const props = withDefaults(defineProps<{ news: News; sizes?: string; priority?: boolean }>(), {
  sizes: '(min-width: 1024px) 440px, 100vw',
  priority: false,
})
const style = computed(() => toneStyle(props.news.tag))
</script>

<template>
  <div class="cover" :style="style">
    <MediaImage v-if="news.cover" :media="news.cover" :sizes="sizes" :priority="priority" fill />
    <div v-else class="cover__poster" aria-hidden="true">
      <GirihPattern tone="light" :size="64" :opacity="0.16" fade="right" />
      <SuzaniRosette class="cover__rosette" tone="light" :opacity="0.42" />
      <span class="cover__word">{{ news.tag || 'Yangilik' }}</span>
    </div>
  </div>
</template>

<style scoped>
.cover {
  position: relative;
  overflow: hidden;
  background: var(--surface-2);
}

.cover__poster {
  position: absolute;
  inset: 0;
  isolation: isolate;
  background:
    radial-gradient(120% 90% at 100% 0%, rgb(255 255 255 / 0.18), transparent 55%),
    linear-gradient(150deg, color-mix(in srgb, var(--tone) 92%, #fff) 0%, color-mix(in srgb, var(--tone) 70%, #0b1330) 100%);
}

.cover__rosette {
  position: absolute;
  z-index: -1;
  right: -18%;
  top: -22%;
  width: 72%;
  aspect-ratio: 1;
}

.cover__word {
  position: absolute;
  left: 7%;
  bottom: 6%;
  font-size: clamp(34px, 4.2vw, 64px);
  font-weight: 700;
  letter-spacing: -0.05em;
  line-height: 1;
  color: #fff;
}
</style>
