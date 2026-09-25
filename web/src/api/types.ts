export interface RouterProfile {
  id: string
  name: string
  host: string
  port: number
  username: string
  useTls: boolean
  hotspotSettings?: {
    hotspotName?: string
    dnsName?: string
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
    dnsName: string
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

export interface SaleEntry {
  at: string    // ISO UTC timestamp
  profile: string
  price: string
  currency: string
  count: number
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
  template?: 'minimal' | 'voucher' | 'card'
}

export interface HotspotSettings {
  hotspotName: string
  dnsName: string
  currency: string
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
