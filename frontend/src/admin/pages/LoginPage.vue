<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import BrandLogo from '@/components/layout/BrandLogo.vue'
import FormField from '@/admin/components/FormField.vue'
import { useAuthStore } from '@/admin/stores/auth'
import { toApiError } from '@/services/http'
import '@/admin/styles/admin.css'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const form = reactive({ email: '', password: '' })
const errors = ref<Record<string, string>>({})
const message = ref('')
const loading = ref(false)
const showPassword = ref(false)

async function submit() {
  errors.value = {}
  message.value = ''
  if (!form.email) errors.value.email = 'Email kiriting'
  if (!form.password) errors.value.password = 'Parol kiriting'
  if (Object.keys(errors.value).length) return

  loading.value = true
  try {
    await auth.login(form.email, form.password)
    const redirect =
      typeof route.query.redirect === 'string' && route.query.redirect.startsWith('/admin') ? route.query.redirect : '/admin/dashboard'
    router.replace(redirect)
  } catch (e) {
    const err = toApiError(e)
    errors.value = err.fields
    message.value = err.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="admin login">
    <div class="login__card">
      <BrandLogo size="lg" />
      <div class="login__head">
        <h1>Admin panelga kirish</h1>
        <p>Kontent, loyihalar va so‘rovlarni boshqarish.</p>
      </div>

      <form class="a-stack" novalidate @submit.prevent="submit">
        <FormField label="Email" for="email" :error="errors.email">
          <input
            id="email"
            v-model="form.email"
            class="a-input"
            type="email"
            autocomplete="username"
            placeholder="admin@uzbekona.dev"
            autofocus
          />
        </FormField>
        <FormField label="Parol" for="password" :error="errors.password">
          <div class="login__pwd">
            <input
              id="password"
              v-model="form.password"
              class="a-input"
              :type="showPassword ? 'text' : 'password'"
              autocomplete="current-password"
              placeholder="••••••••••"
            />
            <button
              type="button"
              class="a-icon-btn"
              :aria-label="showPassword ? 'Parolni yashirish' : 'Parolni ko‘rsatish'"
              @click="showPassword = !showPassword"
            >
              <AppIcon :name="showPassword ? 'eye-off' : 'eye'" :size="18" />
            </button>
          </div>
        </FormField>

        <p v-if="message" class="login__error" role="alert"><AppIcon name="alert" :size="18" /> {{ message }}</p>

        <button type="submit" class="a-btn a-btn--primary login__submit" :disabled="loading">
          {{ loading ? 'Tekshirilmoqda…' : 'Kirish' }}
          <AppIcon v-if="!loading" name="arrow-right" :size="18" />
        </button>
      </form>

      <a href="/" class="login__back"><AppIcon name="arrow-left" :size="16" /> Saytga qaytish</a>
    </div>
  </div>
</template>

<style scoped>
.login {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
  background: radial-gradient(60% 50% at 80% 10%, var(--accent-soft), transparent 70%), var(--bg);
}

.login__card {
  width: min(440px, 100%);
  display: grid;
  gap: 28px;
  padding: 40px;
  border-radius: 28px;
  background: var(--surface);
  border: 1px solid var(--line);
  box-shadow: var(--shadow-lg);
}

.login__head h1 {
  font-size: 28px;
  letter-spacing: -0.035em;
}

.login__head p {
  margin-top: 6px;
  color: var(--ink-2);
}

.login__pwd {
  position: relative;
}

.login__pwd .a-input {
  padding-right: 52px;
}

.login__pwd .a-icon-btn {
  position: absolute;
  right: 6px;
  top: 50%;
  translate: 0 -50%;
}

.login__error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 14px;
  border-radius: 12px;
  background: color-mix(in srgb, var(--danger) 10%, transparent);
  color: var(--danger);
  font-weight: 600;
  font-size: 14px;
}

.login__submit {
  height: 52px;
  width: 100%;
  font-size: 15px;
}

.login__back {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  justify-self: center;
  font-size: 14px;
  font-weight: 600;
  color: var(--ink-3);
}

.login__back:hover {
  color: var(--ink);
}
</style>
