import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';

// Re-using the same interfaces
interface AssignedUser {
    id: number;
    name: string;
    avatar_url?: string;
}

interface TaskDetails {
    id: number;
    title: string;
    description: string | null;
    status: 'todo' | 'in_progress' | 'done';
    due_date: string | null;
    created_at: string;
    assigned_users: AssignedUser[];
}

export default function TaskDetailsPage() {
    const { projectId, taskId } = useParams<{ projectId: string; taskId: string }>();
    const [task, setTask] = useState<TaskDetails | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchTaskDetails = async () => {
            const token = localStorage.getItem('authToken');
            if (!token) {
                navigate('/login');
                return;
            }

            try {
                // NOTE: Use your actual API endpoint to fetch a single task's details
                const response = await fetch(`/api/tasks/${taskId}`, {
                    headers: { 'Authorization': `Bearer ${token}` },
                });

                if (!response.ok) {
                    throw new Error('Failed to load task details.');
                }
                const data: TaskDetails = await response.json();
                setTask(data);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'An error occurred.');
            } finally {
                setLoading(false);
            }
        };

        fetchTaskDetails();
    }, [taskId, navigate]);

    if (loading) return <div className="text-center text-white p-10">Loading task details...</div>;
    if (error) return <div className="text-center text-red-500 p-10">{error}</div>;
    if (!task) return <div className="text-center text-white p-10">Task not found.</div>;

    return (
        <div className="min-h-screen bg-zinc-900 text-white p-8">
            <header className="mb-8">
                <Link to={`/projects/${projectId}/tasks`} className="text-blue-400 hover:underline">&larr; Back to Task List</Link>
            </header>

            <main className="max-w-4xl mx-auto bg-zinc-800 p-8 rounded-lg shadow-lg">
                <div className="flex justify-between items-start">
                    <h1 className="text-4xl font-bold">{task.title}</h1>
                    <span className={`px-3 py-1 text-sm rounded-full ${task.status === 'done' ? 'bg-green-700' : task.status === 'in_progress' ? 'bg-yellow-700' : 'bg-gray-600'}`}>
                        {task.status.replace('_', ' ')}
                    </span>
                </div>
                
                <div className="text-zinc-400 mt-2">
                    <strong>Due Date:</strong> {task.due_date ? new Date(task.due_date).toLocaleDateString() : 'Not specified'}
                </div>
                
                <div className="mt-8 border-t border-zinc-700 pt-6">
                    <h2 className="text-2xl font-semibold mb-4">Description</h2>
                    <p className="text-zinc-300 whitespace-pre-wrap">
                        {task.description || 'No description provided.'}
                    </p>
                </div>
                
                <div className="mt-8 border-t border-zinc-700 pt-6">
                    <h2 className="text-2xl font-semibold mb-4">Assigned To</h2>
                    <div className="flex flex-wrap gap-4">
                        {task.assigned_users.length > 0 ? (
                            task.assigned_users.map(user => (
                                <div key={user.id} className="flex items-center gap-3 bg-zinc-700 p-2 rounded-lg">
                                    <img src={user.avatar_url || 'https://via.placeholder.com/40'} alt={user.name} className="w-10 h-10 rounded-full" />
                                    <span className="font-medium">{user.name}</span>
                                </div>
                            ))
                        ) : (
                            <p className="text-zinc-500">No one is assigned to this task yet.</p>
                        )}
                    </div>
                </div>
            </main>
        </div>
    );
}
