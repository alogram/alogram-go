# Changelog

All notable changes to the Alogram PayRisk Go SDK will be documented in this file.

## [0.1.6-rc.3] - 2026-02-10

### Added
- Standardized "Smart" client architecture with hand-written façade.
- Resilient retry logic (429 & 5xx) with exponential backoff and jitter.
- Context-aware methods for all risk operations.
- Native OpenTelemetry support for spans and attributes.

### Changed
- Optimized `net/http` transport configuration for production keep-alives.
- Synchronized with Payments Risk API v0.1.6-rc.3.

## [0.1.6-rc.1] - 2026-02-10

### Added
- Built-in retry loop for transient failures (Rate Limits, Server Errors).
- OpenTelemetry integration for tracing risk decisions.

### Changed
- Synchronized with Payments Risk API v0.1.6.