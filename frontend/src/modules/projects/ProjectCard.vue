<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import ProjectVisual from '@/components/media/ProjectVisual.vue'
import type { Project } from '@/types/api'

/**
 * Loyiha kartasi: vizual + nom, tavsif va meta.
 * `wide` — to'liq kenglik, `half` — 2 ustun (bosh sahifa), `compact` — 3 ustunli katalog grid'i.
 */
const props = withDefaults(defineProps<{ project: Project; size?: 'wide' | 'half' | 'compact'; priority?: boolean }>(), {
  size: 'wide',
})

const meta = computed(() => [props.project.platforms[0], props.project.industry, props.project.year].filter(Boolean).join(' · '))
const sizes = computed(() => {
  if (props.size === 'wide') return '(min-width: 1440px) 1360px, 100vw'
  if (props.size === 'compact') return '(min-width: 1024px) 440px, (min-width: 640px) 50vw, 100vw'
  return '(min-width: 1024px) 680px, 100vw'
})
</script>

<template>
  <article class="project" :class="`project--${size}`">
    <RouterLink :to="`/projects/${project.slug}`" class="project__link" :aria-label="`${project.title} — case study`">
      <div class="project__visual">
        <ProjectVisual
          :title="project.title"
          :subtitle="project.tagline"
          :accent="project.accent"
          :cover="project.cover"
          :seed="project.slug"
          :sizes="sizes"
          :priority="priority"
        />
      </div>

      <div v-if="size === 'compact'" class="project__info project__info--compact">
        <div class="project__head">
          <h3 class="project__title">{{ project.title }}</h3>
          <span class="project__arrow"><AppIcon name="arrow-up-right" :size="18" /></span>
        </div>
        <p class="project__desc">{{ project.short_description || project.tagline }}</p>
        <ul class="project__tags" role="list">
          <li v-for="p in project.platforms" :key="p" class="project__tag">{{ p }}</li>
          <li v-if="project.industry" class="project__tag project__tag--muted">{{ project.industry }}</li>
          <li v-if="project.year" class="project__year">{{ project.year }}</li>
        </ul>
      </div>

      <div v-else class="project__info">
        <div class="project__text">
          <h3 class="project__title">{{ project.title }}</h3>
          <p class="project__desc">{{ project.short_description || project.tagline }}</p>
        </div>
        <div class="project__meta">
          <span>{{ meta }}</span>
          <span class="project__arrow"><AppIcon name="arrow-up-right" :size="22" /></span>
        </div>
      </div>
    </RouterLink>
  </article>
</template>

<style scoped>
.project__link {
  display: block;
}

.project__visual {
  position: relative;
  border-radius: var(--r-lg);
  overflow: hidden;
  min-height: 520px;
  height: clamp(380px, 44vw, 640px);
  isolation: isolate;
}

.project--half .project__visual {
  min-height: 0;
  height: clamp(320px, 36vw, 520px);
}

.project__visual > :deep(.visual) {
  transition: transform 1100ms var(--ease);
}

.project__link:hover .project__visual > :deep(.visual) {
  transform: scale(1.025);
}

.project__info {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 32px;
  padding-top: 28px;
}

.project__title {
  font-size: var(--fs-card);
  letter-spacing: -0.035em;
  line-height: 1.1;
}

.project__desc {
  margin-top: 6px;
  font-size: var(--fs-body-lg);
  color: var(--ink-2);
}

.project__meta {
  display: flex;
  align-items: center;
  gap: 20px;
  flex-shrink: 0;
  font-size: var(--fs-small);
  font-weight: 600;
  color: var(--ink-2);
}

.project__arrow {
  display: grid;
  place-items: center;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  border: 1px solid var(--line-strong);
  color: var(--ink);
  transition:
    background-color var(--dur) var(--ease),
    color var(--dur) var(--ease),
    border-color var(--dur) var(--ease),
    transform var(--dur) var(--ease);
}

.project__link:hover .project__arrow {
  background: var(--ink);
  color: var(--bg);
  border-color: var(--ink);
  transform: rotate(45deg);
}

.project--half .project__info {
  flex-direction: column;
  gap: 20px;
}

/* ─── Compact: katalog grid'i ─────────────────────── */
.project--compact .project__visual {
  min-height: 0;
  height: auto;
  aspect-ratio: 4 / 3;
  border: 1px solid var(--line);
}

.project__info--compact {
  display: grid;
  justify-content: stretch;
  gap: 8px;
  padding-top: 18px;
}

.project__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.project--compact .project__title {
  font-size: clamp(19px, 1.4vw, 22px);
}

.project--compact .project__desc {
  margin-top: 0;
  font-size: 15px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.project--compact .project__arrow {
  width: 38px;
  height: 38px;
  flex-shrink: 0;
}

.project__tags {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  margin-top: 6px;
}

.project__tag {
  padding: 4px 10px;
  border-radius: var(--r-pill);
  background: var(--surface-2);
  font-size: 12px;
  font-weight: 600;
  color: var(--ink);
}

.project__tag--muted {
  background: transparent;
  border: 1px solid var(--line);
  color: var(--ink-2);
}

.project__year {
  margin-left: auto;
  font-size: 13px;
  font-weight: 600;
  color: var(--ink-3);
}

@media (max-width: 768px) {
  .project__visual,
  .project--half .project__visual {
    min-height: 0;
    height: auto;
    aspect-ratio: 4 / 4.4;
  }

  .project--compact .project__visual {
    aspect-ratio: 16 / 10;
  }

  .project__info:not(.project__info--compact) {
    flex-direction: column;
    gap: 16px;
    padding-top: 20px;
  }

  .project__arrow {
    width: 44px;
    height: 44px;
  }
}
</style>
