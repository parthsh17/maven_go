import type { OrderEvent } from '../types';
import { Badge } from './Badge';

interface StateTimelineProps {
  events: OrderEvent[];
}

export function StateTimeline({ events }: StateTimelineProps) {
  if (events.length === 0) {
    return (
      <div className="empty-state">
        <span>◎</span>
        <p>No state transitions recorded yet</p>
      </div>
    );
  }

  return (
    <div className="timeline">
      {events.map((evt, i) => (
        <div key={i} className="timeline-item">
          <div className="timeline-dot" />
          {i < events.length - 1 && <div className="timeline-line" />}
          <div className="timeline-content">
            <div className="timeline-states">
              <Badge state={evt.previous_state} />
              <span className="timeline-arrow">→</span>
              <Badge state={evt.new_state} />
            </div>
            {evt.message && <p className="timeline-message">{evt.message}</p>}
            <span className="timeline-time">
              {new Date(evt.timestamp).toLocaleString()}
            </span>
          </div>
        </div>
      ))}
    </div>
  );
}
