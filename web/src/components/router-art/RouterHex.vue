<!-- hEX / RB750 family: small plastic box, no WiFi, 5 LAN ports -->
<template>
  <svg :width="size" :height="size" viewBox="0 0 120 100" xmlns="http://www.w3.org/2000/svg">
    <!-- Body shadow -->
    <rect x="14" y="62" width="92" height="12" rx="5" fill="var(--color-border)" opacity="0.4"/>
    <!-- Main body — taller/squarer than hAP -->
    <rect x="12" y="34" width="96" height="32" rx="6" fill="var(--color-muted)" stroke="var(--color-border)" stroke-width="1.5"/>
    <!-- Top face bevel -->
    <rect x="16" y="34" width="88" height="8" rx="4" fill="var(--color-border)" opacity="0.5"/>
    <!-- Brand label area -->
    <rect x="16" y="45" width="30" height="12" rx="2" fill="var(--color-border)" opacity="0.3"/>
    <!-- LEDs on label area -->
    <circle cx="22" cy="51" r="2.5" :fill="powerLed" style="transition:fill 0.4s"/>
    <circle cx="30" cy="51" r="2.5" :fill="wanLed"   style="transition:fill 0.4s"/>
    <!-- 5 Ethernet ports on right side -->
    <rect x="52" y="42" width="52" height="16" rx="2" fill="var(--color-surface)" stroke="var(--color-border)" stroke-width="1"/>
    <rect x="55" y="45" width="7" height="10" rx="1" fill="var(--color-border)"/>
    <rect x="65" y="45" width="7" height="10" rx="1" fill="var(--color-border)"/>
    <rect x="75" y="45" width="7" height="10" rx="1" fill="var(--color-border)"/>
    <rect x="85" y="45" width="7" height="10" rx="1" fill="var(--color-border)"/>
    <rect x="95" y="45" width="7" height="10" rx="1" fill="var(--color-border)"/>
    <!-- Feet -->
    <rect x="20" y="65" width="10" height="5" rx="2" fill="var(--color-border)"/>
    <rect x="90" y="65" width="10" height="5" rx="2" fill="var(--color-border)"/>
  </svg>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ size?: number; powerLed: string; wanLed: string }>(), {
  size: 120,
})
</script>
