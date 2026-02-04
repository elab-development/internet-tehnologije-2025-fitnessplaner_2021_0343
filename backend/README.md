# Fitness Meal Plan Application - Backend

Full-stack aplikacija za praćenje ishrane, treninga i napretka.

## 🚀 Brzi Start

### Automatska Setup (Preporučeno)
Aplikacija automatski kreira bazu i tabele pri pokretanju!

```bash
# 1. Instaliraj zavisnosti
go mod download

# 2. Pokreni backend (kreira bazu i tabele automatski)
go run main.go
```

### Ručna Setup (Opciono)

Ako želiš da ručno kreiraš bazu:

```sql
CREATE DATABASE app_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### Konfiguracija

Postavi environment varijable (opciono):
```bash
export DB_USER=root
export DB_PASSWORD=your_password
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_NAME=app_db
```

Ili koristi default vrednosti iz `utils/database.go`.

## 📁 Struktura

```
backend/
├── auth/              # JWT (1 fajl)
├── controllers/       # 3 kontrolera (user, food, data)
├── middleware/        # 1 fajl (sve middleware)
├── models/           # 4 modela
├── migrations/       # 1 migracija
├── routes/           # Rute
└── utils/            # Database

frontend/
├── src/
│   ├── components/   # Reusable komponente
│   ├── pages/        # Stranice
│   └── services/     # API servisi
```

## 📡 API Endpoints

**Public:** `/api/register`, `/api/login`, `/health`

**Protected (JWT):** `/api/profile`, `/api/logout`, `/api/food/search`, `/api/meal-plan`, `/api/workouts/*`, `/api/progress/*`

## 🔧 Konfiguracija

- Backend: `utils/database.go` (MySQL)
- Frontend: `.env` (API URL)
- Ports: Backend 8080, Frontend 5173
