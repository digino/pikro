<template>
  <div class="rounded-xl border border-border p-4 bg-surface">
    <div class="flex items-center gap-1.5 text-xs font-medium uppercase tracking-widest text-text-muted mb-2">
      <component :is="icon" class="size-3.5 shrink-0" />
      {{ label }}
    </div>
    <p class="text-2xl font-bold text-text-primary truncate">{{ value || '—' }}</p>
  </div>
</template>

<script setup lang="ts">
import type { FunctionalComponent } from 'vue'
defineProps<{ label: string; value: string; icon?: FunctionalComponent }>()
</script>
