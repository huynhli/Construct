import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Link } from 'react-router-dom';

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
	const navigate = useNavigate();

	useEffect(() => {
		const fetchProjects = async () => {
		const token = localStorage.getItem('authToken');
		if (!token) {
			navigate('/login');
			return;
		}

		await fetch("/REPLACE/WITH/ENDPOINT", {
			headers: {
				'Authorization': 'Bearer ${token}'
			}
		})
			.then((res) => res.json())
			.then((json) => setProjects(json))
			.catch((err) => console.error(err));

		setLoading(false);
			
		};

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
				<button
					onClick={handleLogout}
					className="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded-lg focus:outline-none focus:shadow-outline transition duration-300"
				>
					Logout
				</button>
            </header>
            
            <main className="p-8">
                {/* Project Gallery */}
                <div className="w-full max-w-4xl mx-auto space-y-8">
					{projects.length > 0 ? (

						projects.map(project => (
							<Link to={`/projects/${project.id}/tasks`} key={project.id}>
								<div className="bg-zinc-800 rounded-lg shadow-lg overflow-hidden transform hover:-translate-y-2 transition-transform duration-300">
									<img
										src={`https://via.placeholder.com/800x400.png?text=${encodeURIComponent(project.name)}`}
										alt={project.name}
										className="w-full h-64 object-cover"
									/>
									<div className="p-6">
										<h2 className="text-2xl font-semibold mb-2">{project.name}</h2>
										<p className="text-zinc-400 mb-4">
											{project.description || 'No description provided.'}
										</p>
										{project.due_date ? (
											<div className="text-sm text-zinc-500">
												<strong>Due:</strong> {new Date(project.due_date).toLocaleDateString()}
											</div>
										) : (
											<div className="text-sm text-zinc-500">
												<strong>Due:</strong> Not specified
											</div>
										)}
									</div>
								</div>
							</Link>
						))
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