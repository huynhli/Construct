import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

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
	const navigate = useNavigate();

    return (
        <div className="min-h-screen bg-zinc-900 text-white">
            <header className="bg-zinc-800 shadow-md p-4 flex justify-between items-center">
                <h1 className="text-2xl font-bold">My Projects</h1>
            </header>
            
            <main className="p-8">
                {/* Project Gallery */}
                <div className="w-full max-w-4xl mx-auto space-y-8">
					{projects.length > 0 ? (
						projects.map(project => (
							<div key={project.id} className="bg-zinc-800 rounded-lg shadow-lg overflow-hidden transform hover:-translate-y-2 transition-transform duration-300">
								<img
									src={'https://upload.wikimedia.org/wikipedia/commons/1/14/Brr_brr_patapim.jpg'} // Fallback image
									alt={project.name}
									className="w-fit h-64 object-cover"
								/>
								<div className="p-6">
									<h2 className="text-2xl font-semibold mb-2">{project.name}</h2>
									<p className="text-zinc-400 mb-4">
										{project.description || "No description"}
									</p>

									<div className="text-sm text-sinz-500">
										</div>
								</div>
							</div>
						)

						)
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