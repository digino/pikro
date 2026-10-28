import i18n from '@/locales'

// Minified Paul Johnston reference MD5 impl (the same one MikroTik's own
// stock hotspot templates use)
const CHAP_SCRIPT = [
  '<script>',
  'function safe_add(f,i){var n=(65535&f)+(65535&i);return(f>>16)+(i>>16)+(n>>16)<<16|65535&n}function rol(f,i){return f<<i|f>>>32-i}function cmn(f,i,n,h,r,g){return safe_add(rol(safe_add(safe_add(i,f),safe_add(h,g)),r),n)}function ff(f,i,n,h,r,g,e){return cmn(i&n|~i&h,f,i,r,g,e)}function gg(f,i,n,h,r,g,e){return cmn(i&h|n&~h,f,i,r,g,e)}function hh(f,i,n,h,r,g,e){return cmn(i^n^h,f,i,r,g,e)}function ii(f,i,n,h,r,g,e){return cmn(n^(i|~h),f,i,r,g,e)}function coreMD5(f){for(var i=1732584193,n=-271733879,h=-1732584194,r=271733878,g=0;g<f.length;g+=16){var e=i,t=n,a=h,c=r;i=ff(i,n,h,r,f[g+0],7,-680876936),r=ff(r,i,n,h,f[g+1],12,-389564586),h=ff(h,r,i,n,f[g+2],17,606105819),n=ff(n,h,r,i,f[g+3],22,-1044525330),i=ff(i,n,h,r,f[g+4],7,-176418897),r=ff(r,i,n,h,f[g+5],12,1200080426),h=ff(h,r,i,n,f[g+6],17,-1473231341),n=ff(n,h,r,i,f[g+7],22,-45705983),i=ff(i,n,h,r,f[g+8],7,1770035416),r=ff(r,i,n,h,f[g+9],12,-1958414417),h=ff(h,r,i,n,f[g+10],17,-42063),n=ff(n,h,r,i,f[g+11],22,-1990404162),i=ff(i,n,h,r,f[g+12],7,1804603682),r=ff(r,i,n,h,f[g+13],12,-40341101),h=ff(h,r,i,n,f[g+14],17,-1502002290),i=gg(i,n=ff(n,h,r,i,f[g+15],22,1236535329),h,r,f[g+1],5,-165796510),r=gg(r,i,n,h,f[g+6],9,-1069501632),h=gg(h,r,i,n,f[g+11],14,643717713),n=gg(n,h,r,i,f[g+0],20,-373897302),i=gg(i,n,h,r,f[g+5],5,-701558691),r=gg(r,i,n,h,f[g+10],9,38016083),h=gg(h,r,i,n,f[g+15],14,-660478335),n=gg(n,h,r,i,f[g+4],20,-405537848),i=gg(i,n,h,r,f[g+9],5,568446438),r=gg(r,i,n,h,f[g+14],9,-1019803690),h=gg(h,r,i,n,f[g+3],14,-187363961),n=gg(n,h,r,i,f[g+8],20,1163531501),i=gg(i,n,h,r,f[g+13],5,-1444681467),r=gg(r,i,n,h,f[g+2],9,-51403784),h=gg(h,r,i,n,f[g+7],14,1735328473),i=hh(i,n=gg(n,h,r,i,f[g+12],20,-1926607734),h,r,f[g+5],4,-378558),r=hh(r,i,n,h,f[g+8],11,-2022574463),h=hh(h,r,i,n,f[g+11],16,1839030562),n=hh(n,h,r,i,f[g+14],23,-35309556),i=hh(i,n,h,r,f[g+1],4,-1530992060),r=hh(r,i,n,h,f[g+4],11,1272893353),h=hh(h,r,i,n,f[g+7],16,-155497632),n=hh(n,h,r,i,f[g+10],23,-1094730640),i=hh(i,n,h,r,f[g+13],4,681279174),r=hh(r,i,n,h,f[g+0],11,-358537222),h=hh(h,r,i,n,f[g+3],16,-722521979),n=hh(n,h,r,i,f[g+6],23,76029189),i=hh(i,n,h,r,f[g+9],4,-640364487),r=hh(r,i,n,h,f[g+12],11,-421815835),h=hh(h,r,i,n,f[g+15],16,530742520),i=ii(i,n=hh(n,h,r,i,f[g+2],23,-995338651),h,r,f[g+0],6,-198630844),r=ii(r,i,n,h,f[g+7],10,1126891415),h=ii(h,r,i,n,f[g+14],15,-1416354905),n=ii(n,h,r,i,f[g+5],21,-57434055),i=ii(i,n,h,r,f[g+12],6,1700485571),r=ii(r,i,n,h,f[g+3],10,-1894986606),h=ii(h,r,i,n,f[g+10],15,-1051523),n=ii(n,h,r,i,f[g+1],21,-2054922799),i=ii(i,n,h,r,f[g+8],6,1873313359),r=ii(r,i,n,h,f[g+15],10,-30611744),h=ii(h,r,i,n,f[g+6],15,-1560198380),n=ii(n,h,r,i,f[g+13],21,1309151649),i=ii(i,n,h,r,f[g+4],6,-145523070),r=ii(r,i,n,h,f[g+11],10,-1120210379),h=ii(h,r,i,n,f[g+2],15,718787259),n=ii(n,h,r,i,f[g+9],21,-343485551),i=safe_add(i,e),n=safe_add(n,t),h=safe_add(h,a),r=safe_add(r,c)}return[i,n,h,r]}function binl2hex(f){for(var i="0123456789abcdef",n="",h=0;h<4*f.length;h++)n+=i.charAt(f[h>>2]>>h%4*8+4&15)+i.charAt(f[h>>2]>>h%4*8&15);return n}function str2binl(f){for(var i=1+(f.length+8>>6),n=new Array(16*i),h=0;h<16*i;h++)n[h]=0;for(h=0;h<f.length;h++)n[h>>2]|=(255&f.charCodeAt(h))<<h%4*8;return n[h>>2]|=128<<h%4*8,n[16*i-2]=8*f.length,n}function hexMD5(f){return binl2hex(coreMD5(str2binl(f)))}',
  'function doLogin(){document.sendin.username.value=document.login.username.value;document.sendin.password.value=hexMD5(\'$(chap-id)\'+document.login.password.value+\'$(chap-challenge)\');document.sendin.submit();return false}',
  '<\/script>',
].join('\n')

