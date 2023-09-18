package flux

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestResource_LogValue(t *testing.T) {
	out := bytes.NewBufferString("")
	opt := slog.HandlerOptions{ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		// Remove time from the output for predictable test output.
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		return a
	}}
	l := slog.New(slog.NewJSONHandler(out, &opt))

	l.Info("resource found", "resource", Resource{
		Kind:      "HelmRelease",
		Namespace: "default",
		Name:      "foo",
		Conditions: map[string]string{
			"ready":    "False",
			"released": "True",
		},
	})

	assert.Equal(t, `{"level":"INFO","msg":"resource found","resource":{"kind":"HelmRelease","namespace":"default","name":"foo","conditions":{"ready":"False","released":"True"}}}
`, out.String())

}

func BenchmarkResource_LogValue(b *testing.B) {
	r := Resource{
		Kind:      "HelmRelease",
		Namespace: "default",
		Name:      "foo",
		Conditions: map[string]string{
			"ready":    "False",
			"released": "True",
		},
	}

	for i := 0; i < b.N; i++ {
		v := r.LogValue()
		if len(v.Group()) != 4 {
			b.Fail()
		}
	}
}
