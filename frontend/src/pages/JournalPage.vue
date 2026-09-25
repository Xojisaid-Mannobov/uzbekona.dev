<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import PageHero from '@/components/layout/PageHero.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import UiButton from '@/components/ui/UiButton.vue'
import ArticleCard from '@/modules/journal/ArticleCard.vue'
import { useAsync } from '@/composables/useAsync'
import { useSeo } from '@/composables/useSeo'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'
import type { Article } from '@/types/api'

const route = useRoute()
const router = useRouter()

// Kategoriya URL query'da saqlanadi (?category=engineering) — havolani ulashish mumkin
const category = computed(() => (typeof route.query.category === 'string' ? route.query.category : ''))
const page = ref(1)
const items = ref<Article[]>([])

const { data: categories } = useAsync(publicApi.categories)
const { data, loading, error, reload } = useAsync(
  () => publicApi.articles({ category: category.value || undefined, page: page.value, limit: 9 }),
  { watch: [category, page] },
)

watch(category, () => {
  page.value = 1
  items.value = []
})
watch(data, (d) => {
  if (!d) return
  items.value = d.meta.page === 1 ? d.items : [...items.value, ...d.items]
})

const hasMore = computed(() => !!data.value && data.value.meta.page * data.value.meta.limit < data.value.meta.total)
const featured = computed(() => (!category.value && items.value.length ? items.value[0] : null))
const rest = computed(() => (featured.value ? items.value.slice(1) : items.value))

function select(slug: string) {
  router.replace({ query: slug ? { category: slug } : {} })
}

useSeo({ title: 'Journal', description: 'Engineering, product, dizayn, AI va infratuzilma haqida Uzbekona.dev jamoasi tajribasi.' })
</script>

<template>
  <div>
    <PageHero
      label="Journal"
      title="Tajriba va kuzatuvlar."
      lead="Loyihalar davomida o‘rganganlarimiz: arxitektura qarorlari, mahsulot tajribasi va texnologiyalar."
    >
      <div class="cats" role="group" aria-label="Kategoriyalar">
        <button type="button" class="cat" :class="{ 'is-active': !category }" :aria-pressed="!category" @click="select('')">
          Barchasi
        </button>
        <button
          v-for="c in categories"
          :key="c.id"
          type="button"
          class="cat"
          :class="{ 'is-active': category === c.slug }"
          :aria-pressed="category === c.slug"
          @click="select(c.slug)"
        >
          {{ c.name }}
          <span>{{ c.count }}</span>
        </button>
      </div>
    </PageHero>

    <section class="section section--flush-top">
      <div class="container">
        <div v-if="loading && !items.length" class="grid">
          <div v-for="n in 3" :key="n" class="skeleton sk" />
        </div>

        <StateMessage
          v-else-if="error && !items.length"
          tone="error"
          title="Maqolalarni yuklab bo‘lmadi"
          :text="error.message"
          action-label="Qayta urinish"
          @action="reload"
        />

        <StateMessage
          v-else-if="!items.length"
          title="Bu kategoriyada hozircha maqola yo‘q"
          text="Boshqa kategoriyani tanlang yoki keyinroq qaytib keling."
        >
          <UiButton variant="secondary" size="md" @click="select('')">Barcha maqolalar</UiButton>
        </StateMessage>

        <template v-else>
          <div v-if="featured" class="featured" v-reveal>
            <ArticleCard :article="featured" large />
          </div>
          <div class="grid">
            <div v-for="(a, i) in rest" :key="a.id" v-reveal="(i % 3) * 80">
              <ArticleCard :article="a" />
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
.cats {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.cat {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 52px;
  padding: 0 22px;
  border-radius: var(--r-pill);
  border: 1px solid var(--line-strong);
  font-weight: 600;
  transition:
    background-color var(--dur-fast) var(--ease),
    color var(--dur-fast) var(--ease),
    border-color var(--dur-fast) var(--ease);
}

.cat span {
  font-size: 13px;
  color: var(--ink-3);
}

.cat:hover {
  border-color: var(--ink);
}

.cat.is-active {
  background: var(--ink);
  color: var(--bg);
  border-color: var(--ink);
}

.featured {
  margin-bottom: clamp(72px, 8vw, 112px);
  max-width: 1100px;
}

.grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 72px clamp(20px, 2.4vw, 40px);
}

.sk {
  aspect-ratio: 4 / 4.4;
  border-radius: var(--r-lg);
}

.more {
  display: flex;
  justify-content: center;
  margin-top: 80px;
}

@media (max-width: 1024px) {
  .grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
