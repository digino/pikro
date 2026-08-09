<template>
  <PageLayout title="Hotspot" subtitle="Setup Wizard">
    <div
      v-if="phase === 'loading'"
      class="flex flex-col items-center py-16 gap-3"
    >
      <span class="spinner" />
      <p class="text-sm text-text-secondary">Checking router configuration…</p>
    </div>

    <!-- ── Phase: blocked ─────────────────────────────────────────────────── -->
    <div v-else-if="phase === 'blocked'" class="space-y-4">
      <div
        class="flex text-sm items-start gap-3 p-4 border rounded-xl bg-amber/8 border-amber/20"
      >
        <ExclamationTriangleIcon class="size-5 shrink-0 mt-0.5 text-amber" />
        <div class="flex-1 min-w-0">
          <p class="font-medium text-amber">Hotspot already configured</p>
          <p class="text-text-secondary mt-1">
            A hotspot is active on this router. Reset it below to reconfigure,
            or go to Users to manage it.
          </p>
        </div>
      </div>

      <div
        class="border border-border rounded-xl p-4 text-sm space-y-2 bg-surface"
      >
        <p class="font-semibold text-text-secondary uppercase mb-3">
          Detected hotspot
        </p>
        <div class="grid grid-cols-2 gap-x-4 gap-y-2">
          <span class="text-text-secondary">Server name</span>
          <span class="font-mono text-text-primary">{{
            preflight?.hotspotName || "—"
          }}</span>
          <span class="text-text-secondary">Interface</span>
          <span class="font-mono text-text-primary">{{
            preflight?.hotspotOnIface || "—"
          }}</span>
          <span class="text-text-secondary">Profile</span>
          <span class="font-mono text-text-primary">{{
            preflight?.hotspotProfile || "—"
          }}</span>
          <span class="text-text-secondary">Address pool</span>
          <span class="font-mono text-text-primary">{{
            preflight?.hotspotAddressPool || "—"
          }}</span>
          <span class="text-text-secondary">DNS name</span>
          <span class="font-mono text-text-primary">{{
            preflight?.hotspotDnsName || "(none)"
          }}</span>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <RouterLink to="/hotspot/settings" class="btn btn-primary mr-auto">
          <Cog6ToothIcon class="size-3.5" />
          Configure branding &amp; vouchers
        </RouterLink>
        <button class="btn btn-ghost" @click="runPreflight">
          <ArrowPathIcon class="size-3.5" />
          Re-check
        </button>
        <button
          class="btn border-red/30 text-red hover:bg-red/8"
          :disabled="tearing"
          @click="runTeardown"
        >
          <TrashIcon class="size-3.5" />
          {{ tearing ? "Resetting…" : "Reset config" }}
        </button>
      </div>

      <div
        v-if="teardownResult"
        class="border border-border rounded-xl p-4 bg-surface"
      >
        <p
          class="text-xs font-semibold text-text-muted uppercase tracking-wide mb-2"
        >
          Reset result
        </p>
        <ul class="space-y-1.5">
          <li
            v-for="step in teardownResult.steps"
            :key="step.name"
            class="flex items-center gap-2 text-sm"
          >
            <CheckCircleIcon
              v-if="step.ok && !step.skipped"
              class="size-3.5 shrink-0 text-green"
            />
            <MinusCircleIcon
              v-else-if="step.skipped"
              class="size-3.5 text-text-muted shrink-0"
            />
            <XCircleIcon v-else class="size-3.5 shrink-0 text-red" />
            <span
              :class="step.skipped ? 'text-text-muted' : 'text-text-secondary'"
              >{{ step.name }}</span
            >
            <span v-if="step.error" class="font-mono text-red"
              >— {{ step.error }}</span
            >
          </li>
        </ul>
      </div>
    </div>

    <!-- ── Phase: form ────────────────────────────────────────────────────── -->
    <div v-else-if="phase === 'form'" class="space-y-5">
      <!-- Step indicator -->
      <StepperRoot
        :model-value="step"
        linear
        class="flex items-center gap-2"
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
                :class="
                  step === i + 1
                    ? 'text-text-primary font-medium'
                    : 'text-text-muted'
                "
              >
                {{ label }}
              </StepperTitle>
            </StepperTrigger>
            <StepperSeparator
              v-if="i < STEPS.length - 1"
              class="flex-1 h-px bg-border"
            />
          </StepperItem>
        </template>
      </StepperRoot>

      <div
        v-if="preflightError"
        class="flex items-start gap-2 p-3 border rounded-lg text-xs bg-red/8 border-red/20 text-red"
      >
        <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5" />
        <span
          >Could not read interfaces: {{ preflightError }}. Enter interface
          names manually.</span
        >
      </div>

      <form @submit.prevent="onFormSubmit" class="space-y-5">
        <!-- Step: Bridge (conditional) -->
        <div
          v-if="STEPS[step - 1] === 'Bridge'"
          class="border border-border rounded-xl p-5 space-y-4 bg-surface"
        >
          <div
            class="flex items-start gap-2 p-3 rounded-lg border text-xs bg-amber/8 border-amber/20 text-amber"
          >
            <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5" />
            <span
              >No bridge or wireless interface was found on this router. Create
              a bridge to use as the hotspot LAN.</span
            >
          </div>

          <label class="flex flex-col gap-1">
            <span class="text-sm font-medium">Bridge name</span>
            <input
              v-model="bridge.name"
              class="input"
              placeholder="bridge1"
              required
            />
          </label>

          <div class="flex flex-col gap-1.5">
            <span class="text-sm font-medium">Ports to add</span>
            <p class="text-xs text-text-muted">
              Select the physical ports clients will connect through. Leave your
              uplink port (usually ether1) unchecked.
            </p>
            <div class="grid grid-cols-2 gap-1.5 mt-1">
              <label
                v-for="i in etherInterfaces"
                :key="i.name"
                class="flex items-center gap-2 cursor-pointer select-none text-sm"
              >
                <span
                  class="relative inline-flex items-center justify-center size-4 rounded border shrink-0 transition-colors"
                  :style="
                    bridge.ports.includes(i.name)
                      ? 'background: var(--color-accent); border-color: var(--color-accent)'
                      : 'background: transparent; border-color: var(--color-border)'
                  "
                >
                  <svg
                    v-if="bridge.ports.includes(i.name)"
                    viewBox="0 0 10 8"
                    fill="none"
                    class="size-2.5"
                  >
                    <path
                      d="M1 4l2.5 2.5L9 1"
                      stroke="#09090b"
                      stroke-width="1.5"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                </span>
                <input
                  type="checkbox"
                  :value="i.name"
                  v-model="bridge.ports"
                  class="sr-only"
                />
                <span class="text-text-secondary"
                  >{{ i.name }}{{ !i.running ? " — down" : "" }}</span
                >
              </label>
            </div>
          </div>
        </div>

        <!-- Step: Network -->
        <div
          v-else-if="STEPS[step - 1] === 'Network'"
          class="border border-border rounded-xl p-5 space-y-4 bg-surface"
        >
          <h2 class="font-semibold text-text-primary">Network configuration</h2>

          <label class="flex flex-col gap-1">
            <span class="font-medium text-red"
              >LAN / Hotspot interface <span>*</span></span
            >
            <AppSelect
              v-if="lanOptions.length"
              v-model="form.lanIface"
              :options="lanOptions"
              placeholder="Select interface…"
            />
            <input
              v-else
              v-model="form.lanIface"
              class="input"
              placeholder="e.g. bridge1 or ether2"
            />
            <p class="text-sm text-text-muted">
              The interface clients connect to (usually a bridge or wlan)
            </p>
          </label>

          <label class="flex flex-col gap-1">
            <span class="font-medium text-red"
              >WAN interface <span>*</span></span
            >
            <AppSelect
              v-if="wanOptions.length"
              v-model="form.wanIface"
              :options="wanOptions"
              placeholder="Select interface…"
            />
            <input
              v-else
              v-model="form.wanIface"
              class="input"
              placeholder="e.g. ether1"
            />
            <p class="text-sm text-text-muted">
              The interface connected to the internet (usually ether1)
            </p>
          </label>

          <label class="flex flex-col gap-1">
            <span class="font-medium">Hotspot name</span>
            <div class="join">
              <input
                v-model="form.hotspotName"
                class="input flex-1 min-w-0"
                placeholder="myhotspot"
                pattern="[a-zA-Z0-9\-]+"
              />
              <AppSelect
                v-model="form.hotspotExtension"
                :options="extensionOptions"
                trigger-class="w-auto shrink-0"
              />
            </div>
            <p class="text-sm text-text-muted">
              Clients are redirected to
              <span
                >{{ form.hotspotName || "myhotspot"
                }}{{ form.hotspotExtension }}</span
              >
              — leave blank to use the router IP instead
            </p>
          </label>
        </div>

        <!-- Step: Subnet -->
        <div
          v-else-if="STEPS[step - 1] === 'Subnet'"
          class="border border-border rounded-xl p-5 space-y-4 bg-surface"
        >
          <div
            v-if="dhcpChecking"
            class="flex items-center gap-2 text-sm text-text-muted"
          >
            <span
              class="inline-block size-3.5 border-2 rounded-full animate-spin border-border border-t-text-secondary"
            />
            Checking for an existing DHCP server on {{ form.lanIface }}…
          </div>
          <div
            v-else-if="dhcpCheck?.exists"
            class="flex items-start gap-2 p-3 rounded-lg border text-sm bg-amber/8 border-amber/20 text-amber"
          >
            <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5" />
            <span
              >{{ form.lanIface }} already has a DHCP server for
              <span class="font-mono">{{ dhcpCheck.subnet }}</span> — reusing
              it instead of creating a new one, so the subnet is fixed to
              this network.</span
            >
          </div>

          <label class="flex flex-col gap-1">
            <span class="font-medium text-red"
              >Subnet (CIDR) <span>*</span></span
            >
            <input
              v-model="form.subnet"
              class="input font-mono"
              placeholder="192.168.88.0/24"
              :disabled="dhcpCheck?.exists"
            />
            <p v-if="subnetError" class="text-sm text-red">{{ subnetError }}</p>
          </label>

          <div
            v-if="derived"
            class="border border-border rounded-lg p-3 text-sm bg-base"
          >
            <div class="grid grid-cols-2 gap-x-4 gap-y-1.5">
              <span class="text-text-secondary">Router IP (gateway)</span>
              <span class="font-mono text-text-primary font-medium">{{
                derived.gateway
              }}</span>
              <span class="text-text-secondary">DHCP pool</span>
              <span class="font-mono text-text-primary font-medium"
                >{{ derived.poolStart }} – {{ derived.poolEnd }}</span
              >
              <span class="text-text-secondary">Hotspot profile name</span>
              <span class="font-mono text-text-primary font-medium"
                >pikro-profile</span
              >
            </div>
          </div>
        </div>

        <!-- Step: Review -->
        <template v-else-if="STEPS[step - 1] === 'Review'">
          <div
            class="flex items-start gap-2 text-sm p-3 border rounded-lg bg-amber/8 border-amber/20 text-amber"
          >
            <ExclamationTriangleIcon class="size-5 shrink-0 mt-0.5" />
            <span
              >Ensure no hotspot is already configured on this router. Running
              setup on an existing configuration may cause conflicts.</span
            >
          </div>

          <div class="border border-border rounded-xl p-5 space-y-4 bg-surface">
            <h2 class="text-sm font-semibold text-text-primary">
              What will be configured
            </h2>
            <ul class="space-y-1 text-sm text-text-secondary">
              <li v-if="needsBridge">
                • Create bridge {{ bridge.name || "(bridge)"
                }}{{
                  bridge.ports.length
                    ? ` with ports ${bridge.ports.join(", ")}`
                    : ""
                }}
              </li>
              <li>
                • Assign {{ derived?.gateway ?? "…" }} to
                {{ form.lanIface || "(LAN interface)" }}
              </li>
              <li>
                • Create IP pool <span class="font-mono">hotspot-pool</span>
              </li>
              <li>• Create DHCP server + network</li>
              <li>• Create hotspot server profile with Pikro login page</li>
              <li>
                • Enable hotspot on {{ form.lanIface || "(LAN interface)" }}
              </li>
              <li>
                • Add NAT masquerade on {{ form.wanIface || "(WAN interface)" }}
              </li>
              <li>• Enable DNS remote requests</li>
            </ul>
          </div>
        </template>

        <!-- Step: Cleanup -->
        <template v-else-if="STEPS[step - 1] === 'Cleanup'">
          <div class="border border-border rounded-xl p-5 space-y-4 bg-surface">
            <h2 class="text-sm font-semibold text-text-primary">
              Protect your vouchers
            </h2>
            <p class="text-sm text-text-secondary">
              RouterOS doesn't remove expired vouchers on its own — Pikro installs
              a cleanup scheduler on the router so expired and quota-exhausted
              vouchers get removed automatically, on the schedule you pick below.
            </p>
            <label class="flex flex-col gap-1">
              <span class="text-sm font-medium text-text-secondary">Run every</span>
              <AppSelect
                v-model="form.cleanupInterval"
                :options="CLEANUP_INTERVAL_OPTIONS"
              />
            </label>
          </div>
        </template>

        <div class="flex items-center gap-3">
          <button type="button" class="btn btn-ghost" @click="onBack">
            {{ step === 1 ? "Re-check router" : "Back" }}
          </button>
          <button
            v-if="step < STEPS.length"
            type="button"
            class="btn btn-primary"
            :disabled="!canAdvance"
            @click="step++"
          >
            Next
          </button>
          <button
            v-else
            type="submit"
            class="btn btn-primary"
            :disabled="!canSubmit"
          >
            Run Setup
          </button>
          <!-- <button
            type="button"
            class="btn ml-auto border-red/30 text-red hover:bg-red/8"
            :disabled="tearing"
            @click="runTeardown"
          >
            <TrashIcon class="size-3.5" />
            {{ tearing ? "Resetting…" : "Reset config" }}
          </button> -->
        </div>
      </form>

      <div
        v-if="teardownResult"
        class="border border-border rounded-xl p-4 bg-surface"
      >
        <p
          class="text-xs font-semibold text-text-muted uppercase tracking-wide mb-2"
        >
          Reset result
        </p>
        <ul class="space-y-1.5">
          <li
            v-for="step in teardownResult.steps"
            :key="step.name"
            class="flex items-center gap-2 text-sm"
          >
            <CheckCircleIcon
              v-if="step.ok && !step.skipped"
              class="size-3.5 shrink-0 text-green"
            />
            <MinusCircleIcon
              v-else-if="step.skipped"
              class="size-3.5 text-text-muted shrink-0"
            />
            <XCircleIcon v-else class="size-3.5 shrink-0 text-red" />
            <span
              :class="step.skipped ? 'text-text-muted' : 'text-text-secondary'"
              >{{ step.name }}</span
            >
            <span v-if="step.error" class="font-mono text-red"
              >— {{ step.error }}</span
            >
          </li>
        </ul>
      </div>
    </div>

    <!-- ── Phase: running ─────────────────────────────────────────────────── -->
    <div
      v-else-if="phase === 'running'"
      class="flex flex-col items-center py-16 gap-3"
    >
      <span class="spinner" />
      <p class="text-sm font-medium text-text-primary">Configuring hotspot…</p>
      <p class="text-xs text-text-muted">
        This takes a few seconds. Please wait.
      </p>
    </div>

    <!-- ── Phase: result ──────────────────────────────────────────────────── -->
    <div v-else-if="phase === 'result' && result" class="space-y-4">
      <div
        class="flex items-center gap-2 p-3 rounded-lg text-sm font-medium border"
        :class="
          result.success
            ? 'bg-green/8 border-green/20 text-green'
            : 'bg-red/8 border-red/20 text-red'
        "
      >
        <CheckCircleIcon v-if="result.success" class="size-4 shrink-0" />
        <ExclamationTriangleIcon v-else class="size-4 shrink-0" />
        {{
          result.success
            ? "Hotspot configured successfully!"
            : "Setup completed with errors"
        }}
      </div>

      <div class="border border-border rounded-xl p-4 bg-surface">
        <ul class="space-y-2">
          <li
            v-for="step in result.steps"
            :key="step.name"
            class="flex items-start gap-2 text-sm"
          >
            <CheckCircleIcon
              v-if="step.ok && !step.skipped"
              class="size-3.5 shrink-0 mt-0.5 text-green"
            />
            <MinusCircleIcon
              v-else-if="step.skipped"
              class="size-3.5 text-text-muted shrink-0 mt-0.5"
            />
            <XCircleIcon v-else class="size-3.5 shrink-0 mt-0.5 text-red" />
            <div class="flex-1 min-w-0">
              <span
                :class="
                  step.ok
                    ? 'text-text-secondary'
                    : step.skipped
                      ? 'text-text-muted'
                      : 'font-medium text-red'
                "
              >
                {{ step.name }}
                <span v-if="step.skipped" class="text-text-muted font-normal">
                  — already configured</span
                >
              </span>
              <p v-if="step.error" class="font-mono mt-0.5 break-all text-red">
                {{ step.error }}
              </p>
            </div>
          </li>
        </ul>
      </div>

      <div
        v-if="result.success"
        class="border border-border rounded-xl p-4 space-y-3 bg-surface"
      >
        <p
          class="text-xs font-semibold text-text-muted uppercase tracking-wide"
        >
          What's next
        </p>
        <ol class="space-y-2.5 text-sm text-text-secondary">
          <li class="flex items-start gap-2">
            <span
              class="size-5 shrink-0 rounded-full flex items-center justify-center text-xs font-bold mt-0.5 bg-muted text-text-primary"
              >1</span
            >
            <div>
              <span
                >Create your first hotspot user with a time or data limit.</span
              >
              <span class="block mt-0.5 text-amber">
                RouterOS added a default
                <span class="font-mono">admin</span> user with no password —
                delete or secure it from the Users page.
              </span>
            </div>
          </li>
          <li class="flex items-start gap-2">
            <span
              class="size-5 shrink-0 rounded-full flex items-center justify-center text-xs font-bold mt-0.5 bg-muted text-text-primary"
              >2</span
            >
            <span
              >Connect a client device to the router's WiFi or LAN — it will be
              redirected to
              <span class="font-mono">{{
                form.hotspotName
                  ? form.hotspotName + form.hotspotExtension
                  : "the login page"
              }}</span
              >.</span
            >
          </li>
          <li class="flex items-start gap-2">
            <span
              class="size-5 shrink-0 rounded-full flex items-center justify-center text-xs font-bold mt-0.5 bg-muted text-text-primary"
              >3</span
            >
            <span
              >Share the credentials — or create multiple users in bulk for
              voucher-style access.</span
            >
          </li>
          <li class="flex items-start gap-2">
            <span
              class="size-5 shrink-0 rounded-full flex items-center justify-center text-xs font-bold mt-0.5 bg-muted text-text-primary"
              >4</span
            >
            <span>
              Configure your business name, currency, voucher layout, and login
              page in
              <RouterLink
                to="/hotspot/settings"
                class="underline text-text-primary hover:text-text-secondary transition-colors"
                >Settings</RouterLink
              >.
            </span>
          </li>
        </ol>
      </div>

      <div
        v-if="!result.success"
        class="text-sm border border-border rounded-lg p-3 space-y-1 bg-surface"
      >
        <p class="font-medium text-text-primary">Some steps failed</p>
        <p class="text-text-muted">
          Steps marked as "already configured" are fine — they were skipped
          because the resource exists. Fix any red steps manually in RouterOS,
          then run setup again.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <RouterLink
          v-if="result.success"
          to="/hotspot/users"
          class="btn btn-primary"
        >
          <UsersIcon class="size-4" />
          Add users
        </RouterLink>
        <button class="btn btn-ghost btn-sm border-transparent" @click="reset">
          <ArrowPathIcon class="size-3.5" />
          Run again
        </button>
      </div>
    </div>

    <!-- ── Teardown confirmation dialog ─────────────────────────────────── -->
    <AppDialog
      :open="showTeardownConfirm"
      title="Reset hotspot config?"
      @update:open="showTeardownConfirm = $event"
    >
      <div class="space-y-4">
        <div
          class="border border-border rounded-lg p-3 text-sm space-y-1.5 bg-base"
        >
          <div class="flex justify-between font-medium">
            <span class="text-text-secondary">Router</span>
            <span class="text-text-primary">{{
              store.routers.find((r) => r.id === store.activeId)?.name
            }}</span>
          </div>
          <div class="flex justify-between font-medium">
            <span class="text-text-secondary">Host</span>
            <span class="font-mono text-text-primary">{{
              store.routers.find((r) => r.id === store.activeId)?.host
            }}</span>
          </div>
          <div class="flex justify-between font-medium">
            <span class="text-text-secondary">Hotspot interface</span>
            <span class="font-mono text-text-primary">{{
              preflight?.hotspotOnIface || form.lanIface || "—"
            }}</span>
          </div>
        </div>

        <div>
          <p class="font-medium text-text-primary text-sm mb-1.5">
            This will remove:
          </p>
          <ul class="text-sm text-text-secondary space-y-1">
            <li class="flex items-center gap-1.5 shrink-0 text-red">
              <TrashIcon class="size-4" /> Hotspot server
            </li>
            <li class="flex items-center gap-1.5 shrink-0 text-red">
              <TrashIcon class="size-4" /> Hotspot server profile
            </li>
            <li class="flex items-center gap-1.5 shrink-0 text-red">
              <TrashIcon class="size-4" /> NAT masquerade rule
            </li>
            <li class="flex items-center gap-1.5 shrink-0 text-red">
              <TrashIcon class="size-4" /> Walled garden rules
            </li>
          </ul>
          <p class="text-sm text-text-secondary font-medium mt-2">
            IP address, DHCP server, IP pool and DNS settings are kept.
          </p>
        </div>

        <div class="flex justify-end gap-2 pt-1">
          <button class="btn btn-ghost" @click="showTeardownConfirm = false">
            Cancel
          </button>
          <button class="btn btn-danger" @click="confirmTeardown">
            Reset config
          </button>
        </div>
      </div>
    </AppDialog>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import {
  StepperRoot,
  StepperItem,
  StepperTrigger,
  StepperIndicator,
  StepperTitle,
  StepperSeparator,
} from "reka-ui";
import {
  ExclamationTriangleIcon,
  ArrowPathIcon,
  CheckCircleIcon,
  XCircleIcon,
  MinusCircleIcon,
  UsersIcon,
  TrashIcon,
  Cog6ToothIcon,
} from "@heroicons/vue/24/outline";
import { RouterLink } from "vue-router";
import {
  hotspotPreflight,
  setupHotspot,
  teardownHotspot,
  checkExistingDHCP,
  type PreflightResult,
  type SetupResult,
  type TeardownResult,
  type InterfaceInfo,
  type DHCPCheckResult,
} from "@/api";
import { useRoutersStore } from "@/stores/routers";
import { friendlyError } from "@/utils/errors";
import { HOTSPOT_EXTENSION_OPTIONS } from "@/utils/hotspotExtensions";
import { CLEANUP_INTERVAL_OPTIONS } from "@/utils/cleanupIntervals";
import AppDialog from "@/components/AppDialog.vue";
import PageLayout from "@/components/PageLayout.vue";
import AppSelect, { type SelectOption } from "@/components/AppSelect.vue";

