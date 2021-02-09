#!/usr/bin/python
# encoding: utf-8
#
# Copyright (c) 2015 deanishe@deanishe.net
#
# MIT Licence. See http://opensource.org/licenses/MIT
#
# Created on 2015-04-26
#

"""Workflow configuration."""

from __future__ import print_function, absolute_import


import os

DB_FILE = './characters.sqlite'
DELIMITER = u'⟩'
KEYWORD = os.getenv('keyword', 'unicode')
# The font face to generate icons from by default.
# ArialUnicodeMS is currently the best choice I've found in terms of
# coverage.
FONT = os.getenv('preview_font', 'ArialUnicodeMS')
# The size (in pixels) of the generated icons. 256 is a decent size:
# it looks okay in Quick Look, and there should be max ~300 MB of icons.
ICON_SIZE = os.getenv('icon_size', '256')

# settings.json will be populated from these
DEFAULT_SETTINGS = {}

# For auto-updates
UPDATE_SETTINGS = {}

# URL shown in debugger if workflow fails
HELP_URL = None

# Lower to make icon generation faster
MAX_RESULTS = os.getenv('max_results', '50')

# Both of these variables will be adjusted in main() to be under the
# data and workflow directories respectively.

# Name of icons directory
ICON_DIR = 'icons'

# Name of IconGen binary (should be in workflow root)
ICONGEN = 'IconGen'

# Path to logfile
LOGFILE = None

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
    'python': 'Python or Go literal string',
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
    """Make config values relative to appropriate Workflow dirs."""
    global DB_FILE, ICON_DIR, ICONGEN, ICON_SIZE, \
        LOGFILE, KEYWORD, MAX_RESULTS
    DB_FILE = wf.workflowfile(DB_FILE)
    ICON_DIR = wf.datafile(ICON_DIR)
    ICONGEN = wf.workflowfile(ICONGEN)
    LOGFILE = wf.logfile
    KEYWORD = wf.decode(KEYWORD)
    ICON_SIZE = int(ICON_SIZE)
    MAX_RESULTS = int(MAX_RESULTS)
