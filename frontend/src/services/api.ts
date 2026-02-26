import axios from 'axios';

// Get API URL iz environment-a ili korist default
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

// Kreiraj axios instancu
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Dodaj token zahtevu ako je dostupan
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Handle token expiration
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      // Redirect to login if needed
      if (window.location.pathname !== '/login') {
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);

// Auth API
export const authAPI = {
  register: async (data: {
    name: string;
    email: string;
    password: string;
    goal: 'lose_weight' | 'hypertrophy';
  }) => {
    const response = await api.post('/api/register', data);
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    return response.data;
  },

  login: async (data: { email: string; password: string }) => {
    const response = await api.post('/api/login', data);
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
    }
    return response.data;
  },

  logout: async () => {
    try {
      await api.post('/api/logout');
    } catch (error) {
      console.error('Logout API error:', error);
    } finally {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
  },

  getProfile: async () => {
    const response = await api.get('/api/profile');
    return response.data;
  },
};

// Food API
export const foodAPI = {
  search: async (barcode: string) => {
    const response = await api.post('/api/food/search', { barcode });
    return response.data;
  },
};

// Meal Plan API
export const mealPlanAPI = {
  generate: async () => {
    const response = await api.get('/api/meal-plan');
    return response.data;
  },
};

// Workout API
export const workoutAPI = {
  getAll: async () => {
    const response = await api.get('/api/workouts');
    return response.data;
  },
  create: async (data: {
    name: string;
    description?: string;
    duration: number;
    calories_burned: number;
    workout_date: string;
  }) => {
    const response = await api.post('/api/workouts/create', data);
    return response.data;
  },
  update: async (id: number, data: {
    name: string;
    description?: string;
    duration: number;
    calories_burned: number;
    workout_date: string;
  }) => {
    const response = await api.put(`/api/workouts/update?id=${id}`, data);
    return response.data;
  },
  delete: async (id: number) => {
    const response = await api.delete(`/api/workouts/delete?id=${id}`);
    return response.data;
  },
};

// Progress API
export const progressAPI = {
  getAll: async () => {
    const response = await api.get('/api/progress');
    return response.data;
  },
  create: async (data: {
    weight: number;
    body_fat?: number;
    muscle_mass?: number;
    notes?: string;
    progress_date: string;
  }) => {
    const response = await api.post('/api/progress/create', data);
    return response.data;
  },
  update: async (id: number, data: {
    weight: number;
    body_fat?: number;
    muscle_mass?: number;
    notes?: string;
    progress_date: string;
  }) => {
    const response = await api.put(`/api/progress/update?id=${id}`, data);
    return response.data;
  },
  delete: async (id: number) => {
    const response = await api.delete(`/api/progress/delete?id=${id}`);
    return response.data;
  },
};

// Health check
export const healthCheck = async () => {
  const response = await api.get('/health');
  return response.data;
};

export default api;
