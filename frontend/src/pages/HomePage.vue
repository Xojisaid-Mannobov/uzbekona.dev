<script setup lang="ts">
import { computed } from 'vue'
import HeroSection from '@/modules/home/HeroSection.vue'
import PreviewSection from '@/modules/home/PreviewSection.vue'
import AboutStatement from '@/modules/home/AboutStatement.vue'
import SelectedProjects from '@/modules/home/SelectedProjects.vue'
import MetricsSection from '@/modules/home/MetricsSection.vue'
import ServicesSection from '@/modules/home/ServicesSection.vue'
import TechSection from '@/modules/home/TechSection.vue'
import ProcessSection from '@/modules/home/ProcessSection.vue'
import WhySection from '@/modules/home/WhySection.vue'
import TeamSection from '@/modules/home/TeamSection.vue'
import LabsSection from '@/modules/home/LabsSection.vue'
import JournalSection from '@/modules/home/JournalSection.vue'
import { useAsync } from '@/composables/useAsync'
import { useSeo } from '@/composables/useSeo'
import { publicApi } from '@/services/public'

// TZ 56-band: Hero → Preview → About → Projects → Metrics → Services → Tech → Process → Why → Team → Labs → Journal
const { data: projects, loading, error, reload } = useAsync(() => publicApi.projects({ featured: true, limit: 7 }))
const heroProject = computed(() => projects.value?.[0] ?? null)

useSeo({})
</script>

<template>
  <div>
    <HeroSection />
    <PreviewSection :project="heroProject" :loading="loading" />
    <AboutStatement />
    <SelectedProjects :projects="projects" :loading="loading" :error="error" @retry="reload" />
    <MetricsSection />
    <ServicesSection />
    <TechSection />
    <ProcessSection />
    <WhySection />
    <TeamSection />
    <LabsSection />
    <JournalSection />
  </div>
</template>
