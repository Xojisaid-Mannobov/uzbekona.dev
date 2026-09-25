-- Uzbekona.dev — boshlang'ich kontent.
-- Barcha matnlarni admin paneldan tahrirlash mumkin.
-- Metrikalar va kontaktlar real ma'lumotlarga almashtirilishi shart.

-- ─── Sozlamalar ──────────────────────────────────────────────
INSERT INTO settings (key, value) VALUES
('site', $j${
  "name": "Uzbekona.dev",
  "tagline": "Digital Product Studio",
  "email": "hello@uzbekona.dev",
  "phone": "",
  "telegram": "@uzbekona_dev",
  "address": "Toshkent, O‘zbekiston",
  "available": true
}$j$),
('metrics', $j$[
  {"value": "75K+", "label": "Foydalanuvchilar"},
  {"value": "12+", "label": "Mahsulot"},
  {"value": "1M+", "label": "So‘rov"},
  {"value": "99.9%", "label": "Uptime"}
]$j$),
('seo', $j${
  "title": "Uzbekona.dev — Digital Product Studio",
  "description": "Uzbekona.dev web platformalar, mobil ilovalar, Telegram tizimlari va biznes jarayonlarini avtomatlashtiruvchi raqamli mahsulotlar ishlab chiqadi."
}$j$);

INSERT INTO social_links (platform, label, url, position) VALUES
('github', 'GitHub', 'https://github.com/uzbekona-dev', 1),
('telegram', 'Telegram', 'https://t.me/uzbekona_dev', 2),
('instagram', 'Instagram', 'https://instagram.com/uzbekona.dev', 3),
('linkedin', 'LinkedIn', 'https://linkedin.com/company/uzbekona-dev', 4);

-- ─── Xizmatlar ───────────────────────────────────────────────
INSERT INTO services (slug, title, summary, description, features, stack, position) VALUES
('web-platformalar', 'Web platformalar',
 'Universitetlar, kompaniyalar va startaplar uchun yuklamaga chidamli web tizimlar.',
 'Kabinetlar, portallar, marketplace va SaaS mahsulotlar. Biz interfeysdan tortib ma’lumotlar bazasigacha bo‘lgan butun zanjirni loyihalaymiz: foydalanuvchi rollari, integratsiyalar, hisobotlar va admin panel bir tizim sifatida ishlaydi.',
 ARRAY['Foydalanuvchi kabinetlari', 'Admin panel va rollar', 'Tashqi API integratsiyalari', 'SEO va performance optimizatsiya'],
 ARRAY['Vue.js', 'TypeScript', 'Go', 'PostgreSQL'], 1),
('mobil-ilovalar', 'Mobil ilovalar',
 'iOS va Android uchun bitta kod bazasidagi tezkor, tabiiy his beruvchi ilovalar.',
 'Mobil ilova — bu web saytning kichraytirilgan nusxasi emas. Offline rejim, push-bildirishnomalar, to‘lovlar va qurilma imkoniyatlarini hisobga olgan holda mahsulotni telefon uchun qayta o‘ylaymiz.',
 ARRAY['iOS va Android', 'Push-bildirishnomalar', 'Offline rejim', 'App Store / Google Play’ga chiqarish'],
 ARRAY['React Native', 'TypeScript', 'Go'], 2),
('telegram-tizimlari', 'Telegram tizimlari',
 'Botlar, Mini App’lar va kanallarni boshqaruvchi to‘liq tizimlar.',
 'Telegram — O‘zbekistondagi eng katta raqamli kanal. Biz oddiy javob beruvchi botlar emas, balki to‘lov, CRM, admin panel va analitikaga ega Telegram mahsulotlarini quramiz.',
 ARRAY['Telegram botlar', 'Telegram Mini Apps', 'To‘lov integratsiyasi', 'Admin panel va analitika'],
 ARRAY['Go', 'Node.js', 'PostgreSQL', 'Redis'], 3),
('ichki-boshqaruv-tizimlari', 'Ichki boshqaruv tizimlari',
 'Excel va qog‘ozni almashtiradigan CRM, ERP va hujjat aylanish tizimlari.',
 'Kompaniyangiz ichidagi jarayonlarni o‘rganamiz va ularni tushunarli interfeysga ega tizimga aylantiramiz: arizalar, tasdiqlash zanjirlari, ombor, xodimlar va hisobotlar.',
 ARRAY['CRM va ERP', 'Hujjat aylanishi', 'Rollar va ruxsatlar', 'Hisobot va dashboardlar'],
 ARRAY['Vue.js', 'Go', 'PostgreSQL'], 4),
