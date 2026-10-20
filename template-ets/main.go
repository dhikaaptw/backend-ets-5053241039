package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type SumberDaya struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Jenis string `json:"jenis"`
	Stok  int    `json:"stok"`
}

type Peminjaman struct {
	ID           int    `json:"id"`
	SumberDayaID int    `json:"sumber_daya_id"`
	Peminjam     string `json:"peminjam"`
	Tanggal      string `json:"tanggal"`
	SlotMulai    int    `json:"slot_mulai"`
	SlotSelesai  int    `json:"slot_selesai"`
	Jumlah       int    `json:"jumlah"`
	Status       string `json:"status"`
}

type Aplikasi struct {
	sumberDaya   []SumberDaya
	idSumberDaya int
	mu           sync.Mutex
}

func kirimError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func (a *Aplikasi) tambahSumberDaya(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nama  *string `json:"nama"`
		Jenis *string `json:"jenis"`
		Stok  *int    `json:"stok"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		kirimError(w, 400, "json invalid")
		return
	}

	if input.Nama == nil || input.Jenis == nil || input.Stok == nil {
		kirimError(w, 400, "missing required fields")
		return
	}

	var baru SumberDaya
	baru.Nama = *input.Nama
	baru.Jenis = *input.Jenis
	baru.Stok = *input.Stok

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

	a.mu.Lock()
	defer a.mu.Unlock()

	baru.ID = a.idSumberDaya
	a.idSumberDaya++

	a.sumberDaya = append(a.sumberDaya, baru)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(baru)
}

func (a *Aplikasi) semuaSumberDaya(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a.sumberDaya)
}

func (a *Aplikasi) detailSumberDaya(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		kirimError(w, 400, "id invalid")
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	for i := 0; i < len(a.sumberDaya); i++ {
		if a.sumberDaya[i].ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(a.sumberDaya[i])
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
	a := &Aplikasi{
		sumberDaya:   []SumberDaya{},
		idSumberDaya: 1,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /sumber-daya", a.tambahSumberDaya)
	mux.HandleFunc("GET /sumber-daya", a.semuaSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", a.detailSumberDaya)

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
