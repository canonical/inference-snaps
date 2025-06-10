# sphinx-docs-starter-pack changelog

## 1.1.0

* Adds sitemap support.
* Simplifies vale binary download & install.
* Leaves vale install output in STDOUT to reveal potential problems.
* Improves update logic.

### Changed

* `docs/conf.py` [#389](https://github.com/canonical/sphinx-docs-starter-pack/pull/389)
* `docs/requirements.txt`(https://github.com/canonical/sphinx-docs-starter-pack/pull/389)

## 1.0.1

Fixes an issue with Vale implementation, and adds words to main wordlist.

### Changed

* `docs/Makefile` [852c19b](https://github.com/canonical/sphinx-docs-starter-pack/commit/852c19bf162e4697d7f36b49e8bc36ad71302216)
* `docs/.sphinx/.wordlist.txt` [#367](https://github.com/canonical/sphinx-docs-starter-pack/pull/367)
* `docs/.sphinx/get_vale_conf.py` [#358](https://github.com/canonical/sphinx-docs-starter-pack/pull/358)

## 1.0.0

First versioned release. Adds an update command to better facilitate updates to
starter pack based documentation sets.

### Added

* `CHANGELOG.md` [#357](https://github.com/canonical/sphinx-docs-starter-pack/pull/357)
* `.github/pull_request_template.md` [#357](https://github.com/canonical/sphinx-docs-starter-pack/pull/357)
* `docs/.sphinx/version` [#357](https://github.com/canonical/sphinx-docs-starter-pack/pull/357)
* `docs/.sphinx/update_sp.py` [#357](https://github.com/canonical/sphinx-docs-starter-pack/pull/357)

### Changed

* `.readthedocs.yaml` [#357](https://github.com/canonical/sphinx-docs-starter-pack/pull/357)
* `.github/workflows/sphinx-python-dependency-build-checks.yml` [#357](https://github.com/canonical/sphinx-docs-starter-pack/pull/357)
* `docs/.sphinx/.markdownlint.json` [#357](https://github.com/canonical/sphinx-docs-starter-pack/pull/357)
* `docs/Makefile` [#357](https://github.com/canonical/sphinx-docs-starter-pack/pull/357)
* `docs/conf.py` [#357](https://github.com/canonical/sphinx-docs-starter-pack/pull/357)
* `docs/requirements.txt` [#357](https://github.com/canonical/sphinx-docs-starter-pack/pull/357)

### Removed

* `.wokeignore` [#363](https://github.com/canonical/sphinx-docs-starter-pack/pull/363)
* `docs/.sphinx/_static/project_specific.css` [#357](https://github.com/canonical/sphinx-docs-starter-pack/pull/357)

## pre-version

This version is the initial versioned release, supporting the implementation of
updates.

### Added

* All files

## VERSION

{Summary of features}

### Added

* {File} {[Commit number](https://www.github.com) or [PR](https://www.github.com)}
* {File} {[Commit number](https://www.github.com) or [PR](https://www.github.com)}
* {File} {[Commit number](https://www.github.com) or [PR](https://www.github.com)}

### Changed

* {File} {[Commit number](https://www.github.com) or [PR](https://www.github.com)}
* {File} {[Commit number](https://www.github.com) or [PR](https://www.github.com)}
* {File} {[Commit number](https://www.github.com) or [PR](https://www.github.com)}

### Removed

* {File} {[Commit number](https://www.github.com) or [PR](https://www.github.com)}
* {File} {[Commit number](https://www.github.com) or [PR](https://www.github.com)}
* {File} {[Commit number](https://www.github.com) or [PR](https://www.github.com)}