const CHAP_HIDDEN_FORM = [
  '$(if chap-id)',
  '<form name="sendin" action="$(link-login-only)" method="post" style="display:none">',
  '  <input type="hidden" name="username" />',
  '  <input type="hidden" name="password" />',
  '  <input type="hidden" name="dst" value="$(link-orig)" />',
  '  <input type="hidden" name="popup" value="true" />',
  '<\/form>',
  CHAP_SCRIPT,
  '$(endif)',
].join('\n')

// ── Guest pages ──────────────────────────────────────────────────────────────
// One Ink & Amber style for login, status, logout and redirect. The operator
// picks how guests log in and the page language; alogin/error stay the fixed
// Go defaults. Text comes from the `guestPage` locale keys, in the page's
// language rather than the language Pikro itself is shown in.

export type LoginMode = 'code' | 'password'
export type GuestLanguage = 'en' | 'fr'

export interface GuestPageOptions {
  title: string
  subtitle: string
  mode: LoginMode
  language: GuestLanguage
  prices: { name: string; price: string }[]
  contact: string
}

export interface GuestPages {
  login: string
  status: string
  logout: string
  redirect: string
  alogin: string
}

function esc(s: string) {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
}

const STYLE = [
  '<style>',
  '*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }',
  'body { font-family: system-ui, -apple-system, "Segoe UI", Roboto, sans-serif; background: #f2f2ef; color: #111; min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 24px 20px; }',
  '.wrap { width: 100%; max-width: 340px; display: flex; flex-direction: column; align-items: center; gap: 14px; text-align: center; }',
  '.tile { width: 64px; height: 64px; border-radius: 18px; background: #111; color: #fff; display: flex; align-items: center; justify-content: center; }',
  'h1 { font-size: 22px; font-weight: 800; line-height: 1.15; letter-spacing: -.01em; }',
  '.sub { font-size: 13px; color: #555; }',
  'form { width: 100%; display: flex; flex-direction: column; gap: 10px; margin-top: 10px; }',
  'label { text-align: left; font-size: 12px; font-weight: 600; color: #555; }',
  'input { width: 100%; height: 52px; padding: 0 14px; border: 2px solid #111; border-radius: 12px; background: #fff; color: #111; font-size: 16px; font-weight: 500; font-family: inherit; -webkit-appearance: none; }',
  'input.code { text-align: center; font: 600 20px ui-monospace, "SF Mono", Menlo, Consolas, monospace; letter-spacing: .12em; }',
  'input:focus { outline: 3px solid #ffbf00; outline-offset: 1px; }',
  '.btn { display: flex; align-items: center; justify-content: center; width: 100%; height: 50px; border: 0; border-radius: 12px; background: #111; color: #fff; font-size: 15px; font-weight: 700; font-family: inherit; cursor: pointer; text-decoration: none; }',
  '.hint { font-size: 12px; color: #666; }',
  '.hint a { color: #111; font-weight: 600; }',
  '.err { width: 100%; padding: 10px 12px; border: 2px solid #b42318; border-radius: 12px; background: #fdecea; color: #b42318; font-size: 13px; font-weight: 600; }',
  '.box { width: 100%; border: 2px solid #111; border-radius: 12px; background: #fff; padding: 12px 14px; text-align: left; }',
  '.box h2 { font-size: 12px; font-weight: 600; color: #555; margin-bottom: 6px; }',
  '.row { display: flex; justify-content: space-between; gap: 12px; padding: 5px 0; font-size: 14px; }',
  '.row + .row { border-top: 1px dashed #ccc; }',
  '.row b { font-variant-numeric: tabular-nums; white-space: nowrap; }',
  '.buy { display: flex; align-items: center; justify-content: center; gap: 8px; width: 100%; min-height: 48px; border: 2px solid #111; border-radius: 12px; background: #ffbf00; color: #111; font-weight: 700; font-size: 14px; text-decoration: none; }',
  '.ok { display: inline-flex; align-items: center; gap: 6px; font-size: 14px; font-weight: 700; color: #157f3b; }',
  '.ok::before { content: ""; width: 9px; height: 9px; border-radius: 50%; background: #157f3b; }',
  '.foot { font-size: 11px; color: #999; margin-top: 6px; }',
  '<\/style>',
].join('\n')

