<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import UiButton from '@/components/ui/UiButton.vue'
import BrandLogo from './BrandLogo.vue'
import { navigation } from '@/content/site'

const route = useRoute()
const scrolled = ref(false)
const menuOpen = ref(false)

function onScroll() {
  scrolled.value = window.scrollY > 12
}

onMounted(() => {
  onScroll()
  window.addEventListener('scroll', onScroll, { passive: true })
})
onBeforeUnmount(() => window.removeEventListener('scroll', onScroll))

// Sahifa almashganda mobil menyu yopiladi
watch(
  () => route.fullPath,
  () => (menuOpen.value = false),
)
watch(menuOpen, (open) => document.documentElement.classList.toggle('menu-open', open))

function onKey(e: KeyboardEvent) {
  if (e.key === 'Escape') menuOpen.value = false
}
</script>

<template>
  <header class="nav" :class="{ 'nav--scrolled': scrolled || menuOpen, 'nav--open': menuOpen }" @keydown="onKey">
    <div class="container nav__inner">
      <RouterLink to="/" class="nav__logo" aria-label="Uzbekona.dev — bosh sahifa">
        <BrandLogo />
      </RouterLink>

      <nav class="nav__links" aria-label="Asosiy navigatsiya">
        <RouterLink v-for="item in navigation" :key="item.to" :to="item.to" class="nav__link">
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="nav__actions">
        <UiButton to="/contact" size="md" icon="arrow-up-right" class="nav__cta">Bog‘lanish</UiButton>
        <button
          class="nav__burger"
          type="button"
          :aria-expanded="menuOpen"
          aria-controls="mobile-menu"
          :aria-label="menuOpen ? 'Menyuni yopish' : 'Menyuni ochish'"
          @click="menuOpen = !menuOpen"
        >
          <AppIcon :name="menuOpen ? 'close' : 'menu'" :size="24" />
        </button>
      </div>
    </div>

    <Transition name="menu">
      <div v-if="menuOpen" id="mobile-menu" class="menu">
        <nav class="container menu__links" aria-label="Mobil navigatsiya">
          <RouterLink v-for="(item, i) in navigation" :key="item.to" :to="item.to" class="menu__link" :style="{ '--i': i }">
            <span class="menu__index">0{{ i + 1 }}</span>
            {{ item.label }}
          </RouterLink>
          <UiButton to="/contact" icon="arrow-up-right" block class="menu__cta">Bog‘lanish</UiButton>
        </nav>
      </div>
    </Transition>
  </header>
</template>

<style scoped>
.nav {
  position: sticky;
  top: 0;
  z-index: 50;
  height: var(--nav-h);
  border-bottom: 1px solid transparent;
  transition:
    background-color var(--dur) var(--ease),
    border-color var(--dur) var(--ease),
    backdrop-filter var(--dur) var(--ease);
}

.nav--scrolled {
  background: color-mix(in srgb, var(--bg) 78%, transparent);
  backdrop-filter: saturate(160%) blur(18px);
  -webkit-backdrop-filter: saturate(160%) blur(18px);
  border-bottom-color: var(--line);
}

.nav__inner {
  height: 100%;
  display: flex;
  align-items: center;
  gap: 40px;
}

.nav__logo {
  display: inline-flex;
  margin-right: auto;
}

.nav__links {
  display: flex;
  align-items: center;
  gap: 4px;
}

.nav__link {
  position: relative;
  padding: 12px 18px;
  font-size: 16px;
  font-weight: 500;
  color: var(--ink-2);
  border-radius: var(--r-pill);
  transition:
    color var(--dur-fast) var(--ease),
    background-color var(--dur-fast) var(--ease);
}

.nav__link:hover {
  color: var(--ink);
  background: var(--surface-2);
}

.nav__link.router-link-active {
  color: var(--ink);
}

.nav__link.router-link-active::after {
  content: '';
  position: absolute;
  left: 50%;
  bottom: 4px;
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: var(--accent);
  translate: -50% 0;
}

.nav__actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.nav__burger {
  display: none;
  width: 52px;
  height: 52px;
  border-radius: 50%;
  place-items: center;
  background: var(--surface-2);
}

/* ─── Mobil menyu ────────────────────────────────── */
.menu {
  position: fixed;
  inset: var(--nav-h) 0 0;
  background: var(--bg);
  overflow-y: auto;
}

.menu__links {
  display: grid;
  padding-top: 24px;
  padding-bottom: 40px;
}

.menu__link {
  display: flex;
  align-items: baseline;
  gap: 16px;
  padding: 20px 0;
  font-size: 36px;
  font-weight: 600;
  letter-spacing: -0.035em;
  border-bottom: 1px solid var(--line);
  animation: menu-in 600ms var(--ease) both;
  animation-delay: calc(var(--i) * 50ms + 80ms);
}

.menu__index {
  font-size: 14px;
  color: var(--ink-3);
  letter-spacing: 0;
}

.menu__cta {
  margin-top: 32px;
}

@keyframes menu-in {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
}

.menu-enter-active,
.menu-leave-active {
  transition: opacity var(--dur) var(--ease);
}

.menu-enter-from,
.menu-leave-to {
  opacity: 0;
}

@media (max-width: 1100px) {
  .nav__links {
    display: none;
  }

  .nav__burger {
    display: grid;
  }
}

@media (max-width: 560px) {
  .nav__cta {
    display: none;
  }
}
</style>

<style>
html.menu-open {
  overflow: hidden;
}
</style>
