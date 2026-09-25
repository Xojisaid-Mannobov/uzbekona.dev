-- Demo yangiliklar (admin paneldan tahrirlanadi yoki o'chiriladi)
INSERT INTO news (slug, title, excerpt, tag, content, status, pinned, published_at)
SELECT n.slug, n.title, n.excerpt, n.tag, n.content::jsonb, 'published', n.pinned, now() - n.age
FROM (VALUES
  ('uzbekona-dev-yangi-sayt', 'Uzbekona.dev yangi saytini ishga tushirdi', 'Milliy naqshlar, Registon ruhidagi dizayn va o‘zimiz yozgan kontent boshqaruv tizimi — yangi sayt qanday qurilgani haqida.', 'Studiya', true, interval '1 day',
   $j$[
     {"type": "text", "data": {"text": "Yangi sayt — bizning vizit kartamiz. Shuning uchun uni mijozlarimiz uchun qanday qursak, xuddi shunday qurdik: tadqiqot, dizayn, ishlab chiqish va sinov."}},
     {"type": "heading", "data": {"text": "Milliy ruh va zamonaviy texnologiya"}},
     {"type": "text", "data": {"text": "Dizaynda girih va suzani naqshlari, Registon silueti va bayroq ranglari ishlatildi. Barcha bezaklar SVG’da chizilgan — sayt tez ochiladi va har qanday ekranda aniq ko‘rinadi."}},
     {"type": "quote", "data": {"text": "Raqamli O‘zbekistonni birga quramiz.", "author": "Uzbekona.dev"}},
     {"type": "heading", "data": {"text": "O‘zimizning boshqaruv paneli"}},
     {"type": "text", "data": {"text": "Loyihalar, maqolalar, yangiliklar va tashriflar statistikasi — hammasi bitta admin panelda. Tayyor shablonlar ishlatilmadi."}}
   ]$j$),
  ('kuaf-platformasi-yangilandi', 'KUAF raqamli ekotizimining yangi versiyasi', 'Talabalar kabineti, dars jadvali va hujjat aylanishi bitta tizimda — platformaning yangi imkoniyatlari.', 'Loyiha', false, interval '6 days',
   $j$[
     {"type": "text", "data": {"text": "Universitet uchun yaratilgan raqamli ekotizimning yangi versiyasi foydalanuvchilarga taqdim etildi."}},
     {"type": "heading", "data": {"text": "Nimalar o‘zgardi"}},
     {"type": "text", "data": {"text": "Talabalar kabineti tezlashdi, dars jadvali mobil qurilmalarga moslashdi, hujjatlar esa endi elektron tarzda imzolanadi."}}
   ]$j$),
  ('telekit-open-source', 'Telekit — Telegram botlar uchun ochiq kodli toolkit', 'Uzbekona Labs’dagi birinchi open-source loyiha: Go’da Telegram botlarni tez va tartibli qurish uchun vositalar.', 'Labs', false, interval '12 days',
   $j$[
     {"type": "text", "data": {"text": "Mijozlar uchun o‘nlab Telegram tizimlarini qurish jarayonida takrorlanuvchi yechimlarni bitta kutubxonaga jamladik va uni ochiq kodli qildik."}},
     {"type": "text", "data": {"text": "Routing, holat (state) boshqaruvi va middleware — hammasi tayyor. Hamjamiyatdan fikr va hissa kutamiz."}}
   ]$j$),
  ('jamoaga-backend-dasturchi', 'Jamoamizga backend dasturchi qidiryapmiz', 'Go va PostgreSQL bilan ishlashni yaxshi ko‘rasizmi? Katta miqyosdagi loyihalarda birga ishlaymiz.', 'Jamoa', false, interval '20 days',
   $j$[
     {"type": "text", "data": {"text": "Loyihalar ko‘payishi bilan jamoamizni kengaytiryapmiz. Go, PostgreSQL va Docker bilan tajribaga ega backend dasturchini kutamiz."}},
     {"type": "text", "data": {"text": "Rezyume va GitHub profilingizni bog‘lanish sahifasi orqali yuboring."}}
   ]$j$)
) AS n(slug, title, excerpt, tag, pinned, age, content)
ON CONFLICT (slug) DO NOTHING;
