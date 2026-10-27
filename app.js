'use strict';

/* =========================================================
   Настройки
   ========================================================= */
const CONFIG = {
  API_BASE: '/api/v1',   // префикс API на Go-бэкенде
  REFRESH_MS: 30000,     // как часто обновлять список
  // ?mock=1 в адресе — принудительно демо-данные.
  // Если API недоступен при первой загрузке — демо включится само.
  FORCE_MOCK: new URLSearchParams(location.search).has('mock'),
};

/* =========================================================
   DOM и состояние
   ========================================================= */
const $ = (s) => document.querySelector(s);
const els = {
  list: $('#list'), summary: $('#summary'), updatedAt: $('#updatedAt'), demo: $('#demoBadge'),
  refresh: $('#refreshBtn'), filter: $('#filterBtn'),
  modal: $('#modal'), mDot: $('#mDot'), mTitle: $('#mTitle'), mUrl: $('#mUrl'),
  mUptime: $('#mUptime'), mAvg: $('#mAvg'), mIncidents: $('#mIncidents'),
  chart: $('#chart'), strip: $('#strip'),
};

const state = {
  sites: [],
  loaded: false,
  demo: CONFIG.FORCE_MOCK,
  onlyDown: false,
  openId: null,
  history: null,
  lastFocus: null,
};

/* =========================================================
   API
   ========================================================= */
async function request(path) {
  const res = await fetch(CONFIG.API_BASE + path, { headers: { Accept: 'application/json' } });
  if (!res.ok) throw new Error('HTTP ' + res.status);
  return res.json();
}

async function fetchSites() {
  if (!state.demo) {
    try {
      const data = await request('/sites');
      return Array.isArray(data) ? data : data.sites || [];
    } catch (e) {
      if (state.loaded) throw e;            // бэкенд упал во время работы — не подменяем данные
      console.warn('API недоступен, показываю демо-данные:', e.message);
      state.demo = true;
    }
  }
  return mock.sites();
}

function fetchHistory(id) {
  if (state.demo) return mock.history(id);
  return request(`/sites/${encodeURIComponent(id)}/history?period=24h`);
}

/* =========================================================
   Утилиты
   ========================================================= */
const esc = (s) => String(s ?? '').replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[c]));

function plural(n, one, few, many) {
  const a = n % 10, b = n % 100;
  if (a === 1 && b !== 11) return one;
  if (a >= 2 && a <= 4 && (b < 12 || b > 14)) return few;
  return many;
}

const isUp = (s) => s.status === 'up';
const fmtMs = (ms) => (Number.isFinite(ms) ? `${ms} мс` : '—');
const fmtTime = (t) => new Date(t).toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
const fmtPct = (v) => (v >= 100 ? '100%' : v.toFixed(2).replace('.', ',') + '%');
const STATUS_TEXT = { up: 'работает', down: 'лежит' };

function speedClass(site) {
  if (!isUp(site)) return 't-down';
  const ms = site.response_time_ms;
  return ms < 300 ? 't-fast' : ms < 1000 ? 't-mid' : 't-slow';
}

function niceCeil(v) {
  const p = 10 ** Math.floor(Math.log10(v));
  const f = v / p;
  return (f <= 1 ? 1 : f <= 2 ? 2 : f <= 5 ? 5 : 10) * p;
}

/* =========================================================
   Список сайтов
   ========================================================= */
function rowHtml(s, i) {
  const st = s.status === 'up' || s.status === 'down' ? s.status : 'unknown';
  return `
    <li>
      <button class="row" type="button" data-id="${esc(s.id)}" style="--i:${i}">
        <span class="dot is-${st}"></span>
        <span class="row__info">
          <span class="row__name">${esc(s.name)}</span>
          <span class="row__url">${esc(s.url)}</span>
        </span>
        <span class="row__status">${STATUS_TEXT[st] || 'нет данных'}</span>
        <span class="row__time ${speedClass(s)}">${isUp(s) ? fmtMs(s.response_time_ms) : '—'}</span>
        <svg class="row__chev" viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M9 6l6 6-6 6"/></svg>
      </button>
    </li>`;
}

