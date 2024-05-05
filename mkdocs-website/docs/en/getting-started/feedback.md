# Feedback

We welcome (and encourage) your feedback! Please search for existing tickets or
posts before creating new ones. Here are the different ways to provide feedback:

=== "Bugs"

    If you find a bug, please let us know by posting into the [v3 Alpha Feedback](https://discord.gg/Vgff2p8gsy) channel on Discord. 
    
    - The post should clearly state what the bug is and have a simple reproducible example. If the docs are unclear what *should* happen, please include that in the post.
    - The post should be given the `Bug` tag.
    - Please include the output of `wails doctor` in your post.
    - If the bug is behaviour that does not align with current documentation, e.g. a window does not resize properly, please do the following:
      - Update an existing example in the `v3/example` directory or create a new example in the `v3/examples` folder that clearly shows the issue.
      - Open a [PR](https://github.com/wailsapp/wails/pulls) with the title `[v3 alpha test] <description of bug>`.
      - Please include a link to the PR in your post.

    !!! warning
        *Remember*, unexpected behaviour isn't necessarily a bug - it might just not do what you expect it to do. Use [Suggestions](#suggestions) for this.


=== "Fixes"

    If you have a fix for a bug or an update for documentation, please do the following:

    - Open a pull request on the [Wails repository](https://github.com/wailsapp/wails). The title of the PR should start with `[v3 alpha]`.
    - Create a post in the [v3 Alpha Feedback](https://discord.gg/Vgff2p8gsy) channel.
    - The post should be given the `PR` tag.
    - Please include a link to the PR in your post.

=== "Suggestions"

    If you have a suggestion, please let us know by posting into the [v3 Alpha Feedback](https://discord.gg/Vgff2p8gsy) channel on Discord:

    - The post should be given the `Suggestion` tag.

    Please feel free to reach out to us on [Discord](https://discord.gg/Vgff2p8gsy) if you have any questions.

=== "Upvoting"

    - Posts can be "upvoted" by using the :thumbsup: emoji. Please apply to any posts that are a priority for you.
    - Please *don't* just add comments like "+1" or "me too".
    - Please feel free to comment if there is more to add to the post, such as "this bug also affect ARM builds" or "Another option would be to ....."

There is a list of known issues & work in progress can be found
[here](https://github.com/orgs/wailsapp/projects/6).

## Things we are looking for feedback on

- The API
  - Is it easy to use?
  - Does it do what you expect?
  - Is it missing anything?
  - Is there anything that should be removed?
  - Is it consistent between Go and JS?
- The build system
  - Is it easy to use?
  - Can we improve it?
- The examples
  - Are they clear?
  - Do they cover the basics?
- Features
  - What features are missing?
  - What features are not needed?
- Documentation
  - What could be clearer?
