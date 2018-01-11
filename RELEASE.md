Release process
===============

This document outlines how to create a release of thriftrw-go

1.  `git checkout master`

2.  `git pull`

3.  `git merge <branch>` where `<branch>` is the branch we want to cut the
    release on (most likely `dev`)

4.  Alter CHANGELOG.md from `[Unreleased]` to `[X.Y.Z] - YYY-MM-DD` and add
    a reference at the bottom of the document in the form of:
    `[X.Y.Z]: https://github.com/thriftrw/thriftrw-go/compare/vU.V.W...vX.Y.Z`,
    where X.Y.Z is the version to release and U.V.W is the prior.

5.  Alter `version/version.go` to have the same version as `version_to_release`

6.  Run `make verifyversion`

7.  Create a commit with the title `Prepare for release <version_to_release>`

8.  Create a git tag for the version using
    `git tag -a v<version_to_release> -m v<version_to_release` (e.g.,
    `git tag -a v1.0.0 -m v1.0.0`)

9.  Push the tag to origin `git push --tags origin v<version_to_release>`

10. `git push origin master`

11. Go to https://github.com/thriftrw/thriftrw-go/tags and edit the release notes of
    the new tag (copy the changelog into the release notes and make the release
    name the version number)

12. `git checkout dev`

13. `git merge master`

14. Update `CHANGELOG.md` and `version/version.go` to have a new
    `[Unreleased]` (`- No changes yet`) block, run `make generate`, and put into a commit
    with title `Back to development`

15. Run `make verifyversion`

16. `git push origin dev`
