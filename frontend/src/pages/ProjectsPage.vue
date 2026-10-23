<script setup lang="ts">
import { computed, ref } from 'vue'
import PageHero from '@/components/layout/PageHero.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import ProjectCard from '@/modules/projects/ProjectCard.vue'
import { useAsync } from '@/composables/useAsync'
import { useSeo } from '@/composables/useSeo'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'

const { data: projects, loading, error, reload } = useAsync(() => publicApi.projects())

// Platforma bo'yicha filtr (Web, Telegram …) — ma'lumotlardan avtomatik, eng ko'pidan boshlab
const filter = ref('all')
const platforms = computed(() => {
  const counts = new Map<string, number>()
  projects.value?.forEach((p) => p.platforms.forEach((x) => counts.set(x, (counts.get(x) ?? 0) + 1)))
  return [...counts].sort((a, b) => b[1] - a[1]).map(([name, count]) => ({ name, count }))
})
const visible = computed(() =>
  filter.value === 'all' ? (projects.value ?? []) : (projects.value ?? []).filter((p) => p.platforms.includes(filter.value)),
)

useSeo({
  title: 'Loyihalar',
  description: 'Uzbekona.dev tomonidan ishlab chiqilgan web platformalar, mobil ilovalar va Telegram tizimlari.',
})
</script>

<template>
  <div>
    <PageHero
      label="Loyihalar"
      title="Ishlagan mahsulotlarimiz."
      lead="Web platformalar, ichki boshqaruv tizimlari va Telegram botlar — hammasi bugun real foydalanuvchilar uchun ishlab turibdi."
    >
      <div v-if="platforms.length > 1" class="filters" role="group" aria-label="Platforma bo‘yicha filtr">
        <button
          type="button"
          class="filter"
          :class="{ 'is-active': filter === 'all' }"
          :aria-pressed="filter === 'all'"
          @click="filter = 'all'"
        >
          Barchasi <span>{{ projects?.length }}</span>
        </button>
        <button
          v-for="p in platforms"
          :key="p.name"
          type="button"
          class="filter"
          :class="{ 'is-active': filter === p.name }"
          :aria-pressed="filter === p.name"
          @click="filter = p.name"
        >
          {{ p.name }} <span>{{ p.count }}</span>
        </button>
      </div>
    </PageHero>

    <section class="section section--flush-top">
      <div class="container">
        <div v-if="loading && !projects" class="grid">
          <div v-for="n in 6" :key="n" class="skeleton grid__skeleton" />
        </div>
        <StateMessage
          v-else-if="error"
          tone="error"
          title="Loyihalarni yuklab bo‘lmadi"
          :text="error.message"
          action-label="Qayta urinish"
          @action="reload"
        />
        <StateMessage
          v-else-if="!visible.length"
          title="Hozircha loyiha yo‘q"
          text="Bu yo‘nalishdagi case-study’lar tez orada qo‘shiladi."
        />
        <div v-else class="grid">
          <div v-for="(p, i) in visible" :key="p.id" v-reveal="(i % 3) * 80">
            <ProjectCard :project="p" size="compact" :priority="i < 3" />
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.filter {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 46px;
  padding: 0 22px;
  border-radius: var(--r-pill);
  background: var(--bg);
  border: 1px solid var(--line-strong);
  font-weight: 600;
  transition:
    background-color var(--dur-fast) var(--ease),
    color var(--dur-fast) var(--ease),
    border-color var(--dur-fast) var(--ease);
}

.filter span {
  font-size: 13px;
  color: var(--ink-3);
}

.filter:hover {
  border-color: var(--ink);
}

.filter.is-active {
  background: var(--ink);
  border-color: var(--ink);
  color: var(--bg);
}

.grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: clamp(40px, 4vw, 56px) clamp(20px, 2vw, 28px);
}

.grid__skeleton {
  aspect-ratio: 4 / 3;
  border-radius: var(--r-lg);
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
