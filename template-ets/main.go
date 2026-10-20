package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
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
	peminjaman   []Peminjaman
	idPeminjaman int
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

func (a *Aplikasi) tambahPeminjaman(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SumberDayaID *int    `json:"sumber_daya_id"`
		Peminjam     *string `json:"peminjam"`
		Tanggal      *string `json:"tanggal"`
		SlotMulai    *int    `json:"slot_mulai"`
		SlotSelesai  *int    `json:"slot_selesai"`
		Jumlah       *int    `json:"jumlah"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		kirimError(w, 400, "json invalid")
		return
	}

	if input.SumberDayaID == nil || input.Peminjam == nil ||
		input.Tanggal == nil || input.SlotMulai == nil ||
		input.SlotSelesai == nil || input.Jumlah == nil {
		kirimError(w, 400, "missing required fields")
		return
	}

	var baru Peminjaman
	baru.SumberDayaID = *input.SumberDayaID
	baru.Peminjam = *input.Peminjam
	baru.Tanggal = *input.Tanggal
	baru.SlotMulai = *input.SlotMulai
	baru.SlotSelesai = *input.SlotSelesai
	baru.Jumlah = *input.Jumlah

	baru.Peminjam = strings.TrimSpace(baru.Peminjam)
	baru.Tanggal = strings.TrimSpace(baru.Tanggal)

	if baru.Peminjam == "" || len([]rune(baru.Peminjam)) > 50 {
		kirimError(w, 400, "name invalid")
		return
	}

	tanggal, err := time.Parse("2006-01-02", baru.Tanggal)
	if err != nil || tanggal.Format("2006-01-02") != baru.Tanggal {
		kirimError(w, 400, "date invalid")
		return
	}

	if baru.SlotMulai < 1 || baru.SlotMulai > 8 ||
		baru.SlotSelesai < 1 || baru.SlotSelesai > 8 ||
		baru.SlotMulai > baru.SlotSelesai {
		kirimError(w, 400, "slot invalid")
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	ada := false
	stok := 0

	for i := 0; i < len(a.sumberDaya); i++ {
		if a.sumberDaya[i].ID == baru.SumberDayaID {
			ada = true
			stok = a.sumberDaya[i].Stok
			break
		}
	}

	if !ada {
		kirimError(w, 400, "resource not found")
		return
	}

	if baru.Jumlah < 1 || baru.Jumlah > stok {
		kirimError(w, 400, "invalid amount")
		return
	}

	baru.ID = a.idPeminjaman
	a.idPeminjaman++
	baru.Status = "aktif"

	a.peminjaman = append(a.peminjaman, baru)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(baru)
}

func (a *Aplikasi) semuaPeminjaman(w http.ResponseWriter, r *http.Request) {
	filterID := r.URL.Query().Get("sumber_daya_id")
	filterTanggal := r.URL.Query().Get("tanggal")
	filterStatus := r.URL.Query().Get("status")

	id := 0

	if filterID != "" {
		var err error
		id, err = strconv.Atoi(filterID)

		if err != nil {
			kirimError(w, 400, "id invalid")
			return
		}
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	hasil := []Peminjaman{}

	for i := 0; i < len(a.peminjaman); i++ {
		data := a.peminjaman[i]

		if filterID != "" && data.SumberDayaID != id {
			continue
		}

		if filterTanggal != "" && data.Tanggal != filterTanggal {
			continue
		}

		if filterStatus != "" && data.Status != filterStatus {
			continue
		}

		hasil = append(hasil, data)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hasil)
}

func (a *Aplikasi) detailPeminjaman(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil {
		kirimError(w, 400, "id invalid")
		return
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	for i := 0; i < len(a.peminjaman); i++ {
		if a.peminjaman[i].ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(a.peminjaman[i])
			return
		}
	}

	kirimError(w, 404, "loan not found")
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
		peminjaman:   []Peminjaman{},
		idPeminjaman: 1,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /sumber-daya", a.tambahSumberDaya)
	mux.HandleFunc("GET /sumber-daya", a.semuaSumberDaya)
	mux.HandleFunc("GET /sumber-daya/{id}", a.detailSumberDaya)

	mux.HandleFunc("POST /peminjaman", a.tambahPeminjaman)

	mux.HandleFunc("GET /peminjaman", a.semuaPeminjaman)
	mux.HandleFunc("GET /peminjaman/{id}", a.detailPeminjaman)

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
