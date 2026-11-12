# Applicant Tracking System API

A lightweight, thread-safe Applicant Tracking System REST API built in Go using only the standard library. The API manages recruitment pipelines, handles applicant registrations, enforces stage-based state machines, and provides validated filtering and pagination.

---

## 1. Overview and How to Run

### Overview

This service provides an API for hiring teams to manage candidates throughout the entire recruitment lifecycle—from initial CV screening through HR interview, technical assignment, user interview, job offering, and onboarding.

**Key Features:**

- **Zero Third-Party Dependencies**: Built exclusively with Go's standard library (`net/http`, `sync`, `encoding/json`, etc.).
- **Deterministic State Machine**: Strict stage and status transition matrix preventing invalid recruitment steps.
- **Thread-Safe In-Memory Store**: Synchronized read/write operations with email uniqueness constraints.
- **Robust Input Validation**: Strict JSON body parsing, URL validation, platform-specific social media prefix checking, word-count motivation checks, and timezone-aware date validation (`Asia/Jakarta`).

### Prerequisites

- **Go**: Version 1.22+ (tested with Go 1.27.4)
- **cURL** or any REST client (e.g., Postman, Bruno)

### Running the Application

1. **Clone and navigate to the project directory**:

   ```bash
   cd NadirSeptianHadiansyah_BackEndEngineer_SKBatch2
   ```

2. **Run tests**:

   ```bash
   go test -v ./...
   ```

3. **Start the server**:

   ```bash
   # Default port: 8080
   go run .

   # Or run with a custom port
   PORT=9000 go run .
   ```

4. **Health check**:
   ```bash
   curl http://localhost:8080/health
   # Response: {"status":"ok"}
   ```

---

## 2. Top‑level Structure

```
NadirSeptianHadiansyah_BackEndEngineer_SKBatch2/
├── main.go         # Server entrypoint, HTTP mux routing, and graceful configuration
├── models.go       # Domain models, request DTOs, and response structures
├── rules.go        # Stage, status, position definitions, and transition state machine
├── validate.go     # Field validators (email, URL, social media, motivation, dates)
├── store.go        # Thread-safe in-memory store with mutex locking
├── handlers.go     # HTTP endpoint handlers and query listing/filtering logic
├── errors.go       # Standardized JSON error response writer
├── rules_test.go   # Table-driven unit tests for the transition state machine
├── go.mod          # Go module specification
├── .gitignore      # Git exclusion rules
└── README.md       # Project documentation
```

---

## 3. Endpoints

### 3.1 Health Check

Checks service availability.

- **Route**: `GET /health`
- **Response** (`200 OK`):
  ```json
  {
    "status": "ok"
  }
  ```

---

### 3.2 Create Applicant

Registers a new candidate into the recruitment pipeline. New candidates automatically start at the `cv_screening` stage with `pending` status.

- **Route**: `POST /applicants`
- **Headers**: `Content-Type: application/json`
- **Payload Constraints**:
  - `name`: String, 2 to 100 characters.
  - `email`: Valid email format, maximum 254 characters, unique across all applicants (case-insensitive).
  - `position`: Must be one of the 28 predefined positions (e.g., `backend_engineer`, `frontend_engineer`).
  - `motivation`: String between 20 and 200 words.
  - `cv_url`: Valid `http` or `https` URL (max 2048 characters).
  - `portfolio_pdf_url`: Valid `http` or `https` URL (max 2048 characters).
  - `portfolio_url` _(optional)_: Valid `http` or `https` URL if provided.
  - `social_media` _(optional)_: Key-value map for `linkedin`, `instagram`, `tiktok`, or `threads`. Must match specific platform URL patterns (e.g., `https://linkedin.com/in/{username}`).

- **Sample Request**:

  ```bash
  curl -X POST http://localhost:8080/applicants \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Alex Doe",
      "email": "alex.doe@example.com",
      "position": "backend_engineer",
      "motivation": "I am deeply passionate about distributed systems and building resilient backend architectures. Over the past few years, I have gained solid experience working with Go and high-throughput microservices. I am eager to contribute my backend engineering skills to this mission-driven team.",
      "cv_url": "https://storage.example.com/cv/alex-doe.pdf",
      "portfolio_pdf_url": "https://storage.example.com/portfolio/alex-doe.pdf",
      "portfolio_url": "https://alexdoe.dev",
      "social_media": {
        "linkedin": "https://linkedin.com/in/alexdoe",
        "instagram": "https://instagram.com/alexdoe"
      }
    }'
  ```

