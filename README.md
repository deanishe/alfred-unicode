
<div align="center" class="centred">
    <img width="128" height="128" src="https://raw.githubusercontent.com/deanishe/alfred-unicode/master/icons/icon-large.png" alt="workflow icon" title="workflow icon">
</div>


# Alfred Unicode #

Search Unicode and copy/paste characters as HTML entities, Go or Python string literals, hexadecimal numbers etc. in Alfred 4+.

<div align="center" class="centred">
    <img width="570" height="505" src="https://raw.githubusercontent.com/deanishe/alfred-unicode/master/demo.gif" alt="animated workflow demo" title="animated workflow demo">
</div>

<!-- MarkdownTOC autolink="true" bracket="round" levels="2,3,4,5" autoanchor="true" -->

- [Installation](#installation)
- [Usage](#usage)
    - [Input](#input)
- [Configuration](#configuration)
    - [Modifiers](#modifiers)
- [Formats](#formats)
    - [Notes](#notes)
- [Licensing & thanks](#licensing--thanks)

<!-- /MarkdownTOC -->


<a id="installation"></a>
## Installation ##

Download the latest version of the workflow from the [releases page][releases], then double-click the `Alfred Unicode-X.X.X.alfredworkflow` file to install.


<a id="usage"></a>
## Usage ##

- `.u <query>` — Search Unicode characters (see [input](#input) below).
    - `⇧` or `⌘Y` — Quick Look image of selected character.
    - `↩` or `⇥` — View selected character in various formats (CSS, HTML, Python etc.).
        - `↩` or `⌘C` — Copy selected format to clipboard.
        - `⌘↩` — Paste selected format to frontmost application.
    - `⌘↩` ([configurable](#modifiers)) — Copy/paste UTF-8 character.
    - `⌥↩` ([configurable](#modifiers)) — Copy/paste Python Unicode escape.
    - `^↩` ([configurable](#modifiers)) — Copy/paste Unicode ID (U+00AB).
    - `⇧↩` ([configurable](#modifiers)) — Copy/paste HTML hexadecimal entity.
    - `...` — Add your own hotkeys (see [Modifiers](#modifiers)).
- `.u help` — Open workflow help file
- `.u prefs` — View workflow configuration
    - `search` — Enter keyword `search` to return to Unicode search.
    - **Workflow Is Up To Date** / **Update Available!** — Whether a new version of the workflow is available. Action this item to force an update check.
    - **Formatters** — Disable formatters you don't want.
        - `↩` — To view formatter list and to toggle formatters on/off.
    - **ES6 Mode** — Use ES6-style `\u{1F680}` escapes in JS for characters above 0xFFFF.
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

Configure the included Hotkey to allow you to search directly from other applications (select a character or Unicode escape and hit the Hotkey).


<a id="input"></a>
### Input ###

Use the special search queries `search`, `prefs` and `help` to navigate to the workflow's search, preferences and help file respectively. These work throughout the workflow.

**Single-character queries** are treated as exact matches for their own codepoint: if you enter `p`, you will only see the result for U+0070 `LATIN SMALL LETTER P`. To search for *all* P's, add a space after `p` (i.e. `p `).

You can search for characters based on their **names** (e.g. `LATIN CAPITAL LETTER U WITH DIAERESIS`), **aliases** (e.g. `EOL`/`LF`/`NL`/`NEW LINE` for `LINE FEED (LF)`), **hexadecimal** (e.g. `fe`, `00fe`, `0xfe`) or **decimal** (e.g. `32`), or HTML entity name (e.g. `raquo`, `Uuml`).

The workflow understands all formats it generates, so you can also enter a Go string escape (e.g. `\U0001F37B`) or Unicode ID (e.g. `U+1F4A9`). Configure the included Hotkey to be able to look up escapes directly from your editor.

You can restrict search to a particular field with the syntax `<field>:<query>`, e.g. `entity:copy` to find characters whose HTML entity name matches `copy`. The available fields are:

|   Field   |                     Description                      |
|-----------|------------------------------------------------------|
| `name`    | Character name, e.g. `PILE OF POO`                   |
| `aliases` | Alternate names, e.g. `BS` for `BACKSPACE`           |
| `hex`     | 4+ character hexadecimal representation, e.g. `00FE` |
| `dec`     | Decimal representation, e.g. `254`                   |
| `entity`  | HTML entity name (if any), e.g. `szlig`              |

**Note**: All fields are treated as text, so `hex:20` will return characters `0x20HH`, *not* `0x20` (`dec:32` will match `0x20`, but also `32NN`). Use `0x` format to specify that a query is hexadecimal, i.e. `0x20`. Append a space to `dec:NN` queries to only match that specific character (i.e. `dec:32 `).

Finally, you can also use `AND` or `OR` in your queries, e.g. `arrow AND left OR right` (note `AND` and `OR` **must** be uppercase).


<a id="configuration"></a>
## Configuration ##

The workflow can be configured via several [workflow environment variables][config-sheet]. The most important ones can be set in the workflow by entering the query "prefs" in the Unicode search.

|       Variable      |  Default  |                                           Description                                           |
|---------------------|-----------|-------------------------------------------------------------------------------------------------|
| `BACKGROUND_COLOUR` | `<empty>` | Background colour of generated preview images (search result icons). Uses transparent if unset. |
| `DISABLED`          | `<empty>` | Comma-separated list of disabled formatters.                                                    |
| `ES6_MODE`          | `<empty>` | Use ES6-style escapes for characters over 0xFFFF. If unset, surrogate pairs are used.           |
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
## Formats ##

You can copy (and paste) each character in some or all of the following formats:

|         Description         |   Short name  |       Example        |
|-----------------------------|---------------|----------------------|
| UTF-8 character             | `character`   | `🍻`                  |
| Character name              | `name`        | `CLINKING BEER MUGS` |
| Unicode ID                  | `unicode`     | `U+1F37B`            |
| Hex literal                 | `hexadecimal` | `0x1F37B`            |
| Decimal literal             | `decimal`     | `127867`             |
| Python or Go escaped string | `python`      | `\U0001F37B`         |
| URL-encoded                 | `urlencoded`  | `%F0%9F%8D%BB`       |
| JavaScript escaped string   | `javascript`  | `\u{1F37B}`          |
| CSS entity                  | `css`         | `\01F37B`            |
| HTML entity, hexadecimal    | `htmlhex`     | `&#x1F37B;`          |
| HTML entity, decimal        | `htmldec`     | `&#127867;`          |
| HTML entity, named          | `htmlnamed`   | `&szlig;`            |

**TIP**: You can get the short name for a format by using `⌘C` on its item in the detail view.


<a id="notes"></a>
### Notes ###

- The Python/Go formatter outputs `\u00AB`-style escapes for characters up to 0xFFFF and `\U000ABCDE`-style escapes for higher ones.
- The JavaScript formatter outputs `\xAB`-style escapes for characters up to 0xFF, `\uABCD`-style escapes for characters up to 0xFFFF, and surrogate pairs for characters above that. Turn on **ES6 Mode** in the workflow settings to use ES6 `\u{1ABCD}`-style escapes instead of surrogate pairs.


<a id="licensing--thanks"></a>
## Licensing & thanks ##

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
