<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import ServicesSection from '@/modules/home/ServicesSection.vue'
import ProcessSection from '@/modules/home/ProcessSection.vue'
import TechSection from '@/modules/home/TechSection.vue'
import { useSeo } from '@/composables/useSeo'

useSeo({
  title: 'Xizmatlar',
  description: 'Web platformalar, mobil ilovalar, Telegram tizimlari, avtomatlashtirish, AI integratsiyalar va infratuzilma.',
})
</script>

<template>
  <div>
    <PageHero
      label="Xizmatlar"
      :title="'Biznes uchun\nraqamli mahsulotlar.'"
      lead="Bitta jamoa — tadqiqotdan production’gacha. Kerakli yo‘nalishni tanlang yoki vazifangizni aytib bering, yechimni birga loyihalaymiz."
    />
    <ServicesSection :show-header="false" class="section--flush-top" />
    <ProcessSection index="01" />
    <TechSection index="02" />
  </div>
</template>
