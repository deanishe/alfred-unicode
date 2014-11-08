#!/usr/bin/env python
# encoding: utf-8
#
# Copyright © 2014 deanishe@deanishe.net
#
# MIT Licence. See http://opensource.org/licenses/MIT
#
# Created on 2014-11-06
#

"""

Generate an sqlite database from `characters.tsv`

"""

from __future__ import print_function, unicode_literals, absolute_import

import csv
import logging
import os
import sqlite3
import sys
from time import time

logging.basicConfig(
    format='%(asctime)s %(levelname)8s %(filename)s:%(lineno)d %(message)s',
    datefmt='%H:%M:%S',
    level=logging.DEBUG)

log = logging.getLogger('')


CHARS_TSV = os.path.join(os.path.dirname(__file__), 'characters.tsv')
ICONS_TSV = os.path.join(os.path.dirname(__file__), 'icons.tsv')
DB_FILE = 'characters.sqlite'


def create_index_db():
    """Create a "virtual" table, which sqlite users for full-text search"""
    log.info('Creating index database at `{}`'.format(DB_FILE))
    con = sqlite3.connect(DB_FILE)
    with con:
        cur = con.cursor()
        cur.execute("""CREATE VIRTUAL TABLE
                    chars USING
                    fts3(name, hex, entity, icon)""")


def update_index_db():
    """Read in data file and add it to index database"""
    start = time()
    icons = {}

    log.info('Updating search database at `{}`'.format(DB_FILE))

    with open(ICONS_TSV, 'rb') as fp:
        reader = csv.reader(fp, delimiter=b'\t')
        for row in reader:
            h, icon = row
            icons[h] = icon

    con = sqlite3.connect(DB_FILE)
    count = 0
    with con:
        cur = con.cursor()
        with open(CHARS_TSV, 'rb') as fp:
            reader = csv.reader(fp, delimiter=b'\t')
            for row in reader:
                name, h, entity = [v.decode('utf-8') for v in row]
                icon = icons.get(h, '')
                cur.execute("""INSERT OR IGNORE INTO
                            chars (name, hex, entity, icon)
                            VALUES (?, ?, ?, ?)
                            """, (name, h, entity, icon))
                count += 1
    log.info('{:d} characters added/updated in {:0.3f} seconds'.format(
             count, time() - start))


def main():
    if not os.path.exists(DB_FILE):
        create_index_db()
    update_index_db()

if __name__ == '__main__':
    sys.exit(main())
