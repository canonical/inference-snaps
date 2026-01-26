# Contributing
We are committed to continuously improving our code and documentation and making it more useful for both ourselves and others.

We welcome community contributions, suggestions, and fixes.
Please review the following guidelines before contributing to the project.

(contribute-to-our-code)=
## Contribute to our code
You can help improve our code by:
- Proposing new features or improvements
- Identifying and fixing bugs
- Reporting bugs and issues by opening an issue on GitHub

### Inference snaps code bases
Our projects' source code can be found at:
- [DeepSeek R1 snap](https://github.com/canonical/deepseek-r1-snap)
- [Qwen VL snap](https://github.com/canonical/qwen-vl-snap)
- [Gemma3 snap](https://github.com/canonical/gemma3-snap)
- [Nemotron 3 Nano snap](https://github.com/canonical/nemotron-3-nano-snap)

### Tooling
Inference snaps use the following reusable building blocks:
- [Inference Snaps CLI](https://github.com/canonical/inference-snaps-cli) – the command-line interface at the heart of each inference snap
- [Inference Snaps Developer Tools](https://github.com/canonical/inference-snaps-dev) – GitHub Actions and utility scripts for building, testing, and publishing inference snaps


(contribute-to-our-docs)=
## Contribute to our docs
You can help improve this documentation by:

- Documenting new features or improvements you contribute to the code.
- Clarifying concepts or common questions based on your own experience.
- Fixing typos, formatting errors, and other mistakes.
- Reporting documentation issues by opening an issue on GitHub.

### License and copyright
Unless explicitly stated in the license headers of source files, all contributions to the Inference Snaps documentation are licensed under the [Creative Commons Attribution-Share Alike (CC-BY-SA) 3.0 Unported License](https://creativecommons.org/licenses/by-sa/3.0/). 
You can also get a copy of this license by sending a letter to Creative Commons, 171 Second Street, Suite 300, San Francisco, California, 94105, USA.

All contributors must sign the [Canonical contributor license agreement](https://canonical.com/legal/contributors), which grants Canonical permission to use the contributions. 
The author of a change remains the copyright owner of their code (no copyright assignment occurs).

### Style and language
This documentation is aimed at helping both beginners and experts. 
We try to use appropriate writing for the subject, use inclusive language, and assume the reader has as little prior knowledge as possible.

The navigation, structure, style, and content of our documentation follows the Diátaxis framework for technical documentation authoring. 
This splits pages into tutorials, how-to guides, reference material, and explanations.

Some general tips to ensure high quality documentation are:

- use a spell checker

- resist being overly formal

- resist being overly verbose

- Use links and examples from reputable sources, e.g., official upstream documentation.

We also recommend using the [Ubuntu style guide].

You are welcome to suggest changes to any documentation topic with updated or more insightful information. 
We aim for consistency, but don’t let formality prevent you from contributing.

### Inference Snaps documentation overview

The Inference Snaps documentation is built on top of [Canonical's Sphinx starter pack] and hosted on GitHub.
It is published on [Read the Docs].

The documentation is organized according to the [Diátaxis](https://diataxis.fr/) framework and written in reStructuredText or MyST Markdown.

- For general guidance, refer to the [starter pack guide]
- For syntax help and guidelines, refer to the [MyST syntax guide]


### Build the docs

The sphinx starterpack allows you to build the Inference Snaps documentation locally.

#### Fork and clone the repository

If you are working with this documentation set for the first time, fork the repository to your GitHub account, and clone the forked repository.

#### Install dependencies

Go into the `docs` directory and run:

```
make install
```

This will create a virtual environment with all the tooling and packages
required to validate and render the documentation.

#### Build and serve the docs

To get a live preview and validation of the documentation, run:

```
make run
```

This will serve a locally rendered version of the documentation at:
<http://localhost:8000> by default.

The rendered preview will get updated every time a document is modified and
saved.

#### Validate the docs

Before you submit a PR, validate your changes with the following recipes:

```
make spelling  
make linkcheck  
make woke       # to check for non-inclusive language
```

If the spell check reports any false positives:
1. Add common terms to the `.custom_wordlist.txt` file to ignore them.
2. Ignore document-specific terms using the `` {spellexception}`special.word` `` markup.
3. Sort the list. You can use `LC_ALL=POSIX sort --stable --ignore-case .custom_wordlist.txt --output=.custom_wordlist.txt`.

For more information about the inclusive language check, see [woke].

### Propose changes

Submit your documentation changes via a pull request on GitHub.

Your changes will be reviewed in due time; if approved, they will eventually be
merged.

#### Describing pull requests

To be properly considered, reviewed, and merged, your pull request must provide
the following details:

- **Title**: Summarize the change in a short, descriptive title.
- **Description**: Explain the problem that your pull request solves. Mention
any new features, bug fixes, or refactoring.
- **Relevant issues**: Reference any [related issues, pull requests, and
repositories].

#### Commit structure and messages

Use separate commits for each logical change, and for changes to different
sections in the Inference Snaps documentation.
Prefix your commit messages with the names of sections or pages that they
affect, using the code tree structure. For example, start a commit that updates
the reference page kernel boot parameters with `ref/kernel-boot:`.

Use [conventional commits] to ensure consistency across the project:

```none
docs(ref/kernel-boot): Add example for nohz parameter

* Additional description goes here
```

Such structures makes it easier to review contributions and simplifies porting
fixes to other branches.

#### Sign off on commits

To improve contribution tracking, we use the developer certificate of origin
([DCO 1.1](https://developercertificate.org/)) and require a "sign-off" for any
changes going into each branch.

The sign-off is a simple line at the end of the commit message certifying that
you wrote it or have the right to commit it as an open-source contribution.

To sign off on a commit, use the `--signoff` or `-s` option in `git commit`.
For example:

```
git commit -m "docs(ref/kernel-boot): Add example for nohz parameter" -s
```
#### Automatic documentation checks
When you submit a pull request, GitHub will automatically run checks on the documentation to verify the spelling, validity of links, correct formatting of Markdown files, and use of inclusive language.

### Open Documentation Academy
The Inference Snaps project is a proud member of the Canonical Open Documentation Academy (CODA).

The Open Documentation Academy is an initiative that aims to encourage open source contributions from the community.
We use documentation as the gateway, thus lowering the barrier to entry.

Through the academy, community members get access to curated and well-defined beginner-friendly tasks along with mentorship and advice in a friendly and encouraging environment. 
This allows them to make meaningful contributions to established open-source projects, such as Launchpad, from the start.

For experts, the academy can be a great place to share knowledge and get involved with new developments.

For more information on how you can get involved, visit the [CODA website](https://documentation.academy). 
You can also browse through the types of issues you’ll be working on in the [CODA repository](https://github.com/canonical/open-documentation-academy/issues).

Other channels to get involved with the academy include:

The CODA [discussion forum](https://discourse.ubuntu.com/c/community/open-documentation-academy/166) on the Ubuntu Community Hub

The Matrix channel **#documentation:ubuntu.com** for interactive chat

You can also [follow Canonical’s documentation team on Fosstodon](https://fosstodon.org/@CanonicalDocumentation) for the latest updates and information on upcoming events.

<!-- LINKS -->

[Ubuntu Code of Conduct]: https://ubuntu.com/community/ethos/code-of-conduct
[Ubuntu style guide]: https://docs.ubuntu.com/styleguide/en/
[Canonical contributor license agreement]: https://ubuntu.com/legal/contributors
[Canonical's Sphinx starter pack]: https://github.com/canonical/sphinx-docs-starter-pack
[Read the Docs]: https://about.readthedocs.com/
[Sphinx]: https://www.sphinx-doc.org/
[reStructuredText]: https://www.sphinx-doc.org/en/master/usage/restructuredtext/index.html
[MyST Markdown]: https://myst-parser.readthedocs.io/en/latest/
[starter pack guide]: https://canonical-starter-pack.readthedocs-hosted.com/stable/
[MyST syntax guide]: https://canonical-starter-pack.readthedocs-hosted.com/stable/reference/myst-syntax-reference/
[Diátaxis]: https://diataxis.fr/
[woke]: https://github.com/get-woke/woke
[related issues, pull requests, and repositories]: https://docs.github.com/en/get-started/writing-on-github/working-with-advanced-formatting/autolinked-references-and-urls
[conventional commits]: https://www.conventionalcommits.org/
