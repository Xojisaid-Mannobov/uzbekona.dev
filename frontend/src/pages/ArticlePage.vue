<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import UiButton from '@/components/ui/UiButton.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import MediaImage from '@/components/media/MediaImage.vue'
import BlockRenderer from '@/components/content/BlockRenderer.vue'
import ArticleCard from '@/modules/journal/ArticleCard.vue'
import { useAsync } from '@/composables/useAsync'
import { useSeo } from '@/composables/useSeo'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'
import { formatDate } from '@/utils/format'
import { focusCss, ratioCss } from '@/utils/cover'

const route = useRoute()
const slug = computed(() => String(route.params.slug))
const { data: article, loading, error, reload } = useAsync(() => publicApi.article(slug.value), { watch: [slug] })

useSeo(() => ({
  title: article.value?.seo.title || article.value?.title || 'Journal',
  description: article.value?.seo.description || article.value?.excerpt,
  image: article.value?.cover?.url,
  type: 'article',
  noindex: error.value?.status === 404,
}))
</script>

<template>
  <div>
    <div v-if="loading && !article" class="container reading" style="padding-block: 96px" aria-busy="true">
      <div class="skeleton" style="height: 20px; width: 240px" />
      <div class="skeleton" style="height: 160px; margin-top: 32px" />
      <div class="skeleton" style="height: 420px; margin-top: 64px" />
    </div>

    <div v-else-if="error" class="container section">
      <StateMessage
        tone="error"
        :title="error.status === 404 ? 'Maqola topilmadi' : 'Yuklab bo‘lmadi'"
        :text="error.status === 404 ? undefined : error.message"
        :action-label="error.status === 404 ? undefined : 'Qayta urinish'"
        @action="reload"
      >
        <UiButton to="/journal" variant="secondary" icon-left="arrow-left">Journal’ga qaytish</UiButton>
      </StateMessage>
    </div>

    <article v-else-if="article">
      <header class="container article-hero">
        <div class="reading">
          <RouterLink to="/journal" class="back"><AppIcon name="arrow-left" :size="18" /> Journal</RouterLink>
          <p class="meta" v-reveal>
            <RouterLink v-if="article.category" :to="{ path: '/journal', query: { category: article.category.slug } }" class="meta__cat">
              {{ article.category.name }}
            </RouterLink>
            <time :datetime="article.published_at ?? undefined">{{ formatDate(article.published_at) }}</time>
            <span>{{ article.reading_time }} daqiqa o‘qish</span>
          </p>
          <h1 class="article-title" v-reveal="60">{{ article.title }}</h1>
          <p v-if="article.excerpt" class="t-lead article-excerpt" v-reveal="120">{{ article.excerpt }}</p>
          <p v-if="article.author_name" class="author">{{ article.author_name }}</p>
        </div>
      </header>

      <div v-if="article.cover" class="container">
        <div
          class="article-cover"
          :style="article.cover_ratio !== 'auto' ? { aspectRatio: ratioCss(article.cover_ratio, article.cover, '16 / 9') } : undefined"
          v-reveal="{ variant: 'scale' }"
        >
          <MediaImage
            :media="article.cover"
            priority
            sizes="(min-width: 1440px) 1360px, 100vw"
            :fill="article.cover_ratio !== 'auto'"
            :position="focusCss(article.cover_focus)"
          />
        </div>
      </div>

      <div class="container article-body">
        <div class="reading">
          <BlockRenderer :blocks="article.content ?? []" mode="article" />
        </div>
      </div>

      <section v-if="article.related.length" class="section related">
        <div class="container">
          <h2 class="t-h3 related__title">O‘qishni davom ettiring</h2>
          <div class="related__grid">
            <ArticleCard v-for="a in article.related" :key="a.id" :article="a" />
          </div>
        </div>
      </section>
    </article>
  </div>
</template>

<style scoped>
.reading {
  max-width: 820px;
  margin-inline: auto;
}

.article-hero {
  padding-top: clamp(48px, 6vw, 96px);
  padding-bottom: clamp(48px, 5vw, 80px);
}

.back {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--ink-2);
  margin-bottom: 48px;
}

.meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 20px;
  font-size: var(--fs-small);
  color: var(--ink-3);
  font-weight: 500;
}

.meta__cat {
  color: var(--accent);
  font-weight: 700;
}

.article-title {
  margin-top: 24px;
  font-size: clamp(32px, 3.8vw, 56px);
  line-height: 1.02;
  letter-spacing: -0.045em;
}

.article-excerpt {
  margin-top: 28px;
}

.author {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--line);
  font-weight: 600;
}

.article-cover {
  position: relative;
  max-height: 82vh;
  border-radius: var(--r-lg);
  overflow: hidden;
  margin-bottom: clamp(48px, 5vw, 80px);
}

.article-body {
  padding-bottom: clamp(80px, 8vw, 140px);
}

.related {
  border-top: 1px solid var(--line);
}

.related__title {
  margin-bottom: 56px;
}

.related__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 56px clamp(20px, 2.4vw, 40px);
}

@media (max-width: 1024px) {
  .related__grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 640px) {
  .related__grid {
    grid-template-columns: 1fr;
  }
}
</style>
