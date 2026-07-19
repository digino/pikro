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

const props = defineProps<{
  history: BwPoint[]
}>()

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
      borderColor: 'rgba(250,250,250,0.9)',
      backgroundColor: 'rgba(250,250,250,0.08)',
      borderWidth: 2,
      pointRadius: 0,
      fill: true,
      tension: 0.4,
    },
    {
      label: 'Upload',
      data: props.history.map(p => p.up),
      borderColor: 'rgba(82,82,91,0.9)',
      backgroundColor: 'transparent',
      borderWidth: 1.5,
      pointRadius: 0,
      fill: false,
      tension: 0.4,
    },
  ],
}))

const chartOptions = computed<ChartOptions<'line'>>(() => {
  const allVals = props.history.flatMap(p => [p.down, p.up])
  const peak = Math.max(...allVals, 0.01)
  const exp = Math.pow(10, Math.floor(Math.log10(peak)))
  const ceil = Math.ceil(peak / exp) * exp

  return {
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
    interaction: { mode: 'index', intersect: false },
    plugins: {
      legend: { display: false },
      tooltip: {
        backgroundColor: 'rgba(14,15,17,0.95)',
        borderColor: 'rgba(255,255,255,0.1)',
        borderWidth: 1,
        titleColor: 'rgba(255,255,255,0.5)',
        bodyColor: 'rgba(255,255,255,0.9)',
        padding: 8,
        callbacks: {
          label: ctx => `${ctx.dataset.label}: ${(ctx.parsed.y ?? 0).toFixed(2)} Mbps`,
        },
      },
    },
    scales: {
      x: {
        grid: { color: 'rgba(255,255,255,0.05)', drawTicks: false },
        border: { display: false },
        ticks: {
          color: 'rgba(255,255,255,0.3)',
          font: { family: 'monospace', size: 10 },
          maxRotation: 0,
          autoSkip: false,
        },
      },
      y: {
        min: 0,
        max: ceil,
        grid: { color: 'rgba(255,255,255,0.05)', drawTicks: false },
        border: { display: false },
        ticks: {
          color: 'rgba(255,255,255,0.3)',
          font: { family: 'monospace', size: 10 },
          maxTicksLimit: 4,
          callback: v => `${v} Mbps`,
        },
      },
    },
  }
})
</script>
