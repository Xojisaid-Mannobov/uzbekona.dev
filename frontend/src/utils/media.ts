import type { MediaRef } from '@/types/api'

/** Backend yaratgan variantlardan srcset yasaydi: "a-w640.jpg 640w, …, a.jpg 3000w" */
export function srcset(media: MediaRef | null | undefined): string | undefined {
  if (!media || media.kind !== 'image' || !media.variants?.length) return undefined
  const entries = media.variants.map((v) => `${v.url} ${v.width}w`)
  if (media.width) entries.push(`${media.url} ${media.width}w`)
  return entries.join(', ')
}

/** Rasmning tabiiy nisbati (layout shift bo'lmasligi uchun). */
export function aspect(media: MediaRef | null | undefined, fallback = '16 / 10'): string {
  if (media?.width && media.height) return `${media.width} / ${media.height}`
  return fallback
}
