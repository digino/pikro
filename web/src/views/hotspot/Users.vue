<template>
  <PageLayout title="Hotspot" subtitle="Users">
    <template #actions>
      <template v-if="tab === 'users' && selected.size > 0">
        <span class="text-sm text-text-secondary font-semibold">{{ selected.size }} selected</span>
        <button
          class="text-sm text-text-secondary hover:text-text-primary transition-colors cursor-pointer underline"
          @click="selected = new Set()"
        >
          Clear
        </button>
        <button
          class="btn btn-danger"
          :disabled="bulkDeleting"
          @click="removeSelected"
        >
          <span
            v-if="bulkDeleting"
            class="size-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin"
          />
          <TrashIcon v-else class="size-3.5" />
          Delete {{ selected.size }}
        </button>
      </template>
      <RouterLink to="/hotspot/vouchers" class="btn btn-ghost">
        <TicketIcon class="size-4" />
        Generate vouchers
      </RouterLink>
      <button
        v-if="tab === 'users'"
        class="btn btn-ghost"
        :disabled="filteredUsers.length === 0"
        @click="openPrint"
      >
        <PrinterIcon class="size-4" />
        {{ selected.size > 0 ? `Print ${selected.size} selected` : "Print" }}
      </button>
      <button class="btn btn-primary" @click="openAdd">
        <PlusIcon class="size-4" />
        New user
      </button>
    </template>

    <!-- Tabs -->
    <div class="flex items-center justify-between border-b border-border -mt-2">
      <div class="flex items-center gap-1">
        <button
          v-for="t in tabs"
          :key="t.key"
          class="px-3 py-2 text-sm font-semibold border-b-2 transition-colors"
          :class="
            tab === t.key
              ? 'border-text-primary text-text-primary'
              : 'border-transparent text-text-muted hover:text-text-secondary'
          "
          @click="switchTab(t.key)"
        >
          {{ t.label }}
          <span
            v-if="t.key === 'users' || active.length > 0"
            class="ml-1 text-text-secondary"
          >
            ({{ t.key === "users" ? filteredUsers.length : active.length }})
          </span>
        </button>
      </div>

      <div class="flex items-center gap-2 pb-2">
        <span class="text-sm text-text-secondary font-semibold"
          >Auto-cleanup</span
        >
        <SwitchRoot
          :model-value="cleanupInstalled === true"
          :disabled="cleanupToggling"
          class="relative inline-flex h-5 w-9 shrink-0 items-center rounded-full transition-colors data-[state=checked]:bg-green data-[state=unchecked]:bg-border disabled:opacity-50"
          @update:model-value="toggleCleanup"
        >
          <SwitchThumb
            class="pointer-events-none block size-4 rounded-full bg-white shadow transform transition-transform translate-x-0.5 data-[state=checked]:translate-x-4"
          />
        </SwitchRoot>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-10">
      <span class="spinner" />
    </div>

    <div
      v-else-if="error"
      class="flex items-center gap-2 p-4 border rounded-xl text-sm bg-red/8 border-red/20 text-red"
    >
      <ExclamationTriangleIcon class="size-4 shrink-0" />
      {{ error }}
      <button
        class="ml-auto text-xs underline"
        @click="tab === 'users' ? loadUsers() : loadActive()"
      >
        Retry
      </button>
    </div>

    <!-- Users tab -->
    <div v-else-if="tab === 'users'" class="space-y-2">
      <!-- Search & filters -->
      <div class="flex items-center gap-2.5 flex-wrap">
        <div class="relative flex-1 max-w-xs bg-wh">
          <input
            v-model="searchQuery"
            class="input"
            placeholder="Search username or comment…"
          />
        </div>
        <SelectRoot
          :model-value="filterProfile || undefined"
          @update:model-value="filterProfile = $event ?? ''"
        >
          <SelectTrigger
            class="input-select flex items-center gap-1.5 cursor-pointer"
          >
            <SelectValue placeholder="All profiles" class="flex-1 text-left" />
            <ChevronDownIcon class="size-3.5 text-text-secondary shrink-0" />
          </SelectTrigger>
          <SelectPortal>
            <SelectContent
              class="z-50 min-w-(--reka-select-trigger-width) bg-surface border border-border rounded-lg shadow-xl overflow-hidden"
              position="popper"
              :side-offset="3"
            >
              <SelectViewport class="p-1">
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
        <SelectRoot
          :model-value="filterStatus || undefined"
          @update:model-value="filterStatus = $event ?? ''"
        >
          <SelectTrigger
            class="input-select flex items-center gap-1.5 cursor-pointer"
          >
            <SelectValue placeholder="All statuses" class="flex-1 text-left" />
            <ChevronDownIcon class="size-3.5 text-text-secondary shrink-0" />
          </SelectTrigger>
          <SelectPortal>
            <SelectContent
              class="z-50 min-w-(--reka-select-trigger-width) bg-surface border border-border rounded-lg shadow-xl overflow-hidden"
              position="popper"
              :side-offset="3"
            >
              <SelectViewport class="p-1">
                <SelectItem
                  v-for="opt in STATUS_OPTS"
                  :key="opt.value"
                  :value="opt.value"
                  class="flex items-center justify-between p-2 text-sm rounded-md cursor-pointer text-text-secondary transition-colors hover:bg-muted hover:text-text-primary data-highlighted:bg-muted data-highlighted:text-text-primary data-[state=checked]:text-text-primary data-[state=checked]:font-medium"
                >
                  <SelectItemText>{{ opt.label }}</SelectItemText>
                  <SelectItemIndicator
                    ><CheckCircleIcon class="size-4 text-green"
                  /></SelectItemIndicator>
                </SelectItem>
              </SelectViewport>
            </SelectContent>
          </SelectPortal>
        </SelectRoot>

        <button
          v-if="searchQuery || filterProfile || filterStatus"
          class="btn"
          @click="
            searchQuery = '';
            filterProfile = '';
            filterStatus = '';
          "
        >
          Clear
        </button>
      </div>

      <div class="border border-border rounded-xl overflow-hidden">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-border bg-surface">
              <th class="px-4 py-3 w-8">
                <CheckboxRoot
                  :model-value="someSelected ? 'indeterminate' : allSelected"
                  class="size-4 rounded border border-border bg-base flex items-center justify-center transition-colors data-[state=checked]:bg-text-primary data-[state=checked]:border-text-primary data-[state=indeterminate]:bg-text-primary data-[state=indeterminate]:border-text-primary hover:border-muted"
                  @update:model-value="toggleSelectAll"
                >
                  <CheckboxIndicator>
                    <CheckIcon class="size-3 text-base" />
                  </CheckboxIndicator>
                </CheckboxRoot>
              </th>
              <th
                class="text-left px-4 py-3 text-sm font-semibold text-text-primary"
              >
                Username
              </th>
              <th class="text-left px-4 py-3 font-semibold text-text-primary">
                Profile
              </th>
              <th class="text-left px-4 py-3 font-semibold text-text-primary">
                Time limit
              </th>
              <th class="text-left px-4 py-3 font-semibold text-text-primary">
                Data limit
              </th>
              <th class="text-left px-4 py-3 font-semibold text-text-primary">
                Comment
              </th>
              <th class="text-left px-4 py-3 font-semibold text-text-primary">
                Expires
              </th>
              <th class="text-left px-4 py-3 font-semibold text-text-primary">
                Status
              </th>
              <th class="px-4 py-3"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredUsers.length === 0">
              <td colspan="9" class="text-center text-text-muted text-sm py-10">
                {{
                  searchQuery || filterProfile
                    ? "No users match your search"
                    : "No users found"
                }}
              </td>
            </tr>
            <tr
              v-for="u in pagedUsers"
              :key="u['.id']"
              class="border-b border-border last:border-0 transition-colors"
              :class="selected.has(u['.id']) ? 'bg-surface' : ''"
            >
              <td class="px-4 py-3">
                <CheckboxRoot
                  :model-value="selected.has(u['.id'])"
                  class="size-4 rounded border border-border bg-base flex items-center justify-center transition-colors data-[state=checked]:bg-text-primary data-[state=checked]:border-text-primary hover:border-muted"
                  @update:model-value="toggleSelect(u['.id'])"
                >
                  <CheckboxIndicator>
                    <CheckIcon class="size-3 text-base" />
                  </CheckboxIndicator>
                </CheckboxRoot>
              </td>
              <td
                class="px-4 py-3 font-mono font-semibold text-text-primary text-sm"
              >
                {{ u.name }}
              </td>
              <td class="px-4 py-3">
                <span
                  class="text-xs text-text-secondary font-medium border border-border px-2 py-0.5 rounded bg-white"
                >
                  {{ u.profile || "default" }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-text-secondary">
                {{ u["limit-uptime"] || "—" }}
              </td>
              <td class="px-4 py-3 text-sm text-text-secondary">
                {{ formatBytes(u["limit-bytes-total"]) }}
              </td>
              <td class="px-4 py-3 text-sm text-text-muted font-mono">
                {{ u.comment || "—" }}
              </td>
              <td class="px-4 py-3 text-xs" :class="expiryClass(u)">
                {{ expiryLabel(u) }}
              </td>
              <td class="px-4 py-3">
                <StatusBadge
                  variant="pill"
                  :color="statusColor(u)"
                  :label="statusLabel(u)"
                />
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex items-center justify-end gap-1">
                  <button
                    class="btn btn-ghost btn-sm"
                    @click="toggleDisabled(u)"
                  >
                    <component
                      :is="
                        u.disabled === 'true' ? CheckCircleIcon : NoSymbolIcon
                      "
                      class="size-3.5"
                    />
                    {{ u.disabled === "true" ? "Enable" : "Disable" }}
                  </button>
                  <button
                    class="btn btn-sm btn-ghost hover:text-red hover:bg-red/10"
                    @click="removeUser(u['.id'])"
                  >
                    <TrashIcon class="size-3.5" />
                    Delete
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
        <!-- Pagination footer -->
        <div
          v-if="filteredUsers.length > 0"
          class="flex items-center justify-between px-4 py-2.5 border-t border-border bg-surface"
        >
          <div class="flex items-center gap-1">
            <span class="text-sm mr-2">Per page</span>
            <button
              v-for="n in PAGE_SIZES"
              :key="n"
              class="btn btn-sm btn-ghost"
              :class="
                usersPageSize === n
                  ? 'bg-muted text-text-primary'
                  : 'text-text-muted hover:text-text-secondary'
              "
              @click="setUsersPageSize(n)"
            >
              {{ n }}
            </button>
          </div>
          <PaginationRoot
            v-slot="{ page }"
            :page="usersPage"
            :items-per-page="usersPageSize"
            :total="filteredUsers.length"
            @update:page="usersPage = $event"
          >
            <div class="flex items-center gap-1 text-sm text-text-muted">
              <span
                >{{ (page - 1) * usersPageSize + 1 }}–{{
                  Math.min(page * usersPageSize, filteredUsers.length)
                }}
                of {{ filteredUsers.length }}</span
              >
              <PaginationPrev
                class="btn btn-sm btn-primary disabled:opacity-30"
              >
                <ChevronLeftIcon class="size-4" />
              </PaginationPrev>
              <PaginationNext
                class="btn btn-sm btn-primary disabled:opacity-30"
              >
                <ChevronRightIcon class="size-4" />
              </PaginationNext>
            </div>
          </PaginationRoot>
        </div>
      </div>
    </div>

    <!-- Active sessions tab -->
    <div v-else class="border border-border rounded-xl overflow-hidden">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b border-border bg-surface">
            <th
              class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
            >
              User
            </th>
            <th
              class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
            >
              IP
            </th>
            <th
              class="text-left px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
            >
              MAC
            </th>
            <th
              class="text-right px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
            >
              Uptime
            </th>
            <th
              class="text-right px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
            >
              Down
            </th>
            <th
              class="text-right px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
            >
              Up
            </th>
            <th
              class="text-right px-4 py-3 text-xs font-semibold text-text-primary uppercase tracking-wide"
            >
              Left
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="active.length === 0">
            <td colspan="7" class="text-center text-text-muted text-sm py-10">
              No active sessions
            </td>
          </tr>
          <tr
            v-for="s in pagedActive"
            :key="s['.id']"
            class="border-b border-border last:border-0 transition-colors hover:bg-surface"
          >
            <td
              class="px-4 py-3 font-mono font-semibold text-text-primary text-xs"
            >
              {{ s.user }}
            </td>
            <td class="px-4 py-3 font-mono text-xs text-text-secondary">
              {{ s.address }}
            </td>
            <td class="px-4 py-3 font-mono text-xs text-text-muted">
              {{ s["mac-address"] || "—" }}
            </td>
            <td
              class="px-4 py-3 text-right font-mono text-xs text-text-secondary"
            >
              {{ s.uptime || "—" }}
            </td>
            <td
              class="px-4 py-3 text-right font-mono text-xs text-text-primary"
            >
              {{ formatBytes(s["bytes-in"]) }}
            </td>
            <td
              class="px-4 py-3 text-right font-mono text-xs text-text-secondary"
            >
              {{ formatBytes(s["bytes-out"]) }}
            </td>
            <td
              class="px-4 py-3 text-right font-mono text-xs text-text-secondary"
            >
              {{ s["session-time-left"] || "—" }}
            </td>
          </tr>
        </tbody>
      </table>
      <!-- Active sessions pagination footer -->
      <div
        v-if="active.length > 0"
        class="flex items-center justify-between px-4 py-2.5 border-t border-border bg-surface"
      >
        <div class="flex items-center gap-1">
          <span class="text-xs text-text-muted mr-1">Per page</span>
          <button
            v-for="n in PAGE_SIZES"
            :key="n"
            class="px-2 py-0.5 text-xs rounded transition-colors"
            :class="
              activePageSize === n
                ? 'bg-muted text-text-primary'
                : 'text-text-muted hover:text-text-secondary'
            "
            @click="setActivePageSize(n)"
          >
            {{ n }}
          </button>
        </div>
        <PaginationRoot
          v-slot="{ page }"
          :page="activePage"
          :items-per-page="activePageSize"
          :total="active.length"
          @update:page="activePage = $event"
        >
          <div class="flex items-center gap-1 text-xs text-text-muted">
            <span
              >{{ (page - 1) * activePageSize + 1 }}–{{
                Math.min(page * activePageSize, active.length)
              }}
              of {{ active.length }}</span
            >
            <PaginationPrev
              class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
            >
              <ChevronLeftIcon class="size-3.5" />
            </PaginationPrev>
            <PaginationNext
              class="p-1 rounded hover:bg-surface disabled:opacity-30 transition-colors"
            >
              <ChevronRightIcon class="size-3.5" />
            </PaginationNext>
          </div>
        </PaginationRoot>
      </div>
    </div>

    <!-- Add single user dialog -->
    <AppDialog
      :open="showAdd"
      title="New Hotspot User"
      @update:open="showAdd = $event"
    >
      <form @submit.prevent="submitAdd" class="space-y-4">
        <div class="grid grid-cols-2 gap-3">
          <label class="flex flex-col gap-1">
            <span class="font-medium text-red">Username <span>*</span></span>
            <input v-model="form.name" class="input" required />
          </label>
          <label class="flex flex-col gap-1">
            <span class="font-medium text-text-secondary">Password</span>
            <input
              v-model="form.password"
              class="input"
              placeholder="leave blank for none"
            />
          </label>
        </div>

        <label class="flex flex-col gap-1">
          <span class="font-medium text-text-secondary">Profile</span>
          <select v-model="form.profile" class="input">
            <option value="">default</option>
            <option v-for="p in profiles" :key="p['.id']" :value="p.name">
              {{ p.name }}
            </option>
          </select>
        </label>

        <div class="border-t border-border pt-3 space-y-3">
          <p
            class="text-sm font-semibold text-text-muted uppercase tracking-wide"
          >
            Limits (override profile)
          </p>

          <div class="flex flex-col gap-1">
            <span class="font-medium text-text-secondary">Time limit</span>
            <input
              v-model="form.limitUptimeRaw"
              class="input font-mono"
              placeholder="e.g. 1h, 1d, 1w, 1d12h — blank for unlimited"
              @blur="normalizeFormUptime"
            />
            <p
              v-if="form.limitUptimeRaw && !formUptimePreview"
              class="text-sm text-red"
            >
              Invalid format — use: 30m, 2h, 1d, 1w or combinations like 1d12h
            </p>
            <p
              v-else-if="formUptimePreview"
              class="text-sm"
              :class="formUptimeWarning ? 'text-amber' : 'text-text-muted'"
            >
              <span v-if="formUptimeWarning"
                >⚠ Time limit exceeds the profile's validity — user may never
                hit this limit.</span
              >
              <span v-else class="font-mono"
                >Sends to router: <span>{{ formUptimePreview }}</span></span
              >
            </p>
          </div>

          <div class="flex flex-col gap-1">
            <span class="font-medium text-text-secondary">Data limit</span>
            <div
              class="flex w-full overflow-hidden rounded-lg border border-border focus-within:outline-2 focus-within:outline-accent focus-within:outline-offset-1"
            >
              <input
                v-model.number="form.limitBytesTotalValue"
                type="number"
                min="0"
                class="input-bare flex-1 min-w-0"
                placeholder="0"
              />
              <select
                v-model="form.limitBytesTotalUnit"
                class="input-bare border-l border-border shrink-0 w-16 text-xs"
              >
                <option value="M">MB</option>
                <option value="G">GB</option>
              </select>
            </div>
          </div>
        </div>

        <label class="flex flex-col gap-1">
          <span class="font-medium text-text-secondary">Comment</span>
          <input
            v-model="form.comment"
            class="input"
            placeholder="optional note"
          />
        </label>

        <p v-if="addError" class="text-sm text-red">{{ addError }}</p>

        <div class="flex justify-end gap-2 pt-1">
          <button type="button" class="btn btn-ghost" @click="showAdd = false">
            Cancel
          </button>
          <button
            type="submit"
            class="btn btn-primary"
            :disabled="adding || (!!form.limitUptimeRaw && !formUptimePreview)"
          >
            <span
              v-if="adding"
              class="size-4 border-2 border-black/20 border-t-black rounded-full animate-spin"
            />
            Create
          </button>
        </div>
      </form>
    </AppDialog>

    <PrintTemplateDialog
      :open="showPrintDialog"
      :entries="printEntries"
      :business-name="hotspotSettings.hotspotName ?? ''"
      :currency="hotspotSettings.currency ?? ''"
      :profile-metas="profileMetas"
      :login-url="loginUrl"
      :login-url-supports-credentials="true"
      @update:open="showPrintDialog = $event"
    />
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from "vue";
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
  PaginationRoot,
  PaginationPrev,
  PaginationNext,
  SwitchRoot,
  SwitchThumb,
} from "reka-ui";
import {
  PlusIcon,
  TrashIcon,
  ExclamationTriangleIcon,
  TicketIcon,
  PrinterIcon,
  CheckIcon,
  NoSymbolIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
  ChevronDownIcon,
  CheckCircleIcon,
} from "@heroicons/vue/24/outline";
import { RouterLink, useRoute } from "vue-router";
import { useRoutersStore } from "@/stores/routers";
import {
  listHotspotUsers,
  listHotspotActive,
  listHotspotProfiles,
  createHotspotUser,
  toggleHotspotUser,
  deleteHotspotUser,
  getProfileMetas,
  getHotspotSettings,
  getCleanupScheduler,
  putCleanupScheduler,
  type ProfileMeta,
  type HotspotSettings,
} from "@/api";
import { friendlyError } from "@/utils/errors";
import AppDialog from "@/components/AppDialog.vue";
import PageLayout from "@/components/PageLayout.vue";
import StatusBadge from "@/components/StatusBadge.vue";
import PrintTemplateDialog from "@/components/PrintTemplateDialog.vue";
import { type VoucherEntry } from "@/utils/vouchers";
import { normalizeDnsName } from "@/utils/dnsName";

