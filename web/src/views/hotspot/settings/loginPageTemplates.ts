import i18n from "@/locales";

// Minified Paul Johnston reference MD5 impl (the same one MikroTik's own
// stock hotspot templates use)
const MD5_JS =
  'function safe_add(f,i){var n=(65535&f)+(65535&i);return(f>>16)+(i>>16)+(n>>16)<<16|65535&n}function rol(f,i){return f<<i|f>>>32-i}function cmn(f,i,n,h,r,g){return safe_add(rol(safe_add(safe_add(i,f),safe_add(h,g)),r),n)}function ff(f,i,n,h,r,g,e){return cmn(i&n|~i&h,f,i,r,g,e)}function gg(f,i,n,h,r,g,e){return cmn(i&h|n&~h,f,i,r,g,e)}function hh(f,i,n,h,r,g,e){return cmn(i^n^h,f,i,r,g,e)}function ii(f,i,n,h,r,g,e){return cmn(n^(i|~h),f,i,r,g,e)}function coreMD5(f){for(var i=1732584193,n=-271733879,h=-1732584194,r=271733878,g=0;g<f.length;g+=16){var e=i,t=n,a=h,c=r;i=ff(i,n,h,r,f[g+0],7,-680876936),r=ff(r,i,n,h,f[g+1],12,-389564586),h=ff(h,r,i,n,f[g+2],17,606105819),n=ff(n,h,r,i,f[g+3],22,-1044525330),i=ff(i,n,h,r,f[g+4],7,-176418897),r=ff(r,i,n,h,f[g+5],12,1200080426),h=ff(h,r,i,n,f[g+6],17,-1473231341),n=ff(n,h,r,i,f[g+7],22,-45705983),i=ff(i,n,h,r,f[g+8],7,1770035416),r=ff(r,i,n,h,f[g+9],12,-1958414417),h=ff(h,r,i,n,f[g+10],17,-42063),n=ff(n,h,r,i,f[g+11],22,-1990404162),i=ff(i,n,h,r,f[g+12],7,1804603682),r=ff(r,i,n,h,f[g+13],12,-40341101),h=ff(h,r,i,n,f[g+14],17,-1502002290),i=gg(i,n=ff(n,h,r,i,f[g+15],22,1236535329),h,r,f[g+1],5,-165796510),r=gg(r,i,n,h,f[g+6],9,-1069501632),h=gg(h,r,i,n,f[g+11],14,643717713),n=gg(n,h,r,i,f[g+0],20,-373897302),i=gg(i,n,h,r,f[g+5],5,-701558691),r=gg(r,i,n,h,f[g+10],9,38016083),h=gg(h,r,i,n,f[g+15],14,-660478335),n=gg(n,h,r,i,f[g+4],20,-405537848),i=gg(i,n,h,r,f[g+9],5,568446438),r=gg(r,i,n,h,f[g+14],9,-1019803690),h=gg(h,r,i,n,f[g+3],14,-187363961),n=gg(n,h,r,i,f[g+8],20,1163531501),i=gg(i,n,h,r,f[g+13],5,-1444681467),r=gg(r,i,n,h,f[g+2],9,-51403784),h=gg(h,r,i,n,f[g+7],14,1735328473),i=hh(i,n=gg(n,h,r,i,f[g+12],20,-1926607734),h,r,f[g+5],4,-378558),r=hh(r,i,n,h,f[g+8],11,-2022574463),h=hh(h,r,i,n,f[g+11],16,1839030562),n=hh(n,h,r,i,f[g+14],23,-35309556),i=hh(i,n,h,r,f[g+1],4,-1530992060),r=hh(r,i,n,h,f[g+4],11,1272893353),h=hh(h,r,i,n,f[g+7],16,-155497632),n=hh(n,h,r,i,f[g+10],23,-1094730640),i=hh(i,n,h,r,f[g+13],4,681279174),r=hh(r,i,n,h,f[g+0],11,-358537222),h=hh(h,r,i,n,f[g+3],16,-722521979),n=hh(n,h,r,i,f[g+6],23,76029189),i=hh(i,n,h,r,f[g+9],4,-640364487),r=hh(r,i,n,h,f[g+12],11,-421815835),h=hh(h,r,i,n,f[g+15],16,530742520),i=ii(i,n=hh(n,h,r,i,f[g+2],23,-995338651),h,r,f[g+0],6,-198630844),r=ii(r,i,n,h,f[g+7],10,1126891415),h=ii(h,r,i,n,f[g+14],15,-1416354905),n=ii(n,h,r,i,f[g+5],21,-57434055),i=ii(i,n,h,r,f[g+12],6,1700485571),r=ii(r,i,n,h,f[g+3],10,-1894986606),h=ii(h,r,i,n,f[g+10],15,-1051523),n=ii(n,h,r,i,f[g+1],21,-2054922799),i=ii(i,n,h,r,f[g+8],6,1873313359),r=ii(r,i,n,h,f[g+15],10,-30611744),h=ii(h,r,i,n,f[g+6],15,-1560198380),n=ii(n,h,r,i,f[g+13],21,1309151649),i=ii(i,n,h,r,f[g+4],6,-145523070),r=ii(r,i,n,h,f[g+11],10,-1120210379),h=ii(h,r,i,n,f[g+2],15,718787259),n=ii(n,h,r,i,f[g+9],21,-343485551),i=safe_add(i,e),n=safe_add(n,t),h=safe_add(h,a),r=safe_add(r,c)}return[i,n,h,r]}function binl2hex(f){for(var i="0123456789abcdef",n="",h=0;h<4*f.length;h++)n+=i.charAt(f[h>>2]>>h%4*8+4&15)+i.charAt(f[h>>2]>>h%4*8&15);return n}function str2binl(f){for(var i=1+(f.length+8>>6),n=new Array(16*i),h=0;h<16*i;h++)n[h]=0;for(h=0;h<f.length;h++)n[h>>2]|=(255&f.charCodeAt(h))<<h%4*8;return n[h>>2]|=128<<h%4*8,n[16*i-2]=8*f.length,n}function hexMD5(f){return binl2hex(coreMD5(str2binl(f)))}';

