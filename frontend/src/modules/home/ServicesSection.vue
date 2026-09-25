<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import MediaImage from '@/components/media/MediaImage.vue'
import { useAsync } from '@/composables/useAsync'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'
import { pad2 } from '@/utils/format'
import type { Service } from '@/types/api'

withDefaults(defineProps<{ index?: string; showHeader?: boolean }>(), { index: '02', showHeader: true })

const { data: services, loading, error, reload } = useAsync(publicApi.services)

// Hover paytida kursor ortidan suzuvchi preview (faqat preview rasmi bor xizmatlar uchun)
const hovered = ref<Service | null>(null)
const pos = ref({ x: 0, y: 0 })

function onMove(e: MouseEvent) {
  pos.value = { x: e.clientX, y: e.clientY }
}
</script>

<template>
  <section class="section services" aria-labelledby="services-title">
    <div class="container">
      <SectionHeader
        title-id="services-title"
        v-if="showHeader"
        :index="index"
        label="Xizmatlar"
        title="Biznes uchun raqamli mahsulotlar."
        lead="G‘oyadan production’gacha: tadqiqot, dizayn, ishlab chiqish va qo‘llab-quvvatlash bitta jamoada."
      />

      <div v-if="loading && !services" class="services__list">
        <div v-for="n in 5" :key="n" class="skeleton services__skeleton" />
      </div>

      <StateMessage
        v-else-if="error"
        tone="error"
        title="Xizmatlarni yuklab bo‘lmadi"
        :text="error.message"
        action-label="Qayta urinish"
        @action="reload"
      />

      <ol v-else role="list" class="services__list" @mousemove="onMove" @mouseleave="hovered = null">
        <li v-for="(s, i) in services" :key="s.id" v-reveal="Math.min(i, 4) * 60">
          <RouterLink :to="`/services/${s.slug}`" class="service" @mouseenter="hovered = s">
            <span class="service__num">{{ pad2(i + 1) }}</span>
            <span class="service__title">{{ s.title }}</span>
            <span class="service__summary">{{ s.summary }}</span>
            <span class="service__arrow"><AppIcon name="arrow-up-right" :size="26" /></span>
          </RouterLink>
        </li>
      </ol>
    </div>

    <Transition name="float">
      <div
        v-if="hovered?.preview"
        class="services__float"
        :style="{ transform: `translate3d(${pos.x + 28}px, ${pos.y - 120}px, 0)` }"
        aria-hidden="true"
      >
        <MediaImage :media="hovered.preview" sizes="320px" fill />
      </div>
    </Transition>
  </section>
</template>

<style scoped>
.services__list {
  border-top: 1px solid var(--line-strong);
}

.services__skeleton {
  height: 104px;
  margin-top: 12px;
}

.service {
  display: grid;
  grid-template-columns: 96px minmax(0, 1.2fr) minmax(0, 1fr) 56px;
  align-items: center;
  gap: 24px;
  min-height: 112px;
  padding: 20px 8px;
  border-bottom: 1px solid var(--line-strong);
  transition:
    padding var(--dur) var(--ease),
    background-color var(--dur) var(--ease);
}

.service__num {
  font-size: var(--fs-small);
  font-weight: 700;
  color: var(--ink-3);
  font-variant-numeric: tabular-nums;
}

.service__title {
  font-size: var(--fs-row);
  font-weight: 600;
  letter-spacing: -0.035em;
  line-height: 1.1;
  transition: transform var(--dur) var(--ease);
}

.service__summary {
  font-size: var(--fs-body);
  color: var(--ink-2);
  opacity: 0;
  transform: translateY(8px);
  transition:
    opacity var(--dur) var(--ease),
    transform var(--dur) var(--ease);
}

.service__arrow {
  display: grid;
  place-items: center;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  justify-self: end;
  transition:
    background-color var(--dur) var(--ease),
    color var(--dur) var(--ease),
    transform var(--dur) var(--ease);
}

.service:hover,
.service:focus-visible {
  padding-inline: 24px;
  background: var(--surface);
}

.service:hover .service__title {
  transform: translateX(6px);
}

.service:hover .service__summary,
.service:focus-visible .service__summary {
  opacity: 1;
  transform: none;
}

.service:hover .service__arrow {
  background: var(--accent);
  color: var(--on-accent);
  transform: rotate(45deg);
}

.services__float {
  position: fixed;
  left: 0;
  top: 0;
  z-index: 40;
  width: 320px;
  height: 220px;
  border-radius: var(--r-md);
  overflow: hidden;
  pointer-events: none;
  box-shadow: var(--shadow-lg);
  transition: transform 180ms linear;
}

.float-enter-active,
.float-leave-active {
  transition:
    opacity var(--dur) var(--ease),
    scale var(--dur) var(--ease);
}

.float-enter-from,
.float-leave-to {
  opacity: 0;
  scale: 0.9;
}

/* Sensorli ekranlarda hover yo'q — tavsif doim ko'rinadi */
@media (hover: none), (max-width: 1024px) {
  .service {
    grid-template-columns: 56px minmax(0, 1fr) 48px;
    row-gap: 8px;
  }

  .service__summary {
    grid-column: 2 / 3;
    grid-row: 2;
    opacity: 1;
    transform: none;
    font-size: 16px;
  }

  .service__arrow {
    grid-row: 1 / span 2;
    grid-column: 3;
    width: 48px;
    height: 48px;
  }
}
</style>
