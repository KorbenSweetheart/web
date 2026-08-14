# Go Code Commenting & Documentation Guidelines

This guide establishes the required standards for documenting Go codebases. These rules are mandatory for developers and automated agents.

---

## 1. General Principles

* **Language:** All comments, package documentation, and annotations must be written in English.
* **Explain "Why", Not "What":** Comments must clarify intent, constraints, domain rationale, and non-obvious trade-offs. Do not restate what the code already expresses.
* **Necessity:** If the code is completely self-explanatory, omit the comment. If semantic context or invariants are non-obvious, document them.
* **Accuracy:** Outdated comments are strictly prohibited. When refactoring or modifying behavior, update or delete associated comments immediately.
* **No Masking Bad Code:** Never use comments to explain poorly named variables, convoluted functions, or code smells. Refactor the code first.

---

## 2. Package-Level Comments (`doc.go` or Package Headers)

Package comments are mandatory for any package that:
* Exports public APIs or domain logic.
* Implements infrastructure components (storage, transport, external integrations).
* Serves as an entry point or shared library.

### Requirements:
* Place in a dedicated `doc.go` file (preferred for complex packages) or directly above the `package <name>` statement in the primary file.
* Must begin with `// Package <name> ...`.
* Focus on the package's core responsibility, role within the architecture, and high-level boundaries—never on internal implementation details.

Example:
// Package minio provides storage adapters and bucket management
// for interacting with MinIO object storage.
package minio

---

## 3. Exported Types

All exported types (`struct`, `interface`, type definitions, and enum-like constants) must have a doc comment.

### Requirements:
* Must begin with the type identifier.
* Describe the semantic role of the type in the domain or system.
* Avoid low-effort boilerplate (e.g., `// User represents a user`).

Example:
// UserRepository manages the persistence and retrieval of domain User entities.
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
}

---

## 4. Exported Struct Fields

Exported fields require comments if:
* Their domain meaning or boundary constraint is not obvious from the field name.
* The field represents an external contract (HTTP payload, gRPC field mapping, event schema).
* The field has strict format or value restrictions.

Standard self-evident fields (e.g., `ID`, `CreatedAt`, `UpdatedAt`) do not require comments unless custom constraints apply.

Example:
type RateLimitConfig struct {
    // RequestsPerSecond specifies the maximum allowed incoming requests
    // before the transport middleware begins throttling with 429 Too Many Requests.
    RequestsPerSecond int

    // BurstCapacity defines the token bucket size for brief traffic spikes.
    BurstCapacity int
}

---

## 5. Exported Functions and Methods

All exported functions and methods must have a doc comment.

### Requirements:
* Must begin with the function/method name.
* Describe what the function accomplishes from the caller's perspective.
* Explicitly document side effects, panics, concurrency constraints, or non-obvious error conditions.

Example:
// Authenticate validates user credentials and issues a signed JWT token.
// Returns errs.ErrUnauthorized if the credentials do not match or the user is disabled.
func (s *AuthService) Authenticate(ctx context.Context, email, password string) (string, error) {
    // ...
}

---

## 6. Unexported Complex Logic & Non-Linear Flow

Comments are required for unexported code containing:
* Non-trivial algorithms or complex business calculations.
* Workarounds for third-party driver bugs or platform-specific quirks.
* Performance optimizations that sacrifice standard readability.

Example:
// We use a custom bitmask here instead of a map to eliminate allocations
// on the critical hot path of packet evaluation.
if mask&flagActive != 0 {
    // ...
}

---

## 7. Inline Comments

* Permitted only to clarify branching logic, non-obvious state transitions, or external preconditions.
* **Prohibited:** Step-by-step narration of linear code (e.g., `// check if err is not nil`).

---

## 8. Annotations (`TODO` & `FIXME`)

* `TODO`: Used strictly for planned improvements or missing features. Must include an explanation and, where possible, a tracking reference or ticket ID.
* `FIXME`: Used strictly for temporary workarounds or known edge-case defects that require remediation.

Example:
// TODO(PROJ-123): Replace with redis-backed distributed lock once cluster mode is enabled.
// FIXME: Transient network timeouts cause dropped events under high load; add exponential backoff.

---

## 9. Architectural Deviations & Trade-offs

If code intentionally violates established standards, relies on an architectural trade-off, or implements a temporary shortcut, place a comment explaining:
1. Why the deviation is necessary.
2. The trade-offs accepted.
3. Conditions under which it should be refactored or removed.

---

## 10. Prohibited Practices

* **Never** write redundant or cosmetic comments.
* **Never** leave commented-out code blocks in the repository; rely on Git history instead.
* **Never** duplicate comments across implementations without updating their specific context.
* **Never** let comments fall out of sync with code modifications.