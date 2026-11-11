<template>
  <div>
    <div
      v-if="phase === 'loading'"
      class="card max-w-3xl flex flex-col items-center py-16 gap-3"
    >
      <span class="spinner" />
      <p class="m-0 text-sm text-text-secondary">{{ t('settings.setup.checking') }}</p>
    </div>

    <!-- ── Phase: blocked ─────────────────────────────────────────────────── -->
    <div v-else-if="phase === 'blocked'" class="max-w-3xl flex flex-col gap-7">
      <div class="card p-7 flex flex-col gap-5">
        <div
          class="flex text-sm items-start gap-3 p-4 rounded-control border-2 border-ink bg-primary-wash"
        >
          <ExclamationTriangleIcon class="size-5 shrink-0 mt-0.5" />
          <div class="flex-1 min-w-0">
            <p class="m-0 font-semibold">{{ t('settings.setup.blocked.title') }}</p>
            <p class="m-0 mt-1 text-text-secondary">
              {{ t('settings.setup.blocked.body') }}
            </p>
          </div>
        </div>

        <div class="bg-ground rounded-control p-4 text-sm flex flex-col gap-3">
          <p class="eyebrow m-0">{{ t('settings.setup.blocked.detected') }}</p>
          <div class="grid grid-cols-2 gap-x-4 gap-y-2">
            <span class="text-text-secondary">{{ t('settings.setup.blocked.serverName') }}</span>
            <span class="font-mono text-ink">{{
              preflight?.hotspotName || "—"
            }}</span>
            <span class="text-text-secondary">{{ t('settings.setup.blocked.interface') }}</span>
            <span class="font-mono text-ink">{{
              preflight?.hotspotOnIface || "—"
            }}</span>
            <span class="text-text-secondary">{{ t('settings.setup.blocked.profile') }}</span>
            <span class="font-mono text-ink">{{
              preflight?.hotspotProfile || "—"
            }}</span>
            <span class="text-text-secondary">{{ t('settings.setup.blocked.addressPool') }}</span>
            <span class="font-mono text-ink">{{
              preflight?.hotspotAddressPool || "—"
            }}</span>
            <span class="text-text-secondary">{{ t('settings.setup.blocked.dnsName') }}</span>
            <span class="font-mono text-ink">{{
              preflight?.hotspotDnsName || t('settings.setup.blocked.none')
            }}</span>
          </div>
        </div>

        <div class="flex flex-wrap items-center gap-3">
          <RouterLink to="/hotspot/profiles" class="btn btn-ink">
            {{ t('settings.setup.blocked.addProfiles') }}
          </RouterLink>
          <button class="btn" @click="runPreflight">
            <ArrowPathIcon class="size-4" />
            {{ t('settings.setup.blocked.recheck') }}
          </button>
          <!-- <button
            class="btn border-danger-text text-danger-text"
            :disabled="tearing"
            @click="runTeardown"
          >
            <TrashIcon class="size-3.5" />
            {{ tearing ? "Resetting…" : "Reset config" }}
          </button> -->
        </div>
      </div>

      <div v-if="teardownResult" class="card p-6 flex flex-col gap-3">
        <p class="eyebrow m-0">{{ t('settings.setup.resetResult') }}</p>
        <ul class="m-0 p-0 list-none flex flex-col gap-2">
          <li
            v-for="step in teardownResult.steps"
            :key="step.name"
            class="flex items-center gap-2 text-sm"
          >
            <CheckCircleIcon
              v-if="step.ok && !step.skipped"
              class="size-4 shrink-0 text-online"
            />
            <MinusCircleIcon
              v-else-if="step.skipped"
              class="size-4 text-text-muted shrink-0"
            />
            <XCircleIcon v-else class="size-4 shrink-0 text-danger-text" />
            <span
              :class="step.skipped ? 'text-text-muted' : 'text-text-secondary'"
              >{{ step.name }}</span
            >
            <span v-if="step.error" class="font-mono text-danger-text"
              >— {{ step.error }}</span
            >
          </li>
        </ul>
      </div>
    </div>

    <!-- ── Phase: form ────────────────────────────────────────────────────── -->
    <div v-else-if="phase === 'form'" class="max-w-3xl flex flex-col gap-7">
      <div class="card p-7 flex flex-col gap-6">
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
              <StepperTrigger class="flex items-center gap-2" as="div">
                <StepperIndicator
                  class="size-7 rounded-full border-2 border-ink flex items-center justify-center text-xs font-bold transition-colors shrink-0 data-[state=active]:bg-ink data-[state=active]:text-paper data-[state=completed]:bg-online-wash data-[state=completed]:text-online data-[state=inactive]:bg-ground data-[state=inactive]:border-hairline data-[state=inactive]:text-text-muted"
                >
                  {{ step > i + 1 ? "✓" : i + 1 }}
                </StepperIndicator>
                <StepperTitle
                  class="text-sm hidden sm:inline"
                  :class="
                    step === i + 1
                      ? 'text-ink font-bold'
                      : 'text-text-muted font-medium'
                  "
                >
                  {{ t(`settings.setup.steps.${label.toLowerCase()}`) }}
                </StepperTitle>
              </StepperTrigger>
              <StepperSeparator
                v-if="i < STEPS.length - 1"
                class="flex-1 h-0.5 bg-hairline"
              />
            </StepperItem>
          </template>
        </StepperRoot>

        <div
          v-if="preflightError"
          role="alert"
          class="flex items-start gap-2 p-3 rounded-control border-2 border-danger-text bg-danger-wash text-sm text-danger-text"
        >
          <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5" />
          <span>{{ t('settings.setup.preflightError', { error: preflightError }) }}</span>
        </div>

        <form @submit.prevent="onFormSubmit" class="flex flex-col gap-6">
          <!-- Step: Bridge (conditional) -->
          <div
            v-if="STEPS[step - 1] === 'Bridge'"
            class="flex flex-col gap-5"
          >
            <div
              class="flex items-start gap-2 p-3 rounded-control border-2 border-ink bg-primary-wash text-sm"
            >
              <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5" />
              <span>{{ t('settings.setup.bridge.notFound') }}</span>
            </div>

            <label class="flex flex-col gap-1">
              <span class="font-medium">{{ t('settings.setup.bridge.name') }}</span>
              <input
                v-model="bridge.name"
                class="input"
                placeholder="bridge1"
                required
              />
            </label>

            <div class="flex flex-col gap-1.5">
              <span class="font-medium">{{ t('settings.setup.bridge.ports') }}</span>
              <p class="m-0 text-xs text-text-muted">
                {{ t('settings.setup.bridge.portsHint') }}
              </p>
              <div class="grid grid-cols-2 gap-2.5 mt-1">
                <label
                  v-for="i in etherInterfaces"
                  :key="i.name"
                  class="flex items-center gap-2.5 cursor-pointer select-none text-sm"
                >
                  <span
                    class="relative inline-flex items-center justify-center size-5 rounded-md border-2 border-ink shrink-0 transition-colors"
                    :class="bridge.ports.includes(i.name) ? 'bg-ink' : 'bg-paper'"
                  >
                    <CheckIcon
                      v-if="bridge.ports.includes(i.name)"
                      class="size-3.5 text-paper"
                    />
                  </span>
                  <input
                    type="checkbox"
                    :value="i.name"
                    v-model="bridge.ports"
                    class="sr-only"
                  />
                  <span class="text-text-secondary"
                    >{{ i.name }}{{ !i.running ? " — " + t('settings.setup.down') : "" }}</span
                  >
                </label>
              </div>
            </div>
          </div>

          <!-- Step: Network -->
          <div
            v-else-if="STEPS[step - 1] === 'Network'"
            class="flex flex-col gap-5"
          >
            <h2 class="type-display m-0 text-2xl">{{ t('settings.setup.network.title') }}</h2>

            <label class="flex flex-col gap-1">
              <span class="font-medium"
                >{{ t('settings.setup.network.lan') }} <span class="text-danger-text">*</span></span
              >
              <AppSelect
                v-if="lanOptions.length"
                v-model="form.lanIface"
                :options="lanOptions"
                :placeholder="t('settings.setup.network.selectInterface')"
              />
              <input
                v-else
                v-model="form.lanIface"
                class="input"
                :placeholder="t('settings.setup.network.lanPlaceholder')"
              />
              <p class="m-0 text-xs text-text-muted">
                {{ t('settings.setup.network.lanHint') }}
              </p>
            </label>

            <label class="flex flex-col gap-1">
              <span class="font-medium"
                >{{ t('settings.setup.network.wan') }} <span class="text-danger-text">*</span></span
              >
              <AppSelect
                v-if="wanOptions.length"
                v-model="form.wanIface"
                :options="wanOptions"
                :placeholder="t('settings.setup.network.selectInterface')"
              />
              <input
                v-else
                v-model="form.wanIface"
                class="input"
                :placeholder="t('settings.setup.network.wanPlaceholder')"
              />
              <p class="m-0 text-xs text-text-muted">
                {{ t('settings.setup.network.wanHint') }}
              </p>
            </label>

            <label class="flex flex-col gap-1">
              <span class="font-medium"
                >{{ t('settings.setup.network.hotspotName') }} <span class="text-danger-text">*</span></span
              >
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
              <p class="m-0 text-xs text-text-muted">
                {{ t('settings.setup.network.redirectedTo') }}
                <span class="font-mono"
                  >{{ form.hotspotName || "myhotspot"
                  }}{{ form.hotspotExtension }}</span
                >
              </p>
            </label>
          </div>

          <!-- Step: Subnet -->
          <div
            v-else-if="STEPS[step - 1] === 'Subnet'"
            class="flex flex-col gap-5"
          >
            <div
              v-if="dhcpChecking"
              class="flex items-center gap-2.5 text-sm text-text-muted"
            >
              <span class="spinner spinner--sm" />
              {{ t('settings.setup.subnet.checkingDhcp', { iface: form.lanIface }) }}
            </div>
            <div
              v-else-if="dhcpCheck?.exists"
              class="flex items-start gap-2 p-3 rounded-control border-2 border-ink bg-primary-wash text-sm"
            >
              <ExclamationTriangleIcon class="size-4 shrink-0 mt-0.5" />
              <i18n-t keypath="settings.setup.subnet.dhcpExists" tag="span">
                <template #iface>{{ form.lanIface }}</template>
                <template #subnet>
                  <span class="font-mono">{{ dhcpCheck.subnet }}</span>
                </template>
              </i18n-t>
            </div>

            <label class="flex flex-col gap-1">
              <span class="font-medium"
                >{{ t('settings.setup.subnet.label') }} <span class="text-danger-text">*</span></span
              >
              <input
                v-model="form.subnet"
                class="input font-mono"
                placeholder="192.168.88.0/24"
                :disabled="dhcpCheck?.exists"
              />
              <p v-if="subnetError" role="alert" class="m-0 text-sm text-danger-text">
                {{ subnetError }}
              </p>
            </label>

            <div v-if="derived" class="bg-ground rounded-control p-4 text-sm">
              <div class="grid grid-cols-2 gap-x-4 gap-y-2">
                <span class="text-text-secondary">{{ t('settings.setup.subnet.gateway') }}</span>
                <span class="font-mono text-ink font-medium">{{
                  derived.gateway
                }}</span>
                <span class="text-text-secondary">{{ t('settings.setup.subnet.dhcpPool') }}</span>
                <span class="font-mono text-ink font-medium"
                  >{{ derived.poolStart }} – {{ derived.poolEnd }}</span
                >
                <span class="text-text-secondary">{{ t('settings.setup.subnet.profileName') }}</span>
                <span class="font-mono text-ink font-medium"
                  >pikro-profile</span
                >
              </div>
            </div>
          </div>

          <!-- Step: Review -->
          <template v-else-if="STEPS[step - 1] === 'Review'">
            <div class="flex flex-col gap-5">
              <div
                class="flex items-start gap-2 text-sm p-3 rounded-control border-2 border-ink bg-primary-wash"
              >
                <ExclamationTriangleIcon class="size-5 shrink-0 mt-0.5" />
                <span>{{ t('settings.setup.review.warning') }}</span>
              </div>

              <div class="bg-ground rounded-control p-4 flex flex-col gap-3">
                <h2 class="m-0 font-bold">{{ t('settings.setup.review.heading') }}</h2>
                <ul class="m-0 p-0 list-none flex flex-col gap-1 text-sm text-text-secondary">
                  <li v-if="needsBridge">
                    •
                    {{
                      bridge.ports.length
                        ? t("settings.setup.review.createBridgeWithPorts", {
                            name: bridgeName,
                            ports: bridge.ports.join(", "),
                          })
                        : t("settings.setup.review.createBridge", { name: bridgeName })
                    }}
                  </li>
                  <li>
                    •
                    {{
                      t("settings.setup.review.assign", {
                        ip: derived?.gateway ?? "…",
                        iface: lanIfaceName,
                      })
                    }}
                  </li>
                  <li>
                    • {{ t('settings.setup.review.createPool') }} <span class="font-mono">hotspot-pool</span>
                  </li>
                  <li>• {{ t('settings.setup.review.createDhcp') }}</li>
                  <li>• {{ t('settings.setup.review.createProfile') }}</li>
                  <li>
                    • {{ t('settings.setup.review.enableHotspot', { iface: lanIfaceName }) }}
                  </li>
                  <li>
                    • {{
                      t("settings.setup.review.nat", {
                        iface: form.wanIface || t("settings.setup.review.wanFallback"),
                      })
                    }}
                  </li>
                  <li>• {{ t('settings.setup.review.dns') }}</li>
                </ul>
              </div>
            </div>
          </template>

          <!-- Step: Cleanup -->
          <template v-else-if="STEPS[step - 1] === 'Cleanup'">
            <div class="flex flex-col gap-4">
              <h2 class="type-display m-0 text-2xl">{{ t('settings.setup.cleanup.title') }}</h2>
              <p class="m-0 text-sm leading-normal text-text-secondary">
                {{ t('settings.setup.cleanup.body') }}
              </p>
              <label class="flex flex-col gap-1">
                <span class="eyebrow">{{ t('settings.setup.cleanup.runEvery') }}</span>
                <AppSelect
                  v-model="form.cleanupInterval"
                  :options="CLEANUP_INTERVAL_OPTIONS"
                />
              </label>
            </div>
          </template>

          <div class="flex flex-wrap items-center gap-3">
            <button type="button" class="btn" @click="onBack">
              {{ step === 1 ? t('settings.setup.recheckRouter') : t('settings.setup.back') }}
            </button>
            <button
              v-if="step < STEPS.length"
              type="button"
              class="btn btn-ink"
              :disabled="!canAdvance"
              @click="step++"
            >
              {{ t('settings.setup.next') }}
            </button>
            <button
              v-else
              type="submit"
              class="btn btn-primary"
              :disabled="!canSubmit"
            >
              {{ t('settings.setup.run') }}
            </button>
            <!-- <button
              type="button"
              class="btn ml-auto border-danger-text text-danger-text"
              :disabled="tearing"
              @click="runTeardown"
            >
              <TrashIcon class="size-3.5" />
              {{ tearing ? "Resetting…" : "Reset config" }}
            </button> -->
          </div>
        </form>
      </div>

      <div v-if="teardownResult" class="card p-6 flex flex-col gap-3">
        <p class="eyebrow m-0">{{ t('settings.setup.resetResult') }}</p>
        <ul class="m-0 p-0 list-none flex flex-col gap-2">
          <li
            v-for="step in teardownResult.steps"
            :key="step.name"
            class="flex items-center gap-2 text-sm"
          >
            <CheckCircleIcon
              v-if="step.ok && !step.skipped"
              class="size-4 shrink-0 text-online"
            />
            <MinusCircleIcon
              v-else-if="step.skipped"
              class="size-4 text-text-muted shrink-0"
            />
            <XCircleIcon v-else class="size-4 shrink-0 text-danger-text" />
            <span
              :class="step.skipped ? 'text-text-muted' : 'text-text-secondary'"
              >{{ step.name }}</span
            >
            <span v-if="step.error" class="font-mono text-danger-text"
              >— {{ step.error }}</span
            >
          </li>
        </ul>
      </div>
    </div>

    <!-- ── Phase: running ─────────────────────────────────────────────────── -->
    <div
      v-else-if="phase === 'running'"
      class="card max-w-3xl flex flex-col items-center py-16 gap-3"
    >
      <span class="spinner" />
      <p class="m-0 font-bold">{{ t('settings.setup.running.title') }}</p>
      <p class="m-0 text-sm text-text-muted">
        {{ t('settings.setup.running.wait') }}
      </p>
    </div>

    <!-- ── Phase: result ──────────────────────────────────────────────────── -->
    <div v-else-if="phase === 'result' && result" class="max-w-3xl flex flex-col gap-7">
      <div class="card p-7 flex flex-col gap-5">
        <div
          class="flex items-center gap-2.5 p-3.5 rounded-control border-2 text-sm font-semibold"
          :class="
            result.success
              ? 'bg-online-wash border-ink'
              : 'bg-danger-wash border-danger-text text-danger-text'
          "
        >
          <CheckCircleIcon v-if="result.success" class="size-5 shrink-0 text-online" />
          <ExclamationTriangleIcon v-else class="size-5 shrink-0" />
          {{
            result.success
              ? t('settings.setup.result.success')
              : t('settings.setup.result.withErrors')
          }}
        </div>

        <ul class="m-0 p-0 list-none flex flex-col gap-2.5">
          <li
            v-for="step in result.steps"
            :key="step.name"
            class="flex items-start gap-2 text-sm"
          >
            <CheckCircleIcon
              v-if="step.ok && !step.skipped"
              class="size-4 shrink-0 mt-0.5 text-online"
            />
            <MinusCircleIcon
              v-else-if="step.skipped"
              class="size-4 text-text-muted shrink-0 mt-0.5"
            />
            <XCircleIcon v-else class="size-4 shrink-0 mt-0.5 text-danger-text" />
            <div class="flex-1 min-w-0">
              <span
                :class="
                  step.ok
                    ? 'text-text-secondary'
                    : step.skipped
                      ? 'text-text-muted'
                      : 'font-medium text-danger-text'
                "
              >
                {{ step.name }}
                <span v-if="step.skipped" class="text-text-muted font-normal">
                  — {{ t('settings.setup.result.alreadyConfigured') }}</span
                >
              </span>
              <p v-if="step.error" class="m-0 font-mono mt-0.5 break-all text-danger-text">
                {{ step.error }}
              </p>
            </div>
          </li>
        </ul>
      </div>

      <div v-if="result.success" class="card p-6 flex flex-col gap-4">
        <h3 class="m-0 font-bold">{{ t('settings.setup.result.nextHeading') }}</h3>
        <ol class="m-0 p-0 list-none flex flex-col gap-3 text-sm text-text-secondary">
          <li class="flex items-start gap-2.5">
            <span
              class="size-6 shrink-0 rounded-full border-2 border-ink flex items-center justify-center text-xs font-bold bg-primary-wash text-ink"
              >1</span
            >
            <div>
              <span>{{ t('settings.setup.result.next1') }}</span>
              <span class="block mt-1 p-2.5 rounded-control border-2 border-ink bg-primary-wash text-ink">
                <i18n-t keypath="settings.setup.result.adminNote">
                  <template #admin><span class="font-mono">admin</span></template>
                </i18n-t>
              </span>
            </div>
          </li>
          <li class="flex items-start gap-2.5">
            <span
              class="size-6 shrink-0 rounded-full border-2 border-ink flex items-center justify-center text-xs font-bold bg-primary-wash text-ink"
              >2</span
            >
            <i18n-t keypath="settings.setup.result.next2" tag="span">
              <template #target>
                <span class="font-mono">{{
                  form.hotspotName
                    ? form.hotspotName + form.hotspotExtension
                    : t('settings.setup.result.loginPageFallback')
                }}</span>
              </template>
            </i18n-t>
          </li>
          <li class="flex items-start gap-2.5">
            <span
              class="size-6 shrink-0 rounded-full border-2 border-ink flex items-center justify-center text-xs font-bold bg-primary-wash text-ink"
              >3</span
            >
            <span>{{ t('settings.setup.result.next3') }}</span>
          </li>
          <li class="flex items-start gap-2.5">
            <span
              class="size-6 shrink-0 rounded-full border-2 border-ink flex items-center justify-center text-xs font-bold bg-primary-wash text-ink"
              >4</span
            >
            <i18n-t keypath="settings.setup.result.next4" tag="span">
              <template #link>
                <button
                  type="button"
                  class="underline font-semibold text-ink hover:text-text-secondary transition-colors cursor-pointer"
                  @click="emit('openLoginPage')"
                  >{{ t("settings.tabs.login") }}</button
                >
              </template>
            </i18n-t>
          </li>
        </ol>
      </div>

      <div v-if="!result.success" class="card p-6 text-sm flex flex-col gap-1">
        <p class="m-0 font-bold">{{ t('settings.setup.result.failedHeading') }}</p>
        <p class="m-0 text-text-secondary">
          {{ t('settings.setup.result.failedBody') }}
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-3">
        <RouterLink
          v-if="result.success"
          to="/hotspot/users"
          class="btn btn-primary"
        >
          {{ t('settings.setup.result.addUsers') }}
        </RouterLink>
        <button class="btn" @click="reset">
          <ArrowPathIcon class="size-4" />
          {{ t('settings.setup.result.runAgain') }}
        </button>
      </div>
    </div>

    <!-- ── Teardown confirmation dialog ─────────────────────────────────── -->
    <AppDialog
      :open="showTeardownConfirm"
      :title="t('settings.setup.teardown.title')"
      @update:open="showTeardownConfirm = $event"
    >
      <div class="flex flex-col gap-4">
        <div class="bg-ground rounded-control p-3.5 text-sm flex flex-col gap-1.5">
          <div class="flex justify-between font-medium">
            <span class="text-text-secondary">{{ t('settings.setup.teardown.router') }}</span>
            <span class="text-ink">{{
              store.routers.find((r) => r.id === store.activeId)?.name
            }}</span>
          </div>
          <div class="flex justify-between font-medium">
            <span class="text-text-secondary">{{ t('settings.setup.teardown.host') }}</span>
            <span class="font-mono text-ink">{{
              store.routers.find((r) => r.id === store.activeId)?.host
            }}</span>
          </div>
          <div class="flex justify-between font-medium">
            <span class="text-text-secondary">{{ t('settings.setup.teardown.hotspotInterface') }}</span>
            <span class="font-mono text-ink">{{
              preflight?.hotspotOnIface || form.lanIface || "—"
            }}</span>
          </div>
        </div>

        <div>
          <p class="m-0 font-bold text-sm mb-1.5">{{ t('settings.setup.teardown.willRemove') }}</p>
          <ul class="m-0 p-0 list-none text-sm text-text-secondary flex flex-col gap-1">
            <li class="flex items-center gap-1.5 shrink-0 text-danger-text">
              <TrashIcon class="size-4" /> {{ t('settings.setup.teardown.server') }}
            </li>
            <li class="flex items-center gap-1.5 shrink-0 text-danger-text">
              <TrashIcon class="size-4" /> {{ t('settings.setup.teardown.serverProfile') }}
            </li>
            <li class="flex items-center gap-1.5 shrink-0 text-danger-text">
              <TrashIcon class="size-4" /> {{ t('settings.setup.teardown.nat') }}
            </li>
            <li class="flex items-center gap-1.5 shrink-0 text-danger-text">
              <TrashIcon class="size-4" /> {{ t('settings.setup.teardown.walledGarden') }}
            </li>
          </ul>
          <p class="m-0 text-sm text-text-secondary font-medium mt-2">
            {{ t('settings.setup.teardown.kept') }}
          </p>
        </div>

        <div class="flex justify-end gap-2.5 pt-1">
          <button class="btn" @click="showTeardownConfirm = false">
            {{ t('settings.setup.teardown.cancel') }}
          </button>
          <button class="btn btn-danger" @click="confirmTeardown">
            {{ t('settings.setup.teardown.confirm') }}
          </button>
        </div>
      </div>
    </AppDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { useI18n } from "vue-i18n";
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
  TrashIcon,
  Cog6ToothIcon,
  CheckIcon,
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
import { renderGuestPages } from "./loginPageTemplates";
import AppDialog from "@/components/AppDialog.vue";
import AppSelect, { type SelectOption } from "@/components/AppSelect.vue";

