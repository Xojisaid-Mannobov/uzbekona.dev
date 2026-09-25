// Yangilik teglari uchun rang ohanglari — bayroq va milliy bezak ranglari
export type NewsTone = 'blue' | 'green' | 'gold' | 'teal' | 'red'

const KNOWN: Record<string, NewsTone> = {
  studiya: 'blue',
  loyiha: 'green',
  labs: 'gold',
  jamoa: 'teal',
  tadbir: 'red',
  hamkorlik: 'teal',
  mahsulot: 'green',
}
const ORDER: NewsTone[] = ['blue', 'green', 'gold', 'teal', 'red']

export function toneOf(tag: string): NewsTone {
  const key = tag.trim().toLowerCase()
  if (KNOWN[key]) return KNOWN[key]
  let h = 0
  for (const ch of key) h = (h * 31 + ch.charCodeAt(0)) >>> 0
  return ORDER[h % ORDER.length]
}

const SHORT_MONTHS = ['yan', 'fev', 'mar', 'apr', 'may', 'iyun', 'iyul', 'avg', 'sen', 'okt', 'noy', 'dek']

/** Taqvim varag'i uchun: kun, qisqa oy nomi va yil */
export function leafDate(iso: string | null) {
  const d = iso ? new Date(iso) : new Date()
  return { day: d.getDate(), month: SHORT_MONTHS[d.getMonth()], year: d.getFullYear() }
}

export const TONE_COLOR: Record<NewsTone, string> = {
  blue: 'var(--accent)',
  green: '#16994a',
  gold: 'var(--gold)',
  teal: '#0f8f86',
  red: '#d9262c',
}

/** Teg rangi — kartalar va taqvim varag'i uchun CSS o'zgaruvchisi */
export const toneStyle = (tag: string) => ({ '--tone': TONE_COLOR[toneOf(tag)] })
