# OpenID Connect (OIDC) single sign-on

Stash can delegate authentication to an external OpenID Connect provider (e.g.
Authelia, Authentik, Keycloak, Google, Okta, Auth0). When OIDC is configured a
**Login with SSO** button is shown on the login page.

OIDC can be used on its own or alongside the built-in username/password login:

- **OIDC only** – configure OIDC and leave the username/password blank. The
  password form is hidden and only SSO login is offered.
- **OIDC + password** – configure both. Users may sign in with either method.

> When OIDC is configured it counts as configured authentication, so the same
> protections that apply to password authentication (for example, blocking
> access from public IP addresses when no authentication is configured) apply.

## Configuration

OIDC is configured through the stash `config.yml` file under the `oidc` key, or
via environment variables. At minimum an **issuer** and a **client ID** are
required.

Register a client with your provider and set its redirect URI to:

```
https://<your-stash-host>/oidc/callback
```

### `config.yml`

```yaml
oidc:
  # required
  client_id: stash
  issuer: https://auth.example.com

  # optional
  client_secret: super-secret          # omit for public clients (PKCE is always used)
  redirect_url: https://stash.example.com/oidc/callback  # derived from the request if omitted
  scopes:                              # "openid" is always requested; defaults add profile + email
    - profile
    - email
    - groups
  username_claim: preferred_username   # ID token claim to use as the username (default: sub)

  # restrict access to members of specific groups (optional)
  groups_claim: groups                 # claim containing the user's groups (default: groups)
  allowed_groups:
    - stash-users
    - admins
```

### Environment variables

Scalar settings may also be provided via environment variables (useful for
Docker). List-valued settings (`scopes`, `allowed_groups`) must be set in
`config.yml`.

| Variable                   | Config key            |
| -------------------------- | --------------------- |
| `STASH_OIDC_CLIENT_ID`     | `oidc.client_id`      |
| `STASH_OIDC_CLIENT_SECRET` | `oidc.client_secret`  |
| `STASH_OIDC_ISSUER`        | `oidc.issuer`         |
| `STASH_OIDC_REDIRECT_URL`  | `oidc.redirect_url`   |
| `STASH_OIDC_USERNAME_CLAIM`| `oidc.username_claim` |
| `STASH_OIDC_GROUPS_CLAIM`  | `oidc.groups_claim`   |

## How it works

1. The user clicks **Login with SSO** and is redirected to the provider's
   authorization endpoint. A random `state` (CSRF), `nonce` and a PKCE code
   verifier are stored in a short-lived, signed cookie.
2. The provider authenticates the user and redirects back to `/oidc/callback`.
3. Stash validates the `state`, exchanges the authorization code for tokens
   (using PKCE), verifies the ID token signature and `nonce`, and reads the
   username (and, if configured, groups) from the ID token claims.
4. If an `allowed_groups` list is configured, the user must belong to at least
   one of the listed groups. Otherwise any user the provider authenticates is
   allowed.
5. A normal stash session cookie is issued and the user is redirected to the
   page they originally requested.

## Notes

- Stash uses a single-user model. When password credentials are also
  configured, the OIDC session is bound to the configured username so that
  features tied to that identity (such as signed media URLs) keep working.
- The redirect URI registered with the provider must exactly match the URL
  stash uses. When running behind a reverse proxy, either set `redirect_url`
  explicitly or ensure the proxy forwards `X-Forwarded-Proto`,
  `X-Forwarded-Prefix` and the `Host` header correctly (or set `external_host`).
