# OneTool - Project Context for Claude

## Overview
OneTool adalah aplikasi bisnis all-in-one yang mencakup manajemen klien, sales, proyek, keuangan, dan komunikasi. Dibangun untuk kebutuhan tim IT Audit.

**Versi saat ini: v1.0.2**

## Tech Stack

### Frontend
- **Framework**: React 18 + TypeScript 5
- **Build Tool**: Vite 5
- **Styling**: Tailwind CSS 3.4
- **Routing**: React Router v6
- **State Management**: Redux Toolkit + React Redux
- **HTTP Client**: Axios (centralized di `src/services/api.ts`)
- **UI Libraries**: Recharts, @hello-pangea/dnd, react-big-calendar, Lucide React, react-toastify
- **Server**: Nginx (production)
- **Containerization**: Docker

### Backend
- **Language**: Go 1.23
- **Framework**: Gin (HTTP router)
- **ORM**: GORM v2 + PostgreSQL 15
- **Auth**: JWT (golang-jwt/v5) + bcrypt
- **Structure**: Clean Architecture (cmd / internal)
- **Layers**: handlers → middleware → models → database → config
- **Containerization**: Docker
- **Orchestration**: docker-compose

## Project Structure

```
root/
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
│   │   │   ├── Tasks/        # TasksPage (List/Kanban/Gantt)
│   │   │   ├── Todo/         # TodoPage
│   │   │   ├── Messages/     # MessagesPage (stub - coming soon)
│   │   │   ├── Leads/        # LeadsPage (List/Kanban)
│   │   │   ├── Expenses/     # ExpensesPage
│   │   │   ├── Notes/        # NotesPage
│   │   │   ├── Events/       # EventsPage
│   │   │   ├── Files/        # FilesPage
│   │   │   ├── Reports/      # ReportsPage (+ Export CSV)
│   │   │   ├── Team/         # TeamMembersPage, TimeCardsPage,
│   │   │   │                 # LeavePage, AnnouncementsPage, HelpPage
│   │   │   └── Settings/     # UsersPage, AuditLogPage
│   │   ├── components/
│   │   │   ├── layout/       # Layout.tsx
│   │   │   └── common/       # ProtectedRoute, ManageLabelsModal, index
│   │   ├── services/
│   │   │   └── api.ts        # Semua API calls terpusat di sini
│   │   ├── store/
│   │   │   ├── index.ts
│   │   │   └── slices/       # authSlice, uiSlice
│   │   ├── utils/
│   │   │   └── format.ts     # isValidEmail, toISODate, formatNumber, parseNumber, formatIDR
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
│       ├── middleware/auth.go  # AuthRequired, AdminRequired, CORS
│       ├── models/models.go    # 22 entitas, semua pakai soft delete (+ AuditLog)
│       ├── handlers/
│       │   ├── auth.go
│       │   └── handlers.go    # Semua handler bisnis di sini
│       └── server/server.go   # Route setup (~65 endpoints)
├── docker-compose.yml
├── .gitignore
└── README.md
```

## Fitur Utama
- **Auth**: Login, Forgot Password (stub), Protected Routes, JWT Middleware
- **Dashboard**: Stats real-time (tasks, projects, invoices, team, clock in/out)
- **Clients**: Manajemen klien & detail klien (10 tab: overview, contacts, projects, invoices, dll)
- **Leads**: Pipeline leads/prospek (List + Kanban drag-drop), **konversi Lead → Client**
- **Sales**: Store, Items, Orders, Contracts, Payments, Invoices (dengan line items & payments)
- **Projects**: Manajemen proyek & detail proyek (List + Gantt)
- **Tasks**: Manajemen tugas (List + Kanban + Gantt)
- **Todo**: Per-user todo list
- **Events**: Kalender events
- **Files**: File manager (metadata only, storage belum diimplementasi)
- **Expenses**: Pencatatan pengeluaran per-user
- **Notes**: Catatan per-user
- **Team**: Members, Time Cards (clock in/out dengan duration), Leave Management, Announcements
- **Reports**: Invoice summary, Projects summary, Leads funnel, Expenses total, **Export CSV (Excel)**
- **Settings**: Manajemen users, **Audit Trail**

## Business Logic Penting

### Invoice Auto-Status
`recalcInvoice()` di `handlers.go` dipanggil setiap kali payment ditambah/dihapus atau item berubah.
Otomatis update `paid_amount`, `due_amount`, dan `status`:
- `paid >= total` → `fully_paid`
- `paid > 0` → `partially_paid`
- `paid == 0 && status != draft && due_date lewat` → `overdue`

### Lead → Client Conversion
`POST /api/v1/leads/:id/convert` — membuat Client baru dari data Lead, auto-set lead status ke `won`.
Di frontend: tombol "→ Client" di setiap baris list dan kanban card.

