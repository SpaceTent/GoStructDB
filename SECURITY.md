# Security Policy

## Supported versions

GoStructDB is pre-1.0. Security fixes are made on `main` and released in the next
`v0.0.x` tag. Please make sure you're on the [latest release](https://github.com/SpaceTent/GoStructDB/releases/latest)
before reporting.

| Version | Supported |
| --- | --- |
| Latest `v0.0.x` release | ✅ |
| Older releases | ❌ |

## Reporting a vulnerability

**Please do not report security vulnerabilities through public GitHub issues, discussions
or pull requests.**

Report them privately using GitHub's
[private vulnerability reporting](https://github.com/SpaceTent/GoStructDB/security/advisories/new)
(**Security** tab → **Report a vulnerability**). If that isn't available, email
**mark@ackroyd.net** with `[GoStructDB security]` in the subject.

Please include:

- A description of the issue and its impact (for example, SQL injection through generated SQL).
- Steps or a minimal Go program to reproduce it, including struct tags and table schema.
- The GSDB version, database and Go version.

We'll acknowledge your report as soon as we can, keep you updated while we work on a fix,
and credit you in the release notes unless you'd rather stay anonymous.

## Scope notes

`DB.Insert`, `DB.Update` and `DB.Save` build SQL text from struct values (strings are
hex-encoded). Reports showing that a value can break out of that encoding are in scope.
SQL you write yourself and pass to `Execute`/`Query` is your responsibility; use `?`
parameters for untrusted input.
