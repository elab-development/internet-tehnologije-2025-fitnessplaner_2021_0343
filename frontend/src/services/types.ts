// User types
export interface User {
  id: number;
  name: string;
  email: string;
  goal: 'lose_weight' | 'hypertrophy';
  role?: 'admin' | 'user' | 'premium';
  height?: number;
  weight?: number;
  created_at?: string;
  updated_at?: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
  goal: 'lose_weight' | 'hypertrophy';
  height?: number;
  weight?: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  token: string;
}

// Food types
export interface Food {
  id?: number;
  name: string;
  barcode?: string;
  calories: number;
  protein: number;
  carbs: number;
  fat: number;
}

export interface FoodSearchRequest {
  barcode: string;
}

// Meal Plan types
export interface MealPlan {
  id?: number;
  user_id: number;
  goal: string;
  foods: Food[];
  total_calories: number;
  total_protein: number;
  total_carbs: number;
  total_fat: number;
}

// Workout types
export interface Workout {
  id?: number;
  user_id: number;
  name: string;
  description?: string;
  duration: number;
  calories_burned: number;
  workout_date: string;
  created_at?: string;
  updated_at?: string;
}

// Progress types
export interface Progress {
  id?: number;
  user_id: number;
  weight: number;
  body_fat?: number;
  muscle_mass?: number;
  notes?: string;
  progress_date: string;
  created_at?: string;
  updated_at?: string;
}

