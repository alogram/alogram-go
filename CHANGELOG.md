# Changelog

All notable changes to the Alogram PayRisk Go SDK will be documented in this file.

## [0.2.8] - 2026-04-11
- Added support for external assessments and enhanced identity context.
- Introduced detailed risk breakdowns and reason codes in the decision response.
- Standardized fraud scoring with the new `decisionScore` field (`riskScore` is now deprecated).
- Improved connection resilience with built-in automatic retries.

## [0.1.6-rc.3] - 2026-02-10
- Hardened MockRiskClient with thread-safe mock decisions.
- Unified production API defaults across all clients.
- Added comprehensive AI Agent integration instructions.
- Standardized project-local build caches for hermetic testing.

## [0.1.6-rc.2] - 2026-02-10
- Synchronized with Payments Risk API v0.1.6-rc.3.
- Enhanced README with premium documentation sections.

## [0.1.6-rc.3] - 2026-02-09
- Initial release candidate.
