1. What did you ask the AI to do, and what did you write or decide yourself?
AI was helpful in code generation for almost the entire module, but I had to drive it to be specific in many scenarios.
For example,
(a) Code structure, AI doesnt stick to a constant code structure (can be solved via claude rules but still need to be vetted)
(b) Error handling (AI follows different error handling conventions at different places)


2. Where did you override, correct, or throw away the AI’s output — and why?
(a) Overrided domain segregation of services as Im more aware of the services may be split later into microservices than AI does
(b) Had to correct AI many times explaining error handling, adding missing logs, etc..
(c) AI created its own random unique ID, steered it to ULID or UUID on a case to case basis for me to explain different scenarios


3. The two or three biggest trade-offs you made, and the alternatives you considered.
(a) Stuck only to backend and used integration tests alone without any unit tests as its faster to test E2E flows. I would include more robust unit tests as an alternative.
(b) Used in-memory datastore with partially vetted concurrency logic, in production database handles the concurrency problem for me (alternative in live systems)

4. What’s missing, or what you’d do with another day?
(a) More unit tests
(b) Current architecture has only backend services, I can include BFF (backend for frontend) for much better orchestration. For example, ideally we need a users service, instead I combined users in auth for now as I only needed to store email and password as user details. If I have to split user details and user credentials, BFF and new user service will have to be created