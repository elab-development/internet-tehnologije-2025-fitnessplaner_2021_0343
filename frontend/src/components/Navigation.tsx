import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import Button from './Button';

const Navigation: React.FC = () => {
  const { isAuthenticated, user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <nav className="bg-white shadow-md mb-6">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center space-x-8">
            <Link to="/" className="text-2xl font-bold bg-gradient-to-r from-purple-600 to-purple-800 bg-clip-text text-transparent">
              Fitness App
            </Link>
            {isAuthenticated && (
              <div className="flex space-x-4">
                <Link to="/dashboard" className="text-gray-700 hover:text-purple-600 transition-colors">
                  Dashboard
                </Link>
                <Link to="/workouts" className="text-gray-700 hover:text-purple-600 transition-colors">
                  Workouts
                </Link>
                <Link to="/progress" className="text-gray-700 hover:text-purple-600 transition-colors">
                  Progress
                </Link>
                <Link to="/videos" className="text-gray-700 hover:text-purple-600 transition-colors">
                  Videos
                </Link>
                <Link to="/profile" className="text-gray-700 hover:text-purple-600 transition-colors">
                  Profile
                </Link>
              </div>
            )}
          </div>
          <div className="flex items-center space-x-4">
            {isAuthenticated ? (
              <>
                <span className="text-sm text-gray-600">
                  {user?.name} ({user?.role})
                </span>
                <Button variant="secondary" onClick={handleLogout}>
                  Logout
                </Button>
              </>
            ) : (
              <>
                <Link to="/login">
                  <Button variant="secondary">Login</Button>
                </Link>
                <Link to="/register">
                  <Button>Register</Button>
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navigation;

