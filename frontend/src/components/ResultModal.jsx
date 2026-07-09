import { BadgeCheck, RotateCcw, TriangleAlert } from 'lucide-react';
import { safeText } from '../utils/sanitize.js';

export default function ResultModal({ result, onNext }) {
  if (!result) return null;
  return (
    <aside className="result-panel">
      <div>
        <span className={result.isCorrect ? 'result-pill ok' : 'result-pill bad'}>
          {result.isCorrect ? 'Верный разбор' : 'Разбор ошибки'}
        </span>
        <h2>{result.correctAnswer === 'phishing' ? 'Это фишинг' : 'Это легитимное письмо'}</h2>
        <div className={result.isCorrect ? 'feedback-box ok' : 'feedback-box bad'}>
          {result.isCorrect ? <BadgeCheck size={20} /> : <TriangleAlert size={20} />}
          <div>
            <strong>{safeText(result.feedbackTitle || (result.isCorrect ? 'Верно' : 'Почему неверно'))}</strong>
            <p>{safeText(result.feedbackMessage || result.explanation)}</p>
          </div>
        </div>
      </div>
      <div className="flags">
        {result.redFlags?.length ? result.redFlags.map((flag, index) => (
          <div className="flag" key={index}>
            <strong>{safeText(flag.type)} · {safeText(flag.field)}</strong>
            <span>{safeText(flag.explanation)}</span>
          </div>
        )) : <div className="flag">Красных флагов нет.</div>}
      </div>
      <button className="primary-action" onClick={onNext}>
        <RotateCcw size={18} />
        <span>Следующее письмо</span>
      </button>
    </aside>
  );
}
