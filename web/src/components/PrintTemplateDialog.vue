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
            v-for="t in VOUCHER_TEMPLATES"
            :key="t.key"
            type="button"
            class="flex items-start gap-3 border rounded-lg p-3 text-left transition-colors"
            :class="
              template === t.key
                ? 'border-accent bg-accent/10'
                : 'border-border hover:border-muted'
            "
            @click="template = t.key"
          >
            <div>
              <p class="text-sm font-semibold text-text-primary">{{ t.label }}</p>
              <p class="text-xs text-text-secondary">{{ t.description }}</p>
            </div>
          </button>
        </div>
      </div>

      <p v-if="printError" class="text-xs text-red">{{ printError }}</p>

      <div class="flex justify-end gap-2 pt-1 border-t border-border">
        <button type="button" class="btn btn-ghost" @click="emit('update:open', false)">
          Cancel
        </button>
        <button type="button" class="btn btn-primary" :disabled="printing" @click="handlePrint">
          <span v-if="printing" class="size-3.5 border-2 border-black/20 border-t-black rounded-full animate-spin" />
          <PrinterIcon v-else class="size-4" /> Print
        </button>
      </div>
    </div>
  </AppDialog>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { PrinterIcon } from "@heroicons/vue/24/outline";
import AppDialog from "@/components/AppDialog.vue";
import { VOUCHER_TEMPLATES, type VoucherTemplate } from "@/utils/voucherTemplates";
import { printVouchers, type VoucherEntry, type PrintVouchersOptions } from "@/utils/vouchers";

const props = defineProps<{
  open: boolean;
  entries: VoucherEntry[];
  businessName: string;
  currency: string;
  profileMetas: PrintVouchersOptions["profileMetas"];
  loginUrl?: string;
  loginUrlSupportsCredentials?: boolean;
}>();

const emit = defineEmits<{
  "update:open": [open: boolean];
}>();

const count = ref(0);
const template = ref<VoucherTemplate["key"]>("classic");
const printing = ref(false);
const printError = ref("");

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      count.value = props.entries.length;
      printError.value = "";
    }
  },
);

async function handlePrint() {
  printing.value = true;
  printError.value = "";
  try {
    await printVouchers(props.entries, {
      template: template.value,
      businessName: props.businessName,
      showValidity: true,
      showPrice: true,
      currency: props.currency,
      profileMetas: props.profileMetas,
      loginUrl: props.loginUrl,
      loginUrlSupportsCredentials: props.loginUrlSupportsCredentials,
    });
    emit("update:open", false);
  } catch (e: any) {
    printError.value = e?.message ?? "Failed to prepare vouchers for printing";
  } finally {
    printing.value = false;
  }
}
</script>
