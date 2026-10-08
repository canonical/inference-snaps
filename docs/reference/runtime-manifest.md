(runtime-manifest)=

# Runtime manifest

Each runtime is accompanied by a `runtime.yaml` manifest file identifying its
server endpoints, environment variables, filesystem layout, and snap components.
The {ref}`engine manifest <engine-manifest>` references a runtime by name in its
`runtime` field.

The following sections describe the accepted fields and their validation
requirements. Unknown fields fail validation.

## `name`

**Type:** string. **Required:** yes.

A non-empty, unique name for the runtime. It must match the name of the directory
containing the `runtime.yaml` file.

## `servers`

**Type:** mapping. **Required:** yes.

Server endpoint descriptions, keyed by server name. At least one entry is
required. Server names must not contain `/`. The name `openai` identifies the
OpenAI-compatible endpoint used by CLI features that require an OpenAI server.

Each entry is a mapping with the following fields.

### `protocol`

**Type:** string. **Required:** yes.

The server protocol. Accepted values are:

| Value | Protocol and transport |
| --- | --- |
| `http` | HTTP over TCP. |
| `https` | HTTPS over TCP. |
| `http+unix` | HTTP over a Unix socket. |
| `https+unix` | HTTPS over a Unix socket. |
| `ws` | WebSocket over TCP. |
| `wss` | Secure WebSocket over TCP. |
| `ws+unix` | WebSocket over a Unix socket. |
| `wss+unix` | Secure WebSocket over a Unix socket. |

TCP endpoints obtain their host and port from `http.host` and `http.port` for
HTTP protocols, or `ws.host` and `ws.port` for WebSocket protocols. The
`namespace` field determines the configuration key prefix.

Unix socket endpoints use a socket in the shared provider directory. The socket
filename is determined by `namespace`.

### `base-path`

**Type:** string. **Required:** no. **Default:** empty string.

The URL path included in the endpoint address, for example `/v1`. A non-empty
value must start with `/` and must be a URL path, not a full URL. Query strings
and fragments are not accepted.

### `namespace`

**Type:** string. **Required:** no. **Default:** empty string.

The namespace for the server's configuration keys and Unix socket filename.
A non-empty value must match `^[a-z][a-z0-9-]*$`: a lowercase letter followed by
zero or more lowercase letters, digits, or hyphens.

For TCP endpoints, a non-empty namespace prefixes configuration keys with
`<namespace>.`. For example, `namespace: api` selects `api.http.host` and
`api.http.port` for HTTP. An empty namespace uses the keys without a prefix.

For Unix socket endpoints, the filename is `<namespace>.sock`, or `server.sock`
when the namespace is empty.

## `environment`

**Type:** list of strings. **Required:** no.

Environment variable assignments in `NAME=value` form. Names must match
`^[A-Za-z_][A-Za-z0-9_]*$`. Values may be empty.

Assignments are applied in list order, before the
{ref}`model manifest <model-manifest>` assignments. Later assignments override
earlier values of the same variable. References in values, such as `$PATH` or
`${SNAP_COMPONENTS}`, are expanded from the current environment when each
assignment is applied. Unset variables expand to empty strings.

## `layout`

**Type:** mapping. **Required:** no.

Temporary symbolic links used while the engine environment is loaded. Each key
is a non-empty link path, and its value is a mapping with a required, non-empty
string field `symlink` identifying the link target.

Link paths must be within `/tmp`; targets may be elsewhere. Environment variable
references in both paths and targets are expanded after all runtime and model
environment assignments have been applied.

Model layout entries override runtime entries with the same key. Links are
removed when the engine environment is unloaded.

## `components`

**Type:** list of strings. **Required:** no.

The snap components required by the runtime. Each component name must be
non-empty. When `snap/snapcraft.yaml` or `snapcraft.yaml` is found at the snap
root, validation also requires each component to be declared in that file.

## Example YAML serialization of a runtime manifest

```yaml
# runtimes/llamacpp/runtime.yaml
name: llamacpp

servers:
  openai:
    protocol: http
    base-path: /v1

environment:
  - PATH=$PATH:$SNAP_COMPONENTS/llamacpp/bin
  - LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$SNAP_COMPONENTS/llamacpp/lib
  - LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$SNAP_COMPONENTS/llamacpp/usr/lib/$ARCH_TRIPLET

layout:
  /tmp/llamacpp:
    symlink: $SNAP_COMPONENTS/llamacpp

components:
  - llamacpp
```
