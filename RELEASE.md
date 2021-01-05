Release process
===============

Prerequisites
-------------

Make sure you have `hub` installed.

```
brew install hub
```

Releasing
---------

To release new versions of ThriftRW Go, follow these instructions.

1.  Set up some environment variables for use later.

        # The version being released.
        VERSION=1.2.3

        # This is the branch from which $VERSION will be released.
        # This is almost always dev.
        BRANCH=dev

2.  Set up a release branch. We will propose the release from this branch.

        git fetch origin master
        git checkout -b $(whoami)/release origin/master

3.  Merge changes to be released into this branch.

        git merge $BRANCH

4.  Verify that there are no changes to the generated code.

        make generate
        git diff

    If the diff is non-empty, abort the release and figure out why the
    generated code has changed.

5.  Alter the Unreleased entry in CHANGELOG.md to point to `$VERSION` and
    update the link at the bottom of the file. Use the format `YYYY-MM-DD` for
    the year.

    ```diff
    -## [Unreleased]
    +## [1.2.3] - 2020-01-03
    ```

    ```diff
    -[Unreleased]: https://github.com/thriftrw/thriftrw-go/compare/v1.2.2...HEAD
    +[1.2.3]: https://github.com/thriftrw/thriftrw-go/compare/v1.2.2...v1.2.3
    ```

6.  Update the version number in version/version.go and verify that it matches
    what is in the changelog.

        sed -i '' -e "s/^const Version =.*/const Version = \"$VERSION\"/" version/version.go
        make verifyversion

7.  Create a commit for the release.

        git add version/version.go CHANGELOG.md
        git commit -m "Preparing release v$VERSION"

8.  Make a pull request with these changes against `master`.

        hub pull-request -b master --push

9.  Land the pull request after approval as a **merge commit**. To do this,
    select **Create a merge commit** from the pull-down next to the merge
    button and click **Merge pull request**. Make sure you delete that branch
    after it has been merged with **Delete Branch**.

10. Once the change has been landed, pull it locally.

        git checkout master
        git pull

11. Tag a release.

        hub release create -o -m v$VERSION -t master v$VERSION

12. Copy the changelog entries for this release into the release description
    in the newly opened browser window.

13. Switch back to development.

        git checkout $BRANCH
        git merge master
        git checkout -b back-to-dev

14. Add a placeholder for the next version to CHANGELOG.md and a new link at
    the bottom.

    ```diff
    +## [Unreleased]
    +- No changes yet.
    +
     ## [1.2.3] - 2020-01-03
    ```

    ```diff
    +[Unreleased]: https://github.com/thriftrw/thriftrw-go/compare/v1.2.3...HEAD
     [1.2.3]: https://github.com/thriftrw/thriftrw-go/compare/v1.2.2...v1.2.3
    ```

15. Update the version number in version/version.go to the next minor version.

    ```diff
    -const Version = "1.2.3"
    +const Version = "1.3.0"
    ```

16. Verify the version number matches.

        make verifyversion

17. Update the generated code.

        make generate

18. Open a PR with your changes against `dev` to back to development.

        git commit -am "Back to development"
        hub pull-request -b dev --push