const { t, locale } = useI18n();

const emit = defineEmits<{ openLoginPage: [] }>();

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
  cleanupInterval: "1h",
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

const bridgeName = computed(
  () => bridge.value.name || t("settings.setup.review.bridgeFallback"),
);
const lanIfaceName = computed(
  () => form.value.lanIface || t("settings.setup.review.lanFallback"),
);

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
  return `${i.name} [${i.type}]${!i.running ? " — " + t('settings.setup.down') : ""}${i.comment ? ` (${i.comment})` : ""}`;
}

const lanOptions = computed<SelectOption[]>(() => {
  if (!interfaces.value.length) return [];
  if (needsBridge.value && bridge.value.name) {
    return [
      { value: bridge.value.name, label: t('settings.setup.newBridge', { name: bridge.value.name }) },
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
  if (!m) return t('settings.setup.subnet.invalidCidr');
  const octets = [m[1], m[2], m[3], m[4]].map(Number);
  if (octets.some((o) => o > 255)) return t('settings.setup.subnet.octets');
  const prefix = parseInt(m[5]);
  if (prefix < 16 || prefix > 30) {
    return t('settings.setup.subnet.prefix');
  }
  if (octets[3] !== 0) {
    return t('settings.setup.subnet.networkAddress');
  }
  if (octets[0] === 127 || (octets[0] === 169 && octets[1] === 254)) {
    return t('settings.setup.subnet.reserved');
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
    form.value.hotspotName.trim() !== "" &&
    derived.value !== null &&
    (!needsBridge.value || bridge.value.name.trim() !== ""),
);

const canAdvance = computed(() => {
  const current = STEPS.value[step.value - 1];
  if (current === "Bridge") return bridge.value.name.trim() !== "";
  if (current === "Network")
    return (
      form.value.lanIface.trim() !== "" &&
      form.value.wanIface.trim() !== "" &&
      form.value.hotspotName.trim() !== ""
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
    preflightError.value = friendlyError(e, t('settings.setup.connectionFailed'));
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
      // Pikro's guest pages with their defaults, as Settings > Login page
      // makes them until the operator changes them there.
      pages: renderGuestPages({
        title: "Wi-Fi",
        subtitle: "$(hostname)",
        mode: "code",
        language: locale.value === "fr" ? "fr" : "en",
        prices: [],
      }),
      ...(needsBridge.value
        ? { newBridgeName: bridge.value.name, bridgePorts: bridge.value.ports }
        : {}),
    });
  } catch (e: any) {
    result.value = {
      success: false,
      steps: [
        {
          name: t('settings.setup.setupRequest'),
          ok: false,
          error: friendlyError(e, t('settings.setup.requestFailed')),
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
          name: t('settings.setup.teardown.step'),
          ok: false,
          error: friendlyError(e, t('settings.setup.teardown.failed')),
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
