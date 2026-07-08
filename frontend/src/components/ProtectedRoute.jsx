import { useAuth } from '../context/AuthContext.jsx';

export default function ProtectedRoute({ children, fallback }) {
  const { user, loading } = useAuth();
  if (loading) return <div className="page">Загрузка...</div>;
  if (!user) return fallback;
  return children;
}

