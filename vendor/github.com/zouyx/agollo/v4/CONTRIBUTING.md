## Contributing to agollo

Apollo is released under Apache 2.0 license, and follows a very standard Github development process, using Github tracker for issues and merging pull requests into master. If you want to contribute even something trivial please do not hesitate, but follow the guidelines below.

### Sign the Contributor License Agreement

Before we accept a non-trivial patch or pull request we will need you to sign the Contributor License Agreement. Signing the contributor’s agreement does not grant anyone commit rights to the main repository, but it does mean that we can accept your contributions, and you will get an author credit if we do. Active contributors might be asked to join the core team, and given the ability to merge pull requests.

### Code Conventions

Our code style as below

* Make sure all new .go files have a simple comment with at least an `@author` tag identifying you, and preferably at least a paragraph on what the class is for.

* Add yourself as an @author to the .go files that you modify substantially (more than cosmetic changes).

* A few unit tests should be added for a new feature or an important bug fix.

* If no-one else is using your branch, please rebase it against the current master (or other target branch in the main project).

* When writing a commit message please follow these conventions: if you are fixing an existing issue, please add Fixes #XXX at the end of the commit message (where XXX is the issue number).

* Use ```go fmt``` to format your code , For intellij you can use **File Watching** to trigger format , If you use other IDEs, then you may use command(```go fmt ./...```) before commit.
