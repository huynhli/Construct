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
	const [pageLoading, setPageLoading] = useState<boolean>(true);
	const [pageLoadError, setPageLoadError] = useState<string | null>(null);
	const navigate = useNavigate();

	const [showProjectCreationModal, setShowProjectCreationModal] = useState<boolean>(false)
	const [projectCreationForm, setProjectCreationForm] = useState<Pick<Project, 'company_id' | 'name' | 'description' | 'due_date'>>({
		company_id: 0,
		name: '',
		description: null,
		due_date: null,
	})
	const [projectCreationFormError, setProjectCreationFormError] = useState<string | null>(null)
	const [projectCreationFormDisabled, setProjectCreationFormDisabled] = useState<boolean>(false)
	const [showProjectCreated, setShowProjectCreated] = useState<0 | 1 | 2>(0)

	useEffect(() => {
		const fetchProjects = async () => {
			// Check for auth token
			const token = localStorage.getItem('authToken');
			console.log(token)
			if (!token) {
				navigate('/login');
				return;
			}
			
			// Fetch projects from API
			try {
				const response = await apiClient.get('projects');
				setProjects(response);
			} catch (err) {
				setPageLoadError(err instanceof Error ? err.message : 'An unexpected error occurred.');
			} finally {
				setPageLoading(false);
			}
		};

		fetchProjects();
	}, [navigate])

	const handleLogout = () => {
        localStorage.removeItem('authToken');
        navigate('/login');
    };

	const createProject = () => {
		setShowProjectCreationModal(true)
	}

	const disableProjectCreationModal = () => {
		setShowProjectCreationModal(false)
	}

	const handleProjectCreationFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
		// Reject past dates if selected
		if (e.target.name === 'due_date' && e.target.value) {
			const today = new Date();
			today.setHours(0, 0, 0, 0);
			if (new Date(e.target.value) < today) setProjectCreationFormError("Invalid due date. Please select a future date.");
		}
		setProjectCreationForm((prev: Pick<Project, 'company_id' | 'name' | 'description' | 'due_date'>) => ({...prev, [e.target.name]: e.target.value}))
	}
	
	const handleProjectCreationFormSubmit = (e: React.FormEvent) => {
		e.preventDefault()
		setProjectCreationFormDisabled(true)

		submitProjectCreationForm()
		setProjectCreationFormDisabled(false)
		
	}

	const submitProjectCreationForm = async () => {
		try {
			await apiClient.post('projects', projectCreationForm);
			
			setShowProjectCreated(1)
			// cause projects page to refresh
		} catch (err) {
			setShowProjectCreated(2)
		}
	}


	if (pageLoadError) [
		<div className="flex items-center justify-center min-h-screen bg-red-400 text-white">Error: {pageLoadError}</div>
	]
    
    if (pageLoading) {
        return <div className="flex items-center justify-center min-h-screen bg-zinc-900 text-white">Loading projects...</div>
    }

    return (
        <div className="min-h-screen bg-zinc-900 text-white flex flex-col">
			
			{/* Header */}
            <header className="bg-zinc-800 shadow-md p-4 flex justify-between items-center">
                <h1 className="text-2xl font-bold">My Projects</h1>
				{/* {pageLoadError && <div className="bg-red-900 text-red-200 p-3 rounded-lg text-center mb-4">{pageLoadError}</div>} */}
				<button
					onClick={handleLogout}
					className="bg-red-600 hover:bg-red-700 text-white font-bold py-2 px-4 rounded-lg focus:outline-none focus:shadow-outline transition duration-300"
				>
					Logout
				</button>
            </header>
            
			{/* Main */}
            <main className="p-8 flex-1 relative">
				
				{/* Project banners */}
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

				{/* Project creation button */}
				<button 
					className='
						absolute bottom-10 right-10 
						[clip-path:circle(50%)] p-[1%]
						bg-zinc-200 hover:bg-black
						hover:scale-110
						transition-all duration-100
						active:bg-zinc-400
						'
					onClick={createProject}
				>
					<svg xmlns="http://www.w3.org/2000/svg" className="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
						<path d="M12 5v14M5 12h14"/>
					</svg>
				</button>

				{/* Project creation modal */}
				{ showProjectCreationModal && 

					// Backdrop
					<div 
						className='fixed inset-0 bg-black/50 z-50 flex items-center justify-center'
						// onClick={disableProjectCreationModal} --> everything outside modal disables, need modal to have onClick={(e) => e.stopPropagation()}
					>
						
						{/* Modal box */}
						<div className='bg-zinc-700 p-6 rounded-xl w-[85%] h-[90%]'>
							<div className='flex justify-end w-full'>
								<button 
									className='
										[clip-path:circle(50%)] p-2
										bg-red-400 hover:bg-red-500 active:bg-red-600
										hover:scale-120
										transition-all duration-100'
									onClick={disableProjectCreationModal}
								>
									<svg xmlns="http://www.w3.org/2000/svg" className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
										<path d="M18 6L6 18M6 6l12 12"/>
									</svg>
								</button>
							</div>
							
							
							<form className='w-[100%] flex h-[90%] text-[clamp(0.875rem,1.5vw,1.25rem)] overflow-y-auto' onSubmit={handleProjectCreationFormSubmit}>
								<fieldset disabled={projectCreationFormDisabled} className='w-full h-full flex flex-col items-center'>
									<h1 className='text-2xl sm:text-4xl font-bold mb-4 w-[90%] sm:w-[80%]'>Create a Project</h1>
									<div className='flex w-[95%] sm:w-[80%] my-2'>
										<h2 className='h-full mr-6 sm:mr-2 w-[20%]'>Company ID:</h2>
										<input className='sm:w-[75%] max-w-[75%] bg-white text-black px-1' 
											name='company_id' value={projectCreationForm.company_id} onChange={handleProjectCreationFormChange} 
											type='number' placeholder='Enter Company ID...' 
											/>
									</div>
									<div className='flex w-[95%] sm:w-[80%] m-2'>
										<h2 className='h-full mr-6 sm:mr-2 w-[20%]'>Project name:</h2>
										<input className='sm:w-[75%] max-w-[75%] bg-white text-black px-1'
											name='name' value={projectCreationForm.name} onChange={handleProjectCreationFormChange} 
											type='text' placeholder='Enter project name...' required
											/>
									</div>
									<div className='flex w-[95%] sm:w-[80%] h-[40%] m-2'>
										<h2 className='h-full mr-6 sm:mr-2 w-[20%]'>Description (optional):</h2>
										<textarea
											className='sm:w-[75%] max-w-[75%] bg-white text-black px-1 resize-none overflow-y-auto'
											name='description'
											value={projectCreationForm.description ?? ''}
											onChange={handleProjectCreationFormChange}
											placeholder='Enter description... (optional)'
											rows={4}
											/>
									</div>
									<div className='flex w-[95%] sm:w-[80%] mb-6 m-2'>
										<h2 className='h-full mr-6 sm:mr-2 w-[20%]'>Due date (optional):</h2>
										<input className='sm:w-[75%] max-w-[75%] bg-white text-black px-1' 
											name='due_date' value={projectCreationForm.due_date ?? ''} onChange={handleProjectCreationFormChange} 
											type='date'
											/>
									</div>
									{projectCreationFormError && 
										<div className="bg-red-200 text-black p-3 rounded-lg text-center mb-4">{projectCreationFormError}</div>
									}
									<button type='submit' className='bg-gray-400 hover:bg-gray-500 active:bg-gray-600 py-2 px-4 border-gray-100 border-1 rounded-md'>
										Submit
									</button>
								</fieldset>
							</form>
						</div>
					</div>
				}

				{/* Project created popup */}
				{showProjectCreated !== 0 &&
					<div className='bg-purple-200 p-6 rounded-xl w-[65%] h-[60%]'></div>
				}

            </main>
        </div>
    );
}