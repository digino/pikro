<template>
  <div class="space-y-3">
    <p class="text-xs text-text-secondary -mt-1">
      {{ t('routers.hotspotFields.hint') }}
    </p>

    <label class="flex flex-col gap-1">
      <span class="font-medium">{{ t('routers.hotspotFields.name') }}</span>
      <input v-model="hotspotName" class="input" :placeholder="t('routers.hotspotFields.namePlaceholder')" />
    </label>

    <div class="flex flex-col gap-1">
      <span class="font-medium">{{ t('routers.hotspotFields.currency') }}</span>
      <AppSelect
        :model-value="currency || undefined"
        :options="currencyOptions"
        :placeholder="t('routers.hotspotFields.none')"
        @update:model-value="(v) => (currency = String(v))"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import AppSelect from "@/components/AppSelect.vue";
import { CURRENCIES } from "@/utils/currencies";

const hotspotName = defineModel<string>("hotspotName", { required: true });
const currency = defineModel<string>("currency", { required: true });

const { t } = useI18n();

const currencyOptions = CURRENCIES.map((c) => ({ value: c.value, label: c.label }));
</script>
