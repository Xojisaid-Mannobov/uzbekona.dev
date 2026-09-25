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

// Platforma bo'yicha filtr (Web, Mobile, Telegram …) — ma'lumotlardan avtomatik
const filter = ref('all')
const platforms = computed(() => {
  const set = new Set<string>()
  projects.value?.forEach((p) => p.platforms.forEach((x) => set.add(x)))
  return [...set]
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
      lead="Har bir loyiha — real foydalanuvchilar ishlatadigan tizim. Case-study’larda muammo, yechim va natijani ko‘rsatamiz."
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
          :key="p"
          type="button"
          class="filter"
          :class="{ 'is-active': filter === p }"
          :aria-pressed="filter === p"
          @click="filter = p"
        >
          {{ p }}
        </button>
      </div>
    </PageHero>

    <section class="section section--flush-top">
      <div class="container">
        <div v-if="loading && !projects" class="grid">
          <div v-for="n in 2" :key="n" class="skeleton grid__skeleton" />
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
          <div v-for="(p, i) in visible" :key="p.id" :class="{ grid__wide: i % 3 === 0 }" v-reveal="(i % 3 === 2 ? 1 : 0) * 120">
            <ProjectCard :project="p" :size="i % 3 === 0 ? 'wide' : 'half'" :priority="i === 0" />
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
  grid-template-columns: 1fr 1fr;
  gap: clamp(72px, 8vw, 120px) clamp(24px, 2.4vw, 40px);
}

.grid__wide {
  grid-column: 1 / -1;
}

.grid__skeleton {
  height: 600px;
  border-radius: var(--r-lg);
}

@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
