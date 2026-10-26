# Project Specification — Vibe Coding

> **Tugas 1:** Isi spesifikasi aplikasi sebelum mulai melakukan vibe coding.
>
> Tentukan kebutuhan aplikasi dan pilih teknologi yang akan digunakan. Jangan mulai membuat kode sebelum spesifikasi selesai.

## 1. Informasi Aplikasi

**Nama aplikasi:**  
`Protein Pangan`

**Deskripsi aplikasi:**  
`Aplikasi katalog dan eksplorasi data protein pangan sehari-hari (hewani dan nabati) yang dilengkapi visualisasi struktur 3D interaktif, analisis profil asam amino, serta perbandingan antar-protein untuk mendukung riset biologi dan nutrisi pangan.`

**Target pengguna:**  
`Peneliti Biologi, Akademisi, Mahasiswa Sains, dan Praktisi Nutrisi/Gizi`

**Masalah yang ingin diselesaikan:**  
`Membantu dan memudahkan peneliti serta praktisi gizi dalam menemukan informasi komprehensif, terstruktur, dan visual mengenai karakteristik protein dalam bahan pangan sehari-hari secara cepat dan terstandar.`

---

## 2. Kebutuhan Aplikasi

**Fitur utama aplikasi:**

1. **Katalog Protein Pangan:** Daftar protein pangan hewani dan nabati dengan fungsionalitas pencarian kata kunci (*search*) dan filter kategori (*kategori sumber: hewani/nabati*).
2. **Halaman Detail Protein:** Informasi mendalam mencakup sumber pangan, fungsi biologis, berat molekul, kelarutan, nilai biologis, dan tautan referensi ilmiah (UniProt ID & PDB ID).
3. **Visualisasi Struktur 3D Interaktif:** Viewer struktur molekul protein 3D langsung pada web browser menggunakan library `3Dmol.js` (mendukung rotasi, zoom, dan perubahan mode tampilan seperti pita/ribbon, stick, surface).
4. **Perbandingan Protein & Profil Asam Amino:** Fitur komparasi berdampingan (*side-by-side*) antara dua atau lebih protein serta visualisasi diagram grafik komposisi asam amino esensial & non-esensial menggunakan `Chart.js`.

**Halaman aplikasi:**

- **Katalog / Beranda (*Home & Catalog Page*):** Menampilkan daftar kartu (*cards*) protein pangan dengan pencarian, filter kategori hewani/nabati, dan ringkasan singkat.
- **Halaman Detail Protein (*Detail Page*):** Menampilkan data komprehensif protein yang dipilih, viewer struktur molekul 3D, serta tabel profil asam amino.
- **Halaman Perbandingan (*Comparison Page*):** Memungkinkan pemilihan 2 protein untuk dibandingkan karakteristik, berat molekul, dan grafik radar/bar komposisi asam aminonya.

---

## 3. Pilihan Teknologi

Pilih satu opsi untuk setiap komponen yang diperlukan.

### A. Backend Framework

- [ ] 1. Go Standard Library (`net/http`)
- [x] 2. Gin
- [ ] 3. Lainnya sebutkan

**Pilihan saya:** 2 (Gin)  
**Ketentuan:** Menggunakan framework Gin untuk membangun REST API yang ringan, cepat, dan terstruktur rapi.

### B. ORM (Object-Relational Mapping)

- [ ] 1. Tidak menggunakan ORM (SQL manual)
- [x] 2. GORM
- [ ] 3. Lainnya sebutkan

**Pilihan saya:** 2 (GORM)  
**Ketentuan:** Menggunakan GORM sebagai ORM untuk pemodelan entitas, relasi, auto-migrasi skema, dan seeding data protein ke database.

### C. Database

- [x] 1. PostgreSQL
- [ ] 2. MySQL
- [ ] 3. SQLite

**Pilihan saya:** 1 (PostgreSQL)  
**Ketentuan:** Menggunakan database PostgreSQL (di-host pada Neon DB) yang andal dalam menyimpan data relasional dan skema JSON jika diperlukan.

### D. Frontend

- [x] 1. HTML, CSS, JavaScript murni
- [ ] 2. React
- [ ] 3. Vue

**Pilihan saya:** 1 (HTML, CSS, JavaScript murni)  
**Ketentuan:** Menggunakan Vanilla HTML5, modern responsive CSS, dan Vanilla JavaScript (ES6+ fetch API). Visualisasi 3D menggunakan `3Dmol.js` dan grafik visualisasi menggunakan `Chart.js` yang dimuat via CDN.

### E. Sistem Autentikasi

- [ ] 1. Session
- [x] 2. JWT
- [ ] 3. Tidak diperlukan (Aplikasi Publik / Terbuka)

**Pilihan saya:** 2 (JWT)  
**Ketentuan:** Menggunakan JSON Web Token (JWT) dengan rahasia kunci `JWT_SECRET=35zmwq80E3/CRrsGzzUHG8B8tlP26P11sWocAuCzRIk=` dan enkripsi kata sandi `bcrypt` untuk otentikasi pendaftaran (*register*) dan masuk (*login*) akun pengguna.

### F. Penyimpanan Sesi atau Token

- [ ] 1. Cookie
- [x] 2. localStorage
- [ ] 3. Tidak diperlukan

**Pilihan saya:** 2 (localStorage)  
**Ketentuan:** Token JWT disimpan pada `localStorage` di sisi klien web browser dan dikirimkan melalui header `Authorization: Bearer <token>` pada setiap permintaan endpoint terproteksi (`/api/auth/me`).

