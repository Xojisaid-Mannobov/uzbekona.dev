<script setup lang="ts">
import HeroSection from '@/modules/home/HeroSection.vue'
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
import NewsSection from '@/modules/home/NewsSection.vue'
import NewsTicker from '@/modules/news/NewsTicker.vue'
import GhostMarquee from '@/components/ornament/GhostMarquee.vue'
import { useAsync } from '@/composables/useAsync'
import { useSeo } from '@/composables/useSeo'
import { publicApi } from '@/services/public'

// Tartib ishontirish mantig'i bo'yicha: kimmiz → isbot (loyihalar, raqamlar) → nega biz → nima qilamiz → qanday → vositalar.
// Preview bloki olib tashlangan: hero'da illyustratsiya bor, birinchi loyiha esa "Tanlangan loyihalar"da darhol ko'rinadi.
const { data: projects, loading, error, reload } = useAsync(() => publicApi.projects({ featured: true, limit: 7 }))

useSeo({})
</script>

<template>
  <div>
    <HeroSection />
    <NewsTicker />
    <AboutStatement />
    <SelectedProjects :projects="projects" :loading="loading" :error="error" @retry="reload" />
    <MetricsSection />
    <WhySection index="02" />
    <ServicesSection index="03" />
    <GhostMarquee
      :rows="[
        ['Web platformalar', 'Mobil ilovalar', 'Telegram tizimlari', 'AI integratsiya'],
        ['Avtomatlashtirish', 'Product design', 'Infratuzilma', 'Raqamli O‘zbekiston'],
      ]"
    />
    <ProcessSection index="04" />
    <TechSection index="05" />
    <TeamSection />
    <NewsSection index="07" />
    <LabsSection index="08" />
    <JournalSection />
  </div>
</template>