const WIFI_TILE = '<div class="tile"><svg aria-hidden="true" width="34" height="34" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.25" stroke-linecap="round"><path d="M2 8.8a15.5 15.5 0 0 1 20 0M5 12.6a11 11 0 0 1 14 0M8.5 16.2a6 6 0 0 1 7 0"/><circle cx="12" cy="20" r="1.4" fill="currentColor"/><\/svg><\/div>'

function head(lang: GuestLanguage, title: string, extra = '') {
  return [
    '<!DOCTYPE html>',
    `<html lang="${lang}">`,
    '<head>',
    extra,
    '<meta charset="UTF-8" />',
    '<meta http-equiv="pragma" content="no-cache" />',
    '<meta http-equiv="expires" content="-1" />',
    '<meta name="viewport" content="width=device-width, initial-scale=1.0" />',
    `<title>${title}<\/title>`,
    STYLE,
  ].join('\n')
}

// RouterOS only reports English errors; swap the common ones for plain text
// in the page language and leave anything else as RouterOS wrote it.
function errorScript(msgs: Record<'wrong' | 'used' | 'data' | 'devices', string>) {
  return [
    '<script>',
    '(function(){var e=document.getElementById("err");if(!e)return;var m=e.textContent.toLowerCase();',
    `var M=${JSON.stringify(msgs)};`,
    'if(/invalid (username|password)|no such user/.test(m))e.textContent=M.wrong;',
    'else if(/uptime limit|session limit|expired/.test(m))e.textContent=M.used;',
    'else if(/traffic limit/.test(m))e.textContent=M.data;',
    'else if(/no more sessions|already logged|not yours/.test(m))e.textContent=M.devices;',
    '})();',
    '<\/script>',
  ].join('\n')
}

