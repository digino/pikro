// Hand-inked multicolour icons (ink + white + amber). Render as <img>, not
// inline: each SVG carries its own filter id. Heroicons stay for small glyphs.
export type InkIcon =
  | 'dashboard' | 'routers' | 'users' | 'profiles' | 'vouchers' | 'reports'
  | 'settings' | 'logs' | 'monitor' | 'hosts' | 'dhcp' | 'speed-test'

export function inkIcon(name: InkIcon): string {
  return new URL(`../assets/icons/ink/${name}.svg`, import.meta.url).href
}
