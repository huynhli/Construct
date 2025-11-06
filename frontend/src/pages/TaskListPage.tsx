import { useState, useEffect } from 'react';
import { useLocation, useNavigate, Link } from 'react-router-dom';
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


export default function TaskListPage() {
    const [tasks, setTasks] = useState<Task[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    const location = useLocation();
    const queryParams = new URLSearchParams(location.search);
    const projectID = queryParams.get("projectID");
    const projectName = queryParams.get("projectName");

    useEffect(() => {

        const fetchTasks = async () => {
            const token = localStorage.getItem('authToken');
            if (!token) {
                navigate('/login');
                return;
            }

            try {
                const response = await apiClient.get(`tasks?projectID=${projectID}`)
                console.log(response)
                setTasks(response);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'An error occurred.');
            } finally {
                setLoading(false);
            }
        };

        fetchTasks();
    }, [projectID, navigate]);

    if (loading) return <div className="text-center text-white p-10">Loading tasks...</div>;
    if (error) return <div className="text-center text-red-500 p-10">{error}</div>;
    if (!tasks) return <div className="text-center text-white p-10">Tasks not found.</div>;

    return (
        <div className="min-h-screen bg-zinc-900 text-white p-8">
            <header className="mb-8">
                <Link to="/" className="text-blue-400 hover:underline">&larr; Back to Projects</Link>
                <h1 className="text-4xl font-bold mt-2">Tasks for: {projectName}</h1>
            </header>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {tasks.length > 0 ? (
                    tasks.map(task => (
                        <Link
                            to={`/taskDetails?projectID=${projectID}&projectName=${projectName}&taskID=${task.id}`}
                            key={task.id}
                            state={{ task: task }}
                            className="block bg-zinc-800 p-6 rounded-lg shadow-lg hover:bg-zinc-700 transition-colors duration-300">
                            <h2 className="text-xl font-semibold truncate">{task.title}</h2>
                            <p className={`mt-2 px-2 py-1 text-xs inline-block rounded-full ${task.status === 'done' ? 'bg-green-700' : task.status === 'in_progress' ? 'bg-yellow-700' : 'bg-gray-600'}`}>
                                {task.status.replace('_', ' ')}
                            </p>
                            <div className="mt-4 text-sm text-zinc-400">
                                <strong>Due:</strong> {task.due_date ? new Date(task.due_date).toLocaleDateString() : 'Not set'}
                            </div>
                        </Link>
                    ))
                ) : (
                    <p className="text-zinc-500 col-span-full text-center">No tasks found for this project.</p>
                )}
            </div>
        </div>
    );
}