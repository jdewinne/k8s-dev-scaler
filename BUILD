load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/jdewinne/k8s-dev-scaler
gazelle(name = "gazelle")

go_library(
    name = "k8s-dev-scaler_lib",
    srcs = ["main.go"],
    importpath = "github.com/jdewinne/k8s-dev-scaler",
    visibility = ["//visibility:private"],
    deps = [
        "//scaler",
        "@io_k8s_client_go//kubernetes",
        "@io_k8s_client_go//rest",
        "@io_k8s_client_go//tools/clientcmd",
    ],
)

go_binary(
    name = "k8s-dev-scaler",
    embed = [":k8s-dev-scaler_lib"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "k8s-dev-scaler_test",
    srcs = ["main_test.go"],
    embed = [":k8s-dev-scaler_lib"],
)
