<template>
  <!-- Step indicator -->
  <div class="flex items-center gap-2 mb-5">
    <template v-for="(label, i) in STEPS" :key="i">
      <div class="flex items-center gap-1.5">
        <span
          class="size-6 rounded-full flex items-center justify-center text-xs font-bold transition-colors"
          :class="
            step === i + 1
              ? 'bg-accent text-base'
              : step > i + 1
                ? 'bg-green/20 text-green'
                : 'bg-muted text-text-muted'
          "
          >{{ step > i + 1 ? "✓" : i + 1 }}</span
        >
        <span
          class="text-sm hidden sm:inline"
          :class="
            step === i + 1 ? 'text-text-primary font-medium' : 'text-text-muted'
          "
          >{{ label }}</span
        >
      </div>
      <div v-if="i < STEPS.length - 1" class="flex-1 h-px bg-border" />
    </template>
  </div>

  <!-- Step 1: Quantity & format -->
  <div v-if="step === 1" class="space-y-4">
    <p class="font-semibold">How many vouchers do you need?</p>
    <input
      v-model.number="form.count"
      type="number"
      min="1"
      max="10000"
      class="input text-2xl font-mono text-center py-3"
      placeholder="10"
      autofocus
    />
    <p v-if="form.count > 100" class="text-xs text-amber">
      Large batch — may take {{ Math.round(form.count * 0.3) }}s+
    </p>

    <div class="space-y-2">
      <p class="font-medium">Username length</p>
      <div class="grid grid-cols-4 gap-2">
        <button
          v-for="n in [4, 5, 6, 8]"
          :key="n"
          type="button"
          class="border rounded-lg py-2 text-sm font-mono transition-colors"
          :class="
            form.nameLength === n
              ? 'border-accent bg-accent/10 text-text-primary font-semibold'
              : 'border-border text-text-secondary hover:border-muted'
          "
          @click="form.nameLength = n"
        >
          {{ n }} chars
        </button>
      </div>
    </div>

    <div class="space-y-2">
      <p class="font-medium">Characters</p>
      <div class="flex gap-4">
        <label
          class="flex items-center gap-1.5 text-sm text-text-secondary cursor-pointer"
        >
          <CheckboxRoot
            :checked="form.charsLetters"
            class="size-4 rounded border border-border bg-base flex items-center justify-center transition-colors data-[state=checked]:bg-accent data-[state=checked]:border-accent"
            @update:checked="form.charsLetters = $event"
          >
            <CheckboxIndicator>
              <CheckIcon class="size-3 text-base" />
            </CheckboxIndicator>
          </CheckboxRoot>
          Letters (a–z)
        </label>
        <label
          class="flex items-center gap-1.5 text-sm text-text-secondary cursor-pointer"
        >
          <CheckboxRoot
            :checked="form.charsDigits"
            class="size-4 rounded border border-border bg-base flex items-center justify-center transition-colors data-[state=checked]:bg-accent data-[state=checked]:border-accent"
            @update:checked="form.charsDigits = $event"
          >
            <CheckboxIndicator>
              <CheckIcon class="size-3 text-base" />
            </CheckboxIndicator>
          </CheckboxRoot>
          Digits (0–9)
        </label>
      </div>
      <p class="text-sm text-text-muted font-mono">
        Preview: <span class="text-text-primary">{{ namePreview }}</span>
      </p>
    </div>
  </div>

  <!-- Step 2: Password -->
  <div v-else-if="step === 2" class="space-y-4">
    <p class="font-semibold">How should passwords be set?</p>
    <div class="grid gap-2">
      <button
        v-for="opt in PASSWORD_OPTS"
        :key="opt.value"
        type="button"
        class="flex items-start gap-3 border rounded-lg p-3 text-left transition-colors"
        :class="
          form.passwordMode === opt.value
            ? 'border-accent bg-accent/10'
            : 'border-border hover:border-muted'
        "
        @click="form.passwordMode = opt.value"
      >
        <span
          class="mt-0.5 size-5 rounded-full border-2 shrink-0 flex items-center justify-center transition-colors"
          :class="
            form.passwordMode === opt.value ? 'border-accent' : 'border-border'
          "
        >
          <span
            v-if="form.passwordMode === opt.value"
            class="size-2 rounded-full bg-accent"
          />
        </span>
        <div>
          <p class="text-sm font-semibold text-text-primary">{{ opt.label }}</p>
          <p class="text-xs text-text-secondary">{{ opt.desc }}</p>
        </div>
      </button>
    </div>
    <input
      v-if="form.passwordMode === 'fixed'"
      v-model="form.fixedPassword"
      class="input"
      placeholder="Enter fixed password"
    />
  </div>

  <!-- Step 3: Profile & limits -->
  <div v-else-if="step === 3" class="space-y-4">
    <p class="font-semibold">Profile and limits</p>

    <div class="flex flex-col gap-1">
      <span class="text-sm font-medium text-text-secondary">Profile</span>
      <SelectRoot
        :model-value="form.profile || undefined"
        @update:model-value="form.profile = $event ?? ''"
      >
        <SelectTrigger
          class="flex items-center gap-1.5 w-full p-2 text-sm bg-base border border-border rounded-lg text-text-primary cursor-pointer transition-colors hover:border-muted focus-visible:outline-2 focus-visible:outline-accent focus-visible:outline-offset-1"
        >
          <SelectValue
            placeholder="default"
            class="flex-1 truncate text-left"
          />
          <ChevronDownIcon class="size-4 text-text-secondary shrink-0" />
        </SelectTrigger>
        <SelectPortal>
          <SelectContent
            class="z-50 min-w-(--reka-select-trigger-width) bg-surface border border-border rounded-lg shadow-xl overflow-hidden"
            position="popper"
            :side-offset="3"
          >
            <SelectViewport class="p-1">
              <SelectItem
                value=""
                class="flex items-center justify-between p-2 text-sm rounded-md cursor-pointer text-text-secondary transition-colors hover:bg-muted hover:text-text-primary data-highlighted:bg-muted data-highlighted:text-text-primary data-[state=checked]:text-text-primary data-[state=checked]:font-medium"
              >
                <SelectItemText>default</SelectItemText>
                <SelectItemIndicator
                  ><CheckCircleIcon class="size-4 text-green"
                /></SelectItemIndicator>
              </SelectItem>
              <SelectItem
                v-for="p in profiles"
                :key="p['.id']"
                :value="p.name"
                class="flex items-center justify-between p-2 text-sm rounded-md cursor-pointer text-text-secondary transition-colors hover:bg-muted hover:text-text-primary data-highlighted:bg-muted data-highlighted:text-text-primary data-[state=checked]:text-text-primary data-[state=checked]:font-medium"
              >
                <SelectItemText>{{ p.name }}</SelectItemText>
                <SelectItemIndicator
                  ><CheckCircleIcon class="size-4 text-green"
                /></SelectItemIndicator>
              </SelectItem>
            </SelectViewport>
          </SelectContent>
        </SelectPortal>
      </SelectRoot>
    </div>

    <div class="flex flex-col gap-1">
      <span class="text-sm font-medium text-text-secondary">
        Time limit
        <span class="text-text-muted font-normal">(overrides profile)</span>
      </span>
      <input
        v-model="form.limitUptimeRaw"
        class="input font-mono"
        placeholder="e.g. 1h, 1d, 1w — blank to use profile"
        @blur="normalizeUptime"
      />
      <p v-if="form.limitUptimeRaw && !uptimePreview" class="text-xs text-red">
        Invalid — use: 30m, 2h, 1d, 1w or combinations like 1d12h
      </p>
      <p
        v-else-if="uptimePreview"
        class="text-xs"
        :class="uptimeWarning ? 'text-amber' : 'text-text-muted'"
      >
        <span v-if="uptimeWarning"
          >⚠ Exceeds profile validity — user may never hit this limit.</span
        >
        <span v-else class="font-mono"
          >Sends to router: {{ uptimePreview }}</span
        >
      </p>
    </div>

    <div class="flex flex-col gap-1">
      <span class="text-sm font-medium text-text-secondary">
        Data limit
        <span class="text-text-muted font-normal">(overrides profile)</span>
      </span>
      <div
        class="flex w-full overflow-hidden rounded-lg border border-border focus-within:outline-2 focus-within:outline-accent focus-within:outline-offset-1"
      >
        <input
          v-model.number="form.limitBytesTotalValue"
          type="number"
          min="0"
          class="input-bare flex-1 min-w-0"
          placeholder="0 = use profile"
        />
        <SelectRoot
          :model-value="form.limitBytesTotalUnit"
          @update:model-value="form.limitBytesTotalUnit = $event"
        >
          <SelectTrigger
            class="flex items-center gap-1 px-2 border-l border-border text-sm text-text-secondary bg-base focus:outline-none w-20"
          >
            <SelectValue />
            <ChevronDownIcon class="size-3.5 shrink-0" />
          </SelectTrigger>
          <SelectPortal>
            <SelectContent
              class="z-50 bg-surface border border-border rounded-lg shadow-xl overflow-hidden"
              position="popper"
              :side-offset="3"
            >
              <SelectViewport class="p-1">
                <SelectItem
                  v-for="u in ['M', 'G']"
                  :key="u"
                  :value="u"
                  class="flex items-center justify-between p-2 text-sm rounded-md cursor-pointer text-text-secondary hover:bg-muted hover:text-text-primary data-highlighted:bg-muted data-highlighted:text-text-primary data-[state=checked]:text-text-primary data-[state=checked]:font-medium"
                >
                  <SelectItemText>{{ u === "M" ? "MB" : "GB" }}</SelectItemText>
                </SelectItem>
              </SelectViewport>
            </SelectContent>
          </SelectPortal>
        </SelectRoot>
      </div>
    </div>
  </div>

  <!-- Step 4: Review -->
  <div v-else-if="step === 4" class="space-y-4">
    <p class="font-semibold">Ready to generate</p>
    <div class="rounded-xl border border-border divide-y divide-border text-sm">
      <div class="flex justify-between px-4 py-2.5">
        <span class="text-text-muted">Quantity</span>
        <span class="font-semibold text-text-primary font-mono">{{
          form.count
        }}</span>
      </div>
      <div class="flex justify-between px-4 py-2.5">
        <span class="text-text-muted">Username</span>
        <span class="font-mono text-text-primary"
          >{{ form.nameLength }} chars · {{ namePreview }}</span
        >
      </div>
      <div class="flex justify-between px-4 py-2.5">
        <span class="text-text-muted">Password</span>
        <span class="text-text-primary">{{
          PASSWORD_OPTS.find((o) => o.value === form.passwordMode)?.label
        }}</span>
      </div>
      <div class="flex justify-between px-4 py-2.5">
        <span class="text-text-muted">Profile</span>
        <span class="text-text-primary">{{ form.profile || "default" }}</span>
      </div>
      <div v-if="uptimePreview" class="flex justify-between px-4 py-2.5">
        <span class="text-text-muted">Time limit</span>
        <span class="font-mono text-text-primary">{{ uptimePreview }}</span>
      </div>
      <div
        v-if="form.limitBytesTotalValue"
        class="flex justify-between px-4 py-2.5"
      >
        <span class="text-text-muted">Data limit</span>
        <span class="font-mono text-text-primary"
          >{{ form.limitBytesTotalValue }}
          {{ form.limitBytesTotalUnit === "G" ? "GB" : "MB" }}</span
        >
      </div>
    </div>
    <label class="flex flex-col gap-1">
      <span class="text-sm font-medium text-text-secondary">
        Comment <span class="text-text-muted font-normal">(optional)</span>
      </span>
      <input
        v-model="form.comment"
        class="input"
        placeholder="e.g. July batch, event name…"
      />
    </label>
    <p v-if="error" class="text-xs text-red">{{ error }}</p>
  </div>

  <!-- Navigation -->
  <div
    class="flex items-center justify-between pt-4 mt-1 border-border"
  >
    <button
      type="button"
      class="btn btn-ghost"
      @click="step > 1 ? step-- : emit('cancel')"
    >
      {{ step > 1 ? "Back" : "Cancel" }}
    </button>
    <button
      v-if="step < 4"
      type="button"
      class="btn btn-primary"
      :disabled="!step1Valid"
      @click="step++"
    >
      Next
    </button>
    <button
      v-else
      type="button"
      class="btn btn-primary"
      :disabled="!!form.limitUptimeRaw && !uptimePreview"
      @click="emit('submit', { ...form, uptimePreview })"
    >
      Generate {{ form.count }} vouchers
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import {
  SelectRoot,
  SelectTrigger,
  SelectValue,
  SelectPortal,
  SelectContent,
  SelectViewport,
  SelectItem,
  SelectItemText,
  SelectItemIndicator,
  CheckboxRoot,
  CheckboxIndicator,
} from "reka-ui";
import { ChevronDownIcon, CheckCircleIcon, CheckIcon } from "@heroicons/vue/24/outline";

