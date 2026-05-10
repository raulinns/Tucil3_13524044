# Ice Sliding Puzzle Solver

Program Go untuk menyelesaikan Ice Sliding Puzzle dengan UCS, GBFS, A*, BFS, dan DFS. GUI web React tersedia untuk upload testcase, memilih algoritma dan heuristic, melihat playback solusi, melihat trace graf pencarian, serta menyimpan hasil solusi ke file `.txt`.

Semua perintah dijalankan dari root repository:

```txt
Tugas-Kecil-3/
```

## Requirement

- Go 1.24 atau lebih baru
- Node.js dan npm untuk menjalankan frontend tanpa Docker
- Docker dan Docker Compose untuk menjalankan GUI dalam container

## Kompilasi

```bash
go build -o bin/ice-sliding-solver ./src
```

Jalankan executable hasil kompilasi:

```bash
./bin/ice-sliding-solver
```

## Menjalankan CLI

```bash
go run ./src
```

CLI akan meminta path file testcase `.txt`, pilihan algoritma, pilihan heuristic untuk GBFS/A*, lalu menawarkan penyimpanan hasil solusi ke `.txt`.

Contoh path testcase saat prompt CLI muncul:

```txt
test/testcase_ucs.txt
```

## Menjalankan GUI dengan Docker

```bash
docker compose down
docker compose up --build
```

Buka:

```txt
http://localhost:8080
```

## Menjalankan GUI tanpa Docker

Terminal 1:

```bash
go run ./src/web/backend
```

Terminal 2:

```bash
cd src/web/frontend
npm install
npm run dev
```

Frontend Vite berjalan pada URL yang ditampilkan terminal, umumnya `http://localhost:5173`.

## Cara Menggunakan GUI

1. Upload file testcase `.txt`.
2. Pilih algoritma `UCS`, `GBFS`, `A*`, `BFS`, atau `DFS`.
3. Pilih heuristic `H1`, `H2`, atau `H3` untuk `GBFS` atau `A*`.
4. Klik `Solve`.
5. Gunakan playback untuk maju, mundur, atau mengatur kecepatan visualisasi.
6. Lihat `Search Graph` untuk node dan edge yang diekspansi.
7. Klik `Save .txt` untuk menyimpan ringkasan solusi.

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

## Testcase

Testcase laporan berada di folder `test/`:

- `testcase_ucs.txt`: testcase utama untuk UCS.
- `testcase_gbfs.txt`: testcase papan lebih kompleks untuk GBFS.
- `testcase_astar.txt`: testcase utama untuk A*.
- `testcase_dfs.txt`: testcase papan non-persegi untuk DFS.
- `testcase_bfs.txt`: testcase papan non-persegi untuk BFS.
- `testcase_sanity.txt`: sanity check sederhana.
- `edge_no_solution.txt`: input valid tanpa solusi.
- `edge_invalid_*.txt`: testcase error handling parser.

## Author

Narendra Dharma Wistara M. - 13524044
