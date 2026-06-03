# Jakarta BBM Tracker - Backend API

Backend Go untuk Jakarta BBM Tracker. Melakukan scraping harga BBM dari situs resmi brand (Pertamina, Shell, Vivo, BP) saat server dinyalakan, menyimpan hasil cache ke `data/fuel_prices.json`, dan menyediakan REST API untuk frontend.

## Menjalankan Backend

```bash
cd backend
go run ./cmd/server
```

Server berjalan di `http://localhost:8080` (default).

## Konfigurasi Environment

- `PORT`: Port HTTP server, default `8080`
- `PERTAMINA_PRICE_URL`: Override URL scrape Pertamina
- `SHELL_PRICE_URL`: Override URL scrape Shell
- `BP_PRICE_URL`: Override URL scrape BP
- `VIVO_PRICE_URL`: Override URL scrape Vivo
- `SCRAPE_TIMEOUT_SECONDS`: Timeout untuk scraping request, default `20`

## Endpoints API

### Health Check
```
GET /health
```
Response: `{"status": "ok"}`

### Ambil Harga BBM (Cache)
```
GET /api/v1/prices
```
Response: JSON array harga BBM dari cache `data/fuel_prices.json`

### Trigger Scrape Ulang
```
POST /api/v1/scrape
```
Melakukan scraping harga BBM dari semua sumber dan update cache.

## Catatan Penting

- Beberapa situs resmi brand dapat mengubah struktur HTML atau membatasi akses otomatis
- Jika scrape gagal, API tetap mengembalikan status error per sumber di field `errors`
- Cache JSON (`data/fuel_prices.json`) tetap diupdate dengan hasil yang berhasil di-scrape
- Gunakan endpoint `/api/v1/scrape` untuk refresh harga terbaru secara manual