// The hash script is the router's md5.js, not inlined: RouterOS 6 refuses
// to store a page over 4 KB through the API.
const CHAP_SCRIPT = [
  '<script src="/md5.js"><\/script>',
  "<script>",
  "function doLogin(){document.sendin.username.value=document.login.username.value;document.sendin.password.value=hexMD5('$(chap-id)'+document.login.password.value+'$(chap-challenge)');document.sendin.submit();return false}",
  "<\/script>",
].join("\n");

const CHAP_HIDDEN_FORM = [
  "$(if chap-id)",
  '<form name="sendin" action="$(link-login-only)" method="post" style="display:none">',
  '  <input type="hidden" name="username" />',
  '  <input type="hidden" name="password" />',
  '  <input type="hidden" name="dst" value="$(link-orig)" />',
  '  <input type="hidden" name="popup" value="true" />',
  "<\/form>",
  CHAP_SCRIPT,
  "$(endif)",
].join("\n");

// ── Guest pages ──────────────────────────────────────────────────────────────
// One Ink & Amber style for login, status, logout and redirect. The operator
// picks how guests log in and the page language; error stays the fixed
// Go default. Text comes from the `guestPage` locale keys, in the page's
// language rather than the language Pikro itself is shown in.

export type LoginMode = "code" | "password";
export type GuestLanguage = "en" | "fr";
export type GuestBackground = "light" | "tint" | "dark";

export const DEFAULT_ACCENT = "#111111";
export const ACCENT_PRESETS = [
  { key: "ink", hex: DEFAULT_ACCENT },
  { key: "amber", hex: "#ffbf00" },
  { key: "ocean", hex: "#1d4ed8" },
  { key: "forest", hex: "#15803d" },
  { key: "berry", hex: "#be185d" },
  { key: "sunset", hex: "#ea580c" },
] as const;
export const CONTACT_MAX = 60;

export interface GuestPageOptions {
  title: string;
  subtitle: string;
  mode: LoginMode;
  language: GuestLanguage;
  prices: { name: string; price: string }[];
  /** Button colour, "#rrggbb"; Pikro's ink when absent or malformed. */
  accent?: string;
  background?: GuestBackground;
  /** Short line under the login button, e.g. a phone number. */
  contact?: string;
}

