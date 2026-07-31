<template>
  <Bar :data="chartData" :options="chartOptions" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  BarElement,
  BarController,
  LinearScale,
  CategoryScale,
  Tooltip,
  type ChartData,
  type ChartOptions,
} from 'chart.js'
import { formatCompactAmount } from '@/utils/currencies'

ChartJS.register(BarElement, BarController, LinearScale, CategoryScale, Tooltip)

export interface SalesBarPoint {
  label: string
  count: number
  revenue: number
}

const props = defineProps<{ points: SalesBarPoint[]; currency?: string }>()

function cssVar(name: string) {
  return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
}

const COUNT_COLOR = '#22d3ee'   // cyan-400: vouchers generated
const REVENUE_COLOR = '#f59e0b' // amber-400: revenue

const chartData = computed<ChartData<'bar'>>(() => ({
  labels: props.points.map(p => p.label),
  datasets: [
    {
      label: 'Vouchers generated',
      data: props.points.map(p => p.count),
      backgroundColor: `${COUNT_COLOR}cc`,
      borderRadius: 3,
      yAxisID: 'yCount',
    },
    {
      label: 'Revenue',
      data: props.points.map(p => p.revenue),
      backgroundColor: `${REVENUE_COLOR}cc`,
      borderRadius: 3,
      yAxisID: 'yRevenue',
    },
  ],
}))

const chartOptions = computed<ChartOptions<'bar'>>(() => {
  const border = cssVar('--color-border')
  const textMuted = cssVar('--color-text-muted')
  const textBody = cssVar('--color-text-primary')
  const surface = cssVar('--color-surface')

  return {
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
    plugins: {
      legend: {
        display: true,
        position: 'top',
        align: 'end',
        labels: { color: textMuted, boxWidth: 10, boxHeight: 10, font: { size: 11 } },
      },
      tooltip: {
        backgroundColor: surface,
        borderColor: border,
        borderWidth: 1,
        titleColor: textMuted,
        bodyColor: textBody,
        padding: 8,
        callbacks: {
          label: ctx => {
            if (ctx.dataset.label === 'Revenue') {
              return `Revenue: ${formatCompactAmount(ctx.parsed.y ?? 0, props.currency)}`
            }
            return `Vouchers: ${ctx.parsed.y ?? 0}`
          },
        },
      },
    },
    scales: {
      x: {
        grid: { display: false },
        border: { display: false },
        ticks: { color: textMuted, font: { size: 10 }, maxRotation: 0, autoSkip: true, maxTicksLimit: 10 },
      },
      yCount: {
        position: 'left',
        min: 0,
        grid: { color: `color-mix(in srgb, ${border} 60%, transparent)`, drawTicks: false },
        border: { display: false },
        ticks: { color: textMuted, font: { size: 10 }, maxTicksLimit: 5, precision: 0 },
      },
      yRevenue: {
        position: 'right',
        min: 0,
        grid: { display: false },
        border: { display: false },
        ticks: {
          color: textMuted,
          font: { size: 10 },
          maxTicksLimit: 5,
          callback: v => formatCompactAmount(Number(v)),
        },
      },
    },
  }
})
</script>
