<template>
  <div v-if="session && !session.configured" class="card p-6 text-sm text-text-secondary">
    {{ t('settings.cloud.notConfigured') }}
  </div>

  <div v-else-if="session" class="grid grid-cols-1 xl:grid-cols-[minmax(0,1.3fr)_minmax(0,1fr)] gap-7 items-start">
    <!-- This router -->
    <section class="card p-7 flex flex-col gap-4.5">
      <div class="flex flex-wrap items-center gap-3">
        <h2 class="type-display m-0 text-3xl">{{ t('settings.cloud.heading') }}</h2>
        <span class="chip bg-primary text-primary-fg">{{ t('settings.beta') }}</span>
      </div>
      <p class="leading-relaxed text-text-secondary">{{ t('settings.cloud.intro') }}</p>

      <div v-if="loadingRouter" class="flex justify-center py-6"><span class="spinner" /></div>
      <p v-else-if="routerError" role="alert" class="m-0 text-sm text-danger-text">{{ routerError }}</p>

      <template v-else-if="routerCloud">
        <!-- Connected: summary -->
        <div v-if="routerCloud.status.connected && !editing" class="flex flex-col gap-3">
          <div class="flex items-center gap-2.5 p-3.5 rounded-control border-2 border-ink bg-online-wash font-semibold">
            <CheckIcon class="size-5 shrink-0 stroke-3 text-online" aria-hidden="true" />
            <span>{{ t('settings.cloud.connected', { n: routerCloud.status.poolVouchers }) }}</span>
          </div>
          <p class="m-0 text-xs text-text-muted">{{ t('settings.cloud.routerId', { id: routerCloud.status.cloudRouterId }) }}</p>

          <div v-if="routerCloud.status.buyButton !== 'present'" role="status" class="flex flex-col gap-3 p-3.5 rounded-control border-2 border-ink bg-amber/10 text-sm">
            <div class="flex items-start gap-2.5">
              <ExclamationTriangleIcon class="size-5 shrink-0 text-amber" aria-hidden="true" />
              <span>{{ t('settings.cloud.noButton') }}</span>
            </div>
            <div class="flex flex-wrap gap-2.5">
              <button type="button" class="btn" @click="emit('openLoginPage')">{{ t('settings.cloud.usePikroPage') }}</button>
              <button type="button" class="btn" :disabled="retryingButton" @click="retryButton">
                <span v-if="retryingButton" class="spinner spinner--sm" />
                {{ t('settings.cloud.tryAgain') }}
              </button>
            </div>
          </div>

          <div class="flex flex-wrap items-start gap-4">
            <div v-if="qr" class="size-32 shrink-0 p-2 bg-white rounded-control border-2 border-ink [&>svg]:size-full" v-html="qr" />
            <div class="flex flex-col gap-2 min-w-0 flex-1">
              <span class="font-medium">{{ t('settings.cloud.buyLink') }}</span>
              <div class="flex items-center gap-2.5">
                <code class="text-xs break-all flex-1">{{ buyLink }}</code>
                <button type="button" class="btn shrink-0" @click="copyBuyUrl">
                  {{ copied ? t('settings.cloud.copied') : t('settings.cloud.copy') }}
                </button>
              </div>
              <p class="m-0 text-xs text-text-muted">{{ t('settings.cloud.buyLinkHint') }}</p>
            </div>
          </div>
          <div class="flex flex-wrap gap-2.5">
            <button v-if="qr" type="button" class="btn" @click="printCard">
              <PrinterIcon class="size-4.5" /> {{ t('settings.cloud.printCard') }}
            </button>
            <button type="button" class="btn" :disabled="!canSetUp" @click="editing = true">{{ t('settings.cloud.changePlans') }}</button>
            <button type="button" class="btn" @click="confirmingDisconnect = true">{{ t('settings.cloud.disconnect') }}</button>
          </div>
        </div>

        <!-- Setup form -->
        <form v-else class="flex flex-col gap-4" @submit.prevent="connect">
          <label class="flex flex-col gap-1">
            <span class="font-medium">{{ t('settings.cloud.ownerCountry') }}</span>
            <select v-model="ownerCountry" class="input-select" required>
              <option v-for="c in countries.main" :key="c.code" :value="c.code">{{ c.name }}</option>
              <optgroup :label="t('settings.cloud.allCountries')">
                <option v-for="c in countries.other" :key="c.code" :value="c.code">{{ c.name }}</option>
              </optgroup>
            </select>
          </label>
          <label class="flex flex-col gap-1">
            <span class="font-medium">{{ t('settings.cloud.ownerPhone') }}</span>
            <span class="join">
              <span class="inline-flex items-center px-3 border-2 border-ink rounded-control bg-ground font-mono text-sm">+{{ getCountryCallingCode(ownerCountry) }}</span>
              <input v-model="ownerPhone" class="input min-w-0" inputmode="tel" :placeholder="phoneExample" required />
            </span>
            <span v-if="ownerPhone && !parsedPhone?.isValid()" class="text-xs text-danger-text">
              {{ t('settings.cloud.invalidPhone', { country: countryName(ownerCountry) }) }}
            </span>
            <span class="text-xs text-text-muted">{{ t('settings.cloud.ownerPhoneHint') }}</span>
          </label>

          <fieldset class="m-0 p-0 border-0 flex flex-col gap-2">
            <legend class="font-medium mb-2">{{ t('settings.cloud.plans') }}</legend>
            <label
              v-for="p in routerCloud.plans" :key="p.profile"
              class="flex items-center gap-2.5 select-none"
              :class="p.price ? 'cursor-pointer' : 'opacity-50 cursor-not-allowed'"
            >
              <CheckboxRoot
                :model-value="selected.includes(p.profile)"
                :disabled="!p.price"
                class="size-5 rounded-md border-2 border-ink bg-paper flex items-center justify-center cursor-pointer data-[state=checked]:bg-ink disabled:cursor-not-allowed"
                @update:model-value="toggle(p.profile, $event === true)"
              >
                <CheckboxIndicator><CheckIcon class="size-3.5 text-paper" /></CheckboxIndicator>
              </CheckboxRoot>
              <span class="text-sm flex-1">{{ p.profile }}</span>
              <span class="text-sm font-semibold">
                {{ p.price ? `${p.price} ${routerCloud.currency}` : t('settings.cloud.noPrice') }}
              </span>
            </label>
            <p class="m-0 text-xs text-text-muted">{{ t('settings.cloud.plansHint') }}</p>
          </fieldset>

          <p v-if="!routerCloud.currency" role="alert" class="m-0 text-sm text-danger-text">{{ t('settings.cloud.noCurrency') }}</p>
          <p v-if="!canSetUp" class="m-0 text-sm text-text-secondary">{{ t('settings.cloud.signInFirst') }}</p>

          <div class="flex flex-wrap gap-2.5">
            <button
              type="submit"
              class="btn btn-primary min-h-12.5 text-[15px]"
              :disabled="!canSetUp || connecting || !selected.length || !parsedPhone?.isValid() || !routerCloud.currency"
            >
              <span v-if="connecting" class="spinner spinner--sm" />
              {{ connecting ? t('settings.cloud.connecting') : t('settings.cloud.connect') }}
            </button>
            <button v-if="editing" type="button" class="btn" :disabled="connecting" @click="editing = false">
              {{ t('settings.cloud.cancel') }}
            </button>
          </div>
        </form>

        <!-- Result of the last connect -->
        <ul v-if="result" class="m-0 p-0 list-none flex flex-col gap-2 text-sm">
          <li v-for="s in result.steps" :key="s.name" class="flex items-start gap-2.5">
            <ExclamationTriangleIcon v-if="s.warning" class="size-5 shrink-0 stroke-[2.75] text-amber" aria-hidden="true" />
            <CheckIcon v-else-if="s.ok" class="size-5 shrink-0 stroke-[2.75] text-online" aria-hidden="true" />
            <XMarkIcon v-else class="size-5 shrink-0 stroke-[2.75] text-danger-text" aria-hidden="true" />
            <span>
              {{ s.name }}
              <span v-if="s.skipped" class="text-text-muted">({{ t('settings.cloud.alreadyDone') }})</span>
              <span v-if="s.warning" class="block text-amber">{{ s.warning }}</span>
              <span v-if="s.error" class="block text-danger-text">{{ s.error }}</span>
            </span>
          </li>
        </ul>
        <p v-if="connectError" role="alert" class="m-0 text-sm text-danger-text">{{ connectError }}</p>
      </template>
    </section>

    <!-- The technician -->
    <section class="card p-6 flex flex-col gap-4">
      <template v-if="session.signedIn">
        <div>
          <h3 class="m-0 font-bold">{{ t('settings.cloud.signedInAs', { name: session.name || session.email }) }}</h3>
          <p class="m-0 mt-1 text-sm text-text-secondary">
            {{ t('settings.cloud.sessionUntil', { date: formatDate(session.expiresAt) }) }}
          </p>
        </div>
        <p v-if="session.revoked" role="alert" class="m-0 text-sm text-danger-text">{{ t('settings.cloud.revoked') }}</p>
        <button type="button" class="btn self-start" @click="signOut">{{ t('settings.cloud.signOut') }}</button>
      </template>

      <form v-else class="flex flex-col gap-4" @submit.prevent="codeSent ? signIn() : sendCode()">
        <div>
          <h3 class="m-0 font-bold">{{ t('settings.cloud.signInTitle') }}</h3>
          <p class="m-0 mt-1 text-sm text-text-secondary">{{ t('settings.cloud.signInBody') }}</p>
        </div>
        <label class="flex flex-col gap-1">
          <span class="font-medium">{{ t('settings.cloud.email') }}</span>
          <input v-model="email" type="email" class="input" autocomplete="email" placeholder="name@example.com" :disabled="codeSent" required />
        </label>
        <label v-if="codeSent" class="flex flex-col gap-1">
          <span class="font-medium">{{ t('settings.cloud.code') }}</span>
          <input v-model="code" class="input" inputmode="numeric" autocomplete="one-time-code" maxlength="6" required />
          <span class="text-xs text-text-muted">{{ t('settings.cloud.codeHint') }}</span>
        </label>
        <p v-if="signInError" role="alert" class="m-0 text-sm text-danger-text">{{ signInError }}</p>
        <div class="flex flex-wrap gap-2.5">
          <button type="submit" class="btn btn-ink" :disabled="busy">
            <span v-if="busy" class="spinner spinner--sm" />
            {{ codeSent ? t('settings.cloud.signIn') : t('settings.cloud.sendCode') }}
          </button>
          <button v-if="codeSent" type="button" class="btn" :disabled="busy" @click="codeSent = false; code = ''">
            {{ t('settings.cloud.otherEmail') }}
          </button>
        </div>
      </form>
    </section>

    <AppDialog v-model:open="confirmingDisconnect" :title="t('settings.cloud.disconnectTitle')">
      <p class="m-0 text-sm leading-normal text-text-secondary">{{ t('settings.cloud.disconnectBody') }}</p>
      <p v-if="disconnectError" role="alert" class="m-0 mt-3 text-sm text-danger-text">{{ disconnectError }}</p>
      <div class="flex justify-end gap-2.5 mt-6">
        <button type="button" class="btn" @click="confirmingDisconnect = false">{{ t('settings.cloud.cancel') }}</button>
        <button type="button" class="btn btn-primary" :disabled="disconnecting" @click="disconnect">
          <span v-if="disconnecting" class="spinner spinner--sm" />
          {{ t('settings.cloud.disconnect') }}
        </button>
      </div>
    </AppDialog>
  </div>

  <div v-else class="flex justify-center py-10"><span class="spinner" /></div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { CheckboxIndicator, CheckboxRoot } from 'reka-ui'
