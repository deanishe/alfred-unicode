
<div align="center">
    <img width="128" height="128" src="https://raw.githubusercontent.com/deanishe/alfred-unicode/master/icons/icon-large.png" alt="workflow icon" title="workflow icon">
</div>

Alfred Unicode
==============

Search Unicode and copy/paste codepoints as HTML entities, Go or Python string literals, hexadecimal numbers etc.

<div align="center">
    <img width="570" height="505" src="https://raw.githubusercontent.com/deanishe/alfred-unicode/master/demo.png" alt="workflow demo" title="workflow demo">
</div>

<!-- MarkdownTOC autolink="true" bracket="round" levels="2,3,4,5" autoanchor="true" -->

- [Installation](#installation)
- [Usage](#usage)
    - [Input](#input)
- [Configuration](#configuration)
    - [Modifiers](#modifiers)
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
    - **Font: Arial Unicode MS (default)** — Font used to generate preview images.
        - `↩` — Choose a different font
        - `⌘↩` — Reset to default
    - **Text Colour: #FF2FA9 (default)** — Colour of text in preview images.
        - `↩` — Choose a different colour
        - `⌘↩` — Reset to default
    - **Background Colour: transparent (default)** — Colour of background of preview images.
        - `↩` — Choose a different colour
        - `⌘↩` — Reset to default
    - **Documentation** — Open this page in your browser.
    - **Report Bug** — Open [workflow issue tracker][issues] in your browser.
    - **Delete Cached Data** — Delete Unicode database and cached preview icons, forcing the workflow to re-create them.


<a id="input"></a>
### Input



<a id="configuration"></a>
Configuration
-------------

The workflow can be configured via several [workflow environment variables][config-sheet]. The most important ones can be set in the workflow by entering the query "prefs" in the Unicode search.

|       Variable      |  Default  |                                           Description                                           |
|---------------------|-----------|-------------------------------------------------------------------------------------------------|
| `BACKGROUND_COLOUR` | `<empty>` | Background colour of generated preview images (search result icons). Uses transparent if unset. |
| `FONT_NAME`         | `<empty>` | Font used to generate preview images. Uses "Arial Unicode MS" if unset.                         |
| `ICON_SIZE`         | `128`     | Width/height of preview images.                                                                 |
| `RESIZE_FACTOR`     | `1.0`     | Use larger/smaller font size in preview images.                                                 |
| `TEXT_COLOUR`       | `<empty>` | Text colour for preview images. Uses `#FF2FA9` if unset.                                        |
| `TEXT_SIZE`         | `<empty>` | Font size for preview images. Calculated automatically if unset.                                |

By default, the font size is calculated automatically based on the size of the image and the text. If the text is consistently too small or large with a particular font, you can increase/decrease `RESIZE_FACTOR` to increase/decrease the size of the text.

Alternatively, you can force a specific font size for all images by setting `TEXT_SIZE`.


<a id="modifiers"></a>
### Modifiers ###

You can assign your own modifiers to the search results to copy/paste quickly in your preferred format.

To assign a format to a modifier, add a new [workflow environment variable][config-sheet] with the name `MOD_<KEY1>[_<KEY2>[_<KEY3>]]` and set its value to the [short name](#formats) of the desired format.

For example, to override `⌘↩` to copy/paste JavaScript instead of the UTF-8 character, add a new variable named `MOD_CMD` with the value `javascript`. To re-assign UTF-8 character to `⌥⌘↩`, add a new variable named `MOD_CMD_OPT` (or `MOD_OPT_CMD` … doesn't matter) with the value `character`.

Valid modifiers are `CMD`, `OPT`, `SHIFT`, `CTRL`.


<a id="formats"></a>
Formats
-------

You can copy (and paste) each codepoint in some or all of the following formats:

|         Description         |   Short name  |       Example        |
|-----------------------------|---------------|----------------------|
| UTF-8 character             | `character`   | `🍻`                  |
| Codepoint name              | `name`        | `CLINKING BEER MUGS` |
| Unicode ID                  | `unicode`     | `U+1F37B`            |
| Hex literal (e.g. `0x00FF`) | `hexadecimal` | `0x1F37B`            |
| Python or Go escaped string | `python`      | `\U0001F37B`         |
| URL-encoded                 | `url-encoded` | `%F0%9F%8D%BB`       |
| JavaScript escaped string   | `javascript`  | `\u{1F37B}`          |
| CSS entity                  | `css`         | `\01F37B`            |
| HTML entity, hexadecimal    | `html-hex`    | `&#x1F37B;`          |
| HTML entity, decimal        | `html-dec`    | `&#127867;`          |
| HTML entity, named          | `html-named`  | `&szlig;`            |

**TIP**: You can get the short name for a format by using `⌘C` on its item in the detail view.

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
