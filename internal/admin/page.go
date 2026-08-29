package admin

import (
	"html"
	"strings"
)

func PageHTML() string {
	return `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<title>CodeBuddy Proxy · Console</title>
<style>
:root{
  --sand:#F4EEE4;
  --sand-2:#EFE6D8;
  --paper:#FFFCF7;
  --ink:#1B2430;
  --muted:#667384;
  --line:rgba(27,36,48,.10);
  --teal:#0F7C74;
  --teal-deep:#0A5C57;
  --coral:#E86F3A;
  --coral-deep:#C95524;
  --sky:#3B82A0;
  --ok:#1F8A5B;
  --bad:#C94444;
  --warn:#B7791F;
  --shadow:0 18px 40px -24px rgba(27,36,48,.35);
  --ease:cubic-bezier(.32,.72,0,1);
}
*{box-sizing:border-box}
html,body{margin:0;min-height:100%}
body{
  font-family:system-ui,-apple-system,"Segoe UI",sans-serif;
  color:var(--ink);
  background:
    radial-gradient(900px 520px at 8% -8%, rgba(15,124,116,.18), transparent 55%),
    radial-gradient(760px 480px at 96% 4%, rgba(232,111,58,.16), transparent 50%),
    radial-gradient(680px 420px at 70% 90%, rgba(59,130,160,.12), transparent 55%),
    linear-gradient(180deg, var(--sand) 0%, #F8F3EA 48%, #F1EBE1 100%);
}
body::before{
  content:"";
  position:fixed; inset:0; pointer-events:none; z-index:0; opacity:.035;
  background-image:url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='.85' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)'/%3E%3C/svg%3E");
}
a{color:inherit}
.shell{position:relative;z-index:1;max-width:1180px;margin:0 auto;padding:28px 20px 72px}
.topbar{
  display:flex;align-items:center;justify-content:space-between;gap:16px;flex-wrap:wrap;
  margin-bottom:28px;padding:12px 14px;border-radius:999px;
  background:rgba(255,252,247,.72);border:1px solid var(--line);
  box-shadow:var(--shadow);backdrop-filter:blur(14px);
}
.brand{display:flex;align-items:center;gap:12px;min-width:0}
.mark{
  width:38px;height:38px;border-radius:14px;flex:none;
  background:
    linear-gradient(145deg, rgba(255,255,255,.35), transparent 42%),
    linear-gradient(135deg, var(--teal) 0%, var(--sky) 55%, var(--coral) 120%);
  box-shadow:inset 0 1px 0 rgba(255,255,255,.35), 0 8px 18px -10px rgba(15,124,116,.55);
}
.brand-copy{min-width:0}
.brand-copy strong{display:block;font-size:15px;letter-spacing:-.02em}
.brand-copy span{display:block;font-size:12px;color:var(--muted)}
.pillrow{display:flex;gap:8px;flex-wrap:wrap;align-items:center}
.pill{
  display:inline-flex;align-items:center;gap:7px;
  padding:7px 12px;border-radius:999px;font-size:12px;font-weight:500;
  background:rgba(255,252,247,.9);border:1px solid var(--line);color:var(--muted);
}
.pill .dot{width:7px;height:7px;border-radius:50%;background:var(--teal);box-shadow:0 0 0 3px rgba(15,124,116,.15)}
.pill .dot.warn{background:var(--warn);box-shadow:0 0 0 3px rgba(183,121,31,.16)}
.pill .dot.bad{background:var(--bad);box-shadow:0 0 0 3px rgba(201,68,68,.15)}
.hero{
  display:grid;grid-template-columns:1.15fr .85fr;gap:18px;margin-bottom:18px;
}
@media (max-width:920px){.hero{grid-template-columns:1fr}}
.panel{
  position:relative;border-radius:28px;padding:4px;
  background:linear-gradient(180deg, rgba(255,255,255,.75), rgba(255,255,255,.28));
  border:1px solid rgba(255,255,255,.55);box-shadow:var(--shadow);
}
.panel-inner{
  height:100%;border-radius:24px;background:var(--paper);
  border:1px solid var(--line);padding:22px 22px 20px;
  box-shadow:inset 0 1px 0 rgba(255,255,255,.8);
}
.eyebrow{
  display:inline-flex;align-items:center;gap:8px;
  padding:5px 10px;border-radius:999px;margin-bottom:14px;
  font-size:11px;font-weight:600;letter-spacing:.14em;text-transform:uppercase;
  color:var(--teal-deep);background:rgba(15,124,116,.08);border:1px solid rgba(15,124,116,.14);
}
h1{
  margin:0 0 10px;font-size:clamp(30px,4vw,42px);line-height:1.05;
  letter-spacing:-.04em;font-weight:700;max-width:12ch;
}
.lede{margin:0;color:var(--muted);font-size:15px;line-height:1.55;max-width:42ch}
.metrics{
  display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:10px;margin-top:22px;
}
@media (max-width:720px){.metrics{grid-template-columns:repeat(2,minmax(0,1fr))}}
.metric{
  padding:14px 14px 12px;border-radius:18px;
  background:linear-gradient(180deg,#FFFEFB,var(--sand-2));
  border:1px solid var(--line);
}
.metric .k{font-size:11px;color:var(--muted);letter-spacing:.04em;text-transform:uppercase;font-weight:600}
.metric .v{margin-top:6px;font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:22px;letter-spacing:-.03em;font-weight:500}
.metric .v.sm{font-size:14px;line-height:1.35}
.section-head{
  display:flex;align-items:flex-end;justify-content:space-between;gap:12px;flex-wrap:wrap;
  margin:0 0 14px;
}
.section-head h2{margin:0;font-size:18px;letter-spacing:-.02em}
.section-head p{margin:4px 0 0;color:var(--muted);font-size:13px}
.actions{display:flex;gap:8px;flex-wrap:wrap}
.btn,.linkbtn,button,select,input{
  font:inherit;
}
button,.btn,.linkbtn{
  appearance:none;border:0;cursor:pointer;
  display:inline-flex;align-items:center;justify-content:center;gap:8px;
  padding:11px 16px;border-radius:999px;font-weight:600;font-size:13.5px;
  transition:transform .35s var(--ease), background .35s var(--ease), box-shadow .35s var(--ease), color .35s var(--ease);
}
button:active,.btn:active,.linkbtn:active{transform:scale(.98)}
button.primary,.btn.primary{
  color:#fff;background:linear-gradient(180deg, #F08A55, var(--coral));
  box-shadow:0 12px 24px -14px rgba(232,111,58,.85), inset 0 1px 0 rgba(255,255,255,.28);
}
button.primary:hover{background:linear-gradient(180deg, #F49A6B, var(--coral-deep))}
button.teal{
  color:#fff;background:linear-gradient(180deg, #1A9A90, var(--teal));
  box-shadow:0 12px 24px -14px rgba(15,124,116,.8), inset 0 1px 0 rgba(255,255,255,.25);
}
button.ghost,.linkbtn{
  color:var(--ink);background:rgba(255,252,247,.9);border:1px solid var(--line);
}
button.ghost:hover,.linkbtn:hover{background:#fff}
button.danger{color:#fff;background:linear-gradient(180deg,#D65A5A,var(--bad));box-shadow:0 10px 20px -14px rgba(201,68,68,.7)}
.linkbtn{text-decoration:none}
.field-grid{display:grid;grid-template-columns:1fr 1.4fr;gap:10px;margin:4px 0 14px}
@media (max-width:640px){.field-grid{grid-template-columns:1fr}}
label{display:block;font-size:12px;font-weight:600;color:var(--muted);margin:0 0 6px}
select,input{
  width:100%;padding:12px 14px;border-radius:14px;
  border:1px solid var(--line);background:#FFFEFB;color:var(--ink);
  outline:none;transition:border-color .3s var(--ease), box-shadow .3s var(--ease);
}
select:focus,input:focus{border-color:rgba(15,124,116,.45);box-shadow:0 0 0 4px rgba(15,124,116,.12)}
.oauth-status{
  margin-top:14px;padding:14px;border-radius:18px;
  background:linear-gradient(135deg, rgba(15,124,116,.07), rgba(232,111,58,.06));
  border:1px dashed rgba(15,124,116,.22);
}
.oauth-status .title{font-size:12px;font-weight:600;color:var(--muted);margin-bottom:6px;text-transform:uppercase;letter-spacing:.08em}
.oauth-status .msg{font-size:14px;line-height:1.45;min-height:1.45em}
.stack{display:grid;gap:16px}
.account{
  display:grid;gap:12px;padding:16px;border-radius:20px;
  background:linear-gradient(180deg,#FFFEFB, #F7F1E7);
  border:1px solid var(--line);
}
.account-top{display:flex;justify-content:space-between;gap:12px;flex-wrap:wrap;align-items:flex-start}
.account-title{display:flex;flex-wrap:wrap;gap:8px;align-items:center}
.account-title strong{font-size:15px;letter-spacing:-.02em}
.badge{
  display:inline-flex;align-items:center;padding:4px 9px;border-radius:999px;
  font-size:11px;font-weight:600;border:1px solid var(--line);background:#fff;color:var(--muted);
}
.badge.site{color:var(--teal-deep);background:rgba(15,124,116,.08);border-color:rgba(15,124,116,.16)}
.badge.muted{color:var(--ink-soft);background:rgba(47,53,48,.05);border-color:rgba(47,53,48,.1)}
.seg{display:inline-flex;gap:6px;padding:4px;border-radius:999px;background:rgba(47,53,48,.04);border:1px solid rgba(47,53,48,.08)}
.seg button{border:0;background:transparent;color:var(--ink-soft);padding:8px 14px;border-radius:999px;cursor:pointer;font:inherit}
.seg button.active{background:var(--teal);color:#fff;box-shadow:0 8px 18px rgba(15,124,116,.18)}
.seg-hint{margin:8px 0 0;color:var(--ink-soft);font-size:12px}
.badge.on{color:var(--ok);background:rgba(31,138,91,.08);border-color:rgba(31,138,91,.18)}
.badge.off{color:var(--bad);background:rgba(201,68,68,.08);border-color:rgba(201,68,68,.16)}
.meta{color:var(--muted);font-size:13px;line-height:1.5}
.meta code, .mono{
  font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:12px;
}
.err{margin-top:2px;color:var(--bad);font-size:12.5px}
.chips{display:flex;flex-wrap:wrap;gap:8px}
.chip{
  padding:8px 12px;border-radius:999px;font-size:12.5px;font-weight:500;
  background:#FFFEFB;border:1px solid var(--line);color:var(--ink);
}
.empty{
  padding:28px 18px;border-radius:20px;text-align:center;
  border:1px dashed rgba(27,36,48,.16);color:var(--muted);background:rgba(255,252,247,.55);
}
details.raw{
  margin-top:8px;border-radius:18px;border:1px solid var(--line);background:rgba(255,252,247,.7);overflow:hidden;
}
details.raw summary{
  cursor:pointer;list-style:none;padding:12px 14px;font-size:13px;font-weight:600;color:var(--muted);
}
details.raw summary::-webkit-details-marker{display:none}
pre{
  margin:0;padding:0 14px 14px;white-space:pre-wrap;word-break:break-word;
  font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:11.5px;line-height:1.5;color:#334155;
  max-height:280px;overflow:auto;
}
.footer-note{margin-top:22px;color:var(--muted);font-size:12px;text-align:center}

.copyline{display:grid;grid-template-columns:1fr auto;gap:8px;align-items:center}
.copyline input{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:12.5px}
.secret-hint{margin-top:8px;font-size:12px;color:var(--muted);line-height:1.45}
.toast{
  position:fixed;right:18px;bottom:18px;z-index:40;min-width:180px;max-width:min(420px,92vw);
  padding:12px 14px;border-radius:14px;background:#1B2430;color:#FFFCF7;font-size:13px;font-weight:600;
  box-shadow:0 16px 30px -18px rgba(27,36,48,.55);opacity:0;transform:translateY(8px);
  transition:opacity .35s var(--ease), transform .35s var(--ease);pointer-events:none;
}
.toast.show{opacity:1;transform:translateY(0)}
.toast.err{background:var(--bad)}
.chip.btnish{cursor:pointer}
.chip.btnish:hover{border-color:rgba(15,124,116,.35);background:rgba(15,124,116,.06)}
.usage-line{margin-top:4px;font-size:12.5px}
.usage-line .pill{display:inline-flex;padding:3px 8px;border-radius:999px;font-weight:600;font-size:12px;border:1px solid var(--line)}
.usage-line .pill.good{color:var(--ok);background:rgba(31,138,91,.08);border-color:rgba(31,138,91,.18)}
.usage-line .pill.warn{color:var(--warn);background:rgba(183,121,31,.10);border-color:rgba(183,121,31,.20)}
.usage-line .pill.bad{color:var(--bad);background:rgba(201,68,68,.08);border-color:rgba(201,68,68,.16)}

#statusBox,#oauthBox,#modelsBox{display:none}
</style>
</head>
<body>
<div class="shell">
  <header class="topbar">
    <div class="brand">
      <div class="mark" aria-hidden="true"></div>
      <div class="brand-copy">
        <strong>CodeBuddy Proxy</strong>
        <span>protocol_direct · OpenAI-compatible</span>
      </div>
    </div>
    <div class="pillrow">
      <span class="pill"><span class="dot" id="healthDot"></span><span id="healthText">检查中</span></span>
      <span class="pill">transport · <span class="mono" id="pillTransport">—</span></span>
      <span class="pill">site · <span class="mono" id="pillSite">—</span></span>
    </div>
  </header>

  <div class="hero">
    <section class="panel">
      <div class="panel-inner">
        <div class="eyebrow">Gateway Console</div>
        <h1>账号、模型与网关状态</h1>
        <p class="lede">聚焦 OAuth 接入、账号池与请求健康。配色面向未来 GEO 站点宣发，可直接沿用品牌主色。</p>
        <div class="metrics">
          <div class="metric"><div class="k">登录态</div><div class="v sm" id="mLogin">未登录</div></div>
          <div class="metric"><div class="k">Credits 余额</div><div class="v sm" id="mCredits">—</div></div>
          <div class="metric"><div class="k">启用账号</div><div class="v" id="mEnabled">0</div></div>
          <div class="metric"><div class="k">成功 / 失败</div><div class="v sm" id="mSF">0 / 0</div></div>
        </div>
        <div class="actions" style="margin-top:16px">
          <button class="ghost" id="btnRefresh" type="button">刷新状态</button>
          <button class="ghost" id="btnModels" type="button">拉取模型</button>
        </div>
      </div>
    </section>

    <section class="panel" id="codebuddy">
      <div class="panel-inner">
        <div class="section-head">
          <div>
            <h2>OAuth 登录</h2>
            <p>选择站点并开始认证，完成后账号自动进入池中。</p>
          </div>
        </div>
        <div class="field-grid">
          <div>
            <label for="site">站点</label>
            <select id="site">
              <option value="domestic">国内站 · domestic</option>
              <option value="global">国际站 · global</option>
            </select>
          </div>
          <div>
            <label for="label">账号标签</label>
            <input id="label" placeholder="例如：办公主力号" value="CodeBuddy OAuth"/>
          </div>
        </div>
        <div class="actions">
          <button class="primary" id="btnStart" type="button">开始认证</button>
          <button class="teal" id="btnPoll" type="button">检查登录</button>
          <a class="linkbtn" id="launchLink" href="#" target="_blank" rel="noreferrer">打开登录页</a>
        </div>
        <div class="oauth-status">
          <div class="title">会话状态</div>
          <div class="msg" id="oauthMsg">OAuth 会话空闲（账号池登录态见上方「登录态」；这里只反映进行中的认证流程）</div>
        </div>
        <details class="raw">
          <summary>原始 OAuth 响应</summary>
          <pre id="oauthRaw">idle</pre>
        </details>
      </div>
    </section>
  </div>


  <section class="panel" id="client-config">
      <div class="panel-inner">
        <div class="section-head">
          <div>
            <h2>OpenAI 兼容接入</h2>
            <p>给下游客户端填 Base URL + API Key；不会暴露 CodeBuddy OAuth token。</p>
          </div>
          <div class="actions">
            <button class="ghost" id="btnRefreshClient" type="button">刷新接入信息</button>
            <button class="teal" id="btnGenerateKey" type="button">生成 API Key</button>
          </div>
        </div>
        <div class="field-grid" style="grid-template-columns:1fr 1fr">
          <div>
            <label for="openAiBaseUrl">Base URL</label>
            <div class="copyline">
              <input id="openAiBaseUrl" readonly placeholder="加载中…"/>
              <button class="ghost" id="copyBaseUrl" type="button">复制</button>
            </div>
          </div>
          <div>
            <label for="openAiChatUrl">Chat Completions</label>
            <div class="copyline">
              <input id="openAiChatUrl" readonly placeholder="加载中…"/>
              <button class="ghost" id="copyChatUrl" type="button">复制</button>
            </div>
          </div>
          <div>
            <label for="openAiApiKey">API Key（网关层）</label>
            <div class="copyline">
              <input id="openAiApiKey" class="secret-input" type="password" readonly placeholder="未配置"/>
              <button class="ghost" id="copyApiKey" type="button">复制</button>
            </div>
          </div>
          <div>
            <label for="openAiModel">推荐模型</label>
            <div class="copyline">
              <input id="openAiModel" readonly value="auto"/>
              <button class="ghost" id="copyModel" type="button">复制</button>
            </div>
          </div>
        </div>
        <div class="secret-hint" id="clientConfigHint">API Key 点击复制时才会读取明文。生成的新 Key 仅当前进程生效，重启需写入 CODEBUDDY_PROXY_API_KEY。</div>
      </div>
    </section>

  <div class="stack">
    <section class="panel">
      <div class="panel-inner">
        <div class="section-head">
          <div>
            <h2>账号池</h2>
            <p>一键切换国内 / 国际号池；请求只会使用当前区域账号。</p>
          </div>
          <div class="actions">
            <div class="seg" id="poolSiteSeg" role="group" aria-label="号池区域">
              <button type="button" data-site="domestic" id="btnPoolDomestic">国内</button>
              <button type="button" data-site="global" id="btnPoolGlobal">国际</button>
            </div>
          </div>
        </div>
        <p class="seg-hint" id="poolSiteHint">当前号池：—</p>
        <div id="accounts"></div>
      </div>
    </section>

    <section class="panel">
      <div class="panel-inner">
        <div class="section-head">
          <div>
            <h2>可用模型</h2>
            <p>展示上游模型名（不含 codebuddy/ 前缀）；点击模型芯片可复制。</p>
          </div>
        </div>
        <div class="chips" id="modelChips"><div class="empty">点击「拉取模型」加载</div></div>
        <details class="raw">
          <summary>原始模型响应</summary>
          <pre id="modelsRaw">[]</pre>
        </details>
      </div>
    </section>

    <section class="panel">
      <div class="panel-inner">
        <div class="section-head">
          <div>
            <h2>运行快照</h2>
            <p>完整 status JSON，便于排障。</p>
          </div>
        </div>
        <details class="raw" open>
          <summary>status payload</summary>
          <pre id="statusRaw">loading…</pre>
        </details>
      </div>
    </section>
  </div>

  <p class="footer-note">Brand palette · sand / teal / coral · ready for GEO landing reuse</p>
</div>

<div class="toast" id="toast" role="status" aria-live="polite"></div>

<!-- hidden compatibility targets for existing helpers -->
<pre id="statusBox"></pre>
<pre id="oauthBox"></pre>
<pre id="modelsBox"></pre>

<script>
const $ = (id) => document.getElementById(id);

async function api(path, opts={}) {
  const headers = Object.assign({'content-type':'application/json'}, opts.headers||{});
  const res = await fetch(path, Object.assign({}, opts, {headers: headers, credentials:'same-origin'}));
  const text = await res.text();
  let data = {};
  try { data = text ? JSON.parse(text) : {}; } catch (e) { data = {raw:text}; }
  if (!res.ok) throw new Error((data && data.error && data.error.message) || data.message || data.error || ('HTTP '+res.status));
  return data;
}

function escapeHtml(s){
  return String(s||'').replace(/[&<>"']/g, function(c){
    return ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]);
  });
}

function fmtUptime(ms){
  ms = Number(ms||0);
  if (!ms || ms < 0) return '—';
  const s = Math.floor(ms/1000);
  const h = Math.floor(s/3600);
  const m = Math.floor((s%3600)/60);
  const r = s%60;
  if (h > 0) return h + 'h ' + m + 'm';
  if (m > 0) return m + 'm ' + r + 's';
  return r + 's';
}

function setHealth(ok, text){
  const dot = $('healthDot');
  dot.className = 'dot' + (ok ? '' : ' bad');
  $('healthText').textContent = text;
}

const usageByAccount = {};
function renderAccounts(summary, activeSite) {
  const box = $('accounts');
  const accounts = (summary && summary.accounts) || [];
  activeSite = normalizeSite(activeSite || (summary && summary.activeSite) || '');
  if (!accounts.length) {
    box.innerHTML = '<div class="empty">暂无账号。请先在右侧完成 OAuth 登录。</div>';
    return;
  }
  box.innerHTML = accounts.map(function(a) {
    const name = escapeHtml(a.userNickname || a.userName || a.userId || '未命名用户');
    const label = escapeHtml(a.label || a.id);
    const logged = a.loggedIn && a.hasCredentials;
    const site = normalizeSite(a.site);
    const inPool = !activeSite || site === activeSite;
    const usage = usageByAccount[a.id];
    let usageHtml = '<div class="usage-line"><button class="ghost" data-act="usage" data-id="' + escapeHtml(a.id) + '" type="button">查余额</button></div>';
    if (usage && usage.error) {
      usageHtml = '<div class="usage-line"><span class="pill bad" title="' + escapeHtml(usage.error) + '">查询失败</span> ' +
        '<button class="ghost" data-act="usage" data-id="' + escapeHtml(a.id) + '" type="button">重试</button></div>';
    } else if (usage && usage.credits) {
      const c = usage.credits;
      const remaining = c.unlimited ? '不限量' : (c.remaining == null ? '-' : String(c.remaining));
      const total = c.unlimited ? '不限量' : (c.total == null ? '-' : String(c.total));
      let level = 'good';
      if (!c.unlimited && Number(c.total) > 0) {
        const ratio = Number(c.remaining) / Number(c.total);
        if (ratio <= 0.05 || Number(c.remaining) <= 0) level = 'bad';
        else if (ratio <= 0.2) level = 'warn';
      }
      if (usage.notify && usage.notify.level === 'bad') level = 'bad';
      else if (usage.notify && usage.notify.level === 'warn' && level === 'good') level = 'warn';
      usageHtml = '<div class="usage-line"><span class="pill ' + level + '">' + escapeHtml(remaining + ' / ' + total) + ' Credits</span> ' +
        '<button class="ghost" data-act="usage" data-id="' + escapeHtml(a.id) + '" type="button">刷新余额</button> ' +
        '<a href="' + escapeHtml(usage.officialUsageUrl || 'https://www.codebuddy.cn/profile/plan') + '" target="_blank" rel="noreferrer">官网套餐</a></div>';
    }
    return '<div class="account">' +
      '<div class="account-top">' +
        '<div class="account-title">' +
          '<strong>' + label + '</strong>' +
          '<span class="badge site">' + escapeHtml(siteLabel(site)) + '</span>' +
          '<span class="badge ' + (inPool?'on':'muted') + '">' + (inPool?'当前号池':'其他区域') + '</span>' +
          '<span class="badge ' + (logged?'on':'off') + '">' + (logged?'已登录':'未登录') + '</span>' +
          '<span class="badge ' + (a.enabled?'on':'off') + '">' + (a.enabled?'enabled':'disabled') + '</span>' +
        '</div>' +
        '<div class="actions">' +
          '<button class="ghost" data-act="usage" data-id="' + escapeHtml(a.id) + '" type="button">查余额</button>' +
          '<button class="ghost" data-act="toggle" data-id="' + escapeHtml(a.id) + '" data-enabled="' + (a.enabled?0:1) + '" type="button">' + (a.enabled?'禁用':'启用') + '</button>' +
          '<button class="ghost" data-act="refresh" data-id="' + escapeHtml(a.id) + '" type="button">刷新 Token</button>' +
          '<button class="danger" data-act="delete" data-id="' + escapeHtml(a.id) + '" type="button">删除</button>' +
        '</div>' +
      '</div>' +
      '<div class="meta">' + name +
        ' · <span class="mono">' + escapeHtml(a.authType||'') + '</span>' +
        ' · success ' + (a.successRequests||0) + ' / fail ' + (a.failedRequests||0) +
        (a.tokenExpired ? ' · <span class="err">token expired</span>' : '') +
      '</div>' +
      usageHtml +
      (a.lastError ? ('<div class="err">' + escapeHtml(a.lastError) + '</div>') : '') +
    '</div>';
  }).join('');
  box.querySelectorAll('button[data-act]').forEach(function(btn){ btn.addEventListener('click', onAccountAction); });
}

function bareModelId(raw){
  const id = String(raw||'').trim();
  if (!id) return 'auto';
  return id.replace(/^codebuddy[/:]/i, '') || 'auto';
}
function renderModels(data){
  const chips = $('modelChips');
  const list = Array.isArray(data) ? data
    : (Array.isArray(data && data.models) ? data.models
    : (Array.isArray(data && data.data) ? data.data : []));
  if (!list.length) {
    chips.innerHTML = '<div class="empty">暂无模型数据</div>';
    return;
  }
  chips.innerHTML = list.map(function(m){
    const raw = typeof m === 'string' ? m : (m.modelId || m.upstreamId || m.id || m.name || m.model || 'model');
    const id = bareModelId(raw);
    const baseLabel = typeof m === 'string' ? id : bareModelId(m.displayName || m.name || id);
    const credits = (typeof m === 'object' && m && m.credits) ? String(m.credits) : '';
    const free = typeof m === 'object' && m && (m.free === true || /x0(\.0+)?\s*credits/i.test(credits));
    const mult = (typeof m === 'object' && m && m.creditMultiplier != null) ? m.creditMultiplier : null;
    const badge = credits ? (' · ' + (free ? '免费' : credits.replace(/\s*credits$/i,''))) : '';
    const tip = [id, credits ? ('倍率 ' + credits) : '', (m && m.description) ? m.description : ''].filter(Boolean).join(' | ');
    return '<button type="button" class="chip btnish" data-copy="' + escapeHtml(id) + '" title="' + escapeHtml(tip || '点击复制') + '">' + escapeHtml(baseLabel + badge) + (free ? ' <span class="badge on">免费</span>' : (mult!=null && credits ? ' <span class="badge site">' + escapeHtml(String(mult)+'x') + '</span>' : '')) + '</button>';
  }).join('');
  chips.querySelectorAll('[data-copy]').forEach(function(btn){
    btn.addEventListener('click', function(){ copyText(btn.getAttribute('data-copy'), '模型', btn); });
  });
}


function normalizeSite(site){
  site = String(site||'').toLowerCase().trim();
  if (site === 'domestic' || site === 'cn' || site === 'china' || site === 'internal') return 'domestic';
  return 'global';
}
function siteLabel(site){
  return normalizeSite(site) === 'domestic' ? '国内' : '国际';
}
function paintPoolSite(site, accounts){
  site = normalizeSite(site);
  const domesticBtn = $('btnPoolDomestic');
  const globalBtn = $('btnPoolGlobal');
  if (domesticBtn) domesticBtn.className = site === 'domestic' ? 'active' : '';
  if (globalBtn) globalBtn.className = site === 'global' ? 'active' : '';
  const domestic = accounts && accounts.domesticCount != null ? accounts.domesticCount : '—';
  const global = accounts && accounts.globalCount != null ? accounts.globalCount : '—';
  const activeEnabled = accounts && accounts.activeEnabledCount != null ? accounts.activeEnabledCount : '—';
  if ($('poolSiteHint')) {
    $('poolSiteHint').textContent = '当前号池：' + siteLabel(site) + ' · 可用启用账号 ' + activeEnabled + ' · 国内账号 ' + domestic + ' / 国际账号 ' + global + '（仅当前区域会参与请求）';
  }
}
async function switchPoolSite(site){
  site = normalizeSite(site);
  const data = await api('/direct-admin/api/pool-site', {method:'POST', body: JSON.stringify({site: site})});
  paintStatus(data);
  if (data.note) showToast(data.note);
  // Refresh models for the new region primary account.
  refreshModels().catch(function(){});
  return data;
}

function paintStatus(data){
  const stats = data.stats || {};
  const accounts = data.accounts || {};
  const cfg = data.config || {};
  const enabledCount = accounts.enabledCount != null ? accounts.enabledCount : ((accounts.accounts||[]).filter(function(a){return a.enabled;}).length);
  $('mEnabled').textContent = String(enabledCount);
  $('mSF').textContent = (stats.successRequests||0) + ' / ' + (stats.failedRequests||0);
  const loggedIn = !!accounts.loggedIn;
  const primary = accounts.primary || ((accounts.accounts||[])[0] || null);
  $('mLogin').textContent = loggedIn
    ? ('已登录' + (primary && (primary.userNickname || primary.userName || primary.userId) ? (' · ' + (primary.userNickname || primary.userName || primary.userId)) : ''))
    : '未登录';
  $('pillTransport').textContent = data.transport || 'protocol_direct';
  const poolSite = normalizeSite(data.poolSite || cfg.poolSite || cfg.site || 'global');
  $('pillSite').textContent = siteLabel(poolSite);
  paintPoolSite(poolSite, accounts);
  // Keep OAuth login site aligned with active pool by default.
  if ($('site') && !$('site').dataset.userTouched) {
    $('site').value = poolSite === 'domestic' ? 'domestic' : 'global';
  }
  const activeEnabled = accounts.activeEnabledCount != null ? accounts.activeEnabledCount : enabledCount;
  $('mEnabled').textContent = String(activeEnabled);
  setHealth(!!data.ok, data.ok ? (loggedIn ? ('服务正常 · ' + siteLabel(poolSite) + '号池') : ('服务正常 · ' + siteLabel(poolSite) + '号池未登录')) : '状态异常');
  // Auto-fetch balance for primary account once.
  if (primary && primary.id && primary.hasCredentials && !usageByAccount[primary.id] && !paintStatus._usageKick) {
    paintStatus._usageKick = true;
    fetchAccountUsage(primary.id, true).catch(function(){});
  }
  const usage = primary && usageByAccount[primary.id];
  if (usage && usage.credits) {
    const c = usage.credits;
    $('mCredits').textContent = c.unlimited ? '不限量' : ((c.remaining == null ? '-' : c.remaining) + ' / ' + (c.total == null ? '-' : c.total));
  } else if (usage && usage.error) {
    $('mCredits').textContent = '查询失败';
  } else {
    $('mCredits').textContent = loggedIn ? '点击查余额' : '—';
  }
  const raw = JSON.stringify(data, null, 2);
  $('statusBox').textContent = raw;
  $('statusRaw').textContent = raw;
  renderAccounts(accounts, poolSite);
}

function paintOAuth(data){
  const session = (data && data.session) || {};
  const login = (data && data.login) || {};
  const status = session.status || login.status || (data && data.ok ? 'ok' : 'idle');
  const err = session.error || login.message || '';
  let msg = '状态：' + status;
  if (session.label) msg += ' · ' + session.label;
  if (session.site) msg += ' · ' + session.site;
  if (err) msg += ' · ' + err;
  $('oauthMsg').textContent = msg;
  const raw = JSON.stringify(data, null, 2);
  $('oauthBox').textContent = raw;
  $('oauthRaw').textContent = raw;
}


let clientConfig = { apiKey: '', baseUrl: '', chatCompletionsUrl: '', recommendedModel: 'auto' };
function showToast(msg, kind){
  const el = $('toast');
  el.textContent = msg;
  el.className = 'toast show' + (kind === 'error' ? ' err' : '');
  clearTimeout(showToast._t);
  showToast._t = setTimeout(function(){ el.className = 'toast'; }, 2200);
}
async function copyText(text, label, button){
  const value = String(text||'');
  if (!value) { showToast((label||'内容') + '为空', 'error'); return; }
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(value);
    } else {
      const ta = document.createElement('textarea');
      ta.value = value; document.body.appendChild(ta); ta.select();
      document.execCommand('copy'); ta.remove();
    }
    const old = button ? button.textContent : '';
    if (button) button.textContent = '已复制';
    showToast((label||'内容') + ' 已复制');
    if (button) setTimeout(function(){ button.textContent = old; }, 1200);
  } catch (e) {
    showToast('复制失败：' + (e && e.message ? e.message : e), 'error');
  }
}
function paintClientConfig(cfg){
  clientConfig = cfg || clientConfig;
  const base = cfg.baseUrl || cfg.apiBase || '';
  const chat = cfg.chatCompletionsUrl || (base ? (base.replace(/\/$/,'') + '/chat/completions') : '');
  const model = bareModelId(cfg.recommendedModel || 'auto');
  $('openAiBaseUrl').value = base;
  $('openAiChatUrl').value = chat;
  $('openAiModel').value = model;
  const configured = !!cfg.apiKeyConfigured || !!cfg.apiKey;
  $('openAiApiKey').value = configured ? (cfg.apiKeyPreview || '已配置 · 点击复制') : '';
  $('openAiApiKey').placeholder = configured ? '' : 'API Key 未配置';
  $('copyApiKey').disabled = !configured;
  if (cfg.note) $('clientConfigHint').textContent = cfg.note;
  clientConfig.apiKey = cfg.apiKey || '';
  clientConfig.baseUrl = base;
  clientConfig.chatCompletionsUrl = chat;
  clientConfig.recommendedModel = model;
}
async function refreshClientConfig(){
  const data = await api('/direct-admin/api/client-config');
  paintClientConfig(data);
  return data;
}
async function generateApiKey(){
  if (!confirm('生成新的网关 API Key？\n\n旧 Key 会立即失效，新 Key 会写入 .env 并在重启后继续生效。\n请同步更新 ZCode / NewAPI 等客户端配置。')) return;
  const data = await api('/direct-admin/api/client-config/generate-key', {method:'POST', body:'{}'});
  paintClientConfig(data);
  if (data.apiKey) await copyText(data.apiKey, '新 API Key', $('btnGenerateKey'));
  if (data.note) showToast(data.note);
}

async function refreshStatus(){
  const data = await api('/direct-admin/api/status');
  paintStatus(data);
}

async function refreshModels(){
  const data = await api('/direct-admin/api/codebuddy/models?fresh=1');
  const raw = JSON.stringify(data, null, 2);
  $('modelsBox').textContent = raw;
  $('modelsRaw').textContent = raw;
  renderModels(data);
}

async function startOAuth(){
  $('oauthMsg').textContent = '正在创建登录会话…';
  const data = await api('/direct-admin/api/codebuddy/oauth/start', {
    method:'POST',
    body: JSON.stringify({site:$('site').value, label:$('label').value, reuseExisting:false})
  });
  paintOAuth(data);
  const url = (data.login && data.login.url) || (data.session && data.session.url) || '#';
  const launch = (data.login && data.login.launchUrl) || (data.session && data.session.launchUrl) || url;
  $('launchLink').href = launch;
  if (url && url !== '#') window.open(url, '_blank', 'noopener');
}

async function pollOAuth(){
  $('oauthMsg').textContent = '正在检查登录…';
  const data = await api('/direct-admin/api/codebuddy/oauth/poll', {method:'POST', body:'{}'});
  paintOAuth(data);
  await refreshStatus();
}

async function fetchAccountUsage(accountId, silent){
  try {
    const result = await api('/direct-admin/api/codebuddy/accounts/'+encodeURIComponent(accountId)+'/usage');
    usageByAccount[accountId] = result;
    if (!silent) {
      const credits = result.credits || {};
      showToast(credits.display || credits.label || '余额已更新');
    }
  } catch (e) {
    usageByAccount[accountId] = { error: e.message || String(e) };
    if (!silent) showToast('查余额失败：' + (e.message || e), 'error');
  }
  // refresh UI from last status snapshot by re-pulling status
  await refreshStatus();
}
async function onAccountAction(ev){
  const btn = ev.currentTarget;
  const id = btn.getAttribute('data-id');
  const act = btn.getAttribute('data-act');
  if (act === 'usage') {
    await fetchAccountUsage(id, false);
    return;
  }
  if (act === 'delete') {
    if (!confirm('确认删除该账号？')) return;
    await api('/direct-admin/api/codebuddy/accounts/'+encodeURIComponent(id), {method:'DELETE'});
    delete usageByAccount[id];
  } else if (act === 'toggle') {
    const enabled = btn.getAttribute('data-enabled') === '1';
    await api('/direct-admin/api/codebuddy/accounts/'+encodeURIComponent(id)+'/'+(enabled?'enable':'disable'), {method:'POST', body:'{}'});
  } else if (act === 'refresh') {
    await api('/direct-admin/api/codebuddy/accounts/'+encodeURIComponent(id)+'/refresh-token', {method:'POST', body:'{}'});
  }
  await refreshStatus();
}

if ($('btnPoolDomestic')) $('btnPoolDomestic').onclick = function(){ switchPoolSite('domestic').catch(function(e){ showToast(e.message, 'error'); }); };
if ($('btnPoolGlobal')) $('btnPoolGlobal').onclick = function(){ switchPoolSite('global').catch(function(e){ showToast(e.message, 'error'); }); };
if ($('site')) $('site').addEventListener('change', function(){ $('site').dataset.userTouched = '1'; });
$('btnRefresh').onclick = function(){ refreshStatus().catch(function(e){ $('statusRaw').textContent = e.message; setHealth(false, '刷新失败'); }); };
$('btnModels').onclick = function(){ refreshModels().catch(function(e){ $('modelsRaw').textContent = e.message; $('modelChips').innerHTML = '<div class="empty">' + escapeHtml(e.message) + '</div>'; }); };
$('btnStart').onclick = function(){ startOAuth().catch(function(e){ $('oauthMsg').textContent = e.message; $('oauthRaw').textContent = e.message; }); };
$('btnPoll').onclick = function(){ pollOAuth().catch(function(e){ $('oauthMsg').textContent = e.message; $('oauthRaw').textContent = e.message; }); };
$('btnRefreshClient').onclick = function(){ refreshClientConfig().catch(function(e){ showToast(e.message, 'error'); }); };
$('btnGenerateKey').onclick = function(){ generateApiKey().catch(function(e){ showToast(e.message, 'error'); }); };
$('copyBaseUrl').onclick = function(){ copyText($('openAiBaseUrl').value, 'Base URL', $('copyBaseUrl')); };
$('copyChatUrl').onclick = function(){ copyText($('openAiChatUrl').value, 'Chat Completions', $('copyChatUrl')); };
$('copyModel').onclick = function(){ copyText($('openAiModel').value, '模型', $('copyModel')); };
$('copyApiKey').onclick = function(){
  refreshClientConfig().then(function(cfg){
    if (!cfg.apiKey) { showToast('API Key 未配置', 'error'); return; }
    $('openAiApiKey').value = cfg.apiKeyPreview || '已配置 · 点击复制';
    return copyText(cfg.apiKey, 'API Key', $('copyApiKey'));
  }).catch(function(e){ showToast(e.message, 'error'); });
};
refreshStatus().catch(function(e){ $('statusRaw').textContent = e.message; setHealth(false, '无法连接'); });
refreshClientConfig().catch(function(e){ showToast(e.message, 'error'); });
refreshModels().catch(function(){});
setInterval(function(){ refreshStatus().catch(function(){}); }, 15000);
</script>
</body>
</html>`
}