- **Sample Response** (`201 Created`):
  ```json
  {
    "id": 1,
    "name": "Alex Doe",
    "email": "alex.doe@example.com",
    "position": "backend_engineer",
    "stage": "cv_screening",
    "status": "pending",
    "motivation": "I am deeply passionate about distributed systems...",
    "cv_url": "https://storage.example.com/cv/alex-doe.pdf",
    "portfolio_pdf_url": "https://storage.example.com/portfolio/alex-doe.pdf",
    "portfolio_url": "https://alexdoe.dev",
    "social_media": {
      "instagram": "https://instagram.com/alexdoe",
      "linkedin": "https://linkedin.com/in/alexdoe"
    },
    "scheduled_at": null,
    "start_date": null,
    "created_at": "2026-11-12T12:00:00Z",
    "updated_at": "2026-11-12T12:00:00Z"
  }
  ```

---

### 3.3 List Applicants

Retrieves a paginated, filterable, and sortable list of applicants.

- **Route**: `GET /applicants`
- **Query Parameters**:
  | Parameter | Type | Default | Description |
  |---|---|---|---|
  | `stage` | string | _(none)_ | Filter by stage (e.g., `cv_screening`, `hr_interview`, `assignment`, `user_interview`, `offering`, `onboarding`) |
  | `status` | string | _(none)_ | Filter by status (e.g., `pending`, `offered`, `scheduled`, `passed`, `rejected`) |
  | `position` | string | _(none)_ | Filter by position (case-insensitive) |
  | `page` | integer | `1` | Page number (`>= 1`) |
  | `limit` | integer | `20` | Items per page (`1` to `100`) |
  | `sort` | string | `created_at` | Field to sort by: `created_at` or `name` |
  | `order` | string | `desc` | Direction: `asc` or `desc` |

- **Sample Request**:

  ```bash
  curl "http://localhost:8080/applicants?stage=cv_screening&status=pending&sort=created_at&order=desc&page=1&limit=10"
  ```

- **Sample Response** (`200 OK`):
  ```json
  {
    "data": [
      {
        "id": 1,
        "name": "Alex Doe",
        "email": "alex.doe@example.com",
        "position": "backend_engineer",
        "stage": "cv_screening",
        "status": "pending",
        "motivation": "I am deeply passionate about distributed systems...",
        "cv_url": "https://storage.example.com/cv/alex-doe.pdf",
        "portfolio_pdf_url": "https://storage.example.com/portfolio/alex-doe.pdf",
        "portfolio_url": "https://alexdoe.dev",
        "social_media": {
          "instagram": "https://instagram.com/alexdoe",
          "linkedin": "https://linkedin.com/in/alexdoe"
        },
        "scheduled_at": null,
        "start_date": null,
        "created_at": "2026-11-12T12:00:00Z",
        "updated_at": "2026-11-12T12:00:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "totalPages": 1
    }
  }
  ```

---

### 3.4 Update Applicant Stage & Status

Transitions an applicant to a new status according to the recruitment state machine. When a candidate reaches an advancing status (`passed` or `accepted`), they automatically progress to the next stage with its initial status.

- **Route**: `PATCH /applicants/{id}/stage`
- **Headers**: `Content-Type: application/json`
- **Payload Parameters**:
  - `status` _(string, required)_: Desired target status.
  - `scheduled_at` _(string, required when status is `scheduled` or `rescheduled`)_: Date in `YYYY-MM-DD` format; must be today or future date in WIB (`Asia/Jakarta`).
  - `start_date` _(string, required when status is `accepted`)_: Date in `YYYY-MM-DD` format; must be today or future date in WIB (`Asia/Jakarta`).