('avtomatlashtirish', 'Avtomatlashtirish',
 'Takrorlanuvchi qo‘l mehnatini ishonchli avtomatik jarayonlarga aylantiramiz.',
 'Hisobotlarni yig‘ish, ma’lumotlarni tizimlar orasida ko‘chirish, bildirishnomalar va integratsiyalar — xodimlar vaqtini qimmatli ishlarga bo‘shatadigan avtomatlashtirish.',
 ARRAY['Integratsiyalar', 'Rejalashtirilgan vazifalar', 'Ma’lumotlar sinxronizatsiyasi', 'Bildirishnomalar'],
 ARRAY['Go', 'Python', 'Redis'], 5),
('ai-va-llm', 'AI va LLM integratsiyalar',
 'Kompaniya ma’lumotlari ustida ishlaydigan AI yordamchilar va RAG tizimlar.',
 'LLM’larni real biznes jarayonlariga ulaymiz: hujjatlar bo‘yicha qidiruv, mijozlarga javob beruvchi yordamchilar, matnlarni tasniflash va ma’lumot ajratib olish.',
 ARRAY['RAG tizimlar', 'AI yordamchilar', 'Hujjatlarni tahlil qilish', 'LLM integratsiyasi'],
 ARRAY['Python', 'Go', 'PostgreSQL', 'LLM'], 6),
('product-design', 'Product Design',
 'Tadqiqotdan prototipgacha: foydalanuvchi uchun tushunarli interfeyslar.',
 'Dizayn — bu bezak emas, qaror. Foydalanuvchi ssenariylarini, ma’lumotlar oqimini va interfeys ierarxiyasini koddan oldin loyihalaymiz.',
 ARRAY['UX tadqiqot', 'Wireframe va prototip', 'Dizayn tizimi', 'UI dizayn'],
 ARRAY['Figma'], 7),
('infrastructure', 'Infrastructure',
 'Mahsulot barqaror ishlashi uchun server, CI/CD va monitoring.',
 'Docker, Nginx va Cloudflare asosidagi infratuzilma, avtomatik deploy, zaxira nusxalar va monitoring — launch’dan keyin ham tizim tinch ishlashi uchun.',
 ARRAY['Docker va CI/CD', 'Monitoring va loglar', 'Zaxira nusxalar', 'Xavfsizlik'],
 ARRAY['Docker', 'Nginx', 'Cloudflare', 'Linux'], 8);

-- ─── Journal kategoriyalari ──────────────────────────────────
INSERT INTO article_categories (name, slug, position) VALUES
('Engineering', 'engineering', 1),
('Product', 'product', 2),
('Design', 'design', 3),
('AI', 'ai', 4),
('Infrastructure', 'infrastructure', 5);

-- ─── Loyihalar ───────────────────────────────────────────────
INSERT INTO projects (slug, title, tagline, short_description, full_description, year, client, industry, platforms, services, stack, accent, status, featured, position, published_at) VALUES
('kuaf', 'KUAF', 'Universitet uchun yagona raqamli ekotizim.',
 'University Digital Ecosystem',
 'Talabalar, o‘qituvchilar va ma’muriyat uchun bitta platforma: qabul, o‘quv jarayoni, hujjatlar va xabarlar yagona tizimda.',
 2026, 'KUAF', 'Education', ARRAY['Web'], ARRAY['Web platforma', 'Product Design', 'Infrastructure'],
 ARRAY['Vue.js', 'TypeScript', 'Go', 'PostgreSQL', 'Redis'], '#2350F5', 'published', true, 1, now()),
('otmedu', 'OTMEDU', 'Oliy ta’lim uchun raqamli boshqaruv platformasi.',
 'Higher Education Platform',
 'Oliy ta’lim muassasalari uchun o‘quv jarayoni va ma’lumotlarni boshqarish platformasi.',
 2025, 'OTMEDU', 'Education', ARRAY['Web'], ARRAY['Web platforma', 'Ichki boshqaruv tizimi'],
 ARRAY['Vue.js', 'Go', 'PostgreSQL'], '#0F766E', 'published', true, 2, now()),
('mylearn', 'MyLearn', 'O‘rganishni kundalik odatga aylantiruvchi platforma.',
 'Learning Platform',
 'Kurslar, topshiriqlar va natijalarni kuzatish uchun onlayn ta’lim platformasi.',
 2025, 'MyLearn', 'EdTech', ARRAY['Web', 'Mobile'], ARRAY['Web platforma', 'Mobil ilova', 'Product Design'],
 ARRAY['Vue.js', 'TypeScript', 'Go'], '#B4461E', 'published', true, 3, now());

