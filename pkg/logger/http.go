package logger

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type transport struct {
	transport http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	input := map[string]any{
		"method":  req.Method,
		"host":    req.URL.Scheme + "://" + req.URL.Host,
		"path":    req.URL.Path,
		"query":   SpreadMaps(req.URL.Query()),
		"headers": SpreadMaps(req.Header),
	}
	b, _ := io.ReadAll(req.Body)
	req.Body = io.NopCloser(bytes.NewReader(b))
	input["body"] = Compress(b)
	begin := time.Now()
	resp, err := t.transport.RoundTrip(req)
	if err == nil {
		b, _ = io.ReadAll(resp.Body)
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
		Transport: &transport{transport: http.DefaultTransport},
		Timeout:   timeout,
	}
	return client
}

func NewTransport(tsp http.RoundTripper) http.RoundTripper {
	return &transport{transport: tsp}
}

func SpreadMaps(h map[string][]string) map[string]string {
	res := make(map[string]string, len(h))
	for k, v := range h {
		res[k] = strings.Join(v, "|")
	}
	return res
}

// Compress 超过2048字节返回截断中间的内容
func Compress(b []byte) string {
	if l := len(b); l > 2048 {
		buf := bytes.NewBuffer(nil)
		buf.Write(b[:1000])
		buf.WriteString("***省略{")
		buf.WriteString(strconv.Itoa(l - 2000))
		buf.WriteString("}字符***")
		buf.Write(b[l-1000:])
		return buf.String()
	}
	return string(b)
}
