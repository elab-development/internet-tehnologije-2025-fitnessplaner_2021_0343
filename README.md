# Fitness App

Full-stack fitness tracking application with React frontend and Go backend.

## Prerequisites

- **Go** 1.21 or higher
- **Node.js** 18 or higher
- **MySQL** 8.0 or higher
- **npm** or **yarn**

## Setup Instructions

### 1. Clone the Repository

```bash
git clone <your-repo-url>
cd proba123
```

### 2. Database Setup

#### Option A: Automatic Setup (Recommended)
The application will automatically create the database and tables when you start the backend.

#### Option B: Manual Setup
If you prefer to set up manually:

1. Start MySQL server
2. Create database:
```sql
CREATE DATABASE app_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. Backend Setup

1. Navigate to backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables (optional):
Create a `.env` file or set environment variables:
```bash
export DB_USER=root
export DB_PASSWORD=your_password
export DB_HOST=127.0.0.1
export DB_PORT=3306
export DB_NAME=app_db
```

**Windows (PowerShell):**
```powershell
$env:DB_USER="root"
$env:DB_PASSWORD="your_password"
$env:DB_HOST="127.0.0.1"
$env:DB_PORT="3306"
$env:DB_NAME="app_db"
```

**Windows (CMD):**
```cmd
set DB_USER=root
set DB_PASSWORD=your_password
set DB_HOST=127.0.0.1
set DB_PORT=3306
set DB_NAME=app_db
```

4. Run the backend:
```bash
go run main.go
```

The backend will:
- ✅ Connect to MySQL server
- ✅ Create database if it doesn't exist
- ✅ Create all required tables automatically
- ✅ Start server on `http://localhost:8080`

### 4. Frontend Setup

1. Navigate to frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Configure API URL (optional):
Create a `.env` file:
```env
VITE_API_URL=http://localhost:8080
```

4. Run the frontend:
```bash
npm run dev
```

The frontend will start on `http://localhost:5173` (or another port if 5173 is busy).

## Default Configuration

- **Backend Port**: 8080
- **Frontend Port**: 5173 (Vite default)
- **Database**: app_db
- **Database User**: root
- **Database Password**: Vojislav123! (change this in production!)

## Features

- ✅ User registration and authentication
- ✅ Workout tracking
- ✅ Progress tracking (weight, body fat, muscle mass)
- ✅ YouTube workout videos based on user goals
- ✅ Profile management
- ✅ Automatic database migrations

## Troubleshooting

### Database Connection Issues

1. **MySQL not running**: Start MySQL service
2. **Wrong credentials**: Check environment variables or defaults in `backend/utils/database.go`
3. **Database doesn't exist**: The app will create it automatically, but ensure MySQL user has CREATE DATABASE permission

### Port Already in Use

- Backend: Change port in `backend/main.go`
- Frontend: Vite will automatically use next available port

### Migration Issues

The app automatically creates tables on startup. If you see table errors:
1. Check MySQL logs
2. Ensure user has CREATE TABLE permissions
3. Check database connection in logs

## Project Structure

```
proba123/
├── backend/
│   ├── controllers/     # API handlers
│   ├── models/          # Data models
│   ├── routes/          # Route definitions
│   ├── utils/           # Database and utilities
│   ├── migrations/      # SQL migration files
│   └── main.go          # Entry point
├── frontend/
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── pages/       # Page components
│   │   ├── services/    # API services
│   │   └── contexts/    # React contexts
│   └── package.json
└── README.md
```

## License

[Your License Here]
