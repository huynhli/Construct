import { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import apiClient from '../api/apiClient';

interface Project {
	id: number;
	company_id: number;
	name: string;
	description: string | null;
	due_date: string | null;
	created_at: string;
	updated_at: string;
}

export default function HomePage() {
	const [projects, setProjects] = useState<Project[]>([]);
	const [loading, setLoading] = useState(true);
	const [error, setError] = useState<string | null>(null);
	const navigate = useNavigate();

	useEffect(() => {
		const fetchProjects = async () => {
		const token = localStorage.getItem('authToken');
		console.log(token)
		if (!token) {
			navigate('/login');
			return;
		}

		try {
            const response = await apiClient.get('projects');
			setProjects(response);

        } catch (err) {
            setError(err instanceof Error ? err.message : 'An unexpected error occurred.');
        } finally {
            setLoading(false);
        }};

		fetchProjects();
	}, [navigate])

	const handleLogout = () => {
        localStorage.removeItem('authToken');
        navigate('/login');
    };
    
    if (loading) {
        return <div className="flex items-center justify-center min-h-screen bg-zinc-900 text-white">Loading projects...</div>;
    }

    return (
        <div className="min-h-screen bg-zinc-900 text-white">
            <header className="bg-zinc-800 shadow-md p-4 flex justify-between items-center">
                <h1 className="text-2xl font-bold">My Projects</h1>
				{error && <div className="bg-red-900 text-red-200 p-3 rounded-lg text-center mb-4">{error}</div>}
				<button
					onClick={handleLogout}
					className="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded-lg focus:outline-none focus:shadow-outline transition duration-300"
				>
					Logout
				</button>
            </header>
            
            <main className="p-8">
                {/* Project Gallery */}
                <div className="w-full max-w-4xl mx-auto">
					{projects.length > 0 ? (
						<div className="grid grid-cols-1 gap-6">
							{projects.map(project => (
								<Link to={`/tasks?projectID=${project.id}&projectName=${project.name}`} key={project.id}>
									<div className="bg-zinc-800 rounded-lg shadow-lg overflow-hidden transform hover:-translate-y-1 transition-transform duration-300 h-80 flex flex-col"> {/* Added fixed height and flex column */}
										<img
											src={`https://via.placeholder.com/800x400.png?text=${encodeURIComponent(project.name)}`}
											alt={project.name}
											className="w-full h-40 object-cover flex-shrink-0"
										/>
										<div className="p-4 flex-grow flex flex-col justify-between">
											<div>
												<h2 className="text-xl font-semibold mb-2">
													{project.name}
												</h2>
												<p className="text-zinc-400 text-sm line-clamp-2">
													{project.description || 'No description provided.'}
												</p>
											</div>
											{project.due_date ? (
												<div className="text-xs text-zinc-500 mt-2">
													<strong>Due:</strong> {new Date(project.due_date).toLocaleDateString()}
												</div>
											) : (
												<div className="text-xs text-zinc-500 mt-2">
													<strong>Due:</strong> Not specified
												</div>
											)}
										</div>
									</div>
								</Link>
							))}
						</div>
					) : (
						<div className="text-center text-zinc-300">
							<h2 className="text-2xl">No projects assigned</h2>
						</div>
					)}
                </div>
            </main>
        </div>
    );
}