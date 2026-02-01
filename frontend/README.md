# Frontend - Fitness Meal Plan App

React + TypeScript + Vite aplikacija za plan ishrane.

## 🚀 Pokretanje

### 1. Instalacija zavisnosti

```bash
npm install
```

### 2. Konfiguracija

Kreiraj `.env` fajl (ili koristi `env.example`):
```
VITE_API_URL=http://localhost:8080
```

### 3. Pokretanje development servera

```bash
npm run dev
```

Frontend će biti dostupan na `http://localhost:5173` (ili drugom portu ako je 5173 zauzet).

## 📡 Povezivanje sa Backend-om

Backend API servis je već konfigurisan u `src/services/api.ts` i automatski se povezuje sa backend-om na `http://localhost:8080`.

### API Funkcije

**Autentifikacija:**
- `authAPI.register()` - Registracija korisnika
- `authAPI.login()` - Login
- `authAPI.logout()` - Logout
- `authAPI.getProfile()` - Dobijanje profila

**Hrana:**
- `foodAPI.search(barcode)` - Pretraga hrane po barcodu

**Plan ishrane:**
- `mealPlanAPI.generate()` - Generisanje plana ishrane

### Korišćenje Auth Context-a

```tsx
import { useAuth } from './contexts/AuthContext';

function MyComponent() {
  const { user, login, logout, isAuthenticated } = useAuth();
  
  // ...
}
```

## 📁 Struktura

```
src/
├── components/    # React komponente
├── contexts/      # React Context (AuthContext)
├── pages/         # Stranice aplikacije
└── services/      # API servisi
    ├── api.ts     # Axios konfiguracija i API funkcije
    └── types.ts   # TypeScript tipovi
```

## 🔗 Backend Endpoints

Frontend koristi sledeće backend endpoint-e:

- `POST /api/register` - Registracija
- `POST /api/login` - Login
- `GET /api/profile` - Profil (zaštićeno)
- `POST /api/food/search` - Pretraga hrane (zaštićeno)
- `GET /api/meal-plan` - Plan ishrane (zaštićeno)

Svi zaštićeni endpoint-i zahtevaju JWT token u `Authorization` header-u.
