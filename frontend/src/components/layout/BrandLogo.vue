<script setup lang="ts">
import BrandMark from './BrandMark.vue'

// Belgi (lola shaklidagi "U") + wordmark: "Uzbekona" + accent rangli ".dev"
withDefaults(defineProps<{ size?: 'md' | 'lg' }>(), { size: 'md' })
</script>

<template>
  <span class="brand" :class="`brand--${size}`">
    <BrandMark class="brand__mark" :code="false" />
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
  height: 1.5em;
}

.brand__tld {
  color: var(--accent);
}
</style>
