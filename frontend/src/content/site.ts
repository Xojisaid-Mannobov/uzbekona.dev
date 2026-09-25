// Saytning doimiy (kod ichidagi) matnlari. Dinamik kontent — admin panel orqali API'dan keladi.

export const navigation = [
  { to: '/projects', label: 'Loyihalar' },
  { to: '/services', label: 'Xizmatlar' },
  { to: '/about', label: 'Biz haqimizda' },
  { to: '/team', label: 'Jamoa' },
  { to: '/news', label: 'Yangiliklar' },
  { to: '/journal', label: 'Journal' },
]

export const footerNavigation = [
  { to: '/projects', label: 'Projects' },
  { to: '/services', label: 'Services' },
  { to: '/about', label: 'About' },
  { to: '/team', label: 'Team' },
  { to: '/news', label: 'News' },
  { to: '/journal', label: 'Journal' },
]

export const hero = {
  tagline: 'Raqamli O‘zbekistonni birga quramiz',
  services: ['Raqamli mahsulotlar', 'Avtomatlashtirish', 'Dasturlash'],
  lead: 'G‘oyalarni ishlaydigan raqamli mahsulotlarga aylantiramiz: web platformalar, mobil ilovalar, Telegram tizimlari va biznes avtomatlashtirish.',
}

export const aboutStatement = {
  label: 'Biz haqimizda',
  lead: 'Biz sayt emas — natija topshiramiz.',
  text: 'Avval biznesingiz qanday ishlashini tushunamiz, keyin uni odamlar zavq bilan ishlatadigan raqamli tizimga aylantiramiz.',
}

export const technologies = [
  { group: 'Frontend', items: ['Vue.js', 'React', 'TypeScript'] },
  { group: 'Backend', items: ['Go', 'Node.js', 'Python', 'Laravel'] },
  { group: 'Database', items: ['PostgreSQL', 'Redis'] },
  { group: 'Infrastructure', items: ['Docker', 'Nginx', 'Cloudflare'] },
  { group: 'AI', items: ['LLM', 'RAG', 'Automation'] },
]

export const processSteps = [
  {
    title: 'Tushunamiz',
    text: 'Biznesingiz, mijozlaringiz va hozirgi jarayonlaringizni o‘rganamiz. Muammo aniq bo‘lmaguncha bitta qator kod ham yozilmaydi.',
  },
  {
    title: 'Loyihalaymiz',
    text: 'Arxitektura va interfeysni chizib, prototipda sinab ko‘rasiz — katta xarajatdan oldin natijani ko‘rasiz.',
  },
  { title: 'Quramiz', text: 'Ikki haftalik bosqichlarda ishlab chiqamiz. Har bosqich oxirida ishlaydigan versiya qo‘lingizda bo‘ladi.' },
  {
    title: 'Sinovdan o‘tkazamiz',
    text: 'Avtomatik testlar, yuklama sinovlari va real foydalanuvchilar bilan tekshiruv — xatolarni mijozlaringizdan oldin topamiz.',
  },
  {
    title: 'Ishga tushiramiz',
    text: 'Monitoring, zaxira nusxalar va xavfsizlik sozlangan holda launch qilamiz. Birinchi kun ham tinch o‘tadi.',
  },
  {
    title: 'O‘stiramiz',
    text: 'Launch — faqat boshlanish. Kelishuv asosida analitika va foydalanuvchilar fikriga tayanib mahsulotingizni rivojlantirib boramiz.',
  },
]

/** "Nega aynan biz" — har bir karta mijozning bitta xavotiriga javob beradi */
export const whyUs = {
  title: 'Nega aynan Uzbekona.dev?',
  lead: 'Chunki biz uchun har bir loyiha — shunchaki shartnoma emas, balki nomimiz va mas’uliyatimiz.',
  cards: [
    {
      key: 'care',
      title: 'Har bir loyihaga o‘zimiznikidek mehr beramiz',
      text: 'Sizning mahsulotingiz biz uchun portfoliodagi navbatdagi qator emas. Uni xuddi o‘zimiz har kuni ishlatadigandek — sinchkovlik bilan, eng kichik tafsilotigacha o‘ylab quramiz.',
    },
    {
      key: 'scale',
      title: 'Katta miqyosga tayyor tizimlar quramiz',
      text: 'Universitetlar va tashkilotlar uchun minglab foydalanuvchi bir vaqtda ishlaydigan platformalar yaratamiz. Yuklama oshsa — tizim to‘xtamaydi.',
    },
    {
      key: 'sprint',
      title: 'Har ikki haftada — ishlaydigan natija',
      text: 'Oylab kutib o‘tirmaysiz. Har bosqich oxirida mahsulotni o‘z ko‘zingiz bilan ko‘rasiz, sinaysiz va yo‘nalishni birga belgilaymiz.',
    },
    {
      key: 'ownership',
      title: 'Kod, server va hujjatlar — to‘liq sizniki',
      text: 'Hech qanday yashirin «garov» yo‘q. Manba kodi, kirish huquqlari va texnik hujjatlar boshidanoq sizning qo‘lingizda.',
    },
    {
      key: 'support',
      title: 'Ishga tushirgandan keyin ham yoningizdamiz',
      text: 'Launch — xayrlashuv emas. Kelishuv asosida monitoring, yangilanishlar va mahsulotni o‘stirishni o‘z zimmamizga olamiz.',
    },
    {
      key: 'trust',
      title: 'Mijozlarimiz ishonchi — bizning eng katta yutug‘imiz',
      text: 'Muvaffaqiyatni topshirilgan loyihalar soni bilan emas, bizga qayta murojaat qilgan va bizni boshqalarga tavsiya qilgan mijozlar bilan o‘lchaymiz.',
    },
  ],
} as const

export const labStageLabels: Record<string, string> = {
  open_source: 'Open Source',
  experimental: 'Experimental',
  in_development: 'In Development',
}

export const projectTypes = [
  'Web platforma',
  'Mobil ilova',
  'Telegram tizimi',
  'Ichki boshqaruv tizimi',
  'Avtomatlashtirish',
  'AI integratsiya',
  'Boshqa',
]

export const budgets = ['$3 000 gacha', '$3 000 – $10 000', '$10 000 – $30 000', '$30 000+', 'Hali aniq emas']
