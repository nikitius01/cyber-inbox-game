import { useMemo, useState } from 'react';
import { Home, LogIn, Moon, Shield, Sun, UserRound } from 'lucide-react';
import { useAuth } from './context/AuthContext.jsx';
import { useTheme } from './context/ThemeContext.jsx';
import { useToast } from './context/ToastContext.jsx';
import HomePage from './pages/HomePage.jsx';
import LoginPage from './pages/LoginPage.jsx';
import RegisterPage from './pages/RegisterPage.jsx';
import GamePage from './pages/GamePage.jsx';
import ProfilePage from './pages/ProfilePage.jsx';

const routes = {
  home: HomePage,
  login: LoginPage,
  register: RegisterPage,
  game: GamePage,
  profile: ProfilePage,
};

export default function App() {
  const [route, setRoute] = useState('home');
  const { user, logout } = useAuth();
  const { resolvedTheme } = useTheme();
  const { showToast } = useToast();
  const Page = useMemo(() => routes[route] || HomePage, [route]);

  function navigate(nextRoute) {
    if (nextRoute === 'game' && !user) {
      showToast('Сначала войдите или создайте аккаунт.', 'warning');
      setRoute('login');
      return;
    }
    setRoute(nextRoute);
  }

  async function handleLogout() {
    await logout();
    showToast('Вы вышли из аккаунта.', 'info');
    setRoute('home');
  }

  return (
    <div className="app-shell">
      <nav className="nav">
        <button className="brand" onClick={() => navigate('home')}>
          <Shield size={22} />
          <span>Инспектор входящих</span>
        </button>
        <div className="nav-actions">
          <button className="icon-button" title={resolvedTheme === 'dark' ? 'Темная тема' : 'Светлая тема'}>
            {resolvedTheme === 'dark' ? <Moon size={18} /> : <Sun size={18} />}
          </button>
          <button className="icon-button" title="Игра" onClick={() => navigate('game')}>
            <Home size={18} />
          </button>
          {user ? (
            <>
              <button className="icon-button" title="Профиль" onClick={() => navigate('profile')}>
                <UserRound size={18} />
              </button>
              <button className="text-button" onClick={handleLogout}>Выйти</button>
            </>
          ) : (
            <button className="icon-button" title="Войти" onClick={() => navigate('login')}>
              <LogIn size={18} />
            </button>
          )}
        </div>
      </nav>
      <Page navigate={navigate} />
    </div>
  );
}
