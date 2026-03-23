import { Routes, Route, Outlet } from 'react-router-dom'
import HomePage from './pages/ProjectsPage.tsx'
import LoginPage from './pages/LoginPage.tsx'
import SignUpPage from './pages/SignUpPage.tsx'
import TaskListPage from './pages/TaskListPage.tsx'
import TaskDetailsPage from './pages/TaskDetailsPage.tsx'

export default function App() {
  // defining default layout
  const Layout = () => {
    return (
      <div>
        <Outlet/>
      </div>
    )
  }

  return (
    <div>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route path="/" element={<HomePage />} />
          <Route path="/login" element={<LoginPage />} />
          <Route path="/signup" element={<SignUpPage />} />
          <Route path="/tasks" element={<TaskListPage />} />
          <Route path="/taskDetails" element={<TaskDetailsPage />} />
        </Route>
      </Routes>
    </div>
  );
}
