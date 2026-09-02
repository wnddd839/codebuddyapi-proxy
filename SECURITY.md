# Security Policy

## Supported versions

Security fixes are applied on the latest `main` and the newest GitHub Release tag.

## Reporting a vulnerability

Please **do not** open a public GitHub issue for security reports.

Email or message the maintainer via GitHub: [@wnddd839](https://github.com/wnddd839)

Include:

- affected version / commit
- reproduction steps
- impact (auth bypass, secret leak, RCE, etc.)

You should get an acknowledgement within a few days.

## Scope notes

This project is primarily a **local / private-network** reverse proxy:

- Prefer binding to `127.0.0.1`
- Keep `CODEBUDDY_PROXY_API_KEY` set when exposing `/v1`
- Do **not** put admin secrets in URL query strings (unsupported)
- Do not commit `.env` or account JSON
