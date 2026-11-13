(available-snaps)=

# Available snaps

This page contains a list of inference snaps and the available optimizations.
Search for a snap, or select a model family to see its inference snaps.

<!--
This page is generated from snaps.yaml in this directory using Jinja templating.
To add or update a snap, edit snaps.yaml, not this file.
-->

``````{tab-set}
:class: family-picker

{% for family in families %}
`````{tab-item} {{ family.name }}

{% for s in family.snaps %}
({{ s.anchor }})=

````{card}
:class-card: snap-card

```{rubric} {{ s.name }}
:class: snap-name
```

[![{{ s.snap }} snap](https://snapcraft.io/{{ s.snap }}/badge.svg)](https://snapcraft.io/{{ s.snap }})
[![{{ s.snap }} source][gh-badge]]({{ s.source }})

{{ s.description }}

{{ s.hardware_intro }}

| Arch | Optimization | Description |
|--------------|--------------|-------------|
{% for h in s.hardware %}
| {{ h.arch }} | {{ h.optimization }} | {{ h.description }} |
{% endfor %}

{% if s.explore_optimizations | default(false) %}
{{ explore_optimizations }}
{% endif %}
{% if s.wsl | default(false) %}
{{ wsl_support }}
{% endif %}
````

{% endfor %}
`````

{% endfor %}
``````

<!--
## Deprecated snaps

The following inference snaps are no longer maintained:
- [mistral-7b-instruct](https://github.com/canonical/mistral-7b-instruct-snap)
-->

[gh-badge]: https://img.shields.io/badge/GitHub-repo?logo=github&color=%23104C35