const store = useRoutersStore();
const route = useRoute();

const STATUS_OPTS = [
  { value: "active", label: "Active" },
  { value: "expires-soon", label: "Expires soon" },
  { value: "waiting", label: "Waiting" },
  { value: "limit-reached", label: "Limit reached" },
  { value: "expired", label: "Expired" },
  { value: "disabled", label: "Disabled" },
] as const;

const tab = ref<"users" | "active">(route.query.tab === "active" ? "active" : "users");
const tabs = [
  { key: "users" as const, label: "All Users" },
  { key: "active" as const, label: "Active Sessions" },
];
const users = ref<Record<string, string>[]>([]);
const active = ref<Record<string, string>[]>([]);

const searchQuery = ref("");
const filterProfile = ref("");
const filterStatus = ref("");

const filteredUsers = computed(() => {
  let list = users.value;
  const q = searchQuery.value.trim().toLowerCase();
  if (q)
    list = list.filter(
      (u) =>
        u.name?.toLowerCase().includes(q) ||
        displayComment(u.comment).toLowerCase().includes(q),
    );
  if (filterProfile.value)
    list = list.filter((u) => (u.profile || "default") === filterProfile.value);
  if (filterStatus.value)
    list = list.filter((u) => userStatus(u) === filterStatus.value);
  return list;
});

