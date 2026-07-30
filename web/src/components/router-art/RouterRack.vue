<!-- CCR / RB4011 / CRS family: 1U rack-mount unit -->
<template>
  <svg :width="size" :height="size" viewBox="0 0 120 100" xmlns="http://www.w3.org/2000/svg">
    <!-- Rack ears left -->
    <rect x="2" y="38" width="10" height="24" rx="2" fill="var(--color-border)"/>
    <circle cx="7" cy="50" r="2.5" fill="var(--color-surface)"/>
    <!-- Rack ears right -->
    <rect x="108" y="38" width="10" height="24" rx="2" fill="var(--color-border)"/>
    <circle cx="113" cy="50" r="2.5" fill="var(--color-surface)"/>
    <!-- Body shadow -->
    <rect x="13" y="63" width="94" height="8" rx="2" fill="var(--color-border)" opacity="0.4"/>
    <!-- Main 1U body -->
    <rect x="12" y="36" width="96" height="28" rx="3" fill="var(--color-muted)" stroke="var(--color-border)" stroke-width="1.5"/>
    <!-- Top bevel -->
    <rect x="14" y="36" width="92" height="5" rx="2" fill="var(--color-border)" opacity="0.4"/>
    <!-- Left panel: LEDs + model -->
    <circle cx="20" cy="50" r="2.5" :fill="powerLed" style="transition:fill 0.4s"/>
    <circle cx="28" cy="50" r="2.5" :fill="actLed"   style="transition:fill 0.4s"/>
    <!-- Lots of ports — SFP + Ethernet group -->
    <!-- SFP cages (2) -->
    <rect x="38" y="42" width="12" height="16" rx="1.5" fill="var(--color-surface)" stroke="var(--color-border)" stroke-width="1"/>
    <rect x="53" y="42" width="12" height="16" rx="1.5" fill="var(--color-surface)" stroke="var(--color-border)" stroke-width="1"/>
    <!-- Ethernet bank (10 ports) -->
    <rect x="70" y="42" width="38" height="16" rx="2" fill="var(--color-surface)" stroke="var(--color-border)" stroke-width="1"/>
    <rect x="72" y="44" width="5" height="12" rx="1" fill="var(--color-border)"/>
    <rect x="79" y="44" width="5" height="12" rx="1" fill="var(--color-border)"/>
    <rect x="86" y="44" width="5" height="12" rx="1" fill="var(--color-border)"/>
    <rect x="93" y="44" width="5" height="12" rx="1" fill="var(--color-border)"/>
    <rect x="100" y="44" width="5" height="12" rx="1" fill="var(--color-border)"/>
    <!-- Vent slots far left -->
    <rect x="16" y="41" width="16" height="2" rx="1" fill="var(--color-border)" opacity="0.5"/>
    <rect x="16" y="45" width="16" height="2" rx="1" fill="var(--color-border)" opacity="0.5"/>
    <rect x="16" y="49" width="16" height="2" rx="1" fill="var(--color-border)" opacity="0.5"/>
    <rect x="16" y="53" width="16" height="2" rx="1" fill="var(--color-border)" opacity="0.5"/>
    <rect x="16" y="57" width="16" height="2" rx="1" fill="var(--color-border)" opacity="0.5"/>
  </svg>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ size?: number; powerLed: string; actLed: string }>(), {
  size: 120,
})
</script>
