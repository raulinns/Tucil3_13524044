# Ice Sliding Puzzle Solver

Program Go untuk menyelesaikan Ice Sliding Puzzle dengan UCS, GBFS, A*, BFS, dan DFS. GUI web React tersedia untuk upload testcase, memilih algoritma, melihat animasi solusi, dan melihat trace graf pencarian.

Semua perintah di bawah dijalankan dari root repository:

```txt
Tugas-Kecil-3/
```

## Menjalankan CLI

```bash
go run ./src
```

CLI akan meminta path file testcase `.txt`, pilihan algoritma, dan pilihan heuristic jika algoritma memakai heuristic.

Contoh path testcase saat prompt CLI muncul:

```txt
test/sample.txt
```

## Menjalankan GUI dengan Docker

Cara ini menjalankan backend Go dan frontend React dari satu container.

```bash
docker compose down
docker compose up --build
```

Buka:

```txt
http://localhost:8080
```

Di GUI:

1. Upload file testcase `.txt`.
2. Pilih algoritma `UCS`, `GBFS`, `A*`, `BFS`, atau `DFS`.
3. Pilih heuristic untuk `GBFS` atau `A*`.
4. Klik `Solve`.
5. Gunakan playback untuk melihat perpindahan pada board.
6. Lihat bagian `Search Graph` untuk node yang diekspansi dan edge yang ditemukan solver.

## Menjalankan GUI tanpa Docker

Terminal 1, jalankan backend:

```bash
go run ./src/web/backend
```

Terminal 2, jalankan frontend:

```bash
cd src/web/frontend
npm install
npm run dev
```

Buka:

```txt
http://localhost:5173
```

Vite akan meneruskan request `/solve` ke backend di `http://localhost:8080`.

Untuk membuat build frontend dari root tanpa masuk permanen ke folder frontend:

```bash
npm --prefix src/web/frontend install
npm --prefix src/web/frontend run build
```

Untuk menjalankan test Go dari root:

```bash
go test ./src/...
```

## Format Input

File harus berekstensi `.txt`. Format umum:

```txt
N M
<N baris grid>
<N x M nilai cost>
```

Karakter grid:

- `X`: tembok
- `*`: ice / tile biasa
- `L`: lava
- `Z`: start
- `O`: goal
- `0` sampai `9`: checkpoint berurutan

## Endpoint Backend

`POST /solve`

Request `multipart/form-data`:

- `file`: file testcase `.txt`
- `algorithm`: `UCS`, `GBFS`, `ASTAR`, `BFS`, atau `DFS`
- `heuristic`: `H1`, `H2`, atau `H3`

Response utama:

- `found`: apakah solusi ditemukan
- `moves`: urutan gerakan, misalnya `RDLR`
- `totalCost`: total cost solusi
- `iterations`: jumlah node yang diekspansi
- `executionMs`: waktu eksekusi solver
- `steps`: visualisasi board per langkah solusi
- `trace`: node dan edge untuk visualisasi graf pencarian

## Catatan Implementasi

- Solver CLI memakai API `UCS`, `GBFS`, `AStar`, `BFS`, dan `DFS`.
- GUI memakai wrapper `UCSWithTrace`, `GBFSWithTrace`, `AStarWithTrace`, `BFSWithTrace`, dan `DFSWithTrace`.
- Data trace diletakkan di package `searchtrace`, sehingga struktur `solver.Result` tetap fokus pada hasil pencarian.
- `BFS` optimal terhadap jumlah slide, sedangkan `DFS` tidak menjamin optimalitas.
