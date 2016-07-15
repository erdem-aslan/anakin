package main

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

//@todo: define the Audit payload
type Audit struct {
	Date   time.Time
	Method string
	Url    *url.URL
	Source string
}

// AuditHttpHandler implements http.Handler interface
type AuditHttpHandler struct {
	decoratedHandler http.Handler
	dispatcher       io.Writer
}

func (a *AuditHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.generateAudit(r)
	a.decoratedHandler.ServeHTTP(w, r)
}

func (a *AuditHttpHandler) generateAudit(r *http.Request) {
	//@todo: implement audit generation and dispatching to various stores.
}
