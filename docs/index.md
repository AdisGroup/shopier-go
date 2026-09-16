---
layout: home

hero:
  name: "Shopier Go SDK"
  text: "Production-ready Go client for the Shopier REST API"
  tagline: "Zero external dependencies. Full API coverage. Go 1.23+ iterators. Automatic retries."
  actions:
    - theme: brand
      text: Get Started
      link: /getting-started
    - theme: alt
      text: API Reference
      link: /resources/orders
    - theme: alt
      text: GitHub
      link: https://github.com/AdisGroup/shopier-go

features:
  - title: Zero Dependencies
    details: Built entirely with Go standard library (net/http, crypto/hmac, log/slog). No third-party vulnerabilities or bloat.
  - title: Resilient by Design
    details: Automated exponential backoff, sliding-window rate limit awareness, and Retry-After header parsing.
  - title: Modern Go Ergonomics
    details: Native Go 1.23+ Range-over-func iterators (client.Orders.All) alongside classic pagination metadata.
  - title: Webhook Verification
    details: Constant-time HS256 HMAC signature verification with built-in http.Handler and multi-framework adapters.
---
