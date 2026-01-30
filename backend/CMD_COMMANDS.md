# 🚀 CMD Komande za Pokretanje

## 1. Migracije Baze (Prvo jednom)

### Opcija A - Batch fajl:
```cmd
MIGRATE.bat
```

### Opcija B - Ručno u CMD:
```cmd
mysql -u root -p < backend\migrations\001_init.sql
```

Ili u MySQL klijentu:
```sql
CREATE DATABASE IF NOT EXISTS app_db;
USE app_db;
```
Zatim kopiraj i pokreni SQL iz `backend\migrations\001_init.sql`

---

## 2. Pokretanje Aplikacije

### Opcija A - Batch fajl (Preporučeno):
```cmd
START.bat
```

### Opcija B - Ručno u CMD:

**Terminal 1 - Backend:**
```cmd
cd backend
go mod tidy
go run main.go
```

**Terminal 2 - Frontend (novi CMD prozor):**
```cmd
cd frontend
npm install
npm run dev
```

---

## 3. Pristup Aplikaciji

- **Backend API:** http://localhost:8080
- **Frontend:** http://localhost:5173
- **Health Check:** http://localhost:8080/health

---

## 4. Zaustavljanje

Pritisni `Ctrl+C` u CMD prozoru gde server radi.

---

## ⚠️ Troubleshooting

### Backend ne radi:
```cmd
# Proveri Go
go version

# Proveri MySQL kredencijale u backend\utils\database.go
```

### Frontend ne radi:
```cmd
# Proveri Node.js
node --version
npm --version

# Reinstaliraj dependencies
cd frontend
npm install
```

### Greška sa bazom:
```cmd
# Proveri da li MySQL radi
# Proveri da li je baza kreirana
# Pokreni migracije ponovo
```

