<script setup lang="ts">
import AppIcon from '@/components/ui/AppIcon.vue'
import BrandIcon from '@/components/ui/BrandIcon.vue'
import PageHero from '@/components/layout/PageHero.vue'
import ContactForm from '@/modules/contact/ContactForm.vue'
import { useSettingsStore } from '@/stores/settings'
import { useSeo } from '@/composables/useSeo'

const settings = useSettingsStore()

const steps = [
  { title: 'So‘rov', text: 'Formani to‘ldirasiz yoki Telegram’da yozasiz.' },
  { title: 'Suhbat', text: '30 daqiqalik qo‘ng‘iroqda vazifani aniqlaymiz.' },
  { title: 'Taklif', text: 'Yechim, muddat va narx bilan aniq taklif yuboramiz.' },
]

useSeo({ title: 'Bog‘lanish', description: 'Loyihangiz haqida yozing — Uzbekona.dev jamoasi bir ish kuni ichida javob beradi.' })
</script>

<template>
  <div>
    <PageHero
      label="Bog‘lanish"
      title="Loyihani boshlaymiz."
      lead="Vazifangizni qisqacha yozing. Birinchi suhbat bepul va hech narsaga majburlamaydi."
    />

    <section class="section section--flush-top">
      <div class="container contact">
        <ContactForm />

        <aside class="contact__aside">
          <div class="contact__block">
            <p class="t-label">To‘g‘ridan-to‘g‘ri</p>
            <a v-if="settings.site.email" :href="`mailto:${settings.site.email}`" class="contact__big">
              {{ settings.site.email }} <AppIcon name="arrow-up-right" :size="22" />
            </a>
            <a v-if="settings.telegramUrl" :href="settings.telegramUrl" class="contact__big" target="_blank" rel="noopener noreferrer">
              <BrandIcon name="telegram" :size="26" /> {{ settings.site.telegram }} <AppIcon name="arrow-up-right" :size="22" />
            </a>
            <a v-if="settings.site.phone" :href="`tel:${settings.site.phone.replace(/\s/g, '')}`" class="contact__big">
              {{ settings.site.phone }}
            </a>
            <p v-if="settings.site.address" class="t-muted">{{ settings.site.address }}</p>
          </div>

          <div class="contact__block">
            <p class="t-label">Keyin nima bo‘ladi</p>
            <ol role="list" class="contact__steps">
              <li v-for="(s, i) in steps" :key="s.title">
                <span class="contact__num">0{{ i + 1 }}</span>
                <div>
                  <strong>{{ s.title }}</strong>
                  <p>{{ s.text }}</p>
                </div>
              </li>
            </ol>
          </div>
        </aside>
      </div>
    </section>
  </div>
</template>

<style scoped>
.contact {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(0, 1fr);
  gap: clamp(48px, 6vw, 112px);
  align-items: start;
}

.contact__aside {
  display: grid;
  gap: 56px;
  position: sticky;
  top: calc(var(--nav-h) + 32px);
}

.contact__block {
  display: grid;
  gap: 16px;
  justify-items: start;
}

.contact__big {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: clamp(19px, 1.6vw, 24px);
  font-weight: 600;
  letter-spacing: -0.03em;
  border-bottom: 1px solid var(--line-strong);
  padding-bottom: 4px;
  transition:
    color var(--dur-fast) var(--ease),
    border-color var(--dur-fast) var(--ease);
}

.contact__big:hover {
  color: var(--accent);
  border-color: var(--accent);
}

.contact__steps {
  display: grid;
  gap: 24px;
  width: 100%;
}

.contact__steps li {
  display: flex;
  gap: 20px;
  padding-top: 20px;
  border-top: 1px solid var(--line);
}

.contact__num {
  font-size: 26px;
  font-weight: 600;
  letter-spacing: -0.05em;
  color: var(--accent);
  line-height: 1;
}

.contact__steps strong {
  font-size: 18px;
  letter-spacing: -0.02em;
}

.contact__steps p {
  color: var(--ink-2);
}

@media (max-width: 1024px) {
  .contact {
    grid-template-columns: 1fr;
  }

  .contact__aside {
    position: static;
  }
}
</style>
