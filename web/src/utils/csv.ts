// Escapes a value for CSV: wraps in quotes and doubles any internal quotes
// if the value contains a comma, quote, or newline.
function csvCell(value: string | number): string {
  const s = String(value)
  if (/[",\n]/.test(s)) return `"${s.replace(/"/g, '""')}"`
  return s
}

export function toCsv(headers: string[], rows: (string | number)[][]): string {
  const lines = [headers, ...rows].map(row => row.map(csvCell).join(','))
  return lines.join('\r\n')
}

export function downloadCsv(filename: string, headers: string[], rows: (string | number)[][]) {
  const csv = toCsv(headers, rows)
  const blob = new Blob(['﻿' + csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}
