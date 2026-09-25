<script setup lang="ts">
import SectionHeader from '@/components/ui/SectionHeader.vue'
import UiButton from '@/components/ui/UiButton.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import ArticleCard from '@/modules/journal/ArticleCard.vue'
import { useAsync } from '@/composables/useAsync'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'

const { data, loading, error, reload } = useAsync(() => publicApi.articles({ limit: 3 }))
</script>

<template>
  <section v-if="loading || error || data?.items.length" class="section journal" aria-labelledby="journal-title">
    <div class="container">
      <SectionHeader title-id="journal-title" index="08" label="Journal" title="Bilganimizni yashirmaymiz.">
        <template #aside>
          <UiButton to="/journal" variant="secondary" icon="arrow-up-right">Barcha maqolalar</UiButton>
        </template>
      </SectionHeader>

      <div v-if="loading && !data" class="journal__grid">
        <div v-for="n in 3" :key="n" class="skeleton journal__skeleton" />
      </div>

      <StateMessage
        v-else-if="error"
        tone="error"
        title="Maqolalarni yuklab bo‘lmadi"
        :text="error.message"
        action-label="Qayta urinish"
        @action="reload"
      />

      <div v-else class="journal__grid">
        <div v-for="(a, i) in data?.items" :key="a.id" v-reveal="i * 90">
          <ArticleCard :article="a" />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.journal__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 56px clamp(20px, 2.4vw, 40px);
}

.journal__skeleton {
  aspect-ratio: 4 / 4.4;
  border-radius: var(--r-lg);
}

@media (max-width: 1024px) {
  .journal__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .journal__grid {
    grid-template-columns: 1fr;
  }
}
</style>
