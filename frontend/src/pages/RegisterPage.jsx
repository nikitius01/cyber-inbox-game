import { useState } from 'react';
import { useAuth } from '../context/AuthContext.jsx';
import { useToast } from '../context/ToastContext.jsx';

export default function RegisterPage({ navigate }) {
  const { register } = useAuth();
  const { showToast } = useToast();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  async function submit(event) {
    event.preventDefault();
    setError('');
    try {
      await register(username, email, password);
      showToast('Аккаунт создан. Добро пожаловать.', 'success');
      navigate('game');
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <main className="page auth-page">
      <form className="auth-form" onSubmit={submit}>
        <h1>Регистрация</h1>
        <label>Имя<input value={username} onChange={(e) => setUsername(e.target.value)} minLength="3" required /></label>
        <label>Email<input value={email} onChange={(e) => setEmail(e.target.value)} type="email" required /></label>
        <label>Пароль<input value={password} onChange={(e) => setPassword(e.target.value)} type="password" minLength="8" required /></label>
        {error && <div className="form-error">{error}</div>}
        <button className="primary-action">Создать аккаунт</button>
        <button type="button" className="link-button" onClick={() => navigate('login')}>Уже есть аккаунт</button>
      </form>
    </main>
  );
}
