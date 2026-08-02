import type { VoucherTemplate } from "./types";
import { qrSvg } from "./qr";

const QR_SIZE_PX = 60;

// Business — same idea as Modern but roomier, with a QR code alongside the
// credentials. The QR encodes a login link: if the hotspot profile uses
// http-pap auth (loginUrlSupportsCredentials), the link includes the
// username/password so scanning logs the user straight in; otherwise it
// just opens the login page and the user types the credentials themselves.
const render: VoucherTemplate["render"] = async (items, opts) => {
  const { businessName, loginUrl, loginUrlSupportsCredentials } = opts;

  const cards = await Promise.all(
    items.map(async ({ name, password, validity, price }, i) => {
      let qr = "";
      if (loginUrl) {
        const target = loginUrlSupportsCredentials
          ? `${loginUrl}${loginUrl.includes("?") ? "&" : "?"}username=${encodeURIComponent(name)}&password=${encodeURIComponent(password)}`
          : loginUrl;
        qr = await qrSvg(target, QR_SIZE_PX);
      }

      return `
        <div class="v">
          <div class="left">
            ${businessName ? `<div class="biz">${businessName}</div>` : ""}
            <div class="row"><span class="lbl">Username</span><span class="val">${name}</span></div>
            <div class="row"><span class="lbl">Password</span><span class="val">${password}</span></div>
            ${validity || price ? `<div class="row sub"><span class="lbl">${validity ? "Valid" : "Price"}</span><span class="val-sm">${[validity, price].filter(Boolean).join(" · ")}</span></div>` : ""}
            <div class="num">#${i + 1}</div>
          </div>
          ${qr ? `<div class="right">${qr}</div>` : ""}
        </div>`;
    }),
  );

  const css = `
    *{box-sizing:border-box;margin:0;padding:0}
    body{font-family:ui-sans-serif,system-ui,sans-serif;background:#fff;color:#000;padding:6mm}
    .grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(40mm,1fr));gap:2.5mm}
    .v{border:1.25pt solid #000;border-radius:2px;padding:2mm;display:flex;align-items:center;gap:1.5mm;page-break-inside:avoid}
    .left{flex:1;min-width:0;display:flex;flex-direction:column;gap:0.6mm}
    .right{shrink:0;display:flex;align-items:center;justify-content:center}
    .biz{font-size:6pt;font-weight:700;text-transform:uppercase;letter-spacing:.02em;border-bottom:0.75pt solid #000;padding-bottom:0.6mm;margin-bottom:0.2mm}
    .row{display:flex;justify-content:space-between;align-items:baseline;gap:1mm}
    .lbl{font-size:5pt;text-transform:uppercase;letter-spacing:.02em;color:#333}
    .val{font-size:8.5pt;font-weight:700;font-family:ui-monospace,monospace}
    .sub{border-top:0.75pt solid #000;padding-top:0.6mm;margin-top:0.2mm}
    .val-sm{font-size:6pt;font-weight:600}
    .num{font-size:5pt;color:#666;margin-top:auto}
    @media print{body{padding:4mm}.grid{gap:2mm}}
  `;

  return `<!DOCTYPE html><html><head><meta charset="UTF-8"/><title>Vouchers — Business</title><style>${css}</style></head><body><div class="grid">${cards.join("")}</div></body></html>`;
};

export const BUSINESS_TEMPLATE: VoucherTemplate = {
  key: "business",
  label: "Business",
  description: "Cards with a scannable QR login code, ~40 per page",
  render,
};
