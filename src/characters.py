#!/usr/bin/python
# encoding: utf-8
#
# Copyright © 2014 deanishe@deanishe.net
#
# MIT Licence. See http://opensource.org/licenses/MIT
#
# Created on 2014-11-06
#

"""
Workflow Script Filter to search and view Unicode characters.

Selecting a character will display a variety of encoded forms, e.g.
HTML entities, CSS entities, Python Unicode literal.
"""

from __future__ import print_function, unicode_literals, absolute_import


import os
import sqlite3
import struct
import subprocess
import sys
from time import time
from urllib import quote_plus

from workflow import Workflow, ICON_WARNING

log = None
DB_FILE = 'characters.sqlite'
DELIMITER = '⟩'
KEYWORD = 'unicode'


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


def run_alfred(query):
    """Call Alfred with ``query``"""
    subprocess.call([
        'osascript', '-e',
        'tell application "Alfred 2" to search "{}"'.format(query)])


# Search ranking function
# Adapted from http://goo.gl/4QXj25 and http://goo.gl/fWg25i
def make_rank_func(weights):
    """`weights` is a list or tuple of the relative ranking per column.

    Use floats (1.0 not 1) for more accurate results. Use 0 to ignore a
    column.
    """
    def rank(matchinfo):
        # matchinfo is defined as returning 32-bit unsigned integers
        # in machine byte order
        # http://www.sqlite.org/fts3.html#matchinfo
        # and struct defaults to machine byte order
        bufsize = len(matchinfo)  # Length in bytes.
        matchinfo = [struct.unpack(b'I', matchinfo[i:i+4])[0]
                     for i in range(0, bufsize, 4)]
        it = iter(matchinfo[2:])
        return sum(x[0]*w/x[1]
                   for x, w in zip(zip(it, it, it), weights)
                   if x[1])
    return rank


def charinfo(name):
    con = sqlite3.connect(DB_FILE)
    cursor = con.cursor()
    cursor.execute("""SELECT name, hex, entity, icon FROM
                   chars WHERE name = ?""", (name,))
    result = cursor.fetchone()
    if result is None:
        return None
    name, h, entity, icon = result
    dec = int(h, 16)
    s = '\\U{:08x}'.format(dec)
    u = s.decode('unicode-escape')
    # log.debug('Unicode : {}'.format(u))
    info = {
        'icon': icon,
        'char': u,
        'name': name,
        'hex': 'U+{}'.format(h),
        'entity_named': None,
        'entity_hex': '&#x{};'.format(h),
        'entity_dec': '&#{};'.format(dec),
        'entity_css': '\\{:0>6}'.format(h),
        'url_encoded': '{}'.format(quote_plus(u.encode('utf-8'))),
        'python': '\\u{}'.format(h),
    }
    if entity:
        info['entity_named'] = '&{};'.format(entity)
    log.debug('charinfo : {}'.format(info))
    return info


def search(query, cursor):
    try:
        cursor.execute("""SELECT name, hex, entity, icon FROM
                    (SELECT rank(matchinfo(chars))
                    AS r, name, hex, entity, icon
                    FROM chars WHERE chars MATCH ?)
                    ORDER BY r DESC LIMIT 100""", (query, ))
        results = cursor.fetchall()
    except sqlite3.OperationalError as err:
        if b'malformed MATCH' in err.message:
            wf.add_item('Invalid query', icon=ICON_WARNING)
            wf.send_feedback()
            return
        else:
            raise err
    log.debug('{:d} results for `{}`'.format(len(results), query))
    return results


def main(wf):
    query = None
    if len(wf.args):
        query = wf.args[0]

    # Back up
    if query.endswith(DELIMITER):
        query = wf.cached_data('last-query', max_age=0)
        if query:
            query = '{} {}'.format(KEYWORD, query)
        else:
            query = '{} '.format(KEYWORD)
        run_alfred(query)
        return

    if DELIMITER in query:  # Show character details
        name, query = [s.strip() for s in query.split(DELIMITER)]
        log.info('Details for character `{}`'.format(name))
        info = charinfo(name)
        if not info:
            wf.add_item('Unknown character', name, icon=ICON_WARNING)
        else:
            output = []
            for key in display_keys:
                name = key_name_map[key]
                value = info[key]
                if not value:
                    continue
                output.append((name, value))
            if query:
                output = wf.filter(query, output, lambda t: t[0], min_score=40)
            if not output:
                wf.add_item('No matching representations', icon=ICON_WARNING)
            for (name, value) in output:
                wf.add_item(
                    value,
                    name,
                    valid=True,
                    copytext=value,
                    largetext=value,
                    arg=value,
                    icon='icons/{}'.format(info['icon']),
                )
        wf.send_feedback()
        return
    else:  # Cache user-entered query
        wf.cache_data('last-query', query)

    start = time()
    con = sqlite3.connect(DB_FILE)
    con.create_function('rank', 1, make_rank_func((1.0, 0, 1.5, 0)))
    cursor = con.cursor()
    results = search(query, cursor)
    if not results:
        results = search(query + '*', cursor)
        if not results:
            results = search('*{}*'.format(query), cursor)

    if not results:
        wf.add_item('No matches', 'Try a different query', icon=ICON_WARNING)

    log.info('{:d} results for `{}` in {:0.3f} seconds'.format(
             len(results), query, time() - start))

    for (name, h, entity, icon) in results:
        icon = os.path.join('icons', icon)
        if not os.path.exists(icon):
            icon = os.path.join('icons', 'unknown.png')
        subtitle = 'U+{}'.format(h)
        wf.add_item(name,
                    subtitle,
                    autocomplete='{} {} '.format(name, DELIMITER),
                    copytext=subtitle,
                    icon=icon)

    wf.send_feedback()


if __name__ == '__main__':
    wf = Workflow()
    log = wf.logger
    sys.exit(wf.run(main))
