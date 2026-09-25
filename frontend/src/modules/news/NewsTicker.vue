<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import OrnamentStar from '@/components/ornament/OrnamentStar.vue'
import { useAsync } from '@/composables/useAsync'
import { publicApi } from '@/services/public'
import { toneStyle } from './tone'

// "Yangiliklar lentasi" — hero ostida yuguruvchi sarlavhalar (televizion yangiliklar uslubida).
// Lenta ikki nusxada chiziladi: birinchisi o'qiladi, ikkinchisi faqat uzluksiz aylanish uchun (aria-hidden).
const { data } = useAsync(() => publicApi.news({ limit: 6 }))
const items = computed(() => data.value?.items ?? [])
const duration = computed(() => `${Math.max(24, items.value.length * 9)}s`)
</script>

<template>
  <aside v-if="items.length" class="ticker" aria-label="So‘nggi yangiliklar">
    <RouterLink to="/news" class="ticker__label"><span class="ticker__live" aria-hidden="true" /> Yangiliklar</RouterLink>
    <div class="ticker__viewport">
      <div class="ticker__track" :style="{ '--dur': duration }">
        <ul v-for="copy in 2" :key="copy" class="ticker__list" role="list" :aria-hidden="copy === 2 ? 'true' : undefined">
          <li v-for="n in items" :key="n.id" class="ticker__item" :style="toneStyle(n.tag)">
            <RouterLink :to="`/news/${n.slug}`" :tabindex="copy === 2 ? -1 : undefined">
              <span v-if="n.tag" class="ticker__tag">{{ n.tag }}</span>
              {{ n.title }}
            </RouterLink>
            <OrnamentStar :size="12" class="ticker__sep" />
          </li>
        </ul>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.ticker {
  display: flex;
  align-items: stretch;
  height: 56px;
  border-block: 1px solid var(--line);
  background: var(--surface);
  overflow: hidden;
}

.ticker__label {
  position: relative;
  z-index: 1;
  display: inline-flex;
  align-items: center;
  gap: 10px;
  flex-shrink: 0;
  padding: 0 22px 0 max(var(--gutter), (100vw - var(--container)) / 2);
  background: var(--ink);
  color: var(--bg);
  font-size: var(--fs-xs);
  font-weight: 700;
  letter-spacing: 0.16em;
  text-transform: uppercase;
}

/* Yorliq o'ng chetida qiya kesim */
.ticker__label::after {
  content: '';
  position: absolute;
  top: 0;
  right: -18px;
  width: 36px;
  height: 100%;
  background: inherit;
  transform: skewX(-18deg);
  z-index: -1;
}

.ticker__live {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #e5262c;
  box-shadow: 0 0 0 0 rgb(229 38 44 / 0.6);
  animation: live 1.8s ease-out infinite;
}

@keyframes live {
  to {
    box-shadow: 0 0 0 9px rgb(229 38 44 / 0);
  }
}

.ticker__viewport {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  mask-image: linear-gradient(90deg, transparent 0, #000 48px, #000 calc(100% - 48px), transparent 100%);
}

.ticker__track {
  display: flex;
  width: max-content;
  height: 100%;
  padding-left: 36px;
  animation: ticker var(--dur) linear infinite;
}

.ticker:hover .ticker__track,
.ticker:focus-within .ticker__track {
  animation-play-state: paused;
}

@keyframes ticker {
  to {
    transform: translateX(-50%);
  }
}

.ticker__list {
  display: flex;
  align-items: center;
}

.ticker__item {
  display: inline-flex;
  align-items: center;
  gap: 28px;
  padding-right: 28px;
  white-space: nowrap;
  font-weight: 600;
  font-size: 15px;
}

.ticker__item a {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  transition: color var(--dur-fast) var(--ease);
}

.ticker__item a:hover {
  color: var(--accent);
}

.ticker__tag {
  padding: 3px 9px;
  border-radius: var(--r-pill);
  background: color-mix(in srgb, var(--tone) 14%, transparent);
  color: var(--tone);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
}

.ticker__sep {
  color: var(--gold);
}

@media (max-width: 560px) {
  .ticker {
    height: 50px;
  }

  .ticker__label {
    padding-right: 14px;
    letter-spacing: 0.1em;
  }
}

@media (prefers-reduced-motion: reduce) {
  .ticker__viewport {
    overflow-x: auto;
  }

  .ticker__track {
    animation: none;
  }
}
</style>
