<script setup lang="ts">
import { computed, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import BrandIcon from '@/components/ui/BrandIcon.vue'
import UiButton from '@/components/ui/UiButton.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import BlockRenderer from '@/components/content/BlockRenderer.vue'
import GirihPattern from '@/components/ornament/GirihPattern.vue'
import DateLeaf from '@/modules/news/DateLeaf.vue'
import NewsCard from '@/modules/news/NewsCard.vue'
import NewsCover from '@/modules/news/NewsCover.vue'
import { toneStyle } from '@/modules/news/tone'
import { useAsync } from '@/composables/useAsync'
import { useSeo } from '@/composables/useSeo'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'
import { formatDate } from '@/utils/format'

const route = useRoute()
const slug = computed(() => String(route.params.slug))
const { data: item, loading, error, reload } = useAsync(() => publicApi.newsItem(slug.value), { watch: [slug] })

useSeo(() => ({
  title: item.value?.seo.title || item.value?.title || 'Yangiliklar',
  description: item.value?.seo.description || item.value?.excerpt,
  image: item.value?.cover?.url,
  type: 'article',
  noindex: error.value?.status === 404,
}))

// Ulashish: Telegram (O'zbekistonda asosiy messenjer) va havolani nusxalash
const pageUrl = () => window.location.origin + window.location.pathname
const telegramShare = computed(
  () => `https://t.me/share/url?url=${encodeURIComponent(pageUrl())}&text=${encodeURIComponent(item.value?.title ?? '')}`,
)
const copied = ref(false)
async function copyLink() {
  try {
    await navigator.clipboard.writeText(pageUrl())
    copied.value = true
    setTimeout(() => (copied.value = false), 2000)
  } catch {
    // Clipboard ruxsati yo'q — foydalanuvchi manzilni qo'lda nusxalaydi
  }
}
</script>

<template>
  <div>
    <div v-if="loading && !item" class="container reading" style="padding-block: 96px" aria-busy="true">
      <div class="skeleton" style="height: 20px; width: 240px" />
      <div class="skeleton" style="height: 160px; margin-top: 32px" />
      <div class="skeleton" style="height: 420px; margin-top: 64px" />
    </div>

    <div v-else-if="error" class="container section">
      <StateMessage
        tone="error"
        :title="error.status === 404 ? 'Yangilik topilmadi' : 'Yuklab bo‘lmadi'"
        :text="error.status === 404 ? undefined : error.message"
        :action-label="error.status === 404 ? undefined : 'Qayta urinish'"
        @action="reload"
      >
        <UiButton to="/news" variant="secondary" icon-left="arrow-left">Yangiliklarga qaytish</UiButton>
      </StateMessage>
    </div>

    <article v-else-if="item" :style="toneStyle(item.tag)">
      <header class="news-hero">
        <GirihPattern :color="'var(--tone)'" :size="80" :opacity="0.12" fade="right" />
        <div class="container reading news-hero__inner">
          <RouterLink to="/news" class="back"><AppIcon name="arrow-left" :size="18" /> Yangiliklar</RouterLink>
          <div class="news-hero__head" v-reveal>
            <DateLeaf :date="item.published_at" size="lg" class="news-hero__leaf" />
            <div class="news-hero__meta">
              <span v-if="item.tag" class="tag">{{ item.tag }}</span>
              <time :datetime="item.published_at ?? undefined">{{ formatDate(item.published_at) }}</time>
            </div>
          </div>
          <h1 class="news-title" v-reveal="60">{{ item.title }}</h1>
          <p v-if="item.excerpt" class="t-lead news-excerpt" v-reveal="120">{{ item.excerpt }}</p>

          <div class="share" v-reveal="160">
            <span class="share__label">Ulashish:</span>
            <a :href="telegramShare" target="_blank" rel="noopener noreferrer" class="share__btn">
              <BrandIcon name="telegram" :size="20" /> Telegram
            </a>
            <button type="button" class="share__btn" @click="copyLink">
              <AppIcon :name="copied ? 'check' : 'link'" :size="18" /> {{ copied ? 'Nusxalandi' : 'Havolani nusxalash' }}
            </button>
          </div>
        </div>
      </header>

      <div class="container">
        <NewsCover
          :news="item"
          class="news-cover"
          :fallback="item.cover ? '16 / 8' : '16 / 5.5'"
          sizes="(min-width: 1440px) 1360px, 100vw"
          priority
          v-reveal="{ variant: 'scale' }"
        />
      </div>

      <div class="container news-body">
        <div class="reading">
          <BlockRenderer :blocks="item.content ?? []" mode="article" />
        </div>
      </div>

      <section v-if="item.related.length" class="section related">
        <div class="container">
          <div class="related__head">
            <h2 class="t-h3">Boshqa yangiliklar</h2>
            <UiButton to="/news" variant="secondary" size="md" icon="arrow-up-right">Barchasi</UiButton>
          </div>
          <div class="related__grid">
            <NewsCard v-for="n in item.related" :key="n.id" :news="n" />
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

.news-hero {
  position: relative;
  overflow: hidden;
  isolation: isolate;
  padding-top: clamp(40px, 5vw, 80px);
  padding-bottom: clamp(40px, 4.4vw, 72px);
}

.news-hero__inner {
  position: relative;
  z-index: 1;
}

.back {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--ink-2);
  margin-bottom: 40px;
}

.news-hero__head {
  display: flex;
  align-items: center;
  gap: 20px;
}

.news-hero__leaf {
  --leaf-accent: var(--tone);
}

.news-hero__meta {
  display: grid;
  gap: 8px;
  font-size: var(--fs-small);
  font-weight: 500;
  color: var(--ink-3);
  justify-items: start;
}

.tag {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 5px 12px;
  border-radius: var(--r-pill);
  background: color-mix(in srgb, var(--tone) 14%, transparent);
  color: var(--tone);
  font-size: var(--fs-xs);
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.news-title {
  margin-top: 28px;
  font-size: clamp(32px, 3.8vw, 56px);
  line-height: 1.04;
  letter-spacing: -0.045em;
}

.news-excerpt {
  margin-top: 24px;
}

.share {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 10px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--line);
}

.share__label {
  font-size: var(--fs-small);
  font-weight: 600;
  color: var(--ink-3);
  margin-right: 4px;
}

.share__btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 42px;
  padding: 0 16px;
  border-radius: var(--r-pill);
  border: 1px solid var(--line-strong);
  background: var(--surface);
  font-size: var(--fs-small);
  font-weight: 600;
  transition:
    border-color var(--dur-fast) var(--ease),
    color var(--dur-fast) var(--ease);
}

.share__btn:hover {
  border-color: var(--tone);
  color: var(--tone);
}

.news-cover {
  max-height: 82vh;
  border-radius: var(--r-lg);
  margin-bottom: clamp(48px, 5vw, 80px);
}

.news-body {
  padding-bottom: clamp(72px, 7vw, 120px);
}

.related {
  border-top: 1px solid var(--line);
}

.related__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  margin-bottom: 48px;
}

.related__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 48px clamp(20px, 2.4vw, 40px);
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
