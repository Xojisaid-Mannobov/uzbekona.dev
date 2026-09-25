<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import AppIcon from '@/components/ui/AppIcon.vue'
import type { IconName } from '@/components/ui/icons'
import BrandLogo from '@/components/layout/BrandLogo.vue'
import ToastHost from '@/admin/components/ToastHost.vue'
import ConfirmDialog from '@/admin/components/ConfirmDialog.vue'
import { useAuthStore } from '@/admin/stores/auth'
import { useTheme } from '@/composables/useTheme'
import { adminApi } from '@/admin/services/adminApi'
import '@/admin/styles/admin.css'

const auth = useAuthStore()
const route = useRoute()
const router = useRouter()
const theme = useTheme()
const sidebarOpen = ref(false)
const newRequests = ref(0)
const newApplications = ref(0)

const nav: { to: string; label: string; icon: IconName; badge?: 'requests' | 'applications' }[][] = [
  [
    { to: '/admin/dashboard', label: 'Dashboard', icon: 'home' },
    { to: '/admin/analytics', label: 'Statistika', icon: 'pulse' },
  ],
  [
    { to: '/admin/projects', label: 'Loyihalar', icon: 'folder' },
    { to: '/admin/services', label: 'Xizmatlar', icon: 'layers' },
    { to: '/admin/team', label: 'Jamoa', icon: 'users' },
    { to: '/admin/labs', label: 'Labs', icon: 'flask' },
    { to: '/admin/news', label: 'Yangiliklar', icon: 'news' },
    { to: '/admin/articles', label: 'Maqolalar', icon: 'file-text' },
    { to: '/admin/media', label: 'Media', icon: 'image' },
  ],
  [
    { to: '/admin/requests', label: 'So‘rovlar', icon: 'inbox', badge: 'requests' },
    { to: '/admin/applications', label: 'Nomzodlar', icon: 'briefcase', badge: 'applications' },
    { to: '/admin/settings', label: 'Sozlamalar', icon: 'settings' },
    { to: '/admin/users', label: 'Adminlar', icon: 'user' },
  ],
]

const title = computed(() => route.meta.title ?? 'Admin')

// Yangi so'rovlar soni — sidebar'dagi badge uchun
async function refreshCounts() {
  try {
    const [req, apps] = await Promise.all([
      adminApi.requests.list({ status: 'new', limit: 1 }),
      adminApi.applications.list({ status: 'new', limit: 1 }),
    ])
    newRequests.value = req.meta.total
    newApplications.value = apps.meta.total
  } catch {
    /* badge ixtiyoriy */
  }
}

onMounted(refreshCounts)
watch(
  () => route.fullPath,
  () => {
    sidebarOpen.value = false
    if (/^\/admin\/(requests|applications|dashboard)/.test(route.path)) refreshCounts()
  },
)

async function logout() {
  await auth.logout()
  router.push('/admin/login')
}
</script>

<template>
  <div class="admin shell" :class="{ 'is-open': sidebarOpen }">
    <aside class="side" aria-label="Admin navigatsiya">
      <RouterLink to="/admin/dashboard" class="side__brand"><BrandLogo /></RouterLink>

      <nav class="side__nav">
        <div v-for="(group, gi) in nav" :key="gi" class="side__group">
          <RouterLink v-for="item in group" :key="item.to" :to="item.to" class="side__link">
            <AppIcon :name="item.icon" :size="19" />
            <span>{{ item.label }}</span>
            <span v-if="item.badge === 'requests' && newRequests" class="side__badge">{{ newRequests }}</span>
            <span v-if="item.badge === 'applications' && newApplications" class="side__badge">{{ newApplications }}</span>
          </RouterLink>
        </div>
      </nav>

      <div class="side__foot">
        <a href="/" target="_blank" class="side__link"><AppIcon name="external" :size="19" /><span>Saytni ochish</span></a>
      </div>
    </aside>

    <div v-if="sidebarOpen" class="scrim" @click="sidebarOpen = false" />

    <div class="main">
      <header class="top">
        <button class="a-icon-btn top__burger" type="button" aria-label="Menyu" @click="sidebarOpen = !sidebarOpen">
          <AppIcon name="menu" />
        </button>
        <p class="top__title">{{ title }}</p>
        <div class="top__right">
          <button class="a-icon-btn" type="button" :aria-label="theme.isDark() ? 'Yorug‘ rejim' : 'Qorong‘i rejim'" @click="theme.toggle()">
            <AppIcon :name="theme.isDark() ? 'sun' : 'moon'" :size="19" />
          </button>
          <div class="top__user">
            <span class="top__avatar">{{ auth.admin?.name?.charAt(0) ?? 'A' }}</span>
            <span class="top__name">{{ auth.admin?.name }}</span>
          </div>
          <button class="a-icon-btn" type="button" aria-label="Chiqish" title="Chiqish" @click="logout">
            <AppIcon name="logout" :size="19" />
          </button>
        </div>
      </header>

      <main class="content">
        <RouterView v-slot="{ Component }">
          <Transition name="fade" mode="out-in">
            <component :is="Component" :key="route.path" />
          </Transition>
        </RouterView>
      </main>
    </div>

    <ToastHost />
    <ConfirmDialog />
  </div>
