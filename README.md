# 🧬 Protein Pangan — Platform Analisis & Visualisasi 3D Protein Pangan

Aplikasi web modern dan interaktif untuk katalog, eksplorasi karakteristik fisikokimia, visualisasi struktur molekuler 3D, serta perbandingan profil asam amino protein pangan sehari-hari (hewani dan nabati). Dirancang khusus untuk mendukung riset biologi, teknologi pangan, dan analisis gizi.

---

## 🌟 Fitur Utama

1. **Katalog & Pencarian Cepat (*Search & Filter*):**
   * Pencarian kata kunci secara instan (nama protein, nama umum, sumber pangan, kode PDB, atau UniProt ID).
   * Filter kategori sumber: **Semua**, **🥩 Hewani**, dan **🌱 Nabati**.
   * Pengurutan (*sorting*) berdasarkan nama, berat molekul (kDa), atau skor daya cerna.
   * Metrik ringkasan jumlah protein di katalog.

2. **Halaman Detail & Viewer Struktur 3D Interaktif:**
   * Visualisasi interaktif struktur kristalografi 3D molekul protein langsung di browser menggunakan library `3Dmol.js`.
   * Kontrol viewer: Rotasi otomatis (*auto-spin*), ganti gaya visual (*Cartoon/Ribbon*, *Stick*, *Sphere*, *Surface*), dan *reset camera*.
   * Data ilmiah terperinci: Berat molekul (kDa), Titik Isoelektrik (pI), Daya Cerna (skor DIAAS/PDCAAS), dan fungsi biologis.
   * Tautan langsung ke basis data internasional: [RCSB PDB](https://www.rcsb.org/) dan [UniProt Knowledgebase](https://www.uniprot.org/).

3. **Analisis Profil Asam Amino:**
   * Visualisasi diagram batang (*Bar Chart*) dan diagram radar (*Radar Chart*) menggunakan `Chart.js`.
   * Klasifikasi asam amino esensial vs non-esensial secara visual.
   * Tabel rincian gram kandungan per 100g protein murni.

4. **Komparasi Antar-Protein (*Side-by-Side Comparison*):**
   * Pemilihan 2 atau lebih protein dari katalog untuk dibandingkan langsung.
   * Tabel perbandingan komparatif parameter fisikokimia.
   * Diagram grafik multi-dataset yang menumpuk profil asam amino antar-protein untuk analisis nutrisi komparatif.

5. **Sistem Autentikasi & Otorisasi Pengguna (JWT):**
   * Pendaftaran akun pengguna (*Register*) dengan enkripsi kata sandi menggunakan `bcrypt`.
   * Masuk sistem (*Login*) via Username atau Email dengan penerbitan token JSON Web Token (JWT).
   * **Pengalihan Otomatis (*Auto-Redirect*):** Setelah login/register berhasil, sistem secara otomatis mengalihkan browser langsung ke endpoint aplikasi katalog protein (`/`).
   * Halaman login/register mandiri (`/login` & `/register`) dengan deteksi sesi aktif otomatis.
   * Middleware autentikasi Bearer Token untuk pengamanan rute-rute API terproteksi (`/api/auth/me`).
   * Modal dialog autentikasi terintegrasi di navbar frontend beserta penyimpanan sesi di `localStorage`.

---

## 🛠️ Tech Stack

* **Backend:** [Go](https://go.dev/) (v1.22+) dengan framework [Gin Web Framework](https://github.com/gin-gonic/gin)
* **ORM:** [GORM](https://gorm.io/) (PostgreSQL Driver)
* **Database:** [PostgreSQL](https://www.postgresql.org/) (Serverless Cloud Neon DB)
* **Frontend:** Modern Vanilla HTML5, CSS3 kustom responsif, dan Vanilla JavaScript (ES6+ Fetch API)
* **Visualisasi & Grafik:** 
  * Struktur 3D: [3Dmol.js](https://3dmol.csb.pitt.edu/) (via CDN)
  * Diagram Grafik: [Chart.js](https://www.chartjs.org/) (via CDN)
* **Data Seeding:** Embed data ilmiah tervalidasi (`//go:embed` dari `seed_proteins.json`)

---

## 📁 Struktur Direktori

```text
go_vibe/
├── .env.example              # Template variabel lingkungan
├── .env                      # Konfigurasi aktif (diabaikan oleh git)
├── .gitignore                # Aturan file yang diabaikan git
├── go.mod                    # Modul Go & dependensi
├── go.sum                    # Checksum dependensi Go
├── main.go                   # Titik masuk utama aplikasi (entry point)
├── README.md                 # Dokumentasi proyek
├── config/
│   └── database.go           # Inisialisasi koneksi PostgreSQL Neon dengan GORM
├── models/
│   ├── protein.go            # Model entitas Protein & AminoAcid
│   └── user.go               # Model entitas Pengguna & hash kata sandi bcrypt
├── utils/
│   └── jwt.go                # Utilitas pembuatan & verifikasi JWT
├── middleware/
│   └── auth.go               # Middleware proteksi Bearer JWT token
├── data/
│   ├── seed_proteins.json    # Data awal literatur protein ilmiah
│   └── seeder.go             # Logika auto-seed database saat pertama jalan
├── handlers/
│   ├── protein_handler.go    # Logika handler REST API Protein
│   ├── protein_handler_test.go # Pengujian unit & integrasi Protein
│   ├── auth_handler.go       # Logika handler Register, Login, & Profile
│   └── auth_handler_test.go  # Pengujian unit & integrasi Autentikasi JWT
├── routes/
│   ├── routes.go             # Definisi routing Gin & serving halaman statis
│   └── routes_test.go        # Pengujian endpoint routing
├── web/
│   ├── index.html            # Halaman katalog utama aplikasi protein
│   ├── login.html            # Halaman portal login & registrasi mandiri
│   ├── detail.html           # Halaman detail & viewer 3D
│   ├── compare.html          # Halaman komparasi antar-protein
│   ├── css/
│   │   └── style.css         # Styling antarmuka modern & modal auth responsif
│   └── js/
│       ├── main.js           # Logika interaksi katalog & state komparasi
│       ├── detail.js         # Logika 3Dmol.js viewer & Chart.js detail
│       ├── compare.js        # Logika tabel komparasi & grafik multi-dataset
│       └── auth.js           # Manajer sesi JWT & interaksi modal auth
└── docs/
    └── prompts/
        └── 00-project-spec.md # Dokumen spesifikasi proyek
```

---

## 🚀 Panduan Menjalankan Aplikasi

### 1. Prasyarat
* Telah terpasang **Go 1.22** atau yang lebih baru.
* Koneksi internet (untuk Neon DB dan CDN 3Dmol/Chart.js).

### 2. Konfigurasi Variabel Lingkungan
Salin file `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```
Pastikan file `.env` memiliki variabel `DATABASE_URL` dan `PORT`:
```dotenv
PORT=8080
DATABASE_URL=postgresql://neondb_owner:npg_8HquofB3izDy@ep-raspy-leaf-ao5zeti1-pooler.c-2.ap-southeast-1.aws.neon.tech/neondb?sslmode=require
GIN_MODE=debug
```

### 3. Menjalankan Server
Jalankan perintah berikut di direktori proyek:
```bash
go run main.go
```
*Aplikasi akan melakukan koneksi ke Neon DB, migrasi skema tabel secara otomatis, menyuntikkan data awal (seeding), dan menjalankan server di `http://localhost:8080`.*

Buka peramban (browser) dan akses:
* **Katalog:** [http://localhost:8080](http://localhost:8080)
* **Perbandingan:** [http://localhost:8080/compare](http://localhost:8080/compare)

---

## 🧪 Menjalankan Pengujian (Testing)

Untuk menjalankan seluruh unit test dan test integrasi handler:
```bash
go test -v ./...
```

---

## 📡 Dokumentasi Endpoint REST API

| Method | Endpoint | Deskripsi | Query Parameters |
|---|---|---|---|
| `GET` | `/api/health` | Pemeriksaan status server | - |
| `GET` | `/api/stats` | Statistik total data dan kategori | - |
| `GET` | `/api/proteins` | Mengambil daftar katalog protein | `search`, `category`, `sort` |
| `GET` | `/api/proteins/:id` | Mengambil rincian protein berdasarkan ID | - |
| `GET` | `/api/proteins/compare` | Komparasi beberapa protein sekaligus | `ids` (cth: `?ids=1,2,5`) |

---
