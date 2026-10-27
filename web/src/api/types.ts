export interface RouterProfile {
  id: string
  name: string
  host: string
  port: number
  username: string
  useTls: boolean
  hotspotSettings?: {
    hotspotName?: string
    currency?: string
  }
}

export interface NewRouterProfile {
  name: string
  host: string
  port: number
  username: string
  password: string
  useTls: boolean
  hotspotSettings?: {
    hotspotName: string
    currency: string
  }
}

export interface NewHotspotUser {
  name: string
  password: string
  profile: string
  limitUptime: string
  limitBytesTotal: string
  rateLimit: string
  comment: string
  expiryComment?: string
  price?: string
  currency?: string
}

// Sales ledger: vouchers are recorded when generated (a batch) and counted as
// sold once activated (first used). Revenue only counts activations.
export interface SalesBatch {
  id: string
  at: string    // ISO UTC timestamp
  profile: string
  price: string
  currency: string
  count: number
}

export interface SalesActivation {
  batch: string
  name: string
  at: string    // ISO UTC first-login time
  profile: string
  price: string
  currency: string
}

export interface SalesData {
  batches: SalesBatch[]
  activations: SalesActivation[]
}

export interface NewHotspotUsersBatch {
  /** Same id for every chunk of one generation run. */
  batchId: string
  profile: string
  limitUptime: string
  limitBytesTotal: string
  comment: string
  price: string
  currency: string
  users: { name: string; password: string }[]
}

export interface BatchCreateResult {
  name: string
  ok: boolean
  error?: string
}

export interface UpdateHotspotUser {
  password?: string
  profile?: string
  limitUptime: string
  limitBytesTotal: string
}

export interface ProfileMeta {
  validity: string
  price: string
}

export interface LoginPageSettings {
  title?: string
  subtitle?: string
  /** Legacy design picker, read once to derive loginMode. */
  template?: 'minimal' | 'voucher' | 'card'
  loginMode?: 'code' | 'password'
  language?: 'en' | 'fr'
  showPrices?: boolean
  contact?: string
  uploadedAt?: string
}

export interface HotspotSettings {
  hotspotName: string
  currency: string
  profileMetas?: Record<string, ProfileMeta>
  loginPage?: LoginPageSettings
}

export interface CleanupSchedulerStatus {
  installed: boolean
  interval: string
}

export interface HotspotProfileParams {
  name: string
  addressPool: string
  sharedUsers: string
  rateLimit: string
  validity: string
  price: string
}

export interface InterfaceInfo {
  name: string
  type: string
  running: boolean
  comment: string
}

export interface PreflightResult {
  interfaces: InterfaceInfo[]
  hotspotExists: boolean
  hotspotOnIface: string
  hotspotName: string
  hotspotProfile: string
  hotspotDnsName: string
  hotspotAddressPool: string
  hotspotAddress: string
}

export interface SetupStepResult {
  name: string
  ok: boolean
  error?: string
  skipped?: boolean
}

export interface SetupResult {
  steps: SetupStepResult[]
  success: boolean
}

export interface TeardownResult {
  steps: SetupStepResult[]
}

export interface DiscoveredDevice {
  ip: string
  mac: string
  identity: string
  version: string
  board: string
  platform: string
  uptime: string
  iface: string
}
