let currentProtein = null;
let glviewer = null;
let isSpinning = true;
let currentStyle = 'cartoon';
let aminoChart = null;

// Helper to get compare list
function getCompareIds() {
  try {
    const raw = localStorage.getItem('compare_proteins');
    return raw ? JSON.parse(raw) : [];
  } catch (e) {
    return [];
  }
}

function saveCompareIds(ids) {
  localStorage.setItem('compare_proteins', JSON.stringify(ids));
}

async function loadProteinDetail() {
  const urlParams = new URLSearchParams(window.location.search);
  const id = urlParams.get('id');

  const loader = document.getElementById('loader');
  const content = document.getElementById('detailContent');
  const errorEl = document.getElementById('detailError');

  if (!id) {
    alert('ID protein tidak ditemukan pada URL');
    window.location.href = '/';
    return;
  }

  // Reset state tiap kali dipanggil
  loader.style.display = 'block';
  content.style.display = 'none';
  if (errorEl) errorEl.style.display = 'none';

  try {
    const res = await fetch(`/api/proteins/${id}`);

    // Tangani respons HTTP error (404, 500, dst) sebelum parse JSON
    if (!res.ok) {
      let errorMsg = `Gagal memuat detail protein (HTTP ${res.status}).`;
      try {
        const errJson = await res.json();
        if (errJson.error) errorMsg += ' ' + errJson.error;
      } catch (e) {
        errorMsg += ' ' + res.statusText;
      }
      throw new Error(errorMsg);
    }

    const json = await res.json();

    loader.style.display = 'none';
    if (json.success && json.data) {
      currentProtein = json.data;
      renderDetail(json.data);
      content.style.display = 'block';
      init3DViewer(json.data.pdb_id);
      renderAminoAcidChart(json.data.parsed_amino_acids || [], 'bar');
      renderAminoAcidTable(json.data.parsed_amino_acids || []);
    } else {
      alert('Protein tidak ditemukan: ' + (json.error || ''));
      window.location.href = '/';
    }
  } catch (err) {
    loader.style.display = 'none';
    console.error('Failed to load detail:', err);

    // Tampilkan error inline dengan tombol Coba Lagi, atau alert fallback
    if (errorEl) {
      errorEl.querySelector('.error-message').textContent = err.message;
      errorEl.style.display = 'block';
    } else {
      alert('Gagal memuat detail protein: ' + err.message);
    }
  }
}

// Expose untuk tombol Coba Lagi di HTML
window.retryLoadProteinDetail = loadProteinDetail;

function renderDetail(p) {
  document.title = `${p.name} — Protein Pangan`;
  document.getElementById('proteinName').textContent = p.name;
  document.getElementById('commonName').textContent = p.common_name || p.food_source;

  const badgeCategory = document.getElementById('badgeCategory');
  badgeCategory.textContent = p.source_category;
  badgeCategory.className = 'badge ' + (p.source_category.toLowerCase() === 'hewani' ? 'badge-hewani' : 'badge-nabati');

  document.getElementById('uniprotBadge').textContent = `UniProt: ${p.uniprot_id || '-'}`;
  document.getElementById('pdbBadge').textContent = `PDB: ${p.pdb_id || '-'}`;

  document.getElementById('valMw').textContent = p.molecular_weight ? `${p.molecular_weight} kDa` : '-';
  document.getElementById('valPi').textContent = p.isoelectric_point ? `${p.isoelectric_point}` : '-';
  document.getElementById('valScore').textContent = p.digestibility_score ? `${p.digestibility_score}%` : '-';
  document.getElementById('valSource').textContent = p.food_source || '-';
  document.getElementById('valCategory').textContent = p.source_category || '-';
  document.getElementById('valPdb').textContent = p.pdb_id || '-';

  document.getElementById('proteinDesc').textContent = p.description || '';
  document.getElementById('bioFunc').textContent = p.biological_function || '-';

  // External Links
  const pdbUrl = `https://www.rcsb.org/structure/${p.pdb_id}`;
  const uniprotUrl = p.reference_url || `https://www.uniprot.org/uniprotkb/${p.uniprot_id}/entry`;

  document.getElementById('refPdbLink').href = pdbUrl;
  document.getElementById('refPdbLink').textContent = `RCSB PDB ${p.pdb_id} (${p.name})`;

  document.getElementById('refUniprotLink').href = uniprotUrl;
  document.getElementById('refUniprotLink').textContent = `UniProt ${p.uniprot_id} (${p.name})`;
}

