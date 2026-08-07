<template>
  <div
    class="text-center"
    :class="[sizeClasses, bordered ? 'rounded-xl border border-border bg-surface' : '']"
  >
    <slot>
      <p v-if="title" class="font-semibold text-text-primary mb-1">{{ title }}</p>
      <p class="text-text-secondary" :class="textClasses">{{ message }}</p>
    </slot>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";

const props = withDefaults(
  defineProps<{
    title?: string;
    message?: string;
    size?: "sm" | "md" | "lg";
    bordered?: boolean;
  }>(),
  { size: "md" },
);

const sizeClasses = computed(
  () => ({ sm: "py-2", md: "py-8", lg: "py-12" })[props.size],
);
const textClasses = computed(
  () => ({ sm: "text-xs", md: "text-sm", lg: "text-lg" })[props.size],
);
</script>
