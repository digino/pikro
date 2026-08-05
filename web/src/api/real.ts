import axios from 'axios'
import type {
  RouterProfile, NewRouterProfile,
  NewHotspotUser, UpdateHotspotUser,
  ProfileMeta, HotspotSettings, LoginPageSettings,
  CleanupSchedulerStatus, HotspotProfileParams,
  PreflightResult, SetupResult, TeardownResult,
  DiscoveredDevice, SaleEntry,
} from './types'

export type * from './types'

const api = axios.create({ baseURL: '/api' })

export const getAppVersion = () =>
  api.get<{ version: string }>('/version').then(r => r.data.version)

// Routers
export const listRouters = () =>
  api.get<RouterProfile[]>('/routers').then(r => r.data)

export const addRouter = (profile: NewRouterProfile) =>
  api.post<{ id: string }>('/routers', profile).then(r => r.data)

export interface UpdateRouterProfile {
  name: string
  host: string
  port: number
  username: string
  password?: string
  useTls: boolean
  hotspotSettings?: {
    hotspotName: string
    dnsName: string
    currency: string
  }
}

export const updateRouter = (id: string, profile: UpdateRouterProfile) =>
  api.patch(`/routers/${id}`, profile)

export const deleteRouter = (id: string) =>
  api.delete(`/routers/${id}`)

export const testRouter = (id: string) =>
  api.get<{ ok: boolean }>(`/routers/${id}/test`).then(r => r.data)

// System
export const getSystemResource = (id: string) =>
  api.get<Record<string, string>>(`/routers/${id}/system/resource`).then(r => r.data)

export const getWanIP = (id: string) =>
  api.get<{ ip: string }>(`/routers/${id}/system/wan-ip`).then(r => r.data.ip)

export const getInterfaceTraffic = (id: string) =>
  api.get<Record<string, string>[]>(`/routers/${id}/system/traffic`).then(r => r.data)

export const getSystemClock = (id: string) =>
  api.get<Record<string, string>>(`/routers/${id}/system/clock`).then(r => r.data)

export interface PollSnapshot {
  resource: Record<string, string>
  traffic: Record<string, string>[]
  addresses: Record<string, string>[]
  clock: Record<string, string>
  interfaces: Record<string, string>[]
}

export const getPollSnapshot = (id: string) =>
  api.get<PollSnapshot>(`/routers/${id}/system/poll`).then(r => r.data)

export const getSystemLogs = (id: string) =>
  api.get<Record<string, string>[]>(`/routers/${id}/system/logs`).then(r => r.data)

export const getDHCPLeases = (id: string) =>
  api.get<Record<string, string>[]>(`/routers/${id}/system/dhcp-leases`).then(r => r.data)

export interface MikhmonMigrationResult {
  usersScanned: number
  usersUnused: number
  usersConverted: number
  usersSkipped: number
  profilesScanned: number
  profilesConverted: number
}

export const migrateFromMikhmon = (id: string) =>
  api.post<MikhmonMigrationResult>(`/routers/${id}/hotspot/migrate-mikhmon`).then(r => r.data)

// Hotspot users
export const listHotspotUsers = (id: string) =>
  api.get<Record<string, string>[]>(`/routers/${id}/hotspot/users`).then(r => r.data)

export const listHotspotActive = (id: string) =>
  api.get<Record<string, string>[]>(`/routers/${id}/hotspot/active`).then(r => r.data)

export const listHotspotHosts = (id: string) =>
  api.get<Record<string, string>[]>(`/routers/${id}/hotspot/hosts`).then(r => r.data)

export const listHotspotProfiles = (id: string) =>
  api.get<Record<string, string>[]>(`/routers/${id}/hotspot/profiles`).then(r => r.data)

export const listAddressPools = (id: string) =>
  api.get<Record<string, string>[]>(`/routers/${id}/address-pools`).then(r => r.data)

export const createHotspotUser = (id: string, user: NewHotspotUser) =>
  api.post(`/routers/${id}/hotspot/users`, user).then(r => r.data)

export const updateHotspotUser = (id: string, userID: string, user: UpdateHotspotUser) =>
  api.patch(`/routers/${id}/hotspot/users/${userID}`, user)

export const toggleHotspotUser = (id: string, userID: string, disabled: boolean) =>
  api.post(`/routers/${id}/hotspot/users/${userID}/toggle`, { disabled })

export const getProfileMetas = (id: string) =>
  api.get<Record<string, ProfileMeta>>(`/routers/${id}/hotspot/profile-metas`).then(r => r.data)

export const deleteHotspotUser = (routerId: string, userId: string) =>
  api.delete(`/routers/${routerId}/hotspot/users/${userId}`)

export const disconnectHotspotActive = (routerId: string, sessionId: string) =>
  api.delete(`/routers/${routerId}/hotspot/active/${sessionId}`)

// Hotspot settings
export const getHotspotSettings = (id: string) =>
  api.get<HotspotSettings>(`/routers/${id}/hotspot/settings`).then(r => r.data)

export const putHotspotSettings = (id: string, s: HotspotSettings) =>
  api.put<HotspotSettings>(`/routers/${id}/hotspot/settings`, s).then(r => r.data)

export const uploadLoginPage = (id: string, p: LoginPageSettings & { html: string }) =>
  api.put(`/routers/${id}/hotspot/login-page`, p).then(r => r.data)

export const getLoginPageHTML = (id: string) =>
  api.get<{ html: string }>(`/routers/${id}/hotspot/login-page`).then(r => r.data.html)

// Cleanup scheduler
export const getCleanupScheduler = (id: string) =>
  api.get<CleanupSchedulerStatus>(`/routers/${id}/hotspot/cleanup`).then(r => r.data)

export const putCleanupScheduler = (id: string, enabled: boolean, interval: string) =>
  api.put<CleanupSchedulerStatus>(`/routers/${id}/hotspot/cleanup`, { enabled, interval }).then(r => r.data)

// Hotspot profiles
export const createHotspotProfile = (id: string, p: HotspotProfileParams) =>
  api.post(`/routers/${id}/hotspot/profiles`, p).then(r => r.data)

export const updateHotspotProfile = (id: string, profileId: string, p: HotspotProfileParams) =>
  api.patch(`/routers/${id}/hotspot/profiles/${profileId}`, p)

export const deleteHotspotProfile = (id: string, profileId: string, profileName?: string) =>
  api.delete(`/routers/${id}/hotspot/profiles/${profileId}`, { data: { name: profileName } })

// Hotspot setup
export const hotspotPreflight = (id: string) =>
  api.get<PreflightResult>(`/routers/${id}/hotspot/preflight`).then(r => r.data)

export const setupHotspot = (id: string, body: { lanIface: string; wanIface: string; subnet: string; hotspotName: string }) =>
  api.post<SetupResult>(`/routers/${id}/hotspot/setup`, body).then(r => r.data)

export const teardownHotspot = (id: string) =>
  api.delete<TeardownResult>(`/routers/${id}/hotspot/setup`).then(r => r.data)

// Sales ledger
export const getSalesLedger = (id: string, year: number) =>
  api.get<SaleEntry[]>(`/routers/${id}/sales`, { params: { year } }).then(r => r.data)

// Discovery
export const discoverRouters = () =>
  api.get<DiscoveredDevice[]>('/discover').then(r => r.data)

// Speed test
export const runSpeedTest = (id: string, target = '10') =>
  api.post<Record<string, string>>(`/routers/${id}/speedtest`, { target }).then(r => r.data)
