<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import UiButton from '@/components/ui/UiButton.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import ProjectVisual from '@/components/media/ProjectVisual.vue'
import MediaImage from '@/components/media/MediaImage.vue'
import BlockRenderer from '@/components/content/BlockRenderer.vue'
import { useAsync } from '@/composables/useAsync'
import { useSeo } from '@/composables/useSeo'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'
import { paragraphs } from '@/utils/format'

const route = useRoute()
const slug = computed(() => String(route.params.slug))
const { data: project, loading, error, reload } = useAsync(() => publicApi.project(slug.value), { watch: [slug] })

// Client / Industry / Year / Services / Technology
const facts = computed(() => {
  const p = project.value
  if (!p) return []
  return [
    { label: 'Client', value: p.client },
    { label: 'Industry', value: p.industry },
    { label: 'Year', value: p.year ? String(p.year) : '' },
    { label: 'Services', value: p.services.join(', ') },
    { label: 'Technology', value: p.stack.join(', ') },
  ].filter((f) => f.value)
})

useSeo(() => ({
  title: project.value ? project.value.seo.title || project.value.title : error.value ? 'Loyiha topilmadi' : 'Loyiha',
  description: project.value?.seo.description || project.value?.tagline,
  image: project.value?.cover?.url,
  noindex: error.value?.status === 404,
}))
</script>

<template>
  <div>
    <!-- Yuklanmoqda -->
    <div v-if="loading && !project" class="container detail-skeleton" aria-busy="true">
      <div class="skeleton" style="height: 24px; width: 180px" />
      <div class="skeleton" style="height: 120px; width: 60%; margin-top: 32px" />
      <div class="skeleton" style="height: 700px; margin-top: 80px; border-radius: var(--r-lg)" />
    </div>

    <!-- Xato / topilmadi -->
    <div v-else-if="error" class="container section">
      <StateMessage
        tone="error"
        :title="error.status === 404 ? 'Loyiha topilmadi' : 'Loyihani yuklab bo‘lmadi'"
        :text="error.status === 404 ? 'Havola eskirgan yoki loyiha hali e’lon qilinmagan.' : error.message"
        :action-label="error.status === 404 ? undefined : 'Qayta urinish'"
        @action="reload"
      >
        <UiButton v-if="error.status === 404" to="/projects" variant="secondary" icon-left="arrow-left">Barcha loyihalar</UiButton>
      </StateMessage>
    </div>

    <article v-else-if="project" :key="project.id">
      <header class="container case-hero">
        <RouterLink to="/projects" class="case-hero__back"><AppIcon name="arrow-left" :size="18" /> Loyihalar</RouterLink>
        <h1 class="t-display case-hero__title" v-reveal>{{ project.title }}</h1>
        <p class="t-statement case-hero__tagline" v-reveal="80">{{ project.tagline }}</p>

        <dl class="facts" v-reveal="140">
          <div v-for="f in facts" :key="f.label" class="facts__item">
            <dt>{{ f.label }}</dt>
            <dd>{{ f.value }}</dd>
          </div>
        </dl>
      </header>

      <div class="container">
        <div class="case-cover" v-reveal="{ variant: 'scale' }">
          <ProjectVisual
            :title="project.title"
            :subtitle="project.tagline"
            :accent="project.accent"
            :cover="project.cover"
            :seed="project.slug"
            variant="dashboard"
            priority
          />
        </div>
      </div>

      <section v-if="project.full_description || project.live_url" class="container section intro">
        <h2 class="t-label">Loyiha haqida</h2>
        <div class="intro__body">
          <p v-for="(p, i) in paragraphs(project.full_description)" :key="i" class="intro__text" v-reveal>{{ p }}</p>
          <UiButton v-if="project.live_url" :href="project.live_url" variant="secondary" icon="arrow-up-right">Saytni ochish</UiButton>
        </div>
      </section>

      <section v-if="project.metrics.length" class="container section section--flush-top">
        <dl class="case-metrics">
          <div v-for="(m, i) in project.metrics" :key="i" class="case-metrics__item" v-reveal="i * 80">
            <dt>{{ m.label }}</dt>
            <dd>{{ m.value }}</dd>
          </div>
        </dl>
      </section>

      <section v-if="project.blocks?.length" class="container section section--flush-top">
        <h2 class="visually-hidden">Case study</h2>
        <BlockRenderer :blocks="project.blocks" mode="case" />
      </section>

      <section v-if="project.gallery?.length" class="container section section--flush-top">
        <div class="case-gallery">
          <figure
            v-for="(g, i) in project.gallery"
            :key="g.media_id"
            class="case-gallery__item"
            :class="{ 'is-wide': i % 3 === 0 }"
            v-reveal="(i % 3 === 2 ? 1 : 0) * 100"
          >
            <div class="case-gallery__frame">
              <MediaImage
                v-if="g.media"
                :media="g.media"
                :sizes="i % 3 === 0 ? '(min-width: 1440px) 1360px, 100vw' : '(min-width: 1024px) 680px, 100vw'"
              />
            </div>
            <figcaption v-if="g.caption">{{ g.caption }}</figcaption>
          </figure>
        </div>
      </section>

      <!-- Keyingi loyiha -->
      <RouterLink v-if="project.next" :to="`/projects/${project.next.slug}`" class="next">
        <div class="container next__inner">
          <p class="t-label">Keyingi loyiha</p>
          <div class="next__row">
            <span class="t-display next__title">{{ project.next.title }}</span>
            <span class="next__arrow"><AppIcon name="arrow-right" :size="32" /></span>
          </div>
          <p class="next__tagline">{{ project.next.tagline }}</p>
        </div>
      </RouterLink>
    </article>
  </div>