const PAGE_SIZES = [20, 50, 100] as const;
const usersPageSize = ref<20 | 50 | 100>(20);
const usersPage = ref(1);
const activePageSize = ref<20 | 50 | 100>(20);
const activePage = ref(1);

const pagedUsers = computed(() => {
  const start = (usersPage.value - 1) * usersPageSize.value;
  return filteredUsers.value.slice(start, start + usersPageSize.value);
});
const pagedActive = computed(() => {
  const start = (activePage.value - 1) * activePageSize.value;
  return active.value.slice(start, start + activePageSize.value);
});

function setUsersPageSize(n: 20 | 50 | 100) {
  usersPageSize.value = n;
  usersPage.value = 1;
}
function setActivePageSize(n: 20 | 50 | 100) {
  activePageSize.value = n;
  activePage.value = 1;
}

watch(users, () => {
  usersPage.value = 1;
});
watch(active, () => {
  activePage.value = 1;
});
watch([searchQuery, filterProfile, filterStatus], () => {
  usersPage.value = 1;
});

const profiles = ref<Record<string, string>[]>([]);
const profileMetas = ref<Record<string, ProfileMeta>>({});
const hotspotSettings = ref<HotspotSettings>({
  hotspotName: "",
  dnsName: "",
  currency: "",
});

