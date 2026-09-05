"""Bazel-owned development deployable metadata."""

def _impl(ctx):
    output = ctx.actions.declare_file(ctx.label.name + ".dev.json")
    ctx.actions.write(output, json.encode({
        "binaryPath": ctx.attr.binary_path,
        "buildTarget": ctx.attr.build_target,
        "containerPath": ctx.attr.container_path,
        "devCommand": ctx.attr.dev_command,
        "imageName": ctx.attr.image_name,
        "imageTarget": ctx.attr.image_target,
        "kind": ctx.attr.kind,
        "name": ctx.attr.service_name,
        "port": ctx.attr.port,
        "pushTarget": ctx.attr.push_target,
        "publicPath": ctx.attr.public_path,
        "internalHost": ctx.attr.internal_host,
        "readinessPath": ctx.attr.readiness_path,
        "sourceRoots": ctx.attr.source_roots,
        "syncPaths": ctx.attr.sync_paths,
        "dependencyPaths": ctx.attr.dependency_paths,
        "dependencyCommand": ctx.attr.dependency_command,
    }) + "\n")
    return [DefaultInfo(files = depset([output]))]

dev_deployable = rule(
    implementation = _impl,
    attrs = {
        "binary_path": attr.string(),
        "build_target": attr.string(mandatory = True),
        "container_path": attr.string(),
        "dev_command": attr.string_list(),
        "dependency_paths": attr.string_list(),
        "dependency_command": attr.string_list(),
        "deps": attr.label_list(allow_files = True),
        "image_name": attr.string(mandatory = True),
        "image_target": attr.string(mandatory = True),
        "kind": attr.string(mandatory = True, values = ["frontend", "backend"]),
        "port": attr.int(mandatory = True),
        "public_path": attr.string(),
        "internal_host": attr.string(),
        "readiness_path": attr.string(default = "/readyz"),
        "push_target": attr.string(mandatory = True),
        "service_name": attr.string(mandatory = True),
        "source_roots": attr.string_list(mandatory = True),
        "sync_paths": attr.string_list(),
    },
)
