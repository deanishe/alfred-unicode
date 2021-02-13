
<div align="center">
    <img width="128" height="128" src="https://raw.githubusercontent.com/deanishe/alfred-unicode/master/icons/icon-large.png" alt="workflow icon" title="workflow icon">
</div>

Alfred Unicode
==============

Explore Unicode characters in Alfred.

<!-- MarkdownTOC autolink="true" bracket="round" levels="2,3,4,5" autoanchor="true" -->

- [Installation](#installation)
- [Usage](#usage)
    - [Input](#input)
- [Configuration](#configuration)
    - [Modifiers](#modifiers)
        - [Short names](#short-names)
- [Formats](#formats)
- [Licensing & thanks](#licensing--thanks)

<!-- /MarkdownTOC -->

<a id="installation"></a>
Installation
------------

Download the latest version of the workflow from the [releases page][releases], then double-click the `Alfred Unicode-X.X.X.alfredworkflow` file to install.


<a id="usage"></a>
Usage
-----

- `unicode <query>` — Search Unicode codepoints (see [input](#input) below).
    - `⇧` or `⌘Y` — Quick Look image of selected codepoint.
    - `↩` or `⇥` — View selected codepoint in various formats (CSS, HTML, Python etc.).
        - `↩` or `⌘C` — Copy selected format to clipboard.
        - `⌘↩` — Paste selected format to frontmost application.
    - `⌘↩` ([configurable](#modifiers)) — Copy/paste Unicode character.
    - `⌥↩` ([configurable](#modifiers)) — Copy/paste Python Unicode escape.
    - `^↩` ([configurable](#modifiers)) — Copy/paste Unicode ID (U+00AB).
    - `⇧↩` ([configurable](#modifiers)) — Copy/paste HTML hexadecimal entity.
    - `...` — Add your own hotkeys (see [Modifiers](#modifiers)).
- `unicode prefs` — View workflow configuration
    - `search` — Enter keyword `search` to return to Unicode search.
    - **Workflow Is Up To Date** / **Update Available!** — Whether a new version of the workflow is available. Action this item to force an update check.
    - **Documentation** — Open this page in your browser.
    - **Report Bug** — Open [workflow issue tracker][issues] in your browser.
    - **Delete Cached Data** — Delete Unicode database and cached preview icons, forcing the workflow to re-create them.


<a id="input"></a>
### Input



<a id="configuration"></a>
Configuration
-------------




<a id="modifiers"></a>
### Modifiers ###

You can assign your own modifiers to the search results to copy/paste quickly in your preferred format.

To assign a format to a modifier, add a new [workflow environment variable][config-sheet] with the name `MOD_<KEY1>[_<KEY2>[_<KEY3>]]` and set its value to the [short name](#short-names) of the desired format.

For example, to override `⌘↩` to copy/paste JavaScript instead of Python, add a new variable named `MOD_CMD` with the value `javascript`. To re-assign Python to `⌥⌘↩`, add a new variable named `MOD_CMD_OPT` (or `MOD_OPT_CMD` … doesn't matter) with the value `python`.

Valid modifiers are `CMD`, `OPT`, `SHIFT`, `CTRL`.


<a id="short-names"></a>
#### Short names ####

You can get the short name for a format by using `⌘C` on its item in the detail view.

|   Short name  |         Description         |
|---------------|-----------------------------|
| `character`   | UTF-8 character             |
| `name`        | Codepoint name              |
| `unicode`     | Unicode ID                  |
| `python`      | Python or Go escaped string |
| `url-encoded` | URL-encoded                 |
| `javascript`  | JavaScript escaped string   |
| `css`         | CSS entity                  |
| `html-hex`    | HTML entity, hexadecimal    |
| `html-dec`    | HTML entity, decimal        |
| `html-named`  | HTML entity, named          |


<a id="formats"></a>
Formats
-------

You can copy (and paste) each codepoint in some or all of the following formats:

- UTF-8 character (e.g. `🍻`)
- Codepoint name (e.g. `CLINKING BEER MUGS`)
- Unicode ID (e.g. `U+1F37B`)
- Python or Go escaped string (e.g. `\U0001F37B`)
- URL-encoded (e.g. `%F0%9F%8D%BB`)
- JavaScript escaped string (e.g. `\u{1F37B}`)
- CSS entity (e.g. `\01F37B`)
- HTML entity, hexadecimal (e.g. `&#x1F37B;`)
- HTML entity, decimal (e.g. `&#127867;`)
- HTML entity, named (e.g. `&szlig;`)

**NOTE:** The JavaScript formatter outputs `\uABCD`-style escapes for codepoints up to 0xFFFF, and ES6 `\u{1ABCD}`-style escapes for codepoints above that.


<a id="licensing--thanks"></a>
Licensing & thanks
------------------

- The workflow is released under the [MIT licence][mit].
- It is heavily based on the [AwGo][awgo] library, also MIT licensed.
- The workflow icons are based on [Font Awesome][awesome], released under the [Creative Commons Attribution 4.0 licence][ccby40].

[releases]: https://github.com/deanishe/alfred-unicode/releases/latest
[issues]: https://github.com/deanishe/alfred-unicode/issues
[mit]: LICENCE.txt
[awesome]: https://github.com/FortAwesome/Font-Awesome
[ccby40]: https://creativecommons.org/licenses/by/4.0/legalcode
[awgo]: https://github.com/deanishe/awgo
[config-sheet]: https://www.alfredapp.com/help/workflows/advanced/variables/#environment
