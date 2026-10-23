<script setup lang="ts">
import SectionHeader from '@/components/ui/SectionHeader.vue'
import UiButton from '@/components/ui/UiButton.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import ProjectCard from '@/modules/projects/ProjectCard.vue'
import { vReveal } from '@/composables/reveal'
import type { Project } from '@/types/api'
import type { ApiError } from '@/services/http'

defineProps<{ projects: Project[] | null; loading: boolean; error: ApiError | null }>()
defineEmits<{ retry: [] }>()
</script>

<template>
  <section id="projects" class="section projects" aria-labelledby="projects-title">
    <div class="container">
      <SectionHeader title-id="projects-title" index="01" label="Tanlangan loyihalar" title="Biz qurgan va bugun ishlayotgan mahsulotlar.">
        <template #aside>
          <UiButton to="/projects" variant="secondary" icon="arrow-up-right">Barcha loyihalar</UiButton>
        </template>
      </SectionHeader>

      <div v-if="loading && !projects" class="projects__grid" aria-busy="true">
        <div v-for="n in 3" :key="n" class="skeleton projects__skeleton" />
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

      <div v-else class="projects__grid">
        <div v-for="(p, i) in projects" :key="p.id" v-reveal="(i % 3) * 80">
          <ProjectCard :project="p" size="compact" />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.projects__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: clamp(40px, 4vw, 56px) clamp(20px, 2vw, 28px);
}

.projects__skeleton {
  aspect-ratio: 4 / 3;
  border-radius: var(--r-lg);
}

@media (max-width: 1024px) {
  .projects__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .projects__grid {
    grid-template-columns: 1fr;
  }
}
</style>
