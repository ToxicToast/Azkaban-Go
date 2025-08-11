package healthmon

import (
	"context"
	"sync"
	"time"
)

type RedisPinger interface {
	Ping(ctx context.Context) error
}
type KafkaChecker interface {
	Check(ctx context.Context) error
}

type Status struct {
	OK        bool
	Error     string
	CheckedAt time.Time
}

type Monitor struct {
	mu sync.RWMutex

	// cached statuses
	service map[string]Status // gRPC downstreams, key = service name (z.B. "warcraft")
	redis   Status
	kafka   Status

	// deps
	services      []string
	grpcPingFn    func(ctx context.Context, service string) error
	redisPinger   RedisPinger
	kafkaChecker  KafkaChecker
	intervalGRPC  time.Duration
	intervalRedis time.Duration
	intervalKafka time.Duration
}

func NewMonitor(services []string, grpcPingFn func(ctx context.Context, service string) error,
	intervalGRPC, intervalRedis, intervalKafka time.Duration) *Monitor {

	return &Monitor{
		service:       make(map[string]Status),
		services:      services,
		grpcPingFn:    grpcPingFn,
		intervalGRPC:  intervalGRPC,
		intervalRedis: intervalRedis,
		intervalKafka: intervalKafka,
	}
}

func (m *Monitor) WithRedis(p RedisPinger) *Monitor  { m.redisPinger = p; return m }
func (m *Monitor) WithKafka(k KafkaChecker) *Monitor { m.kafkaChecker = k; return m }

func (m *Monitor) Start(ctx context.Context) {
	// gRPC targets
	if m.intervalGRPC > 0 && m.grpcPingFn != nil && len(m.services) > 0 {
		go m.loop(ctx, m.intervalGRPC, func(c context.Context) {
			for _, s := range m.services {
				err := m.grpcPingFn(c, s)
				m.setService(s, err)
			}
		})
	}
	// Redis
	if m.intervalRedis > 0 && m.redisPinger != nil {
		go m.loop(ctx, m.intervalRedis, func(c context.Context) {
			err := m.redisPinger.Ping(c)
			m.setRedis(err)
		})
	}
	// Kafka
	if m.intervalKafka > 0 && m.kafkaChecker != nil {
		go m.loop(ctx, m.intervalKafka, func(c context.Context) {
			err := m.kafkaChecker.Check(c)
			m.setKafka(err)
		})
	}
}

func (m *Monitor) loop(ctx context.Context, d time.Duration, fn func(ctx context.Context)) {
	t := time.NewTicker(d)
	defer t.Stop()
	// sofortiger erster Lauf
	fn(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			fn(ctx)
		}
	}
}

func (m *Monitor) setService(name string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.service[name] = Status{OK: err == nil, Error: errString(err), CheckedAt: time.Now()}
}
func (m *Monitor) setRedis(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.redis = Status{OK: err == nil, Error: errString(err), CheckedAt: time.Now()}
}
func (m *Monitor) setKafka(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.kafka = Status{OK: err == nil, Error: errString(err), CheckedAt: time.Now()}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

type Snapshot struct {
	Services map[string]Status `json:"services"`
	Redis    *Status           `json:"redis,omitempty"`
	Kafka    *Status           `json:"kafka,omitempty"`
}

func (m *Monitor) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cp := make(map[string]Status, len(m.service))
	for k, v := range m.service {
		cp[k] = v
	}
	s := Snapshot{Services: cp}
	if m.redis.CheckedAt.After(time.Time{}) {
		s.Redis = &m.redis
	}
	if m.kafka.CheckedAt.After(time.Time{}) {
		s.Kafka = &m.kafka
	}
	return s
}

func (m *Monitor) AllOK(requiredServices []string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, s := range requiredServices {
		st, ok := m.service[s]
		if !ok || !st.OK {
			return false
		}
	}
	// Redis/Kafka nur, wenn konfiguriert
	if m.redisPinger != nil && !m.redis.OK {
		return false
	}
	if m.kafkaChecker != nil && !m.kafka.OK {
		return false
	}
	return true
}
