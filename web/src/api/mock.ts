// Mock data for UI development — activated by VITE_MOCK=true.
// All functions match the signature of api/index.ts exports exactly.

export type * from './types'
import type {
  RouterProfile, NewRouterProfile,
  NewHotspotUser, UpdateHotspotUser,
  ProfileMeta, HotspotSettings, LoginPageSettings,
  CleanupSchedulerStatus, HotspotProfileParams,
  PreflightResult, SetupResult, TeardownResult,
  DiscoveredDevice,
} from './types'

const delay = (ms = 300) => new Promise(r => setTimeout(r, ms))

// ── Fixed IDs ──────────────────────────────────────────────────────────────
const ROUTER_ID = 'mock-router-1'
const now = Math.floor(Date.now() / 1000)

// ── State (mutations persist within the session) ───────────────────────────
let routers: RouterProfile[] = [
  { id: ROUTER_ID, name: 'CyberCafé Abidjan', host: '192.168.88.1', port: 8728, username: 'admin', useTls: false },
]

let users: Record<string, string>[] = [
  { '.id': '*1', name: 'abc123', password: 'abc123', profile: 'Standard 1h',  'limit-uptime': '1h',  'limit-bytes-total': '',           comment: `exp:${now + 3200}`,  uptime: '00:22:10', disabled: 'false' },
  { '.id': '*2', name: 'xyz789', password: 'xyz789', profile: 'Standard 1h',  'limit-uptime': '1h',  'limit-bytes-total': '',           comment: `exp:${now - 100}`,   uptime: '01:00:05', disabled: 'false' },
  { '.id': '*3', name: 'pqr456', password: 'pqr456', profile: 'Data 500MB',   'limit-uptime': '',    'limit-bytes-total': '524288000',  comment: `exp:${now + 72000}`, uptime: '00:00:00', disabled: 'false' },
  { '.id': '*4', name: 'lmn000', password: 'lmn000', profile: 'Standard 1h',  'limit-uptime': '1h',  'limit-bytes-total': '',           comment: '',                   uptime: '00:00:00', disabled: 'false' },
  { '.id': '*5', name: 'zzz999', password: 'zzz999', profile: 'VIP 24h',      'limit-uptime': '24h', 'limit-bytes-total': '1073741824', comment: `exp:${now + 82000}`, uptime: '03:14:00', disabled: 'true'  },
  { '.id': '*6', name: 'ttt111', password: 'ttt111', profile: 'Data 500MB',   'limit-uptime': '',    'limit-bytes-total': '524288000',  comment: `exp:${now + 3600}`,  uptime: '00:45:00', disabled: 'false' },
  { '.id': '*7', name: 'mmm222', password: 'mmm222', profile: 'Standard 1h',  'limit-uptime': '1h',  'limit-bytes-total': '',           comment: '',                   uptime: '00:00:00', disabled: 'false' },
  { '.id': '*8', name: 'bbb333', password: 'bbb333', profile: 'VIP 24h',      'limit-uptime': '24h', 'limit-bytes-total': '1073741824', comment: `exp:${now + 50000}`, uptime: '12:00:00', disabled: 'false' },
]

const profiles: Record<string, string>[] = [
  { '.id': '*1', name: 'Standard 1h', 'rate-limit': '2M/4M', 'shared-users': '1', 'address-pool': 'hotspot-pool' },
  { '.id': '*2', name: 'Data 500MB',  'rate-limit': '4M/8M', 'shared-users': '1', 'address-pool': 'hotspot-pool' },
  { '.id': '*3', name: 'VIP 24h',     'rate-limit': '10M/10M','shared-users': '2', 'address-pool': 'hotspot-pool' },
]

const profileMetas: Record<string, ProfileMeta> = {
  'Standard 1h': { validity: '1h',  price: '200' },
  'Data 500MB':  { validity: '3d',  price: '500' },
  'VIP 24h':     { validity: '24h', price: '1000' },
}

let active: Record<string, string>[] = [
  { '.id': '*A', user: 'abc123', address: '192.168.99.10', 'mac-address': 'AA:BB:CC:DD:EE:01', uptime: '00:22:10', 'session-time-left': '00:37:50', 'bytes-in': '15728640', 'bytes-out': '4194304' },
  { '.id': '*B', user: 'ttt111', address: '192.168.99.15', 'mac-address': 'AA:BB:CC:DD:EE:02', uptime: '00:45:00', 'session-time-left': '00:15:00', 'bytes-in': '52428800', 'bytes-out': '10485760' },
  { '.id': '*C', user: 'bbb333', address: '192.168.99.22', 'mac-address': 'AA:BB:CC:DD:EE:03', uptime: '12:00:00', 'session-time-left': '12:00:00', 'bytes-in': '314572800','bytes-out': '104857600' },
]

let hotspotSettings: HotspotSettings = {
  hotspotName: 'Wifi Zone Nord-Ouest',
  dnsName: 'nordouest.spot',
  currency: 'XOF',
  loginPage: { title: 'Wifi Zone Nord-Ouest', subtitle: 'Bienvenue — connectez-vous pour accéder à internet', template: 'minimal' },
}

let cleanup: CleanupSchedulerStatus = { installed: true, interval: '1d' }
let liveLoginPageHTML = '<!DOCTYPE html><html><body><p>No login page uploaded yet.</p></body></html>'

let nextId = 9

// ── Helpers ────────────────────────────────────────────────────────────────
function ok<T>(data: T): Promise<T> { return delay().then(() => data) }

