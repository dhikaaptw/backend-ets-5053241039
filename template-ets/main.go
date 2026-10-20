package main

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func kirimError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func tambahSumberDaya(w http.ResponseWriter, r *http.Request) {
	var baru SumberDaya

	err := json.NewDecoder(r.Body).Decode(&baru)
	if err != nil {
		kirimError(w, 400, "json invalid")
		return
	}

	baru.Nama = strings.TrimSpace(baru.Nama)
	baru.Jenis = strings.TrimSpace(baru.Jenis)

	if baru.Nama == "" || len([]rune(baru.Nama)) > 50 {
		kirimError(w, 400, "name invalid")
		return
	}

	if baru.Jenis != "ruang" && baru.Jenis != "alat" {
		kirimError(w, 400, "jenis must be ruang/alat")
		return
	}

	if baru.Stok < 1 {
		kirimError(w, 400, "stock must be at least 1")
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
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		kirimError(w, 400, "invalid id")
		return
	}

	for i := 0; i < len(sumberDaya); i++ {
		if sumberDaya[i].ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(sumberDaya[i])
			return
		}
	}

	kirimError(w, 404, "resource not found")
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
