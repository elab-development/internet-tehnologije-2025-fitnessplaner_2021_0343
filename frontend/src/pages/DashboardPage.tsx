import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { foodAPI, mealPlanAPI } from '../services/api';
import { Food, MealPlan } from '../services/types';

const DashboardPage: React.FC = () => {
  const { user, logout } = useAuth();
  const [barcode, setBarcode] = useState('');
  const [food, setFood] = useState<Food | null>(null);
  const [mealPlan, setMealPlan] = useState<MealPlan | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSearchFood = async () => {
    if (!barcode.trim()) {
      setError('Please enter a barcode');
      return;
    }

    setError('');
    setLoading(true);
    setFood(null);

    try {
      const result = await foodAPI.search(barcode);
      setFood(result);
    } catch (err: any) {
      setError(err.response?.data || err.message || 'Failed to search food');
    } finally {
      setLoading(false);
    }
  };

  const handleGenerateMealPlan = async () => {
    setError('');
    setLoading(true);
    setMealPlan(null);

    try {
      const result = await mealPlanAPI.generate();
      setMealPlan(result);
    } catch (err: any) {
      setError(err.response?.data || err.message || 'Failed to generate meal plan');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <div>
          <h1>Welcome, {user?.name}!</h1>
          <p>Your goal: <strong>{user?.goal === 'lose_weight' ? 'Lose Weight' : 'Hypertrophy'}</strong></p>
        </div>
        <button onClick={logout} className="logout-btn">Logout</button>
      </header>

      <div className="dashboard-content">
        {error && <div className="error">{error}</div>}

        {/* Food Search Section */}
        <section className="card">
          <h2>Search Food by Barcode</h2>
          <div className="search-form">
            <input
              type="text"
              value={barcode}
              onChange={(e) => setBarcode(e.target.value)}
              placeholder="Enter barcode (e.g., 3274080005003)"
              onKeyPress={(e) => e.key === 'Enter' && handleSearchFood()}
            />
            <button onClick={handleSearchFood} disabled={loading}>
              {loading ? 'Searching...' : 'Search'}
            </button>
          </div>

          {food && (
            <div className="food-result">
              <h3>{food.name}</h3>
              <div className="nutrition-grid">
                <div className="nutrition-item">
                  <span className="label">Calories:</span>
                  <span className="value">{food.calories.toFixed(2)} kcal</span>
                </div>
                <div className="nutrition-item">
                  <span className="label">Protein:</span>
                  <span className="value">{food.protein.toFixed(2)}g</span>
                </div>
                <div className="nutrition-item">
                  <span className="label">Carbs:</span>
                  <span className="value">{food.carbs.toFixed(2)}g</span>
                </div>
                <div className="nutrition-item">
                  <span className="label">Fat:</span>
                  <span className="value">{food.fat.toFixed(2)}g</span>
                </div>
              </div>
            </div>
          )}
        </section>

        {/* Meal Plan Section */}
        <section className="card">
          <h2>Your Meal Plan</h2>
          <button onClick={handleGenerateMealPlan} disabled={loading} className="generate-btn">
            {loading ? 'Generating...' : 'Generate Meal Plan'}
          </button>

          {mealPlan && (
            <div className="meal-plan">
              <div className="meal-plan-summary">
                <h3>Summary</h3>
                <div className="nutrition-grid">
                  <div className="nutrition-item">
                    <span className="label">Total Calories:</span>
                    <span className="value">{mealPlan.total_calories.toFixed(2)} kcal</span>
                  </div>
                  <div className="nutrition-item">
                    <span className="label">Total Protein:</span>
                    <span className="value">{mealPlan.total_protein.toFixed(2)}g</span>
                  </div>
                  <div className="nutrition-item">
                    <span className="label">Total Carbs:</span>
                    <span className="value">{mealPlan.total_carbs.toFixed(2)}g</span>
                  </div>
                  <div className="nutrition-item">
                    <span className="label">Total Fat:</span>
                    <span className="value">{mealPlan.total_fat.toFixed(2)}g</span>
                  </div>
                </div>
              </div>

              <div className="meal-plan-foods">
                <h3>Foods ({mealPlan.foods.length})</h3>
                {mealPlan.foods.map((food, index) => (
                  <div key={index} className="food-item">
                    <h4>{food.name}</h4>
                    <div className="nutrition-grid">
                      <span>Calories: {food.calories.toFixed(2)}</span>
                      <span>Protein: {food.protein.toFixed(2)}g</span>
                      <span>Carbs: {food.carbs.toFixed(2)}g</span>
                      <span>Fat: {food.fat.toFixed(2)}g</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </section>
      </div>
    </div>
  );
};

export default DashboardPage;
