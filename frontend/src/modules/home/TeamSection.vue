<script setup lang="ts">
import SectionHeader from '@/components/ui/SectionHeader.vue'
import UiButton from '@/components/ui/UiButton.vue'
import TeamGrid from '@/modules/team/TeamGrid.vue'
</script>

<template>
  <section class="section" aria-labelledby="team-title">
    <div class="container">
      <SectionHeader title-id="team-title" index="06" label="Jamoa" title="Mahsulotni quradigan odamlar.">
        <template #aside>
          <UiButton to="/team" variant="secondary" icon="arrow-up-right">Jamoa haqida</UiButton>
        </template>
      </SectionHeader>
      <TeamGrid />
    </div>
  </section>
</template>
