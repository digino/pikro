import { qrSvg } from "./voucherTemplates/qr";

const escapeHtml = (s: string) =>
  s.replace(/[&<>"]/g, c => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" })[c]!);

/**
 * The A4 card an owner puts on the wall: guests scan it to buy a ticket.
 * In French, like the buy page. No prices, so it stays right when they change.
 */
export async function cloudCardHtml(hotspotName: string, buyLink: string): Promise<string> {
  const qr = await qrSvg(buyLink, 600);
  const name = hotspotName ? `<p class="name">${escapeHtml(hotspotName)}</p>` : "";
  return `<!doctype html>
<html lang="fr"><head><meta charset="utf-8"><title>Acheter un ticket WiFi</title>
<style>
@page { size: A4; margin: 14mm; }
* { box-sizing: border-box; -webkit-print-color-adjust: exact; print-color-adjust: exact; }
body { margin: 0; font-family: system-ui, -apple-system, "Segoe UI", Roboto, sans-serif; color: #111; }
.card { border: 3px solid #111; box-shadow: 8px 8px 0 #111; padding: 12mm 14mm; text-align: center; }
.name { margin: 0 0 4mm; font-size: 20pt; font-weight: 700; color: #3b3a36; }
h1 { margin: 0 0 8mm; font-size: 34pt; line-height: 1.1; font-weight: 800; letter-spacing: -0.02em; }
h1 span { background: #ffbf00; padding: 0 2mm; box-decoration-break: clone; -webkit-box-decoration-break: clone; }
.qr { width: 105mm; height: 105mm; margin: 0 auto 8mm; padding: 4mm; border: 3px solid #111; background: #fff; }
.qr svg { width: 100%; height: 100%; display: block; }
ol { display: inline-block; margin: 0 0 8mm; padding: 0; list-style: none; text-align: left; font-size: 15pt; line-height: 1.5; counter-reset: step; }
li { counter-increment: step; display: flex; gap: 3mm; align-items: baseline; margin-bottom: 2mm; }
li::before { content: counter(step); flex: none; width: 8mm; height: 8mm; line-height: 8mm; text-align: center; font-weight: 800; font-size: 13pt; background: #ffbf00; border: 2px solid #111; }
.foot { margin: 0; font-size: 11pt; color: #5c5a55; }
</style></head>
<body><div class="card">
${name}
<h1><span>Achetez votre ticket WiFi</span><br>par mobile money</h1>
<div class="qr">${qr}</div>
<ol>
<li>Connectez-vous au WiFi</li>
<li>Scannez ce code avec l'appareil photo</li>
<li>Choisissez un forfait et payez</li>
<li>Touchez « Se connecter au WiFi »</li>
</ol>
<p class="foot">Paiement par mobile money avec Pikro</p>
</div></body></html>`;
}