export interface GuestPages {
  login: string;
  status: string;
  logout: string;
  redirect: string;
  alogin: string;
  css: string;
  md5Js: string;
}

function esc(s: string) {
  return s
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}

// Uploaded as the router's css/style.css and linked, for the same 4 KB reason.
// Colours come from the :root variables guestThemeVars writes in front.
const GUEST_CSS = [
  "*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }",
  'body { font-family: system-ui, -apple-system, "Segoe UI", Roboto, sans-serif; background: var(--bg); color: var(--ink); min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 24px 20px; }',
  ".wrap { width: 100%; max-width: 340px; display: flex; flex-direction: column; align-items: center; gap: 14px; text-align: center; }",
  "h1 { font-size: 26px; font-weight: 800; line-height: 1.15; letter-spacing: -.01em; }",
  ".sub { font-size: 13px; color: var(--muted); }",
  "form { width: 100%; display: flex; flex-direction: column; gap: 10px; margin-top: 10px; }",
  "label { text-align: left; font-size: 12px; font-weight: 600; color: var(--muted); }",
  "input { width: 100%; height: 52px; padding: 0 14px; border: 2px solid var(--ink); border-radius: 12px; background: var(--card); color: var(--ink); font-size: 16px; font-weight: 500; font-family: inherit; -webkit-appearance: none; }",
  'input.code { text-align: center; font: 600 20px ui-monospace, "SF Mono", Menlo, Consolas, monospace; letter-spacing: .12em; }',
  "input:focus { outline: 3px solid var(--focus); outline-offset: 1px; }",
  ".btn { display: flex; align-items: center; justify-content: center; width: 100%; height: 50px; border: 0; border-radius: 12px; background: var(--accent); color: var(--on); font-size: 15px; font-weight: 700; font-family: inherit; cursor: pointer; text-decoration: none; }",
  ".hint { font-size: 12px; color: var(--muted); }",
  ".hint a { color: var(--ink); font-weight: 600; }",
  ".contact { font-size: 13px; font-weight: 600; }",
  ".err { width: 100%; padding: 10px 12px; border: 2px solid #b42318; border-radius: 12px; background: #fdecea; color: #b42318; font-size: 13px; font-weight: 600; }",
  ".box { width: 100%; border: 2px solid var(--ink); border-radius: 12px; background: var(--card); padding: 12px 14px; text-align: left; }",
  ".box h2 { font-size: 12px; font-weight: 600; color: var(--muted); margin-bottom: 6px; }",
  ".row { display: flex; justify-content: space-between; gap: 12px; padding: 5px 0; font-size: 14px; }",
  ".row + .row { border-top: 1px dashed var(--rule); }",
  ".row b { font-variant-numeric: tabular-nums; white-space: nowrap; }",
  ".ok { display: inline-flex; align-items: center; gap: 6px; font-size: 14px; font-weight: 700; color: var(--ok); }",
  '.ok::before { content: ""; width: 9px; height: 9px; border-radius: 50%; background: var(--ok); }',
  ".foot { font-size: 11px; color: var(--faint); margin-top: 6px; }",
].join("\n");

const HEX = /^#[0-9a-f]{6}$/i;

function rgb(hex: string) {
  return [1, 3, 5].map((i) => parseInt(hex.slice(i, i + 2), 16));
}

// WCAG relative luminance and contrast ratio.
function luminance(hex: string) {
  const [r, g, b] = rgb(hex).map((v) => {
    v /= 255;
    return v <= 0.03928 ? v / 12.92 : ((v + 0.055) / 1.055) ** 2.4;
  });
  return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}

function contrast(a: string, b: string) {
  const [x, y] = [luminance(a), luminance(b)].sort((p, q) => q - p);
  return (x + 0.05) / (y + 0.05);
}

// `weight` of a over b.
function mix(a: string, b: string, weight: number) {
  const bb = rgb(b);
  return (
    "#" +
    rgb(a)
      .map((v, i) => Math.round(v * weight + bb[i] * (1 - weight)).toString(16).padStart(2, "0"))
      .join("")
  );
}

