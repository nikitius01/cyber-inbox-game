import { useEffect, useState } from 'react';
import { FileArchive, Link2, Loader2, MailCheck, Rows3 } from 'lucide-react';
import SwipeDeck from '../components/SwipeDeck.jsx';
import ResultModal from '../components/ResultModal.jsx';
import { taskApi } from '../api/taskApi.js';
import ProtectedRoute from '../components/ProtectedRoute.jsx';
import { useToast } from '../context/ToastContext.jsx';

const categories = ['easy', 'medium', 'hard', 'nightmare', 'AI'];

export default function GamePage() {
  return (
    <ProtectedRoute fallback={<LoginRequired />}>
      <GameContent />
    </ProtectedRoute>
  );
}

function GameContent() {
  const [category, setCategory] = useState('easy');
  const [task, setTask] = useState(null);
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const { showToast } = useToast();

  useEffect(() => {
    loadTask(category);
  }, [category]);

  async function loadTask(nextCategory = category) {
    setLoading(true);
    setError('');
    setResult(null);
    try {
      setTask(await taskApi.random(nextCategory));
    } catch (err) {
      setTask(null);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function answer(value) {
    if (!task || result) return;
    try {
      setResult(await taskApi.answer(task, value));
      showToast(value === 'phishing' ? 'Ответ принят: фишинг.' : 'Ответ принят: легитимное.', 'info');
    } catch (err) {
      setError(err.message);
    }
  }

  return (
    <main className="page game-layout">
      <aside className="side-panel">
        <h2>Категория</h2>
        <div className="category-list">
          {categories.map((item) => (
            <button key={item} className={item === category ? 'selected' : ''} onClick={() => setCategory(item)}>
              {item}
            </button>
          ))}
        </div>
        <div className="inspection-guide">
          <h3>Проверяй</h3>
          <GuideItem icon={MailCheck} text="тон письма и срочность" />
          <GuideItem icon={Rows3} text="SPF, DKIM, DMARC" />
          <GuideItem icon={Link2} text="видимый и реальный URL" />
          <GuideItem icon={FileArchive} text="расширение и MIME вложений" />
        </div>
      </aside>

      <section className="game-stage">
        {loading && <div className="loading"><Loader2 className="spin" />Загрузка письма...</div>}
        {error && <div className="form-error">{error}</div>}
        {task && !loading && <SwipeDeck task={task} result={result} onAnswer={answer} />}
      </section>

      <ResultModal result={result} onNext={() => loadTask()} />
    </main>
  );
}

function GuideItem({ icon: Icon, text }) {
  return (
    <div className="guide-item">
      <Icon size={16} />
      <span>{text}</span>
    </div>
  );
}

function LoginRequired() {
  return (
    <main className="page auth-page">
      <div className="auth-form pop-in">
        <h1>Сначала регистрация</h1>
        <p>Игра сохраняет streak и статистику, поэтому проверка писем доступна только после входа.</p>
      </div>
    </main>
  );
}