- **Stage Transition Rules**:
  | Current Stage | Current Status | Allowed Requested Statuses | Resulting Stage & Status |
  |---|---|---|---|
  | `cv_screening` | `pending` | `passed`, `rejected` | `passed` advances to `hr_interview` (`offered`); `rejected` stays in stage (terminal). |
  | `hr_interview` | `offered` | `scheduled`, `rejected` | `scheduled` requires `scheduled_at`. |
  | `hr_interview` | `scheduled` | `rescheduled`, `passed`, `rejected` | `rescheduled` requires `scheduled_at`; `passed` advances to `assignment` (`offered`). |
  | `hr_interview` | `rescheduled` | `rescheduled`, `passed`, `rejected` | `rescheduled` requires `scheduled_at`; `passed` advances to `assignment` (`offered`). |
  | `assignment` | `offered` | `passed`, `rejected` | `passed` advances to `user_interview` (`offered`); `rejected` stays in stage (terminal). |
  | `user_interview` | `offered` | `scheduled`, `rejected` | `scheduled` requires `scheduled_at`. |
  | `user_interview` | `scheduled` | `rescheduled`, `passed`, `rejected` | `rescheduled` requires `scheduled_at`; `passed` advances to `offering` (`offered`). |
  | `user_interview` | `rescheduled` | `rescheduled`, `passed`, `rejected` | `rescheduled` requires `scheduled_at`; `passed` advances to `offering` (`offered`). |
  | `offering` | `offered` | `accepted`, `declined` | `accepted` requires `start_date` and advances to `onboarding` (`scheduled`); `declined` is terminal. |
  | `onboarding` | `scheduled` | `completed`, `declined` | Both are terminal states. |

- **Sample Request (Advancing from CV Screening to HR Interview)**:

  ```bash
  curl -X PATCH http://localhost:8080/applicants/1/stage \
    -H "Content-Type: application/json" \
    -d '{
      "status": "passed"
    }'
  ```

- **Sample Response** (`200 OK`):

  ```json
  {
    "id": 1,
    "name": "Alex Doe",
    "email": "alex.doe@example.com",
    "position": "backend_engineer",
    "stage": "hr_interview",
    "status": "offered",
    "motivation": "I am deeply passionate about distributed systems...",
    "cv_url": "https://storage.example.com/cv/alex-doe.pdf",
    "portfolio_pdf_url": "https://storage.example.com/portfolio/alex-doe.pdf",
    "portfolio_url": "https://alexdoe.dev",
    "social_media": {
      "instagram": "https://instagram.com/alexdoe",
      "linkedin": "https://linkedin.com/in/alexdoe"
    },
    "scheduled_at": null,
    "start_date": null,
    "created_at": "2026-11-12T12:00:00Z",
    "updated_at": "2026-11-12T12:05:00Z"
  }
  ```

- **Sample Request (Scheduling an Interview with Date Validation)**:
  ```bash
  curl -X PATCH http://localhost:8080/applicants/1/stage \
    -H "Content-Type: application/json" \
    -d '{
      "status": "scheduled",
      "scheduled_at": "2026-11-20"
    }'
  ```

---

## 4. Error Format and Status Code Table

### Error Response Schema

All error responses return a unified JSON payload:

```json
{
  "error": "Short description of the error",
  "statusCode": 400,
  "details": [
    {
      "field": "field_name",
      "message": "Specific reason for the validation failure"
    }
  ]
}
```

_(Note: If there are no granular field details, the `details` key is omitted)._

### Status Code Table

| Status Code | HTTP Status             | Description / When Returned                                                                                                                                                                               |
| ----------- | ----------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **200**     | `OK`                    | Successful retrieval (`GET /applicants`, `GET /health`) or successful stage update (`PATCH /applicants/{id}/stage`).                                                                                      |
| **201**     | `Created`               | Successfully registered a new applicant (`POST /applicants`).                                                                                                                                             |
| **400**     | `Bad Request`           | Malformed JSON syntax, unexpected JSON fields (`DisallowUnknownFields`), failed field validation, unknown status/stage/position query filters, or missing mandatory dates (`scheduled_at`, `start_date`). |
| **404**     | `Not Found`             | Requested route does not exist (`/*`), or applicant ID does not exist in the store.                                                                                                                       |
| **409**     | `Conflict`              | Attempted to create an applicant with an email that is already registered (case-insensitive check).                                                                                                       |
| **422**     | `Unprocessable Entity`  | Illegal stage transition requested (e.g., attempting to jump from `offered` directly to `passed` in an interview stage, or requesting an invalid status for the current stage).                           |
| **423**     | `Locked`                | Applicant is in a terminal status (`rejected`, `declined`, `completed`) and cannot be transitioned further.                                                                                               |
| **500**     | `Internal Server Error` | Unexpected server or storage failure.                                                                                                                                                                     |

