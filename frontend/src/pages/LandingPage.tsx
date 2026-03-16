import { Link } from 'react-router-dom';

export function LandingPage() {
  return (
    <div className="landing">
      <div className="landing-bg" />
      <header className="landing-header">
        <div className="brand-logo">
          <span className="brand-icon-lg">◈</span>
          <span className="brand-name-lg">Maven</span>
        </div>
      </header>

      <main className="hero">
        <div className="hero-badge">Order Lifecycle Platform</div>
        <h1 className="hero-title">
          Institutional-Grade<br />
          <span className="gradient-text">Order Management</span>
        </h1>
        <p className="hero-subtitle">
          Maven is a concurrency-safe, state-driven order lifecycle engine built in Go —
          with an operational control plane you're looking at right now.
        </p>
        <div className="hero-actions">
          <Link to="/signup" className="btn btn-primary btn-lg">Get Started</Link>
          <Link to="/login" className="btn btn-outline btn-lg">Sign In</Link>
        </div>

        <div className="feature-grid">
          {[
            { icon: '⚡', title: 'Concurrent Processing', desc: 'Worker pool with goroutines, channels & WaitGroups' },
            { icon: '🔄', title: 'State Machine', desc: 'Strict lifecycle: CREATED → VALIDATED → QUEUED → EXECUTING → COMPLETED' },
            { icon: '🛡️', title: 'Race-Free', desc: 'Mutex-protected in-memory store, zero data races' },
            { icon: '📊', title: 'Live Metrics', desc: 'Real-time counters for all states and worker activity' },
            { icon: '↩️', title: 'Auto Retry', desc: 'Failed orders automatically retried up to 3× with state tracking' },
            { icon: '📋', title: 'Event Audit Log', desc: 'Every state transition logged with timestamp and message' },
          ].map((f) => (
            <div className="feature-card" key={f.title}>
              <div className="feature-icon">{f.icon}</div>
              <h3>{f.title}</h3>
              <p>{f.desc}</p>
            </div>
          ))}
        </div>
      </main>

      <footer className="landing-footer">
        <p>© 2026 Maven · Concurrency-Safe Order Lifecycle Platform</p>
      </footer>
    </div>
  );
}
