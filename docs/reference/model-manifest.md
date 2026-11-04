(model-manifest)=

# Model manifest

Each model is accompanied by a `model.yaml` manifest file identifying its
metadata, disk requirements, snap components, environment variables, and
filesystem layout. The {ref}`engine manifest <engine-manifest>` references models
by name in its `model.default` and `model.options` fields.

The following sections describe the accepted fields and their validation
requirements. Unknown fields fail validation.

## `name`

**Type:** string. **Required:** yes.

A non-empty, unique name for the model. It must match the name of the directory
containing the `model.yaml` file.

## `alias`

**Type:** string. **Required:** no.

An alternative model identifier accepted by CLI commands that resolve models by
name or alias. Resolution is restricted to models listed in the active engine's
`model.options`.

A non-empty alias must be unique among the models listed in any one engine's
`model.options`. Models belonging to different engines may share an alias.

## `description`

**Type:** string. **Required:** yes.

A non-empty description of the model.

## `model-card-url`

**Type:** string. **Required:** no.

The URL of the model card. A non-empty value must be an absolute URL with a
scheme and host.

## `format`

**Type:** string. **Required:** no.

The model file format. Accepted non-empty values are `GGUF`, `OpenVINO_IR`,
`CTranslate2`, and `MediaTek_DLA`. Values are case-sensitive.

## `quantization`

**Type:** string. **Required:** no.

The model's quantization, for example `Q4_K_M`. There is no fixed set of accepted
values.

## `capabilities`

**Type:** list of strings. **Required:** no.

The capabilities declared by the model. Accepted values are `text`, `vision`,
`tools`, `thinking`, `text-embedding`, `transcription`, and
`realtime-transcription`. Values are case-sensitive.

## `disk-size`

**Type:** string. **Required:** yes.

The disk space required by the model, used when determining whether the model
fits the available disk space.

Values without a suffix are in bytes. The suffixes `M` and `G` denote multiples
of 1024 squared and 1024 cubed bytes, respectively, for example `101M` or `4G`.
Fractional values are accepted; the resulting byte count is truncated to an
integer. Values must be finite and non-negative. Scientific notation and other
unit suffixes are not supported.

## `components`

**Type:** list of strings. **Required:** no.

The snap components required by the model. Each component name must be non-empty.
When `snap/snapcraft.yaml` or `snapcraft.yaml` is found at the snap root,
validation also requires each component to be declared in that file.

## `environment`

**Type:** list of strings. **Required:** no.

Environment variable assignments in `NAME=value` form. Names must match
`^[A-Za-z_][A-Za-z0-9_]*$`. Values may be empty.

Assignments are applied in list order, after the
{ref}`runtime manifest <runtime-manifest>` assignments. Later assignments
override earlier values of the same variable. References in values, such as
`$SNAP_COMPONENTS` or `${MODEL_DIR}`, are expanded from the current environment
when each assignment is applied. Unset variables expand to empty strings.

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

## Example YAML serialization of a model manifest

```yaml
# models/smollm2-135m/model.yaml
name: smollm2-135m
alias: small
description: SmolLM2-135M-Instruct with Q4_K_M quantization.
model-card-url: https://huggingface.co/unsloth/SmolLM2-135M-Instruct-GGUF
format: GGUF
quantization: Q4_K_M
capabilities:
  - text
disk-size: 101M

components:
  - smollm2-135m-q4-k-m-gguf

environment:
  - MODEL_DIR=$SNAP_COMPONENTS/smollm2-135m-q4-k-m-gguf
  - MODEL_FILE=$MODEL_DIR/SmolLM2-135M-Instruct-Q4_K_M.gguf

layout:
  /tmp/model.gguf:
    symlink: $MODEL_FILE
```