export function renderGuestPages(o: GuestPageOptions): GuestPages {
  const t = (key: string, named: Record<string, unknown> = {}) =>
    i18n.global.t(key, named, { locale: o.language })
  const title = esc(o.title)
  const subtitle = esc(o.subtitle)
  const foot = `<p class="foot">${t('guestPage.poweredBy')}<\/p>`
  const header = [WIFI_TILE, `<h1>${title}<\/h1>`, `<p class="sub">${subtitle}<\/p>`].join('\n')

  const prices = o.prices.length
    ? [
        '<div class="box">',
        `<h2>${t('guestPage.prices')}<\/h2>`,
        ...o.prices.map(p => `<div class="row"><span>${esc(p.name)}<\/span><b>${esc(p.price)}<\/b><\/div>`),
        '<\/div>',
      ].join('\n')
    : ''
  const tel = o.contact.replace(/[^\d+]/g, '')
  const contact = tel
    ? `<a class="buy" href="tel:${tel}">${t('guestPage.buy')} · ${esc(o.contact)}<\/a>`
    : ''

  const form = o.mode === 'code'
    ? [
        '<form name="sendin" action="$(link-login-only)" method="post" style="display:none">',
        '  <input type="hidden" name="username" />',
        '  <input type="hidden" name="password" />',
        '  <input type="hidden" name="dst" value="$(link-orig)" />',
        '  <input type="hidden" name="popup" value="true" />',
        '<\/form>',
        '$(if chap-id)',
        CHAP_SCRIPT,
        '$(endif)',
        // The code is both username and password, so guests type it once.
        '<script>',
        'function doCodeLogin(){',
        '  var code=document.login.code.value.trim();',
        '  document.sendin.username.value=code;',
        '  document.sendin.password.value=$(if chap-id)hexMD5(\'$(chap-id)\'+code+\'$(chap-challenge)\')$(else)code$(endif);',
        '  document.sendin.submit();',
        '  return false;',
        '}',
        '<\/script>',
        '<form name="login" onsubmit="return doCodeLogin()">',
        `  <label for="c">${t('guestPage.code')}<\/label>`,
        '  <input id="c" class="code" type="text" name="code" autocomplete="off" autocapitalize="none" autocorrect="off" spellcheck="false" required />',
        `  <button type="submit" class="btn">${t('guestPage.connect')}<\/button>`,
        '<\/form>',
      ].join('\n')
    : [
        CHAP_HIDDEN_FORM,
        '<form name="login" action="$(link-login-only)" method="post" $(if chap-id)onsubmit="return doLogin()"$(endif)>',
        '  <input type="hidden" name="dst" value="$(link-orig)" />',
        '  <input type="hidden" name="popup" value="true" />',
        `  <label for="u">${t('tickets.username')}<\/label>`,
        '  <input id="u" type="text" name="username" value="$(username)" autocomplete="username" autocapitalize="none" autocorrect="off" spellcheck="false" required />',
        `  <label for="p">${t('tickets.password')}<\/label>`,
        '  <input id="p" type="password" name="password" autocomplete="current-password" />',
        `  <button type="submit" class="btn">${t('guestPage.connect')}<\/button>`,
        '<\/form>',
      ].join('\n')

  const login = [
    head(o.language, title),
    '<\/head>',
    '<body>',
    '<main class="wrap">',
    header,
    '$(if error)<p class="err" id="err" role="alert">$(error)<\/p>$(endif)',
    form,
    `<p class="hint">${t(o.mode === 'code' ? 'guestPage.hint' : 'guestPage.hintPassword')}<\/p>`,
    prices,
    contact,
    foot,
    '<\/main>',
    errorScript({
      wrong: t('guestPage.errors.wrong'),
      used: t('guestPage.errors.used'),
      data: t('guestPage.errors.data'),
      devices: t('guestPage.errors.devices'),
    }),
    '<\/body>',
    '<\/html>',
  ].join('\n')

  const status = [
    head(o.language, t('guestPage.connected'), '$(if refresh-timeout)<meta http-equiv="refresh" content="$(refresh-timeout-secs)">$(endif)'),
    '<script>',
    '$(if advert-pending == \'yes\')',
    'var popup = \'\'; function focusAdvert() { if (window.focus) popup.focus(); } function openAdvert() { popup = open(\'$(link-advert)\', \'hotspot_advert\', \'\'); setTimeout("focusAdvert()", 1000); }',
    '$(endif)',
    'function openLogout() { if (window.name !== \'hotspot_status\') return true; open(\'$(link-logout)\', \'hotspot_logout\', \'toolbar=0,location=0,directories=0,status=0,menubars=0,resizable=1,width=360,height=400\'); window.close(); return false; }',
    '<\/script>',
    '<\/head>',
    '<body $(if advert-pending == \'yes\')onload="openAdvert()"$(endif)>',
    '<main class="wrap">',
    header,
    `<p class="ok">${t('guestPage.connected')}<\/p>`,
    '<div class="box">',
    '$(if session-time-left)',
    `<div class="row"><span>${t('guestPage.timeLeft')}<\/span><b>$(session-time-left)<\/b><\/div>`,
    '$(endif)',
    `<div class="row"><span>${t('guestPage.timeUsed')}<\/span><b>$(uptime)<\/b><\/div>`,
    `<div class="row"><span>${t('guestPage.dataUsed')}<\/span><b>&darr; $(bytes-out-nice) / &uarr; $(bytes-in-nice)<\/b><\/div>`,
    `<div class="row"><span>${t('guestPage.ip')}<\/span><b>$(ip)<\/b><\/div>`,
    '<\/div>',
    '$(if login-by-mac != \'yes\')',
    '<form action="$(link-logout)" name="logout" onsubmit="return openLogout()">',
    `  <button type="submit" class="btn">${t('guestPage.logout')}<\/button>`,
    '<\/form>',
    '$(endif)',
    contact,
    foot,
    '<\/main>',
    '<\/body>',
    '<\/html>',
  ].join('\n')

  const logout = [
    head(o.language, t('guestPage.loggedOut')),
    '<script>',
    'function openLogin() { if (window.name !== \'hotspot_logout\') return true; open(\'$(link-login)\', \'_blank\', \'\'); window.close(); return false; }',
    '<\/script>',
    '<\/head>',
    '<body>',
    '<main class="wrap">',
    header,
    `<p class="sub">${t('guestPage.loggedOut')}<\/p>`,
    '<div class="box">',
    `<div class="row"><span>${t('guestPage.timeUsed')}<\/span><b>$(uptime)<\/b><\/div>`,
    '$(if session-time-left)',
    `<div class="row"><span>${t('guestPage.timeLeft')}<\/span><b>$(session-time-left)<\/b><\/div>`,
    '$(endif)',
    `<div class="row"><span>${t('guestPage.dataUsed')}<\/span><b>&darr; $(bytes-out-nice) / &uarr; $(bytes-in-nice)<\/b><\/div>`,
    '<\/div>',
    '<form action="$(link-login)" name="login" onsubmit="return openLogin()">',
    `  <button type="submit" class="btn">${t('guestPage.signInAgain')}<\/button>`,
    '<\/form>',
    prices,
    contact,
    foot,
    '<\/main>',
    '<\/body>',
    '<\/html>',
  ].join('\n')

  const redirect = [
    head(o.language, t('guestPage.connected'), '<meta http-equiv="refresh" content="2; url=$(link-redirect)" />'),
    '<script>',
    'function startClock() { $(if popup == \'true\') open(\'$(link-status)\', \'hotspot_status\', \'toolbar=0,location=0,directories=0,status=0,menubars=0,resizable=1,width=360,height=400\'); $(endif) location.href = unescape(\'$(link-redirect-esc)\'); }',
    '<\/script>',
    '<\/head>',
    '<body onload="startClock()">',
    '<main class="wrap">',
    header,
    `<p class="ok">${t('guestPage.connected')}<\/p>`,
    `<p class="hint"><a href="$(link-redirect)">${t('guestPage.tapIfStuck')}<\/a><\/p>`,
    foot,
    '<\/main>',
    '<\/body>',
    '<\/html>',
  ].join('\n')

  // Shown right after a successful login: confirm, then carry on to the site
  // the guest first asked for, as RouterOS's stock page does.
  const alogin = [
    head(o.language, t('guestPage.connected'), '<meta http-equiv="refresh" content="3; url=$(link-redirect)" />'),
    '<\/head>',
    '<body>',
    '<main class="wrap">',
    header,
    `<p class="ok">${t('guestPage.connected')}<\/p>`,
    `<a class="btn" href="$(link-status)">${t('guestPage.seeConnection')}<\/a>`,
    `<p class="hint"><a href="$(link-redirect)">${t('guestPage.continue')}<\/a><\/p>`,
    foot,
    '<\/main>',
    '<\/body>',
    '<\/html>',
  ].join('\n')

  return { login, status, logout, redirect, alogin }
}

// The iframe preview has no RouterOS to expand $(...) tokens: show the page
// as a guest first sees it (no error, empty fields, sample session figures).
const PREVIEW_SAMPLES: Record<string, string> = {
  'session-time-left': '22h47m',
  'uptime': '1h13m',
  'bytes-in-nice': '3.1 MiB',
  'bytes-out-nice': '48.2 MiB',
  'ip': '10.5.50.23',
}

export function previewHtml(html: string) {
  return html
    .replace(/\$\(if error\)[\s\S]*?\$\(endif\)/g, '')
    .replace(/\$\((?:if|elif) [^)]*\)|\$\((?:else|endif)\)/g, '')
    .replace(/\$\(([\w-]+)\)/g, (_, token: string) => PREVIEW_SAMPLES[token] ?? '')
}
