package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	HTTP          *HTTPMetrics
	LLM           *LLMMetrics
	Scoring       *ScoringMetrics
}

func New(service string) *Metrics {
	return &Metrics{
		HTTP:    NewHTTPMetrics(service),
		LLM:     NewLLMMetrics(service),
		Scoring: NewScoringMetrics(service),
	}
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.Handler()
}

type HTTPMetrics struct {
	requestsTotal      *prometheus.CounterVec
	requestDuration    *prometheus.HistogramVec
	requestsInflight   prometheus.Gauge
	requestSizeBytes   *prometheus.HistogramVec
	responseSizeBytes  *prometheus.HistogramVec
	errorsTotal        *prometheus.CounterVec
}

func NewHTTPMetrics(service string) *HTTPMetrics {
	return &HTTPMetrics{
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{"method", "path", "status"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{"method", "path", "status"},
		),
		requestsInflight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_inflight",
				Help: "Number of HTTP requests currently in flight",
				ConstLabels: prometheus.Labels{"service": service},
			},
		),
		requestSizeBytes: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "HTTP request body size in bytes",
				Buckets: prometheus.ExponentialBuckets(16, 4, 8),
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{"method", "path"},
		),
		responseSizeBytes: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response body size in bytes",
				Buckets: prometheus.ExponentialBuckets(64, 4, 8),
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{"method", "path", "status"},
		),
		errorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_errors_total",
				Help: "Total number of HTTP errors by type",
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{"method", "path", "status"},
		),
	}
}

func (m *HTTPMetrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.requestsInflight.Inc()
		defer m.requestsInflight.Dec()

		ww := wrapResponseWriter(w)
		next.ServeHTTP(ww, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ww.statusCode)
		path := getRoutePattern(r)

		m.requestsTotal.WithLabelValues(r.Method, path, status).Inc()
		m.requestDuration.WithLabelValues(r.Method, path, status).Observe(duration)
		m.responseSizeBytes.WithLabelValues(r.Method, path, status).Observe(float64(ww.bodySize))

		if ww.statusCode >= 400 {
			m.errorsTotal.WithLabelValues(r.Method, path, status).Inc()
		}
	})
}

func (m *HTTPMetrics) ObserveRequestBody(r *http.Request, size int) {
	path := getRoutePattern(r)
	m.requestSizeBytes.WithLabelValues(r.Method, path).Observe(float64(size))
}

type LLMMetrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewLLMMetrics(service string) *LLMMetrics {
	return &LLMMetrics{
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "llm_requests_total",
				Help: "Total number of LLM inference requests",
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{"status"},
		),
		requestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "llm_request_duration_seconds",
				Help:    "LLM inference duration in seconds",
				Buckets: prometheus.DefBuckets,
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{},
		),
	}
}

func (m *LLMMetrics) Observe(duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	m.requestsTotal.WithLabelValues(status).Inc()
	m.requestDuration.WithLabelValues().Observe(duration.Seconds())
}

type ScoringMetrics struct {
	requestsTotal       *prometheus.CounterVec
	scoreValue          *prometheus.HistogramVec
	promptLengthBytes   *prometheus.HistogramVec
	responseLengthBytes *prometheus.HistogramVec
	wordCount           *prometheus.HistogramVec
}

func NewScoringMetrics(service string) *ScoringMetrics {
	return &ScoringMetrics{
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "score_requests_total",
				Help: "Total number of scoring requests",
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{"model"},
		),
		scoreValue: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "score_value",
				Help:    "Score value distribution (0-100)",
				Buckets: []float64{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{"model"},
		),
		promptLengthBytes: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "score_prompt_length_bytes",
				Help:    "Prompt length in bytes",
				Buckets: prometheus.ExponentialBuckets(16, 4, 8),
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{"model"},
		),
		responseLengthBytes: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "score_response_length_bytes",
				Help:    "Response length in bytes",
				Buckets: prometheus.ExponentialBuckets(16, 4, 8),
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{"model"},
		),
		wordCount: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "score_word_count",
				Help:    "Response word count",
				Buckets: []float64{1, 5, 10, 20, 50, 100, 200, 500},
				ConstLabels: prometheus.Labels{"service": service},
			},
			[]string{"model"},
		),
	}
}

func (m *ScoringMetrics) Observe(model string, score int, promptLen, respLen, words int) {
	lbls := prometheus.Labels{"model": model}
	m.requestsTotal.With(lbls).Inc()
	m.scoreValue.With(lbls).Observe(float64(score))
	m.promptLengthBytes.With(lbls).Observe(float64(promptLen))
	m.responseLengthBytes.With(lbls).Observe(float64(respLen))
	m.wordCount.With(lbls).Observe(float64(words))
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	bodySize   int
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bodySize += n
	return n, err
}

func getRoutePattern(r *http.Request) string {
	routeCtx := chi.RouteContext(r.Context())
	if routeCtx != nil {
		if pattern := routeCtx.RoutePattern(); pattern != "" {
			return pattern
		}
	}
	return r.URL.Path
}