const store = useRoutersStore();

type Phase = "loading" | "form" | "blocked" | "running" | "result";
const phase = ref<Phase>("loading");
const preflight = ref<PreflightResult | null>(null);
const preflightError = ref("");
const result = ref<SetupResult | null>(null);
const interfaces = ref<InterfaceInfo[]>([]);
const form = ref({
  lanIface: "",
  wanIface: "",
  subnet: "192.168.88.0/24",
  hotspotName: "",
  hotspotExtension: ".spot",
  cleanupInterval: "7d",
});

const extensionOptions = HOTSPOT_EXTENSION_OPTIONS;
const tearing = ref(false);
const teardownResult = ref<TeardownResult | null>(null);
const showTeardownConfirm = ref(false);

const bridge = ref({ name: "bridge1", ports: [] as string[] });

const bridgeAndWlan = computed(() =>
  interfaces.value.filter(
    (i) => i.type === "bridge" || i.type === "wlan" || i.type === "vlan",
  ),
);

const etherInterfaces = computed(() =>
  interfaces.value.filter((i) => i.type === "ether"),
);

const needsBridge = ref(false);

const STEPS = ref<string[]>(["Network", "Subnet", "Review"]);
const step = ref(1);

const dhcpCheck = ref<DHCPCheckResult | null>(null);
const dhcpChecking = ref(false);
const dhcpCheckedIface = ref("");

