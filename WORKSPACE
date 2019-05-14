load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

http_archive(
    name = "io_bazel_rules_go",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.18.3/rules_go-0.18.3.tar.gz"],
    sha256 = "86ae934bd4c43b99893fc64be9d9fc684b81461581df7ea8fc291c816f5ee8c5",
    patch_args = ["-p1"],
    patches = [
        "//patches:rulesgo-env-attr.patch",
        "//patches:rulesgo-gogo.patch",  # T2921479
    ],
)

local_repository(
    name = "io_bazel_rules_goo",
    path ="/Users/rhang/gocode/src/github.com/rules_go"
)

gazelle_revision = "a6448000532153e49f7e0f401428e9c42337c6e1"

http_archive(
    name = "bazel_gazelle",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/archive/{}.zip".format(gazelle_revision)],
    strip_prefix = "bazel-gazelle-{}".format(gazelle_revision),
    sha256 = "ccf4e4c00dfa6334ab362b28d287448bf71bea267c462861cd8305c56788752d",
    patch_args = ["-p1"],
    patches = [
        "//patches:gazelle-dep.patch",
        "//patches:gazelle-regex.patch",
    ],
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "aed1c249d4ec8f703edddf35cbe9dfaca0b5f5ea6e4cd9e83e99f3b0d1136c3d",
    strip_prefix = "rules_docker-0.7.0",
    urls = ["https://github.com/bazelbuild/rules_docker/archive/v0.7.0.tar.gz"],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

