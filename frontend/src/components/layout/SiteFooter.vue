<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import OrnamentBorder from '@/components/ornament/OrnamentBorder.vue'
import Doppi from '@/components/ornament/Doppi.vue'
import BrandIcon from '@/components/ui/BrandIcon.vue'
import { useSettingsStore } from '@/stores/settings'
import { useTheme } from '@/composables/useTheme'
import { footerNavigation } from '@/content/site'

const settings = useSettingsStore()
const theme = useTheme()
const year = new Date().getFullYear()

const socials = computed(() => settings.socials.filter((s) => s.url))

// "UZBEKONA" dagi O harfiga do'ppi: harf o'rni shrift yuklangandan keyin o'lchanadi
// (getExtentOfChar — gorizontal, canvas measureText — harfning yuqori chegarasi)
const BASELINE = 148
const wordText = ref<SVGTextElement>()
const doppi = ref<{ x: number; y: number; w: number; h: number } | null>(null)

async function placeDoppi() {
  await document.fonts?.ready
  const t = wordText.value
  const ctx = document.createElement('canvas').getContext('2d')
  if (!t || !ctx) return
  try {
    const ext = t.getExtentOfChar(5)
    const cs = getComputedStyle(t)
    ctx.font = `${cs.fontWeight} ${cs.fontSize} ${cs.fontFamily}`
    const top = BASELINE - ctx.measureText('O').actualBoundingBoxAscent
    const w = ext.width * 1.04
    const h = (w * 80) / 120
    // qiya burilish chapga siljitadi — markazni biroz o'ngga surib to'g'rilaymiz
    doppi.value = { x: ext.x + (ext.width - w) / 2 + w * 0.06, y: top - h * 0.74, w, h }
  } catch {
    // Brauzer SVG matn o'lchamini bermasa — do'ppisiz qoladi
  }
}
onMounted(placeDoppi)
</script>

<template>
  <footer class="footer">
    <!-- Hoshiya — an'anaviy chegara naqshi CTA va footer orasida -->
    <OrnamentBorder class="footer__hoshiya" :opacity="0.45" />
    <div class="container">
      <div class="footer__grid">
        <div class="footer__brand">
          <p class="footer__name">Uzbekona<span>.dev</span></p>
          <p class="footer__tagline">Digital Product Studio. Web platformalar, mobil ilovalar, Telegram tizimlari va avtomatlashtirish.</p>
          <a v-if="settings.site.email" :href="`mailto:${settings.site.email}`" class="footer__email">
            {{ settings.site.email }}
            <AppIcon name="arrow-up-right" :size="22" />
          </a>
        </div>

        <nav class="footer__col" aria-label="Footer navigatsiya">
          <p class="footer__heading">Navigation</p>
          <RouterLink v-for="item in footerNavigation" :key="item.to" :to="item.to" class="footer__link">{{ item.label }}</RouterLink>
        </nav>

        <div class="footer__col">
          <p class="footer__heading">Social</p>
          <a v-for="s in socials" :key="s.url" :href="s.url" class="footer__link footer__tg" target="_blank" rel="noopener noreferrer">
            <BrandIcon v-if="s.platform.toLowerCase() === 'telegram'" name="telegram" :size="20" />
            {{ s.label || s.platform }}
          </a>
        </div>

        <div class="footer__col">
          <p class="footer__heading">Kontakt</p>
          <a
            v-if="settings.telegramUrl"
            :href="settings.telegramUrl"
            class="footer__link footer__tg"
            target="_blank"
            rel="noopener noreferrer"
          >
            <BrandIcon name="telegram" :size="20" /> {{ settings.site.telegram }}
          </a>
          <a v-if="settings.site.phone" :href="`tel:${settings.site.phone.replace(/\s/g, '')}`" class="footer__link">{{
            settings.site.phone
          }}</a>
          <p v-if="settings.site.address" class="footer__muted">{{ settings.site.address }}</p>
        </div>
      </div>

      <div class="footer__bottom">
        <p>© {{ year }} Uzbekona.dev</p>
        <button type="button" class="footer__theme" @click="theme.toggle()">
          <AppIcon :name="theme.isDark() ? 'sun' : 'moon'" :size="18" />
          {{ theme.isDark() ? 'Yorug‘ rejim' : 'Qorong‘i rejim' }}
        </button>
      </div>
    </div>

    <!-- Dekorativ wordmark: SVG — konteyner kengligiga aniq moslashadi -->
    <svg class="footer__wordmark" viewBox="0 -100 1000 250" preserveAspectRatio="xMidYMax meet" aria-hidden="true" focusable="false">
      <text ref="wordText" x="500" :y="BASELINE" text-anchor="middle" textLength="980" lengthAdjust="spacingAndGlyphs">UZBEKONA</text>
      <!-- O harfiga kiydirilgan do'ppi -->
      <g v-if="doppi" class="footer__doppi" :transform="`rotate(-9 ${doppi.x + doppi.w / 2} ${doppi.y + doppi.h})`">
        <Doppi :x="doppi.x" :y="doppi.y" :width="doppi.w" :height="doppi.h" />
      </g>
    </svg>
  </footer>
