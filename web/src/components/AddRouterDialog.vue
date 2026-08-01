<template>
  <AppDialog :open="open" title="Add Router" @update:open="onDialogUpdate">
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
            <span class="font-medium text-text-secondary">Name</span>
            <input v-model="form.name" class="input" placeholder="Home router" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="font-medium">Host / IP</span>
            <input v-model="form.host" class="input" placeholder="192.168.88.1" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="font-medium">Port</span>
            <input v-model.number="form.port" type="number" class="input" placeholder="8728" />
          </label>
          <label class="flex flex-col gap-1">
            <span class="font-medium">Username</span>
            <input v-model="form.username" class="input" placeholder="admin" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="font-medium">Password</span>
            <input v-model="form.password" type="password" class="input" />
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
          Add router
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
import { useRoutersStore } from '@/stores/routers'
import type { DiscoveredDevice } from '@/api'
import AppDialog from '@/components/AppDialog.vue'
import RouterHotspotInfoFields from '@/components/RouterHotspotInfoFields.vue'

const props = defineProps<{ open: boolean; prefill?: Partial<DiscoveredDevice> }>()
const emit = defineEmits<{ 'update:open': [value: boolean]; added: [ip: string] }>()

const store = useRoutersStore()
const loading = ref(false)
const submitError = ref('')

const STEPS = ['Connection', 'Hotspot info']
const step = ref(1)

function emptyForm() {
  return {
    name: '', host: '', port: 8728, username: 'admin', password: '', useTls: false,
    hotspotName: '', dnsName: '', currency: '',
  }
}
const form = ref(emptyForm())

// Pre-fill form when opened from the scan dialog
watch(() => props.prefill, (d) => {
  if (!d) return
  form.value.name = d.identity || d.board || 'MikroTik'
  form.value.host = d.ip ?? ''
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
  loading.value = true
  submitError.value = ''
  try {
    const f = form.value
    await store.add({
      name: f.name, host: f.host, port: f.port, username: f.username,
      password: f.password, useTls: f.useTls,
      hotspotSettings: {
        hotspotName: f.hotspotName,
        dnsName: f.dnsName,
        currency: f.currency,
        voucher: { showValidity: true, showPrice: true },
      },
    })
    emit('added', f.host)
    emit('update:open', false)
    step.value = 1
    form.value = emptyForm()
  } catch (e: any) {
    submitError.value = e?.response?.data?.error ?? 'Failed to add router'
  } finally {
    loading.value = false
  }
}
</script>
