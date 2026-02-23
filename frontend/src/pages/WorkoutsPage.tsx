import React, { useState, useEffect } from 'react';
import { workoutAPI } from '../services/api';
import { Workout } from '../services/types';
import Card from '../components/Card';
import Button from '../components/Button';
import Modal from '../components/Modal';
import Input from '../components/Input';

const WorkoutsPage: React.FC = () => {
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingWorkout, setEditingWorkout] = useState<Workout | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    duration: '',
    calories_burned: '',
    workout_date: new Date().toISOString().split('T')[0],
  });

<<<<<<< HEAD
  // Fetch workouts on component mount
=======
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
  useEffect(() => {
    fetchWorkouts();
  }, []);

  const fetchWorkouts = async () => {
    try {
      setLoading(true);
      setError('');
      const data = await workoutAPI.getAll();
      setWorkouts(data || []);
    } catch (err: any) {
      console.error('Error fetching workouts:', err);
<<<<<<< HEAD
      // Handle JSON error response from backend
=======
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
      const errorData = err.response?.data;
      let errorMessage = 'Failed to fetch workouts';
      if (errorData?.message) {
        errorMessage = errorData.message;
      } else if (errorData?.error) {
        errorMessage = errorData.error;
      } else if (typeof errorData === 'string') {
        errorMessage = errorData;
      } else if (err.message) {
        errorMessage = err.message;
      }
      setError(errorMessage);
<<<<<<< HEAD
      setWorkouts([]); // Set empty array on error
=======
      setWorkouts([]);
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
    } finally {
      setLoading(false);
    }
  };

  const handleOpenModal = (workout?: Workout) => {
    if (workout) {
      setEditingWorkout(workout);
      setFormData({
        name: workout.name,
        description: workout.description || '',
        duration: workout.duration.toString(),
        calories_burned: workout.calories_burned.toString(),
        workout_date: workout.workout_date.split('T')[0],
      });
    } else {
      setEditingWorkout(null);
      setFormData({
        name: '',
        description: '',
        duration: '',
        calories_burned: '',
        workout_date: new Date().toISOString().split('T')[0],
      });
    }
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingWorkout(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingWorkout) {
        await workoutAPI.update(editingWorkout.id!, {
          name: formData.name,
          description: formData.description,
          duration: parseInt(formData.duration),
          calories_burned: parseFloat(formData.calories_burned),
          workout_date: formData.workout_date,
        });
      } else {
        await workoutAPI.create({
          name: formData.name,
          description: formData.description,
          duration: parseInt(formData.duration),
          calories_burned: parseFloat(formData.calories_burned),
          workout_date: formData.workout_date,
        });
      }
      handleCloseModal();
      fetchWorkouts();
    } catch (err: any) {
<<<<<<< HEAD
      // Handle JSON error response from backend
=======
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
      const errorData = err.response?.data;
      if (errorData?.message) {
        setError(errorData.message);
      } else if (typeof errorData === 'string') {
        setError(errorData);
      } else {
        setError(err.message || 'Failed to save workout');
      }
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this workout?')) return;
    
    try {
      await workoutAPI.delete(id);
      fetchWorkouts();
    } catch (err: any) {
      // Handle JSON error response from backend
      const errorData = err.response?.data;
      if (errorData?.message) {
        setError(errorData.message);
      } else if (typeof errorData === 'string') {
        setError(errorData);
      } else {
        setError(err.message || 'Failed to delete workout');
      }
    }
  };

  return (
    <div className="w-full" style={{ position: 'relative', zIndex: 1 }}>
      {loading && (
        <div className="text-center py-8 text-gray-600">Loading workouts...</div>
      )}
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-800">My Workouts</h1>
        <Button onClick={() => handleOpenModal()}>Add Workout</Button>
      </div>

      {error && (
        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
          <strong>Error:</strong> {error}
        </div>
      )}

      {!loading && workouts.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {workouts.map((workout) => (
            <Card key={workout.id} title={workout.name}>
              <p className="text-gray-600 mb-2">{workout.description || 'No description'}</p>
              <div className="space-y-1 text-sm">
                <p><strong>Duration:</strong> {workout.duration} minutes</p>
                <p><strong>Calories:</strong> {workout.calories_burned.toFixed(0)} kcal</p>
                <p><strong>Date:</strong> {new Date(workout.workout_date).toLocaleDateString()}</p>
              </div>
              <div className="flex gap-2 mt-4">
                <Button variant="secondary" onClick={() => handleOpenModal(workout)}>
                  Edit
                </Button>
                <Button variant="danger" onClick={() => handleDelete(workout.id!)}>
                  Delete
                </Button>
              </div>
            </Card>
          ))}
        </div>
      )}

      {!loading && workouts.length === 0 && !error && (
        <Card>
          <p className="text-center text-gray-500 py-8">No workouts yet. Add your first workout!</p>
        </Card>
      )}

      <Modal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        title={editingWorkout ? 'Edit Workout' : 'Add Workout'}
        size="md"
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Workout Name"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            required
            placeholder="e.g., Morning Run"
          />
          <Input
            label="Description"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="Workout description"
          />
          <Input
            label="Duration (minutes)"
            type="number"
            value={formData.duration}
            onChange={(e) => setFormData({ ...formData, duration: e.target.value })}
            required
            min={1}
          />
          <Input
            label="Calories Burned"
            type="number"
            value={formData.calories_burned}
            onChange={(e) => setFormData({ ...formData, calories_burned: e.target.value })}
            required
            min={0}
            step="0.1"
          />
          <Input
            label="Date"
            type="date"
            value={formData.workout_date}
            onChange={(e) => setFormData({ ...formData, workout_date: e.target.value })}
            required
          />
          <div className="flex gap-2 justify-end">
            <Button type="button" variant="secondary" onClick={handleCloseModal}>
              Cancel
            </Button>
            <Button type="submit">
              {editingWorkout ? 'Update' : 'Create'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default WorkoutsPage;
<<<<<<< HEAD

=======
>>>>>>> 4dcc7f38d3ca50ba631e57486728f6fe45021608
