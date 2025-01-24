### **Backend - `README.md`**

# Crypto Wallet Tracker Backend

Backend untuk Crypto Wallet Tracker menyediakan API untuk manajemen wallet, autentikasi menggunakan JWT, dan komunikasi real-time dengan WebSocket.

## Fitur

- API RESTful untuk manajemen wallet.
- Server WebSocket untuk chat real-time.
- Autentikasi JWT untuk akses aman.

## Teknologi yang Digunakan

- Gofiber (Golang)
- Gorilla WebSocket
- JWT untuk autentikasi
- PostgreSQL (opsional jika menggunakan database)

## Cara Install dan Menjalankan Aplikasi

### 1. Clone Repository
```bash
git clone https://github.com/your-repo/backend.git
cd backend
```

### 2. Install Dependencies
Pastikan Go sudah terinstall, lalu jalankan:
```bash
go mod tidy
```

### 3. Konfigurasi Environment
Buat file `.env` di root project dengan isi berikut:
```plaintext
PORT=8080
JWT_SECRET=your-secret-key
```

### 4. Menjalankan Aplikasi
Jalankan perintah berikut:
```bash
go run main.go
```
Backend akan berjalan di `http://localhost:8080`.

### 5. Build untuk Produksi
Untuk membuat binary produksi, jalankan:
```bash
go build -o backend
```

## Endpoint API

### Wallets
- `GET /wallets/username/:username`  
  Mendapatkan wallet berdasarkan username.
- `POST /wallets`  
  Membuat wallet baru.
- `PUT /wallets/:id`  
  Mengupdate wallet berdasarkan ID.
- `DELETE /wallets/:id`  
  Menghapus wallet berdasarkan ID.
### Announcements
- GET /announcements
Mendapatkan daftar semua pengumuman.
- POST /announcements
Membuat pengumuman baru.

### Chat
- `GET /api/chat-history`  
  Mendapatkan riwayat chat.

### Autentikasi
- `POST /login`  
  Menghasilkan JWT berdasarkan username.

## Penggunaan WebSocket

Hubungkan ke server WebSocket:
```
ws://localhost:8080/ws
```

Kirim dan terima pesan secara real-time:
```json
{ "type": "message", "username": "user1", "message": "Hello, world!" }
```