const props = defineProps<{
  profiles: Record<string, string>[];
  profileMetas: Record<string, { validity?: string; price?: string }>;
  error?: string;
}>();

const emit = defineEmits<{
  cancel: [];
  submit: [config: BatchConfig];
}>();

export interface BatchConfig {
  count: number;
  nameLength: number;
  charsLetters: boolean;
  charsDigits: boolean;
  passwordMode: "same" | "random" | "fixed";
  fixedPassword: string;
  profile: string;
  limitUptimeRaw: string;
  uptimePreview: string;
  limitBytesTotalValue: number;
  limitBytesTotalUnit: "M" | "G";
  comment: string;
}

const STEPS = ["Quantity", "Password", "Profile", "Review"];

const PASSWORD_OPTS = [
  {
    value: "same" as const,
    label: "Same as username",
    desc: "Easiest to print and hand out",
  },
  {
    value: "random" as const,
    label: "Random (separate)",
    desc: "More secure, each voucher unique",
  },
  {
    value: "fixed" as const,
    label: "Fixed password",
    desc: "One shared password for all",
  },
];

const step = ref(1);

const form = ref<Omit<BatchConfig, "uptimePreview">>({
  count: 10,
  nameLength: 6,
  charsLetters: true,
  charsDigits: true,
  passwordMode: "same",
  fixedPassword: "",
  profile: "",
  limitUptimeRaw: "",
  limitBytesTotalValue: 0,
  limitBytesTotalUnit: "M",
  comment: "",
});

