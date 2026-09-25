// Vercel demo rejimi: statistika bazasi yo'q, shuning uchun tashrif shunchaki qabul qilinadi (204).
// To'liq statistika Go backend + PostgreSQL ulangan serverda ishlaydi.
export default function handler(req, res) {
  res.setHeader('Cache-Control', 'no-store')
  res.status(204).end()
}
