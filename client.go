// Copyright (c) 2025 Alogram Inc.
// All rights reserved.

package alogram

import (
	"bytes"
	"context"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alogram/payrisk-go/internal/payrisk_v1"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ClientOptions struct {
	BaseURL     string
	APIKey      string
	AccessToken string
	TenantID    string
	ClientID    string
	Debug       bool
}

type baseClient struct {
	cfg    *payrisk_v1.Configuration
	api    *payrisk_v1.APIClient
	opts   ClientOptions
	tracer trace.Tracer
}

func newBaseClient(opts ClientOptions) baseClient {
	cfg := payrisk_v1.NewConfiguration()
	cfg.Debug = opts.Debug

	if opts.BaseURL != "" {
		u, err := url.Parse(opts.BaseURL)
		if err == nil {
			cfg.Host = u.Host
			cfg.Scheme = u.Scheme
			cfg.Servers = payrisk_v1.ServerConfigurations{{URL: u.Path}}
		}
	}

	if opts.APIKey != "" {
		cfg.AddDefaultHeader("x-api-key", opts.APIKey)
	}
	if opts.AccessToken != "" {
		cfg.AddDefaultHeader("Authorization", "Bearer "+opts.AccessToken)
	}
	if opts.TenantID != "" {
		cfg.AddDefaultHeader("x-trusted-tenant-id", opts.TenantID)
	}
	if opts.ClientID != "" {
		cfg.AddDefaultHeader("x-trusted-client-id", opts.ClientID)
	}

	return baseClient{
		cfg:    cfg,
		api:    payrisk_v1.NewAPIClient(cfg),
		opts:   opts,
		tracer: otel.Tracer("alogram.payrisk"),
	}
}

func (c *baseClient) generateID() string {
	return uuid.New().String()
}

func (c *baseClient) mapError(err error, resp *http.Response) error {
	if err == nil {
		return nil
	}
	status := 0
	body := ""
	if resp != nil {
		status = resp.StatusCode
		if resp.Body != nil {
			b, _ := io.ReadAll(resp.Body)
			body = string(b)
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
		}
	}
	return NewAlogramError(err.Error(), status, body)
}

func (c *baseClient) isRetryable(err error) bool {
	if _, ok := err.(*RateLimitError); ok {
		return true
	}
	if _, ok := err.(*InternalServerError); ok {
		return true
	}
	return false
}

// 🏢 AlogramRiskClient (Secret Client)
type AlogramRiskClient struct {
	baseClient
}

func NewAlogramRiskClient(opts ClientOptions) (*AlogramRiskClient, error) {
	if strings.HasPrefix(opts.APIKey, "pk_") {
		return nil, NewAlogramError("Cannot initialize AlogramRiskClient with a Publishable Key (pk_...). Please use AlogramPublicClient.", 403, "")
	}
	return &AlogramRiskClient{newBaseClient(opts)}, nil
}

// CheckRisk evaluates risk for a purchase or entity.
func (c *AlogramRiskClient) CheckRisk(ctx context.Context, req payrisk_v1.CheckRequest, ik string, tid string) (*payrisk_v1.DecisionResponse, error) {
	if ik == "" {
		ik = c.generateID()
	}
	if tid == "" {
		tid = c.generateID()
	}

	ctx, span := c.tracer.Start(ctx, "alogram.check_risk", trace.WithAttributes(
		attribute.String("alogram.idempotency_key", ik),
		attribute.String("alogram.trace_id", tid),
	))
	defer span.End()

	var result *payrisk_v1.DecisionResponse
	var httpResp *http.Response
	var err error

	for i := 0; i < 3; i++ {
		result, httpResp, err = c.api.PayriskAPI.RiskCheck(ctx).
			XIdempotencyKey(ik).
			XTraceId(tid).
			CheckRequest(req).
			Execute()

		if err == nil {
			span.SetStatus(codes.Ok, "Success")
			if result.Decision != "" {
				span.SetAttributes(attribute.String("alogram.decision", result.Decision))
			}
			return result, nil
		}

		mappedErr := c.mapError(err, httpResp)
		if !c.isRetryable(mappedErr) {
			span.SetStatus(codes.Error, mappedErr.Error())
			return nil, mappedErr
		}

		backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(backoff):
		}
	}

	return nil, c.mapError(err, httpResp)
}

func (c *AlogramRiskClient) IngestSignals(ctx context.Context, req payrisk_v1.SignalsRequest, ik string, tid string) error {
	return c.ingestSignals(ctx, req, ik, tid)
}

func (c *baseClient) ingestSignals(ctx context.Context, req payrisk_v1.SignalsRequest, ik string, tid string) error {
	if ik == "" {
		ik = c.generateID()
	}
	if tid == "" {
		tid = c.generateID()
	}

	ctx, span := c.tracer.Start(ctx, "alogram.ingest_signals", trace.WithAttributes(
		attribute.String("alogram.idempotency_key", ik),
		attribute.String("alogram.trace_id", tid),
	))
	defer span.End()

	httpResp, err := c.api.PayriskAPI.IngestSignals(ctx).
		XIdempotencyKey(ik).
		XTraceId(tid).
		SignalsRequest(req).
		Execute()

	if err != nil {
		mappedErr := c.mapError(err, httpResp)
		span.SetStatus(codes.Error, mappedErr.Error())
		return mappedErr
	}

	span.SetStatus(codes.Ok, "Success")
	return nil
}

func (c *AlogramRiskClient) IngestEvent(ctx context.Context, event payrisk_v1.PaymentEvent, ik string, tid string) error {
	if ik == "" {
		ik = c.generateID()
	}
	if tid == "" {
		tid = c.generateID()
	}

	ctx, span := c.tracer.Start(ctx, "alogram.ingest_event", trace.WithAttributes(
		attribute.String("alogram.idempotency_key", ik),
		attribute.String("alogram.trace_id", tid),
	))
	defer span.End()

	httpResp, err := c.api.PayriskAPI.IngestPaymentEvent(ctx).
		XIdempotencyKey(ik).
		XTraceId(tid).
		PaymentEvent(event).
		Execute()

	if err != nil {
		mappedErr := c.mapError(err, httpResp)
		span.SetStatus(codes.Error, mappedErr.Error())
		return mappedErr
	}

	span.SetStatus(codes.Ok, "Success")
	return nil
}

// 🌐 AlogramPublicClient (Public Client)
type AlogramPublicClient struct {
	baseClient
}

func NewAlogramPublicClient(opts ClientOptions) (*AlogramPublicClient, error) {
	if strings.HasPrefix(opts.APIKey, "sk_") {
		return nil, NewAlogramError("Cannot initialize AlogramPublicClient with a Secret Key (sk_...). Please use AlogramRiskClient.", 403, "")
	}
	return &AlogramPublicClient{newBaseClient(opts)}, nil
}

func (c *AlogramPublicClient) IngestSignals(ctx context.Context, req payrisk_v1.SignalsRequest, ik string, tid string) error {
	return c.ingestSignals(ctx, req, ik, tid)
}
