<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import FormField from '@/admin/components/FormField.vue'
import MetricsEditor from '@/admin/components/MetricsEditor.vue'
import SeoFields from '@/admin/components/SeoFields.vue'
import ErrorSummary from '@/admin/components/ErrorSummary.vue'
import { useUiStore } from '@/admin/stores/ui'
import { useSettingsStore } from '@/stores/settings'
import { adminApi, type SettingsInput } from '@/admin/services/adminApi'
import { toApiError, type ApiError } from '@/services/http'

const ui = useUiStore()
const publicSettings = useSettingsStore()

const form = ref<SettingsInput | null>(null)
const snapshot = ref('')
const loading = ref(true)
const loadError = ref<ApiError | null>(null)
const saving = ref(false)
const errors = ref<Record<string, string>>({})
const dirty = computed(() => !!form.value && JSON.stringify(form.value) !== snapshot.value)

async function load() {
  loading.value = true
  loadError.value = null
  try {
    const s = await adminApi.settings.get()
    form.value = {
      site: { ...s.site },
      metrics: s.metrics.map((m) => ({ ...m })),
      seo: { ...s.seo },
      socials: s.socials.map((x) => ({ ...x })),
    }
    snapshot.value = JSON.stringify(form.value)
  } catch (e) {
    loadError.value = toApiError(e)
  } finally {
    loading.value = false
  }
}

async function save() {
  if (!form.value) return
  saving.value = true
  errors.value = {}
  try {
    const payload = {
      ...form.value,
      metrics: form.value.metrics.filter((m) => m.value || m.label),
      socials: form.value.socials.filter((s) => s.url),
    }
    const s = await adminApi.settings.update(payload)
    publicSettings.apply(s)
    form.value = {
      site: { ...s.site },
      metrics: s.metrics.map((m) => ({ ...m })),
      seo: { ...s.seo },
      socials: s.socials.map((x) => ({ ...x })),
    }
    snapshot.value = JSON.stringify(form.value)
    ui.success('Sozlamalar saqlandi')
  } catch (e) {
    const err = toApiError(e)
    errors.value = err.fields
    ui.error(err.message)
  } finally {
    saving.value = false
  }
}

onMounted(load)

// ─── Parolni o'zgartirish ───────────────────────────
const pwd = reactive({ current: '', next: '', confirm: '' })
const pwdErrors = ref<Record<string, string>>({})
const pwdSaving = ref(false)

async function changePassword() {
  pwdErrors.value = {}
  if (pwd.next.length < 10) pwdErrors.value.new_password = 'Kamida 10 ta belgi'
  if (pwd.next !== pwd.confirm) pwdErrors.value.confirm = 'Parollar mos emas'
  if (Object.keys(pwdErrors.value).length) return
  pwdSaving.value = true
  try {
    await adminApi.auth.changePassword(pwd.current, pwd.next)
    Object.assign(pwd, { current: '', next: '', confirm: '' })
    ui.success('Parol yangilandi. Boshqa qurilmalardagi sessiyalar yopildi.')
  } catch (e) {
    const err = toApiError(e)
    pwdErrors.value = err.fields
    if (!Object.keys(err.fields).length) ui.error(err.message)
  } finally {
    pwdSaving.value = false
  }
}

const platforms = ['github', 'telegram', 'instagram', 'linkedin', 'x', 'youtube', 'behance', 'dribbble']
</script>

