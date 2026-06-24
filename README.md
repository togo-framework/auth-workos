<!-- togo-brand -->
<p align="center">
  <img src=".github/assets/togo-mark.svg" width="96" alt="togo" />
</p>
<h1 align="center">auth-workos</h1>
<p align="center"><sub>part of the <a href="https://github.com/togo-framework">togo-framework</a> — the full-stack Go + React framework</sub></p>

[WorkOS](https://workos.com) SSO/SAML driver for togo [auth](https://github.com/togo-framework/auth).
Enterprise single sign-on; on success the user is find-or-created and issued a togo session.

```bash
togo install togo-framework/auth
togo install togo-framework/auth-workos
```

Env: `APP_URL`, `WORKOS_CLIENT_ID`, `WORKOS_API_KEY`, and one of
`WORKOS_CONNECTION` / `WORKOS_ORGANIZATION` / `WORKOS_PROVIDER`.
Flow: GET `/api/auth/workos` -> callback -> session.


---

## 💎 Premium sponsors

togo is proudly sponsored by **ID8 Media** and **One Studio**.

<p align="center">
  <a href="https://id8media.com"><img src=".github/assets/id8media.svg" height="44" alt="ID8 Media" /></a>
  &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
  <a href="https://one-studio.co"><img src=".github/assets/one-studio.jpeg" height="44" alt="One Studio" /></a>
</p>
