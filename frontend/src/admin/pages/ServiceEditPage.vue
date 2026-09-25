<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import PageHeader from '@/admin/components/PageHeader.vue'
import AdminEmpty from '@/admin/components/AdminEmpty.vue'
import FormField from '@/admin/components/FormField.vue'
import TagInput from '@/admin/components/TagInput.vue'
import MediaField from '@/admin/components/MediaField.vue'
import SeoFields from '@/admin/components/SeoFields.vue'
import EditorSidebar from '@/admin/components/EditorSidebar.vue'
import ErrorSummary from '@/admin/components/ErrorSummary.vue'
import { useEditor } from '@/admin/composables/useEditor'
import { adminApi, type ServiceInput } from '@/admin/services/adminApi'
import type { MediaRef, Service } from '@/types/api'

type ServiceForm = Omit<ServiceInput, 'preview_id'> & { preview: MediaRef | null }

const route = useRoute()
const id = computed(() => (route.params.id ? Number(route.params.id) : null))

const { form, entity, loading, loadError, saving, errors, dirty, isNew, save, load } = useEditor<Service, ServiceForm, ServiceInput>({
  id,
  empty: () => ({
    title: '',
    slug: '',
    summary: '',
    description: '',
    features: [],
    stack: [],
    preview: null,
    status: 'published',
    seo: { title: '', description: '' },
  }),
  fromEntity: (s) => ({
    title: s.title,
    slug: s.slug,
    summary: s.summary,
    description: s.description,
    features: [...s.features],
    stack: [...s.stack],
    preview: s.preview,
    status: s.status,
    seo: { ...s.seo },
  }),
  toInput: ({ preview, ...rest }) => ({ ...rest, preview_id: preview?.id ?? null }),
  load: adminApi.services.get,
  create: adminApi.services.create,
  update: adminApi.services.update,
  editRoute: (s) => `/admin/services/${s.id}`,
})

const publicUrl = computed(() => (entity.value?.status === 'published' ? `/services/${entity.value.slug}` : null))
</script>

<template>
  <div>
    <PageHeader :title="isNew ? 'Yangi xizmat' : form.title || 'Xizmat'" back="/admin/services" />

    <AdminEmpty v-if="loadError" error title="Yuklab bo‘lmadi" :text="loadError.message" @retry="load" />
    <div v-else-if="loading" class="skeleton" style="height: 480px; border-radius: 20px" />

    <form v-else class="a-edit" novalidate @submit.prevent="save">
      <div>
        <ErrorSummary :errors="errors" />
        <section class="a-card a-stack">
          <FormField label="Nomi" for="title" :error="errors.title">
            <input id="title" v-model="form.title" class="a-input" placeholder="Web platformalar" />
          </FormField>
          <FormField label="Slug" for="slug" :error="errors.slug" hint="Bo‘sh bo‘lsa nomidan yasaladi" optional>
            <input id="slug" v-model="form.slug" class="a-input" placeholder="web-platformalar" />
          </FormField>
          <FormField label="Qisqa tavsif" for="summary" :error="errors.summary" hint="Ro‘yxatda hover paytida ko‘rinadi">
            <textarea id="summary" v-model="form.summary" class="a-textarea" />
          </FormField>
          <FormField label="To‘liq tavsif" for="desc" :error="errors.description" hint="Bo‘sh qator — yangi paragraf">
            <textarea id="desc" v-model="form.description" class="a-textarea a-textarea--lg" />
          </FormField>
          <FormField label="Nima qilamiz (natijalar ro‘yxati)" :error="errors.features">
            <TagInput v-model="form.features" placeholder="Admin panel va rollar…" />
          </FormField>
          <FormField label="Texnologiyalar" :error="errors.stack">
            <TagInput v-model="form.stack" :suggestions="['Vue.js', 'Go', 'PostgreSQL', 'Docker']" />
          </FormField>
        </section>
        <section class="a-card">
          <h2 class="a-card__title">SEO</h2>
          <SeoFields
            v-model="form.seo"
            :fallback-title="form.title"
            :fallback-description="form.summary"
            :path="`/services/${form.slug || '…'}`"
          />
        </section>
      </div>

      <aside class="a-edit__side">
        <EditorSidebar :saving="saving" :dirty="dirty" :is-new="isNew" :public-url="publicUrl" @save="save">
          <FormField label="Holat" for="status">
            <select id="status" v-model="form.status" class="a-select">
              <option value="published">E’lon qilingan</option>
              <option value="draft">Qoralama</option>
            </select>
          </FormField>
        </EditorSidebar>
        <section class="a-card">
          <h2 class="a-card__title">Preview rasm</h2>
          <MediaField v-model="form.preview" />
          <p class="a-hint" style="margin-top: 10px">
            Xizmatlar ro‘yxatida kursor ortidan ko‘rinadi va xizmat sahifasida katta ko‘rsatiladi.
          </p>
        </section>
      </aside>
    </form>
  </div>
</template>
