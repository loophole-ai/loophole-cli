# Loophole Notice

Loophole CLI is a fork of `blu-code`.

## License

Loophole is licensed under the GNU Affero General Public License Version 3.0
(AGPL-3.0). See [LICENSE.txt](./LICENSE.txt) for the full license text.

## Copyright

Loophole is Copyright 2026 Loophole AI. All rights reserved.

## Original Projects

This project incorporates code from the following original projects:

- [blu-code](https://github.com/getblu/blu-code) - Licensed under the MIT License
  - This is the root project. Loophole CLI is a rebranded and reworked derivative of it.
  - Originally created by Kujtim Hoxha and later developed by Garv Agnihotri (Get-Blu).
  - See [LICENSE-blu-code.txt](./LICENSE-blu-code.txt) for the full MIT License text.

## License Compatibility

The MIT License is permissive and permits relicensing of derivative works, so
distributing this project under AGPL-3.0 is permitted. AGPL-3.0 is compatible
with the MIT License.

In keeping with the conditions of the MIT License, the original copyright
notices and permission notice from `blu-code` are retained verbatim in
[LICENSE-blu-code.txt](./LICENSE-blu-code.txt). Portions of this project that
originate from that work remain subject to the MIT License; the remainder is
distributed under AGPL-3.0.

## What changed in this fork

This fork renames the project, its CLI binary and command, Go module path,
package metadata, configuration file names, internal identifiers, and
documentation from "blu" / "blu-code" / "Get-Blu" to "loophole" /
"loophole-cli" / "loophole-ai" / "Loophole AI". It does not claim original
authorship of the underlying functionality, which remains the work of the
upstream authors.

Beyond the rename, this fork reworks substantial parts of the project,
including replacing the model catalogue with [models.dev](https://models.dev),
adding support for any OpenAI-compatible endpoint, and correcting behaviour on
Windows.

## Contributors

Thank you to everyone who has contributed to this project.
