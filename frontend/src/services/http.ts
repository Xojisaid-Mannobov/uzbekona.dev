import axios, { AxiosError } from 'axios'
import type { ApiErrorBody } from '@/types/api'

// Frontend backend bilan faqat shu REST API orqali ishlaydi
export const http = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? '/api/v1',
  withCredentials: true,
  timeout: 30_000,
  headers: { Accept: 'application/json' },
})

/** Yagona xato turi: UI uchun xabar va maydonlar bo'yicha xatolar. */
export class ApiError extends Error {
  status: number
  code: string
  fields: Record<string, string>

  constructor(status: number, body: Partial<ApiErrorBody>) {
    super(body.message || 'Xatolik yuz berdi. Qayta urinib ko‘ring.')
    this.status = status
    this.code = body.code || 'unknown'
    this.fields = body.fields || {}
  }
}

export function toApiError(err: unknown): ApiError {
  if (err instanceof ApiError) return err
  if (err instanceof AxiosError) {
    if (!err.response) {
      return new ApiError(0, { code: 'network', message: 'Internet bilan bog‘lanishda muammo.' })
    }
    const body = (err.response.data as { error?: ApiErrorBody } | undefined)?.error
    return new ApiError(err.response.status, body ?? {})
  }
  return new ApiError(0, {})
}

http.interceptors.response.use(
  (res) => res,
  (err) => Promise.reject(toApiError(err)),
)
