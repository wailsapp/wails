# Security Policy

## Supported versions

| Version | Security support |
| --- | --- |
| 2.x | :white_check_mark: |
| 3.x | :white_check_mark: |

## Scope

We welcome reports affecting:

- Wails v2 and v3 source code;
- official Wails CLI binaries and downloads;
- the official `@wailsio/runtime` npm package;
- Wails documentation and websites; and
- Wails-operated CI, release, and publishing infrastructure.

Reports about applications built with Wails, or third-party services that Wails
does not operate, are outside this policy's scope. Please report those directly
to the relevant application or service owner.

## Reporting a vulnerability

Use [GitHub Private Vulnerability Reporting](https://github.com/wailsapp/wails/security/advisories/new)
to report a potential vulnerability. Do not open a public issue or disclose
vulnerability details publicly before we have agreed on a disclosure plan.

Please include, where possible:

- a clear description of the issue and its potential impact;
- affected versions, platforms, or configurations;
- reproducible steps or a minimal proof of concept; and
- any mitigation or fix ideas you have identified.

Do not include real credentials, tokens, private keys, personal data, or other
sensitive production information. Redact such information from examples and
proofs of concept.

## What to expect

These are best-effort targets, not guaranteed service-level commitments:

- We aim to acknowledge a report within 48 hours.
- We aim to provide an initial triage or status update within 7 calendar days.
- For confirmed vulnerabilities, we aim to provide an update at least weekly
  until resolution.
- Critical reports may receive faster handling.

If we confirm an issue, we will work on a fix and coordinate disclosure with
the reporter. Where appropriate, we may publish a GitHub Security Advisory or
request a CVE.

## Good-faith research

We will not pursue legal action against researchers who investigate and report
in-scope vulnerabilities in good faith and follow this policy.

Please avoid:

- disrupting services or degrading availability;
- accessing, modifying, or retaining data beyond what is necessary to
  demonstrate the issue;
- testing against other users or systems without their permission; and
- social engineering, phishing, or other attacks against Wails contributors,
  users, or infrastructure providers.

If you are unsure whether planned testing is within scope, submit a private
report first and describe the proposed approach.

## Disclosure and recognition

We ask reporters to give us a reasonable opportunity to investigate and fix an
issue before public disclosure. We will coordinate publication with the
reporter where possible.

Wails does not operate a bug-bounty programme. With the reporter's explicit
permission, we are happy to acknowledge responsible disclosures publicly after
the issue is resolved.
