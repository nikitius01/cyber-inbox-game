import { useEffect, useState } from 'react';
import { Flame, ListChecks, Palette, UserRound } from 'lucide-react';
import { profileApi } from '../api/profileApi.js';
import ProtectedRoute from '../components/ProtectedRoute.jsx';
import { useTheme } from '../context/ThemeContext.jsx';
import { useToast } from '../context/ToastContext.jsx';

const categories = ['easy', 'medium', 'hard', 'nightmare', 'AI'];

export default function ProfilePage({ navigate }) {
  return (
    <ProtectedRoute fallback={<LoginRequired navigate={navigate} />}>
      <ProfileContent />
    </ProtectedRoute>
  );
}

function ProfileContent() {
  const [profile, setProfile] = useState(null);
  const [error, setError] = useState('');
  const { preference, setPreference } = useTheme();
  const { showToast } = useToast();

  useEffect(() => {
    profileApi.get().then(setProfile).catch((err) => setError(err.message));
  }, []);

  if (error) return <main className="page"><div className="form-error">{error}</div></main>;
  if (!profile) return <main className="page">Загрузка...</main>;

  const stats = profile.stats || {};
  const solved = stats.solvedByCategory || {};

  return (
    <main className="page profile-layout">
      <section className="profile-head">
        <UserRound size={44} />
        <div>
          <h1>{profile.username}</h1>
          <span>{profile.email}</span>
        </div>
      </section>

      <section className="stats-grid">
        <Metric icon={Flame} label="Streak" value={`${stats.streak || 0} дней`} />
        <Metric icon={ListChecks} label="Всего решено" value={stats.totalSolved || 0} />
      </section>

      <section className="theme-panel">
        <div>
          <Palette size={24} />
          <h2>Тема интерфейса</h2>
        </div>
        <div className="theme-options">
          {[
            ['light', 'Светлая'],
            ['dark', 'Темная'],
            ['auto', 'Автоматически'],
          ].map(([value, label]) => (
            <button
              key={value}
              className={preference === value ? 'selected' : ''}
              onClick={() => {
                setPreference(value);
                showToast(`Тема: ${label.toLowerCase()}.`, 'success');
              }}
            >
              {label}
            </button>
          ))}
        </div>
      </section>

      <section className="category-stats">
        {categories.map((category) => (
          <div key={category}>
            <span>{category}</span>
            <strong>{solved[category] || 0}</strong>
          </div>
        ))}
      </section>
    </main>
  );
}

function Metric({ icon: Icon, label, value }) {
  return (
    <div className="metric">
      <Icon size={22} />
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

function LoginRequired({ navigate }) {
  return (
    <main className="page auth-page">
      <div className="auth-form">
        <h1>Нужен вход</h1>
        <button className="primary-action" onClick={() => navigate('login')}>Войти</button>
      </div>
    </main>
  );
}
