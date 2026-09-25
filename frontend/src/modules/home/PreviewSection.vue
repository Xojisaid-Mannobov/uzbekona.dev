<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import ProjectVisual from '@/components/media/ProjectVisual.vue'
import type { Project } from '@/types/api'
import { loadMotion, reducedMotion, whenIdle } from '@/composables/useSmoothScroll'

// Hero ostidagi katta mahsulot ko'rinishi (650–800px) — scroll bilan kattalashadi
defineProps<{ project: Project | null; loading: boolean }>()

const frame = ref<HTMLElement>()
let ctx: { revert(): void } | undefined
let alive = true

onMounted(() => {
  if (reducedMotion() || !frame.value) return
  whenIdle(async () => {
    const { gsap } = await loadMotion()
    if (!alive || !frame.value) return
    ctx = gsap.context(() => {
      gsap.fromTo(
        frame.value!,
        { scale: 0.92, y: 40 },
        {
          scale: 1,
          y: 0,
          ease: 'none',
          scrollTrigger: { trigger: frame.value!, start: 'top bottom', end: 'top 25%', scrub: 0.6 },
        },
      )
    })
  })
})

onBeforeUnmount(() => {
  alive = false
  ctx?.revert()
})
</script>

<template>
  <section class="preview" aria-label="Loyiha ko‘rinishi">
    <div class="container">
      <div ref="frame" class="preview__frame">
        <div v-if="loading && !project" class="preview__skeleton skeleton" />
        <ProjectVisual
          v-else
          :title="project?.title ?? 'Uzbekona.dev'"
          :subtitle="project?.tagline"
          :accent="project?.accent"
          :cover="project?.cover"
          variant="dashboard"
          priority
        />
      </div>

      <RouterLink v-if="project" :to="`/projects/${project.slug}`" class="preview__caption">
        <span class="preview__name">{{ project.title }}</span>
        <span class="preview__tagline">{{ project.tagline }}</span>
        <span class="preview__more">Case study <AppIcon name="arrow-up-right" :size="18" /></span>
      </RouterLink>
    </div>
  </section>
</template>

<style scoped>
.preview__frame {
  position: relative;
  height: clamp(420px, 52vw, 800px);
  min-height: 650px;
  border-radius: var(--r-lg);
  overflow: hidden;
  transform-origin: 50% 0;
  box-shadow: var(--shadow-lg);
  will-change: transform;
}

.preview__skeleton {
  position: absolute;
  inset: 0;
  border-radius: 0;
}

.preview__caption {
  display: flex;
  align-items: baseline;
  gap: 20px;
  padding-top: 24px;
  font-size: var(--fs-small);
  color: var(--ink-2);
}

.preview__name {
  font-weight: 700;
  color: var(--ink);
}

.preview__more {
  margin-left: auto;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  color: var(--ink);
}

.preview__caption:hover .preview__more {
  color: var(--accent);
}

@media (max-width: 768px) {
  .preview__frame {
    min-height: 0;
    height: auto;
    aspect-ratio: 4 / 3.6;
  }

  .preview__tagline {
    display: none;
  }
}
</style>
