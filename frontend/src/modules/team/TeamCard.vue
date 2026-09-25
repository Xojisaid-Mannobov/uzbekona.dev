<script setup lang="ts">
import { computed } from 'vue'
import MediaImage from '@/components/media/MediaImage.vue'
import GirihPattern from '@/components/ornament/GirihPattern.vue'
import type { TeamMember } from '@/types/api'

const props = defineProps<{ member: TeamMember }>()

const initials = computed(() =>
  props.member.name
    .split(/\s+/)
    .slice(0, 2)
    .map((p) => p.charAt(0))
    .join(''),
)
</script>

<template>
  <article class="member">
    <div class="member__photo">
      <GirihPattern v-if="!member.photo" :size="64" :opacity="0.16" fade="center" />
      <MediaImage v-if="member.photo" :media="member.photo" sizes="(min-width: 1024px) 420px, 100vw" fill :alt="member.name" />
      <span v-else class="member__initials" aria-hidden="true">{{ initials }}</span>
    </div>
    <h3 class="member__name">{{ member.name }}</h3>
    <p class="member__role">{{ member.role }}</p>
    <p v-if="member.bio" class="member__bio">{{ member.bio }}</p>
    <ul v-if="member.socials.length" role="list" class="member__socials">
      <li v-for="s in member.socials" :key="s.url">
        <a :href="s.url" target="_blank" rel="noopener noreferrer">{{ s.platform }}</a>
      </li>
    </ul>
  </article>
</template>

<style scoped>
.member__photo {
  position: relative;
  aspect-ratio: 4 / 5;
  border-radius: var(--r-lg);
  overflow: hidden;
  background: radial-gradient(90% 70% at 70% 10%, var(--accent-soft), transparent 70%), var(--surface-2);
  display: grid;
  place-items: center;
  margin-bottom: 24px;
}

.member__photo :deep(img) {
  filter: grayscale(0.15);
  transition:
    transform 1s var(--ease),
    filter 1s var(--ease);
}

.member:hover .member__photo :deep(img) {
  transform: scale(1.03);
  filter: none;
}

.member__initials {
  font-size: clamp(56px, 6vw, 96px);
  font-weight: 600;
  letter-spacing: -0.06em;
  color: var(--ink-3);
}

.member__name {
  font-size: clamp(19px, 1.5vw, 23px);
  letter-spacing: -0.03em;
}

.member__role {
  margin-top: 4px;
  font-size: 15px;
  color: var(--ink-2);
}

.member__bio {
  margin-top: 16px;
  color: var(--ink-2);
  max-width: 36ch;
}

.member__socials {
  display: flex;
  gap: 16px;
  margin-top: 16px;
  font-size: var(--fs-small);
  font-weight: 600;
  text-transform: capitalize;
}

.member__socials a:hover {
  color: var(--accent);
}
</style>
