import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link, useLocation } from 'react-router-dom';
import apiClient from '../api/apiClient';

interface AssignedUser {
    id: number;
    username: string;
}

interface Task {
    id: number;
    project_id: number;
    creator_id: number;
    title: string;
    description: string | null;
    status: 'todo' | 'in_progress' | 'done';
    due_date: string | null;
    created_at: string;
    updated_at: string;
    assignees: AssignedUser[];
}

export default function TaskDetailsPage() {
    const navigate = useNavigate();
    const location = useLocation();
    const [task, setTask] = useState<Task | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const queryParams = new URLSearchParams(location.search);
    const taskID = queryParams.get("projectID");
    const projectID = queryParams.get("projectID");
    const projectName = queryParams.get("projectName");

    useEffect(() => {
        const fetchTaskDetails = async () => {
            const token = localStorage.getItem('authToken');
            if (!token) {
                navigate('/login');
                return;
            }

            try {
                const data = await apiClient.get(`tasks?projectID=${projectID}`);
                setTask(data);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'An error occurred.');
            } finally {
                setLoading(false);
            }
        };

        if (location.state?.task) {
            // if yes, use it directly and avoid an API call
            setTask(location.state.task);
            setLoading(false);
        } else {
            // if not (e.g., page was refreshed), fall back to fetching from the api
            fetchTaskDetails();
        }
    }, [taskID, location.state]);

    if (loading) return <div className="text-center text-white p-10">Loading task details...</div>;
    if (error) return <div className="text-center text-red-500 p-10">{error}</div>;
    if (!task) return <div className="text-center text-white p-10">Task not found.</div>;

    return (
        <div className="min-h-screen bg-zinc-900 text-white p-8">
            <header className="mb-8">
                <Link to={`/tasks?projectID=${projectID}&projectName=${projectName}`} className="text-blue-400 hover:underline">&larr; Back to Task List</Link>
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
                    <div className="flex flex-wrap gap-2">
                        {task.assignees.length > 0 ? (
                            task.assignees.map(user => (
                                <div key={user.id} className="bg-zinc-700 text-zinc-200 font-medium px-3 py-1 rounded-full">
                                    {user.username}
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
