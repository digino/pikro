<!-- hAP / RB951 family: flat box with internal antennas (no external sticks) -->
<template>
  <svg :width="size" :height="size" viewBox="0 0 120 100" xmlns="http://www.w3.org/2000/svg">
    <!-- Drop shadow -->
    <rect x="10" y="72" width="100" height="10" rx="6" fill="var(--color-border)" opacity="0.35"/>
    <!-- Main body -->
    <rect x="8" y="36" width="104" height="38" rx="8" fill="var(--color-muted)" stroke="var(--color-border)" stroke-width="1.5"/>
    <!-- Top face inset -->
    <rect x="12" y="40" width="96" height="18" rx="5" fill="var(--color-surface)" stroke="var(--color-border)" stroke-width="1" opacity="0.6"/>
    <!-- LED row -->
    <circle cx="18" cy="62" r="2.5" :fill="powerLed" style="transition:fill 0.4s"/>
    <circle cx="27" cy="62" r="2.5" :fill="wifiLed"  style="transition:fill 0.4s"/>
    <circle cx="36" cy="62" r="2.5" :fill="wanLed"   style="transition:fill 0.4s"/>
    <!-- Ethernet ports -->
    <rect x="58" y="57" width="46" height="13" rx="2" fill="var(--color-surface)" stroke="var(--color-border)" stroke-width="1"/>
    <rect x="61" y="60" width="7" height="7" rx="1" fill="var(--color-border)"/>
    <rect x="71" y="60" width="7" height="7" rx="1" fill="var(--color-border)"/>
    <rect x="81" y="60" width="7" height="7" rx="1" fill="var(--color-border)"/>
    <rect x="91" y="60" width="7" height="7" rx="1" fill="var(--color-border)"/>
    <!-- USB port -->
    <rect x="47" y="58" width="8" height="11" rx="1.5" fill="var(--color-surface)" stroke="var(--color-border)" stroke-width="1"/>
    <!-- Rubber feet -->
    <rect x="16" y="72" width="10" height="4" rx="2" fill="var(--color-border)"/>
    <rect x="94" y="72" width="10" height="4" rx="2" fill="var(--color-border)"/>
  </svg>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ size?: number; powerLed: string; wifiLed: string; wanLed: string }>(), {
  size: 120,
})
</script>
