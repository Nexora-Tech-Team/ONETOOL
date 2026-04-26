# OneTool - Project Context for Claude

## Overview
OneTool adalah aplikasi bisnis all-in-one yang mencakup manajemen klien, sales, proyek, keuangan, dan komunikasi.

## Tech Stack

### Frontend
- **Framework**: React + TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **Routing**: React Router (based on pages structure)
- **Server**: Nginx (production)
- **Containerization**: Docker

### Backend
- **Language**: Golang
- **Structure**: Clean Architecture (cmd / internal)
- **Layers**: handlers → middleware → models → database → config
- **Containerization**: Docker
- **Orchestration**: docker-compose

## Project Structure

```
OneTool/
├── Frontend/
│   ├── src/
│   │   ├── pages/
│   │   │   ├── Auth/         # LoginPage, ForgotPasswordPage
│   │   │   ├── Dashboard/    # DashboardPage
│   │   │   ├── Clients/      # ClientsPage, ClientDetailPage
│   │   │   ├── Sales/        # StorePage, ItemsPage, OrdersPage,
│   │   │   │                 # ContractsPage, PaymentsPage,
│   │   │   │                 # InvoicesPage, InvoiceDetailPage
│   │   │   ├── Projects/     # ProjectsPage, ProjectDetailPage
│   │   │   ├── Tasks/        # TasksPage
│   │   │   ├── Todo/         # TodoPage
│   │   │   ├── Messages/     # MessagesPage
│   │   │   ├── Leads/        # LeadsPage
│   │   │   ├── Expenses/     # ExpensesPage
│   │   │   ├── Notes/        # NotesPage
│   │   │   └── Settings/     # UsersPage
│   │   ├── components/
│   │   │   ├── layout/       # Layout.tsx
│   │   │   └── common/       # ProtectedRoute, ManageLabelsModal, index
│   │   ├── utils/
│   │   │   └── format.ts     # utility functions
│   │   ├── App.tsx
│   │   ├── main.tsx
│   │   └── index.css
│   ├── index.html
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   ├── tsconfig.json
│   ├── nginx.conf
│   └── Dockerfile
├── Backend/
│   ├── cmd/
│   │   ├── api/main.go       # entry point API server
│   │   └── seed/main.go      # database seeder
│   └── internal/
│       ├── config/config.go
│       ├── database/database.go
│       ├── middleware/auth.go
│       ├── models/models.go
│       ├── handlers/
│       │   ├── auth.go
│       │   └── handlers.go
│       └── server/server.go
├── docker-compose.yml
└── README.md
```

## Fitur Utama
- **Auth**: Login, Forgot Password, Protected Routes, JWT Middleware
- **Clients**: Manajemen klien & detail klien
- **Sales**: Store, Items, Orders, Contracts, Payments, Invoices
- **Projects**: Manajemen proyek & detail proyek
- **Tasks & Todo**: Manajemen tugas
- **Leads**: Pipeline leads/prospek
- **Expenses**: Pencatatan pengeluaran
- **Messages**: Pesan internal
- **Notes**: Catatan
- **Settings**: Manajemen users

## Coding Conventions

### Frontend (React/TypeScript)
- Gunakan **functional components** dengan hooks
- File komponen menggunakan ekstensi `.tsx`
- Penamaan file: **PascalCase** (contoh: `ClientDetailPage.tsx`)
- Utility functions di `src/utils/`
- Shared components di `src/components/common/`
- Layout wrapper di `src/components/layout/`
- Jangan gunakan `any` di TypeScript kecuali terpaksa
- Gunakan Tailwind CSS untuk styling, hindari inline style

### Backend (Golang)
- Ikuti struktur **Clean Architecture**: cmd → internal → handlers → models
- Error handling eksplisit, jangan abaikan error
- Gunakan **idiomatic Go**: error sebagai return value kedua
- Config dari environment variables via `internal/config`
- Database logic di `internal/database`
- Business logic di `internal/handlers`
- Model/struct definisi di `internal/models`

## Hal yang JANGAN Dilakukan
- Jangan tambah dependency baru tanpa diskusi
- Jangan gunakan `any` di TypeScript
- Jangan bypass middleware auth untuk route yang butuh autentikasi
- Jangan hardcode config/credentials — gunakan `.env`
- Jangan campur business logic dengan database logic

## Environment
- File `.env` ada di root `OneTool/` (dibaca docker-compose) dan `Backend/.env` (untuk run manual)
- `DB_HOST=postgres` untuk Docker, `DB_HOST=localhost` untuk run manual
- Jangan commit file `.env` — sudah masuk `.gitignore`
- Version app ada di `Backend/internal/server/server.go` (health endpoint) dan `Frontend/src/components/layout/Layout.tsx` (sidebar)

## Cara Jalankan Manual (tanpa Docker)

### Prasyarat (install sekali)
```bash
brew install go postgresql@15
brew services start postgresql@15
psql postgres -c "CREATE USER cbqa WITH PASSWORD 'cbqa123';"
psql postgres -c "CREATE DATABASE cbqa_db OWNER cbqa;"
```

### Jalankan App
```bash
# Tab 1 - Backend
cd Backend
go run cmd/api/main.go

# Tab 2 - Frontend
cd Frontend
npm install   # pertama kali saja
npm run dev
```

Akses di http://localhost:3000

### Seed Database (pertama kali)
```bash
cd Backend
go run cmd/seed/main.go
```

### Login Default
| Role | Email | Password |
|------|-------|----------|
| Admin | admin@cbqa.com | Admin123! |
| Member | fauzi@cbqa.com | Member123! |

### Cek Health Backend
```bash
curl http://localhost:8080/health
# {"status":"ok","version":"1.0.0"}
```

## Docker (Full Stack)
```bash
# Pastikan .env ada di root OneTool/
docker-compose up --build   # build & jalankan semua
docker-compose down          # stop semua
docker-compose --profile tools up  # jalankan + pgAdmin (localhost:5050)
```

## GitHub
- Remote: https://github.com/Nexora-Tech-Team/ONETOOL
- Branch utama: `main`