// Used by the Business voucher template's QR code — points at the
// hotspot's own login page, since that DNS name is what devices on the
// hotspot network actually resolve. normalizeDnsName guards against a
// dnsName saved with a scheme already in it (e.g. from before this was
// validated on input) producing a doubled "http://http://...".
const loginUrl = computed(() => {
  const host = normalizeDnsName(hotspotSettings.value.dnsName);
  return host ? `http://${host}/login` : "";
});

const loading = ref(false);
const error = ref("");
const cleanupInstalled = ref<boolean | null>(null);
const cleanupInterval = ref("7d");
const cleanupToggling = ref(false);

async function toggleCleanup(enabled: boolean) {
  if (!store.activeId || cleanupToggling.value) return;
  if (
    !enabled &&
    !confirm(
      "Turn off auto-cleanup? Expired and quota-exhausted vouchers will no longer be removed automatically — they'll accumulate until you delete them manually.",
    )
  )
    return;
  cleanupToggling.value = true;
  try {
    const result = await putCleanupScheduler(
      store.activeId,
      enabled,
      cleanupInterval.value,
    );
    cleanupInstalled.value = result.installed;
    if (result.interval) cleanupInterval.value = result.interval;
  } catch {
    // non-critical — status stays as-is, user can retry
  } finally {
    cleanupToggling.value = false;
  }
}

