# Wails Enhancement Proposal (WEP) Process

## Introduction

Welcome to the Wails Enhancement Proposal (WEP) process. A WEP is the formal
record for a proposed Wails capability or change to public behaviour. It is
submitted as a pull request, not as a feature-request issue.

The process has two parts:

1. **Creating a WEP**: This part involves documenting your idea, submitting the WEP for review, and discussing it with the community to gather feedback and refine it.
2. **Implementing a WEP**: Once a WEP is accepted, the implementation phase begins. This involves developing the feature, submitting an implementation PR, and iterating based on feedback until the feature is merged and documented.

Following this structured approach ensures transparency, community involvement, and efficient enhancement of the Wails project.

Use GitHub Issues for reproducible bugs and documentation problems. Use GitHub
Discussions for questions or an informal first conversation. A Discussion does
not approve an enhancement: a WEP PR is required before maintainers make a
decision on one.

## Do I Need a WEP?

Not every change needs a proposal. As a rule of thumb:

- **WEP required**: new functionality, a new public API or option, a change to
  public behaviour, breaking changes, cross-platform work, or significant
  platform-specific work.
- **No WEP needed**: a confirmed bug fix, documentation-only change, or
  internal refactoring with no externally observable behaviour change. Use the
  standard PR process.

If you are unsure, ask in the [Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas) category on GitHub Discussions or on [Discord](https://discord.gg/JDdSxwjhGf).

## Creating a WEP

### 1. Idea Initiation

- **Gauge Interest (optional)**: Before writing a WEP, you may float the idea in the [Ideas](https://github.com/wailsapp/wails/discussions/categories/ideas) category on GitHub Discussions or on [Discord](https://discord.gg/JDdSxwjhGf). This can catch overlap with existing plans before you invest time in a write-up.
- **Document Your Idea**: 
  - Create a new directory: `v3/wep/proposals/<name of proposal>` with the name of your WEP.
  - Copy the WEP template located in `v3/wep/WEP_TEMPLATE.md` into `v3/wep/proposals/<name of proposal>/proposal.md`.
  - Include any additional resources (images, diagrams, etc.) in the proposal directory.
  - Fill in the template with the details of your proposal. Do not remove any sections.

### 2. Submit a WEP

- **Open a DRAFT PR**:
  - Submit a DRAFT Pull Request (PR) for the WEP with the title `[WEP] <title>`.
  - It should only contain the proposal file and any additional resources (images, diagrams, etc.).
  - Add a summary of the proposal in the PR description.

### 3. Community Discussion

- **Share Your WEP**: The PR is the official place to discuss the WEP. Present it to the Wails community and try to get support for it to increase the chances of acceptance. If you are on the discord server, share it in the [`#enhancement-proposals`](https://discord.gg/TA8kbQds95) channel.
- **Gather Feedback**: Refine your WEP based on community input. All feedback should be added as comments on the PR so the discussion stays with the document.
- **Show Support**: Reactions can show interest, but maintainers decide on a WEP's technical merits, compatibility, maintenance cost, and fit with Wails.
- **Iterate**: Make changes to the WEP based on feedback.
- **Agree on an Implementor**: To avoid stagnant proposals, we require someone agree to implement it. This could be the proposer.
- **Ready for Review**: Once the proposal is ready for review, change the PR status to `Ready for Review`.

A minimum of 2 weeks should be given for community feedback and discussion.

### 4. Final Decision

- **Decision**: The Wails maintainers will make a final decision on the WEP based on community feedback and the WEP's merits. The decision is recorded as a comment on the PR in this format:

  ```
  **Decision**: Accepted / Rejected
  **Decided by**: @maintainer(s)
  **Rationale**: A short explanation of the decision.
  ```

- **If accepted**:
- The WEP is assigned the next WEP number and its directory renamed to `NNNN-<name of proposal>`.
- The WEP's status is updated to `Accepted`, it is added to the [WEP Index](#wep-index), and the PR is merged.
- **If rejected**: The PR is closed and the WEP recorded in the [WEP Index](#wep-index) with a link to the decision, so the outcome is easy to find if the idea comes up again.

*NOTE*: If a proposal has not met the required support or has been inactive for more than a month, it may be closed as `Withdrawn`.

## Implementation of Proposal

Once a WEP has been accepted and an implementation plan has been decided, the focus shifts to bringing the feature to life. This phase encompasses the actual development, review, and integration of the new feature. Here are the steps involved in the implementation process:

### 1. Develop the Feature

- **Follow Standards**: Implement the feature following Wails coding standards.
- **Document the Feature**: Ensure the feature is well-documented during the development process.
- **Submit a PR**: Once implemented, submit a PR for the feature.

### 2. Feedback and Iteration

- **Gather Feedback**: Collect feedback from the community.
- **Iterate**: Make improvements based on feedback.

### 3. Merging

- **Review of PR**: Address any review comments.
- **Merge**: The PR will be merged after satisfactory review.
- **Update the Status**: When the implementation PR merges, update the WEP's status to `Implemented` and add the implementation PR to the [WEP Index](#wep-index).

*NOTE*: Acceptance requires a committed implementor. If implementation has not started within 3 months of acceptance, the WEP is marked as up for grabs in the WEP Index and anyone may take it on.

## WEP Index

| WEP | Title | Status | Proposal | Implementation |
|-----|-------|--------|----------|----------------|
| [0001](proposals/0001-titlebar-buttons/proposal.md) | Customising Window Controls | Implemented | [#3508](https://github.com/wailsapp/wails/pull/3508) | [#3508](https://github.com/wailsapp/wails/pull/3508) |

**Statuses**: `Draft` (under discussion), `Accepted` (approved, awaiting implementation), `Implemented` (shipped), `Rejected` (declined), `Withdrawn` (closed by the author or for inactivity).

The WEP process ensures structured and collaborative enhancement of Wails. Adhere to this guide to contribute effectively to the project's growth.
