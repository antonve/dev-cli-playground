"""Bazel-owned development deployable metadata."""

def _impl(ctx):
    output = ctx.actions.declare_file(ctx.label.name + ".dev.json")
    ctx.actions.write(output, json.encode({
        "binaryPath": ctx.attr.binary_path,
        "buildTarget": ctx.attr.build_target,
        "image": ctx.attr.image,
        "kind": ctx.attr.kind,
        "name": ctx.attr.service_name,
        "port": ctx.attr.port,
        "readinessPath": ctx.attr.readiness_path,
        "sourceRoots": ctx.attr.source_roots,
        "syncPaths": ctx.attr.sync_paths,
    }) + "\n")
    return [DefaultInfo(files = depset([output]))]

dev_deployable = rule(
    implementation = _impl,
    attrs = {
        "binary_path": attr.string(),
        "build_target": attr.string(mandatory = True),
        "deps": attr.label_list(allow_files = True),
        "image": attr.string(mandatory = True),
        "kind": attr.string(mandatory = True, values = ["frontend", "backend"]),
        "port": attr.int(mandatory = True),
        "readiness_path": attr.string(default = "/readyz"),
        "service_name": attr.string(mandatory = True),
        "source_roots": attr.string_list(mandatory = True),
        "sync_paths": attr.string_list(),
    },
)