function renderList(intro = false) {
  // лежащие — наверх, дальше по алфавиту
  const sorted = [...state.sites].sort((a, b) => isUp(a) - isUp(b) || String(a.name).localeCompare(b.name, 'ru'));
  const items = state.onlyDown ? sorted.filter((s) => !isUp(s)) : sorted;

  els.list.classList.toggle('is-intro', intro);
  els.list.innerHTML = items.length
    ? items.map(rowHtml).join('')
    : `<li class="empty">${state.onlyDown ? 'Все сайты работают' : 'Сайтов пока нет'}</li>`;

  const total = state.sites.length;
  const down = state.sites.filter((s) => s.status === 'down').length;
  els.summary.innerHTML =
    `${total} ${plural(total, 'сайт', 'сайта', 'сайтов')} · ` +
    `<b class="c-up">${total - down}</b> работают · ` +
    `<b class="${down ? 'c-down' : ''}">${down}</b> ${plural(down, 'лежит', 'лежат', 'лежат')}`;
}

async function loadSites() {
  els.refresh.classList.add('is-loading');
  try {
    const sites = await fetchSites();
    const intro = !state.loaded;
    state.sites = sites;
    state.loaded = true;
    renderList(intro);
    els.updatedAt.classList.remove('is-error');
    els.updatedAt.textContent = 'Обновлено в ' + new Date().toLocaleTimeString('ru-RU');
  } catch (e) {
    els.updatedAt.classList.add('is-error');
    els.updatedAt.textContent = 'Нет связи с сервером, показаны данные на ' + new Date().toLocaleTimeString('ru-RU');
  } finally {
    els.demo.hidden = !state.demo;
    setTimeout(() => els.refresh.classList.remove('is-loading'), 400);
  }
}

/* =========================================================
   Модальное окно
   ========================================================= */
async function openModal(id) {
  const site = state.sites.find((s) => String(s.id) === id);
  if (!site) return;

  state.openId = id;
  state.history = null;
  state.lastFocus = document.activeElement;

  els.mDot.className = `dot is-${site.status}`;
  els.mTitle.textContent = site.name;
  els.mUrl.textContent = site.url;
  els.mUrl.href = site.url;
  els.mUptime.textContent = els.mAvg.textContent = els.mIncidents.textContent = '—';
  els.mUptime.className = 'stat__value';
  els.chart.innerHTML = '<div class="chart__msg">Загрузка…</div>';
  els.strip.innerHTML = '';

  els.modal.classList.add('is-open');
  els.modal.setAttribute('aria-hidden', 'false');
  document.body.style.overflow = 'hidden';
  els.modal.querySelector('.modal__close').focus({ preventScroll: true });

  try {
    const h = await fetchHistory(id);
    if (state.openId !== id) return;      // пока грузили — открыли другой сайт или закрыли окно
    state.history = h;
    renderHistory(h);
  } catch (e) {
    if (state.openId === id) els.chart.innerHTML = '<div class="chart__msg">Не удалось загрузить историю</div>';
  }
}

function closeModal() {
  if (!state.openId) return;
  state.openId = null;
  els.modal.classList.remove('is-open');
  els.modal.setAttribute('aria-hidden', 'true');
  document.body.style.overflow = '';
  state.lastFocus?.focus?.({ preventScroll: true });
}

function renderHistory(h) {
  const pts = h.points || [];
  if (!pts.length) {
    els.chart.innerHTML = '<div class="chart__msg">Нет данных за последние 24 часа</div>';
    return;
  }

  const known = pts.filter((p) => p.status === 'up' || p.status === 'down');
  const up = known.filter((p) => p.status === 'up' && Number.isFinite(p.response_time_ms));
  const uptime = known.length ? (known.filter((p) => p.status === 'up').length / known.length) * 100 : 0;
  const avg = up.length ? Math.round(up.reduce((s, p) => s + p.response_time_ms, 0) / up.length) : null;
  let incidents = 0;
  pts.forEach((p, i) => { if (p.status === 'down' && pts[i - 1]?.status !== 'down') incidents++; });

  els.mUptime.textContent = known.length ? fmtPct(uptime) : '—';
  els.mUptime.className = 'stat__value ' + (uptime >= 99.9 ? 't-fast' : uptime >= 95 ? 't-mid' : 't-down');
  els.mAvg.textContent = fmtMs(avg);
  els.mIncidents.textContent = incidents;

  drawChart(pts);
  drawStrip(pts);
}

