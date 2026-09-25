<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import AboutStatement from '@/modules/home/AboutStatement.vue'
import MetricsSection from '@/modules/home/MetricsSection.vue'
import WhySection from '@/modules/home/WhySection.vue'
import ProcessSection from '@/modules/home/ProcessSection.vue'
import LabsSection from '@/modules/home/LabsSection.vue'
import { vReveal } from '@/composables/reveal'
import { useSeo } from '@/composables/useSeo'

const values = [
  { title: 'Aniqlik', text: 'Nima qilayotganimiz, qancha vaqt ketishi va nima uchun — har doim ochiq. Yashirin ish yo‘q.' },
  { title: 'Mas’uliyat', text: 'Kod production’da ishlashi bizning vazifamiz. Kelishuv asosida launch’dan keyin ham yoningizdamiz.' },
  { title: 'Soddalik', text: 'Murakkab muammoga eng sodda ishlaydigan yechim. Ortiqcha texnologiya — ortiqcha xavf.' },
]

useSeo({
  title: 'Biz haqimizda',
  description: 'Uzbekona.dev — real tizimlar ishlab chiqadigan texnologik jamoa. Jarayonlarni tushunamiz va raqamlashtiramiz.',
})
</script>

<template>
  <div>
    <PageHero
      label="Biz haqimizda"
      :title="'Real tizimlar quradigan\ntexnologik jamoa.'"
      lead="Uzbekona.dev — digital product studio. Biz universitetlar, kompaniyalar va startaplar uchun har kuni ishlatiladigan mahsulotlar yaratamiz."
    />

    <AboutStatement />
    <MetricsSection />

    <section class="section values" aria-labelledby="values-title">
      <div class="container values__grid">
        <h2 id="values-title" class="t-h2" v-reveal>Qadriyatlarimiz</h2>
        <div class="values__list">
          <div v-for="(v, i) in values" :key="v.title" class="values__item" v-reveal="i * 80">
            <h3 class="t-h3">{{ v.title }}</h3>
            <p>{{ v.text }}</p>
          </div>
        </div>
      </div>
    </section>

    <WhySection index="01" />
    <ProcessSection index="02" />
    <LabsSection />
  </div>
</template>

<style scoped>
.values__grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 2fr);
  gap: 48px;
}

.values__list {
  border-top: 1px solid var(--line-strong);
}

.values__item {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1.4fr);
  gap: 32px;
  padding-block: 40px;
  border-bottom: 1px solid var(--line);
}

.values__item p {
  font-size: var(--fs-body-lg);
  color: var(--ink-2);
}

@media (max-width: 1024px) {
  .values__grid,
  .values__item {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}
</style>
