package main

import "net/http"

type AcceptRoundTripper struct {
	RoundTripper http.RoundTripper
	Accept       []string
}

func (a *AcceptRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for _, accept := range a.Accept {
		req.Header.Add("Accept", accept)
	}
	return a.RoundTripper.RoundTrip(req)
}