/* ---------- График времени ответа (SVG) ---------- */
function drawChart(pts) {
  const W = Math.max(els.chart.clientWidth, 260);
  const H = els.chart.clientHeight || 200;
  const P = { t: 8, r: 6, b: 22, l: 40 };
  const iw = W - P.l - P.r;
  const ih = H - P.t - P.b;
  const n = pts.length;

  const vals = pts.map((p) => (p.status === 'up' && Number.isFinite(p.response_time_ms) ? p.response_time_ms : null));
  const maxV = niceCeil(Math.max(100, ...vals.filter((v) => v != null)));
  const x = (i) => P.l + (n > 1 ? (i / (n - 1)) * iw : iw / 2);
  const y = (v) => P.t + ih - (v / maxV) * ih;
  const step = n > 1 ? iw / (n - 1) : iw;
  const f = (v) => v.toFixed(1);

  // сетка и подписи по Y
  let grid = '';
  for (let k = 0; k <= 4; k++) {
    const v = (maxV * k) / 4;
    grid += `<line class="grid" x1="${P.l}" x2="${W - P.r}" y1="${f(y(v))}" y2="${f(y(v))}"/>` +
            `<text class="axis" x="${P.l - 8}" y="${f(y(v) + 4)}" text-anchor="end">${Math.round(v)}</text>`;
  }

  // подписи по X
  let xLabels = '';
  for (let k = 0; k <= 4; k++) {
    const i = Math.round((k * (n - 1)) / 4);
    const anchor = k === 0 ? 'start' : k === 4 ? 'end' : 'middle';
    xLabels += `<text class="axis" x="${f(x(i))}" y="${H - 4}" text-anchor="${anchor}">${fmtTime(pts[i].t)}</text>`;
  }

  // красные полосы — периоды, когда сайт лежал
  let bands = '';
  for (let i = 0; i < n; i++) {
    if (pts[i].status !== 'down') continue;
    let j = i;
    while (j + 1 < n && pts[j + 1].status === 'down') j++;
    const x0 = Math.max(P.l, x(i) - step / 2);
    const x1 = Math.min(W - P.r, x(j) + step / 2);
    bands += `<rect class="band" x="${f(x0)}" y="${P.t}" width="${f(Math.max(x1 - x0, 2))}" height="${ih}" rx="2"/>`;
    i = j;
  }

  // линия рвётся в местах, где сайт лежал
  const segs = [];
  let cur = null;
  vals.forEach((v, i) => {
    if (v == null) { cur = null; return; }
    if (!cur) { cur = []; segs.push(cur); }
    cur.push([x(i), y(v)]);
  });
  const base = P.t + ih;
  const lines = segs.map((s) => `<path class="line" pathLength="1" d="M${s.map(([a, b]) => `${f(a)},${f(b)}`).join('L')}"/>`).join('');
  const areas = segs.map((s) =>
    `<path class="area" fill="url(#areaGrad)" d="M${f(s[0][0])},${base}L${s.map(([a, b]) => `${f(a)},${f(b)}`).join('L')}L${f(s[s.length - 1][0])},${base}Z"/>`
  ).join('');

  els.chart.innerHTML = `
    <svg width="${W}" height="${H}" viewBox="0 0 ${W} ${H}">
      <defs>
        <linearGradient id="areaGrad" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0" stop-color="currentColor" stop-opacity=".18"/>
          <stop offset="1" stop-color="currentColor" stop-opacity="0"/>
        </linearGradient>
      </defs>
      ${grid}${bands}${areas}${lines}${xLabels}
      <line class="guide" y1="${P.t}" y2="${base}"/>
      <circle class="hover-dot" r="4"/>
      <rect class="hit" x="${P.l}" y="0" width="${iw}" height="${base}" fill="transparent"/>
    </svg>
    <div class="tip" hidden></div>`;

  // подсказка при наведении
  const svg = els.chart.querySelector('svg');
  const hit = svg.querySelector('.hit');
  const guide = svg.querySelector('.guide');
  const dot = svg.querySelector('.hover-dot');
  const tip = els.chart.querySelector('.tip');

  const move = (e) => {
    const rect = svg.getBoundingClientRect();
    const scale = W / rect.width;
    const mx = (e.clientX - rect.left) * scale;
    const i = Math.max(0, Math.min(n - 1, Math.round(((mx - P.l) / iw) * (n - 1))));
    const v = vals[i];
    const xx = x(i);

    guide.setAttribute('x1', xx);
    guide.setAttribute('x2', xx);
    guide.style.opacity = 1;
    if (v != null) {
      dot.setAttribute('cx', xx);
      dot.setAttribute('cy', y(v));
      dot.style.opacity = 1;
    } else {
      dot.style.opacity = 0;
    }

    tip.innerHTML = `<b>${v != null ? v + ' мс' : 'недоступен'}</b><span>${fmtTime(pts[i].t)}</span>`;
    tip.classList.toggle('is-down', v == null);
    tip.hidden = false;
    const half = tip.offsetWidth / 2;
    const left = Math.max(half, Math.min(rect.width - half, xx / scale));
    tip.style.left = left + 'px';
  };
  const leave = () => { guide.style.opacity = 0; dot.style.opacity = 0; tip.hidden = true; };

  hit.addEventListener('pointermove', move);
  hit.addEventListener('pointerdown', move);
  hit.addEventListener('pointerleave', leave);
}

