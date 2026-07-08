import { ArrowRight, Database, FileSearch, ShieldCheck, Sparkles, Trophy } from 'lucide-react';

export default function HomePage({ navigate }) {
  return (
    <main className="page home-grid">
      <section className="hero-panel">
        <div>
          <span className="eyebrow">Swipe security trainer</span>
          <h1>Инспектор входящих</h1>
          <p>Разбирай письма как аналитик: текст, RAW-заголовки, ссылки и вложения в одном тренировочном режиме.</p>
          <div className="hero-metrics">
            <Metric value="120+" label="готовых писем" />
            <Metric value="5" label="категорий" />
            <Metric value="RAW" label="технический анализ" />
          </div>
        </div>
        <button className="primary-action hero-action" onClick={() => navigate('game')}>
          <span>Начать проверку</span>
          <ArrowRight size={18} />
        </button>
      </section>
      <section className="feature-grid">
        <Feature icon={ShieldCheck} title="Свайп-проверка" text="Влево фишинг, вправо легитимное письмо." />
        <Feature icon={Database} title="RAW-анализ" text="SPF, DKIM, DMARC, Reply-To и Received доступны в карточке." />
        <Feature icon={Sparkles} title="AI-категория" text="Новые письма генерируются на сервере через OpenAI API." />
        <Feature icon={FileSearch} title="Ссылки и вложения" text="Отдельные вкладки показывают реальный URL, домен, MIME-тип и риск." />
        <Feature icon={Trophy} title="Профиль и streak" text="Прогресс сохраняется в PostgreSQL и переживает перезапуск сервера." />
      </section>
    </main>
  );
}

function Metric({ value, label }) {
  return (
    <div>
      <strong>{value}</strong>
      <span>{label}</span>
    </div>
  );
}

function Feature({ icon: Icon, title, text }) {
  return (
    <div className="feature">
      <Icon size={22} />
      <strong>{title}</strong>
      <span>{text}</span>
    </div>
  );
}
