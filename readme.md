# 🔐 Turing Machine Simulator: Iterative vs Recursive Analysis

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green)
![Status](https://img.shields.io/badge/Status-Completed-success)

> **Tugas Besar Analisis Kompleksitas Algoritma** > Studi Kasus: Dekripsi Caesar Cipher menggunakan Simulator Mesin Turing.

## 📖 Deskripsi Proyek

Aplikasi ini adalah **Simulator Mesin Turing** berbasis web yang dibangun menggunakan **Go (Golang)**. Aplikasi ini dirancang untuk mendemonstrasikan dan membandingkan performa dua pendekatan algoritma fundamental—**Iteratif** dan **Rekursif**—dalam memecahkan masalah dekripsi sandi sederhana (Caesar Cipher).

Fokus utama proyek ini bukan hanya pada fungsionalitas dekripsi, melainkan pada **Benchmarking (Pengujian Beban)**. Aplikasi ini mampu melakukan *stress test* dari 100 hingga **100 Juta data input** untuk memvisualisasikan perbedaan Kompleksitas Ruang (*Space Complexity*) antara penggunaan memori Heap (Iteratif) dan Stack (Rekursif).

## ✨ Fitur Utama

1.  **Dual Mode Operation:**
    * ✍️ **Mode Manual:** Input teks sandi Anda sendiri (mendukung huruf, spasi, angka, dan simbol) dan lihat hasil dekripsinya secara instan.
    * 📊 **Mode Benchmark:** Menjalankan pengujian otomatis dengan input *dummy* (100 - 100 Juta karakter) untuk mengukur waktu eksekusi.

2.  **High-Performance Backend:**
    * Ditulis dengan **Go** untuk manajemen memori yang efisien.
    * Menggunakan struktur data `[]byte` (*mutable*) untuk performa maksimal.
    * Dilengkapi *Safety Wrapper* (`defer-recover`) untuk menangani *Stack Overflow* pada algoritma rekursif.

3.  **Real-time Visualization:**
    * Grafik interaktif menggunakan **Chart.js**.
    * Tabel log data presisi (hingga milidetik).
    * Antarmuka modern dengan tema *Dark Mode* dan *Glassmorphism*.

## 📷 Screenshots

### 1. Dashboard & Mode Manual
*(Masukkan screenshot input manual di sini, contoh file: `manual_mode.png`)*
![Manual Mode](https://via.placeholder.com/800x400?text=Screenshot+Mode+Manual)

### 2. Benchmark & Grafik Analisis
*(Masukkan screenshot grafik benchmark di sini, contoh file: `benchmark_mode.png`)*
![Benchmark Mode](https://via.placeholder.com/800x400?text=Screenshot+Grafik+Benchmark)

## 🛠️ Teknologi yang Digunakan

* **Backend:** Go (Golang) Standard Library (`net/http`, `html/template`, `encoding/json`).
* **Frontend:** HTML5, CSS3, JavaScript (Vanilla).
* **Library:** [Chart.js](https://www.chartjs.org/) (CDN).

## 🚀 Cara Menjalankan

Pastikan Anda sudah menginstall [Go](https://go.dev/dl/) di komputer Anda.

1.  **Clone Repository ini:**
    ```bash
    git clone [https://github.com/USERNAME_ANDA/NAMA_REPO.git](https://github.com/USERNAME_ANDA/NAMA_REPO.git)
    cd NAMA_REPO
    ```

2.  **Jalankan Aplikasi:**
    ```bash
    go run main.go
    ```
    *(Atau `go run tubes.go` sesuai nama file Anda)*

3.  **Buka Browser:**
    Akses alamat berikut di browser Anda:
    ```
    http://localhost:8080
    ```

## 📊 Analisis Algoritma

Berdasarkan hasil pengujian, berikut adalah ringkasan perbandingan kedua algoritma:

| Parameter | Algoritma Iteratif | Algoritma Rekursif |
| :--- | :--- | :--- |
| **Metode** | Looping (`for`) | Self-calling function |
| **Time Complexity** | $O(n)$ - Linear | $O(n)$ - Linear |
| **Space Complexity** | **$O(1)$** - Efisien (Heap) | **$O(n)$** - Boros (Stack) |
| **Stabilitas** | ✅ Stabil hingga >100 Juta data | ❌ Crash/Stack Overflow di ~15 Juta data |

**Kesimpulan:**
Algoritma Iteratif jauh lebih unggul untuk pemrosesan data linear dalam jumlah besar karena tidak membebani memori Stack sistem operasi.

## 👤 Author

**Ihsan Dwika Putra**
* **NIM:** 103012400119
* **Institusi:** Telkom University
* **Role:** Full Stack Developer & System Analyst

---
*Dibuat untuk memenuhi Tugas Besar Mata Kuliah Analisis Kompleksitas Algoritma (2025).*