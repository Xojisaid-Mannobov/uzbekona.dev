<script setup lang="ts">
import MediaImage from '@/components/media/MediaImage.vue'
import BeforeAfter from './BeforeAfter.vue'
import { vReveal } from '@/composables/reveal'
import { paragraphs, pad2 } from '@/utils/format'
import type { Block } from '@/types/api'

/**
 * Content builder bloklarini chizadi (loyiha case-study va maqolalar uchun).
 * Matnlar faqat text interpolation orqali chiqadi — XSS xavfi yo'q (v-html ishlatilmaydi).
 */
withDefaults(defineProps<{ blocks: Block[]; mode?: 'case' | 'article' }>(), { mode: 'case' })

function youtubeEmbed(url?: string): string | null {
  if (!url) return null
  const m = url.match(/(?:youtube\.com\/(?:watch\?v=|embed\/|shorts\/)|youtu\.be\/)([\w-]{11})/)
  return m ? `https://www.youtube-nocookie.com/embed/${m[1]}` : null
}
</script>

<template>
  <div class="blocks" :class="`blocks--${mode}`">
    <template v-for="(block, i) in blocks" :key="block.id ?? i">
      <!-- Sarlavha -->
      <div v-if="block.type === 'heading'" class="b b-heading narrow" v-reveal>
        <p v-if="block.data.label" class="t-label">{{ block.data.label }}</p>
        <h2 class="t-h2">{{ block.data.text }}</h2>
      </div>

      <!-- Oddiy matn -->
      <div v-else-if="block.type === 'text'" class="b b-text narrow">
        <p v-for="(p, j) in paragraphs(block.data.text)" :key="j">{{ p }}</p>
      </div>

      <!-- Katta matn -->
      <div v-else-if="block.type === 'large_text'" class="b b-large narrow" v-reveal>
        <p v-for="(p, j) in paragraphs(block.data.text)" :key="j" class="t-statement">{{ p }}</p>
      </div>

      <!-- Rasm (container kengligi) va to'liq ekran rasmi -->
      <figure
        v-else-if="(block.type === 'image' || block.type === 'full_image') && block.data.media"
        class="b b-image"
        :class="{ 'b-image--full': block.type === 'full_image' }"
        v-reveal="{ variant: 'scale' }"
      >
        <div class="b-image__frame">
          <MediaImage :media="block.data.media" :sizes="block.type === 'full_image' ? '100vw' : '(min-width: 1440px) 1360px, 100vw'" />
        </div>
        <figcaption v-if="block.data.caption">{{ block.data.caption }}</figcaption>
      </figure>

      <!-- Galereya -->
      <div v-else-if="block.type === 'gallery' && block.data.items?.length" class="b b-gallery">
        <figure v-for="(item, j) in block.data.items" :key="j" class="b-gallery__item" v-reveal="(j % 2) * 100">
          <div class="b-image__frame">
            <MediaImage :media="item.media" sizes="(min-width: 1024px) 680px, 100vw" />
          </div>
          <figcaption v-if="item.caption">{{ item.caption }}</figcaption>
        </figure>
      </div>

      <!-- Video -->
      <figure v-else-if="block.type === 'video'" class="b b-video" v-reveal>
        <div class="b-video__frame">
          <iframe
            v-if="youtubeEmbed(block.data.url)"
            :src="youtubeEmbed(block.data.url)!"
            title="Video"
            loading="lazy"
            allow="accelerometer; encrypted-media; gyroscope; picture-in-picture"
            allowfullscreen
          />
          <video
            v-else-if="block.data.media"
            :src="block.data.media.url"
            :poster="block.data.poster?.url"
            controls
            playsinline
            preload="metadata"
          />
        </div>
        <figcaption v-if="block.data.caption">{{ block.data.caption }}</figcaption>
      </figure>

      <!-- Raqamlar -->
      <dl v-else-if="block.type === 'stats' && block.data.items?.length" class="b b-stats">
        <div v-for="(s, j) in block.data.items" :key="j" class="b-stats__item" v-reveal="j * 80">
          <dt>{{ s.label }}</dt>
          <dd>{{ s.value }}</dd>
        </div>
      </dl>

      <!-- Iqtibos -->
      <blockquote v-else-if="block.type === 'quote'" class="b b-quote narrow" v-reveal>
        <p>“{{ block.data.text }}”</p>
        <footer v-if="block.data.author">
          <strong>{{ block.data.author }}</strong>
          <span v-if="block.data.role"> · {{ block.data.role }}</span>
        </footer>
      </blockquote>

      <!-- 2 va 3 ustun -->
      <div
        v-else-if="block.type === 'two_columns' || block.type === 'three_columns'"
        class="b b-cols"
        :class="block.type === 'two_columns' ? 'b-cols--2' : 'b-cols--3'"
      >
        <div v-for="(c, j) in block.data.columns" :key="j" class="b-cols__item" v-reveal="j * 90">
          <h3 v-if="c.title">{{ c.title }}</h3>
          <p v-for="(p, k) in paragraphs(c.text)" :key="k">{{ p }}</p>
        </div>
      </div>

      <!-- Texnologiyalar -->
      <div v-else-if="block.type === 'technology' && block.data.items?.length" class="b b-tech" v-reveal>
        <p class="t-label">{{ block.data.title || 'Texnologiyalar' }}</p>
        <ul role="list">
          <li v-for="t in block.data.items" :key="t">{{ t }}</li>
        </ul>
      </div>

      <!-- Jarayon -->
      <ol v-else-if="block.type === 'process' && block.data.steps?.length" role="list" class="b b-process">
        <li v-for="(s, j) in block.data.steps" :key="j" v-reveal="(j % 4) * 80">
          <span class="b-process__num">{{ pad2(j + 1) }}</span>
          <h3>{{ s.title }}</h3>
          <p>{{ s.text }}</p>
        </li>
      </ol>

      <!-- Oldin / keyin -->
      <figure v-else-if="block.type === 'before_after' && block.data.before && block.data.after" class="b b-image" v-reveal>
        <BeforeAfter :before="block.data.before" :after="block.data.after" />
        <figcaption v-if="block.data.caption">{{ block.data.caption }}</figcaption>
      </figure>
    </template>
  </div>
