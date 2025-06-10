# Canonical's Sphinx Starter Pack

*A pre-configured repository to build and publish documentation with Sphinx.*

## Description

The Documentation starter pack includes:

* A bundled [Sphinx] theme, configuration, and extensions
* Support for both reStructuredText (reST) and MyST Markdown
* Build checks for links, spelling, and inclusive language
* Customisation support layered over a core configuration

See the full documentation: https://canonical-starter-pack.readthedocs-hosted.com/

## Structure

This section outlines the structure of this repository, and some key files.

### `docs/`

This directory contains the documentation for the starter pack itself.

To view it in your browser, navigate to this directory and type `make run`.

### `.github/workflows/`

This directory contains files used for documentation build checks via GitHub's CI.

The file `test-starter-pack.yml` tests the functionality of the starter pack project.

<!--Links-->

[Sphinx]: https://www.sphinx-doc.org/
