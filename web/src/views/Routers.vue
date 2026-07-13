<template>
  <PageLayout title="Routers" subtitle="Manage your MikroTik connections">
    <template #actions>
      <button class="btn btn-ghost" :disabled="scanning" @click="scan">
        <MagnifyingGlassIcon class="size-3.5" :class="{ 'animate-pulse': scanning }" />
        {{ scanning ? 'Scanning…' : 'Scan network' }}
      </button>
      <button class="btn btn-primary" @click="showAdd = true">
        <PlusIcon class="size-3.5" />
        Add router
      </button>
    </template>

    <!-- Scan error -->
    <div v-if="scanError" class="flex items-center gap-2 p-3 rounded-lg text-xs border text-red bg-red/8 border-red/20">
      <ExclamationTriangleIcon class="size-4 shrink-0" />
      {{ scanError }}
    </div>

    <!-- Scanning spinner -->
    <div v-if="scanning" class="border border-dashed border-border rounded-xl py-16 text-center">
      <span class="inline-block size-5 border-2 border-border border-t-text-secondary rounded-full animate-spin mb-3" />
      <p class="text-xs text-text-muted">Scanning network… (3 s)</p>
    </div>

    <!-- Discovered devices -->
    <div v-else-if="discovered.length > 0">
      <p class="text-sm font-medium text-text-secondary mb-2">{{ discovered.length }} device{{ discovered.length > 1 ? 's' : '' }} found on network</p>
      <div class="space-y-2 mb-6">
        <div
          v-for="d in discovered"
          :key="d.mac || d.ip"
          class="flex items-center gap-3 p-3 border border-border rounded-xl"
        >
          <div class="size-9 shrink-0 rounded-lg flex items-center justify-center border border-border bg-base">
            <ServerIcon class="size-4 text-text-muted" />
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-text-primary truncate">{{ d.identity || d.board || 'MikroTik' }}</p>
            <p class="text-xs text-text-secondary font-mono">{{ d.ip && d.ip !== '0.0.0.0' ? d.ip : 'No IP yet' }}<span v-if="d.mac"> · {{ d.mac }}</span></p>
            <p v-if="d.board || d.version" class="text-xs text-text-muted mt-0.5">{{ [d.board, d.version].filter(Boolean).join(' · ') }}</p>
            <p v-if="!d.ip || d.ip === '0.0.0.0'" class="text-xs mt-0.5 text-amber">Factory reset — assign an IP via WinBox first</p>
          </div>
          <button
            class="btn btn-sm btn-ghost shrink-0"
            :class="isAlreadySaved(d.ip) ? 'cursor-default' : (!d.ip || d.ip === '0.0.0.0') ? 'opacity-40 cursor-not-allowed' : ''"
            :style="isAlreadySaved(d.ip) ? 'color: var(--color-green); border-color: rgba(76,195,138,0.3)' : ''"
            :disabled="isAlreadySaved(d.ip) || !d.ip || d.ip === '0.0.0.0'"
            @click="openAdd(d)"
          >
            <CheckIcon v-if="isAlreadySaved(d.ip)" class="size-3.5" />
            <PlusIcon v-else class="size-3.5" />
            {{ isAlreadySaved(d.ip) ? 'Saved' : 'Add' }}
          </button>
        </div>
      </div>
    </div>

    <!-- No devices after scan -->
    <div v-else-if="scanned && discovered.length === 0" class="border border-dashed border-border rounded-xl py-12 text-center mb-6">
      <ServerIcon class="size-6 text-text-muted mx-auto mb-2" />
      <p class="text-xs font-medium text-text-secondary">No MikroTik devices found</p>
      <p class="text-xs text-text-muted mt-0.5">Make sure the router is on the same subnet</p>
    </div>

    <!-- Router cards -->
    <div v-if="store.routers.length > 0" class="flex flex-col gap-3">
      <div
        v-for="r in store.routers"
        :key="r.id"
        class="flex items-center gap-4 px-4 py-3 border rounded-xl transition-colors"
        :class="store.activeId === r.id ? 'bg-surface border-border' : 'bg-transparent border-border/50'"
      >
        <!-- Icon + active dot -->
        <div class="relative shrink-0">
          <div class="size-9 rounded-lg flex items-center justify-center border border-border bg-base">
            <ServerIcon class="size-4 text-text-muted" />
          </div>
          <span
            v-if="store.activeId === r.id"
            class="absolute -top-0.5 -right-0.5 size-2.5 rounded-full border-2 bg-green border-surface"
          />
        </div>

        <!-- Info -->
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2">
            <span class="font-semibold text-text-primary">{{ r.name }}</span>
            <span v-if="store.activeId === r.id" class="text-xs px-1.5 py-0.5 rounded text-green bg-green/10">Active</span>
          </div>
          <p class="text-xs text-text-secondary font-mono mt-0.5">
            {{ r.host }}<span class="text-text-muted"> : {{ r.port }}</span> · {{ r.username }}<span v-if="r.useTls" class="ml-1 text-text-muted">TLS</span>
          </p>
        </div>

        <!-- Connection status -->
        <div class="shrink-0 w-28 text-right">
          <span v-if="testResults[r.id] === undefined" class="text-xs text-text-muted">—</span>
          <span v-else-if="testResults[r.id] === 'ok'" class="inline-flex items-center gap-1 text-xs px-2 py-0.5 rounded-full text-green bg-green/10">
            <CheckCircleIcon class="size-3" /> Reachable
          </span>
          <span v-else class="inline-flex items-center gap-1 text-xs px-2 py-0.5 rounded-full text-red bg-red/10">
            <XCircleIcon class="size-3" /> Unreachable
          </span>
        </div>

        <!-- Action buttons with labels -->
        <div class="flex items-center gap-2 shrink-0 border-l border-border pl-3 ml-1">
          <button
            class="btn btn-ghost"
            :disabled="testing[r.id]"
            @click="test(r.id)"
          >
            <ArrowPathIcon class="size-3.5" :class="{ 'animate-spin': testing[r.id] }" />
            {{ testing[r.id] ? 'Testing…' : 'Test' }}
          </button>

          <button
            class="btn btn-ghost"
            :disabled="store.activeId === r.id"
            :class="store.activeId === r.id ? 'opacity-40 cursor-default' : ''"
            @click="selectRouter(r.id)"
          >
            <BoltIcon class="size-3.5" />
            Use
          </button>

          <button class="btn btn-ghost" @click="openEdit(r)">
            <PencilSquareIcon class="size-3.5" />
            Edit
          </button>

          <button
            class="btn border-transparent text-text-muted hover:text-red hover:bg-red/10"
            @click="remove(r.id)"
          >
            <TrashIcon class="size-3.5" />
            Delete
          </button>
        </div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else-if="!scanning && !scanned" class="border border-dashed border-border rounded-xl py-16 text-center">
      <ServerIcon class="size-8 text-text-muted mx-auto mb-3" />
      <p class="text-sm font-medium text-text-secondary">No routers configured</p>
      <p class="text-xs text-text-muted mt-1">Click <strong class="text-text-secondary">Scan</strong> to find devices on your network, or add one manually</p>
    </div>

    <AddRouterDialog v-model:open="showAdd" :prefill="prefill" @added="onAdded" />
    <EditRouterDialog v-model:open="showEdit" :router="editTarget" @saved="store.load()" />
  </PageLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  PlusIcon, ServerIcon, ArrowPathIcon, BoltIcon, TrashIcon,
  CheckCircleIcon, XCircleIcon, MagnifyingGlassIcon,
  CheckIcon, ExclamationTriangleIcon, PencilSquareIcon,
} from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'
import { testRouter, discoverRouters, type DiscoveredDevice, type RouterProfile } from '@/api'
import AddRouterDialog from '@/components/AddRouterDialog.vue'
import EditRouterDialog from '@/components/EditRouterDialog.vue'
import PageLayout from '@/components/PageLayout.vue'

