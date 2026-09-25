import type { CoverFocus, CoverRatio, MediaRef } from '@/types/api'

// Muqova rasmi kartada qaysi nisbatda ko'rinishi va kesilganda qaysi qismi saqlanishi (admin tanlaydi)

export const COVER_RATIOS: { value: CoverRatio; label: string }[] = [
  { value: 'auto', label: 'Asl' },
  { value: '16:9', label: '16:9' },
  { value: '4:3', label: '4:3' },
  { value: '1:1', label: '1:1' },
  { value: '3:4', label: '3:4' },
  { value: '21:9', label: '21:9' },
]

export const COVER_FOCUS: { value: CoverFocus; label: string }[] = [
  { value: 'top', label: 'Yuqori' },
  { value: 'center', label: 'Markaz' },
  { value: 'bottom', label: 'Past' },
]

/** CSS aspect-ratio: tanlangan nisbat → rasmning asl o'lchami → standart (fallback) */
export function ratioCss(ratio: CoverRatio | undefined, media: MediaRef | null | undefined, fallback: string): string {
  if (ratio && ratio !== 'auto') return ratio.replace(':', ' / ')
  if (media?.width && media?.height) return `${media.width} / ${media.height}`
  return fallback
}

/** CSS object-position */
export function focusCss(focus: CoverFocus | undefined): string {
  return focus === 'top' ? 'center top' : focus === 'bottom' ? 'center bottom' : 'center'
}
