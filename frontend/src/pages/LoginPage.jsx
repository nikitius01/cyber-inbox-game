import { useState } from 'react';
import { useAuth } from '../context/AuthContext.jsx';
import { useToast } from '../context/ToastContext.jsx';

export default function LoginPage({ navigate }) {
  const { login } = useAuth();
  const { showToast } = useToast();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  async function submit(event) {
    event.preventDefault();
    setError('');
    try {
      await login(email, password);
      showToast('Вход выполнен. Можно проверять письма.', 'success');
      navigate('game');
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <main className="page auth-page">
      <form className="auth-form" onSubmit={submit}>
        <h1>Вход</h1>
        <label>Email<input value={email} onChange={(e) => setEmail(e.target.value)} type="email" required /></label>
        <label>Пароль<input value={password} onChange={(e) => setPassword(e.target.value)} type="password" required /></label>
        {error && <div className="form-error">{error}</div>}
        <button className="primary-action">Войти</button>
        <button type="button" className="link-button" onClick={() => navigate('register')}>Создать аккаунт</button>
      </form>
    </main>
  );
}