INSERT INTO project_blocks (project_id, type, position, data)
SELECT p.id, b.type, b.position, b.data::jsonb
FROM projects p
JOIN (VALUES
  ('kuaf', 'large_text', 1, $j${"text": "Universitetdagi har bir jarayon — qabul, dars jadvali, baholar, arizalar — alohida jadval va chat’larda yashardi. Maqsad ularni bitta tushunarli tizimga birlashtirish edi."}$j$),
  ('kuaf', 'two_columns', 2, $j${"columns": [
      {"title": "Muammo", "text": "Ma’lumotlar turli joylarda saqlanar, har bir bo‘lim o‘z vositasidan foydalanar va talaba bitta savol uchun bir nechta joyga murojaat qilishi kerak edi."},
      {"title": "Yechim", "text": "Rollarga asoslangan yagona platforma: talaba kabineti, o‘qituvchi paneli va ma’muriyat uchun boshqaruv tizimi bitta ma’lumotlar bazasi ustida ishlaydi."}
  ]}$j$),
  ('kuaf', 'heading', 3, $j${"text": "Qanday ishladik", "label": "Jarayon"}$j$),
  ('kuaf', 'process', 4, $j${"steps": [
      {"title": "Tadqiqot", "text": "Bo‘limlar bilan suhbatlar va mavjud jarayonlar xaritasi."},
      {"title": "Arxitektura", "text": "Rollar, ma’lumotlar modeli va integratsiyalarni loyihalash."},
      {"title": "Ishlab chiqish", "text": "Bosqichma-bosqich reliz va real foydalanuvchilar bilan test."},
      {"title": "Qo‘llab-quvvatlash", "text": "Monitoring, fikr-mulohazalar va doimiy rivojlantirish."}
  ]}$j$),
  ('kuaf', 'technology', 5, $j${"title": "Texnologiyalar", "items": ["Vue.js", "TypeScript", "Go", "PostgreSQL", "Redis", "Docker", "Nginx"]}$j$),
  ('otmedu', 'large_text', 1, $j${"text": "Oliy ta’lim muassasalari uchun o‘quv jarayonini raqamlashtirish — hujjatlardan hisobotlargacha."}$j$),
  ('otmedu', 'technology', 2, $j${"title": "Texnologiyalar", "items": ["Vue.js", "Go", "PostgreSQL", "Docker"]}$j$),
  ('mylearn', 'large_text', 1, $j${"text": "Kurslar, topshiriqlar va progressni bitta qulay interfeysda jamlagan ta’lim platformasi."}$j$),
  ('mylearn', 'technology', 2, $j${"title": "Texnologiyalar", "items": ["Vue.js", "TypeScript", "Go", "PostgreSQL"]}$j$)
) AS b (slug, type, position, data) ON b.slug = p.slug;

-- ─── Jamoa ───────────────────────────────────────────────────
INSERT INTO team_members (name, role, bio, socials, position) VALUES
('Xojisaid Mannopov', 'Founder · Full-stack Engineer',
 'Mahsulot arxitekturasi, backend va frontend. Jarayonni tushunishdan boshlab production’gacha.',
 '[{"platform": "github", "url": "https://github.com/Xojisaid-Mannobov"}]', 1);

-- ─── Uzbekona Labs ───────────────────────────────────────────
INSERT INTO labs (slug, title, description, stage, stack, position) VALUES
('telekit', 'Telekit', 'Telegram botlarni tez va tartibli qurish uchun Go asosidagi toolkit: routing, state, middleware.', 'open_source', ARRAY['Go', 'Telegram Bot API'], 1),
('save-flow', 'Save Flow', 'Kontentni saqlash va tartiblash jarayonini avtomatlashtiruvchi servis.', 'in_development', ARRAY['Go', 'PostgreSQL', 'Redis'], 2),
('ai-tools', 'AI tools', 'Kundalik ish uchun LLM asosidagi kichik vositalar: qisqartirish, tasniflash, matndan ma’lumot ajratish.', 'experimental', ARRAY['Python', 'LLM', 'RAG'], 3),
('dev-utilities', 'Developer utilities', 'Jamoa ichida ishlatiladigan CLI va kutubxonalar: migratsiya, deploy, loglar.', 'open_source', ARRAY['Go', 'Docker'], 4);