/** The :root colours of style.css for a button colour and background. */
export function guestThemeVars(accentIn = DEFAULT_ACCENT, background: GuestBackground = "light") {
  const accent = HEX.test(accentIn) ? accentIn.toLowerCase() : DEFAULT_ACCENT;
  const dark = background === "dark";
  const v = dark
    ? { bg: "#141414", ink: "#f2f2ef", muted: "#b5b3ad", card: "#1f1f1f", rule: "#3a3936", faint: "#8a8883", ok: "#5bd68a" }
    : {
        bg: background === "tint" ? mix(accent, "#ffffff", 0.12) : "#f2f2ef",
        ink: "#111111", muted: "#555555", card: "#ffffff", rule: "#cccccc", faint: "#999999", ok: "#157f3b",
      };
  // A button that would vanish into the page (ink on dark) takes the text colour.
  const button = contrast(accent, v.bg) < 1.2 ? v.ink : accent;
  const on = contrast(button, "#ffffff") >= contrast(button, "#111111") ? "#ffffff" : "#111111";
  // Pikro's amber focus ring unless the operator picked a colour.
  const focus = accent === DEFAULT_ACCENT ? "#ffbf00" : accent;
  const vars = { ...v, accent: button, on, focus };
  return `:root{${Object.entries(vars).map(([k, c]) => `--${k}:${c}`).join(";")}}`;
}
const STYLESHEET = '<link rel="stylesheet" href="/css/style.css" />';

function head(lang: GuestLanguage, title: string, extra = "") {
  return [
    "<!DOCTYPE html>",
    `<html lang="${lang}">`,
    "<head>",
    extra,
    '<meta charset="UTF-8" />',
    '<meta http-equiv="pragma" content="no-cache" />',
    '<meta http-equiv="expires" content="-1" />',
    '<meta name="viewport" content="width=device-width, initial-scale=1.0" />',
    `<title>${title}<\/title>`,
    STYLESHEET,
  ].join("\n");
}

// RouterOS only reports English errors; swap the common ones for plain text
// in the page language and leave anything else as RouterOS wrote it.
function errorScript(
  msgs: Record<"wrong" | "used" | "data" | "devices", string>,
) {
  return [
    "<script>",
    '(function(){var e=document.getElementById("err");if(!e)return;var m=e.textContent.toLowerCase();',
    `var M=${JSON.stringify(msgs)};`,
    "if(/invalid (username|password)|no such user/.test(m))e.textContent=M.wrong;",
    "else if(/uptime limit|session limit|expired/.test(m))e.textContent=M.used;",
    "else if(/traffic limit/.test(m))e.textContent=M.data;",
    "else if(/no more sessions|already logged|not yours/.test(m))e.textContent=M.devices;",
    "})();",
    "<\/script>",
  ].join("\n");
}

