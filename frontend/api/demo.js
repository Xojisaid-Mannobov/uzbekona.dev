// Vercel demo rejimi: backend (Go + PostgreSQL) ulanmagan, shuning uchun yozish amallari
// (kontakt forma, admin panel) aniq xabar bilan rad etiladi — hech narsa jimgina yo'qolmaydi.
export default function handler(req, res) {
  const isAdmin = (req.url || '').includes('/admin')
  res.setHeader('Cache-Control', 'no-store')
  res.status(503).json({
    error: {
      code: 'demo_mode',
      message: isAdmin
        ? 'Demo rejim: admin panel Go backend ulangan serverda ishlaydi (docker compose).'
        : 'Demo rejim: forma hozircha qabul qilinmaydi. Iltimos, Telegram yoki email orqali yozing.',
    },
  })
}
