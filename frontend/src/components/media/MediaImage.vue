<script setup lang="ts">
import { computed, ref } from 'vue'
import type { MediaRef } from '@/types/api'
import { srcset } from '@/utils/media'

/**
 * Responsive rasm: backend variantlaridan srcset, lazy loading, yuklanganda yumshoq paydo bo'lish.
 * `fill` — ota elementni to'liq qoplaydi (object-fit: cover).
 */
const props = withDefaults(
  defineProps<{
    media: MediaRef
    sizes?: string
    alt?: string
    priority?: boolean
    fill?: boolean
    fit?: 'cover' | 'contain'
  }>(),
  { sizes: '100vw', fit: 'cover' },
)

const loaded = ref(false)
const set = computed(() => srcset(props.media))
</script>

<template>
  <video
    v-if="media.kind === 'video'"
    class="media"
    :class="{ 'media--fill': fill }"
    :src="media.url"
    autoplay
    muted
    loop
    playsinline
    preload="metadata"
    :aria-label="alt ?? media.alt"
  />
  <img
    v-else
    class="media"
    :class="{ 'media--fill': fill, 'is-loaded': loaded || priority }"
    :style="{ objectFit: fit }"
    :src="media.url"
    :srcset="set"
    :sizes="set ? sizes : undefined"
    :width="media.width ?? undefined"
    :height="media.height ?? undefined"
    :alt="alt ?? media.alt"
    :loading="priority ? 'eager' : 'lazy'"
    :fetchpriority="priority ? 'high' : undefined"
    decoding="async"
    @load="loaded = true"
  />
</template>

<style scoped>
.media {
  width: 100%;
  height: auto;
  opacity: 0;
  transition: opacity var(--dur-slow) var(--ease);
}

video.media,
.media.is-loaded {
  opacity: 1;
}

.media--fill {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}
</style>
