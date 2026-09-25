<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import PageHero from '@/components/layout/PageHero.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import UiButton from '@/components/ui/UiButton.vue'
import NewsCard from '@/modules/news/NewsCard.vue'
import { useAsync } from '@/composables/useAsync'
import { useSeo } from '@/composables/useSeo'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'
import type { News } from '@/types/api'

const page = ref(1)
const items = ref<News[]>([])
const { data, loading, error, reload } = useAsync(() => publicApi.news({ page: page.value, limit: 10 }), { watch: [page] })

watch(data, (d) => {
  if (!d) return
  items.value = d.meta.page === 1 ? d.items : [...items.value, ...d.items]
})

const hasMore = computed(() => !!data.value && data.value.meta.page * data.value.meta.limit < data.value.meta.total)
// Birinchi (mahkamlangan yoki eng yangi) xabar — keng karta
const featured = computed(() => items.value[0] ?? null)
const rest = computed(() => items.value.slice(1))

useSeo({
  title: 'Yangiliklar',
  description: 'Uzbekona.dev yangiliklari: yangi loyihalar, ochiq kodli vositalar, jamoa va hamkorliklar haqida so‘nggi xabarlar.',
})
</script>

<template>
  <div>
    <PageHero
      label="Yangiliklar"
      :title="'Studiya hayotidan\nso‘nggi xabarlar.'"
      lead="Yangi loyihalar, ochiq kodli vositalar, jamoa va hamkorliklar — Uzbekona.dev’da nimalar bo‘layotganini birinchi bo‘lib biling."
    />

    <section class="section section--flush-top">
      <div class="container">
        <div v-if="loading && !items.length" class="news-grid" aria-busy="true">
          <div v-for="n in 3" :key="n" class="skeleton sk" />
        </div>

        <StateMessage
          v-else-if="error && !items.length"
          tone="error"
          title="Yangiliklarni yuklab bo‘lmadi"
          :text="error.message"
          action-label="Qayta urinish"
          @action="reload"
        />

        <StateMessage v-else-if="!items.length" title="Yangiliklar tez orada" text="Birinchi xabarlarimizni tayyorlayapmiz." />

        <template v-else>
          <div v-if="featured" class="featured" v-reveal>
            <NewsCard :news="featured" variant="feature" :level="2" />
          </div>

          <h2 v-if="rest.length" class="t-label divider">Barcha xabarlar</h2>
          <div class="news-grid">
            <div v-for="(n, i) in rest" :key="n.id" v-reveal="(i % 3) * 80">
              <NewsCard :news="n" />
            </div>
          </div>
          <div v-if="hasMore" class="more">
            <UiButton variant="secondary" :loading="loading" @click="page++">Ko‘proq yuklash</UiButton>
          </div>
        </template>
      </div>
    </section>
  </div>
</template>

<style scoped>
.featured {
  padding-bottom: clamp(48px, 5vw, 80px);
}

.divider {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 40px;
}

.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--line-strong);
}

.news-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 56px clamp(20px, 2.4vw, 40px);
}

.sk {
  aspect-ratio: 16 / 11;
  border-radius: var(--r-lg);
}

.more {
  display: flex;
  justify-content: center;
  margin-top: 64px;
}

@media (max-width: 1024px) {
  .news-grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 640px) {
  .news-grid {
    grid-template-columns: 1fr;
    gap: 48px;
  }
}
</style>