const charset = computed(() => {
  let s = "";
  if (form.value.charsLetters) s += "abcdefghijklmnopqrstuvwxyz";
  if (form.value.charsDigits) s += "0123456789";
  return s;
});

const namePreview = computed(() => {
  if (!charset.value) return "—";
  return Array.from(
    { length: form.value.nameLength },
    () => charset.value[Math.floor(Math.random() * charset.value.length)],
  ).join("");
});

const step1Valid = computed(
  () =>
    step.value !== 1 ||
    (!!form.value.count && form.value.count <= 500 && !!charset.value),
);

function parseShorthand(
  s: string,
): { w: number; d: number; h: number; m: number } | null {
  if (!s.trim()) return null;
  const re = /^(?:(\d+)w)?(?:(\d+)d)?(?:(\d+)h)?(?:(\d+)m)?$/i;
  const match = s.trim().match(re);
  if (!match || !match[0]) return null;
  const [, w, d, h, min] = match;
  if (!w && !d && !h && !min) return null;
  return {
    w: parseInt(w || "0"),
    d: parseInt(d || "0"),
    h: parseInt(h || "0"),
    m: parseInt(min || "0"),
  };
}

function shorthandToSeconds(s: string): number {
  const p = parseShorthand(s);
  if (!p) return 0;
  return p.w * 7 * 86400 + p.d * 86400 + p.h * 3600 + p.m * 60;
}

