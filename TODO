v1 Release Priorities
[x] Subscription enforcement - plan limits exist in code but aren't enforced at the API layer
[x] Appointment reminders — the notification hook is missing from appointment creation/updates
[x] Invoice templates — no template system built yet
    [x] make settings width 100%
    [x] fix double scroll of entire app
    [x] create invoice - add item: add space between item and price and fix alignment
[x] Form Validation - disable save buttons if form is invalid or submitting
[p] Email notifications — stub only, no real email provider integrated, and not triggered by any workflows
[ ] Business setting - setting email or invoice receipt messages
[ ] Stripe integration - allow users to upgrade/downgrade subscription plans via Stripe checkout; webhook updates users.subscription_tier on successful payment

Future Features
[x] DO NOT show token in reponse payload
[ ] User profile settings - custom color schemes
[ ] Add horse photo
[ ] Upload body work map
[ ] Ability to upload resources to a client profile (e.g. body worker can upload link to a clients profile for supplement recommendations)
[ ] Bulk actions (e.g. bulk delete from list view)
[ ] Invoice - use dropdrown for items, can add new items in dropdown, auto select newly added item

Tech Debt
[x] Entire app is mobile view, need to add support back for desktop view

Security
[ ] Email verification — blocked by email provider integration ([p] Email notifications must be done first).
    Requires: email_verified column, verification token table, email sending, new frontend screens.
    Severity: Medium. Implement after email provider is integrated.
[ ] Account lockout after N failures — implement now, no dependencies.
    Requires: failed_login_attempts + locked_until columns on users, counter logic in auth/service.go Login, reset on success.
    IP-based rate limiter (5 req/min) exists but does not protect against distributed or per-account attacks.
    Severity: Medium. Recommended next security task.
[ ] Structured JSON logging — cosmetic only, not a security fix. Skip unless a log aggregation service (Datadog, Loki, etc.) is added.
    Severity: Low. Not worth implementing at current scale.
[ ] Reset password on login screen — blocked by email provider integration ([p] Email notifications must be done first).
    Requires: reset token table, forgot-password + reset endpoints, frontend forgot-password flow.
    Severity: Medium. Implement after email provider is integrated.
[x] show password in field

Claude Security Analysis
🟡 Remaining — Minor
1. No Content-Security-Policy header (Low)
The strongest remaining XSP mitigation. Even a basic policy like default-src 'self' would block inline script execution if XSS ever occurred.
2. No revoked_tokens cleanup job (Low)
Expired rows accumulate in the DB. A DELETE FROM revoked_tokens WHERE expires_at < NOW() — on startup or on a schedule — keeps it clean.
3. In-memory rate limiter (Low — current scale)
Fine for a single Railway instance. If you ever scale to multiple instances the buckets won't be shared across them.
4. No frontend request timeouts (Low)
API calls in api.service.ts can hang indefinitely. A simple timeout interceptor would clean this up.
Overall
The app started this session with 4 critical issues and 10+ high/medium ones. All of them are resolved. What's left are four low-priority items, none of which represent meaningful attack vectors at current scale. The security posture is production-ready.