func LaunchPage(message string, success bool) string {
	tone := "bad"
	title := "认证未完成"
	if success {
		tone = "ok"
		title = "认证成功"
	}
	return `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width, initial-scale=1"/>
<title>CodeBuddy OAuth</title>
<style>
:root{--sand:#F4EEE4;--paper:#FFFCF7;--ink:#1B2430;--muted:#667384;--teal:#0F7C74;--coral:#E86F3A;--ok:#1F8A5B;--bad:#C94444}
*{box-sizing:border-box}
body{
  margin:0;min-height:100dvh;display:grid;place-items:center;padding:24px;
  font-family:system-ui,-apple-system,"Segoe UI",sans-serif;color:var(--ink);
  background:
    radial-gradient(700px 420px at 12% 0%, rgba(15,124,116,.18), transparent 55%),
    radial-gradient(640px 380px at 100% 10%, rgba(232,111,58,.16), transparent 50%),
    var(--sand);
}
.card{
  width:min(520px,100%);border-radius:28px;padding:4px;
  background:linear-gradient(180deg, rgba(255,255,255,.8), rgba(255,255,255,.35));
  border:1px solid rgba(255,255,255,.6);box-shadow:0 18px 40px -24px rgba(27,36,48,.35);
}
.inner{border-radius:24px;background:var(--paper);border:1px solid rgba(27,36,48,.1);padding:28px 26px}
.mark{width:42px;height:42px;border-radius:16px;margin-bottom:16px;background:linear-gradient(135deg,var(--teal),#3B82A0 55%,var(--coral))}
h1{margin:0 0 10px;font-size:28px;letter-spacing:-.03em}
p{margin:0 0 18px;line-height:1.55;color:var(--muted)}
.ok{color:var(--ok);font-weight:600}
.bad{color:var(--bad);font-weight:600}
a{
  display:inline-flex;align-items:center;justify-content:center;
  padding:11px 16px;border-radius:999px;text-decoration:none;font-weight:600;font-size:14px;
  color:#fff;background:linear-gradient(180deg,#1A9A90,var(--teal));
}
</style>
</head>
<body>
<div class="card"><div class="inner">
  <div class="mark" aria-hidden="true"></div>
  <h1>` + html.EscapeString(title) + `</h1>
  <p class="` + tone + `">` + html.EscapeString(message) + `</p>
  <a href="/direct-admin/#codebuddy">返回管理台</a>
</div></div>
</body>
</html>`
}

func Compact(value string) string { return strings.TrimSpace(value) }
