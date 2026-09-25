<script setup lang="ts">
import { reactive, ref } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'
import UiButton from '@/components/ui/UiButton.vue'
import { publicApi } from '@/services/public'
import { toApiError } from '@/services/http'
import { budgets, projectTypes } from '@/content/site'
import type { ContactPayload } from '@/types/api'

// Ta'lim krediti uslubidagi katta, tushunarli forma: 60px inputlar, chip tanlovlar
const initial = (): ContactPayload => ({
  name: '',
  contact: '',
  email: '',
  project_type: '',
  budget: '',
  message: '',
  website: '', // honeypot
})

const form = reactive(initial())
const errors = ref<Record<string, string>>({})
const status = ref<'idle' | 'sending' | 'sent' | 'error'>('idle')
const serverMessage = ref('')

// Brauzer tomonidagi tez tekshiruv (asosiy tekshiruv — backendda)
function validate() {
  const e: Record<string, string> = {}
  if (form.name.trim().length < 2) e.name = 'Ismingizni yozing'
  if (!form.contact.trim() && !form.email.trim()) e.contact = 'Telefon, Telegram yoki email qoldiring'
  if (form.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) e.email = 'Email noto‘g‘ri'
  if (form.message.trim().length < 10) e.message = 'Loyiha haqida kamida bir-ikki gap yozing'
  errors.value = e
  return Object.keys(e).length === 0
}

