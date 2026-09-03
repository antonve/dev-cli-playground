"""Registry-neutral OCI image for a pnpm/Vite application."""

load("@rules_oci//oci:defs.bzl", "oci_image", "oci_push")
load("@rules_pkg//pkg:mappings.bzl", "pkg_attributes", "pkg_files", "strip_prefix")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

def frontend_image(name, srcs, bundle):
    pkg_files(
        name = "_" + name + "_files",
        srcs = srcs,
        attributes = pkg_attributes(gid = 1000, mode = "0644", uid = 1000),
        prefix = "/workspace",
        strip_prefix = strip_prefix.from_root(),
    )
    pkg_tar(
        name = "_" + name + "_layer",
        srcs = [":" + "_" + name + "_files"],
    )
    pkg_files(
        name = "_" + name + "_bundle_files",
        srcs = [bundle],
        attributes = pkg_attributes(gid = 1000, mode = "0644", uid = 1000),
        prefix = "/workspace/apps/hello-web",
    )
    pkg_tar(
        name = "_" + name + "_bundle_layer",
        srcs = [":" + "_" + name + "_bundle_files"],
    )
    oci_image(
        name = name,
        base = "@node_22_alpine",
        cmd = ["apps/hello-web/server.mjs"],
        entrypoint = ["node"],
        tars = [
            ":" + "_" + name + "_layer",
            ":" + "_" + name + "_bundle_layer",
        ],
        user = "1000:1000",
        workdir = "/workspace",
    )
    oci_push(name = name + "_push", image = ":" + name)
