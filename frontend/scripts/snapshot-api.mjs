// Ishlab turgan Go API'dan public ma'lumotlarni statik JSON sifatida saqlaydi.
// Vercel kabi faqat statik/serverless hostingda sayt backend'siz "demo rejim"da ishlashi uchun.
//
//   API_URL=http://localhost:8080 node scripts/snapshot-api.mjs
import { mkdir, rm, writeFile } from 'node:fs/promises'
import { dirname, join } from 'node:path'

const api = (process.env.API_URL ?? 'http://localhost:8080').replace(/\/$/, '') + '/api/v1'
const out = new URL('../public/demo-api/', import.meta.url).pathname

async function fetchJSON(path) {
  const res = await fetch(api + path)
  if (!res.ok) throw new Error(`${path}: ${res.status}`)
  return res.json()
}

async function save(path, body) {
  const file = join(out, `${path}.json`)
  await mkdir(dirname(file), { recursive: true })
  await writeFile(file, JSON.stringify(body))
  console.log('✓', path)
}

await rm(out, { recursive: true, force: true })

for (const p of ['settings', 'team', 'labs', 'article-categories']) await save(p, await fetchJSON(`/${p}`))

const projects = await fetchJSON('/projects')
await save('projects', projects)
for (const p of projects.data) await save(`projects/${p.slug}`, await fetchJSON(`/projects/${p.slug}`))

const services = await fetchJSON('/services')
await save('services', services)
for (const s of services.data) await save(`services/${s.slug}`, await fetchJSON(`/services/${s.slug}`))

const articles = await fetchJSON('/articles?limit=50')
await save('articles', articles)
for (const a of articles.data) await save(`articles/${a.slug}`, await fetchJSON(`/articles/${a.slug}`))

const news = await fetchJSON('/news?limit=50')
await save('news', news)
for (const n of news.data) await save(`news/${n.slug}`, await fetchJSON(`/news/${n.slug}`))

const categories = await fetchJSON('/article-categories')
for (const c of categories.data) await save(`articles-category/${c.slug}`, await fetchJSON(`/articles?limit=50&category=${c.slug}`))