export function renderGuestPages(o: GuestPageOptions): GuestPages {
  const t = (key: string, named: Record<string, unknown> = {}) =>
    i18n.global.t(key, named, { locale: o.language });
  const title = esc(o.title);
  const subtitle = esc(o.subtitle);
  const foot = `<p class="foot">${t("guestPage.poweredBy")}<\/p>`;
  const header = [`<h1>${title}<\/h1>`, `<p class="sub">${subtitle}<\/p>`].join(
    "\n",
  );

  const prices = o.prices.length
    ? [
        '<div class="box">',
        `<h2>${t("guestPage.prices")}<\/h2>`,
        ...o.prices.map(
          (p) =>
            `<div class="row"><span>${esc(p.name)}<\/span><b>${esc(p.price)}<\/b><\/div>`,
        ),
        "<\/div>",
      ].join("\n")
    : "";
  const contactText = (o.contact ?? "").trim().slice(0, CONTACT_MAX);
  const contact = contactText ? `<p class="contact">${esc(contactText)}<\/p>` : "";

  const form =
    o.mode === "code"
      ? [
          '<form name="sendin" action="$(link-login-only)" method="post" style="display:none">',
          '  <input type="hidden" name="username" />',
          '  <input type="hidden" name="password" />',
          '  <input type="hidden" name="dst" value="$(link-orig)" />',
          '  <input type="hidden" name="popup" value="true" />',
          "<\/form>",
          "$(if chap-id)",
          CHAP_SCRIPT,
          "$(endif)",
          // The code is both username and password, so guests type it once.
          "<script>",
          "function doCodeLogin(){",
          "  var code=document.login.code.value.trim();",
          "  document.sendin.username.value=code;",
          "  document.sendin.password.value=$(if chap-id)hexMD5('$(chap-id)'+code+'$(chap-challenge)')$(else)code$(endif);",
          "  document.sendin.submit();",
          "  return false;",
          "}",
          "<\/script>",
          '<form name="login" onsubmit="return doCodeLogin()">',
          `  <label for="c">${t("guestPage.code")}<\/label>`,
          '  <input id="c" class="code" type="text" name="code" autocomplete="off" autocapitalize="none" autocorrect="off" spellcheck="false" required />',
          `  <button type="submit" class="btn">${t("guestPage.connect")}<\/button>`,
          "<\/form>",
        ].join("\n")
      : [
          CHAP_HIDDEN_FORM,
          '<form name="login" action="$(link-login-only)" method="post" $(if chap-id)onsubmit="return doLogin()"$(endif)>',
          '  <input type="hidden" name="dst" value="$(link-orig)" />',
          '  <input type="hidden" name="popup" value="true" />',
          `  <label for="u">${t("tickets.username")}<\/label>`,
          '  <input id="u" type="text" name="username" value="$(username)" autocomplete="username" autocapitalize="none" autocorrect="off" spellcheck="false" required />',
          `  <label for="p">${t("tickets.password")}<\/label>`,
          '  <input id="p" type="password" name="password" autocomplete="current-password" />',
          `  <button type="submit" class="btn">${t("guestPage.connect")}<\/button>`,
          "<\/form>",
        ].join("\n");

  const login = [
    head(o.language, title),
    "<\/head>",
    "<body>",
    '<main class="wrap">',
    header,
    '$(if error)<p class="err" id="err" role="alert">$(error)<\/p>$(endif)',
    form,
    prices,
    contact,
    foot,
    "<\/main>",
    errorScript({
      wrong: t("guestPage.errors.wrong"),
      used: t("guestPage.errors.used"),
      data: t("guestPage.errors.data"),
      devices: t("guestPage.errors.devices"),
    }),
    "<\/body>",
    "<\/html>",
  ].join("\n");

  const status = [
    head(
      o.language,
      t("guestPage.connected"),
      '$(if refresh-timeout)<meta http-equiv="refresh" content="$(refresh-timeout-secs)">$(endif)',
    ),
    "<script>",
    "$(if advert-pending == 'yes')",
    "var popup = ''; function focusAdvert() { if (window.focus) popup.focus(); } function openAdvert() { popup = open('$(link-advert)', 'hotspot_advert', ''); setTimeout(\"focusAdvert()\", 1000); }",
    "$(endif)",
    "function openLogout() { if (window.name !== 'hotspot_status') return true; open('$(link-logout)', 'hotspot_logout', 'toolbar=0,location=0,directories=0,status=0,menubars=0,resizable=1,width=360,height=400'); window.close(); return false; }",
    "<\/script>",
    "<\/head>",
    "<body $(if advert-pending == 'yes')onload=\"openAdvert()\"$(endif)>",
    '<main class="wrap">',
    header,
    `<p class="ok">${t("guestPage.connected")}<\/p>`,
    '<div class="box">',
    "$(if session-time-left)",
    `<div class="row"><span>${t("guestPage.timeLeft")}<\/span><b>$(session-time-left)<\/b><\/div>`,
    "$(endif)",
    `<div class="row"><span>${t("guestPage.timeUsed")}<\/span><b>$(uptime)<\/b><\/div>`,
    `<div class="row"><span>${t("guestPage.dataUsed")}<\/span><b>&darr; $(bytes-out-nice) / &uarr; $(bytes-in-nice)<\/b><\/div>`,
    `<div class="row"><span>${t("guestPage.ip")}<\/span><b>$(ip)<\/b><\/div>`,
    "<\/div>",
    "$(if login-by-mac != 'yes')",
    '<form action="$(link-logout)" name="logout" onsubmit="return openLogout()">',
    `  <button type="submit" class="btn">${t("guestPage.logout")}<\/button>`,
    "<\/form>",
    "$(endif)",
    foot,
    "<\/main>",
    "<\/body>",
    "<\/html>",
  ].join("\n");

  const logout = [
    head(o.language, t("guestPage.loggedOut")),
    "<script>",
    "function openLogin() { if (window.name !== 'hotspot_logout') return true; open('$(link-login)', '_blank', ''); window.close(); return false; }",
    "<\/script>",
    "<\/head>",
    "<body>",
    '<main class="wrap">',
    header,
    `<p class="sub">${t("guestPage.loggedOut")}<\/p>`,
    '<div class="box">',
    `<div class="row"><span>${t("guestPage.timeUsed")}<\/span><b>$(uptime)<\/b><\/div>`,
    "$(if session-time-left)",
    `<div class="row"><span>${t("guestPage.timeLeft")}<\/span><b>$(session-time-left)<\/b><\/div>`,
    "$(endif)",
    `<div class="row"><span>${t("guestPage.dataUsed")}<\/span><b>&darr; $(bytes-out-nice) / &uarr; $(bytes-in-nice)<\/b><\/div>`,
    "<\/div>",
    '<form action="$(link-login)" name="login" onsubmit="return openLogin()">',
    `  <button type="submit" class="btn">${t("guestPage.signInAgain")}<\/button>`,
    "<\/form>",
    prices,
    foot,
    "<\/main>",
    "<\/body>",
    "<\/html>",
  ].join("\n");

  const redirect = [
    head(
      o.language,
      t("guestPage.connected"),
      '<meta http-equiv="refresh" content="2; url=$(link-redirect)" />',
    ),
    "<script>",
    "function startClock() { $(if popup == 'true') open('$(link-status)', 'hotspot_status', 'toolbar=0,location=0,directories=0,status=0,menubars=0,resizable=1,width=360,height=400'); $(endif) location.href = unescape('$(link-redirect-esc)'); }",
    "<\/script>",
    "<\/head>",
    '<body onload="startClock()">',
    '<main class="wrap">',
    header,
    `<p class="ok">${t("guestPage.connected")}<\/p>`,
    `<p class="hint"><a href="$(link-redirect)">${t("guestPage.tapIfStuck")}<\/a><\/p>`,
    foot,
    "<\/main>",
    "<\/body>",
    "<\/html>",
  ].join("\n");

  // Shown right after a successful login: confirm, then carry on to the site
  // the guest first asked for, as RouterOS's stock page does.
  const alogin = [
    head(
      o.language,
      t("guestPage.connected"),
      '<meta http-equiv="refresh" content="3; url=$(link-redirect)" />',
    ),
    "<\/head>",
    "<body>",
    '<main class="wrap">',
    header,
    `<p class="ok">${t("guestPage.connected")}<\/p>`,
    `<a class="btn" href="$(link-status)">${t("guestPage.seeConnection")}<\/a>`,
    `<p class="hint"><a href="$(link-redirect)">${t("guestPage.continue")}<\/a><\/p>`,
    foot,
    "<\/main>",
    "<\/body>",
    "<\/html>",
  ].join("\n");

  return {
    login,
    status,
    logout,
    redirect,
    alogin,
    css: guestThemeVars(o.accent, o.background) + "\n" + GUEST_CSS,
    md5Js: MD5_JS,
  };
}

// The iframe preview has no RouterOS to expand $(...) tokens: show the page
// as a guest first sees it (no error, empty fields, sample session figures).
const PREVIEW_SAMPLES: Record<string, string> = {
  "session-time-left": "22h47m",
  uptime: "1h13m",
  "bytes-in-nice": "3.1 MiB",
  "bytes-out-nice": "48.2 MiB",
  ip: "10.5.50.23",
};

export function previewHtml(html: string, css = guestThemeVars() + "\n" + GUEST_CSS) {
  return html
    .replace(STYLESHEET, `<style>${css}<\/style>`) // the iframe can't fetch the router's file
    .replace(/\$\(if error\)[\s\S]*?\$\(endif\)/g, "")
    .replace(/\$\((?:if|elif) [^)]*\)|\$\((?:else|endif)\)/g, "")
    .replace(
      /\$\(([\w-]+)\)/g,
      (_, token: string) => PREVIEW_SAMPLES[token] ?? "",
    );
}
