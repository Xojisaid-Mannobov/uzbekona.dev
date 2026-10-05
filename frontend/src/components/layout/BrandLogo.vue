<script setup lang="ts">
// Rasmiy logo (belgi + "uzbekona.dev"). Qorong'i mavzu uchun matni yorug' variant ishlatiladi.
// `onDark` — doim qorong'i fonda turuvchi joylar (footer) uchun: mavzudan qat'i nazar yorug' matnli variant.
withDefaults(defineProps<{ size?: 'md' | 'lg'; onDark?: boolean }>(), { size: 'md', onDark: false })
</script>

<template>
  <span class="brand" :class="[`brand--${size}`, { 'brand--on-dark': onDark }]">
    <img class="brand__img brand__img--light" src="/logo.png" alt="Uzbekona.dev" width="1000" height="238" decoding="async" />
    <img class="brand__img brand__img--dark" src="/logo-dark.png" alt="" width="1000" height="238" decoding="async" aria-hidden="true" />
  </span>
</template>

<style scoped>
.brand {
  display: inline-flex;
  align-items: center;
}

.brand__img {
  display: block;
  height: 38px;
  width: auto;
}

.brand--lg .brand__img {
  height: 52px;
}

.brand__img--dark {
  display: none;
}

.brand--on-dark .brand__img--light {
  display: none;
}

.brand--on-dark .brand__img--dark {
  display: block;
}
</style>

<style>
/* Mavzu qoidalari scoped emas: :root holatiga bog'liq */
@media (prefers-color-scheme: dark) {
  :root:not([data-theme='light']) .brand:not(.brand--on-dark) .brand__img--light {
    display: none;
  }

  :root:not([data-theme='light']) .brand:not(.brand--on-dark) .brand__img--dark {
    display: block;
  }
}

:root[data-theme='dark'] .brand:not(.brand--on-dark) .brand__img--light {
  display: none;
}

:root[data-theme='dark'] .brand:not(.brand--on-dark) .brand__img--dark {
  display: block;
}
</style>