// 3Dmol.js Viewer Initialization
function init3DViewer(pdbId) {
  const container = document.getElementById('viewer3d');
  const loading = document.getElementById('viewerLoading');
  if (!container || !window.$3Dmol) {
    console.warn('3Dmol library not loaded yet');
    return;
  }

  container.innerHTML = '';
  if (loading) container.appendChild(loading);
  const config = { backgroundColor: '#090d16' };
  glviewer = $3Dmol.createViewer(container, config);

  if (!pdbId) {
    if (loading) loading.textContent = 'PDB ID tidak tersedia untuk protein ini.';
    return;
  }

  // Download PDB directly from RCSB PDB using 3Dmol helper
  $3Dmol.download(`pdb:${pdbId}`, glviewer, {}, function() {
    if (loading) loading.style.display = 'none';
    applyViewerStyle('cartoon');
    glviewer.zoomTo();
    glviewer.render();
    glviewer.spin(true);
    isSpinning = true;
    const spinBtn = document.getElementById('btnSpin');
    if (spinBtn) spinBtn.classList.add('active');
  });
}

function applyViewerStyle(style) {
  if (!glviewer) return;
  currentStyle = style;
  glviewer.removeAllSurfaces();

  switch (style) {
    case 'stick':
      glviewer.setStyle({}, { stick: { colorscheme: 'amino' } });
      break;
    case 'sphere':
      glviewer.setStyle({}, { sphere: { scale: 0.85, colorscheme: 'spectrum' } });
      break;
    case 'surface':
      glviewer.setStyle({}, { cartoon: { color: 'spectrum' } });
      glviewer.addSurface($3Dmol.SurfaceType.VDW, { opacity: 0.85, colorscheme: 'whiteCarbon' });
      break;
    case 'cartoon':
    default:
      glviewer.setStyle({}, { cartoon: { color: 'spectrum' } });
      break;
  }
  glviewer.render();
}

// Amino Acid Chart Rendering (Chart.js)
function renderAminoAcidChart(aminoAcids, chartType = 'bar') {
  const canvas = document.getElementById('aminoAcidChart');
  if (!canvas || !window.Chart) return;

  if (aminoChart) {
    aminoChart.destroy();
  }

  const labels = aminoAcids.map(a => `${a.name} (${a.code})`);
  const data = aminoAcids.map(a => a.amount);
  const backgroundColors = aminoAcids.map(a =>
    a.is_essential ? 'rgba(13, 148, 136, 0.85)' : 'rgba(59, 130, 246, 0.7)'
  );
  const borderColors = aminoAcids.map(a =>
    a.is_essential ? '#0f766e' : '#2563eb'
  );

  const ctx = canvas.getContext('2d');

  if (chartType === 'radar') {
    aminoChart = new Chart(ctx, {
      type: 'radar',
      data: {
        labels: aminoAcids.map(a => a.code),
        datasets: [{
          label: 'Kandungan Asam Amino (g/100g)',
          data: data,
          backgroundColor: 'rgba(13, 148, 136, 0.3)',
          borderColor: '#0d9488',
          pointBackgroundColor: '#0f766e',
          pointBorderColor: '#fff',
          pointHoverBackgroundColor: '#fff',
          pointHoverBorderColor: '#0d9488',
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          r: {
            angleLines: { color: '#e2e8f0' },
            grid: { color: '#e2e8f0' },
            pointLabels: { font: { size: 11, weight: 'bold' } },
            ticks: { beginAtZero: true }
          }
        }
      }
    });
  } else {
    // Default Bar Chart
    aminoChart = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: labels,
        datasets: [{
          label: 'Kandungan Asam Amino (g / 100g protein)',
          data: data,
          backgroundColor: backgroundColors,
          borderColor: borderColors,
          borderWidth: 1,
          borderRadius: 4
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: { display: false },
          tooltip: {
            callbacks: {
              afterLabel: function(ctx) {
                const item = aminoAcids[ctx.dataIndex];
                return item.is_essential ? 'Kategori: Asam Amino Esensial' : 'Kategori: Asam Amino Non-Esensial';
              }
            }
          }
        },
        scales: {
          y: {
            beginAtZero: true,
            title: { display: true, text: 'g / 100g' }
          },
          x: {
            ticks: {
              maxRotation: 45,
              minRotation: 30
            }
          }
        }
      }
    });
  }
}

