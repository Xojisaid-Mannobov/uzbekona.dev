// Saytning doimiy (kod ichidagi) matnlari. Dinamik kontent — admin panel orqali API'dan keladi.

export const navigation = [
  { to: '/projects', label: 'Loyihalar' },
  { to: '/services', label: 'Xizmatlar' },
  { to: '/about', label: 'Biz haqimizda' },
  { to: '/team', label: 'Jamoa' },
  { to: '/journal', label: 'Journal' },
]

export const footerNavigation = [
  { to: '/projects', label: 'Projects' },
  { to: '/services', label: 'Services' },
  { to: '/about', label: 'About' },
  { to: '/team', label: 'Team' },
  { to: '/journal', label: 'Journal' },
]

export const hero = {
  tagline: 'Raqamli O‘zbekistonni birga quramiz',
  services: ['Raqamli mahsulotlar', 'Avtomatlashtirish', 'Dasturlash'],
  lead: 'G‘oyalarni ishlaydigan raqamli mahsulotlarga aylantiramiz: web platformalar, mobil ilovalar, Telegram tizimlari va biznes avtomatlashtirish.',
  words: ['G‘oyalar', 'Odamlar', 'Texnologiya', 'Yorug‘ kelajak'],
}

export const aboutStatement = {
  label: 'Biz haqimizda',
  lead: 'Biz shunchaki sayt yozmaymiz.',
  text: 'Jarayonlarni tushunamiz, ularni raqamlashtiramiz va odamlar ishlata oladigan tizimlarga aylantiramiz.',
}

export const technologies = [
  { group: 'Frontend', items: ['Vue.js', 'React', 'TypeScript'] },
  { group: 'Backend', items: ['Go', 'Node.js', 'Python', 'Laravel'] },
  { group: 'Database', items: ['PostgreSQL', 'Redis'] },
  { group: 'Infrastructure', items: ['Docker', 'Nginx', 'Cloudflare'] },
  { group: 'AI', items: ['LLM', 'RAG', 'Automation'] },
]

export const processSteps = [
  { title: 'Tushunamiz', text: 'Biznes, foydalanuvchilar va mavjud jarayonlarni o‘rganamiz. Muammo aniq bo‘lmaguncha kod yozilmaydi.' },
  { title: 'Loyihalaymiz', text: 'Arxitektura, ma’lumotlar modeli va interfeys ssenariylarini loyihalaymiz. Prototipda tekshiramiz.' },
  { title: 'Quramiz', text: 'Qisqa iteratsiyalarda ishlab chiqamiz — har bosqich oxirida ishlaydigan natijani ko‘rasiz.' },
  { title: 'Test qilamiz', text: 'Avtomatik testlar, yuklama sinovlari va real foydalanuvchilar bilan tekshiruv.' },
  { title: 'Ishga tushiramiz', text: 'Production infratuzilma, monitoring va zaxira nusxalar bilan xavfsiz launch.' },
  { title: 'Rivojlantiramiz', text: 'Launch — boshlanish. Analitika va fikr-mulohazalar asosida mahsulotni o‘stiramiz.' },
]

export const principles = [
  { title: 'Product thinking', text: 'Koddan oldin muammoni tushunamiz.' },
  { title: 'Engineering', text: 'Production uchun ishlaydigan tizim quramiz.' },
  { title: 'Scalability', text: 'Bugun ishlaydigan emas, ertaga ham o‘sadigan arxitektura.' },
  { title: 'Support', text: 'Launch’dan keyin ham mahsulot rivojlantiriladi.' },
]

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
