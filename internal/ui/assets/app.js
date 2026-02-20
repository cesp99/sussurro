'use strict';

// ---- Bootstrap ----
document.addEventListener('DOMContentLoaded', async () => {
  await reloadSettings();
});

// Re-fetch data and re-render in place — never calls location.reload() which
// destroys the WebKit JS context mid-execution and causes a crash.
async function reloadSettings() {
  try {
    const raw = await window.getInitialData();
    const data = JSON.parse(raw);
    render(data);
  } catch (e) {
    console.error('reloadSettings failed:', e);
  }
}

// ---- Render ----
function render(data) {
  // Status bar labels
  document.getElementById('platform-label').textContent = data.platform;
  document.getElementById('version-label').textContent  = `v${data.version}`;

  // Models
  const whisperItems = data.models.filter(m => m.type === 'whisper');
  const llmItems     = data.models.filter(m => m.type === 'llm');
  renderModelList('whisper-list', whisperItems, 'whisper');
  renderModelList('llm-list',     llmItems,     'llm');

  // Hotkey
  renderHotkey(data.hotkey, data.isWayland);
}

// ---- Model list ----
function renderModelList(containerId, models, groupName) {
  const container = document.getElementById(containerId);
  container.innerHTML = '';

  models.forEach(m => {
    const item = document.createElement('div');
    item.className = 'model-item' + (m.active ? ' active' : '');
    item.dataset.id = m.id;

    item.innerHTML = `
      <input type="radio" name="${groupName}" value="${m.id}" ${m.active ? 'checked' : ''}>
      <div class="model-info">
        <span class="model-name">${m.name}${m.active ? '<span class="model-badge">ACTIVE</span>' : ''}</span>
        <span class="model-desc">${m.desc}</span>
        <span class="model-size">${m.size}</span>
      </div>
      <div class="model-status" id="status-${m.id}">
        ${m.installed ? installedBadge() : downloadArea(m.id)}
      </div>
    `;

    const radio = item.querySelector('input[type="radio"]');

    // LLM has only one model — disable the radio, nothing to switch to
    if (m.type === 'llm') {
      radio.disabled = true;
    } else {
      radio.addEventListener('change', async () => {
        if (!radio.checked) return;
        if (!m.installed) { radio.checked = false; return; }

        const res = await window.setActiveModel(m.id);
        if (res.startsWith('error')) { radio.checked = false; return; }

        // Config written — app will restart momentarily; disable all inputs
        lockUI();
      });
    }

    container.appendChild(item);

    // Attach download handler
    if (!m.installed) {
      const btn = item.querySelector('.download-btn');
      if (btn) btn.addEventListener('click', e => { e.stopPropagation(); startDownload(m.id, m.name); });
    }
  });
}

// Disable all interactive elements while the app restarts
function lockUI() {
  document.querySelectorAll('button, input').forEach(el => el.disabled = true);
  const st = document.querySelector('.status-text');
  if (st) st.textContent = 'Restarting…';
}

function installedBadge() {
  return `<span class="installed-badge">✓ Installed</span>`;
}

function downloadArea(id) {
  return `
    <div class="download-area">
      <button class="download-btn" id="btn-${id}">↓ Download</button>
      <div class="dl-progress-wrap" id="prog-wrap-${id}" hidden>
        <progress class="dl-progress" id="prog-${id}" value="0" max="1"></progress>
        <span class="dl-progress-label" id="pct-${id}">0%</span>
      </div>
    </div>
  `;
}

function startDownload(modelId, modelName) {
  const btn      = document.getElementById(`btn-${modelId}`);
  const progWrap = document.getElementById(`prog-wrap-${modelId}`);

  // Show progress, hide button — never show both at once
  if (btn)      btn.hidden      = true;
  if (progWrap) progWrap.hidden = false;

  window.downloadModel(modelId);
}

// Called from Go via webview.Eval
window.onDownloadProgress = function(name, percent) {
  document.querySelectorAll('.model-item').forEach(item => {
    const nameEl = item.querySelector('.model-name');
    if (nameEl && nameEl.textContent.includes(name.split(' ')[0])) {
      const id   = item.dataset.id;
      const prog = document.getElementById(`prog-${id}`);
      const pct  = document.getElementById(`pct-${id}`);
      if (prog) prog.value = percent / 100;
      if (pct)  pct.textContent = `${Math.round(percent)}%`;
    }
  });
};

window.onDownloadComplete = function(modelId) {
  const statusDiv = document.getElementById(`status-${modelId}`);
  if (statusDiv) statusDiv.innerHTML = installedBadge();
  reloadSettings();
};

window.onDownloadError = function(modelId, err) {
  // Restore the download button on failure
  const btn      = document.getElementById(`btn-${modelId}`);
  const progWrap = document.getElementById(`prog-wrap-${modelId}`);
  if (btn)      { btn.hidden = false; }
  if (progWrap) { progWrap.hidden = true; }
  console.error('Download error:', modelId, err);
};

// ---- Hotkey ----
function renderHotkey(trigger, isWayland) {
  const x11Row     = document.getElementById('hotkey-x11');
  const waylandRow = document.getElementById('hotkey-wayland');

  if (isWayland) {
    if (x11Row)     x11Row.hidden     = true;
    if (waylandRow) waylandRow.hidden = false;
    return;
  }

  if (waylandRow) waylandRow.hidden = true;
  if (!x11Row) return;
  x11Row.hidden = false;

  updateHotkeyDisplay(trigger);

  const editBtn = document.getElementById('hotkey-edit-btn');
  if (editBtn) editBtn.addEventListener('click', () => showRecordModal(trigger));
}

function updateHotkeyDisplay(trigger) {
  const display = document.getElementById('hotkey-display');
  if (!display) return;
  display.innerHTML = trigger.split('+')
    .map(k => `<kbd>${k}</kbd>`)
    .join('<span style="color:var(--muted);font-size:11px;padding:0 2px">+</span>');
}

// ---- Record hotkey modal ----
let _recordingHotkey = false;

function showRecordModal(currentTrigger) {
  const modal = document.getElementById('hotkey-modal');
  if (!modal) return;
  modal.classList.add('visible');
  _recordingHotkey = true;

  const keyHandler = async (e) => {
    e.preventDefault();
    if (!_recordingHotkey) return;

    const parts = [];
    if (e.ctrlKey)  parts.push('ctrl');
    if (e.shiftKey) parts.push('shift');
    if (e.altKey)   parts.push('alt');
    if (e.metaKey)  parts.push('super');

    const k = e.key.toLowerCase();
    if (!['control','shift','alt','meta'].includes(k)) {
      parts.push(k === ' ' ? 'space' : k);
      const trigger = parts.join('+');
      _recordingHotkey = false;
      document.removeEventListener('keydown', keyHandler);

      const res = await window.saveHotkey(trigger);
      modal.classList.remove('visible');
      if (!res.startsWith('error')) {
        updateHotkeyDisplay(trigger);
      }
    }
  };
  document.addEventListener('keydown', keyHandler);

  const cancelBtn = document.getElementById('hotkey-modal-cancel');
  if (cancelBtn) {
    cancelBtn.onclick = () => {
      _recordingHotkey = false;
      document.removeEventListener('keydown', keyHandler);
      modal.classList.remove('visible');
    };
  }
}
