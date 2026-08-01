<template>
  <div class="space-y-3">
    <p class="text-xs text-text-secondary -mt-1">
      Optional — used on the login page and vouchers for this router.
      You can change these later by editing the router.
    </p>

    <div class="grid grid-cols-2 gap-3">
      <label class="flex flex-col gap-1">
        <span class="font-medium">Hotspot name</span>
        <input v-model="hotspotName" class="input" placeholder="e.g. myspot" />
      </label>
      <label class="flex flex-col gap-1">
        <span class="font-medium">DNS name</span>
        <input v-model="dnsName" class="input" placeholder="e.g. myspot.spot" />
      </label>
    </div>

    <div class="flex flex-col gap-1">
      <span class="font-medium">Currency</span>
      <AppSelect
        :model-value="currency || undefined"
        :options="currencyOptions"
        placeholder="None"
        @update:model-value="(v) => (currency = String(v))"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import AppSelect from "@/components/AppSelect.vue";
import { CURRENCIES } from "@/utils/currencies";

const hotspotName = defineModel<string>("hotspotName", { required: true });
const dnsName = defineModel<string>("dnsName", { required: true });
const currency = defineModel<string>("currency", { required: true });

const currencyOptions = CURRENCIES.map((c) => ({ value: c.value, label: c.label }));
</script>
