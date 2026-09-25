<script setup lang="ts">
import { RouterLink } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import DateLeaf from './DateLeaf.vue'
import NewsCover from './NewsCover.vue'
import type { News } from '@/types/api'
import { toneStyle } from './tone'

// Yangilik kartasi: muqova ustida taqvim varag'i, rangli teg, sarlavha.
// feature — keng gorizontal karta (asosiy yangilik), row — ixcham ro'yxat qatori.
// level — sarlavha darajasi (sahifadagi h1/h2 tartibi buzilmasligi uchun)
withDefaults(defineProps<{ news: News; variant?: 'card' | 'feature' | 'row'; level?: 2 | 3 }>(), { variant: 'card', level: 3 })
</script>

<template>
  <article class="news" :class="`news--${variant}`" :style="toneStyle(news.tag)">
    <RouterLink :to="`/news/${news.slug}`" class="news__link">
      <template v-if="variant === 'row'">
        <DateLeaf :date="news.published_at" size="sm" class="news__leaf" />
        <div class="news__body">
          <span v-if="news.tag" class="news__tag">{{ news.tag }}</span>
          <component :is="`h${level}`" class="news__title">{{ news.title }}</component>
        </div>
        <span class="news__arrow" aria-hidden="true"><AppIcon name="arrow-up-right" :size="20" /></span>
      </template>

      <template v-else>
        <div class="news__media">
          <NewsCover
            :news="news"
            class="news__cover"
            :fallback="variant === 'feature' ? '16 / 10.5' : '16 / 11'"
            :sizes="variant === 'feature' ? '(min-width: 1024px) 760px, 100vw' : '(min-width: 1024px) 440px, 100vw'"
          />
          <DateLeaf :date="news.published_at" :size="variant === 'feature' ? 'lg' : 'md'" class="news__leaf" />
        </div>
        <div class="news__body">
          <span v-if="news.tag" class="news__tag">{{ news.tag }}</span>
          <component :is="`h${level}`" class="news__title">{{ news.title }}</component>
          <p v-if="news.excerpt" class="news__excerpt">{{ news.excerpt }}</p>
          <span v-if="variant === 'feature'" class="news__more">O‘qish <AppIcon name="arrow-up-right" :size="18" /></span>
        </div>
      </template>
    </RouterLink>
  </article>
</template>

<style scoped>
.news__link {
  display: block;
  height: 100%;
}

.news__media {
  position: relative;
}

.news__cover {
  border-radius: var(--r-lg);
  max-height: min(78vh, 720px);
}

.news__cover :deep(img),
.news__cover :deep(.cover__poster) {
  transition: transform 1s var(--ease);
}

.news__link:hover .news__cover :deep(img),
.news__link:hover .news__cover :deep(.cover__poster) {
  transform: scale(1.04);
}

.news__media .news__leaf {
  position: absolute;
  left: 18px;
  top: 18px;
  --leaf-accent: var(--tone);
}

.news__body {
  display: grid;
  gap: 10px;
  padding-top: 22px;
  justify-items: start;
}

.news__tag {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: var(--fs-xs);
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--tone);
}

.news__tag::before {
  content: '';
  width: 7px;
  height: 7px;
  rotate: 45deg;
  background: currentColor;
}

.news__title {
  font-size: var(--fs-card);
  font-weight: 600;
  line-height: 1.18;
  letter-spacing: -0.03em;
  transition: color var(--dur-fast) var(--ease);
}

.news__link:hover .news__title {
  color: var(--accent);
}

.news__excerpt {
  color: var(--ink-2);
  line-height: 1.55;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* ─── Keng (asosiy) karta ─── */
.news--feature .news__link {
  display: grid;
  grid-template-columns: minmax(0, 1.35fr) minmax(0, 1fr);
  gap: clamp(28px, 3.4vw, 56px);
  align-items: center;
}

.news--feature .news__body {
  padding-top: 0;
  gap: 16px;
}

.news--feature .news__title {
  font-size: clamp(26px, 2.6vw, 40px);
  line-height: 1.08;
  letter-spacing: -0.04em;
}

.news--feature .news__excerpt {
  font-size: var(--fs-body-lg);
  -webkit-line-clamp: 4;
}

.news__more {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-top: 8px;
  padding-bottom: 4px;
  font-weight: 600;
  border-bottom: 1px solid var(--line-strong);
}

/* ─── Ro'yxat qatori ─── */
.news--row .news__link {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: 20px;
  align-items: center;
  padding: 20px 0;
  border-bottom: 1px solid var(--line);
}

.news--row .news__leaf {
  --leaf-accent: var(--tone);
}

.news--row .news__body {
  padding-top: 0;
  gap: 6px;
}

.news--row .news__title {
  font-size: clamp(17px, 1.3vw, 20px);
  line-height: 1.3;
}

.news__arrow {
  display: grid;
  place-items: center;
  width: 42px;
  height: 42px;
  border-radius: 50%;
  border: 1px solid var(--line-strong);
  transition:
    background-color var(--dur-fast) var(--ease),
    color var(--dur-fast) var(--ease),
    transform var(--dur) var(--ease);
}

.news--row .news__link:hover .news__arrow {
  background: var(--accent);
  border-color: var(--accent);
  color: var(--on-accent);
  transform: rotate(45deg);
}

@media (max-width: 860px) {
  .news--feature .news__link {
    grid-template-columns: 1fr;
    gap: 22px;
  }
}
</style>