const uptimePreview = computed(() => {
  const p = parseShorthand(form.value.limitUptimeRaw);
  if (!p) return "";
  const totalDays = p.w * 7 + p.d;
  if (!totalDays && !p.h && !p.m) return "";
  return [
    totalDays ? `${totalDays}d` : "",
    p.h ? `${p.h}h` : "",
    p.m ? `${p.m}m` : "",
  ]
    .filter(Boolean)
    .join("");
});

const uptimeWarning = computed(() => {
  if (!uptimePreview.value || !form.value.profile) return false;
  const meta = props.profileMetas[form.value.profile];
  if (!meta?.validity) return false;
  return (
    shorthandToSeconds(form.value.limitUptimeRaw) >
    shorthandToSeconds(meta.validity)
  );
});

function normalizeUptime() {
  const p = parseShorthand(form.value.limitUptimeRaw);
  if (!p) return;
  form.value.limitUptimeRaw = [
    p.w ? `${p.w}w` : "",
    p.d ? `${p.d}d` : "",
    p.h ? `${p.h}h` : "",
    p.m ? `${p.m}m` : "",
  ]
    .filter(Boolean)
    .join("");
}

function reset() {
  step.value = 1;
  form.value = {
    count: 10,
    nameLength: 6,
    charsLetters: true,
    charsDigits: true,
    passwordMode: "same",
    fixedPassword: "",
    profile: "",
    limitUptimeRaw: "",
    limitBytesTotalValue: 0,
    limitBytesTotalUnit: "M",
    comment: "",
  };
}

defineExpose({ reset });
</script>
