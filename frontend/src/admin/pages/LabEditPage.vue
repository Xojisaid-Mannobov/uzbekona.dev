<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import FormField from '@/admin/components/FormField.vue'
import TagInput from '@/admin/components/TagInput.vue'
import MediaField from '@/admin/components/MediaField.vue'
import EditorSidebar from '@/admin/components/EditorSidebar.vue'
import ErrorSummary from '@/admin/components/ErrorSummary.vue'
import { useEditor } from '@/admin/composables/useEditor'
import { adminApi, type LabInput } from '@/admin/services/adminApi'
import type { Lab, MediaRef } from '@/types/api'

type LabForm = Omit<LabInput, 'cover_id'> & { cover: MediaRef | null }

const route = useRoute()
const id = computed(() => (route.params.id ? Number(route.params.id) : null))

const { form, loading, loadError, saving, errors, dirty, isNew, save, load } = useEditor<Lab, LabForm, LabInput>({
  id,
  empty: () => ({
    title: '',
    slug: '',
    description: '',
    stage: 'in_development',
    url: '',
    repo_url: '',
    cover: null,
    stack: [],
    is_published: true,
  }),
  fromEntity: (l) => ({
    title: l.title,
    slug: l.slug,
    description: l.description,
    stage: l.stage,
    url: l.url,
    repo_url: l.repo_url,
    cover: l.cover,
    stack: [...l.stack],
    is_published: l.is_published,
  }),
  toInput: ({ cover, ...rest }) => ({ ...rest, cover_id: cover?.id ?? null }),
  load: adminApi.labs.get,
  create: adminApi.labs.create,
  update: adminApi.labs.update,
  editRoute: (l) => `/admin/labs/${l.id}`,
})
</script>

<template>
  <div>
    <PageHeader :title="isNew ? 'Yangi Labs loyihasi' : form.title || 'Labs'" back="/admin/labs" />

    <AdminEmpty v-if="loadError" error title="Yuklab bo‘lmadi" :text="loadError.message" @retry="load" />
    <div v-else-if="loading" class="skeleton" style="height: 420px; border-radius: 20px" />

    <form v-else class="a-edit" novalidate @submit.prevent="save">
      <div>
        <ErrorSummary :errors="errors" />
        <section class="a-card a-stack">
          <div class="a-form-grid">
            <FormField label="Nomi" for="title" :error="errors.title">
              <input id="title" v-model="form.title" class="a-input" placeholder="Telekit" />
            </FormField>
            <FormField label="Slug" for="slug" :error="errors.slug" optional>
              <input id="slug" v-model="form.slug" class="a-input" placeholder="telekit" />
            </FormField>
          </div>
          <FormField label="Tavsif" for="desc" :error="errors.description">
            <textarea id="desc" v-model="form.description" class="a-textarea" />
          </FormField>
          <FormField label="Texnologiyalar" :error="errors.stack">
            <TagInput v-model="form.stack" :suggestions="['Go', 'TypeScript', 'Python', 'Docker']" />
          </FormField>
          <div class="a-form-grid">
            <FormField label="Kod havolasi (GitHub)" for="repo" :error="errors.repo_url" optional>
              <input id="repo" v-model="form.repo_url" class="a-input" placeholder="https://github.com/…" />
            </FormField>
            <FormField label="Sayt / demo havolasi" for="url" :error="errors.url" optional>
              <input id="url" v-model="form.url" class="a-input" placeholder="https://" />
            </FormField>
          </div>
        </section>
      </div>

      <aside class="a-edit__side">
        <EditorSidebar :saving="saving" :dirty="dirty" :is-new="isNew" @save="save">
          <FormField label="Holat" for="stage" :error="errors.stage">
            <select id="stage" v-model="form.stage" class="a-select">
              <option value="open_source">Open Source</option>
              <option value="experimental">Experimental</option>
              <option value="in_development">In Development</option>
            </select>
          </FormField>
          <label class="a-switch">
            <input v-model="form.is_published" type="checkbox" />
            <span class="a-switch__track" />
            Saytda ko‘rsatish
          </label>
        </EditorSidebar>
        <section class="a-card">
          <h2 class="a-card__title">Cover <small class="a-hint">ixtiyoriy</small></h2>
          <MediaField v-model="form.cover" />
        </section>
      </aside>
    </form>
  </div>
</template>
