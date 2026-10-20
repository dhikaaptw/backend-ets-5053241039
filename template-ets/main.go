package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

var sumberDaya = []SumberDaya{}
var idSumberDaya = 1

func tambahSumberDaya(w http.ResponseWriter, r *http.Request) {
	var baru SumberDaya

	err := json.NewDecoder(r.Body).Decode(&baru)
	if err != nil {
		http.Error(w, "json invalid", 400)
		return
	}

	baru.Nama = strings.TrimSpace(baru.Nama)
	baru.Jenis = strings.TrimSpace(baru.Jenis)

	if baru.Nama == "" || len([]rune(baru.Nama)) > 50 {
		http.Error(w, "name invalid", 400)
		return
	}

	if baru.Jenis != "ruang" && baru.Jenis != "alat" {
		http.Error(w, "jenis must be ruang/alat", 400)
		return
	}

	if baru.Stok < 1 {
		http.Error(w, "stock must be at least 1", 400)
		return
	}

	baru.ID = idSumberDaya
	idSumberDaya++

	sumberDaya = append(sumberDaya, baru)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(baru)
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

func NewHandler() http.Handler {
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
	http.ListenAndServe(":8080", NewHandler())
}
