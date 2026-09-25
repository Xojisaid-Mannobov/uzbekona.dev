import { computed, onBeforeUnmount, onMounted, ref, type Ref } from 'vue'
import { onBeforeRouteLeave, useRouter, type RouteLocationRaw } from 'vue-router'
import { useUiStore } from '@/admin/stores/ui'
import { toApiError, type ApiError } from '@/services/http'

interface EditorOptions<E, F, I> {
  /** null — yangi yozuv yaratish */
  id: Ref<number | null>
  empty: () => F
  fromEntity: (e: E) => F
  toInput: (f: F) => I
  load: (id: number) => Promise<E>
  create: (input: I) => Promise<E>
  update: (id: number, input: I) => Promise<E>
  /** Yaratilgandan keyin o'tiladigan tahrirlash sahifasi */
  editRoute: (e: E) => RouteLocationRaw
}

/**
 * Admin tahrirlash sahifalari uchun umumiy mantiq:
 * yuklash → forma → saqlash, server xatolarini maydonlarga bog'lash,
 * saqlanmagan o'zgarishlar haqida ogohlantirish.
 */
export function useEditor<E, F, I>(opts: EditorOptions<E, F, I>) {
  const ui = useUiStore()
  const router = useRouter()

  const form = ref(opts.empty()) as Ref<F>
  const entity = ref<E | null>(null) as Ref<E | null>
  const loading = ref(false)
  const loadError = ref<ApiError | null>(null)
  const saving = ref(false)
  const errors = ref<Record<string, string>>({})
  const snapshot = ref(JSON.stringify(form.value))

  const dirty = computed(() => JSON.stringify(form.value) !== snapshot.value)
  const isNew = computed(() => opts.id.value === null)

  function reset(value: F) {
    form.value = value
    snapshot.value = JSON.stringify(value)
  }

  async function load() {
    if (opts.id.value === null) return
    loading.value = true
    loadError.value = null
    try {
      entity.value = await opts.load(opts.id.value)
      reset(opts.fromEntity(entity.value))
    } catch (e) {
      loadError.value = toApiError(e)
    } finally {
      loading.value = false
    }
  }

  async function save(): Promise<boolean> {
    if (saving.value) return false
    saving.value = true
    errors.value = {}
    try {
      const input = opts.toInput(form.value)
      const result = opts.id.value === null ? await opts.create(input) : await opts.update(opts.id.value, input)
      entity.value = result
      reset(opts.fromEntity(result))
      ui.success(opts.id.value === null ? 'Muvaffaqiyatli yaratildi' : 'Saqlandi ✓')
      if (opts.id.value === null) router.replace(opts.editRoute(result))
      return true
    } catch (e) {
      const err = toApiError(e)
      errors.value = err.fields
      ui.error(Object.keys(err.fields).length ? 'Ma’lumotlarni tekshiring — xatolar belgilangan' : err.message)
      return false
    } finally {
      saving.value = false
    }
  }

  // Ctrl/Cmd + S — tezkor saqlash
  function onKey(e: KeyboardEvent) {
    if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 's') {
      e.preventDefault()
      save()
    }
  }

  function onBeforeUnload(e: BeforeUnloadEvent) {
    if (dirty.value) e.preventDefault()
  }

  onMounted(() => {
    load()
    window.addEventListener('keydown', onKey)
    window.addEventListener('beforeunload', onBeforeUnload)
  })

  onBeforeUnmount(() => {
    window.removeEventListener('keydown', onKey)
    window.removeEventListener('beforeunload', onBeforeUnload)
  })

  onBeforeRouteLeave(async () => {
    if (!dirty.value || saving.value) return true
    return ui.confirm({
      title: 'Saqlanmagan o‘zgarishlar bor',
      text: 'Sahifadan chiqsangiz, kiritilgan o‘zgarishlar yo‘qoladi.',
      confirmLabel: 'Saqlamasdan chiqish',
    })
  })

  return { form, entity, loading, loadError, saving, errors, dirty, isNew, save, load }
}
