import type { VoucherTemplate } from "./types";
import { qrSvg } from "./qr";

const QR_SIZE_PX = 90;

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
    body{font-family:ui-sans-serif,system-ui,sans-serif;background:#fff;color:#000;padding:8mm}
    .grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(70mm,1fr));gap:4mm}
    .v{border:1.25pt solid #000;border-radius:3px;padding:4mm;display:flex;align-items:center;gap:3mm;page-break-inside:avoid}
    .left{flex:1;min-width:0;display:flex;flex-direction:column;gap:1.5mm}
    .right{shrink:0;display:flex;align-items:center;justify-content:center}
    .biz{font-size:8pt;font-weight:700;text-transform:uppercase;letter-spacing:.02em;border-bottom:0.75pt solid #000;padding-bottom:1.5mm;margin-bottom:0.5mm}
    .row{display:flex;justify-content:space-between;align-items:baseline;gap:2mm}
    .lbl{font-size:6.5pt;text-transform:uppercase;letter-spacing:.03em;color:#333}
    .val{font-size:12pt;font-weight:700;font-family:ui-monospace,monospace}
    .sub{border-top:0.75pt solid #000;padding-top:1.5mm;margin-top:0.5mm}
    .val-sm{font-size:8pt;font-weight:600}
    .num{font-size:6pt;color:#666;margin-top:auto}
    @media print{body{padding:5mm}.grid{gap:3mm}}
  `;

  return `<!DOCTYPE html><html><head><meta charset="UTF-8"/><title>Vouchers — Business</title><style>${css}</style></head><body><div class="grid">${cards.join("")}</div></body></html>`;
};

export const BUSINESS_TEMPLATE: VoucherTemplate = {
  key: "business",
  label: "Business",
  description: "Roomier cards with a scannable QR login code",
  render,
};
