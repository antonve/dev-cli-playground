"""Small reproducible, registry-neutral Go OCI image macro."""

load("@rules_go//go:def.bzl", "go_cross_binary")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_push")
load("@rules_pkg//pkg:mappings.bzl", "pkg_attributes", "pkg_files")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

def go_image(name, binary, binary_name):
    app = "_" + name + "_binary"
    go_cross_binary(name = app, platform = "@rules_go//go/toolchain:linux_amd64", target = binary)
    pkg_files(name = "_" + name + "_files", srcs = [":" + app], attributes = pkg_attributes(mode = "0555"), prefix = "/app", renames = {":" + app: binary_name})
    pkg_tar(name = "_" + name + "_layer", srcs = [":" + "_" + name + "_files"])
    oci_image(name = name, base = "@alpine_3_22", entrypoint = ["/app/" + binary_name], tars = [":" + "_" + name + "_layer"], user = "65532:65532")
    oci_push(name = name + "_push", image = ":" + name)