// Amino Acid Table
function renderAminoAcidTable(aminoAcids) {
  const tbody = document.getElementById('aminoTableBody');
  if (!tbody) return;

  const maxAmount = Math.max(...aminoAcids.map(a => a.amount), 20);

  tbody.innerHTML = aminoAcids.map(a => {
    const pct = ((a.amount / maxAmount) * 100).toFixed(0);
    const badgeType = a.is_essential
      ? '<span class="badge" style="background:#ccfbf1; color:#0f766e;">Esensial</span>'
      : '<span class="badge" style="background:#f1f5f9; color:#475569;">Non-Esensial</span>';

    return `
      <tr>
        <td><strong>${escapeHTML(a.name)}</strong></td>
        <td><code>${escapeHTML(a.code)}</code></td>
        <td>${badgeType}</td>
        <td><strong>${a.amount.toFixed(1)} g</strong></td>
        <td style="width: 30%;">
          <div style="background:#e2e8f0; border-radius:999px; height:8px; overflow:hidden;">
            <div style="background:${a.is_essential ? '#0d9488' : '#3b82f6'}; width:${pct}%; height:100%;"></div>
          </div>
        </td>
      </tr>
    `;
  }).join('');
}

function escapeHTML(str) {
  if (!str) return '';
  return str.replace(/[&<>'"]/g, tag => ({
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    "'": '&#39;',
    '"': '&quot;'
  }[tag] || tag));
}

// Event Listeners
document.addEventListener('DOMContentLoaded', () => {
  loadProteinDetail();

  // Viewer Style Select
  const styleSelect = document.getElementById('styleSelect');
  if (styleSelect) {
    styleSelect.addEventListener('change', () => {
      applyViewerStyle(styleSelect.value);
    });
  }

  // Spin Toggle Button
  const spinBtn = document.getElementById('btnSpin');
  if (spinBtn) {
    spinBtn.addEventListener('click', () => {
      if (!glviewer) return;
      isSpinning = !isSpinning;
      glviewer.spin(isSpinning);
      glviewer.render();
      spinBtn.textContent = `🔄 Putar Otomatis: ${isSpinning ? 'ON' : 'OFF'}`;
      spinBtn.classList.toggle('active', isSpinning);
    });
  }

  // Reset Camera Button
  const resetBtn = document.getElementById('btnResetView');
  if (resetBtn) {
    resetBtn.addEventListener('click', () => {
      if (!glviewer) return;
      glviewer.zoomTo();
      glviewer.render();
    });
  }

  // Chart Type Toggles
  const barBtn = document.getElementById('btnChartBar');
  const radarBtn = document.getElementById('btnChartRadar');

  if (barBtn && radarBtn) {
    barBtn.addEventListener('click', () => {
      barBtn.classList.add('active');
      radarBtn.classList.remove('active');
      if (currentProtein) renderAminoAcidChart(currentProtein.parsed_amino_acids || [], 'bar');
    });

    radarBtn.addEventListener('click', () => {
      radarBtn.classList.add('active');
      barBtn.classList.remove('active');
      if (currentProtein) renderAminoAcidChart(currentProtein.parsed_amino_acids || [], 'radar');
    });
  }

});