// ── Exports ────────────────────────────────────────────────────────────────
export const getAppVersion = () => ok('dev-mock')

export const listRouters = () => ok([...routers])
export const addRouter = (p: NewRouterProfile) => {
  const r: RouterProfile = { id: `mock-${Date.now()}`, ...p }
  routers.push(r)
  return ok({ id: r.id })
}
export const deleteRouter = (id: string) => { routers = routers.filter(r => r.id !== id); return ok(undefined) }
export const testRouter = (_id: string) => ok({ ok: true })

export const getSystemResource = (_id: string) => ok({
  'cpu-load': '12',
  'free-memory': '38797312',
  'total-memory': '67108864',
  'uptime': '5d02:14:33',
  'version': '7.19.5 (stable)',
  'board-name': 'RB941-2nD',
  'architecture-name': 'mipsbe',
})

export const listHotspotUsers = (_id: string) => ok([...users])
export const listHotspotActive = (_id: string) => ok([...active])
export const listHotspotProfiles = (_id: string) => ok([...profiles])

export const createHotspotUser = (_id: string, u: NewHotspotUser) => {
  const user: Record<string, string> = {
    '.id': `*${nextId++}`,
    name: u.name,
    password: u.password,
    profile: u.profile || 'Standard 1h',
    'limit-uptime': u.limitUptime,
    'limit-bytes-total': u.limitBytesTotal,
    comment: u.comment,
    uptime: '00:00:00',
    disabled: 'false',
  }
  users.push(user)
  return ok(user)
}

export const updateHotspotUser = (_id: string, userID: string, u: UpdateHotspotUser) => {
  const idx = users.findIndex(x => x['.id'] === userID)
  if (idx !== -1) {
    if (u.password)        users[idx].password = u.password
    if (u.profile)         users[idx].profile = u.profile
    if (u.limitUptime !== undefined)    users[idx]['limit-uptime'] = u.limitUptime
    if (u.limitBytesTotal !== undefined) users[idx]['limit-bytes-total'] = u.limitBytesTotal
  }
  return ok(undefined)
}

export const toggleHotspotUser = (_id: string, userID: string, disabled: boolean) => {
  const idx = users.findIndex(x => x['.id'] === userID)
  if (idx !== -1) users[idx].disabled = String(disabled)
  return ok(undefined)
}

export const getProfileMetas = (_id: string) => ok({ ...profileMetas })

export const deleteHotspotUser = (_routerId: string, userId: string) => {
  const deleted = users.find(u => u['.id'] === userId)
  users = users.filter(u => u['.id'] !== userId)
  if (deleted) active = active.filter(s => s.user !== deleted.name)
  return ok(undefined)
}

export const disconnectHotspotActive = (_routerId: string, sessionId: string) => {
  active = active.filter(s => s['.id'] !== sessionId)
  return ok(undefined)
}

export const getHotspotSettings = (_id: string) => ok({ ...hotspotSettings })
export const putHotspotSettings = (_id: string, s: HotspotSettings) => { hotspotSettings = s; return ok(s) }
export const uploadLoginPage = (_id: string, p: LoginPageSettings & { html: string }) => {
  liveLoginPageHTML = p.html
  return ok(undefined)
}
export const getLoginPageHTML = (_id: string) => ok(liveLoginPageHTML)

export const getCleanupScheduler = (_id: string) => ok({ ...cleanup })
export const putCleanupScheduler = (_id: string, enabled: boolean, interval: string) => {
  cleanup = { installed: enabled, interval }
  return ok(cleanup)
}

export const createHotspotProfile = (_id: string, _p: HotspotProfileParams) => ok({})
export const updateHotspotProfile = (_id: string, _profileId: string, _p: HotspotProfileParams) => ok(undefined)
export const deleteHotspotProfile = (_id: string, _profileId: string, _name?: string) => ok(undefined)

export const hotspotPreflight = (_id: string): Promise<PreflightResult> => ok({
  interfaces: [
    { name: 'ether1', type: 'ether', running: true, comment: 'WAN' },
    { name: 'ether2', type: 'ether', running: true, comment: 'LAN' },
    { name: 'wlan1',  type: 'wlan',  running: true, comment: '' },
  ],
  hotspotExists: true,
  hotspotOnIface: 'wlan1',
  hotspotName: 'Pikro-HS',
  hotspotProfile: 'hsprof1',
  hotspotDnsName: 'hotspot.local',
  hotspotAddressPool: 'hotspot-pool',
})

export const setupHotspot = (_id: string, _body: any): Promise<SetupResult> => ok({
  steps: [
    { name: 'Create address pool', ok: true },
    { name: 'Enable hotspot', ok: true },
    { name: 'Upload login page', ok: true },
  ],
  success: true,
})

export const teardownHotspot = (_id: string): Promise<TeardownResult> => ok({ steps: [] })

export const discoverRouters = (): Promise<DiscoveredDevice[]> => ok([
  { ip: '192.168.88.1', mac: 'AA:BB:CC:DD:EE:FF', identity: 'CyberCafé-Router', version: '7.19.5', board: 'RB941-2nD', platform: 'MikroTik', uptime: '5d02:14:33', iface: 'en0' },
])

export const runSpeedTest = (_id: string, _target?: string) =>
  delay(2000).then(() => ({ 'rx-speed': '10485760', 'tx-speed': '0', duration: '8s120ms', 'test-url': 'http://speedtest.tele2.net/10MB.zip', 'file-size': '10485760' }))
