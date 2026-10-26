// State Management
let currentCategory = 'Semua';
let currentSearch = '';
let currentSort = 'name_asc';
let searchDebounceTimer = null;

// Compare selection from localStorage
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
  updateCompareBar();
}

function toggleCompare(id) {
  let ids = getCompareIds();
  const index = ids.indexOf(id);
  if (index > -1) {
    ids.splice(index, 1);
  } else {
    if (ids.length >= 4) {
      alert('Maksimal 4 protein dapat dibandingkan sekaligus.');
      return;
    }
    ids.push(id);
  }
  saveCompareIds(ids);
  renderProteins(window.latestProteins || []);
}

function updateCompareBar() {
  const ids = getCompareIds();
  const bar = document.getElementById('compareBar');
  const barText = document.getElementById('compareBarText');
  if (!bar) return;

  if (ids.length > 0) {
    bar.classList.remove('hidden');
    barText.textContent = `${ids.length} protein dipilih untuk perbandingan`;
  } else {
    bar.classList.add('hidden');
  }
}

// Fetch Stats for Hero Section
async function fetchStats() {
  try {
    const res = await fetch('/api/stats');
    const json = await res.json();
    if (json.success && json.data) {
      const elTotal = document.getElementById('stat-total');
      const elHewani = document.getElementById('stat-hewani');
      const elNabati = document.getElementById('stat-nabati');
      if (elTotal) elTotal.textContent = json.data.total;
      if (elHewani) elHewani.textContent = json.data.hewani_count;
      if (elNabati) elNabati.textContent = json.data.nabati_count;
    }
  } catch (err) {
    console.error('Failed to load stats:', err);
  }
}

// Fetch Catalog Proteins
async function fetchProteins() {
  const loader = document.getElementById('loader');
  const grid = document.getElementById('proteinGrid');
  const empty = document.getElementById('emptyState');

  if (loader) loader.style.display = 'block';
  if (grid) grid.style.display = 'none';
  if (empty) empty.style.display = 'none';

  const params = new URLSearchParams();
  if (currentSearch) params.append('search', currentSearch);
  if (currentCategory && currentCategory !== 'Semua') params.append('category', currentCategory);
  if (currentSort) params.append('sort', currentSort);

  try {
    const res = await fetch(`/api/proteins?${params.toString()}`);
    const json = await res.json();

    if (loader) loader.style.display = 'none';
    if (json.success && json.data && json.data.length > 0) {
      window.latestProteins = json.data;
      renderProteins(json.data);
      if (grid) grid.style.display = 'grid';
    } else {
      if (empty) empty.style.display = 'block';
    }
  } catch (err) {
    if (loader) loader.style.display = 'none';
    if (empty) empty.style.display = 'block';
    console.error('Failed to load proteins:', err);
  }
}

// Render Protein Cards
function renderProteins(proteins) {
  const grid = document.getElementById('proteinGrid');
  if (!grid) return;
  const compareIds = getCompareIds();

  grid.innerHTML = proteins.map(p => {
    const isSelected = compareIds.includes(p.id);
    const badgeClass = p.source_category.toLowerCase() === 'hewani' ? 'badge-hewani' : 'badge-nabati';

    return `
      <article class="protein-card">
        <div class="card-header">
          <span class="badge ${badgeClass}">${p.source_category}</span>
          <span style="font-size: 0.75rem; color: var(--text-muted); font-weight: 600;">PDB: ${p.pdb_id}</span>
        </div>
        <h2 class="card-title">${escapeHTML(p.name)}</h2>
        <p class="card-subtitle">${escapeHTML(p.common_name || p.food_source)}</p>

        <div class="card-props">
          <div>
            <span class="label">Berat Molekul</span>
            <span class="val">${p.molecular_weight ? p.molecular_weight + ' kDa' : '-'}</span>
          </div>
          <div>
            <span class="label">Daya Cerna</span>
            <span class="val">${p.digestibility_score ? p.digestibility_score + '%' : '-'}</span>
          </div>
          <div>
            <span class="label">Titik Isoelektrik (pI)</span>
            <span class="val">${p.isoelectric_point || '-'}</span>
          </div>
          <div>
            <span class="label">UniProt ID</span>
            <span class="val">${p.uniprot_id || '-'}</span>
          </div>
        </div>

        <p class="card-desc">${escapeHTML(p.biological_function || p.description)}</p>

        <div class="card-actions">
          <a href="/detail?id=${p.id}" class="btn btn-primary" style="flex: 1;">
            🔍 Lihat Detail & 3D
          </a>
          <button class="btn btn-secondary btn-compare-toggle ${isSelected ? 'selected' : ''}" 
                  onclick="toggleCompare(${p.id})" 
                  title="${isSelected ? 'Hapus dari komparasi' : 'Pilih untuk komparasi'}">
            ${isSelected ? '✓ Terpilih' : '⚖️ Bandingkan'}
          </button>
        </div>
      </article>
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

// Initial Bootstrapping
document.addEventListener('DOMContentLoaded', () => {
  fetchStats();
  fetchProteins();
  updateCompareBar();

  // Search input with debounce
  const searchInput = document.getElementById('searchInput');
  if (searchInput) {
    searchInput.addEventListener('input', (e) => {
      clearTimeout(searchDebounceTimer);
      searchDebounceTimer = setTimeout(() => {
        currentSearch = e.target.value.trim();
        fetchProteins();
      }, 300);
    });
  }

  // Category filter tabs
  const filterBtns = document.querySelectorAll('.filter-btn[data-category]');
  filterBtns.forEach(btn => {
    btn.addEventListener('click', () => {
      filterBtns.forEach(b => b.classList.remove('active'));
      btn.classList.add('active');
      currentCategory = btn.getAttribute('data-category');
      fetchProteins();
    });
  });

  // Sort dropdown
  const sortSelect = document.getElementById('sortSelect');
  if (sortSelect) {
    sortSelect.addEventListener('change', (e) => {
      currentSort = e.target.value;
      fetchProteins();
    });
  }

  // Clear compare bar
  const clearBtn = document.getElementById('clearCompareBtn');
  if (clearBtn) {
    clearBtn.addEventListener('click', () => {
      saveCompareIds([]);
      renderProteins(window.latestProteins || []);
    });
  }
});
