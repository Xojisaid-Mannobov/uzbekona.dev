<script setup lang="ts">
import AppIcon from '@/components/ui/AppIcon.vue'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import StateMessage from '@/components/states/StateMessage.vue'
import MediaImage from '@/components/media/MediaImage.vue'
import { useAsync } from '@/composables/useAsync'
import { vReveal } from '@/composables/reveal'
import { publicApi } from '@/services/public'
import { labStageLabels } from '@/content/site'

const { data: labs, loading, error, reload } = useAsync(publicApi.labs)
</script>

<template>
  <section v-if="loading || error || labs?.length" class="section labs" aria-labelledby="labs-title">
    <div class="container">
      <SectionHeader
        title-id="labs-title"
        index="07"
        label="Uzbekona Labs"
        title="Ichki va open-source loyihalar."
        lead="Mijoz loyihalaridan tashqari o‘zimiz uchun va hamjamiyat uchun vositalar yaratamiz."
      />

      <div v-if="loading && !labs" class="labs__grid">
        <div v-for="n in 2" :key="n" class="skeleton labs__skeleton" />
      </div>

      <StateMessage
        v-else-if="error"
        tone="error"
        title="Labs’ni yuklab bo‘lmadi"
        :text="error.message"
        action-label="Qayta urinish"
        @action="reload"
      />

      <div v-else class="labs__grid">
        <article v-for="(lab, i) in labs" :key="lab.id" class="lab" v-reveal="(i % 2) * 110">
          <div v-if="lab.cover" class="lab__cover">
            <MediaImage :media="lab.cover" sizes="(min-width: 1024px) 660px, 100vw" fill />
          </div>
          <div class="lab__head">
            <span class="lab__stage" :class="`lab__stage--${lab.stage}`">{{ labStageLabels[lab.stage] }}</span>
            <span class="lab__num">{{ String(i + 1).padStart(2, '0') }}</span>
          </div>
          <h3 class="lab__title">{{ lab.title }}</h3>
          <p class="lab__desc">{{ lab.description }}</p>
          <div class="lab__foot">
            <ul role="list" class="lab__stack">
              <li v-for="t in lab.stack" :key="t">{{ t }}</li>
            </ul>
            <div class="lab__links">
              <a v-if="lab.repo_url" :href="lab.repo_url" target="_blank" rel="noopener noreferrer" :aria-label="`${lab.title} — kod`">
                Kod <AppIcon name="arrow-up-right" :size="16" />
              </a>
              <a v-if="lab.url" :href="lab.url" target="_blank" rel="noopener noreferrer" :aria-label="`${lab.title} — sayt`">
                Ochish <AppIcon name="arrow-up-right" :size="16" />
              </a>
            </div>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>

<style scoped>
.labs__grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: clamp(16px, 1.6vw, 24px);
}

.labs__skeleton {
  height: 320px;
  border-radius: var(--r-lg);
}

.lab {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: 320px;
  padding: clamp(28px, 3vw, 48px);
  border-radius: var(--r-lg);
  background: var(--surface);
  border: 1px solid var(--line);
  transition:
    transform var(--dur) var(--ease),
    box-shadow var(--dur) var(--ease);
}

.lab:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.lab__cover {
  position: relative;
  aspect-ratio: 16 / 9;
  border-radius: var(--r-sm);
  overflow: hidden;
  margin-bottom: 8px;
}

.lab__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: auto;
  padding-bottom: 40px;
}

.lab__stage {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 36px;
  padding: 0 16px;
  border-radius: var(--r-pill);
  font-size: 14px;
  font-weight: 600;
  background: var(--surface-2);
}

.lab__stage::before {
  content: '';
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: currentColor;
}

.lab__stage--open_source {
  color: var(--success);
  background: color-mix(in srgb, var(--success) 10%, transparent);
}

.lab__stage--experimental {
  color: var(--accent);
  background: var(--accent-soft);
}

.lab__stage--in_development {
  color: var(--warning);
  background: color-mix(in srgb, var(--warning) 10%, transparent);
}

.lab__num {
  font-size: var(--fs-small);
  font-weight: 700;
  color: var(--ink-3);
}

.lab__title {
  font-size: clamp(24px, 2vw, 32px);
  letter-spacing: -0.04em;
}

.lab__desc {
  max-width: 44ch;
  font-size: var(--fs-body-lg);
  color: var(--ink-2);
}

.lab__foot {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  gap: 24px;
  margin-top: 16px;
  flex-wrap: wrap;
}

.lab__stack {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.lab__stack li {
  padding: 6px 14px;
  border-radius: var(--r-pill);
  border: 1px solid var(--line);
  font-size: 14px;
  font-weight: 500;
  color: var(--ink-2);
}

.lab__links {
  display: flex;
  gap: 20px;
  font-weight: 600;
}

.lab__links a {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.lab__links a:hover {
  color: var(--accent);
}

@media (max-width: 768px) {
  .labs__grid {
    grid-template-columns: 1fr;
  }

  .lab {
    min-height: 280px;
  }
}
</style>