async function checkDHCP() {
  if (!store.activeId || !form.value.lanIface) return;
  if (dhcpCheckedIface.value === form.value.lanIface) return;
  dhcpChecking.value = true;
  try {
    dhcpCheck.value = await checkExistingDHCP(store.activeId, form.value.lanIface);
    dhcpCheckedIface.value = form.value.lanIface;
    if (dhcpCheck.value.exists && dhcpCheck.value.subnet) {
      form.value.subnet = dhcpCheck.value.subnet;
    }
  } catch {
    dhcpCheck.value = null;
  } finally {
    dhcpChecking.value = false;
  }
}

watch(step, (s) => {
  if (STEPS.value[s - 1] === "Subnet") checkDHCP();
});

// A different LAN interface invalidates any previous check.
watch(
  () => form.value.lanIface,
  () => {
    dhcpCheckedIface.value = "";
    dhcpCheck.value = null;
  },
);

function ifaceLabel(i: InterfaceInfo): string {
  return `${i.name} [${i.type}]${!i.running ? " — down" : ""}${i.comment ? ` (${i.comment})` : ""}`;
}

const lanOptions = computed<SelectOption[]>(() => {
  if (!interfaces.value.length) return [];
  if (needsBridge.value && bridge.value.name) {
    return [
      { value: bridge.value.name, label: `${bridge.value.name} (new bridge)` },
    ];
  }
  const preferred = bridgeAndWlan.value.map((i) => ({
    value: i.name,
    label: ifaceLabel(i),
  }));
  const rest = interfaces.value
    .filter((i) => !bridgeAndWlan.value.includes(i))
    .map((i) => ({ value: i.name, label: ifaceLabel(i) }));
  return [...preferred, ...rest];
});

