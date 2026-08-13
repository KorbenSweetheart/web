# Design Principles

## SOLID

### SRP (Single Responsibility Principle)
Core Meaning: A module, class, or service should have one, and only one, reason to change (it should be responsible to only one actor).

Guidelines: One module/package = one responsibility; one microservice = one bounded domain context.

### OCP (Open/Closed Principle)
Core Meaning: Software entities should be open for extension, but closed for modification.

Guidelines: Introduce new behavior via extension rather than rewriting existing code; introduce abstractions only when actual variability exists.

### LSP (Liskov Substitution Principle)
Core Meaning: Subtypes must be substitutable for their base types without altering the correctness of the program.

Guidelines: Implementations must strictly adhere to the defined contract, preserving behavior, semantic invariants, error signatures, and preconditions/postconditions.

### ISP (Interface Segregation Principle)
Core Meaning: No client should be forced to depend on methods it does not use.

Guidelines: Keep interfaces narrow and consumer-oriented; prefer multiple focused interfaces over a single fat interface.

### DIP (Dependency Inversion Principle)
Core Meaning: High-level modules should not depend on low-level modules; both should depend on abstractions. Abstractions should not depend on details; details should depend on abstractions.

Guidelines: Core business logic depends on abstractions; infrastructure implements these abstractions and is wired via the composition root.

## DRY (Don't Repeat Yourself)
Core Meaning: Every piece of knowledge must have a single, unambiguous, authoritative representation within a system.

Guidelines:
- One single source of truth for repeating business logic.
- Place shared code in libs/; manage shared transport contracts (gRPC) in proto/ instead of duplicating across services.
- Violation Signals: Duplicated middleware, clients, or loggers across services; copy-pasted SQL, DTOs, or validation logic; "almost identical" proto schemas for the same entity.

## KISS (Keep It Simple, Stupid)
Core Meaning: Systems work best if they are kept simple rather than made complicated; simplicity should be a key goal in design, and unnecessary complexity should be avoided.

Guidelines:
- Prefer the simplest implementation that correctly solves the current problem.
- Avoid premature optimization and speculative abstraction (YAGNI — You Aren't Gonna Need It).
- Prioritize code readability and explicit flow over clever or overly concise tricks.

## Clean Architecture

### Layer Rules
- Dependency Rule: Dependencies point inward strictly: transport -> application -> domain; infrastructure remains on the outside.
- Isolation: The domain layer must not import transport frameworks, DB drivers/ORMs, message queue clients, or vendor SDKs.
- Ports & Adapters: The domain defines interfaces (ports); infrastructure provides implementations (adapters); the composition root injects dependencies.
- Testability Criterion: Core business logic must be fully testable and executable without real databases or queues by mocking or stubbing interface dependencies.