<script setup lang="ts">
import { RouterLink } from 'vue-router'
import MediaImage from '@/components/media/MediaImage.vue'
import GirihPattern from '@/components/ornament/GirihPattern.vue'
import type { Article } from '@/types/api'
import { formatDate } from '@/utils/format'

// Maqola kartasi: katta cover, kategoriya, sana, o'qish vaqti va 26–32px sarlavha
withDefaults(defineProps<{ article: Article; large?: boolean }>(), { large: false })
</script>

<template>
  <article class="article" :class="{ 'article--large': large }">
    <RouterLink :to="`/journal/${article.slug}`" class="article__link">
      <div class="article__cover">
        <MediaImage
          v-if="article.cover"
          :media="article.cover"
          :sizes="large ? '(min-width: 1024px) 900px, 100vw' : '(min-width: 1024px) 440px, 100vw'"
          fill
        />
        <div v-else class="article__typo" aria-hidden="true">
          <GirihPattern :size="72" :opacity="0.16" fade="top" />
          <span class="article__typo-meta">{{ article.reading_time }} daqiqa o‘qish</span>
          <span class="article__typo-word">{{ article.category?.name ?? 'Journal' }}</span>
        </div>
      </div>
      <p class="article__meta">
        <span v-if="article.category" class="article__cat">{{ article.category.name }}</span>
        <time :datetime="article.published_at ?? undefined">{{ formatDate(article.published_at) }}</time>
        <span>{{ article.reading_time }} daqiqa</span>
      </p>
      <h3 class="article__title">{{ article.title }}</h3>
      <p v-if="large && article.excerpt" class="article__excerpt">{{ article.excerpt }}</p>
    </RouterLink>
  </article>
</template>

<style scoped>
.article__link {
  display: block;
}

.article__cover {
  position: relative;
  aspect-ratio: 4 / 3.2;
  border-radius: var(--r-lg);
  overflow: hidden;
  background: var(--surface-2);
  margin-bottom: 24px;
}

.article--large .article__cover {
  aspect-ratio: 16 / 10;
}

.article__cover > :deep(*) {
  transition: transform 1s var(--ease);
}

.article__link:hover .article__cover > :deep(*) {
  transform: scale(1.03);
}

/* Cover bo'lmasa — tipografik muqova */
.article__typo {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: clamp(24px, 2.6vw, 40px);
  background: radial-gradient(100% 80% at 100% 0%, var(--accent-soft), transparent 65%), var(--surface);
  border: 1px solid var(--line);
  border-radius: inherit;
}

.article__typo-meta {
  position: relative;
  align-self: flex-end;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--ink-3);
}

/* Kategoriya nomi — katta, pastki chetdan biroz kesilgan tipografik kompozitsiya */
.article__typo-word {
  position: relative;
  font-size: clamp(44px, 4.8vw, 76px);
  line-height: 0.8;
  font-weight: 700;
  letter-spacing: -0.06em;
  color: var(--accent);
  margin-bottom: -0.16em;
  white-space: nowrap;
}

.article--large .article__typo-word {
  font-size: clamp(56px, 8vw, 128px);
}

.article__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 6px 16px;
  font-size: var(--fs-small);
  color: var(--ink-3);
  font-weight: 500;
}

.article__cat {
  color: var(--ink);
  font-weight: 700;
}

.article__title {
  margin-top: 12px;
  font-size: clamp(19px, 1.5vw, 24px);
  letter-spacing: -0.035em;
  line-height: 1.15;
  transition: color var(--dur-fast) var(--ease);
}

.article--large .article__title {
  font-size: clamp(24px, 2.2vw, 34px);
}

.article__link:hover .article__title {
  color: var(--accent);
}

.article__excerpt {
  margin-top: 16px;
  max-width: 60ch;
  font-size: var(--fs-body-lg);
  color: var(--ink-2);
}
</style>