const wanOptions = computed<SelectOption[]>(() =>
  interfaces.value.map((i) => ({ value: i.name, label: ifaceLabel(i) })),
);

const subnetError = computed(() => {
  const raw = form.value.subnet.trim();
  if (!raw) return "";
  const m = raw.match(
    /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})\/(\d{1,2})$/,
  );
  if (!m) return "Enter a valid CIDR, e.g. 192.168.88.0/24";
  const octets = [m[1], m[2], m[3], m[4]].map(Number);
  if (octets.some((o) => o > 255)) return "Octets must be 0-255";
  const prefix = parseInt(m[5]);
  if (prefix < 16 || prefix > 30) {
    return "Prefix must be between /16 and /30 (too small or too large for a hotspot LAN)";
  }
  if (octets[3] !== 0) {
    return "Use the network address (last octet 0), e.g. 192.168.88.0/24";
  }
  if (octets[0] === 127 || (octets[0] === 169 && octets[1] === 254)) {
    return "That range is reserved and can't be used for a LAN";
  }
  return "";
});

const derived = computed(() => {
  if (subnetError.value) return null;
  const raw = form.value.subnet.trim();
  const m = raw.match(/^(\d+\.\d+\.\d+)\.(\d+)\/(\d+)$/);
  if (!m) return null;
  const base3 = m[1];
  return {
    gateway: `${base3}.1`,
    poolStart: `${base3}.2`,
    poolEnd: `${base3}.254`,
  };
});

