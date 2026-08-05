<template>
  <div
    v-if="total > pageSize"
    class="flex items-center justify-between pt-3 mt-1 border-t border-border"
  >
    <span class="text-sm font-medium text-text-muted">
      {{ (page - 1) * pageSize + 1 }}–{{ Math.min(page * pageSize, total) }} of {{ total }}
    </span>
    <div class="flex items-center gap-1">
      <button
        class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
        :disabled="page === 1"
        @click="$emit('update:page', page - 1)"
      >
        <ChevronLeftIcon class="size-4" />
      </button>
      <button
        class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
        :disabled="page >= pageCount"
        @click="$emit('update:page', page + 1)"
      >
        <ChevronRightIcon class="size-4" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ChevronLeftIcon, ChevronRightIcon } from "@heroicons/vue/24/outline";

defineProps<{
  page: number;
  pageCount: number;
  pageSize: number;
  total: number;
}>();

defineEmits<{
  "update:page": [page: number];
}>();
</script>