</template>

<style scoped>
.footer {
  background: var(--dark);
  color: var(--on-dark);
  padding-top: 0;
  overflow: hidden;
}

.footer__hoshiya {
  margin-bottom: clamp(56px, 6vw, 96px);
}

.footer__grid {
  display: grid;
  grid-template-columns: 1.6fr 1fr 1fr 1fr;
  gap: 48px;
}

.footer__name {
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.04em;
}

.footer__name span {
  color: var(--accent);
}

.footer__tagline {
  margin-top: 20px;
  max-width: 360px;
  color: var(--on-dark-2);
}

.footer__email {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  margin-top: 40px;
  font-size: clamp(19px, 1.6vw, 24px);
  font-weight: 600;
  letter-spacing: -0.03em;
  border-bottom: 1px solid var(--dark-line);
  padding-bottom: 6px;
  transition: border-color var(--dur-fast) var(--ease);
}

.footer__email:hover {
  border-color: var(--on-dark);
}

.footer__col {
  display: grid;
  align-content: start;
  gap: 14px;
}

.footer__heading {
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: var(--on-dark-2);
  margin-bottom: 8px;
}

.footer__link {
  font-size: 16px;
  font-weight: 500;
  width: fit-content;
  transition: opacity var(--dur-fast) var(--ease);
}

.footer__link:hover {
  opacity: 0.6;
}

.footer__muted {
  color: var(--on-dark-2);
}

.footer__bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
  margin-top: clamp(64px, 7vw, 112px);
  padding-block: 28px;
  border-top: 1px solid var(--dark-line);
  color: var(--on-dark-2);
  font-size: var(--fs-small);
}

.footer__theme {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 44px;
  padding: 0 18px;
  border-radius: var(--r-pill);
  border: 1px solid var(--dark-line);
  color: var(--on-dark);
  font-size: 14px;
  font-weight: 600;
  transition: border-color var(--dur-fast) var(--ease);
}

.footer__theme:hover {
  border-color: var(--on-dark);
}

.footer__wordmark {
  display: block;
  width: 100%;
  height: auto;
  /* yuqoridagi bo'sh joy do'ppi uchun — pastki qatorga yopishmasligi uchun manfiy margin */
  margin-top: -9vw;
  margin-bottom: -2.6vw;
  pointer-events: none;
  fill: var(--dark-2);
  font-family: var(--font);
  font-size: 196px;
  font-weight: 700;
  letter-spacing: -0.05em;
  user-select: none;
}

.footer__tg {
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.footer__doppi {
  --doppi-body: #070808;
  --doppi-outline: rgb(255 255 255 / 0.16);
  animation: doppi-drop 900ms cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

@keyframes doppi-drop {
  from {
    opacity: 0;
    translate: 0 -40px;
  }
}

@media (max-width: 1024px) {
  .footer__grid {
    grid-template-columns: 1fr 1fr;
  }

  .footer__brand {
    grid-column: 1 / -1;
  }
}

@media (max-width: 560px) {
  .footer__grid {
    grid-template-columns: 1fr;
    gap: 40px;
  }
}
</style>