### G. API

- [x] 1. REST API
- [ ] 2. Tidak membutuhkan API terpisah

**Pilihan saya:** 1 (REST API)  
**Ketentuan:** Backend menyediakan RESTful JSON API (misalnya: `GET /api/proteins`, `GET /api/proteins/:id`, `GET /api/proteins/compare`). Frontend mengonsumsi API menggunakan fungsi bawaan `fetch()`.

### H. Version Control

- [x] 1. Git
- [ ] 2. Tidak menggunakan Git

**Pilihan saya:** 1 (Git)

---

## 4. Skema Data & Model Entitas (Database Design)

Struktur tabel utama `proteins`:
- `id`: uint (Primary Key, Auto Increment)
- `name`: string (Nama protein, cth: *Ovalbumin*, *Kasein*, *Glutenin*, *Aktin*)
- `source_category`: string (*Hewani* atau *Nabati*)
- `food_source`: string (Bahan makanan asal, cth: *Putih Telur*, *Susu Sapi*, *Gandum*, *Daging Sapi*)
- `molecular_weight`: float64 (Berat molekul dalam satuan kDa)
- `biological_function`: text (Fungsi biologis dan perannya dalam pangan/tubuh)
- `pdb_id`: string (Kode rujukan RCSB PDB untuk visualisasi 3D, cth: `1OVA`, `1JEP`)
- `uniprot_id`: string (Kode rujukan UniProt, cth: `P01012`)
- `reference_url`: string (Tautan rujukan literatur ilmiah/database publik)
- `amino_acid_profile`: json/text (Komposisi persentase asam amino penting untuk diagram Chart.js)
- `created_at` & `updated_at`: timestamps

---

## 5. Detail Akses & Keamanan

**Registrasi & Login Akun:** Tidak diperlukan (Aplikasi Katalog Terbuka).  
**Role Pengguna:** Pengunjung Publik / Peneliti Biologi.

**Aturan Akses:**
1. Semua pengguna memiliki akses penuh untuk menjelajah katalog, memfilter data, melihat detail, dan membandingkan profil protein.
2. Aplikasi pada tahap ini bersifat **baca-saja (read-only)** untuk pengguna akhir, sehingga tidak ada form publik yang mengubah database (mencegah manipulasi data riset).
3. Data awal disuntikkan secara otomatis melalui mekanisme seeder resmi di backend saat server pertama kali dijalankan.

---

## 6. Konfigurasi Environment dan Database

**Ketentuan wajib:**

1. Semua konfigurasi sensitif harus menggunakan environment variable.
2. Gunakan file `.env` untuk konfigurasi lokal dan koneksi ke database.
3. Sediakan file `.env.example` sebagai referensi struktur konfigurasi tanpa mencantumkan kredensial asli.
4. File `.env` wajib dimasukkan ke dalam `.gitignore` dan tidak boleh di-commit ke Git.
5. Connection string PostgreSQL dibaca dari environment variable `DATABASE_URL`.
6. Jangan menulis URL database, username, atau password secara hardcode di dalam kode sumber Go.
7. Aplikasi harus menangani error koneksi database dengan pesan log yang informatif dan aman.

**Contoh `.env.example`:**

```dotenv
PORT=8080
DATABASE_URL=postgresql://username:password@localhost:5432/protein_db?sslmode=disable
GIN_MODE=debug
```

---

## 7. Batasan Pengembangan

- [x] Struktur folder sederhana, bersih, dan berstandar Go idiomatik.
- [x] Kode mudah dipahami, terdokumentasi, dan modular.
- [x] Tidak menambahkan library berlebih di luar kebutuhan spesifikasi.
- [x] Pemisahan yang jelas antara frontend (UI/assets) dan backend (Gin handlers, models, seeders).
- [x] Library visualisasi frontend (`3Dmol.js` dan `Chart.js`) dimuat melalui CDN tanpa dependensi bundler/npm yang rumit.
- [x] Data awal protein bersumber dari literatur ilmiah valid (UniProt / PDB) dan disediakan melalui seeder JSON/GORM.
- [x] Menggunakan environment variable untuk konfigurasi dinamis.
- [x] Menyediakan pengujian dasar (unit test / integration test handler).

---

## 8. Pemeriksaan Sebelum Vibe Coding

- [x] Nama dan tujuan aplikasi sudah jelas (Protein Pangan).
- [x] Fitur utama sudah ditentukan (Katalog, Detail, 3Dmol, Chart.js, Perbandingan).
- [x] Pilihan teknologi konsisten (Go Gin, GORM, PostgreSQL Neon, Vanilla JS + CDN).
- [x] Kebutuhan akses ditentukan (Aplikasi publik read-only tanpa autentikasi).
- [x] Format environment variable sudah benar dan valid.
- [x] Batasan proyek dan skema data sudah terdefinisi.

**Konfigurasi `.env` aktif yang disiapkan:**
```dotenv
PORT=8080
DATABASE_URL=postgresql://neondb_owner:npg_8HquofB3izDy@ep-raspy-leaf-ao5zeti1-pooler.c-2.ap-southeast-1.aws.neon.tech/neondb?sslmode=require
GIN_MODE=debug
```

**Status spesifikasi:** Telah diperbaiki, konsisten, lengkap, dan siap untuk implementasi kode.

---