const selected = ref<Set<string>>(new Set());
const bulkDeleting = ref(false);

const allSelected = computed(
  () =>
    filteredUsers.value.length > 0 &&
    filteredUsers.value.every((u) => selected.value.has(u[".id"])),
);
const someSelected = computed(
  () =>
    !allSelected.value &&
    filteredUsers.value.some((u) => selected.value.has(u[".id"])),
);

function toggleSelect(id: string) {
  const s = new Set(selected.value);
  if (s.has(id)) s.delete(id);
  else s.add(id);
  selected.value = s;
}

function toggleSelectAll() {
  selected.value = allSelected.value
    ? new Set()
    : new Set(filteredUsers.value.map((u) => u[".id"]));
}

async function removeSelected() {
  if (!store.activeId || !selected.value.size) return;
  if (
    !confirm(
      `Delete ${selected.value.size} selected user${selected.value.size > 1 ? "s" : ""}?`,
    )
  )
    return;
  bulkDeleting.value = true;
  const ids = [...selected.value];
  await Promise.allSettled(
    ids.map((id) => deleteHotspotUser(store.activeId!, id)),
  );
  selected.value = new Set();
  bulkDeleting.value = false;
  await loadUsers();
}

const showAdd = ref(false);
const adding = ref(false);
const addError = ref("");

