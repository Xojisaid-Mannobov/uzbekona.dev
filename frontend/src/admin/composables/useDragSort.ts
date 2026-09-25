import { ref, type Ref } from 'vue'

/**
 * HTML5 drag & drop bilan ro'yxatni qayta tartiblash (kutubxonasiz).
 * Klaviatura uchun move(i, ±1) ham bor — har bir qatorda ↑/↓ tugmalari.
 */
export function useDragSort<T>(list: Ref<T[]>, onChange?: () => void) {
  const dragIndex = ref<number | null>(null)
  const overIndex = ref<number | null>(null)

  function move(from: number, to: number) {
    if (to < 0 || to >= list.value.length || from === to) return
    const copy = [...list.value]
    const [item] = copy.splice(from, 1)
    copy.splice(to, 0, item)
    list.value = copy
    onChange?.()
  }

  function handlers(index: number) {
    return {
      draggable: true,
      onDragstart: (e: DragEvent) => {
        dragIndex.value = index
        e.dataTransfer?.setData('text/plain', String(index))
        if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move'
      },
      onDragover: (e: DragEvent) => {
        if (dragIndex.value === null) return
        e.preventDefault()
        overIndex.value = index
      },
      onDrop: (e: DragEvent) => {
        e.preventDefault()
        if (dragIndex.value !== null) move(dragIndex.value, index)
        dragIndex.value = overIndex.value = null
      },
      onDragend: () => {
        dragIndex.value = overIndex.value = null
      },
    }
  }

  return { dragIndex, overIndex, handlers, move }
}
