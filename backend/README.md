# Jakarta BBM Tracker Backend

Backend Go ini melakukan scrape harga BBM saat server dinyalakan, menyimpan hasil sementara ke `data/fuel_prices.json`, lalu menyediakan cache lewat REST API.

## Run

```bash
cd backend
go run .
```

Server default berjalan di `:8080`.

## Env

- `PORT`: port HTTP, default `8080`.
- `PERTAMINA_PRICE_URL`: override URL resmi Pertamina.
- `SHELL_PRICE_URL`: override URL resmi Shell.
- `BP_PRICE_URL`: override URL resmi BP.
- `VIVO_PRICE_URL`: override URL resmi Vivo.
- `SCRAPE_TIMEOUT_SECONDS`: timeout request scrape, default `20`.

## Endpoints

- `GET /health`
- `GET /api/v1/prices`
- `POST /api/v1/scrape`

Catatan: beberapa situs resmi dapat mengubah struktur HTML atau membatasi akses otomatis. Jika scrape gagal, response tetap menyertakan status error per sumber di `errors`, dan cache JSON tetap ditulis dengan hasil yang berhasil.
