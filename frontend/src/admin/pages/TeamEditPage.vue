<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import FormField from '@/admin/components/FormField.vue'
import MediaField from '@/admin/components/MediaField.vue'
import EditorSidebar from '@/admin/components/EditorSidebar.vue'
import ErrorSummary from '@/admin/components/ErrorSummary.vue'
import { useEditor } from '@/admin/composables/useEditor'
import { adminApi, type TeamMemberInput } from '@/admin/services/adminApi'
import type { MediaRef, TeamMember } from '@/types/api'

type TeamForm = Omit<TeamMemberInput, 'photo_id'> & { photo: MediaRef | null }

const route = useRoute()
const id = computed(() => (route.params.id ? Number(route.params.id) : null))

const { form, loading, loadError, saving, errors, dirty, isNew, save, load } = useEditor<TeamMember, TeamForm, TeamMemberInput>({
  id,
  empty: () => ({ name: '', role: '', bio: '', photo: null, socials: [], is_published: true }),
  fromEntity: (m) => ({
    name: m.name,
    role: m.role,
    bio: m.bio,
    photo: m.photo,
    socials: m.socials.map((s) => ({ ...s })),
    is_published: m.is_published,
  }),
  toInput: ({ photo, ...rest }) => ({ ...rest, photo_id: photo?.id ?? null, socials: rest.socials.filter((s) => s.url) }),
  load: adminApi.team.get,
  create: adminApi.team.create,
  update: adminApi.team.update,
  editRoute: (m) => `/admin/team/${m.id}`,
})

const platforms = ['github', 'telegram', 'linkedin', 'instagram', 'x', 'website']
</script>

<template>
  <div>
    <PageHeader :title="isNew ? 'Yangi a’zo' : form.name || 'Jamoa a’zosi'" back="/admin/team" />

    <AdminEmpty v-if="loadError" error title="Yuklab bo‘lmadi" :text="loadError.message" @retry="load" />
    <div v-else-if="loading" class="skeleton" style="height: 420px; border-radius: 20px" />

    <form v-else class="a-edit" novalidate @submit.prevent="save">
      <div>
        <ErrorSummary :errors="errors" :labels="{ socials: 'Ijtimoiy tarmoq' }" />
        <section class="a-card a-stack">
          <div class="a-form-grid">
            <FormField label="Ism familiya" for="name" :error="errors.name">
              <input id="name" v-model="form.name" class="a-input" placeholder="Xojisaid Mannopov" />
            </FormField>
            <FormField label="Lavozim" for="role" :error="errors.role">
              <input id="role" v-model="form.role" class="a-input" placeholder="Founder · Full-stack Engineer" />
            </FormField>
          </div>
          <FormField label="Qisqa bio" for="bio" :error="errors.bio" optional>
            <textarea id="bio" v-model="form.bio" class="a-textarea" />
          </FormField>

          <div class="a-field">
            <span class="a-label">Ijtimoiy tarmoqlar</span>
            <div v-for="(s, i) in form.socials" :key="i" class="social">
              <select v-model="s.platform" class="a-select" aria-label="Platforma">
                <option v-for="p in platforms" :key="p" :value="p">{{ p }}</option>
              </select>
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
              @click="form.socials.push({ platform: 'github', url: '' })"
            >
              <AppIcon name="plus" :size="16" /> Havola qo‘shish
            </button>
          </div>
        </section>
      </div>

      <aside class="a-edit__side">
        <EditorSidebar :saving="saving" :dirty="dirty" :is-new="isNew" @save="save">
          <label class="a-switch">
            <input v-model="form.is_published" type="checkbox" />
            <span class="a-switch__track" />
            Saytda ko‘rsatish
          </label>
        </EditorSidebar>
        <section class="a-card">
          <h2 class="a-card__title">Rasm (4:5)</h2>
          <MediaField v-model="form.photo" ratio="4 / 5" label="Portret tanlash" />
        </section>
      </aside>
    </form>
  </div>
</template>

<style scoped>
.social {
  display: grid;
  grid-template-columns: 150px 1fr auto;
  gap: 8px;
}

.is-err {
  border-color: var(--danger);
}
</style>