</template>

<style scoped>
.shell {
  display: grid;
  grid-template-columns: var(--a-side-w) minmax(0, 1fr);
  min-height: 100vh;
  background: var(--bg);
}

.side {
  position: sticky;
  top: 0;
  height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 20px 14px;
  border-right: 1px solid var(--line);
  background: var(--surface);
}

.side__brand {
  padding: 6px 10px 22px;
}

.side__brand :deep(.brand) {
  font-size: 19px;
}

.side__nav {
  flex: 1;
  display: grid;
  align-content: start;
  gap: 18px;
  overflow-y: auto;
}

.side__group {
  display: grid;
  gap: 2px;
}

.side__group + .side__group {
  padding-top: 18px;
  border-top: 1px solid var(--line);
}

.side__link {
  display: flex;
  align-items: center;
  gap: 12px;
  height: 42px;
  padding: 0 12px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 14px;
  color: var(--ink-2);
  transition:
    background-color var(--dur-fast) var(--ease),
    color var(--dur-fast) var(--ease);
}

.side__link:hover {
  background: var(--surface-2);
  color: var(--ink);
}

.side__link.router-link-active {
  background: var(--ink);
  color: var(--bg);
}

.side__badge {
  margin-left: auto;
  min-width: 22px;
  height: 22px;
  padding: 0 7px;
  border-radius: 99px;
  background: var(--accent);
  color: #fff;
  font-size: 12px;
  display: grid;
  place-items: center;
}

.side__foot {
  padding-top: 14px;
  border-top: 1px solid var(--line);
}

.main {
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.top {
  position: sticky;
  top: 0;
  z-index: 20;
  height: var(--a-top-h);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 28px;
  background: color-mix(in srgb, var(--bg) 82%, transparent);
  backdrop-filter: blur(14px);
  border-bottom: 1px solid var(--line);
}

.top__burger {
  display: none;
}

.top__title {
  font-weight: 700;
  font-size: 15px;
}

.top__right {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 8px;
}

.top__user {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 8px;
}

.top__avatar {
  display: grid;
  place-items: center;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: var(--accent);
  color: #fff;
  font-weight: 700;
  font-size: 14px;
}

.top__name {
  font-weight: 600;
  font-size: 14px;
}

.content {
  padding: 28px;
  max-width: 1440px;
  width: 100%;
}

.scrim {
  display: none;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 180ms var(--ease);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media (max-width: 1024px) {
  .shell {
    grid-template-columns: 1fr;
  }

  .side {
    position: fixed;
    z-index: 60;
    left: 0;
    width: var(--a-side-w);
    transform: translateX(-100%);
    transition: transform var(--dur) var(--ease);
  }

  .is-open .side {
    transform: none;
    box-shadow: var(--shadow-lg);
  }

  .scrim {
    display: block;
    position: fixed;
    inset: 0;
    z-index: 55;
    background: rgb(0 0 0 / 0.3);
  }

  .top__burger {
    display: inline-grid;
  }

  .top {
    padding: 0 16px;
  }

  .top__name {
    display: none;
  }

  .content {
    padding: 20px 16px;
  }
}
</style>
