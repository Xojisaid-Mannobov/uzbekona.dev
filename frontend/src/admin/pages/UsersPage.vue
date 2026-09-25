<script setup lang="ts">
import { reactive, ref } from 'vue'
import AppIcon from '@/components/ui/AppIcon.vue'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import FormField from '@/admin/components/FormField.vue'
import { useAsync } from '@/composables/useAsync'
import { useAuthStore } from '@/admin/stores/auth'
import { useUiStore } from '@/admin/stores/ui'
import { adminApi } from '@/admin/services/adminApi'
import { toApiError } from '@/services/http'
import { formatDateTime } from '@/utils/format'
import type { Admin } from '@/types/api'

const auth = useAuthStore()
const ui = useUiStore()
const { data: admins, loading, error, reload } = useAsync(adminApi.users.list)

const form = reactive({ name: '', email: '', password: '' })
const errors = ref<Record<string, string>>({})
const saving = ref(false)
const showForm = ref(false)

async function create() {
  errors.value = {}
  saving.value = true
  try {
    await adminApi.users.create({ ...form })
    Object.assign(form, { name: '', email: '', password: '' })
    showForm.value = false
    ui.success('Admin qo‘shildi')
    reload()
  } catch (e) {
    const err = toApiError(e)
    errors.value = err.fields
    if (!Object.keys(err.fields).length) ui.error(err.message)
  } finally {
    saving.value = false
  }
}

async function remove(a: Admin) {
  const ok = await ui.confirm({
    title: `${a.name} o‘chirilsinmi?`,
    text: 'Bu admin panelga kira olmaydi. Uning barcha sessiyalari darhol yopiladi.',
  })
  if (!ok) return
  try {
    await adminApi.users.remove(a.id)
    ui.success('O‘chirildi')
    reload()
  } catch (e) {
    ui.error(toApiError(e).message)
  }
}
</script>

<template>
  <div>
    <PageHeader
      title="Adminlar"
      subtitle="Admin panelga kirish huquqiga ega foydalanuvchilar. Public foydalanuvchilardan butunlay alohida."
    >
      <button type="button" class="a-btn a-btn--primary" @click="showForm = !showForm">
        <AppIcon :name="showForm ? 'close' : 'plus'" :size="18" /> {{ showForm ? 'Yopish' : 'Admin qo‘shish' }}
      </button>
    </PageHeader>

    <Transition name="slide">
      <form v-if="showForm" class="a-card new-admin" novalidate @submit.prevent="create">
        <div class="a-form-grid">
          <FormField label="Ism" for="n" :error="errors.name">
            <input id="n" v-model="form.name" class="a-input" placeholder="Khuja Mannopov" />
          </FormField>
          <FormField label="Email" for="e" :error="errors.email">
            <input id="e" v-model="form.email" class="a-input" type="email" autocomplete="off" />
          </FormField>
          <FormField
            label="Vaqtinchalik parol"
            for="p"
            :error="errors.password"
            hint="Kamida 10 ta belgi. Admin keyin Sozlamalardan o‘zgartiradi"
            span
          >
            <input id="p" v-model="form.password" class="a-input" type="text" autocomplete="new-password" />
          </FormField>
        </div>
        <button type="submit" class="a-btn a-btn--primary" style="margin-top: 20px" :disabled="saving">
          {{ saving ? 'Qo‘shilmoqda…' : 'Qo‘shish' }}
        </button>
      </form>
    </Transition>

    <AdminEmpty v-if="error" error title="Yuklab bo‘lmadi" :text="error.message" @retry="reload" />
    <div v-else-if="loading && !admins" class="skeleton" style="height: 200px; border-radius: 20px" />
    <table v-else class="a-table">
      <thead>
        <tr>
          <th>Admin</th>
          <th>Oxirgi kirish</th>
          <th>Qo‘shilgan</th>
          <th />
        </tr>
      </thead>
      <tbody>
        <tr v-for="a in admins" :key="a.id">
          <td>
            <div class="a-row">
              <span class="avatar">{{ a.name.charAt(0) }}</span>
              <span>
                <span class="a-cell-title">{{ a.name }} <small v-if="a.id === auth.admin?.id" class="you">siz</small></span>
                <span class="a-cell-sub" style="display: block">{{ a.email }}</span>
              </span>
            </div>
          </td>
          <td class="a-cell-sub">{{ formatDateTime(a.last_login_at) }}</td>
          <td class="a-cell-sub">{{ formatDateTime(a.created_at) }}</td>
          <td>
            <div class="a-actions">
              <button
                v-if="a.id !== auth.admin?.id"
                type="button"
                class="a-icon-btn a-icon-btn--danger"
                title="O‘chirish"
                @click="remove(a)"
              >
                <AppIcon name="trash" :size="18" />
              </button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.new-admin {
  margin-bottom: 20px;
}

.avatar {
  display: grid;
  place-items: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--accent-soft);
  color: var(--accent);
  font-weight: 700;
}

.you {
  margin-left: 6px;
  padding: 2px 8px;
  border-radius: 99px;
  background: var(--surface-2);
  font-size: 11px;
  color: var(--ink-2);
}

.slide-enter-active,
.slide-leave-active {
  transition:
    opacity 220ms var(--ease),
    transform 220ms var(--ease);
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
