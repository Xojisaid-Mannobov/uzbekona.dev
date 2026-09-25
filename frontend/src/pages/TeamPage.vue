<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import TeamGrid from '@/modules/team/TeamGrid.vue'
import { useSeo } from '@/composables/useSeo'

useSeo({ title: 'Jamoa', description: 'Uzbekona.dev mahsulotlarini quradigan engineer va dizaynerlar.' })
</script>

<template>
  <div>
    <PageHero
      label="Jamoa"
      :title="'Mahsulotni quradigan\nodamlar.'"
      lead="Kichik, lekin to‘liq jamoa: product, dizayn, frontend, backend va infratuzilma. Har bir loyihada bevosita siz bilan ishlaydigan odamlar."
    />
    <section class="section section--flush-top">
      <div class="container">
        <TeamGrid />
      </div>
    </section>
  </div>
</template>
