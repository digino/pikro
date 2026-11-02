// Pikro's public pages. Downloads always go through the Pikro page, which
// holds the current links and counts the clicks.
const REPO = 'https://github.com/digino/pikro'

export const LINKS = {
  page: { en: 'https://sogino.com/pikro', fr: 'https://sogino.com/pikro-fr' },
  repo: REPO,
  issues: `${REPO}/issues/new`,
  releases: `${REPO}/releases`,
}

export function downloadPage(locale: string) {
  return `${locale === 'fr' ? LINKS.page.fr : LINKS.page.en}#download`
}