const emptyForm = () => ({
  name: "",
  password: "",
  profile: "",
  limitUptimeRaw: "",
  limitBytesTotalValue: 0,
  limitBytesTotalUnit: "M" as "M" | "G",
  comment: "",
});
const form = ref(emptyForm());

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

function uptimePreviewFrom(raw: string): string {
  const p = parseShorthand(raw);
  if (!p) return raw ? "" : "";
  const totalDays = p.w * 7 + p.d;
  if (!totalDays && !p.h && !p.m) return "";
  return [
    totalDays ? `${totalDays}d` : "",
    p.h ? `${p.h}h` : "",
    p.m ? `${p.m}m` : "",
  ]
    .filter(Boolean)
    .join("");
}

function normalizeShorthand(raw: string): string {
  const p = parseShorthand(raw);
  if (!p) return raw;
  return [
    p.w ? `${p.w}w` : "",
    p.d ? `${p.d}d` : "",
    p.h ? `${p.h}h` : "",
    p.m ? `${p.m}m` : "",
  ]
    .filter(Boolean)
    .join("");
}

function selectedProfileValiditySeconds(profileName: string): number {
  if (!profileName) return 0;
  const meta = profileMetas.value[profileName];
  if (!meta?.validity) return 0;
  return shorthandToSeconds(meta.validity);
}

