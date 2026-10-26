# Release post drafts

`wails-v3-released.md` is the English source. The locale directories hold its
prepared translations; none of these draft pages are published by M-Press.

When the release is approved, set the actual release date and remove
`draft: true` from the English post and all nine translations. Move the English
source to `../content/blog/wails-v3-released.md`, each translated page to
`../content/<locale>/blog/wails-v3-released.md`, and each corresponding sidecar
from `translation-state/<locale>/blog/` to
`../translation/state/<locale>/blog/`. Preserve matching filenames.

If the English draft changes before publication, update the translations and
validate them against that source before moving the files. Run the documented
translation coverage, strict build and site checks before publishing.
