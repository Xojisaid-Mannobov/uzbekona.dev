<script setup lang="ts">
import UiButton from '@/components/ui/UiButton.vue'
import { useSeo } from '@/composables/useSeo'

useSeo({ title: 'Sahifa topilmadi', noindex: true })
</script>

<template>
  <section class="section nf">
    <div class="container">
      <p class="t-label">Xato 404</p>
      <h1 class="t-hero nf__title">Bu sahifa<br />topilmadi.</h1>
      <p class="t-lead nf__text">Havola eskirgan yoki manzilda xato bor. Bosh sahifadan davom eting yoki loyihalarimizni ko‘ring.</p>
      <div class="nf__actions">
        <UiButton to="/" icon-left="arrow-left">Bosh sahifa</UiButton>
        <UiButton to="/projects" variant="secondary" icon="arrow-up-right">Loyihalar</UiButton>
      </div>
    </div>
  </section>
</template>

<style scoped>
.nf {
  min-height: calc(100vh - var(--nav-h));
  display: flex;
  align-items: center;
}

.nf__title {
  margin-top: 32px;
}

.nf__text {
  margin-top: 32px;
  max-width: 560px;
}

.nf__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-top: 48px;
}
</style>
