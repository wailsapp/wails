# Wails Enhancement Proposal (WEP) Process

## Introduction

Welcome to the Wails Enhancement Proposal (WEP) process. This guide outlines the steps for proposing, discussing, and implementing feature enhancements in Wails. The process is divided into two main parts:

1. **Submission of Proposal**: This part involves documenting your idea, submitting it for review, and discussing it with the community to gather feedback and refine the proposal.
2. **Implementation of Proposal**: Once a proposal is accepted, the implementation phase begins. This involves developing the feature, submitting a PR for the implementation, and iterating based on feedback until the feature is merged and documented.

Following this structured approach ensures transparency, community involvement, and efficient enhancement of the Wails project.

**NOTE**: This process is for proposing new functionality. For bug fixes, documentation improvements, and other minor changes, please follow the standard PR process.

## Submission of Proposal

### 1. Idea Initiation

- **Document Your Idea**: 
  - Create a new directory: `v3/wep/proposals/<name of proposal>` with the name of your proposal. 
  - Copy the WEP template located in `v3/wep/WEP_TEMPLATE.md` into `v3/wep/proposals/<name of proposal>/proposal.md`. 
  - Include any additional resources (images, diagrams, etc.) in the proposal directory.
  - Fill in the template with the details of your proposal. Do not remove any sections.

### 2. Submit Proposal

- **Open a DRAFT PR**:
  - Submit a DRAFT Pull Request (PR) for the proposal with the title `[WEP] <title>`.
  - It should only contain the proposal file and any additional resources (images, diagrams, etc.).
  - Add a summary of the proposal in the PR description.

### 3. Community Discussion

- **Share Your Proposal**: Present your proposal to the Wails community. Try to get support for the proposal to increase the chances of acceptance. If you are on the discord server, create a post in the [`#enhancement-proposals`](https://discord.gg/TA8kbQds95) channel.
- **Gather Feedback**: Refine your proposal based on community input. All feedback should be added as comments in the PR.
- **Show Support**: Agreement with the proposal should be indicated by adding a thumbs-up emoji to the PR. The more thumbs-up emojis, the more likely the proposal will be accepted.
- **Iterate**: Make changes to the proposal based on feedback.
- **Agree on an Implementor**: To avoid stagnant proposals, we require someone agree to implement it. This could be the proposer.
- **Ready for Review**: Once the proposal is ready for review, change the PR status to `Ready for Review`.

A minimum of 2 weeks should be given for community feedback and discussion.

### 4. Final Decision

- **Decision**: The Wails maintainers will make a final decision on the proposal based on community feedback and the proposal's merits. 
  - If accepted, the proposal will be assigned a WEP number and the PR merged.
  - If rejected, the reasons will be provided in the PR comments.

*NOTE*: If a proposal has not met the required support or has been inactive for more than a month, it may be closed.

## Implementation of Proposal

Once a proposal has been accepted and an implementation plan has been decided, the focus shifts to bringing the feature to life. This phase encompasses the actual development, review, and integration of the new feature. Here are the steps involved in the implementation process:

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
The WEP process ensures structured and collaborative enhancement of Wails. Adhere to this guide to contribute effectively to the project's growth.