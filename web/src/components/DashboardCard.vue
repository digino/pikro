<template>
  <div class="rounded-xl border border-border bg-surface flex flex-col h-full">
    <div class="flex items-center justify-between gap-2 px-5 py-3">
      <span class="font-semibold text-text-primary">{{ title }}</span>
      <div class="flex items-center gap-2">
        <slot name="actions" />
      </div>
    </div>
    <div class="border-t border-border" />
    <div class="p-5 flex-1 flex flex-col gap-3 min-h-0">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ title: string }>();
</script>
