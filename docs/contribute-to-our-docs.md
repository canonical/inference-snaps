(contribute-to-our-docs)=

# Contributing to our docs
We are committed to continuously improving our documentation and making it more useful for both ourselves and others.

We welcome community contributions, suggestions, and fixes.

## Ways to contribute
You can help improve this documentation by:

- Documenting new features or improvements you contribute to the code.
- Clarifying concepts or common questions based on your own experience.
- Fixing typos, formatting errors, and other mistakes.
- Reporting documentation issues by opening an issue on GitHub.

## AI Snaps documentation overview

The AI Snaps documentation is hosted on GitHub, built with [Sphinx](https://www.sphinx-doc.org/), and published on [Read the Docs](https://about.readthedocs.com/).
To contribute to these docs, you'll need a GitHub account.

The documentation is organized according to the [Diátaxis](https://diataxis.fr/) framework.
It is written in Markdown with [MyST](https://myst-parser.readthedocs.io/en/latest/) extensions. 

You can build the documentation locally by following the build instructions for the sphinx starterpack.

## Style and language

This documentation is aimed at helping both beginners and experts. 
We try to use appropriate writing for the subject, use inclusive language, and assume the reader has as little prior knowledge as possible.

The navigation, structure, style, and content of our documentation follows the Diátaxis framework for technical documentation authoring. 
This splits pages into tutorials, how-to guides, reference material, and explanations.

Some general tips to ensure high quality documentation are:

- use a spell checker

- resist being overly formal

- resist being overly verbose

- Use links and examples from reputable sources, e.g., official upstream documentation.

We also recommend using the Ubuntu style guide.

You are welcome to suggest changes to any documentation topic with updated or more insightful information. 
We aim for consistency, but don’t let formality prevent you from contributing.

## Automatic documentation checks
GitHub runs automatic checks on the documentation to verify the spelling, the validity of links, correct formatting of the Markdown files, and the use of inclusive language.

You can (and should!) run these tests locally before pushing your changes:

- Check the spelling: make doc-spellcheck (or make doc-spelling to first build the documentation and then check it)

- Check the validity of links: make doc-linkcheck

- Check the Markdown formatting: make doc-lint

- Check for inclusive language: make doc-woke

## Open Documentation Academy
The AI model snaps project is a proud member of the Canonical Open Documentation Academy (CODA).

The Open Documentation Academy is an initiative that aims to encourage open source contributions from the community.
We use documentation as the gateway, thus lowering the barrier to entry.

Through the academy, community members get access to curated and well-defined beginner-friendly tasks along with mentorship and advice in a friendly and encouraging environment. 
This allows them to make meaningful contributions to established open-source projects, such as Launchpad, from the start.

For experts, the academy can be a great place to share knowledge and get involved with new developments.

For more information on how you can get involved, visit the [CODA website](https://documentationacademy.org/). 
You can also browse through the types of issues you’ll be working on in the [CODA repository](https://github.com/canonical/open-documentation-academy/issues).

Other channels to get involved with the academy include:

The CODA [discussion forum](https://discourse.ubuntu.com/c/community/open-documentation-academy/166) on the Ubuntu Community Hub

The Matrix channel **#documentation:ubuntu.com** for interactive chat

You can also [follow Canonical’s documentation team on Fosstodon](https://fosstodon.org/@CanonicalDocumentation) for the latest updates and information on upcoming events.