---

## 5. Design Decisions

1. **Standard Library Only (Zero External Dependencies)**
   - Utilizes Go's modern `net/http` ServeMux routing (Go 1.22+ pattern matching for methods and path parameters like `PATCH /applicants/{id}/stage`).
   - Keeps the binary lightweight, fast to compile, free from dependency drift, and easily audit-ready.

2. **Deterministic State Machine with Guardrails**
   - Stage progression is governed by a declarative transition table (`rules.go`).
   - Automatic advancement: Candidates that pass a milestone automatically advance to the next stage with the correct initial status (e.g. `passed` in `cv_screening` automatically advances to `hr_interview` with status `offered`).
   - Stage change cleanups: Resetting `scheduled_at` to `null` upon advancing stages prevents stale interview schedules from carrying over into subsequent stages.
   - Terminal state protection: Candidates marked as `rejected`, `declined`, or `completed` cannot be modified (HTTP 423 Locked).

3. **Concurrency and Thread Safety**
   - `sync.RWMutex` wraps the in-memory store.
   - Read operations (`Get`, `List`) acquire read locks (`RLock`) to permit concurrent reads.
   - Mutations (`Create`, `Transition`) acquire exclusive locks (`Lock`) and execute state checks inside atomic callbacks to avoid race conditions (check-then-act hazards).

4. **Strict JSON and Input Sanitization**
   - Request bodies are wrapped in `http.MaxBytesReader` (1 MB limit) to guard against unbounded memory allocation.
   - `json.Decoder.DisallowUnknownFields()` is enabled to reject payloads containing typos or undocumented fields.
   - Case normalization: Emails, positions, statuses, and stages are trimmed and lowercased to prevent subtle mismatches.

5. **Timezone-Aware Date Validation**
   - Scheduled dates and start dates are strictly validated against `Asia/Jakarta` (WIB) time.
   - Scheduling dates prior to the current calendar day in Jakarta is rejected.

6. **Deterministic Pagination & Stable Sorting**
   - In-memory listing supports pagination (`page`, `limit`) and multi-criteria sorting (`created_at`, `name`).
   - Tie-breaking is enforced using applicant ID to guarantee deterministic ordering across page boundaries.

---

## 6. Limitations

- **In-Memory Volatility**: Data is stored in memory and will be reset when the application process terminates or restarts.
- **Single-Instance Deployment**: Because state is managed in-process, running multiple instances behind a load balancer without sticky sessions would result in disjoint state. A persistent datastore (such as PostgreSQL) would be required for horizontal scaling.
- **Authentication & Authorization**: The API currently lacks authentication (e.g., Bearer tokens, OAuth2, or JWT) and role-based access control (RBAC). In a production environment, HR, recruiters, and interviewers would have scoped permissions.
- **Dataset Size**: In-memory slice filtering and sorting operates smoothly for thousands of records, but for hundreds of thousands of applicants, offloading pagination and filtering to indexed database queries would be necessary.

---

## 7. AI Usage Note

AI tools were utilized during the development and refinement of this project:

- **Refactoring & Code Tidiness**: Consolidating stray files (`position.go`, `list.go`) into their respective domain modules while preserving clean boundaries and eliminating dead code.
- **Test Matrix Validation**: Formulating comprehensive edge-case test tables in `rules_test.go` covering illegal transitions, terminal states, and stage advancements.
- **Documentation**: Assisting in synthesizing API schema specifications, error mappings, and architectural decision records for this README. All generated code and documentation were manually reviewed, validated, and verified through unit tests and automated smoke tests.
