import React, { useState, useEffect } from 'react';
import { progressAPI } from '../services/api';
import { Progress } from '../services/types';
import Card from '../components/Card';
import Button from '../components/Button';
import Modal from '../components/Modal';
import Input from '../components/Input';

const ProgressPage: React.FC = () => {
  const [progressList, setProgressList] = useState<Progress[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingProgress, setEditingProgress] = useState<Progress | null>(null);
  const [formData, setFormData] = useState({
    weight: '',
    body_fat: '',
    muscle_mass: '',
    notes: '',
    progress_date: new Date().toISOString().split('T')[0],
  });

  // Fetch progress on component mount
  useEffect(() => {
    fetchProgress();
  }, []);

  const fetchProgress = async () => {
    try {
      setLoading(true);
      const data = await progressAPI.getAll();
      setProgressList(data);
      setError('');
    } catch (err: any) {
      setError(err.response?.data || err.message || 'Failed to fetch progress');
    } finally {
      setLoading(false);
    }
  };

  const handleOpenModal = (progress?: Progress) => {
    if (progress) {
      setEditingProgress(progress);
      setFormData({
        weight: progress.weight.toString(),
        body_fat: progress.body_fat?.toString() || '',
        muscle_mass: progress.muscle_mass?.toString() || '',
        notes: progress.notes || '',
        progress_date: progress.progress_date.split('T')[0],
      });
    } else {
      setEditingProgress(null);
      setFormData({
        weight: '',
        body_fat: '',
        muscle_mass: '',
        notes: '',
        progress_date: new Date().toISOString().split('T')[0],
      });
    }
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingProgress(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingProgress) {
        await progressAPI.update(editingProgress.id!, {
          weight: parseFloat(formData.weight),
          body_fat: formData.body_fat ? parseFloat(formData.body_fat) : 0,
          muscle_mass: formData.muscle_mass ? parseFloat(formData.muscle_mass) : 0,
          notes: formData.notes,
          progress_date: formData.progress_date,
        });
      } else {
        await progressAPI.create({
          weight: parseFloat(formData.weight),
          body_fat: formData.body_fat ? parseFloat(formData.body_fat) : 0,
          muscle_mass: formData.muscle_mass ? parseFloat(formData.muscle_mass) : 0,
          notes: formData.notes,
          progress_date: formData.progress_date,
        });
      }
      handleCloseModal();
      fetchProgress();
    } catch (err: any) {
      setError(err.response?.data || err.message || 'Failed to save progress');
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this progress entry?')) return;
    
    try {
      await progressAPI.delete(id);
      fetchProgress();
    } catch (err: any) {
      setError(err.response?.data || err.message || 'Failed to delete progress');
    }
  };

  if (loading) {
    return <div className="text-center py-8">Loading progress...</div>;
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-800">My Progress</h1>
        <Button onClick={() => handleOpenModal()}>Add Progress</Button>
      </div>

      {error && <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">{error}</div>}

      {progressList.length === 0 ? (
        <Card>
          <p className="text-center text-gray-500 py-8">No progress entries yet. Add your first entry!</p>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {progressList.map((progress) => (
            <Card key={progress.id} title={new Date(progress.progress_date).toLocaleDateString()}>
              <div className="space-y-2">
                <p><strong>Weight:</strong> {progress.weight} kg</p>
                {progress.body_fat > 0 && <p><strong>Body Fat:</strong> {progress.body_fat}%</p>}
                {progress.muscle_mass > 0 && <p><strong>Muscle Mass:</strong> {progress.muscle_mass} kg</p>}
                {progress.notes && <p className="text-gray-600">{progress.notes}</p>}
              </div>
              <div className="flex gap-2 mt-4">
                <Button variant="secondary" onClick={() => handleOpenModal(progress)}>
                  Edit
                </Button>
                <Button variant="danger" onClick={() => handleDelete(progress.id!)}>
                  Delete
                </Button>
              </div>
            </Card>
          ))}
        </div>
      )}

      <Modal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        title={editingProgress ? 'Edit Progress' : 'Add Progress'}
        size="md"
      >
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            label="Weight (kg)"
            type="number"
            value={formData.weight}
            onChange={(e) => setFormData({ ...formData, weight: e.target.value })}
            required
            min={0}
            step="0.1"
          />
          <Input
            label="Body Fat (%)"
            type="number"
            value={formData.body_fat}
            onChange={(e) => setFormData({ ...formData, body_fat: e.target.value })}
            min={0}
            max={100}
            step="0.1"
          />
          <Input
            label="Muscle Mass (kg)"
            type="number"
            value={formData.muscle_mass}
            onChange={(e) => setFormData({ ...formData, muscle_mass: e.target.value })}
            min={0}
            step="0.1"
          />
          <Input
            label="Notes"
            value={formData.notes}
            onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
            placeholder="Additional notes"
          />
          <Input
            label="Date"
            type="date"
            value={formData.progress_date}
            onChange={(e) => setFormData({ ...formData, progress_date: e.target.value })}
            required
          />
          <div className="flex gap-2 justify-end">
            <Button type="button" variant="secondary" onClick={handleCloseModal}>
              Cancel
            </Button>
            <Button type="submit">
              {editingProgress ? 'Update' : 'Create'}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default ProgressPage;

