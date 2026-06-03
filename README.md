# Jakarta BBM Tracker

Aplikasi pelacak harga bahan bakar minyak (BBM) modern berbasis Vue.js 3 untuk Jakarta, Indonesia. Pantau harga BBM dari merek-merek besar Indonesia (Pertamina, Shell, Vivo, BP) di berbagai wilayah Jakarta.

## Fitur Utama

✨ **Keunggulan Aplikasi:**
- 📍 Jelajahi stasiun pengisian bahan bakar (SPBU) di 5 wilayah Jakarta (Utara, Selatan, Timur, Barat, Pusat)
- ⛽ Lacak harga untuk BBM RON 90, RON 92, RON 95, RON 98, dan Diesel CN 51/53
- 🏢 Filter berdasarkan merek (Pertamina, Shell, Vivo, BP)
- 📊 Sortir berdasarkan harga, merek, atau wilayah
- 📱 Desain responsif mobile-first yang sempurna
- 🗺️ Peta interaktif Leaflet dengan penanda lokasi stasiun
- 🎨 UI dengan gradien modern dan interaksi halus
- 🔍 Informasi detail stasiun dalam popup peta

## Stack Teknologi

- **Vue.js 3** - Framework JavaScript progresif
- **TypeScript** - JavaScript dengan type safety
- **Vite** - Build tool yang super cepat
- **Tailwind CSS** - Utility-first CSS framework
- **Leaflet.js** - Library peta interaktif
- **PostCSS & Autoprefixer** - CSS processing otomatis

## Struktur Proyek

```
src/
├── components/
│   ├── Header.vue        # Header aplikasi dengan branding
│   ├── FilterBar.vue     # Filter wilayah, merek, dan jenis BBM
│   ├── Map.vue           # Peta interaktif Leaflet dengan penanda SPBU
│   └── StationList.vue   # Daftar stasiun di sidebar dengan harga
├── data/
│   └── stations.ts       # Data dummy stasiun BBM dan region
├── types/
│   └── index.ts          # Type definitions TypeScript
├── App.vue               # Komponen aplikasi utama (layout 2 kolom)
├── main.ts               # Entry point aplikasi
└── style.css             # Global styles
```

## Memulai

### Prasyarat
- Node.js 16+
- npm atau yarn

### Instalasi

```bash
# Install dependencies
npm install

# Jalankan development server
npm run dev

# Build untuk production
npm run build

# Preview production build
npm run preview
```

## Development

Development server berjalan di `http://localhost:5173/`

### Script yang Tersedia

- `npm run dev` - Jalankan development server
- `npm run build` - Build untuk production
- `npm run preview` - Preview production build

## Data

**Catatan**: Data SPBU di `src/data/stations.ts` berasal dari `SPBU_DKI_Jakarta.csv`. Harga BBM tetap memakai nilai default per merek dan dapat ditimpa cache backend saat tersedia.

### Data Mencakup:
- 115 SPBU di 5 wilayah Jakarta
- 4 merek BBM dengan warna branding yang autentik
- 6 jenis BBM (RON 90/92/95/98, Diesel CN 51/53)
- Alamat dan link Google Maps dari CSV
- Informasi alamat dan wilayah lengkap

## Styling & Warna

Aplikasi menggunakan skema warna modern:
- **Gradien Utama**: Purple (#667eea ke #764ba2)
- **Warna Brand**:
  - Pertamina: Navy Blue (#003d7a)
  - Shell: Merah (#e81c1c)
  - Vivo: Biru (#3b82f6)
  - BP: Hijau (#00a651)

## Komponen-Komponen

### Header
Menampilkan judul aplikasi dan branding dengan latar belakang gradien.

### FilterBar
Menyediakan kontrol filter dan sorting yang kompak:
- Filter berdasarkan wilayah
- Filter berdasarkan merek
- Filter berdasarkan jenis BBM
- Tombol reset filter

### Map
Peta interaktif berbasis Leaflet yang menampilkan:
- Penanda SPBU berwarna berdasarkan merek
- Popup dengan informasi stasiun dan harga
- Animasi flyTo saat memilih stasiun dari daftar
- Support untuk geolokasi dan zoom interaktif

### StationList
Daftar stasiun di sidebar yang menampilkan:
- Badge merek dengan warna brand
- Nama stasiun dan alamat
- Daftar harga BBM
- Tombol "Lihat di Peta" untuk menggeser peta ke lokasi stasiun

## Backend API (Go)

Aplikasi juga menyediakan backend Go yang melakukan scraping harga BBM real-time dari situs resmi brand dan menyediakan REST API.

### Menjalankan Backend

```bash
cd backend
go run ./cmd/server
```

Server berjalan di `http://localhost:8080`

### Endpoints Utama

- `GET /health` - Health check API
- `GET /api/v1/prices` - Ambil harga BBM dari cache
- `POST /api/v1/scrape` - Trigger scraping ulang harga BBM

### Konfigurasi

Lihat `backend/README.md` untuk detail lengkap environment variables dan konfigurasi backend.

## Browser Support

- Chrome/Edge (versi terbaru)
- Firefox (versi terbaru)
- Safari (versi terbaru)
- Mobile browsers

## Lisensi

MIT

## Penulis

Dibuat sebagai aplikasi demonstrasi pelacak BBM yang terinspirasi dari FuelWatch WA.
