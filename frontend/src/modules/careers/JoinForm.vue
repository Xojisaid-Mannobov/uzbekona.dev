<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'
import UiButton from '@/components/ui/UiButton.vue'
import { publicApi } from '@/services/public'
import { toApiError } from '@/services/http'
import { careerExperience, careerPositions } from '@/content/site'
import type { ApplicationPayload } from '@/types/api'

// Nomzod o'zi haqida yozadi → ariza admin paneldagi "Nomzodlar" bo'limiga tushadi, jamoa ko'rib chiqib bog'lanadi
const ABOUT_MIN = 50

const initial = (): ApplicationPayload => ({
  full_name: '',
  email: '',
  phone: '',
  telegram: '',
  position: '',
  experience: '',
  portfolio_url: '',
  resume_url: '',
  about: '',
  consent: false,
  website: '', // honeypot
})

const form = reactive(initial())
const errors = ref<Record<string, string>>({})
const status = ref<'idle' | 'sending' | 'sent' | 'error'>('idle')
const serverMessage = ref('')
const aboutLength = computed(() => form.about.trim().length)

const isUrl = (v: string) => !v || /^https?:\/\/\S+\.\S+/.test(v.trim())

// Brauzer tomonidagi tez tekshiruv (asosiy tekshiruv — backendda)
function validate() {
  const e: Record<string, string> = {}
  if (form.full_name.trim().length < 3) e.full_name = 'Ism va familiyangizni yozing'
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email.trim())) e.email = 'Email noto‘g‘ri'
  if (!form.position) e.position = 'Yo‘nalishni tanlang'
  if (!isUrl(form.portfolio_url)) e.portfolio_url = 'Havola https:// bilan boshlansin'
  if (!isUrl(form.resume_url)) e.resume_url = 'Havola https:// bilan boshlansin'
  if (aboutLength.value < ABOUT_MIN) e.about = `O‘zingiz haqingizda kamida ${ABOUT_MIN} ta belgi yozing`
  if (!form.consent) e.consent = 'Rozilik belgisini qo‘ying'
  errors.value = e
  return Object.keys(e).length === 0
}

async function submit() {
  if (status.value === 'sending' || !validate()) return
  status.value = 'sending'
  try {
    const res = await publicApi.apply({ ...form })
    serverMessage.value = res.message
    status.value = 'sent'
    Object.assign(form, initial())
  } catch (e) {
    const err = toApiError(e)
    errors.value = err.fields
    serverMessage.value = err.message
    status.value = 'error'
  }
}
</script>

