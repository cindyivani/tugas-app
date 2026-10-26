let comparedProteins = [];
let compareChart = null;
let currentChartMode = 'bar';

const COLOR_PALETTE = [
  { bg: 'rgba(13, 148, 136, 0.65)', border: '#0d9488', solid: '#0d9488' },
  { bg: 'rgba(37, 99, 235, 0.65)', border: '#2563eb', solid: '#2563eb' },
  { bg: 'rgba(217, 119, 6, 0.65)', border: '#d97706', solid: '#d97706' },
  { bg: 'rgba(225, 29, 72, 0.65)', border: '#e11d48', solid: '#e11d48' },
];

function getCompareIds() {
  const urlParams = new URLSearchParams(window.location.search);
  const paramIds = urlParams.get('ids');
  if (paramIds) {
    return paramIds.split(',').map(s => parseInt(s.trim(), 10)).filter(n => !isNaN(n));
  }
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

async function loadComparison() {
  const ids = getCompareIds();
  const loader = document.getElementById('loader');
  const content = document.getElementById('compareContent');
  const empty = document.getElementById('emptyCompareState');

  if (ids.length === 0) {
    loader.style.display = 'none';
    content.style.display = 'none';
    empty.style.display = 'block';
    renderSelectedPills([]);
    return;
  }

  loader.style.display = 'block';
  content.style.display = 'none';
  empty.style.display = 'none';

  try {
    const res = await fetch(`/api/proteins/compare?ids=${ids.join(',')}`);
    const json = await res.json();

    loader.style.display = 'none';
    if (json.success && json.data && json.data.length > 0) {
      comparedProteins = json.data;
      renderSelectedPills(comparedProteins);
      renderComparisonTable(comparedProteins);
      renderCompareChart(comparedProteins, currentChartMode);
      content.style.display = 'block';
    } else {
      empty.style.display = 'block';
    }
  } catch (err) {
    loader.style.display = 'none';
    empty.style.display = 'block';
    console.error('Failed to load comparison:', err);
  }
}

function removeCompareId(id) {
  let ids = getCompareIds();
  ids = ids.filter(i => i !== id);
  saveCompareIds(ids);
  loadComparison();
}

function renderSelectedPills(proteins) {
  const container = document.getElementById('selectedPills');
  if (!container) return;

  if (proteins.length === 0) {
    container.innerHTML = '<span style="color:var(--text-muted); font-size:0.85rem;">Belum ada protein yang dipilih.</span>';
    return;
  }

  container.innerHTML = proteins.map((p, idx) => {
    const col = COLOR_PALETTE[idx % COLOR_PALETTE.length];
    return `
      <span class="badge" style="background:${col.bg}; color:#ffffff; font-size:0.85rem; padding:0.4rem 0.8rem; display:inline-flex; align-items:center; gap:0.5rem;">
        ${escapeHTML(p.name)}
        <button onclick="removeCompareId(${p.id})" style="background:none; border:none; color:#ffffff; font-weight:bold; cursor:pointer; font-size:0.9rem;">✕</button>
      </span>
    `;
  }).join('');
}

function renderComparisonTable(proteins) {
  const tbody = document.getElementById('comparisonTableBody');
  if (!tbody) return;

  const headerRow = `
    <tr>
      <th class="prop-col">Karakteristik</th>
      ${proteins.map((p, idx) => {
        const col = COLOR_PALETTE[idx % COLOR_PALETTE.length];
        return `
          <th class="protein-header" style="border-top: 4px solid ${col.solid};">
            <h3 style="font-size:1.1rem; margin-bottom:0.25rem;">${escapeHTML(p.name)}</h3>
            <span style="font-size:0.8rem; color:var(--text-muted);">${escapeHTML(p.food_source)}</span>
          </th>
        `;
      }).join('')}
    </tr>
  `;

  const categoryRow = `
    <tr>
      <th class="prop-col">Kategori Sumber</th>
      ${proteins.map(p => `
        <td style="text-align:center;">
          <span class="badge ${p.source_category.toLowerCase() === 'hewani' ? 'badge-hewani' : 'badge-nabati'}">
            ${p.source_category}
          </span>
        </td>
      `).join('')}
    </tr>
  `;

  const mwRow = `
    <tr>
      <th class="prop-col">Berat Molekul (kDa)</th>
      ${proteins.map(p => `
        <td style="text-align:center; font-weight:700; font-size:1.05rem;">
          ${p.molecular_weight ? p.molecular_weight + ' kDa' : '-'}
        </td>
      `).join('')}
    </tr>
  `;

  const piRow = `
    <tr>
      <th class="prop-col">Titik Isoelektrik (pI)</th>
      ${proteins.map(p => `
        <td style="text-align:center; font-weight:600;">
          ${p.isoelectric_point || '-'}
        </td>
      `).join('')}
    </tr>
  `;

  const scoreRow = `
    <tr>
      <th class="prop-col">Skor Daya Cerna (DIAAS)</th>
      ${proteins.map(p => `
        <td style="text-align:center; font-weight:700; color:var(--primary);">
          ${p.digestibility_score ? p.digestibility_score + '%' : '-'}
        </td>
      `).join('')}
    </tr>
  `;

  const pdbRow = `
    <tr>
      <th class="prop-col">PDB ID / Struktur 3D</th>
      ${proteins.map(p => `
        <td style="text-align:center;">
          <a href="https://www.rcsb.org/structure/${p.pdb_id}" target="_blank" class="btn btn-secondary btn-sm" style="font-family:monospace;">
            ${p.pdb_id} ↗
          </a>
        </td>
      `).join('')}
    </tr>
  `;

  const uniprotRow = `
    <tr>
      <th class="prop-col">UniProt Accession</th>
      ${proteins.map(p => `
        <td style="text-align:center;">
          <a href="${p.reference_url || 'https://www.uniprot.org/uniprotkb/' + p.uniprot_id + '/entry'}" target="_blank" style="color:var(--accent); font-weight:600; text-decoration:none;">
            ${p.uniprot_id} ↗
          </a>
        </td>
      `).join('')}
    </tr>
  `;

  const functionRow = `
    <tr>
      <th class="prop-col">Fungsi Biologis & Pangan</th>
      ${proteins.map(p => `
        <td style="font-size:0.875rem; line-height:1.6; color:#334155;">
          ${escapeHTML(p.biological_function)}
        </td>
      `).join('')}
    </tr>
  `;

  const actionRow = `
    <tr>
      <th class="prop-col">Aksi</th>
      ${proteins.map(p => `
        <td style="text-align:center;">
          <a href="/detail?id=${p.id}" class="btn btn-primary btn-sm">
            Buka Detail & 3D Viewer →
          </a>
        </td>
      `).join('')}
    </tr>
  `;

  tbody.innerHTML = headerRow + categoryRow + mwRow + piRow + scoreRow + pdbRow + uniprotRow + functionRow + actionRow;
}

// Chart.js Multi-Dataset Comparison Chart
function renderCompareChart(proteins, mode = 'bar') {
  const canvas = document.getElementById('compareChart');
  if (!canvas || !window.Chart || proteins.length === 0) return;

  if (compareChart) {
    compareChart.destroy();
  }

  // Get common amino acid labels from the first protein
  const firstAminoList = proteins[0].parsed_amino_acids || [];
  const labels = firstAminoList.map(a => `${a.name} (${a.code})`);
  const codeLabels = firstAminoList.map(a => a.code);

  const datasets = proteins.map((p, idx) => {
    const col = COLOR_PALETTE[idx % COLOR_PALETTE.length];
    const aminoMap = {};
    (p.parsed_amino_acids || []).forEach(a => {
      aminoMap[a.code] = a.amount;
    });

    const data = firstAminoList.map(a => aminoMap[a.code] || 0);

    return {
      label: p.name,
      data: data,
      backgroundColor: mode === 'radar' ? col.bg : col.bg,
      borderColor: col.border,
      borderWidth: 2,
      borderRadius: mode === 'bar' ? 3 : 0,
      pointBackgroundColor: col.solid,
      pointBorderColor: '#fff',
      fill: mode === 'radar'
    };
  });

  const ctx = canvas.getContext('2d');

  if (mode === 'radar') {
    compareChart = new Chart(ctx, {
      type: 'radar',
      data: {
        labels: codeLabels,
        datasets: datasets
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        elements: {
          line: { tension: 0.1 }
        },
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
    compareChart = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: labels,
        datasets: datasets
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          tooltip: {
            mode: 'index',
            intersect: false
          }
        },
        scales: {
          y: {
            beginAtZero: true,
            title: { display: true, text: 'g / 100g protein' }
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
  loadComparison();

  // Mode toggles
  const barBtn = document.getElementById('cmpBarBtn');
  const radarBtn = document.getElementById('cmpRadarBtn');

  if (barBtn && radarBtn) {
    barBtn.addEventListener('click', () => {
      barBtn.classList.add('active');
      radarBtn.classList.remove('active');
      currentChartMode = 'bar';
      if (comparedProteins.length > 0) renderCompareChart(comparedProteins, 'bar');
    });

    radarBtn.addEventListener('click', () => {
      radarBtn.classList.add('active');
      barBtn.classList.remove('active');
      currentChartMode = 'radar';
      if (comparedProteins.length > 0) renderCompareChart(comparedProteins, 'radar');
    });
  }

  // Clear all compare
  const clearBtn = document.getElementById('clearAllCompareBtn');
  if (clearBtn) {
    clearBtn.addEventListener('click', () => {
      saveCompareIds([]);
      loadComparison();
    });
  }
});
