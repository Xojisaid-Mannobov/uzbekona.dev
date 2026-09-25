<script setup lang="ts">
import PageHero from '@/components/layout/PageHero.vue'
import AppIcon from '@/components/ui/AppIcon.vue'
import type { IconName } from '@/components/ui/icons'
import GhostWord from '@/components/ornament/GhostWord.vue'
import OrnamentStar from '@/components/ornament/OrnamentStar.vue'
import JoinForm from '@/modules/careers/JoinForm.vue'
import { useSeo } from '@/composables/useSeo'
import { vReveal } from '@/composables/reveal'

// "Jamoaga qo'shilish": nomzod o'zi haqida yozadi, jamoa arizani admin panelda ko'rib chiqib bog'lanadi
const perks: { icon: IconName; title: string; text: string }[] = [
  { icon: 'layers', title: 'Real loyihalar', text: 'Minglab odamlar har kuni ishlatadigan tizimlar — o‘quv mashqlari emas.' },
  { icon: 'users', title: 'Kuchli jamoa', text: 'Kod review, juftlikda ishlash va har bir qarorni birga muhokama qilish.' },
  { icon: 'chart', title: 'O‘sish', text: 'Stajyordan mustaqil muhandisgacha — aniq yo‘l va murabbiy.' },
  { icon: 'star', title: 'Milliy ruh', text: 'Raqamli O‘zbekistonni quradigan mahsulotlarga o‘z hissangizni qo‘shasiz.' },
]

const steps = [
  { title: 'Ariza', text: 'Formani to‘ldirasiz — o‘zingiz haqingizda va ishlaringiz haqida.' },
  { title: 'Ko‘rib chiqamiz', text: '3 ish kuni ichida arizani o‘qib, javob beramiz.' },
  { title: 'Suhbat', text: 'Tanishuv va kichik amaliy topshiriq.' },
  { title: 'Taklif', text: 'Birga ishlashni boshlaymiz.' },
]

useSeo({
  title: 'Jamoaga qo‘shilish',
  description: 'Uzbekona.dev jamoasiga qo‘shiling: o‘zingiz haqingizda yozing — arizani ko‘rib chiqib, siz bilan bog‘lanamiz.',
})
</script>

<template>
  <div>
    <PageHero
      label="Jamoaga qo‘shilish"
      :title="'Biz bilan birga\nquring.'"
      lead="Kuchli dasturchi, dizayner yoki menejer bo‘lsangiz — o‘zingiz haqingizda yozing. Arizani diqqat bilan o‘qib chiqamiz va albatta javob beramiz."
    />

    <section class="section section--flush-top">
      <div class="container join">
        <div class="join__form" v-reveal>
          <JoinForm />
        </div>

        <aside class="join__aside">
          <div class="join__block">
            <p class="t-label join__label"><OrnamentStar :size="14" /> Nega biz bilan</p>
            <ul role="list" class="perks">
              <li v-for="(p, i) in perks" :key="p.title" v-reveal="i * 70">
                <span class="perks__icon"><AppIcon :name="p.icon" :size="20" /></span>
                <div>
                  <strong>{{ p.title }}</strong>
                  <p>{{ p.text }}</p>
                </div>
              </li>
            </ul>
          </div>

          <div class="join__block">
            <p class="t-label join__label"><OrnamentStar :size="14" /> Keyin nima bo‘ladi</p>
            <ol role="list" class="steps">
              <li v-for="(s, i) in steps" :key="s.title">
                <span class="steps__num">0{{ i + 1 }}</span>
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

    <GhostWord text="JAMOA" class="join__ghost" />
  </div>
</template>

<style scoped>
.join {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(0, 1fr);
  gap: clamp(48px, 6vw, 112px);
  align-items: start;
}

.join__aside {
  display: grid;
  gap: 48px;
  position: sticky;
  top: calc(var(--nav-h) + 32px);
}

.join__block {
  display: grid;
  gap: 20px;
}

.join__label {
  display: flex;
  align-items: center;
  gap: 10px;
}

.perks {
  display: grid;
  gap: 12px;
}

.perks li {
  display: flex;
  gap: 16px;
  padding: 18px;
  border-radius: var(--r-md);
  background: var(--surface);
  border: 1px solid var(--line);
}

.perks__icon {
  display: grid;
  place-items: center;
  flex-shrink: 0;
  width: 42px;
  height: 42px;
  border-radius: 12px;
  background: var(--accent-soft);
  color: var(--accent);
}

.perks strong,
.steps strong {
  display: block;
  font-size: 17px;
  letter-spacing: -0.02em;
  margin-bottom: 4px;
}

.perks p,
.steps p {
  color: var(--ink-2);
  font-size: 15px;
  line-height: 1.5;
}

.steps {
  display: grid;
  gap: 18px;
}

.steps li {
  display: flex;
  gap: 18px;
  padding-top: 18px;
  border-top: 1px solid var(--line);
}

.steps__num {
  font-size: 24px;
  font-weight: 600;
  letter-spacing: -0.05em;
  color: var(--accent);
  line-height: 1;
}

.join__ghost {
  margin-top: clamp(24px, 4vw, 64px);
}

@media (max-width: 1024px) {
  .join {
    grid-template-columns: 1fr;
  }

  .join__aside {
    position: static;
  }
}
</style>
