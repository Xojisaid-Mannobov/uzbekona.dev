import type { IconName } from '@/components/ui/icons'
import type { Block, BlockDataMap, BlockType } from '@/types/api'

interface BlockDef<T extends BlockType> {
  type: T
  label: string
  description: string
  icon: IconName
  create: () => BlockDataMap[T]
}

type Defs = { [K in BlockType]: BlockDef<K> }

// Content builder'dagi 14 ta blok turi (TZ 44-band)
export const blockDefs: Defs = {
  heading: {
    type: 'heading',
    label: 'Sarlavha',
    description: 'Bo‘lim sarlavhasi va label',
    icon: 'heading',
    create: () => ({ text: '', label: '' }),
  },
  text: { type: 'text', label: 'Matn', description: 'Oddiy paragraflar', icon: 'text', create: () => ({ text: '' }) },
  large_text: {
    type: 'large_text',
    label: 'Katta matn',
    description: 'Yirik bayonot (48–64px)',
    icon: 'text',
    create: () => ({ text: '' }),
  },
  image: { type: 'image', label: 'Rasm', description: 'Container kengligida', icon: 'image', create: () => ({ media: null, caption: '' }) },
  full_image: {
    type: 'full_image',
    label: 'To‘liq ekran rasm',
    description: 'Viewport kengligida',
    icon: 'image',
    create: () => ({ media: null, caption: '' }),
  },
  gallery: { type: 'gallery', label: 'Galereya', description: 'Bir nechta rasm, 2 ustun', icon: 'gallery', create: () => ({ items: [] }) },
  video: {
    type: 'video',
    label: 'Video',
    description: 'MP4 yoki YouTube havola',
    icon: 'video',
    create: () => ({ media: null, poster: null, url: '', caption: '' }),
  },
  stats: {
    type: 'stats',
    label: 'Raqamlar',
    description: 'Natija metrikalari',
    icon: 'chart',
    create: () => ({ items: [{ value: '', label: '' }] }),
  },
  quote: {
    type: 'quote',
    label: 'Iqtibos',
    description: 'Mijoz yoki jamoa fikri',
    icon: 'quote',
    create: () => ({ text: '', author: '', role: '' }),
  },
  two_columns: {
    type: 'two_columns',
    label: 'Ikki ustun',
    description: 'Masalan: Muammo / Yechim',
    icon: 'columns',
    create: () => ({
      columns: [
        { title: '', text: '' },
        { title: '', text: '' },
      ],
    }),
  },
  three_columns: {
    type: 'three_columns',
    label: 'Uch ustun',
    description: 'Uchta qisqa blok',
    icon: 'columns',
    create: () => ({
      columns: [
        { title: '', text: '' },
        { title: '', text: '' },
        { title: '', text: '' },
      ],
    }),
  },
  technology: {
    type: 'technology',
    label: 'Texnologiyalar',
    description: 'Stack ro‘yxati',
    icon: 'cpu',
    create: () => ({ title: 'Texnologiyalar', items: [] }),
  },
  process: {
    type: 'process',
    label: 'Jarayon',
    description: 'Raqamlangan bosqichlar',
    icon: 'list',
    create: () => ({ steps: [{ title: '', text: '' }] }),
  },
  before_after: {
    type: 'before_after',
    label: 'Oldin / Keyin',
    description: 'Taqqoslash slayderi',
    icon: 'compare',
    create: () => ({ before: null, after: null, caption: '' }),
  },
}

export const blockTypeList = Object.values(blockDefs) as BlockDef<BlockType>[]

export function createBlock(type: BlockType): Block {
  return { type, data: blockDefs[type].create() } as Block
}

/** Yig'ilgan blok sarlavhasida ko'rinadigan qisqa mazmun. */
export function blockSummary(block: Block): string {
  const d = block.data as Record<string, unknown>
  const text = (d.text as string) || (d.caption as string) || (d.title as string) || ''
  if (text) return text.length > 80 ? `${text.slice(0, 80)}…` : text
  if (Array.isArray(d.items)) return `${d.items.length} ta element`
  if (Array.isArray(d.steps)) return `${d.steps.length} ta bosqich`
  if (Array.isArray(d.columns))
    return (d.columns as { title: string }[])
      .map((c) => c.title)
      .filter(Boolean)
      .join(' · ')
  return ''
}
