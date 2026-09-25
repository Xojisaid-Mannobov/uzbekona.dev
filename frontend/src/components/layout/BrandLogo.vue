<script setup lang="ts">
// Wordmark: "Uzbekona" + accent rangli ".dev"
withDefaults(defineProps<{ size?: 'md' | 'lg' }>(), { size: 'md' })
</script>

<template>
  <span class="brand" :class="`brand--${size}`">
    <svg class="brand__mark" viewBox="0 0 32 32" aria-hidden="true">
      <rect width="32" height="32" rx="9" fill="currentColor" />
      <path d="M10 9v8.5a6 6 0 0 0 12 0V9" fill="none" stroke="var(--bg)" stroke-width="3" stroke-linecap="round" />
    </svg>
    <span class="brand__word">Uzbekona<span class="brand__tld">.dev</span></span>
  </span>
</template>

<style scoped>
.brand {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  font-weight: 700;
  letter-spacing: -0.035em;
  color: var(--ink);
}

.brand--md {
  font-size: 20px;
}

.brand--lg {
  font-size: 24px;
}

.brand__mark {
  width: 1.45em;
  height: 1.45em;
}

.brand__tld {
  color: var(--accent);
}
</style>
