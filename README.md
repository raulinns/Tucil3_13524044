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
