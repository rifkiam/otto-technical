# Pair Programming – BE Golang CRUD Sederhana

Tujuan: menguji pemahaman dasar Golang (HTTP, JSON, struktur data, error handling) melalui implementasi CRUD sederhana dengan penyimpanan in-memory.

## Pertanyaan (Soal)

Bangun REST API sederhana untuk resource `Item` (id, name, done) dengan operasi:

- `GET /health` → cek status.
- `POST /items` → buat item baru (name wajib, done default false).
- `GET /items` → daftar semua item.
- `GET /items/{id}` → ambil item by id.
- `PUT /items/{id}` → update parsial (name/done).
- `DELETE /items/{id}` → hapus item.

Catatan bisnis sederhana:

- `name` harus non-empty (trimmed) pada create/update.
- `done` hanya boleh diubah via update, bukan saat create.

## Mulai dari nol atau sudah ada?

- Tidak mulai dari nol. Disediakan skeleton proyek lengkap beserta test.
- Kandidat fokus melengkapi logic penyimpanan in-memory dan handler HTTP agar semua test lulus.

## Apa yang sudah disediakan

- Struktur proyek dengan modul `be_test` (lihat `go.mod`).
- Model dasar di `internal/model` untuk struktur `Item`.
- Kontrak store dan error di `internal/store/store.go`.
- Skeleton implementasi in-memory di `internal/store/memory.go` untuk dilengkapi.
- Skeleton HTTP server dan handler di `internal/httpapi/server.go`.
- Entrypoint server di `cmd/server/main.go` (opsional untuk manual testing).
- Test end-to-end di `tests/server_test.go` meliputi health, CRUD, dan validasi.
- Skrip demo opsional di `scripts/demo.sh` untuk menjalankan alur `curl` cepat.

## Yang harus dikerjakan kandidat

- Lengkapi `internal/store/memory.go` (Create, Get, List, Update, Delete) memakai map + mutex, ID unik, timestamps UTC.
- Lengkapi handler di `internal/httpapi/server.go` untuk endpoint `/items` dan `/items/{id}`.
- Terapkan validasi: `name` harus non-empty; mapping error ke status yang benar (400, 404, 422, 201/200/204).
- Pastikan semua test di `tests/server_test.go` PASS: Health, CRUD, Validation.

## Instruksi Pengerjaan

1. Pastikan Go 1.21+ terpasang.
2. Buka folder proyek ini: `be-test`.
3. Jalankan perintah: `go test ./... -v` untuk melihat requirement (sebagian test GAGAL karena skeleton belum diimplementasikan).
4. Implementasikan:
   - `internal/store/memory.go`: gunakan map untuk menyimpan item, tambahkan ID unik, dan timestamps (`time.Now().UTC()`). Lengkapi method `Create`, `Get`, `List`, `Update`, `Delete`.
   - `internal/httpapi/server.go`: lengkapi handler `/items`, `/items/{id}` agar sesuai ekspektasi test.
5. Jalankan lagi `go test ./... -v` hingga semua test PASS.

## Catatan Penting

- Tes menggunakan `httptest` sehingga tidak perlu server terpisah untuk verifikasi.
- Namun kini tersedia entrypoint: jalankan `go run ./cmd/server` untuk mencoba endpoint secara manual di `http://localhost:8080`.
- Tidak perlu menggunakan database nyata; cukup penyimpanan in-memory (map) sesuai skeleton.

## Cara Menjalankan Server (Opsional)

- `cd be-test`
- `go run ./cmd/server`
- Contoh curl:
  - Health: `curl -s http://localhost:8080/health`
  - Create: `curl -s -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{"name":"Apple"}'`
  - List: `curl -s http://localhost:8080/items`
  - Get by ID: `curl -s http://localhost:8080/items/<id>`
  - Update: `curl -s -X PUT http://localhost:8080/items/<id> -H "Content-Type: application/json" -d '{"name":"Banana","done":true}'`
  - Delete: `curl -i -s -X DELETE http://localhost:8080/items/<id>`

### Skrip Demo (opsional)

- Jalankan: `bash scripts/demo.sh`
- Skrip ini akan mengeksekusi health → create → list → get → update → delete → get(404) secara berurutan, dan menampilkan respons.

## Contoh Output `go test` yang LULUS

Jika implementasi sudah benar, output akan seperti ini:

```
$ go test ./... -v
?       be_test/internal/httpapi     [no test files]
?       be_test/internal/model       [no test files]
?       be_test/internal/store       [no test files]
=== RUN   TestHealthOK
--- PASS: TestHealthOK (0.00s)
=== RUN   TestCRUDSimple
--- PASS: TestCRUDSimple (0.00s)
=== RUN   TestValidation
--- PASS: TestValidation (0.00s)
PASS
ok      be_test/tests        0.6s
```

## Ekspektasi Teknis

- Status code:
  - 200 untuk GET/PUT sukses, 201 untuk create, 204 untuk delete.
  - 404 jika id tidak ditemukan, 422 jika input invalid (name kosong), 400 untuk JSON rusak.
- Response JSON selalu memiliki `Content-Type: application/json`.
- Struktur data sederhana, kode bersih dan mudah dibaca.

## Panduan HTTP Status & Konstanta Go

- Sukses Create: `201 Created` → `http.StatusCreated`
- Sukses Get/List/Update: `200 OK` → `http.StatusOK`
- Sukses Delete: `204 No Content` → `http.StatusNoContent`
- JSON rusak/format salah: `400 Bad Request` → `http.StatusBadRequest`
- Data tidak ditemukan: `404 Not Found` → `http.StatusNotFound`
- Validasi gagal (mis. `name` kosong): `422 Unprocessable Entity` → `http.StatusUnprocessableEntity`
- Method tidak diizinkan: `405 Method Not Allowed` → `http.StatusMethodNotAllowed`

Contoh penggunaan di handler:

- Create: `writeJSON(w, http.StatusCreated, item)`
- Update: `writeJSON(w, http.StatusOK, item)`
- Delete: `w.WriteHeader(http.StatusNoContent)`
- Not found: `writeJSON(w, http.StatusNotFound, map[string]string{"error":"not_found"})`
- Invalid: `writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error":"invalid"})`

## Penilaian

- Semua test lulus.
- Kualitas kode: sederhana, rapi, penanganan error jelas.
- Pemahaman dasar Golang: slice/map/struct, JSON encoding/decoding, net/http.

Selamat mengerjakan!