<template>
  <div class="join-form">
    <Transition name="swap" mode="out-in">
      <div v-if="status === 'sent'" class="sent" role="status">
        <span class="sent__icon"><AppIcon name="check" :size="36" /></span>
        <h2 class="t-h3">Rahmat! Arizangiz bizda.</h2>
        <p class="t-lead">{{ serverMessage }} Odatda 3 ish kuni ichida javob beramiz.</p>
        <UiButton to="/team" variant="secondary" size="md" icon-left="arrow-left">Jamoa bilan tanishish</UiButton>
      </div>

      <form v-else novalidate @submit.prevent="submit">
        <div class="grid">
          <div class="field" :class="{ 'has-error': errors.full_name }">
            <label for="jf-name">Ism va familiya *</label>
            <input
              id="jf-name"
              v-model="form.full_name"
              type="text"
              autocomplete="name"
              placeholder="Aziz Karimov"
              :aria-invalid="!!errors.full_name"
              aria-describedby="jf-name-err"
            />
            <p v-if="errors.full_name" id="jf-name-err" class="field__error">{{ errors.full_name }}</p>
          </div>

          <div class="field" :class="{ 'has-error': errors.email }">
            <label for="jf-email">Email *</label>
            <input
              id="jf-email"
              v-model="form.email"
              type="email"
              autocomplete="email"
              placeholder="aziz@mail.uz"
              :aria-invalid="!!errors.email"
              aria-describedby="jf-email-err"
            />
            <p v-if="errors.email" id="jf-email-err" class="field__error">{{ errors.email }}</p>
          </div>

          <div class="field">
            <label for="jf-phone">Telefon</label>
            <input id="jf-phone" v-model="form.phone" type="tel" autocomplete="tel" placeholder="+998 90 123 45 67" />
          </div>

          <div class="field">
            <label for="jf-tg">Telegram</label>
            <input id="jf-tg" v-model="form.telegram" type="text" placeholder="@username" />
          </div>

          <fieldset class="field field--full" :class="{ 'has-error': errors.position }" aria-describedby="jf-pos-err">
            <legend>Yo‘nalish *</legend>
            <div class="chips">
              <label v-for="p in careerPositions" :key="p" class="chip" :class="{ 'is-active': form.position === p }">
                <input v-model="form.position" type="radio" name="position" :value="p" class="visually-hidden" />
                {{ p }}
              </label>
            </div>
            <p v-if="errors.position" id="jf-pos-err" class="field__error">{{ errors.position }}</p>
          </fieldset>

          <fieldset class="field field--full">
            <legend>Tajriba</legend>
            <div class="chips">
              <label v-for="x in careerExperience" :key="x" class="chip" :class="{ 'is-active': form.experience === x }">
                <input v-model="form.experience" type="radio" name="experience" :value="x" class="visually-hidden" />
                {{ x }}
              </label>
            </div>
          </fieldset>

          <div class="field" :class="{ 'has-error': errors.portfolio_url }">
            <label for="jf-portfolio">Portfolio yoki GitHub</label>
            <input
              id="jf-portfolio"
              v-model="form.portfolio_url"
              type="url"
              inputmode="url"
              placeholder="https://github.com/…"
              :aria-invalid="!!errors.portfolio_url"
              aria-describedby="jf-portfolio-err"
            />
            <p v-if="errors.portfolio_url" id="jf-portfolio-err" class="field__error">{{ errors.portfolio_url }}</p>
          </div>

          <div class="field" :class="{ 'has-error': errors.resume_url }">
            <label for="jf-resume">Rezyume havolasi</label>
            <input
              id="jf-resume"
              v-model="form.resume_url"
              type="url"
              inputmode="url"
              placeholder="Google Drive, LinkedIn, hh.uz…"
              :aria-invalid="!!errors.resume_url"
              aria-describedby="jf-resume-err"
            />
            <p v-if="errors.resume_url" id="jf-resume-err" class="field__error">{{ errors.resume_url }}</p>
          </div>

          <div class="field field--full" :class="{ 'has-error': errors.about }">
            <label for="jf-about">O‘zingiz haqingizda *</label>
            <textarea
              id="jf-about"
              v-model="form.about"
              rows="7"
              placeholder="Nimalar qilgansiz, qaysi texnologiyalarni yaxshi bilasiz, qaysi loyihangiz bilan faxrlanasiz va nega aynan biz bilan ishlamoqchisiz?"
              :aria-invalid="!!errors.about"
              aria-describedby="jf-about-err jf-about-count"
            />
            <p id="jf-about-count" class="counter" :class="{ 'is-low': aboutLength > 0 && aboutLength < ABOUT_MIN }">
              {{ aboutLength }} / {{ ABOUT_MIN }}+ belgi
            </p>
            <p v-if="errors.about" id="jf-about-err" class="field__error">{{ errors.about }}</p>
          </div>

          <div class="field field--full" :class="{ 'has-error': errors.consent }">
            <label class="consent">
              <input v-model="form.consent" type="checkbox" :aria-invalid="!!errors.consent" aria-describedby="jf-consent-err" />
              Ma’lumotlarim arizani ko‘rib chiqish va men bilan bog‘lanish uchun ishlatilishiga roziman *
            </label>
            <p v-if="errors.consent" id="jf-consent-err" class="field__error">{{ errors.consent }}</p>
          </div>

          <!-- Honeypot: odamlarga ko'rinmaydi, botlar to'ldiradi -->
          <div class="visually-hidden" aria-hidden="true">
            <label for="jf-website">Website</label>
            <input id="jf-website" v-model="form.website" type="text" tabindex="-1" autocomplete="off" />
          </div>
        </div>

        <div class="actions">
          <p v-if="status === 'error'" class="form-error" role="alert"><AppIcon name="alert" :size="18" /> {{ serverMessage }}</p>
          <p v-else class="hint">Ariza faqat jamoamizga ko‘rinadi va uchinchi shaxslarga berilmaydi.</p>
          <UiButton type="submit" icon="arrow-up-right" :loading="status === 'sending'">Ariza yuborish</UiButton>
        </div>
      </form>
    </Transition>
  </div>
</template>

<style scoped src="../contact/form.css"></style>