const store = useRoutersStore()
const router = useRouter()

const testing = ref<Record<string, boolean>>({})
const testResults = ref<Record<string, 'ok' | 'error'>>({})

const scanning = ref(false)
const scanned = ref(false)
const scanError = ref('')
const discovered = ref<DiscoveredDevice[]>([])
const addedIps = ref(new Set<string>())

const showAdd = ref(false)
const prefill = ref<Partial<DiscoveredDevice>>({})

const showEdit = ref(false)
const editTarget = ref<RouterProfile | null>(null)

function openEdit(r: RouterProfile) {
  editTarget.value = r
  showEdit.value = true
}

async function scan() {
  scanning.value = true
  scanned.value = false
  scanError.value = ''
  discovered.value = []
  try {
    discovered.value = await discoverRouters()
    scanned.value = true
  } catch (e: any) {
    scanError.value = e?.response?.data?.error ?? 'Scan failed — check network/firewall'
  } finally {
    scanning.value = false
  }
}

function isAlreadySaved(ip: string): boolean {
  return store.routers.some(r => r.host === ip) || addedIps.value.has(ip)
}

function openAdd(d: DiscoveredDevice) {
  prefill.value = d
  showAdd.value = true
}

function onAdded(ip: string) {
  addedIps.value = new Set([...addedIps.value, ip])
}

async function test(id: string) {
  testing.value[id] = true
  try {
    await testRouter(id)
    testResults.value[id] = 'ok'
  } catch {
    testResults.value[id] = 'error'
  } finally {
    testing.value[id] = false
  }
}

function selectRouter(id: string) {
  store.select(id)
  router.push('/dashboard')
}

async function remove(id: string) {
  if (!confirm('Remove this router?')) return
  await store.remove(id)
}
</script>
