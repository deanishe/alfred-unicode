#!/usr/bin/python
# encoding: utf-8
#
# Copyright © 2014 deanishe@deanishe.net
#
# MIT Licence. See http://opensource.org/licenses/MIT
#
# Created on 2014-11-06
#

"""Workflow Script Filter to search and view Unicode characters.

Selecting a character will display a variety of encoded forms, e.g.
HTML entities, CSS entities, Python Unicode literal.
"""

from __future__ import print_function, absolute_import


import os
import shutil
import sqlite3
import struct
import sys
import time
from urllib import quote

from workflow import Workflow3, ICON_WARNING
from workflow import background
from workflow.util import run_trigger

import config


log = None


class InvalidQuery(Exception):
    """Raised if an SQLite query is invalid."""


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
        matchinfo = [struct.unpack(b'I', matchinfo[i:i + 4])[0]
                     for i in range(0, bufsize, 4)]
        it = iter(matchinfo[2:])
        return sum(x[0] * w / x[1]
                   for x, w in zip(zip(it, it, it), weights)
                   if x[1])
    return rank


def charinfo(name):
    """Return ``dict`` of representations of named codepoint."""
    con = sqlite3.connect(config.DB_FILE)
    cursor = con.cursor()
    cursor.execute(
        """SELECT name, hex, entity, icon FROM chars WHERE name = ?""",
        (name,))

    result = cursor.fetchone()

    if result is None:
        return None

    name, h, entity, icon = result

    dec = int(h, 16)
    s = '\\U{:08x}'.format(dec)
    u = s.decode('unicode-escape')
    # log.debug('Unicode : {}'.format(u))
    info = {
        'codepoint': h,
        'icon': icon,
        'char': u,
        'name': name,
        'hex': 'U+{}'.format(h),
        'entity_named': None,
        'entity_hex': '&#x{};'.format(h.lstrip('0')),
        'entity_dec': '&#{};'.format(dec),
        'entity_css': '\\{:0>6}'.format(h),
        'url_encoded': '{}'.format(quote(u.encode('utf-8'))),
        'python': '\\U{}'.format(h),
    }
    if entity:
        info['entity_named'] = '&{};'.format(entity)
    # log.debug('charinfo : {}'.format(info))
    return info


def _search(query):
    pass


def search(query, cursor):
    """Search SQLite database."""
    try:
        cursor.execute("""
            SELECT name, hex, entity, icon
                FROM
                    (SELECT
                        rank(matchinfo(chars)) AS r,
                        name,
                        hex,
                        entity,
                        icon
                     FROM chars
                     WHERE chars MATCH ?)
            ORDER BY r
            DESC
            LIMIT ?
            """, (query, config.MAX_RESULTS))
        results = cursor.fetchall()
    except sqlite3.OperationalError as err:
        if b'malformed MATCH' in err.message:
            raise InvalidQuery(query)
        else:
            raise err

    log.debug(u'%d results for `%s`', len(results), query)
    return results


def generate_all_icons():
    """Magic argument callback to generate all icons."""
    if background.is_running('icongen'):
        return 'Generation already in progress.'

    background.run_in_background(
        'icongen',
        ['/usr/bin/python', wf.workflowfile('icons.py')]
    )
    return 'Starting icon generation. This may take up to 15 minutes.'


def clear_all_icons():
    """Magic argument callback to delete icon directory."""
    if os.path.exists(config.ICON_DIR) and os.path.isdir(config.ICON_DIR):
        shutil.rmtree(config.ICON_DIR)
        return 'Deleted icons'
    else:
        return 'No icons to delete'


def do_details(wf, query):
    """Show details for a single codepoint."""
    from icons import generate_icons
    # Icons that don't yet exist
    icons_to_generate = set()
    name, query = [s.strip() for s in query.split(config.DELIMITER)]
    log.info('Details for character `%s`', name)
    info = charinfo(name)

    if not info:
        wf.add_item('Unknown character', name, icon=ICON_WARNING)

    else:

        results = []
        icon = os.path.join(config.ICON_DIR, info['icon'])
        if not os.path.exists(icon):
            icons_to_generate.add(info['codepoint'])

        # Compile codepoint data
        for key in config.display_keys:
            name = config.key_name_map[key]
            value = info[key]
            if not value:
                continue
            results.append((name, value))

        if query:
            results = wf.filter(query, results, lambda t: t[0], min_score=40)

        if not results:
            wf.add_item('No matching representations', icon=ICON_WARNING)

        for (name, value) in results:
            wf.add_item(value,
                        name,
                        valid=True,
                        copytext=value,
                        largetext=value,
                        arg=value,
                        icon=icon)

    if icons_to_generate:
        wf.rerun = 0.5

    wf.send_feedback()

    # Generate icons if need be
    generate_icons(icons_to_generate, config.ICON_DIR)


def do_search(wf, query):
    """Perform search."""
    from icons import generate_icons
    # Icons that don't yet exist
    icons_to_generate = set()
    # Cache user-entered query, so we can jump back to it later
    wf.cache_data('last-query', query)

    start = time.time()
    con = sqlite3.connect(config.DB_FILE)
    con.create_function('rank', 1, make_rank_func((1.0, 0, 1.5, 0)))
    cursor = con.cursor()
    results = search(query, cursor)
    if not results:
        results = search(query + '*', cursor)
        if not results:
            results = search('*' + query + '*', cursor)

    if not results:
        wf.add_item('No matches', 'Try a different query', icon=ICON_WARNING)

    log.info(u'%d results for `%s` in %0.3f seconds',
             len(results), query, time.time() - start)

    for (name, h, entity, icon) in results:
        log.debug('name=%r, h=%r, entity=%r, icon=%r', name, h, entity, icon)
        icon = os.path.join(config.ICON_DIR, icon)
        if not os.path.exists(icon):
            icons_to_generate.add(h)
            # icon = os.path.join('icons', 'unknown.png')
        subtitle = 'U+{}'.format(h)
        wf.add_item(name,
                    subtitle,
                    autocomplete=u'{} {} '.format(name, config.DELIMITER),
                    copytext=subtitle,
                    arg=icon,
                    type='file',
                    icon=icon)

    # Generate icons if need be
    generate_icons(icons_to_generate, config.ICON_DIR)

    if icons_to_generate:
        wf.rerun = 0.5

    wf.send_feedback()


def main(wf):
    """Run workflow."""
    config.init_config(wf)
    # Add custom magic arguments
    wf.magic_arguments['genicons'] = generate_all_icons
    wf.magic_arguments['delicons'] = clear_all_icons

    query = wf.args[0] if wf.args else ''

    # Back up
    if query.endswith(config.DELIMITER):
        query = wf.cached_data('last-query', max_age=0)
        if query:
            query = u'{} {}'.format(config.KEYWORD, query)
        else:
            query = config.KEYWORD + ' '

        run_trigger('search', arg=query)
        return

    if config.DELIMITER in query:
        return do_details(wf, query)

    # Plain query (no delimiter), so search Unicode codepoints
    return do_search(wf, query)


if __name__ == '__main__':
    wf = Workflow3(default_settings=config.DEFAULT_SETTINGS,
                   update_settings=config.UPDATE_SETTINGS,
                   help_url=config.HELP_URL)
    log = wf.logger
    sys.exit(wf.run(main))
