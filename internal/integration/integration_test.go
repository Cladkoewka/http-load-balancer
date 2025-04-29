package integration

import (
	"net/http"
	"testing"
)

func BenchmarkLoadBalancer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resp, err := http.Get("http://localhost:8080/test")
		if err != nil {
			b.Fatal(err)
		}

		resp.Body.Close()
	}
}