async function submit() {
  if (status.value === 'sending' || !validate()) return
  status.value = 'sending'
  try {
    const res = await publicApi.contact({ ...form })
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

function reset() {
  status.value = 'idle'
  serverMessage.value = ''
}
</script>

<template>
  <div class="contact-form">
    <Transition name="swap" mode="out-in">
      <div v-if="status === 'sent'" class="sent" role="status">
        <span class="sent__icon"><AppIcon name="check" :size="36" /></span>
        <h2 class="t-h3">Rahmat! So‘rovingiz qabul qilindi.</h2>
        <p class="t-lead">{{ serverMessage }} Odatda bir ish kuni ichida javob beramiz.</p>
        <UiButton variant="secondary" size="md" @click="reset">Yana so‘rov yuborish</UiButton>
      </div>

      <form v-else novalidate @submit.prevent="submit">
        <div class="grid">
          <div class="field" :class="{ 'has-error': errors.name }">
            <label for="cf-name">Ismingiz *</label>
            <input
              id="cf-name"
              v-model="form.name"
              type="text"
              autocomplete="name"
              placeholder="Xojisaid Mannopov"
              :aria-invalid="!!errors.name"
              aria-describedby="cf-name-err"
            />
            <p v-if="errors.name" id="cf-name-err" class="field__error">{{ errors.name }}</p>
          </div>

          <div class="field" :class="{ 'has-error': errors.contact }">
            <label for="cf-contact">Telefon yoki Telegram *</label>
            <input
              id="cf-contact"
              v-model="form.contact"
              type="text"
              autocomplete="tel"
              placeholder="+998 90 123 45 67 yoki @khuja"
              :aria-invalid="!!errors.contact"
              aria-describedby="cf-contact-err"
            />
            <p v-if="errors.contact" id="cf-contact-err" class="field__error">{{ errors.contact }}</p>
          </div>

          <div class="field field--full" :class="{ 'has-error': errors.email }">
            <label for="cf-email">Email</label>
            <input
              id="cf-email"
              v-model="form.email"
              type="email"
              autocomplete="email"
              placeholder="khuja@company.uz"
              :aria-invalid="!!errors.email"
              aria-describedby="cf-email-err"
            />
            <p v-if="errors.email" id="cf-email-err" class="field__error">{{ errors.email }}</p>
          </div>

          <fieldset class="field field--full">
            <legend>Loyiha turi</legend>
            <div class="chips">
              <label v-for="t in projectTypes" :key="t" class="chip" :class="{ 'is-active': form.project_type === t }">
                <input v-model="form.project_type" type="radio" name="project_type" :value="t" class="visually-hidden" />
                {{ t }}
              </label>
            </div>
          </fieldset>

          <fieldset class="field field--full">
            <legend>Taxminiy budjet</legend>
            <div class="chips">
              <label v-for="b in budgets" :key="b" class="chip" :class="{ 'is-active': form.budget === b }">
                <input v-model="form.budget" type="radio" name="budget" :value="b" class="visually-hidden" />
                {{ b }}
              </label>
            </div>
          </fieldset>

          <div class="field field--full" :class="{ 'has-error': errors.message }">
            <label for="cf-message">Loyiha haqida *</label>
            <textarea
              id="cf-message"
              v-model="form.message"
              rows="6"
              placeholder="Qanday muammoni hal qilmoqchisiz? Kimlar foydalanadi? Muddat bormi?"
              :aria-invalid="!!errors.message"
              aria-describedby="cf-message-err"
            />
            <p v-if="errors.message" id="cf-message-err" class="field__error">{{ errors.message }}</p>
          </div>

          <!-- Honeypot: odamlarga ko'rinmaydi, botlar to'ldiradi -->
          <div class="visually-hidden" aria-hidden="true">
            <label for="cf-website">Website</label>
            <input id="cf-website" v-model="form.website" type="text" tabindex="-1" autocomplete="off" />
          </div>
        </div>

        <div class="actions">
          <p v-if="status === 'error'" class="form-error" role="alert"><AppIcon name="alert" :size="18" /> {{ serverMessage }}</p>
          <p v-else class="hint">Ma’lumotlaringiz faqat siz bilan bog‘lanish uchun ishlatiladi.</p>
          <UiButton type="submit" icon="arrow-up-right" :loading="status === 'sending'">So‘rov yuborish</UiButton>
        </div>
      </form>
    </Transition>
  </div>
</template>

<style scoped>
.grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 28px 20px;
}

.field {
  display: grid;
  gap: 10px;
  border: 0;
  padding: 0;
  margin: 0;
  min-width: 0;
}

.field--full {
  grid-column: 1 / -1;
}

label,
legend {
  font-size: 15px;
  font-weight: 600;
  color: var(--ink-2);
  padding: 0;
}

legend {
  margin-bottom: 12px;
}

input,
textarea {
  width: 100%;
  height: var(--input-h);
  padding: 0 22px;
  border-radius: var(--r-sm);
  border: 1px solid var(--line-strong);
  background: var(--surface);
  font-size: 17px;
  transition:
    border-color var(--dur-fast) var(--ease),
    box-shadow var(--dur-fast) var(--ease);
}

textarea {
  height: auto;
  min-height: 200px;
  padding: 18px 22px;
  resize: vertical;
  line-height: 1.55;
}

input::placeholder,
textarea::placeholder {
  color: var(--ink-3);
}

input:hover,
textarea:hover {
  border-color: var(--ink-3);
}

input:focus,
textarea:focus {
  outline: none;
  border-color: var(--accent);
  box-shadow: 0 0 0 4px var(--accent-soft);
}

.has-error input,
.has-error textarea {
  border-color: var(--danger);
}

.field__error {
  font-size: 14px;
  color: var(--danger);
  font-weight: 500;
}

.chips {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.chip {
  display: inline-flex;
  align-items: center;
  height: 52px;
  padding: 0 22px;
  border-radius: var(--r-pill);
  border: 1px solid var(--line-strong);
  background: var(--surface);
  font-size: 16px;
  font-weight: 600;
  color: var(--ink);
  cursor: pointer;
  transition:
    background-color var(--dur-fast) var(--ease),
    color var(--dur-fast) var(--ease),
    border-color var(--dur-fast) var(--ease);
}

.chip:hover {
  border-color: var(--ink);
}

.chip.is-active {
  background: var(--ink);
  border-color: var(--ink);
  color: var(--bg);
}

.chip:focus-within {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

.actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
  margin-top: 40px;
  flex-wrap: wrap;
}

.hint {
  font-size: 14px;
  color: var(--ink-3);
}

.form-error {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--danger);
  font-weight: 600;
}

.sent {
  display: grid;
  justify-items: start;
  gap: 20px;
  padding: clamp(40px, 5vw, 72px);
  border-radius: var(--r-lg);
  background: var(--surface);
  border: 1px solid var(--line);
}

.sent__icon {
  display: grid;
  place-items: center;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: color-mix(in srgb, var(--success) 14%, transparent);
  color: var(--success);
}

.swap-enter-active,
.swap-leave-active {
  transition:
    opacity var(--dur) var(--ease),
    transform var(--dur) var(--ease);
}

.swap-enter-from,
.swap-leave-to {
  opacity: 0;
  transform: translateY(12px);
}

@media (max-width: 640px) {
  .grid {
    grid-template-columns: 1fr;
  }

  .actions > :deep(.btn) {
    width: 100%;
  }
}
</style>
