<template>
  <AppDialog :open="open" title="Edit Router" @update:open="onDialogUpdate">
    <!-- Step indicator -->
    <StepperRoot
      :model-value="step"
      linear
      class="flex items-center gap-2 mb-5"
      @update:model-value="step = $event ?? step"
    >
      <template v-for="(label, i) in STEPS" :key="i">
        <StepperItem
          :step="i + 1"
          :completed="step > i + 1"
          class="flex items-center gap-1.5"
          :class="i < STEPS.length - 1 ? 'flex-1' : ''"
        >
          <StepperTrigger class="flex items-center gap-1.5" as="div">
            <StepperIndicator
              class="size-6 bg-primary rounded-full flex items-center justify-center text-xs font-bold transition-colors shrink-0 data-[state=active]:bg-accent data-[state=active]:text-base data-[state=completed]:bg-green/20 data-[state=completed]:text-green data-[state=inactive]:bg-muted data-[state=inactive]:text-text-muted"
            >
              {{ step > i + 1 ? "✓" : i + 1 }}
            </StepperIndicator>
            <StepperTitle
              class="text-sm hidden sm:inline"
              :class="step === i + 1 ? 'text-text-primary font-medium' : 'text-text-muted'"
            >
              {{ label }}
            </StepperTitle>
          </StepperTrigger>
          <StepperSeparator v-if="i < STEPS.length - 1" class="flex-1 h-px bg-border" />
        </StepperItem>
      </template>
    </StepperRoot>

    <form @submit.prevent="onFormSubmit" class="space-y-3">
      <!-- Step 1: Connection -->
      <template v-if="step === 1">
        <div class="grid grid-cols-2 gap-3">
          <label class="col-span-2 flex flex-col gap-1">
            <span class="text-xs font-medium text-text-secondary">Name</span>
            <input v-model="form.name" class="input" placeholder="Home router" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-xs font-medium text-text-secondary">Host / IP</span>
            <input v-model="form.host" class="input" placeholder="192.168.88.1" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-xs font-medium text-text-secondary">Port</span>
            <input v-model.number="form.port" type="number" class="input" placeholder="8728" />
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-xs font-medium text-text-secondary">Username</span>
            <input v-model="form.username" class="input" placeholder="admin" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="text-xs font-medium text-text-secondary">Password</span>
            <input v-model="form.password" type="password" class="input" placeholder="Leave blank to keep current" />
          </label>
        </div>

        <label class="flex items-center gap-2 cursor-pointer select-none">
          <span
            class="relative inline-flex items-center justify-center size-4 rounded border shrink-0 transition-colors"
            :style="form.useTls
              ? 'background: var(--color-accent); border-color: var(--color-accent)'
              : 'background: transparent; border-color: var(--color-border)'"
          >
            <svg v-if="form.useTls" viewBox="0 0 10 8" fill="none" class="size-2.5">
              <path d="M1 4l2.5 2.5L9 1" stroke="#09090b" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" />
            </svg>
          </span>
          <input v-model="form.useTls" type="checkbox" class="sr-only" />
          <span class="text-sm text-text-secondary">Use TLS (port 8729)</span>
        </label>
      </template>

      <!-- Step 2: Hotspot info -->
      <RouterHotspotInfoFields
        v-else
        v-model:hotspot-name="form.hotspotName"
        v-model:dns-name="form.dnsName"
        v-model:currency="form.currency"
      />

      <p v-if="submitError" class="text-sm text-red">{{ submitError }}</p>

      <div class="flex justify-end gap-2 pt-1">
        <button
          type="button"
          class="btn btn-ghost"
          @click="step === 1 ? $emit('update:open', false) : step--"
        >
          {{ step === 1 ? "Cancel" : "Back" }}
        </button>
        <button
          v-if="step === 1"
          type="button"
          class="btn btn-primary"
          :disabled="!form.host || !form.username || !form.name"
          @click="step = 2"
        >
          Next
        </button>
        <button v-else type="submit" class="btn btn-primary" :disabled="loading">
          <span v-if="loading" class="spinner spinner--sm" />
          Save changes
        </button>
      </div>
    </form>
  </AppDialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  StepperRoot, StepperItem, StepperTrigger, StepperIndicator, StepperTitle, StepperSeparator,
} from 'reka-ui'
import { updateRouter, type RouterProfile } from '@/api'
import AppDialog from '@/components/AppDialog.vue'
import RouterHotspotInfoFields from '@/components/RouterHotspotInfoFields.vue'

const props = defineProps<{ open: boolean; router: RouterProfile | null }>()
const emit = defineEmits<{ 'update:open': [value: boolean]; saved: [] }>()

const STEPS = ['Connection', 'Hotspot info']
const step = ref(1)

const loading = ref(false)
const submitError = ref('')

function emptyForm() {
  return {
    name: '', host: '', port: 8728, username: '', password: '', useTls: false,
    hotspotName: '', dnsName: '', currency: '',
  }
}
const form = ref(emptyForm())

watch(() => props.router, (r) => {
  if (!r) return
  form.value = {
    name: r.name, host: r.host, port: r.port, username: r.username, password: '', useTls: r.useTls,
    hotspotName: r.hotspotSettings?.hotspotName ?? '',
    dnsName: r.hotspotSettings?.dnsName ?? '',
    currency: r.hotspotSettings?.currency ?? '',
  }
}, { immediate: true })

function onDialogUpdate(isOpen: boolean) {
  if (!isOpen) step.value = 1
  emit('update:open', isOpen)
}

function onFormSubmit() {
  if (step.value === 1) {
    step.value = 2
    return
  }
  submit()
}

async function submit() {
  if (!props.router) return
  loading.value = true
  submitError.value = ''
  try {
    const f = form.value
    await updateRouter(props.router.id, {
      name: f.name,
      host: f.host,
      port: f.port,
      username: f.username,
      ...(f.password ? { password: f.password } : {}),
      useTls: f.useTls,
      hotspotSettings: {
        hotspotName: f.hotspotName,
        dnsName: f.dnsName,
        currency: f.currency,
      },
    })
    emit('saved')
    emit('update:open', false)
  } catch (e: any) {
    submitError.value = e?.response?.data?.error ?? e?.message ?? 'Failed to save'
  } finally {
    loading.value = false
  }
}
</script>
