<script setup lang="ts">
import { computed } from 'vue'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import UiButton from '@/components/ui/UiButton.vue'
import NewsCard from '@/modules/news/NewsCard.vue'
import { useAsync } from '@/composables/useAsync'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'

// Bosh sahifadagi yangiliklar: chapda asosiy (mahkamlangan) xabar, o'ngda taqvim varaqli ro'yxat
withDefaults(defineProps<{ index?: string }>(), { index: '07' })

const { data, loading } = useAsync(() => publicApi.news({ limit: 4 }))
const featured = computed(() => data.value?.items[0] ?? null)
const rest = computed(() => data.value?.items.slice(1) ?? [])
</script>

<template>
  <section v-if="loading || featured" class="section news-home" aria-labelledby="news-title">
    <div class="container">
      <SectionHeader title-id="news-title" :index="index" label="Yangiliklar" title="Studiyada nimalar bo‘lyapti?">
        <template #aside>
          <UiButton to="/news" variant="secondary" icon="arrow-up-right">Barcha yangiliklar</UiButton>
        </template>
      </SectionHeader>

      <div v-if="loading && !data" class="news-home__grid" aria-busy="true">
        <div class="skeleton" style="aspect-ratio: 16 / 11" />
        <div class="skeleton" style="height: 280px" />
      </div>

      <div v-else-if="featured" class="news-home__grid">
        <NewsCard :news="featured" v-reveal />
        <div class="news-home__list" v-reveal="120">
          <NewsCard v-for="n in rest" :key="n.id" :news="n" variant="row" />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.news-home__grid {
  display: grid;
  grid-template-columns: minmax(0, 1.15fr) minmax(0, 1fr);
  gap: clamp(28px, 4vw, 72px);
  align-items: start;
}

.news-home__list {
  border-top: 1px solid var(--line-strong);
}

@media (max-width: 900px) {
  .news-home__grid {
    grid-template-columns: 1fr;
  }
}
</style>