<template>
  <div>
    <PageHeader title="Sozlamalar" subtitle="Kontaktlar, metrikalar, SEO va ijtimoiy tarmoqlar.">
      <button v-if="form" type="button" class="a-btn a-btn--primary" :disabled="saving || !dirty" @click="save">
        <AppIcon name="check" :size="18" /> {{ saving ? 'Saqlanmoqda…' : 'Saqlash' }}
      </button>
    </PageHeader>

    <AdminEmpty v-if="loadError" error title="Yuklab bo‘lmadi" :text="loadError.message" @retry="load" />
    <div v-else-if="loading" class="skeleton" style="height: 480px; border-radius: 20px" />

    <div v-else-if="form" class="grid">
      <div class="a-stack">
        <ErrorSummary :errors="errors" :labels="{ site: 'Sayt', metrics: 'Metrika', socials: 'Ijtimoiy tarmoq', seo: 'SEO' }" />

        <section class="a-card">
          <h2 class="a-card__title">Kompaniya va kontaktlar</h2>
          <div class="a-form-grid">
            <FormField label="Nomi" for="name" :error="errors['site.name']">
              <input id="name" v-model="form.site.name" class="a-input" />
            </FormField>
            <FormField label="Tagline" for="tagline" :error="errors['site.tagline']">
              <input id="tagline" v-model="form.site.tagline" class="a-input" />
            </FormField>
            <FormField label="Email" for="email" :error="errors['site.email']">
              <input id="email" v-model="form.site.email" class="a-input" type="email" />
            </FormField>
            <FormField label="Telefon" for="phone" :error="errors['site.phone']" optional>
              <input id="phone" v-model="form.site.phone" class="a-input" placeholder="+998 …" />
            </FormField>
            <FormField label="Telegram" for="tg" :error="errors['site.telegram']" hint="@username yoki havola">
              <input id="tg" v-model="form.site.telegram" class="a-input" placeholder="@uzbekona_dev" />
            </FormField>
            <FormField label="Manzil" for="address" :error="errors['site.address']" optional>
              <input id="address" v-model="form.site.address" class="a-input" />
            </FormField>
          </div>
          <label class="a-switch" style="margin-top: 20px">
            <input v-model="form.site.available" type="checkbox" />
            <span class="a-switch__track" />
            “Yangi loyihalar uchun ochiqmiz” belgisini ko‘rsatish
          </label>
        </section>

        <section class="a-card">
          <h2 class="a-card__title">Metrikalar (bosh sahifa)</h2>
          <MetricsEditor v-model="form.metrics" :errors="errors" />
          <p class="a-hint" style="margin-top: 12px">Faqat real ma’lumotlar. Bo‘sh qoldirilsa, bo‘lim saytda ko‘rinmaydi.</p>
        </section>

        <section class="a-card">
          <h2 class="a-card__title">Ijtimoiy tarmoqlar</h2>
          <div class="a-stack" style="gap: 8px">
            <div v-for="(s, i) in form.socials" :key="i" class="social">
              <select v-model="s.platform" class="a-select" aria-label="Platforma">
                <option v-for="p in platforms" :key="p" :value="p">{{ p }}</option>
              </select>
              <input v-model="s.label" class="a-input" placeholder="Ko‘rinadigan nom" aria-label="Nomi" />
              <input
                v-model="s.url"
                class="a-input"
                placeholder="https://"
                aria-label="Havola"
                :class="{ 'is-err': errors[`socials.${i}.url`] }"
              />
              <button type="button" class="a-icon-btn a-icon-btn--danger" aria-label="O‘chirish" @click="form.socials.splice(i, 1)">
                <AppIcon name="trash" :size="18" />
              </button>
            </div>
            <button
              type="button"
              class="a-btn a-btn--outline a-btn--sm"
              style="justify-self: start"
              @click="form.socials.push({ platform: 'github', label: '', url: '' })"
            >
              <AppIcon name="plus" :size="16" /> Qo‘shish
            </button>
          </div>
        </section>
      </div>

      <div class="a-stack">
        <section class="a-card">
          <h2 class="a-card__title">Standart SEO</h2>
          <SeoFields v-model="form.seo" fallback-title="Uzbekona.dev — Digital Product Studio" path="/" />
        </section>

        <section class="a-card">
          <h2 class="a-card__title">Parolni o‘zgartirish</h2>
          <form class="a-stack" @submit.prevent="changePassword">
            <FormField label="Joriy parol" for="cur" :error="pwdErrors.current_password">
              <input id="cur" v-model="pwd.current" class="a-input" type="password" autocomplete="current-password" />
            </FormField>
            <FormField label="Yangi parol" for="new" :error="pwdErrors.new_password" hint="Kamida 10 ta belgi">
              <input id="new" v-model="pwd.next" class="a-input" type="password" autocomplete="new-password" />
            </FormField>
            <FormField label="Yangi parolni takrorlang" for="conf" :error="pwdErrors.confirm">
              <input id="conf" v-model="pwd.confirm" class="a-input" type="password" autocomplete="new-password" />
            </FormField>
            <button type="submit" class="a-btn a-btn--outline" :disabled="pwdSaving || !pwd.current || !pwd.next">
              {{ pwdSaving ? 'Saqlanmoqda…' : 'Parolni yangilash' }}
            </button>
          </form>
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.grid {
  display: grid;
  grid-template-columns: minmax(0, 1.4fr) minmax(0, 1fr);
  gap: 20px;
  align-items: start;
}

.grid .a-card + .a-card {
  margin-top: 0;
}

.social {
  display: grid;
  grid-template-columns: 140px 1fr 1.4fr auto;
  gap: 8px;
}

.is-err {
  border-color: var(--danger);
}

@media (max-width: 1100px) {
  .grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .social {
    grid-template-columns: 1fr 1fr;
  }
}
</style>
