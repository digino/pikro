<template>
  <AppDialog :open="open" title="Print vouchers" @update:open="emit('update:open', $event)">
    <div class="space-y-4">
      <p class="text-sm text-text-secondary">
        {{ count }} voucher{{ count !== 1 ? "s" : "" }} ready to print
      </p>

      <div class="space-y-2">
        <p class="text-sm font-medium text-text-secondary">Choose a template</p>
        <div class="grid gap-2">
          <button
            v-for="l in VOUCHER_LAYOUTS"
            :key="l.key"
            type="button"
            class="flex items-start gap-3 border rounded-lg p-3 text-left transition-colors"
            :class="
              layout === l.key
                ? 'border-accent bg-accent/10'
                : 'border-border hover:border-muted'
            "
            @click="layout = l.key"
          >
            <span class="mt-0.5 text-lg shrink-0">{{ l.icon }}</span>
            <div>
              <p class="text-sm font-semibold text-text-primary">{{ l.label }}</p>
              <p class="text-xs text-text-secondary">{{ l.description }}</p>
            </div>
          </button>
        </div>
      </div>

      <div class="flex justify-end gap-2 pt-1 border-t border-border">
        <button type="button" class="btn btn-ghost" @click="emit('update:open', false)">
          Cancel
        </button>
        <button type="button" class="btn btn-primary" @click="handlePrint">
          <PrinterIcon class="size-4" /> Print
        </button>
      </div>
    </div>
  </AppDialog>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { PrinterIcon } from "@heroicons/vue/24/outline";
import AppDialog from "@/components/AppDialog.vue";
import {
  printVouchers,
  VOUCHER_LAYOUTS,
  type VoucherEntry,
  type PrintVouchersOptions,
} from "@/utils/vouchers";

const props = defineProps<{
  open: boolean;
  entries: VoucherEntry[];
  defaultLayout: PrintVouchersOptions["layout"];
  businessName: string;
  showValidity: boolean;
  showPrice: boolean;
  currency: string;
  profileMetas: PrintVouchersOptions["profileMetas"];
}>();

const emit = defineEmits<{
  "update:open": [open: boolean];
}>();

const count = ref(0);
const layout = ref<PrintVouchersOptions["layout"]>(props.defaultLayout);

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      count.value = props.entries.length;
      layout.value = props.defaultLayout;
    }
  },
);

function handlePrint() {
  printVouchers(props.entries, {
    layout: layout.value,
    businessName: props.businessName,
    showValidity: props.showValidity,
    showPrice: props.showPrice,
    currency: props.currency,
    profileMetas: props.profileMetas,
  });
  emit("update:open", false);
}
</script>
