# Fitness Meal Plan Application

Full-stack aplikacija za praćenje ishrane, treninga i napretka.

## 🚀 Brzi Start

### 1. Migracije baze
```cmd
MIGRATE.bat
```

### 2. Pokretanje
```cmd
START.bat
```

Ili ručno:
```cmd
cd backend && go run main.go
cd frontend && npm install && npm run dev
```

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
