import { useState } from 'react';
import { AlertTriangle, FileArchive, Link2, Mail, Rows3 } from 'lucide-react';
import DifficultyBadge from './DifficultyBadge.jsx';
import { safeText } from '../utils/sanitize.js';

const tabs = [
  { id: 'mail', label: 'Письмо', icon: Mail },
  { id: 'raw', label: 'RAW', icon: Rows3 },
  { id: 'links', label: 'Ссылки', icon: Link2 },
  { id: 'attachments', label: 'Вложения', icon: FileArchive },
];

export default function EmailCard({ task, result }) {
  const [tab, setTab] = useState('mail');

  return (
    <article className="email-card">
      <header className="email-head">
        <div>
          <DifficultyBadge value={task.category} />
          <h1>{safeText(task.subject)}</h1>
        </div>
        {result && (
          <span className={result.isCorrect ? 'result-pill ok' : 'result-pill bad'}>
            {result.isCorrect ? 'Верно' : 'Ошибка'}
          </span>
        )}
      </header>

      <div className="sender-grid">
        <span>От</span>
        <strong>{safeText(task.senderName)} &lt;{safeText(task.senderEmail)}&gt;</strong>
        <span>Сигналы</span>
        <strong>{task.links?.length || 0} ссылок · {task.attachments?.length || 0} вложений · RAW доступен</strong>
      </div>

      <div className="tabs">
        {tabs.map((item) => {
          const Icon = item.icon;
          return (
            <button key={item.id} className={tab === item.id ? 'active' : ''} onClick={() => setTab(item.id)}>
              <Icon size={16} />
              <span>{item.label}</span>
            </button>
          );
        })}
      </div>

      <section className="tab-panel">
        {tab === 'mail' && <p className="mail-body">{safeText(task.body)}</p>}
        {tab === 'raw' && <RawView raw={task.raw} result={result} />}
        {tab === 'links' && <LinksView links={task.links} result={result} />}
        {tab === 'attachments' && <AttachmentsView attachments={task.attachments} result={result} />}
      </section>
    </article>
  );
}

function RawView({ raw, result }) {
  const headers = raw?.headers || {};
  return (
    <div className="raw-view">
      {Object.entries(headers).map(([key, value]) => (
        <div key={key} className={result && ['SPF', 'DKIM', 'DMARC', 'Reply-To'].includes(key) ? 'raw-line inspect' : 'raw-line'}>
          <span>{key}</span>
          <code>{safeText(Array.isArray(value) ? value.join(' | ') : value)}</code>
        </div>
      ))}
      <pre>{safeText(raw?.source)}</pre>
    </div>
  );
}

function LinksView({ links, result }) {
  if (!links?.length) return <Empty text="Ссылок нет" />;
  return (
    <div className="data-list">
      {links.map((link, index) => (
        <div key={index} className={result && link.isSuspicious ? 'data-row risk' : 'data-row'}>
          <Link2 size={18} />
          <div>
            <strong>{safeText(link.visibleText)}</strong>
            <span>{safeText(link.actualUrl)}</span>
            {result && link.riskReason && <em>{safeText(link.riskReason)}</em>}
          </div>
        </div>
      ))}
    </div>
  );
}

function AttachmentsView({ attachments, result }) {
  if (!attachments?.length) return <Empty text="Вложений нет" />;
  return (
    <div className="data-list">
      {attachments.map((file, index) => (
        <div key={index} className={result && file.isSuspicious ? 'data-row risk' : 'data-row'}>
          <AlertTriangle size={18} />
          <div>
            <strong>{safeText(file.displayName || file.fileName)}</strong>
            <span>{safeText(file.mimeType)} · {file.sizeKb} KB · .{safeText(file.extension)}</span>
            {result && file.riskReason && <em>{safeText(file.riskReason)}</em>}
          </div>
        </div>
      ))}
    </div>
  );
}

function Empty({ text }) {
  return <div className="empty-state">{text}</div>;
}
