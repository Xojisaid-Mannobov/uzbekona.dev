<script setup lang="ts">
import { RouterLink } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import TeamCard from './TeamCard.vue'
import { useAsync } from '@/composables/useAsync'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'

// 3 ustunli jamoa grid'i + oxirida "jamoaga qo'shiling" kartasi
const { data: team, loading, error, reload } = useAsync(publicApi.team)
</script>

<template>
  <div v-if="loading && !team" class="team-grid">
    <div v-for="n in 3" :key="n" class="skeleton team-grid__skeleton" />
  </div>

  <StateMessage
    v-else-if="error"
    tone="error"
    title="Jamoani yuklab bo‘lmadi"
    :text="error.message"
    action-label="Qayta urinish"
    @action="reload"
  />

  <div v-else class="team-grid">
    <div v-for="(m, i) in team" :key="m.id" v-reveal="(i % 3) * 90">
      <TeamCard :member="m" />
    </div>
    <RouterLink to="/contact" class="team-grid__join" v-reveal="((team?.length ?? 0) % 3) * 90">
      <span class="team-grid__plus"><AppIcon name="plus" :size="32" /></span>
      <span>
        <strong>Jamoaga qo‘shiling</strong>
        <span class="team-grid__join-text">Kuchli engineer va dizaynerlar bilan ishlashni istaymiz. O‘zingiz haqingizda yozing.</span>
      </span>
    </RouterLink>
  </div>
</template>

<style scoped>
.team-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 64px clamp(20px, 2.4vw, 40px);
}

.team-grid__skeleton {
  aspect-ratio: 4 / 5;
  border-radius: var(--r-lg);
}

.team-grid__join {
  aspect-ratio: 4 / 5;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: clamp(28px, 3vw, 44px);
  border-radius: var(--r-lg);
  border: 1px dashed var(--line-strong);
  transition:
    background-color var(--dur) var(--ease),
    border-color var(--dur) var(--ease);
}

.team-grid__join:hover {
  background: var(--surface);
  border-color: var(--ink);
}

.team-grid__plus {
  display: grid;
  place-items: center;
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: var(--ink);
  color: var(--bg);
  transition: transform var(--dur) var(--ease);
}

.team-grid__join:hover .team-grid__plus {
  transform: rotate(90deg);
}

.team-grid__join strong {
  display: block;
  font-size: clamp(24px, 1.9vw, 28px);
  letter-spacing: -0.03em;
  margin-bottom: 12px;
}

.team-grid__join-text {
  color: var(--ink-2);
}

@media (max-width: 1024px) {
  .team-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 600px) {
  .team-grid {
    grid-template-columns: 1fr;
  }

  .team-grid__join {
    aspect-ratio: auto;
    min-height: 300px;
  }
}
</style>
