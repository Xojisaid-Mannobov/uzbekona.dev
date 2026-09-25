<script setup lang="ts">
// Label + kontrol + xato/izoh. Kontrolning o'zi slot orqali beriladi.
defineProps<{ label?: string; error?: string; hint?: string; for?: string; optional?: boolean; span?: boolean }>()
</script>

<template>
  <div class="a-field" :class="{ 'has-error': !!error, 'a-span-2': span }">
    <label v-if="label" class="a-label" :for="$props.for">
      {{ label }}
      <small v-if="optional">ixtiyoriy</small>
    </label>
    <slot />
    <p v-if="error" class="a-error">{{ error }}</p>
    <p v-else-if="hint" class="a-hint">{{ hint }}</p>
  </div>
</template>