</template>

<style scoped>
.blocks {
  display: grid;
  gap: clamp(56px, 7vw, 112px);
}

.blocks--article {
  gap: 32px;
}

.narrow {
  max-width: 920px;
}

.blocks--case .narrow {
  margin-left: calc(100% / 4);
}

.b-heading {
  display: grid;
  gap: 20px;
}

.blocks--article .b-heading {
  margin-top: 24px;
}

.blocks--article .b-heading .t-h2 {
  font-size: clamp(28px, 2.6vw, 40px);
}

.b-text {
  display: grid;
  gap: 1.1em;
  font-size: var(--fs-body-lg);
  line-height: 1.65;
  color: var(--ink-2);
}

.blocks--article .b-text {
  color: var(--ink);
  font-size: clamp(18px, 1.35vw, 20px);
}

.b-large {
  display: grid;
  gap: 0.6em;
  max-width: 1100px;
}

.blocks--article .b-large .t-statement {
  font-size: clamp(26px, 2.4vw, 36px);
}

.b-image figcaption,
.b-gallery figcaption,
.b-video figcaption {
  margin-top: 16px;
  font-size: var(--fs-small);
  color: var(--ink-3);
}

.b-image__frame {
  position: relative;
  border-radius: var(--r-lg);
  overflow: hidden;
  background: var(--surface-2);
}

.b-image--full {
  /* To'liq viewport kengligi */
  width: 100vw;
  margin-left: calc(50% - 50vw);
}

.b-image--full .b-image__frame {
  border-radius: 0;
}

.b-image--full figcaption {
  padding-inline: var(--gutter);
  max-width: calc(var(--container) + var(--gutter) * 2);
  margin-inline: auto;
}

.b-gallery {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: clamp(16px, 2vw, 32px);
}

.b-video__frame {
  position: relative;
  aspect-ratio: 16 / 9;
  border-radius: var(--r-lg);
  overflow: hidden;
  background: var(--dark);
}

.b-video__frame iframe,
.b-video__frame video {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  border: 0;
  object-fit: cover;
}

.b-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 32px;
}

.b-stats__item {
  display: flex;
  flex-direction: column-reverse;
  gap: 12px;
  padding-top: 28px;
  border-top: 1px solid var(--line-strong);
}

.b-stats dd {
  font-size: clamp(56px, 6vw, 96px);
  line-height: 0.9;
  font-weight: 600;
  letter-spacing: -0.055em;
}

.b-stats dt {
  color: var(--ink-2);
}

.b-quote p {
  font-size: clamp(28px, 3vw, 48px);
  line-height: 1.15;
  letter-spacing: -0.03em;
  font-weight: 500;
}

.blocks--article .b-quote {
  padding-left: 28px;
  border-left: 3px solid var(--accent);
}

.blocks--article .b-quote p {
  font-size: clamp(24px, 2.2vw, 32px);
}

.b-quote footer {
  margin-top: 24px;
  color: var(--ink-2);
}

.b-quote strong {
  color: var(--ink);
}

.b-cols {
  display: grid;
  gap: clamp(32px, 4vw, 64px);
}

.b-cols--2 {
  grid-template-columns: 1fr 1fr;
}

.b-cols--3 {
  grid-template-columns: repeat(3, 1fr);
}

.b-cols__item {
  display: grid;
  align-content: start;
  gap: 16px;
  padding-top: 32px;
  border-top: 1px solid var(--line-strong);
  color: var(--ink-2);
  font-size: var(--fs-body-lg);
}

.b-cols__item h3 {
  font-size: clamp(26px, 2.2vw, 34px);
  letter-spacing: -0.03em;
  color: var(--ink);
}

.b-tech {
  display: grid;
  gap: 24px;
}

.b-tech ul {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.b-tech li {
  display: inline-flex;
  align-items: center;
  height: 60px;
  padding: 0 26px;
  border-radius: var(--r-pill);
  background: var(--surface);
  border: 1px solid var(--line);
  font-size: clamp(18px, 1.6vw, 24px);
  font-weight: 600;
  letter-spacing: -0.025em;
}

.b-process {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 0 32px;
}

.b-process li {
  display: grid;
  align-content: start;
  gap: 12px;
  padding: 32px 0;
  border-top: 1px solid var(--line-strong);
}

.b-process__num {
  font-size: clamp(56px, 5vw, 80px);
  line-height: 0.9;
  font-weight: 600;
  letter-spacing: -0.06em;
  color: var(--accent);
  margin-bottom: 16px;
}

.b-process h3 {
  font-size: clamp(24px, 2vw, 30px);
  letter-spacing: -0.03em;
}

.b-process p {
  color: var(--ink-2);
}

@media (max-width: 1024px) {
  .blocks--case .narrow {
    margin-left: 0;
  }

  .b-cols--3 {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .b-cols--2,
  .b-gallery {
    grid-template-columns: 1fr;
  }
}
</style>
