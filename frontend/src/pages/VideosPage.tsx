import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import Card from '../components/Card';

<<<<<<< HEAD
// Predefined YouTube video IDs for different workout types
// Popular hypertrophy/muscle building videos
=======

// Popularni hypertrophy/muscle building videi
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
const HYPERTROPHY_VIDEOS = [
  { id: 'TLnVgSs1YXY', title: 'Full Body Hypertrophy Workout', description: 'Complete muscle building routine for all muscle groups' },
  { id: 'g_tea8ZNk5A', title: 'Push Pull Legs Split', description: 'PPL workout program for maximum muscle growth' },
  { id: 'UItWltVZZmE', title: 'Back & Biceps Hypertrophy', description: 'Upper body muscle building focused on back and biceps' },
  { id: 'eaLjN8x90kY', title: 'Chest & Triceps Workout', description: 'Push day hypertrophy training for chest and triceps' },
  { id: 'mlR6PBj8dB0', title: 'Legs & Glutes Hypertrophy', description: 'Lower body muscle building for legs and glutes' },
  { id: 'jH1b3vE3XqE', title: 'Shoulders & Arms Hypertrophy', description: 'Upper body specialization for shoulders and arms' },
];

<<<<<<< HEAD
// Popular cardio/weight loss videos
=======
// Popularni cardio/weight loss videi
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
const CARDIO_VIDEOS = [
  { id: 'mlR6PBj8dB0', title: '30 Min Full Body Cardio', description: 'High intensity fat burning workout for weight loss' },
  { id: 'jH1b3vE3XqE', title: 'HIIT Cardio Workout', description: 'High intensity interval training for maximum calorie burn' },
  { id: 'UItWltVZZmE', title: 'Fat Burning Cardio', description: 'Effective weight loss routine to burn calories' },
  { id: 'TLnVgSs1YXY', title: 'Low Impact Cardio', description: 'Beginner friendly cardio workout without jumping' },
  { id: 'g_tea8ZNk5A', title: 'Dance Cardio Workout', description: 'Fun and effective cardio session for weight loss' },
  { id: 'eaLjN8x90kY', title: 'Treadmill Cardio Routine', description: 'Cardio workout perfect for weight loss goals' },
];

const VideosPage: React.FC = () => {
  const { user } = useAuth();
  const [selectedCategory, setSelectedCategory] = useState<'hypertrophy' | 'cardio' | null>(null);

<<<<<<< HEAD
  // Determine which videos to show based on user's goal
=======
  
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
  const userGoal = user?.goal || 'lose_weight';
  const defaultCategory = userGoal === 'hypertrophy' ? 'hypertrophy' : 'cardio';
  const currentCategory = selectedCategory || defaultCategory;
  const videos = currentCategory === 'hypertrophy' ? HYPERTROPHY_VIDEOS : CARDIO_VIDEOS;

  return (
    <div className="w-full" style={{ position: 'relative', zIndex: 1 }}>
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-gray-800 mb-4">Workout Videos</h1>
        <p className="text-gray-600 mb-4">
          Based on your goal: <strong className="text-purple-600">{userGoal === 'hypertrophy' ? 'Muscle Building (Hypertrophy)' : 'Weight Loss (Cardio)'}</strong>
        </p>
        
<<<<<<< HEAD
        {/* Category selector */}
=======
        {/* Category selektor */}
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
        <div className="flex gap-4 mb-6">
          <button
            onClick={() => setSelectedCategory('hypertrophy')}
            className={`px-6 py-2 rounded-lg font-semibold transition-all ${
              currentCategory === 'hypertrophy'
                ? 'bg-gradient-to-r from-purple-600 to-purple-800 text-white shadow-lg'
                : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-300'
            }`}
          >
            Hypertrophy
          </button>
          <button
            onClick={() => setSelectedCategory('cardio')}
            className={`px-6 py-2 rounded-lg font-semibold transition-all ${
              currentCategory === 'cardio'
                ? 'bg-gradient-to-r from-purple-600 to-purple-800 text-white shadow-lg'
                : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-300'
            }`}
          >
            Cardio
          </button>
        </div>
      </div>

<<<<<<< HEAD
      {/* Videos Grid */}
=======
      {/* Grid sa videima */}
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {videos.map((video, index) => (
          <Card key={`${video.id}-${index}`} className="overflow-hidden">
            <div className="aspect-video w-full mb-4">
              <iframe
                width="100%"
                height="100%"
                src={`https://www.youtube.com/embed/${video.id}`}
                title={video.title}
                frameBorder="0"
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                allowFullScreen
                className="rounded-lg"
              ></iframe>
            </div>
            <h3 className="text-lg font-semibold text-gray-800 mb-2">{video.title}</h3>
            <p className="text-gray-600 text-sm">{video.description}</p>
          </Card>
        ))}
      </div>

      {videos.length === 0 && (
        <Card>
          <p className="text-center text-gray-500 py-8">No videos available for this category.</p>
        </Card>
      )}
    </div>
  );
};

export default VideosPage;
<<<<<<< HEAD

=======
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