/* ---------- Полоса аптайма (48 отрезков по 30 минут) ---------- */
function drawStrip(pts) {
  const B = Math.min(48, pts.length);
  const size = pts.length / B;
  let html = '';
  for (let b = 0; b < B; b++) {
    const slice = pts.slice(Math.floor(b * size), Math.floor((b + 1) * size));
    if (!slice.length) continue;
    const known = slice.filter((p) => p.status === 'up' || p.status === 'down');
    const u = known.length ? (known.filter((p) => p.status === 'up').length / known.length) * 100 : null;
    const cls = u == null ? '' : u === 100 ? 'is-up' : u >= 90 ? 'is-warn' : 'is-down';
    const title = `${fmtTime(slice[0].t)}–${fmtTime(slice[slice.length - 1].t)} · ${u == null ? 'нет данных' : fmtPct(u)}`;
    html += `<span class="seg ${cls}" style="--i:${b}" title="${esc(title)}"></span>`;
  }
  els.strip.innerHTML = html;
}

/* =========================================================
   Демо-данные (пока нет бэкенда)
   ========================================================= */
const mock = (() => {
  const SITES = [
    { id: 'main', name: 'Основной сайт', url: 'https://example.com', base: 120, flaky: 0.5 },
    { id: 'api', name: 'API', url: 'https://api.example.com/health', base: 45, flaky: 0 },
    { id: 'shop', name: 'Интернет-магазин', url: 'https://shop.example.com', base: 640, flaky: 2 },
    { id: 'blog', name: 'Блог', url: 'https://blog.example.com', base: 1150, flaky: 1 },
    { id: 'docs', name: 'Документация', url: 'https://docs.example.com', base: 140, flaky: 0 },
    { id: 'cdn', name: 'CDN', url: 'https://cdn.example.com', base: 28, flaky: 0 },
    { id: 'crm', name: 'CRM', url: 'https://crm.example.org', base: 310, flaky: 1 },
    { id: 'legacy', name: 'Старый портал', url: 'https://legacy.example.org', base: 880, flaky: -1 },
    { id: 'auth', name: 'Авторизация', url: 'https://auth.example.com', base: 95, flaky: 0.5 },
    { id: 'status', name: 'Статус-страница', url: 'https://status.example.net', base: 70, flaky: 0 },
    { id: 'mail', name: 'Почта', url: 'https://mail.example.net', base: 410, flaky: 1 },
    { id: 'wiki', name: 'Вики', url: 'https://wiki.example.net', base: 520, flaky: 0.5 },
    { id: 'files', name: 'Файлообменник', url: 'https://files.example.org', base: 760, flaky: 2 },
    { id: 'pay', name: 'Платёжный шлюз', url: 'https://pay.example.com', base: 230, flaky: 0 },
  ];
  const cache = new Map();
  const sleep = (ms) => new Promise((r) => setTimeout(r, ms));

  function rng(seed) {
    return () => {
      seed = (seed + 0x6d2b79f5) | 0;
      let t = Math.imul(seed ^ (seed >>> 15), 1 | seed);
      t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t;
      return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
    };
  }
  const hash = (s) => [...s].reduce((h, c) => Math.imul(h ^ c.charCodeAt(0), 16777619), 2166136261);

  function build(site) {
    const r = rng(hash(site.id));
    const STEP = 300, N = 288;                                  // 5 минут × 288 = 24 часа
    const now = Math.floor(Date.now() / 1000 / STEP) * STEP;
    const down = new Array(N).fill(false);

    const incidents = site.flaky < 0 ? 1 : Math.round(site.flaky * r() * 2);
    for (let k = 0; k < incidents; k++) {
      const start = Math.floor(r() * (N - 12));
      const len = 1 + Math.floor(r() * 7);
      for (let i = start; i < start + len && i < N; i++) down[i] = true;
    }
    if (site.flaky < 0) for (let i = N - 4 - Math.floor(r() * 20); i < N; i++) down[i] = true; // лежит прямо сейчас

    const points = down.map((d, i) => {
      const t = (now - (N - 1 - i) * STEP) * 1000;
      if (d) return { t: new Date(t).toISOString(), status: 'down', response_time_ms: null };
      const wave = 1 + 0.25 * Math.sin((t / 86400000) * Math.PI * 2);
      const spike = r() < 0.025 ? 1.8 + r() * 2 : 1;
      return { t: new Date(t).toISOString(), status: 'up', response_time_ms: Math.round(site.base * wave * (0.8 + r() * 0.4) * spike) };
    });

    return { site_id: site.id, from: points[0].t, to: points[N - 1].t, step_seconds: STEP, points };
  }

  const historyOf = (id) => {
    if (!cache.has(id)) cache.set(id, build(SITES.find((s) => s.id === id)));
    return cache.get(id);
  };

  return {
    async sites() {
      await sleep(300);
      return SITES.map((s) => {
        const last = historyOf(s.id).points.at(-1);
        return {
          id: s.id, name: s.name, url: s.url, status: last.status,
          response_time_ms: last.status === 'up' ? Math.round(last.response_time_ms * (0.9 + Math.random() * 0.2)) : null,
          checked_at: new Date().toISOString(),
        };
      });
    },
    async history(id) {
      await sleep(350);
      return historyOf(id);
    },
  };
})();

/* =========================================================
   Запуск
   ========================================================= */
els.list.addEventListener('click', (e) => {
  const row = e.target.closest('.row');
  if (row) openModal(row.dataset.id);
});

els.modal.addEventListener('click', (e) => {
  if (e.target.closest('[data-close]')) closeModal();
});

document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape') closeModal();
});

els.refresh.addEventListener('click', loadSites);

els.filter.addEventListener('click', () => {
  state.onlyDown = !state.onlyDown;
  els.filter.classList.toggle('is-active', state.onlyDown);
  renderList(true);
});

let resizeTimer;
window.addEventListener('resize', () => {
  clearTimeout(resizeTimer);
  resizeTimer = setTimeout(() => { if (state.openId && state.history) drawChart(state.history.points); }, 120);
});

loadSites();
setInterval(() => { if (!document.hidden) loadSites(); }, CONFIG.REFRESH_MS);