### Clock In/Out
- ClockIn: set `InTime` dan `InDate` ke `time.Now()`
- ClockOut: hitung `duration = now - InTime` dalam jam, simpan ke `TimeCard.Duration`

### Dashboard Stats
`GET /api/v1/dashboard` mengembalikan 22 field real-time:
task breakdown, invoice per-status amounts, project counts, clocked_in_count, on_leave_today, dll.

### Audit Trail (v1.0.2)
`recordAudit()` helper di `handlers.go` dipanggil otomatis di setiap operasi Create/Update/Delete/Convert untuk entitas: **Client, Project, Task, Lead, Invoice, Contract**.
- Model: `AuditLog` — menyimpan `user_id`, `action`, `entity_type`, `entity_id`, `entity_name`, `ip_address`, `created_at`
- Endpoint: `GET /api/v1/audit-logs` (Admin only) — support filter `entity_type`, `action`, `from`, `to`, `user_id`
- Frontend: halaman `Settings > Audit Trail` (`/settings/audit-log`)

### Invoice PDF Export (v1.0.2)
`GET /api/v1/invoices/:id/pdf` — menghasilkan HTML invoice yang di-style rapi.
Browser membuka di tab baru dan auto-trigger `window.print()` → user simpan sebagai PDF.
Tombol printer ikon ada di setiap baris tabel InvoicesPage.

### Reports Export CSV (v1.0.2)
`GET /api/v1/reports/export?type=invoices&year=2026` — generate CSV dengan UTF-8 BOM agar Excel baca karakter Indonesia dengan benar.
Tipe yang tersedia: `invoices`, `expenses`, `leads`, `projects`, `timecards`.
Tombol "Export Excel" ada di header halaman Reports, otomatis download sesuai tab aktif.

## API Services di Frontend (`src/services/api.ts`)
Setiap modul punya service object:
- `authService`, `dashboardService`, `clientService`, `projectService`, `taskService`
- `leadService`, `invoiceService`, `paymentService`, `contractService`, `itemService`
- `orderService`, `eventService`, `noteService`, `expenseService`, `fileService`
- `todoService`, `teamService`, `reportService`, `labelService`
- `auditService` — `list(params)` → `GET /api/v1/audit-logs`
- `invoicePDFService` — `openPDF(id)` → buka PDF di tab baru
- `reportService.exportCSV(type, year)` — download CSV

## Coding Conventions

### Frontend (React/TypeScript)
- Gunakan **functional components** dengan hooks
- File komponen menggunakan ekstensi `.tsx`
- Penamaan file: **PascalCase** (contoh: `ClientDetailPage.tsx`)
- Semua API calls melalui service object di `src/services/api.ts`
- Utility functions di `src/utils/format.ts`
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
- Semua model menggunakan `gorm.DeletedAt` (soft delete) — **kecuali `AuditLog` yang tidak pakai soft delete** (log harus permanen)
- Tanggal menggunakan tipe `FlexTime` (wrapper `time.Time` yang menerima RFC3339 dan `YYYY-MM-DD`)

## Hal yang JANGAN Dilakukan
- Jangan tambah dependency baru tanpa diskusi
- Jangan gunakan `any` di TypeScript
- Jangan bypass middleware auth untuk route yang butuh autentikasi
- Jangan hardcode config/credentials — gunakan `.env`
- Jangan campur business logic dengan database logic
- Jangan commit `Frontend/dist/`, `node_modules/`, `.env`, `Backend/api.exe`
- Jangan hapus atau soft-delete record `AuditLog` — log harus immutable

## Environment
- File `.env` ada di root repo (dibaca docker-compose) dan `Backend/.env` (untuk run manual)
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

Akses di http://localhost:3000 (atau 3001 jika port 3000 sedang dipakai)

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
# {"status":"ok","version":"1.0.2"}
```

## Docker (Full Stack)
```bash
# Pastikan .env ada di root
docker-compose up --build   # build & jalankan semua
docker-compose down          # stop semua
docker-compose --profile tools up  # jalankan + pgAdmin (localhost:5050)
```

## GitHub
- Remote: https://github.com/Nexora-Tech-Team/ONETOOL
- Branch utama: `main`
- Struktur root: `Backend/`, `Frontend/`, `docker-compose.yml`

## Changelog

### v1.0.2 (2026-04-26)
- Audit Trail: rekaman otomatis semua perubahan data (Client, Project, Task, Lead, Invoice, Contract)
- Invoice PDF Export: generate & print invoice sebagai PDF langsung dari browser
- Reports Export CSV: download laporan ke CSV/Excel untuk semua tab Reports
- Fix deploy: port 8080 conflict saat redeploy di GitHub Actions

### v1.0.1
- Initial production release
- Clean repo structure

### v1.0.0
- Init project
