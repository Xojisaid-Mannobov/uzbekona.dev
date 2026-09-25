<script setup lang="ts">
import OrnamentStar from './OrnamentStar.vue'

// Footer'dagi katta "UZBEKONA" uslubida — sekin suzuvchi ulkan yozuvlar (ikki qator, qarama-qarshi yo'nalishda).
// Birinchi qator kulrang "arvoh", ikkinchisi oltin tusli. Faqat bezak: ekran o'quvchilardan yashirilgan,
// so'zlar ::before orqali chiziladi — DOM matni emas, shuning uchun kontrast tekshiruvlariga tushmaydi.
defineProps<{ rows: string[][] }>()
</script>

<template>
  <div class="ghost" aria-hidden="true">
    <div v-for="(row, r) in rows" :key="r" class="ghost__row" :class="{ 'ghost__row--outline': r % 2 === 1 }">
      <div class="ghost__track">
        <span v-for="copy in 2" :key="copy" class="ghost__copy">
          <template v-for="word in row" :key="word">
            <span class="ghost__word" :data-word="word" />
            <OrnamentStar class="ghost__star" :size="64" />
          </template>
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.ghost {
  display: grid;
  gap: clamp(0px, 0.5vw, 8px);
  padding-block: clamp(24px, 4vw, 64px);
  overflow: hidden;
  user-select: none;
  pointer-events: none;
}

.ghost__row {
  --dur: 70s;
  overflow: hidden;
  contain: layout paint;
}

.ghost__track {
  display: flex;
  width: max-content;
  will-change: transform;
  animation: ghost-move var(--dur) linear infinite;
}

.ghost__row--outline .ghost__track {
  animation-direction: reverse;
  --dur: 85s;
}

@keyframes ghost-move {
  to {
    transform: translateX(-50%);
  }
}

.ghost__copy {
  display: flex;
  align-items: center;
}

.ghost__word::before {
  content: attr(data-word);
}

.ghost__word {
  padding-inline: clamp(18px, 2.4vw, 40px);
  font-size: clamp(72px, 11vw, 176px);
  font-weight: 700;
  line-height: 1;
  letter-spacing: -0.05em;
  white-space: nowrap;
  color: color-mix(in srgb, var(--ink) 7%, var(--bg));
}

/* Ikkinchi qator — oltin tusli (variable shriftda kontur ichki chiziqlarni ko'rsatib qo'yadi, shuning uchun to'la rang) */
.ghost__row--outline .ghost__word {
  color: color-mix(in srgb, var(--gold) 16%, var(--bg));
}

.ghost__star {
  width: clamp(32px, 4.6vw, 72px);
  height: auto;
  --star-color: color-mix(in srgb, var(--gold) 45%, var(--bg));
}

@media (prefers-reduced-motion: reduce) {
  .ghost__track {
    animation: none;
  }
}
</style>
