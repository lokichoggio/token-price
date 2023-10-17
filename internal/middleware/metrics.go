package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"token-price/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "service"

	defaultMetricsPath = "/metrics"
)

var (
	labels = []string{"status", "code", "method", "host", "path"}

	uptime = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "uptime",
			Help:      "HTTP service uptime.",
		}, nil,
	)

	reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_request_count_total",
			Help:      "Total number of HTTP requests made.",
		}, labels,
	)

	reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latencies in seconds.",
		}, labels,
	)

	reqSizeBytes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_request_size_bytes",
			Help:      "HTTP request sizes in bytes.",
		}, labels,
	)

	respSizeBytes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_response_size_bytes",
			Help:      "HTTP response sizes in bytes.",
		}, labels,
	)
)

// init registers the prometheus metrics
func init() {
	prometheus.MustRegister(uptime, reqCount, reqDuration, reqSizeBytes, respSizeBytes)
	go recordUptime()
}

// recordUptime increases service uptime per second.
func recordUptime() {
	for range time.Tick(time.Second) {
		uptime.WithLabelValues().Inc()
	}
}

// calcRequestSize returns the size of request object.
func calcRequestSize(r *http.Request) float64 {
	size := 0
	if r.URL != nil {
		size = len(r.URL.String())
	}

	size += len(r.Method)
	size += len(r.Proto)

	for name, values := range r.Header {
		size += len(name)
		for _, value := range values {
			size += len(value)
		}
	}
	size += len(r.Host)

	// r.Form and r.MultipartForm are assumed to be included in r.URL.
	if r.ContentLength != -1 {
		size += int(r.ContentLength)
	}
	return float64(size)
}

type RequestLabelMappingFn func(c *gin.Context) string

// PromOpts represents the Prometheus middleware Options.
// It is used for filtering labels by regex.
type PromOpts struct {
	MetricsPath string

	ExcludeRegexStatus     string
	ExcludeRegexEndpoint   string
	ExcludeRegexMethod     string
	EndpointLabelMappingFn RequestLabelMappingFn
}

// NewDefaultOpts return the default ProOpts
func NewDefaultOpts() *PromOpts {
	return &PromOpts{
		MetricsPath: defaultMetricsPath,
	}
}

// checkLabel returns the match result of labels.
// Return true if regex-pattern compiles failed.
func (po *PromOpts) checkLabel(label, pattern string) bool {
	if pattern == "" {
		return true
	}

	matched, err := regexp.MatchString(pattern, label)
	if err != nil {
		return true
	}
	return !matched
}

type MetricsResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w MetricsResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func promMiddlewareIgnoreRequest(c *gin.Context) bool {
	if c.Request.URL.Path == defaultMetricsPath {
		return true
	}
	if strings.HasPrefix(c.Request.URL.Path, "/debug") {
		return true
	}

	return false
}

// PromMiddleware returns a gin.HandlerFunc for exporting some Web metrics
func PromMiddleware(r *gin.Engine, promOpts *PromOpts) gin.HandlerFunc {
	// make sure promOpts is not nil
	if promOpts == nil {
		promOpts = NewDefaultOpts()
	}

	r.GET(promOpts.MetricsPath, PromHandler(promhttp.Handler()))

	return func(c *gin.Context) {
		if promMiddlewareIgnoreRequest(c) {
			return
		}

		start := time.Now()

		crw := &MetricsResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = crw
		c.Next()

		var code int

		if crw.Status() == http.StatusOK {
			resp := &Response{}
			if err := json.Unmarshal(crw.body.Bytes(), resp); err != nil {
				log.Errorf("metrics Unmarshal response error: %s, body: %s", err, crw.body.String())
			} else {
				code = resp.Code
			}
		}

		lvs := []string{
			strconv.Itoa(crw.Status()),
			strconv.Itoa(code),
			c.Request.Method,
			c.Request.Host,
			c.Request.URL.Path,
		}

		// no response content will return -1
		respSize := c.Writer.Size()
		if respSize < 0 {
			respSize = 0
		}
		reqCount.WithLabelValues(lvs...).Inc()
		reqDuration.WithLabelValues(lvs...).Observe(time.Since(start).Seconds())
		reqSizeBytes.WithLabelValues(lvs...).Observe(calcRequestSize(c.Request))
		respSizeBytes.WithLabelValues(lvs...).Observe(float64(respSize))
	}
}

// PromHandler wrappers the standard http.Handler to gin.HandlerFunc
func PromHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