import { PrinterIcon, CheckIcon, ExclamationTriangleIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import { useRoutersStore } from '@/stores/routers'
import {
  getCloudSession, sendCloudCode, signInToCloud, signOutOfCloud,
  getRouterCloud, connectRouterToCloud, addCloudBuyButton, disconnectRouterFromCloud,
  type CloudSession, type RouterCloud, type CloudConnectResult,
} from '@/api'
import AppDialog from '@/components/AppDialog.vue'
import { qrSvg } from '@/utils/voucherTemplates/qr'
import { hotspotLoginUrl, printHtml } from '@/utils/vouchers'
import { cloudCardHtml } from '@/utils/cloudCard'
import { getCountries, getCountryCallingCode, getExampleNumber, parsePhoneNumberFromString, type CountryCode } from 'libphonenumber-js'
import examples from 'libphonenumber-js/mobile/examples'

const emit = defineEmits<{ openLoginPage: [] }>()

const { t, locale } = useI18n()
const store = useRoutersStore()

const errorText = (e: any, fallback: string) => e?.response?.data?.error ?? e?.message ?? fallback
const formatDate = (iso?: string) =>
  iso ? new Intl.DateTimeFormat(locale.value, { dateStyle: 'medium' }).format(new Date(iso)) : ''

// Technician session
const session = ref<CloudSession | null>(null)
const email = ref('')
const code = ref('')
const codeSent = ref(false)
const busy = ref(false)
const signInError = ref('')

async function loadSession() {
  try {
    session.value = await getCloudSession()
  } catch (e: any) {
    session.value = { configured: true, signedIn: false }
    signInError.value = errorText(e, t('settings.cloud.failed'))
  }
}

async function sendCode() {
  busy.value = true
  signInError.value = ''
  try {
    await sendCloudCode(email.value.trim())
    codeSent.value = true
  } catch (e: any) {
    // 429: a code sent less than a minute ago is still valid — let them type it.
    if (e?.response?.status === 429) codeSent.value = true
    else signInError.value = errorText(e, t('settings.cloud.failed'))
  } finally {
    busy.value = false
  }
}

async function signIn() {
  busy.value = true
  signInError.value = ''
  try {
    session.value = await signInToCloud(email.value.trim(), code.value)
    codeSent.value = false
    code.value = ''
  } catch (e: any) {
    signInError.value = errorText(e, t('settings.cloud.failed'))
  } finally {
    busy.value = false
  }
}

async function signOut() {
  session.value = await signOutOfCloud()
}

const canSetUp = computed(() => !!session.value?.signedIn && !session.value.revoked)

// This router
const routerCloud = ref<RouterCloud | null>(null)
const loadingRouter = ref(false)
const routerError = ref('')
const editing = ref(false)
const ownerCountry = ref<CountryCode>('TG')
// Typed without the country code, which comes from ownerCountry; a number
// typed with its own +code is kept as typed.
const ownerPhone = ref('')
const parsedPhone = computed(() => parsePhoneNumberFromString(ownerPhone.value, ownerCountry.value))
const phoneExample = computed(() => fakeNumber(getExampleNumber(ownerCountry.value, examples)?.formatNational() ?? ''))
// "90 11 23 45" -> "90 00 00 00": the library's examples may be real numbers,
// so only the first group (the prefix) is kept.
const fakeNumber = (s: string) => s.replace(/(\s.*)/, rest => rest.replace(/\d/g, '0'))

// Pikro Cloud's payment markets first, then every other country by name.
const MAIN_COUNTRIES: CountryCode[] = ['TG', 'BJ', 'BF', 'CI', 'SN', 'ML', 'NE', 'GN', 'GH', 'NG', 'CM']
const countryName = (code: string) => new Intl.DisplayNames([locale.value], { type: 'region' }).of(code) ?? code
const countries = computed(() => {
  const named = (codes: CountryCode[]) => codes.map(code => ({ code, name: countryName(code) }))
  return {
    main: named(MAIN_COUNTRIES),
    other: named(getCountries().filter(c => !MAIN_COUNTRIES.includes(c))).sort((a, b) => a.name.localeCompare(b.name, locale.value)),
  }
})
const selected = ref<string[]>([])
const connecting = ref(false)
const connectError = ref('')
const result = ref<CloudConnectResult | null>(null)

function toggle(profile: string, on: boolean) {
  selected.value = on ? [...selected.value, profile] : selected.value.filter(p => p !== profile)
}

async function loadRouter() {
  if (!store.activeId) return
  loadingRouter.value = true
  routerError.value = ''
  result.value = null
  editing.value = false
  try {
    routerCloud.value = await getRouterCloud(store.activeId)
    loginUrl.value = await hotspotLoginUrl(store.activeId)
    selected.value = routerCloud.value.plans.filter(p => p.price).map(p => p.profile)
  } catch (e: any) {
    routerError.value = errorText(e, t('settings.cloud.failed'))
  } finally {
    loadingRouter.value = false
  }
}

// The buy link for the QR code and the copy button. It carries the hotspot's
// login page, so after paying, "Se connecter au WiFi" can log the guest in;
// the login page's own button passes it from the router instead.
const loginUrl = ref('')
const buyLink = computed(() => {
  const url = routerCloud.value?.status.buyUrl ?? ''
  return url && loginUrl.value ? `${url}?login=${encodeURIComponent(loginUrl.value)}` : url
})

const qr = ref('')
watch(buyLink, async url => {
  qr.value = url ? await qrSvg(url, 256) : ''
})

async function printCard() {
  printHtml(await cloudCardHtml(store.active()?.hotspotSettings?.hotspotName ?? '', buyLink.value))
}

const copied = ref(false)
async function copyBuyUrl() {
  await navigator.clipboard.writeText(buyLink.value)
  copied.value = true
  setTimeout(() => (copied.value = false), 2000)
}

const retryingButton = ref(false)
async function retryButton() {
  if (!store.activeId) return
  retryingButton.value = true
  try {
    const step = await addCloudBuyButton(store.activeId)
    await loadRouter()
    result.value = { steps: [step], vouchers: 0, cloudRouterId: '' }
  } catch (e: any) {
    connectError.value = errorText(e, t('settings.cloud.failed'))
  } finally {
    retryingButton.value = false
  }
}

async function connect() {
  if (!store.activeId) return
  connecting.value = true
  connectError.value = ''
  result.value = null
  try {
    result.value = await connectRouterToCloud(store.activeId, parsedPhone.value!.number, selected.value)
    if (result.value.steps.every(s => s.ok)) {
      const kept = result.value
      await loadRouter()
      result.value = kept
    }
  } catch (e: any) {
    connectError.value = errorText(e, t('settings.cloud.failed'))
    if (e?.response?.status === 401) await loadSession()
  } finally {
    connecting.value = false
  }
}

const confirmingDisconnect = ref(false)
const disconnecting = ref(false)
const disconnectError = ref('')

async function disconnect() {
  if (!store.activeId) return
  disconnecting.value = true
  disconnectError.value = ''
  try {
    const { cloudNotified } = await disconnectRouterFromCloud(store.activeId)
    confirmingDisconnect.value = false
    await loadRouter()
    if (!cloudNotified) connectError.value = t('settings.cloud.disconnectNotNotified')
  } catch (e: any) {
    disconnectError.value = errorText(e, t('settings.cloud.failed'))
  } finally {
    disconnecting.value = false
  }
}

loadSession()
watch(() => store.activeId, loadRouter, { immediate: true })
</script>
