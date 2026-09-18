
<!--

*********************************************************************
*               PLEASE READ BEFORE SUBMITTING YOUR PR               *
*     YOUR PR MAY BE REJECTED IF IT DOES NOT FOLLOW THESE STEPS     *
*********************************************************************

- *DO NOT* submit bugs for a source install of v3, ONLY tagged versions, e.g. v3.0.0-beta.1
- For a new capability or change to public behaviour, submit a draft **WEP (Wails Enhancement Proposal)** PR first. Do not start implementation until it is accepted.
  The feedback guide for v3 is here: https://v3.wails.io/feedback/

- For a bug fix, link the bug issue where one exists. For an enhancement implementation, link the accepted WEP PR.
- Please fill in the checklists.

-->

# Description

Please include a summary of the change and which issue is fixed. Please also include relevant motivation and context. List any dependencies that are required for this change.

Fixes # (issue)

## Type of change
  
Please select the option that is relevant.

- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] WEP (proposal only; no implementation)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] This change requires a documentation update

# How Has This Been Tested?
  
Please describe the tests that you ran to verify your changes. Provide instructions so we can reproduce. Please also list any relevant details for your test configuration using `wails doctor`.

- [ ] Windows
- [ ] macOS
- [ ] Linux
      
If you checked Linux, please specify the distro and version.
  
## Test Configuration

Please paste the output of `wails doctor`. If you are unable to run this command, please describe your environment in as much detail as possible.

# Checklist:

- [ ] (v2 only) I have updated `website/src/pages/changelog.mdx` with details of this PR (v3 changelog entries are added automatically)
- [ ] My code follows the general coding style of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
