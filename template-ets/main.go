package main

import (
	"net/http"
)

func tambahSumberDaya(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST sumber daya"))
}

func semuaSumberDaya(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET semua sumber daya"))
}

func detailSumberDaya(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET detail sumber daya"))
}

func tambahPeminjaman(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST peminjaman"))
}

func semuaPeminjaman(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET semua peminjaman"))
}

func detailPeminjaman(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET detail peminjaman"))
}

func ubahPeminjaman(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PUT peminjaman"))
}

func hapusPeminjaman(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DELETE peminjaman"))
}

func tambahBlokir(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("POST blokir"))

}
func semuaBlokir(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET semua blokir"))
}

func hapusBlokir(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("DELETE blokir"))
}

func newHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /sumber-daya", tambahSumberDaya)
	mux.HandleFunc("GET /sumber-daya", semuaSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", detailSumberDaya)

	mux.HandleFunc("POST /peminjaman", tambahPeminjaman)

	mux.HandleFunc("GET /peminjaman", semuaPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", detailPeminjaman)

	mux.HandleFunc("PUT /peminjaman/{id}", ubahPeminjaman)
	mux.HandleFunc("DELETE /peminjaman/{id}", hapusPeminjaman)

	mux.HandleFunc("POST /blokir", tambahBlokir)
	mux.HandleFunc("GET /blokir", semuaBlokir)
	mux.HandleFunc("DELETE /blokir/{id}", hapusBlokir)

	return mux
}

func main() {
	http.ListenAndServe(":8080", newHandler())
}
