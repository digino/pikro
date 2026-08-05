<template>
  <AppDialog :open="open" title="Print vouchers" size="xl" @update:open="emit('update:open', $event)">
    <div class="flex gap-5 items-start">
      <div class="w-56 shrink-0 space-y-4">
        <p>
          <span class="font-semibold">{{ count }} voucher{{ count !== 1 ? "s" : "" }}</span> ready to print
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
      </div>

      <!-- Preview -->
      <div class="flex-1 min-w-0">
        <div class="border border-border rounded-xl overflow-hidden bg-base" style="height: 420px">
          <div v-if="previewLoading" class="h-full flex items-center justify-center">
            <span class="spinner spinner--sm" />
          </div>
          <iframe
            v-else
            :srcdoc="previewHtml"
            class="w-full h-full border-0"
            sandbox=""
            title="Voucher preview"
          />
        </div>
      </div>
    </div>

    <div class="flex justify-end gap-2 pt-4 mt-4 border-t border-border">
      <button type="button" class="btn btn-ghost" @click="emit('update:open', false)">
        Cancel
      </button>
      <button type="button" class="btn btn-primary" :disabled="printing" @click="handlePrint">
        <span v-if="printing" class="size-3.5 border-2 border-black/20 border-t-black rounded-full animate-spin" />
        <PrinterIcon v-else class="size-4" /> Print
      </button>
    </div>
  </AppDialog>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { PrinterIcon } from "@heroicons/vue/24/outline";
import AppDialog from "@/components/AppDialog.vue";
import { VOUCHER_TEMPLATES, getVoucherTemplate, type VoucherTemplate } from "@/utils/voucherTemplates";
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
const previewHtml = ref("");
const previewLoading = ref(false);

// Preview shows a handful of real entries (or sample data if none were
// passed in yet) so the admin sees actual usernames/passwords, not filler.
const PREVIEW_COUNT = 6;
const sampleEntries: VoucherEntry[] = [
  { name: "ab3f", password: "ab3f" },
  { name: "x9kz", password: "x9kz" },
  { name: "p7mw", password: "p7mw" },
  { name: "q2nt", password: "q2nt" },
];

async function renderPreview() {
  previewLoading.value = true;
  try {
    const entries = (props.entries.length ? props.entries : sampleEntries).slice(0, PREVIEW_COUNT);
    const items = entries.map((e) => {
      const meta = props.profileMetas[e.profile || "default"];
      return {
        name: e.name,
        password: e.password,
        validity: meta?.validity ?? "1 day",
        price: meta?.price ? `${meta.price}${props.currency ? " " + props.currency : ""}` : `500${props.currency ? " " + props.currency : ""}`,
      };
    });
    previewHtml.value = await getVoucherTemplate(template.value).render(items, {
      businessName: props.businessName,
      loginUrl: props.loginUrl ?? "",
      loginUrlSupportsCredentials: props.loginUrlSupportsCredentials ?? false,
    });
  } finally {
    previewLoading.value = false;
  }
}

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      count.value = props.entries.length;
      printError.value = "";
      renderPreview();
    }
  },
);

watch(template, () => {
  if (props.open) renderPreview();
});

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
