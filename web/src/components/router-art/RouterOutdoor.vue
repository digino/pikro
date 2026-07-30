<!-- LHG / SXT family: outdoor CPE with dish reflector -->
<template>
  <svg :width="size" :height="size" viewBox="0 0 120 100" xmlns="http://www.w3.org/2000/svg">
    <!-- Mounting pole -->
    <rect x="57" y="72" width="6" height="22" rx="2" fill="var(--color-border)"/>
    <!-- Dish reflector (ellipse arc) -->
    <ellipse cx="60" cy="48" rx="44" ry="30" fill="var(--color-muted)" stroke="var(--color-border)" stroke-width="1.5"/>
    <!-- Inner dish concave -->
    <ellipse cx="60" cy="50" rx="38" ry="24" fill="var(--color-surface)" stroke="var(--color-border)" stroke-width="1" opacity="0.6"/>
    <!-- Feed arm -->
    <line x1="60" y1="50" x2="60" y2="28" stroke="var(--color-border)" stroke-width="2.5" stroke-linecap="round"/>
    <!-- Feed horn (LNB unit) -->
    <rect x="54" y="20" width="12" height="10" rx="3" fill="var(--color-muted)" stroke="var(--color-border)" stroke-width="1.2"/>
    <!-- LED on feed unit -->
    <circle cx="60" cy="24" r="2.5" :fill="powerLed" style="transition:fill 0.4s"/>
    <!-- Dish ribs -->
    <line x1="60" y1="50" x2="20" y2="35" stroke="var(--color-border)" stroke-width="1" opacity="0.5"/>
    <line x1="60" y1="50" x2="100" y2="35" stroke="var(--color-border)" stroke-width="1" opacity="0.5"/>
    <line x1="60" y1="50" x2="32" y2="62" stroke="var(--color-border)" stroke-width="1" opacity="0.5"/>
    <line x1="60" y1="50" x2="88" y2="62" stroke="var(--color-border)" stroke-width="1" opacity="0.5"/>
    <!-- Mount bracket -->
    <rect x="52" y="68" width="16" height="8" rx="2" fill="var(--color-border)"/>
  </svg>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ size?: number; powerLed: string }>(), {
  size: 120,
})
</script>