const canSubmit = computed(
  () =>
    form.value.lanIface.trim() !== "" &&
    form.value.wanIface.trim() !== "" &&
    derived.value !== null &&
    (!needsBridge.value || bridge.value.name.trim() !== ""),
);

const canAdvance = computed(() => {
  const current = STEPS.value[step.value - 1];
  if (current === "Bridge") return bridge.value.name.trim() !== "";
  if (current === "Network")
    return (
      form.value.lanIface.trim() !== "" && form.value.wanIface.trim() !== ""
    );
  if (current === "Subnet") return derived.value !== null;
  return true;
});

async function runPreflight() {
  if (!store.activeId) return;
  phase.value = "loading";
  preflightError.value = "";
  try {
    const pf = await hotspotPreflight(store.activeId);
    preflight.value = pf;
    interfaces.value = pf.interfaces;
    if (pf.hotspotExists) {
      phase.value = "blocked";
      return;
    }
    const bridgeIface = interfaces.value.find((i) => i.type === "bridge");
    const wlan = interfaces.value.find((i) => i.type === "wlan");
    const ether = interfaces.value.find((i) => i.type === "ether");
    needsBridge.value = bridgeAndWlan.value.length === 0;
    STEPS.value = needsBridge.value
      ? ["Bridge", "Network", "Subnet", "Review", "Cleanup"]
      : ["Network", "Subnet", "Review", "Cleanup"];
    if (bridgeIface || wlan) {
      form.value.lanIface = bridgeIface?.name ?? wlan?.name ?? "";
    } else {
      form.value.lanIface = bridge.value.name;
    }
    form.value.wanIface =
      interfaces.value.find(
        (i) => i.type === "ether" && i.name !== form.value.lanIface,
      )?.name ??
      ether?.name ??
      "";
    step.value = 1;
    phase.value = "form";
  } catch (e: any) {
    preflightError.value = friendlyError(e, "Connection failed");
    phase.value = "form";
  }
}

