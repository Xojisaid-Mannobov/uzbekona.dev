<script setup lang="ts">
import { onBeforeUnmount, onMounted } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import SiteNavbar from '@/components/layout/SiteNavbar.vue'
import SiteFooter from '@/components/layout/SiteFooter.vue'
import CtaSection from '@/components/layout/CtaSection.vue'
import { useSettingsStore } from '@/stores/settings'
import { startSmoothScroll, stopSmoothScroll, whenIdle } from '@/composables/useSmoothScroll'

const route = useRoute()
const settings = useSettingsStore()
settings.load()

onMounted(() => whenIdle(startSmoothScroll))
onBeforeUnmount(stopSmoothScroll)
</script>

<template>
  <div class="public">
    <a href="#main" class="skip-link">Asosiy kontentga o‘tish</a>
    <SiteNavbar />
    <main id="main" tabindex="-1">
      <RouterView v-slot="{ Component, route: r }">
        <Transition name="page" mode="out-in">
          <component :is="Component" :key="r.path" />
        </Transition>
      </RouterView>
    </main>
    <CtaSection v-if="!route.meta.hideCta" />
    <SiteFooter />
  </div>
</template>

<style scoped>
.public {
  min-height: 100vh;
  /* to'liq kenglikdagi rasmlar gorizontal scroll yaratmasin */
  overflow-x: clip;
  display: flex;
  flex-direction: column;
}

main {
  flex: 1;
  /* sahifa chunk'i yuklanguncha footer ekranga chiqib, keyin surilmasin (CLS) */
  min-height: 100vh;
  outline: none;
}

.page-enter-active,
.page-leave-active {
  transition:
    opacity 380ms var(--ease),
    transform 380ms var(--ease);
}

.page-enter-from {
  opacity: 0;
  transform: translateY(14px);
}

.page-leave-to {
  opacity: 0;
}
</style>
