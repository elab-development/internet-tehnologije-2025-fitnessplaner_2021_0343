import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import Button from '../components/Button';
import Input from '../components/Input';
import Card from '../components/Card';

const RegisterPage: React.FC = () => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [goal, setGoal] = useState<'lose_weight' | 'hypertrophy'>('lose_weight');
  const [height, setHeight] = useState('');
  const [weight, setWeight] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await register({ 
        name, 
        email, 
        password, 
        goal,
        height: height ? parseFloat(height) : undefined,
        weight: weight ? parseFloat(weight) : undefined,
      });
      navigate('/dashboard');
    } catch (err: any) {
      // Handle JSON error response from backend
      const errorData = err.response?.data;
      if (errorData?.message) {
        setError(errorData.message);
      } else if (typeof errorData === 'string') {
        setError(errorData);
      } else {
        setError(err.message || 'Registration failed');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex justify-center items-center min-h-[calc(100vh-200px)]">
      <Card className="w-full max-w-md" title="Register">
        {error && <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">{error}</div>}
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Name"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            required
            placeholder="Your Name"
          />
          <Input
            label="Email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            placeholder="your@email.com"
          />
          <Input
            label="Password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={6}
            placeholder="••••••••"
          />
          <div className="form-group">
            <label className="block mb-2 text-sm font-medium text-gray-700">
              Goal <span className="text-red-500">*</span>
            </label>
            <select
              value={goal}
              onChange={(e) => setGoal(e.target.value as 'lose_weight' | 'hypertrophy')}
              required
              className="w-full px-4 py-2 border-2 border-gray-300 rounded-lg focus:outline-none focus:border-purple-500"
            >
              <option value="lose_weight">Lose Weight</option>
              <option value="hypertrophy">Hypertrophy</option>
            </select>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <Input
              label="Height (cm)"
              type="number"
              value={height}
              onChange={(e) => setHeight(e.target.value)}
              min={50}
              max={250}
              step="0.1"
              placeholder="170"
            />
            <Input
              label="Weight (kg)"
              type="number"
              value={weight}
              onChange={(e) => setWeight(e.target.value)}
              min={30}
              max={300}
              step="0.1"
              placeholder="70"
            />
          </div>
          <Button type="submit" disabled={loading} className="w-full">
            {loading ? 'Registering...' : 'Register'}
          </Button>
        </form>
        <p className="text-center mt-4 text-gray-600">
          Already have an account? <Link to="/login" className="text-purple-600 hover:underline">Login</Link>
        </p>
      </Card>
    </div>
  );
};

export default RegisterPage;

