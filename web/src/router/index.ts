import { createRouter, createWebHistory } from 'vue-router'
import { useRoutersStore } from '@/stores/routers'
import Home from '@/views/Home.vue'
import Dashboard from '@/views/Dashboard.vue'
import HotspotLayout from '@/views/Hotspot.vue'
import HotspotSetup from '@/views/hotspot/Setup.vue'
import HotspotUsers from '@/views/hotspot/Users.vue'
import HotspotVouchers from '@/views/hotspot/Vouchers.vue'
import HotspotProfiles from '@/views/hotspot/Profiles.vue'
import HotspotSettings from '@/views/hotspot/Settings.vue'
import HotspotLogs from '@/views/hotspot/Logs.vue'
import HotspotReports from '@/views/hotspot/Reports.vue'
import NetworkLayout from '@/views/Network.vue'
import NetworkMonitor from '@/views/network/Monitor.vue'
import NetworkHosts from '@/views/network/Hosts.vue'
import NetworkDhcp from '@/views/network/Dhcp.vue'
import SpeedTest from '@/views/SpeedTest.vue'
import Routers from '@/views/Routers.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: Home,
      beforeEnter: () => {
        const store = useRoutersStore()
        if (store.routers.length === 0) return '/routers'
        if (store.activeId) return '/dashboard'
        return true
      },
    },
    { path: '/dashboard', component: Dashboard },
    {
      path: '/hotspot',
      component: HotspotLayout,
      redirect: '/hotspot/setup',
      children: [
        { path: 'setup',    component: HotspotSetup },
        { path: 'users',    component: HotspotUsers },
        { path: 'vouchers', component: HotspotVouchers },
        { path: 'profiles', component: HotspotProfiles },
        { path: 'settings', component: HotspotSettings },
        { path: 'logs',     component: HotspotLogs },
        { path: 'reports',  component: HotspotReports },
      ],
    },
    {
      path: '/network',
      component: NetworkLayout,
      redirect: '/network/monitor',
      children: [
        { path: 'monitor', component: NetworkMonitor },
        { path: 'hosts',   component: NetworkHosts },
        { path: 'dhcp',    component: NetworkDhcp },
      ],
    },
    { path: '/speedtest', component: SpeedTest },
    { path: '/routers',   component: Routers },
  ],
})
