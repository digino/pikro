<template>
  <span
    v-if="variant === 'pill'"
    class="text-sm font-semibold px-2 py-0.5 rounded-full"
    :class="[colorClasses.bg, colorClasses.text]"
  >
    {{ label }}
  </span>
  <span v-else class="inline-flex items-center gap-1" :class="colorClasses.text">
    <span class="size-1.5 rounded-full shrink-0" :class="colorClasses.dot" />
    {{ label }}
  </span>
</template>

<script setup lang="ts">
import { computed } from "vue";

export type StatusColor = "green" | "amber" | "red" | "blue" | "muted";

const props = withDefaults(
  defineProps<{
    label: string;
    color: StatusColor;
    variant?: "pill" | "dot";
  }>(),
  { variant: "dot" },
);

const COLORS: Record<StatusColor, { bg: string; text: string; dot: string }> = {
  green: { bg: "bg-green/10", text: "text-green", dot: "bg-green" },
  amber: { bg: "bg-amber/10", text: "text-amber", dot: "bg-amber" },
  red: { bg: "bg-red/10", text: "text-red", dot: "bg-red" },
  blue: { bg: "bg-blue-500/10", text: "text-blue-400", dot: "bg-blue-400" },
  muted: { bg: "bg-muted", text: "text-text-muted", dot: "bg-text-muted" },
};

const colorClasses = computed(() => COLORS[props.color]);
</script>
