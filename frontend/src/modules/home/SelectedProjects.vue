<script setup lang="ts">
import { computed } from 'vue'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import UiButton from '@/components/ui/UiButton.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import ProjectCard from '@/modules/projects/ProjectCard.vue'
import { vReveal } from '@/composables/reveal'
import type { Project } from '@/types/api'
import type { ApiError } from '@/services/http'

const props = defineProps<{ projects: Project[] | null; loading: boolean; error: ApiError | null }>()
defineEmits<{ retry: [] }>()

// Ritm: 1 ta keng → 2 ta yarim → 1 ta keng …
const rows = computed(() => {
  const list = props.projects ?? []
  const out: { wide?: Project; pair?: Project[] }[] = []
  let i = 0
  while (i < list.length) {
    out.push({ wide: list[i] })
    i++
    if (i < list.length) {
      out.push({ pair: list.slice(i, i + 2) })
      i += 2
    }
  }
  return out
})
</script>

<template>
  <section id="projects" class="section projects" aria-labelledby="projects-title">
    <div class="container">
      <SectionHeader title-id="projects-title" index="01" label="Tanlangan loyihalar" title="Ishlagan mahsulotlarimiz.">
        <template #aside>
          <UiButton to="/projects" variant="secondary" icon="arrow-up-right">Barcha loyihalar</UiButton>
        </template>
      </SectionHeader>

      <div v-if="loading && !projects" class="projects__list" aria-busy="true">
        <div class="skeleton projects__skeleton" />
      </div>

      <StateMessage
        v-else-if="error"
        tone="error"
        title="Loyihalarni yuklab bo‘lmadi"
        :text="error.message"
        action-label="Qayta urinish"
        @action="$emit('retry')"
      />

      <StateMessage
        v-else-if="!projects?.length"
        title="Loyihalar tez orada"
        text="Case-study’lar tayyorlanmoqda. Hozircha biz bilan bevosita bog‘laning."
      />

      <div v-else class="projects__list">
        <template v-for="(row, i) in rows" :key="i">
          <div v-if="row.wide" v-reveal>
            <ProjectCard :project="row.wide" size="wide" />
          </div>
          <div v-else class="projects__pair">
            <div v-for="(p, j) in row.pair" :key="p.id" v-reveal="j * 120">
              <ProjectCard :project="p" size="half" />
            </div>
          </div>
        </template>
      </div>
    </div>
  </section>
</template>

<style scoped>
.projects__list {
  display: grid;
  gap: clamp(72px, 8vw, 120px);
}

.projects__pair {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: clamp(24px, 2.4vw, 40px);
}

.projects__skeleton {
  height: 650px;
  border-radius: var(--r-lg);
}

@media (max-width: 768px) {
  .projects__pair {
    grid-template-columns: 1fr;
    gap: 72px;
  }

  .projects__skeleton {
    height: 420px;
  }
}
</style>
