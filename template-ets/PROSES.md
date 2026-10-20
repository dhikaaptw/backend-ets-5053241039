# PROSES

## 1. Persona dan masalahnya
Persona E. Masalahnya adalah terkadang laboratorium digunakan untuk praktikum atau perbaikan, sehingga tidak bisa dipinjam pada jadwal tertentu.

## 2. Rancangan struktur data
- SumberDaya: menyimpan id, nama, jenis, dan stok dari ruang atau alat.

- Peminjaman: menyimpan id peminjaman, id sumber daya, nama peminjam, tanggal, slot mulai, slot selesai, jumlah, dan status.

- Aplikasi: digunakan untuk menampung data sumber daya dan peminjaman, termasuk nomor id berikutnya dan mutex.

Awalnya saya menggunakan variabel global untuk menyimpan data, tetapi kemudian dipindahkan ke dalam struct Aplikasi supaya setiap pemanggilan NewHandler() memiliki penyimpanan sendiri.

## 3. Satu keputusan desain
(Satu keputusan yang kamu ambil, pilihan lain yang kamu pertimbangkan, dan alasanmu memilih yang ini. Contoh: kenapa cek bentrok per slot, bukan per rentang.)

## 4. Bug atau kendala (minimal dua)
### Bug 1
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

### Bug 2
- Gejala:
- Cara menemukan:
- Penyebab dan solusi:

## 5. Change request
(Apa yang diminta, bagian kode mana yang kamu ubah, dan satu keputusan desain yang kamu ambil.)