-- ─── Journal ─────────────────────────────────────────────────
INSERT INTO articles (slug, title, excerpt, category_id, content, reading_time, author_name, status, featured, published_at)
SELECT a.slug, a.title, a.excerpt, c.id, a.content::jsonb, a.reading_time, 'Xojisaid Mannopov', 'published', a.featured, now() - a.age
FROM (VALUES
  ('nega-go', 'Nega backend uchun Go tanladik', 'Oddiy deploy, bashorat qilinadigan performance va o‘qilishi oson kod — Go’ni tanlashimizning uch sababi.', 'engineering', true, interval '2 days', 4,
   $j$[
     {"type": "text", "data": {"text": "Studiya sifatida biz bir vaqtning o‘zida bir nechta mahsulotni qo‘llab-quvvatlaymiz. Shuning uchun backend tili tanlovida eng muhim mezon — tizimni bir yildan keyin ham tushunish va o‘zgartirish oson bo‘lishi edi."}},
     {"type": "heading", "data": {"text": "1. Bitta binary — oddiy deploy"}},
     {"type": "text", "data": {"text": "Go dasturi bitta binary faylga yig‘iladi. Docker image kichik, ishga tushirish tez, serverda qo‘shimcha runtime kerak emas."}},
     {"type": "heading", "data": {"text": "2. Bashorat qilinadigan performance"}},
     {"type": "text", "data": {"text": "Goroutine’lar va samarali standart kutubxona tufayli bitta kichik server minglab so‘rovlarni barqaror qayta ishlaydi."}},
     {"type": "quote", "data": {"text": "Yaxshi arxitektura — bu murakkablikni yashirish emas, uni to‘g‘ri joyga qo‘yish.", "author": "Uzbekona.dev"}},
     {"type": "heading", "data": {"text": "3. O‘qilishi oson kod"}},
     {"type": "text", "data": {"text": "Go tili ataylab sodda. Yangi dasturchi loyihaga qo‘shilganda kodni tez tushunadi — bu esa mahsulotni uzoq muddat rivojlantirish uchun eng muhim omil."}}
   ]$j$),
  ('telegram-bot-mahsulot', 'Telegram bot — bu mahsulot, skript emas', 'Production’dagi bot uchun state, admin panel, analitika va rate limit nega kerakligi haqida.', 'product', false, interval '9 days', 5,
   $j$[
     {"type": "text", "data": {"text": "Ko‘pchilik Telegram botni bir kechada yoziladigan skript deb biladi. Lekin bot minglab foydalanuvchiga xizmat qila boshlaganda, u to‘laqonli mahsulotga aylanadi."}},
     {"type": "heading", "data": {"text": "State va ma’lumotlar"}},
     {"type": "text", "data": {"text": "Foydalanuvchi qaysi bosqichda ekanini xotirada emas, bazada saqlash kerak — aks holda har bir restart suhbatni uzib qo‘yadi."}},
     {"type": "heading", "data": {"text": "Admin panel va analitika"}},
     {"type": "text", "data": {"text": "Biznes egasi bot ichida nima bo‘layotganini ko‘rishi, xabarlarni boshqarishi va natijalarni o‘lchashi kerak. Busiz bot — qora quti."}},
     {"type": "heading", "data": {"text": "Limitlar va barqarorlik"}},
     {"type": "text", "data": {"text": "Telegram API limitlari, navbatlar va qayta urinish mexanizmlari — ommaviy xabar yuborishda tizim yiqilmasligining kafolati."}}
   ]$j$),
  ('rag-qachon-kerak', 'RAG qachon kerak va qachon kerak emas', 'Kompaniya hujjatlari ustida AI yordamchi qurishdan oldin javob berish kerak bo‘lgan savollar.', 'ai', false, interval '16 days', 4,
   $j$[
     {"type": "text", "data": {"text": "Retrieval-Augmented Generation (RAG) — LLM’ga savolga javob berishdan oldin kerakli hujjat parchalarini topib beradigan yondashuv. U kuchli, lekin har doim ham zarur emas."}},
     {"type": "heading", "data": {"text": "RAG kerak bo‘lganda"}},
     {"type": "text", "data": {"text": "Ma’lumotlar tez-tez yangilansa, hujjatlar hajmi katta bo‘lsa va javob manbasini ko‘rsatish muhim bo‘lsa."}},
     {"type": "heading", "data": {"text": "RAG ortiqcha bo‘lganda"}},
     {"type": "text", "data": {"text": "Agar vazifa oddiy tasniflash yoki qisqa qoidalar to‘plamiga asoslangan bo‘lsa, yaxshi yozilgan prompt va aniq strukturali ma’lumot yetarli."}}
   ]$j$)
) AS a (slug, title, excerpt, category_slug, featured, age, reading_time, content)
JOIN article_categories c ON c.slug = a.category_slug;
