import { motion, useMotionValue, useTransform } from 'framer-motion';
import { Check, X } from 'lucide-react';
import EmailCard from './EmailCard.jsx';

export default function SwipeDeck({ task, result, onAnswer }) {
  const x = useMotionValue(0);
  const rotate = useTransform(x, [-220, 220], [-8, 8]);

  function handleDragEnd(_, info) {
    if (result) return;
    if (info.offset.x < -120) onAnswer('phishing');
    if (info.offset.x > 120) onAnswer('legitimate');
  }

  return (
    <div className="deck">
      <motion.div
        className="swipe-card"
        drag={!result ? 'x' : false}
        dragConstraints={{ left: 0, right: 0 }}
        style={{ x, rotate }}
        onDragEnd={handleDragEnd}
      >
        <EmailCard task={task} result={result} />
      </motion.div>
      <div className="answer-bar">
        <button className="danger" disabled={!!result} onClick={() => onAnswer('phishing')} title="Фишинг">
          <X size={20} />
          <span>Фишинг</span>
        </button>
        <button className="safe" disabled={!!result} onClick={() => onAnswer('legitimate')} title="Легитимное">
          <Check size={20} />
          <span>Легитимное</span>
        </button>
      </div>
    </div>
  );
}

