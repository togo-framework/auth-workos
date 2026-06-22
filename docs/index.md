# auth-workos

[WorkOS](https://workos.com) SSO/SAML driver for togo [auth](https://github.com/togo-framework/auth).
Enterprise single sign-on; on success the user is find-or-created and issued a togo session.

```bash
togo install togo-framework/auth
togo install togo-framework/auth-workos
```

Env: `APP_URL`, `WORKOS_CLIENT_ID`, `WORKOS_API_KEY`, and one of
`WORKOS_CONNECTION` / `WORKOS_ORGANIZATION` / `WORKOS_PROVIDER`.
Flow: GET `/api/auth/workos` -> callback -> session.
