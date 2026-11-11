<script setup lang="ts">
// A user's time left, as in the Users lists: a bar, amber when low, and the label.
defineProps<{ label: string; pct: number }>();
</script>

<template>
  <div class="flex items-center gap-2.5">
    <div class="w-18 h-2 shrink-0 rounded-full border-[1.5px] border-ink bg-paper overflow-hidden">
      <div class="h-full" :class="pct < 40 ? 'bg-primary' : 'bg-ink'" :style="{ width: `${pct}%` }" />
    </div>
    <span class="text-[13px]">{{ label }}</span>
  </div>
</template>
