#!/usr/bin/python
# encoding: utf-8
#
# Copyright © 2015 deanishe@deanishe.net
#
# MIT Licence. See http://opensource.org/licenses/MIT
#
# Created on 2015-04-26
#

"""
"""

from __future__ import print_function, unicode_literals, absolute_import


import sys

from workflow import Workflow

log = None
DB_FILE = 'characters.sqlite'
DELIMITER = '⟩'
KEYWORD = 'unicode'

# The font face to generate icons from by default.
# ArialUnicodeMS is currently the best choice I've found in terms of
# coverage.
DEFAULT_FONT = 'ArialUnicodeMS'
# The size (in pixels) of the generated icons. 256 is a decent size:
# it looks okay in QuickLook, and there should only be ~30 MB of icons.
DEFAULT_SIZE = 256

# settings.json will be populated from these
DEFAULT_SETTINGS = {
    'font': DEFAULT_FONT,
    'size': DEFAULT_SIZE,
}

# For auto-updates
UPDATE_SETTINGS = {
}

HELP_URL = None

# Lower to make icon generation faster
MAX_RESULTS = 50

# Both of these variables will be adjusted in main() to be under the
# data and workflow directories respectively.

# Name of icon cache directory
ICON_DIR = 'icons'
# Path to IconGen binary (should be in workflow root)
ICONGEN = 'IconGen'

# Map keys in dictionary returned by `charinfo()` to proper names
key_name_map = {
    'char': 'Character',
    'name': 'Name',
    'hex': 'Hexadecimal',
    'entity_named': 'HTML named entity',
    'entity_hex': 'HTML hexadecimal entity',
    'entity_dec': 'HTML decimal entity',
    'entity_css': 'CSS entity',
    'url_encoded': 'URL-encoded',
    'python': 'Python literal string',
}

# Which keys (from `charinfo()` dict) to show in results
display_keys = (
    'char',
    # 'name',
    'hex',
    'entity_named',
    'entity_hex',
    'entity_dec',
    'entity_css',
    'url_encoded',
    'python',
)


def init_config(wf):
    """Make config values relative to appropriate Workflow dirs"""
    global DB_FILE, ICON_DIR, ICONGEN
    DB_FILE = wf.workflowfile(DB_FILE)
    ICON_DIR = wf.datafile(ICON_DIR)
    ICONGEN = wf.workflowfile(ICONGEN)
    ICONGEN = ('/Volumes/Users/daj/Library/Developer/Xcode/DerivedData/'
               'IconGenerationHelper-bbowvswngqbbjugdvgdqtwpjvgdp/'
               'Build/Products/Debug/IconGen')


# ---------------------------------------------------------
# The below doesn't do anything yet...

def main(wf):
    """Configure settings"""
    init_config(wf)
    query = None
    if len(wf.args):
        query = wf.args[0]


if __name__ == '__main__':
    wf = Workflow(default_settings=DEFAULT_SETTINGS,
                  update_settings=UPDATE_SETTINGS,
                  help_url=HELP_URL)
    log = wf.logger
    sys.exit(wf.run(main))
