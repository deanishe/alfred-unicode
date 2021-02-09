#!/usr/bin/env python
# encoding: utf-8
#
# Copyright (c) 2015 deanishe@deanishe.net
#
# MIT Licence. See http://opensource.org/licenses/MIT
#
# Created on 2015-04-26
#

"""Generate icons."""

from __future__ import print_function, absolute_import

import logging
import sqlite3
import subprocess
import time

import config
import util

log = logging.getLogger('workflow.icons')
log.addHandler(logging.NullHandler())


def generate_icons(codepoints, icondir=None, overwrite=False,
                   logfile=None, font=None, size=None):
    """Save icons for `codepoints` to `icondir`.

    Icons are generated using using `IconGen`.
    """
    if not codepoints:
        return

    icondir = icondir or config.ICON_DIR
    font = font or config.FONT
    size = size or config.ICON_SIZE

    log.debug('generating %d icons...', len(codepoints))
    start_time = time.time()

    cmd = [
        config.ICONGEN,
        '--outputdir', icondir,
        '--font', font,
        '--size', unicode(size),
        # '--verbose',
    ]
    if config.LOGFILE:
        cmd += ['--logfile', config.LOGFILE]
    if overwrite:
        cmd += ['--overwrite']
    cmd = [s.encode('utf-8') for s in cmd]
    log.debug('cmd=%r', cmd)

    proc = subprocess.Popen(
        cmd,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
    )
    stdout, _ = proc.communicate('\n'.join(codepoints).encode('utf-8'))
    if proc.returncode != 0:
        log.error('generate icons: %r', codepoints)
        log.error('IconGen output:\n%s', stdout.decode('utf-8'))
        return False

    log.info('%d icon(s) generated in %s', len(codepoints),
             util.readable_time(time.time() - start_time))

    return True


def main(wf):
    """Run when called as main script to generate all icons."""
    import os
    config.init_config(wf)

    con = sqlite3.connect(config.DB_FILE)
    cursor = con.cursor()
    cursor.execute("""SELECT hex, icon FROM chars ORDER BY hex DESC""")
    codepoints = [t[0] for t in cursor.fetchall() if not
                  os.path.exists(os.path.join(config.ICON_DIR, t[1]))]

    generate_icons(codepoints, config.ICON_DIR)


if __name__ == '__main__':
    import sys
    from workflow import Workflow3
    wf = Workflow3()
    log = wf.logger
    sys.exit(wf.run(main, text_errors=True))
