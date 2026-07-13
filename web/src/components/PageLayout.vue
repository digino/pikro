<template>
  <div class="flex flex-col h-full min-h-screen bg-base">
    <!-- Topbar -->
    <div class="flex items-center gap-3 px-6 border-b border-border h-14 shrink-0 bg-base">
      <div class="flex items-center gap-2 min-w-0">
        <span class="text-sm font-semibold text-text-primary truncate">{{ title }}</span>
        <template v-if="subtitle">
          <span class="text-text-muted shrink-0">/</span>
          <span class="text-text-secondary text-sm truncate">{{ subtitle }}</span>
        </template>
      </div>
      <div class="ml-auto flex items-center gap-2 shrink-0">
        <slot name="actions" />
      </div>
    </div>

    <!-- Page body -->
    <div class="flex-1 overflow-y-auto min-w-0 p-6 flex flex-col gap-6">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ title: string; subtitle?: string }>();
</script>