const formUptimePreview = computed(() =>
  uptimePreviewFrom(form.value.limitUptimeRaw),
);
const formUptimeWarning = computed(() => {
  if (!formUptimePreview.value || !form.value.profile) return false;
  const validitySecs = selectedProfileValiditySeconds(form.value.profile);
  if (!validitySecs) return false;
  return shorthandToSeconds(form.value.limitUptimeRaw) > validitySecs;
});
function normalizeFormUptime() {
  form.value.limitUptimeRaw = normalizeShorthand(form.value.limitUptimeRaw);
}

let activeSessionsTimer: ReturnType<typeof setInterval> | undefined;

function switchTab(key: "users" | "active") {
  tab.value = key;
  clearInterval(activeSessionsTimer);
  if (key === "active") {
    loadActive();
    activeSessionsTimer = setInterval(() => {
      if (!document.hidden) loadActive();
    }, 5 * 60_000);
  }
}

onUnmounted(() => clearInterval(activeSessionsTimer));

async function loadUsers() {
  if (!store.activeId) return;
  loading.value = true;
  error.value = "";
  try {
    const [u, p, m, s, cleanup] = await Promise.all([
      listHotspotUsers(store.activeId),
      listHotspotProfiles(store.activeId),
      getProfileMetas(store.activeId).catch(
        () => ({}) as Record<string, ProfileMeta>,
      ),
      getHotspotSettings(store.activeId).catch(
        () =>
          ({ hotspotName: "", dnsName: "", currency: "" }) as HotspotSettings,
      ),
      getCleanupScheduler(store.activeId).catch(() => null),
    ]);
    users.value = u;
    profiles.value = p;
    profileMetas.value = m;
    hotspotSettings.value = s;
    cleanupInstalled.value = cleanup?.installed ?? false;
    if (cleanup?.interval) cleanupInterval.value = cleanup.interval;
  } catch (e: any) {
    error.value = friendlyError(e, "Failed to load users");
  } finally {
    loading.value = false;
  }
}

async function loadActive() {
  if (!store.activeId) return;
  loading.value = true;
  error.value = "";
  try {
    active.value = await listHotspotActive(store.activeId);
  } catch (e: any) {
    error.value = friendlyError(e, "Failed to load sessions");
  } finally {
    loading.value = false;
  }
}

function openAdd() {
  form.value = emptyForm();
  addError.value = "";
  showAdd.value = true;
}

async function removeUser(userId: string) {
  if (!store.activeId || !confirm("Delete this user?")) return;
  await deleteHotspotUser(store.activeId, userId);
  await loadUsers();
}

async function submitAdd() {
  if (!store.activeId) return;
  adding.value = true;
  addError.value = "";
  try {
    const f = form.value;
    const limitUptime = formUptimePreview.value;
    const mul = f.limitBytesTotalUnit === "G" ? 1024 ** 3 : 1024 ** 2;
    const limitBytesTotal = f.limitBytesTotalValue
      ? String(f.limitBytesTotalValue * mul)
      : "";
    await createHotspotUser(store.activeId, {
      name: f.name,
      password: f.password,
      profile: f.profile,
      limitUptime,
      limitBytesTotal,
      rateLimit: "",
      comment: f.comment,
      expiryComment: "",
    });
    showAdd.value = false;
    form.value = emptyForm();
    await loadUsers();
  } catch (e: any) {
    addError.value = friendlyError(e, "Failed to create user");
  } finally {
    adding.value = false;
  }
}

async function toggleDisabled(u: Record<string, string>) {
  if (!store.activeId) return;
  const nextDisabled = u.disabled !== "true";
  try {
    await toggleHotspotUser(store.activeId, u[".id"], nextDisabled);
    await loadUsers();
  } catch (e: any) {
    error.value = friendlyError(e, "Failed to update user");
  }
}

const showPrintDialog = ref(false);
const printEntries = ref<VoucherEntry[]>([]);

// With no selection, prints every currently filtered/visible user instead —
// so the header Print button always does something useful.
function openPrint() {
  const source = selected.value.size > 0
    ? users.value.filter((u) => selected.value.has(u[".id"]))
    : filteredUsers.value;
  printEntries.value = source.map((u) => ({
    name: u.name,
    password: u.password ?? "",
    profile: u.profile,
  }));
  showPrintDialog.value = true;
}

function isUptimeExhausted(u: Record<string, string>): boolean {
  const limit = u["limit-uptime"],
    used = u["uptime"];
  if (!limit || !used) return false;
  return parseUptimeSeconds(used) >= parseUptimeSeconds(limit);
}

