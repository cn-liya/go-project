package logger

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"time"
)

type transport struct {
	transport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	input := map[string]any{
		"method": req.Method,
		"host":   req.URL.Scheme + "://" + req.URL.Host,
		"path":   req.URL.Path,
		"query":  SpreadMaps(req.URL.Query()),
		"header": SpreadMaps(req.Header),
	}
	if req.Body != nil {
		b, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(b))
		input["body"] = Compress(b)
	}
	begin := time.Now()
	resp, err := t.transport.RoundTrip(req)
	if err == nil {
		b, _ := io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewReader(b))
		FromContext(req.Context()).Trace("request", input, map[string]any{
			"body":   Compress(b),
			"status": resp.StatusCode,
		}, begin)
	} else {
		FromContext(req.Context()).Error("request", input, err)
	}
	return resp, err
}

func NewHttpClient(timeout time.Duration) *http.Client {
	client := &http.Client{
		Transport: NewTransport(http.DefaultTransport),
		Timeout:   timeout,
	}
	return client
}

func NewTransport(tsp http.RoundTripper) http.RoundTripper { return &transport{transport: tsp} }

// SpreadMaps 将url.Values或http.Header值的数组展开为字符串
func SpreadMaps(m map[string][]string) map[string]string {
	res := make(map[string]string, len(m))
	for k, v := range m {
		if len(v) > 0 { // 同url.Values和http.Header的Get方法
			res[k] = v[0]
		}
	}
	return res
}

// Compress 超过2048字节返回截断中间的内容
func Compress[T string | []byte](v T) string {
	if l := len(v); l > 2048 {
		cut := strconv.Itoa(l - 2000)
		buf := make([]byte, 0, 2012+len(cut))
		buf = append(buf, v[:1000]...)
		buf = append(buf, "***<(*"...)
		buf = append(buf, cut...)
		buf = append(buf, "*>)***"...)
		buf = append(buf, v[l-1000:]...)
		return string(buf)
	}
	return string(v)
}