function onFormSubmit() {
  if (step.value < STEPS.value.length) {
    step.value++;
    return;
  }
  runSetup();
}

function onBack() {
  if (step.value === 1) {
    runPreflight();
  } else {
    step.value--;
  }
}

async function runSetup() {
  if (!store.activeId || !canSubmit.value) return;
  phase.value = "running";
  try {
    result.value = await setupHotspot(store.activeId, {
      lanIface: form.value.lanIface,
      wanIface: form.value.wanIface,
      subnet: form.value.subnet,
      hotspotName: form.value.hotspotName,
      extension: form.value.hotspotExtension,
      cleanupInterval: form.value.cleanupInterval,
      ...(needsBridge.value
        ? { newBridgeName: bridge.value.name, bridgePorts: bridge.value.ports }
        : {}),
    });
  } catch (e: any) {
    result.value = {
      success: false,
      steps: [
        {
          name: "Setup request",
          ok: false,
          error: friendlyError(e, "Request failed"),
        },
      ],
    };
  }
  phase.value = "result";
}

function reset() {
  result.value = null;
  runPreflight();
}

function runTeardown() {
  showTeardownConfirm.value = true;
}

async function confirmTeardown() {
  if (!store.activeId) return;
  showTeardownConfirm.value = false;
  tearing.value = true;
  teardownResult.value = null;
  try {
    teardownResult.value = await teardownHotspot(store.activeId);
    await runPreflight();
  } catch (e: any) {
    teardownResult.value = {
      steps: [
        {
          name: "Teardown",
          ok: false,
          error: friendlyError(e, "Reset failed"),
        },
      ],
    };
  } finally {
    tearing.value = false;
  }
}

watch(
  () => store.activeId,
  (id) => {
    if (id) runPreflight();
  },
  { immediate: true },
);

// Keep the LAN interface in sync with the bridge name while it's being created.
watch(
  () => bridge.value.name,
  (name) => {
    if (needsBridge.value) form.value.lanIface = name;
  },
);
</script>
