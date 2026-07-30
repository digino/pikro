<template>
  <Line :data="chartData" :options="chartOptions" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  LineElement,
  PointElement,
  LinearScale,
  CategoryScale,
  Filler,
  Tooltip,
  type ChartData,
  type ChartOptions,
} from 'chart.js'

ChartJS.register(LineElement, PointElement, LinearScale, CategoryScale, Filler, Tooltip)

interface BwPoint { down: number; up: number }

const props = defineProps<{ history: BwPoint[] }>()

function cssVar(name: string) {
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
}

// Fixed semantic colors — match the legend dots in the bandwidth card header.
const DOWN_COLOR = '#22d3ee' // cyan-400: download
const UP_COLOR   = '#f59e0b' // amber-400: upload

const N = computed(() => props.history.length)

const chartData = computed<ChartData<'line'>>(() => ({
  labels: Array.from({ length: N.value }, (_, i) => {
    const secondsAgo = (N.value - 1 - i) * 3
    if (secondsAgo === 0) return 'now'
    if (secondsAgo % 15 === 0) return `${secondsAgo}s`
    return ''
  }),
  datasets: [
    {
      label: 'Download',
      data: props.history.map(p => p.down),
      borderColor: DOWN_COLOR,
      backgroundColor: `${DOWN_COLOR}20`,
      borderWidth: 2,
      pointRadius: 0,
      fill: true,
      tension: 0.4,
    },
    {
      label: 'Upload',
      data: props.history.map(p => p.up),
      borderColor: UP_COLOR,
      backgroundColor: `${UP_COLOR}15`,
      borderWidth: 2,
      pointRadius: 0,
      fill: true,
      tension: 0.4,
    },
  ],
}))

const chartOptions = computed<ChartOptions<'line'>>(() => {
  const allVals = props.history.flatMap(p => [p.down, p.up])
  const peak = Math.max(...allVals, 0.01)
  const exp  = Math.pow(10, Math.floor(Math.log10(peak)))
  const ceil = Math.ceil(peak / exp) * exp

  const border    = cssVar('--color-border')
  const textMuted = cssVar('--color-text-muted')
  const textBody  = cssVar('--color-text-primary')
  const surface   = cssVar('--color-surface')

  return {
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
    interaction: { mode: 'index', intersect: false },
    plugins: {
      legend: { display: false },
      tooltip: {
        backgroundColor: surface,
        borderColor: border,
        borderWidth: 1,
        titleColor: textMuted,
        bodyColor: textBody,
        padding: 8,
        callbacks: {
          label: ctx => `${ctx.dataset.label}: ${(ctx.parsed.y ?? 0).toFixed(2)} Mbps`,
        },
      },
    },
    scales: {
      x: {
        grid: { color: `color-mix(in srgb, ${border} 60%, transparent)`, drawTicks: false },
        border: { display: false },
        ticks: {
          color: textMuted,
          font: { family: 'monospace', size: 10 },
          maxRotation: 0,
          autoSkip: false,
        },
      },
      y: {
        min: 0,
        max: ceil,
        grid: { color: `color-mix(in srgb, ${border} 60%, transparent)`, drawTicks: false },
        border: { display: false },
        ticks: {
          color: textMuted,
          font: { family: 'monospace', size: 10 },
          maxTicksLimit: 4,
          callback: v => `${v} Mbps`,
        },
      },
    },
  }
})
</script>