function extractExpEpoch(comment: string | undefined): number | null {
  if (!comment) return null;
  const m = comment.match(/^exp:(\d+)/);
  if (!m) return null;
  return parseInt(m[1]);
}

function userStatus(
  u: Record<string, string>,
):
  | "disabled"
  | "expired"
  | "expires-soon"
  | "limit-reached"
  | "waiting"
  | "active" {
  if (u.disabled === "true") return "disabled";
  // Uptime quota is enforced by RouterOS in real time — once hit, the user is
  // already cut off regardless of what the (possibly stale) exp: comment says.
  // Check it before calendar expiry so the badge reflects the reason that
  // actually applies, not whichever the code happens to check first.
  if (isUptimeExhausted(u)) return "limit-reached";
  const epoch = extractExpEpoch(u.comment);
  const now = Math.floor(Date.now() / 1000);
  if (epoch !== null && epoch < now) return "expired";
  if (epoch !== null && epoch - now < 86400) return "expires-soon";
  if (epoch === null) return "waiting";
  return "active";
}

const STATUS_META: Record<
  ReturnType<typeof userStatus>,
  { label: string; color: "green" | "amber" | "red" | "blue" | "muted" }
> = {
  active: { label: "Active", color: "green" },
  waiting: { label: "Waiting", color: "blue" },
  disabled: { label: "Disabled", color: "red" },
  expired: { label: "Expired", color: "red" },
  "expires-soon": { label: "Expires soon", color: "amber" },
  "limit-reached": { label: "Limit reached", color: "amber" },
};

function statusLabel(u: Record<string, string>): string {
  return STATUS_META[userStatus(u)].label;
}

function statusColor(
  u: Record<string, string>,
): "green" | "amber" | "red" | "blue" | "muted" {
  return STATUS_META[userStatus(u)].color;
}

function expiryLabel(u: Record<string, string>): string {
  // Uptime quota exhaustion takes precedence over the exp: comment — see
  // userStatus() for why: RouterOS enforces it in real time, so it's the
  // reason that actually applies even if the calendar exp: hasn't hit yet.
  if (isUptimeExhausted(u)) return "Limit reached";
  const epoch = extractExpEpoch(u.comment);
  if (epoch === null) return "—";
  const now = Math.floor(Date.now() / 1000);
  const diff = epoch - now;
  if (diff <= 0) return "Expired";
  if (diff < 3600) return `${Math.floor(diff / 60)}m left`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h left`;
  return `${Math.floor(diff / 86400)}d left`;
}

function expiryClass(u: Record<string, string>): string {
  if (isUptimeExhausted(u)) return "text-amber font-medium";
  const epoch = extractExpEpoch(u.comment);
  if (epoch === null) return "text-text-muted";
  const diff = epoch - Math.floor(Date.now() / 1000);
  if (diff <= 0) return "text-red font-medium";
  if (diff < 3600) return "text-amber";
  return "text-text-secondary";
}

function parseUptimeSeconds(uptime: string | undefined): number {
  if (!uptime) return 0;
  const m = uptime.match(
    /(?:(\d+)d\s*)?(?:(\d+)h\s*)?(?:(\d+)m\s*)?(?:(\d+)s)?/,
  );
  if (!m) return 0;
  return (
    parseInt(m[1] || "0") * 86400 +
    parseInt(m[2] || "0") * 3600 +
    parseInt(m[3] || "0") * 60 +
    parseInt(m[4] || "0")
  );
}

function formatSeconds(secs: number): string {
  if (secs <= 0) return "expired";
  const d = Math.floor(secs / 86400),
    h = Math.floor((secs % 86400) / 3600),
    m = Math.floor((secs % 3600) / 60);
  if (d > 0) return `${d}d ${h}h`;
  if (h > 0) return `${h}h ${m}m`;
  return `${m}m`;
}

function displayComment(comment: string | undefined): string {
  if (!comment) return "—";
  return comment.replace(/^exp:\d+\s*/, "") || "—";
}

function formatBytes(val: string | undefined): string {
  const n = parseInt(val ?? "0");
  if (!n) return "—";
  if (n >= 1024 ** 3) return (n / 1024 ** 3).toFixed(1) + " GB";
  if (n >= 1024 ** 2) return (n / 1024 ** 2).toFixed(1) + " MB";
  if (n >= 1024) return (n / 1024).toFixed(1) + " KB";
  return n + " B";
}

watch(() => store.activeId, loadUsers, { immediate: true });

// Deep-linked from Dashboard's "Active sessions" card via ?tab=active.
if (tab.value === "active") switchTab("active");
</script>
