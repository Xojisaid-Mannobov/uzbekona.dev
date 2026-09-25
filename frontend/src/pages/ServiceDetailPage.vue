<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import UiButton from '@/components/ui/UiButton.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import MediaImage from '@/components/media/MediaImage.vue'
import ProjectCard from '@/modules/projects/ProjectCard.vue'
import { useAsync } from '@/composables/useAsync'
import { useSeo } from '@/composables/useSeo'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'
import { paragraphs, pad2 } from '@/utils/format'

const route = useRoute()
const slug = computed(() => String(route.params.slug))
const { data: service, loading, error, reload } = useAsync(() => publicApi.service(slug.value), { watch: [slug] })

// Shu xizmat bilan bog'liq loyihalar (services maydonidagi nom bo'yicha taxminiy moslik)
const { data: projects } = useAsync(() => publicApi.projects())
const related = computed(() => {
  const s = service.value
  if (!s || !projects.value) return []
  const key = s.title.toLowerCase().split(' ')[0].slice(0, 5)
  return projects.value.filter((p) => p.services.some((x) => x.toLowerCase().includes(key))).slice(0, 2)
})

useSeo(() => ({
  title: service.value?.seo.title || service.value?.title || 'Xizmat',
  description: service.value?.seo.description || service.value?.summary,
  noindex: error.value?.status === 404,
}))
</script>

<template>
  <div>
    <div v-if="loading && !service" class="container" style="padding-block: 96px" aria-busy="true">
      <div class="skeleton" style="height: 140px; width: 70%" />
      <div class="skeleton" style="height: 320px; margin-top: 64px" />
    </div>

    <div v-else-if="error" class="container section">
      <StateMessage
        tone="error"
        :title="error.status === 404 ? 'Xizmat topilmadi' : 'Yuklab bo‘lmadi'"
        :text="error.status === 404 ? undefined : error.message"
        :action-label="error.status === 404 ? undefined : 'Qayta urinish'"
        @action="reload"
      >
        <UiButton to="/services" variant="secondary" icon-left="arrow-left">Barcha xizmatlar</UiButton>
      </StateMessage>
    </div>

    <article v-else-if="service">
      <header class="container svc-hero">
        <RouterLink to="/services" class="svc-hero__back"><AppIcon name="arrow-left" :size="18" /> Xizmatlar</RouterLink>
        <h1 class="t-display" v-reveal>{{ service.title }}</h1>
        <p class="t-statement svc-hero__summary" v-reveal="80">{{ service.summary }}</p>
      </header>

      <div v-if="service.preview" class="container">
        <div class="svc-preview" v-reveal="{ variant: 'scale' }">
          <MediaImage :media="service.preview" fill priority />
        </div>
      </div>

      <section class="container section svc-body">
        <p class="t-label">Nima qilamiz</p>
        <div class="svc-body__main">
          <div class="svc-body__text">
            <p v-for="(p, i) in paragraphs(service.description)" :key="i">{{ p }}</p>
          </div>

          <ol v-if="service.features.length" role="list" class="svc-features">
            <li v-for="(f, i) in service.features" :key="f" v-reveal="i * 60">
              <span>{{ pad2(i + 1) }}</span>
              {{ f }}
            </li>
          </ol>

          <div v-if="service.stack.length" class="svc-stack">
            <p class="t-label">Texnologiyalar</p>
            <ul role="list">
              <li v-for="t in service.stack" :key="t">{{ t }}</li>
            </ul>
          </div>

          <UiButton to="/contact" icon="arrow-up-right">Loyihani muhokama qilish</UiButton>
        </div>
      </section>

      <section v-if="related.length" class="container section section--flush-top">
        <p class="t-label" style="margin-bottom: 48px">Shu yo‘nalishdagi loyihalar</p>
        <div class="svc-related">
          <ProjectCard v-for="p in related" :key="p.id" :project="p" size="half" />
        </div>
      </section>
    </article>
  </div>
</template>

<style scoped>
.svc-hero {
  padding-top: clamp(48px, 6vw, 96px);
  padding-bottom: clamp(56px, 6vw, 96px);
}

.svc-hero__back {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--ink-2);
  margin-bottom: 40px;
}

.svc-hero__summary {
  margin-top: 32px;
  max-width: 24ch;
  color: var(--ink-2);
}

.svc-preview {
  position: relative;
  height: clamp(360px, 44vw, 700px);
  border-radius: var(--r-lg);
  overflow: hidden;
}

.svc-body {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 3fr);
  gap: 48px;
}

.svc-body__main {
  display: grid;
  gap: 64px;
  justify-items: start;
}

.svc-body__text {
  display: grid;
  gap: 1em;
  font-size: var(--fs-body-lg);
  color: var(--ink-2);
  max-width: 62ch;
}

.svc-features {
  width: 100%;
  border-top: 1px solid var(--line-strong);
}

.svc-features li {
  display: flex;
  align-items: baseline;
  gap: 32px;
  padding: 28px 0;
  border-bottom: 1px solid var(--line);
  font-size: clamp(19px, 1.6vw, 25px);
  font-weight: 600;
  letter-spacing: -0.03em;
}

.svc-features span {
  font-size: var(--fs-small);
  color: var(--ink-3);
  letter-spacing: 0;
}

.svc-stack {
  display: grid;
  gap: 20px;
}

.svc-stack ul {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.svc-stack li {
  height: 46px;
  display: inline-flex;
  align-items: center;
  padding: 0 22px;
  border-radius: var(--r-pill);
  background: var(--surface);
  border: 1px solid var(--line);
  font-weight: 600;
  font-size: 16px;
}

.svc-related {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: clamp(24px, 2.4vw, 40px);
}

@media (max-width: 1024px) {
  .svc-body {
    grid-template-columns: 1fr;
    gap: 24px;
  }
}

@media (max-width: 768px) {
  .svc-related {
    grid-template-columns: 1fr;
    gap: 64px;
  }
}
</style>