</template>

<style scoped>
.detail-skeleton {
  padding-block: 96px;
}

.case-hero {
  padding-top: clamp(48px, 6vw, 96px);
  padding-bottom: clamp(56px, 6vw, 96px);
}

.case-hero__back {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--ink-2);
  margin-bottom: 40px;
}

.case-hero__back:hover {
  color: var(--ink);
}

.case-hero__title {
  font-size: clamp(64px, 10vw, 160px);
}

.case-hero__tagline {
  margin-top: 28px;
  max-width: 18ch;
  color: var(--ink-2);
}

.facts {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 32px;
  margin-top: clamp(56px, 6vw, 96px);
}

.facts__item {
  padding-top: 20px;
  border-top: 1px solid var(--line-strong);
}

.facts dt {
  font-size: var(--fs-xs);
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--ink-3);
}

.facts dd {
  margin-top: 10px;
  font-size: 18px;
  font-weight: 600;
  line-height: 1.4;
}

.case-cover {
  position: relative;
  height: clamp(420px, 52vw, 800px);
  border-radius: var(--r-lg);
  overflow: hidden;
}

.intro {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 3fr);
  gap: 48px;
}

.intro__body {
  display: grid;
  gap: 28px;
  justify-items: start;
}

.intro__text {
  font-size: var(--fs-statement);
  line-height: 1.15;
  letter-spacing: -0.03em;
  font-weight: 500;
  max-width: 24ch;
}

.case-metrics {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 32px;
}

.case-metrics__item {
  display: flex;
  flex-direction: column-reverse;
  gap: 12px;
  padding-top: 28px;
  border-top: 1px solid var(--line-strong);
}

.case-metrics dd {
  font-size: var(--fs-metric);
  line-height: 0.9;
  font-weight: 600;
  letter-spacing: -0.055em;
}

.case-metrics dt {
  color: var(--ink-2);
}

.case-gallery {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: clamp(16px, 2vw, 32px);
}

.case-gallery__item.is-wide {
  grid-column: 1 / -1;
}

.case-gallery__frame {
  border-radius: var(--r-lg);
  overflow: hidden;
  background: var(--surface-2);
}

.case-gallery figcaption {
  margin-top: 14px;
  font-size: var(--fs-small);
  color: var(--ink-3);
}

.next {
  display: block;
  border-top: 1px solid var(--line);
  padding-block: clamp(80px, 9vw, 140px);
  transition: background-color var(--dur) var(--ease);
}

.next:hover {
  background: var(--surface);
}

.next__row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 32px;
  margin-top: 24px;
}

.next__arrow {
  display: grid;
  place-items: center;
  width: clamp(72px, 7vw, 112px);
  aspect-ratio: 1;
  border-radius: 50%;
  background: var(--ink);
  color: var(--bg);
  flex-shrink: 0;
  transition: transform var(--dur) var(--ease);
}

.next:hover .next__arrow {
  transform: translateX(8px);
  background: var(--accent);
}

.next__tagline {
  margin-top: 16px;
  font-size: var(--fs-body-lg);
  color: var(--ink-2);
}

@media (max-width: 1024px) {
  .facts {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .intro {
    grid-template-columns: 1fr;
    gap: 24px;
  }
}

@media (max-width: 768px) {
  .case-cover {
    height: auto;
    aspect-ratio: 4 / 3.6;
  }

  .case-gallery {
    grid-template-columns: 1fr;
  }
}
</style>
