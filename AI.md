# Alogram PayRisk SDK for Go — AI Integration Guide

> [!IMPORTANT]
> This file is optimized for Large Language Models (LLMs) and AI Agents. It provides a dense mental model and implementation rules to ensure a 100% correct integration.

## 🧠 Mental Model
Alogram PayRisk uses a **"Dual-Trust" Smart Client** architecture.
1. **`AlogramRiskClient`**: For server-side backend logic. High trust. Requires `sk_...` keys.
2. **`AlogramPublicClient`**: For frontend/public signals. Low trust. Requires `pk_...` keys.

## 🛠️ Core Workflows

### 1. Synchronous Purchase Risk Check
**Objective**: Determine if a purchase should be approved, declined, or stepped up.
**Method**: `AlogramRiskClient.CheckRisk(request: CheckRequest)`
**Critical Parameters**: 
- `idempotency_key`: Must be unique per transaction (e.g., `idk_{tx_id}`).
- `trace_id`: Optional but recommended for end-to-end debugging.

### 2. Behavioral Signal Ingestion
**Objective**: Collect non-payment events (login, signup) to train the model.
**Method**: `AlogramPublicClient.IngestSignals(request: SignalsRequest)`

## 🛡️ Error Handling Rules
AI Agents must implement the following logic:
- **`RateLimitError`**: Implement exponential backoff (Handled by SDK, but don't overwhelm).
- **`ValidationError`**: Terminal. Log the `e.body` and do not retry.
- **`ScopedAccessError`**: Occurs if a `pk_` key is used in `AlogramRiskClient`. Instruct the user to check their key prefix.

## 🧩 Common hallucination prevention
- **Field Naming**: Always use `assessmentId` (not `requestId`) when referencing decisions.
- **Decision Values**: Decisions are strictly lowercase: `approve`, `decline`, `review`, `step_up`.
- **Amounts**: Always pass `amount` as a float.

## 🧪 Testing
- **Local Mocking**: Use `github.com/alogram/alogram-go.MockRiskClient` for unit tests.
- **Integration**: Point `base_url` to `http://localhost:8080` if using the Alogram Local Emulator.
