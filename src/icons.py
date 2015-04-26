#!/usr/bin/env python
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

import datetime
import logging
import sqlite3
import subprocess
import time

import config
import util

log = logging.getLogger('workflow.icons')
log.addHandler(logging.NullHandler())


def generate_icons(codepoints, output_directory,
                   overwrite=False, logfile=None,
                   font=None, size=None):
    """Save icons for `codepoints` to `output_directory` using `IconGen`"""
    if not len(codepoints):
        return
    font = font or config.DEFAULT_FONT
    size = size or config.DEFAULT_SIZE
    log.debug('Generating {0} icons...'.format(len(codepoints)))
    start_time = time.time()
    cmd = [
        config.ICONGEN,
        '--outputdir', output_directory,
        '--font', '{}'.format(font),
        '--size', '{}'.format(size),
        # '--verbose',
        ]
    if logfile:
        cmd += ['--logfile', logfile]
    if overwrite:
        cmd += ['--overwrite']
    cmd = [s.encode('utf-8') for s in cmd]
    log.debug('cmd : {}'.format(cmd))
    proc = subprocess.Popen(
        cmd,
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
    )
    stdout, _ = proc.communicate('\n'.join(codepoints).encode('utf-8'))
    if proc.returncode != 0:
        log.error('Error generating icons : {}'.format(codepoints))
        log.error('IconGen output:\n{}'.format(stdout.decode('utf-8')))
        return False
    log.debug('{} icons generated in {:0.3f}s'.format(
              len(codepoints), time.time() - start_time))
    return True


def main(wf):
    """Run when called as main script to generate all icons"""
    import os
    config.init_config(wf)
    start = time.time()
    log.debug('Starting icon generation at {}'.format(
              datetime.datetime.fromtimestamp(start)
              .strftime('%Y-%m-%d %H:%M:%S')))
    con = sqlite3.connect(config.DB_FILE)
    cursor = con.cursor()
    cursor.execute("""SELECT hex, icon FROM chars ORDER BY hex DESC""")
    codepoints = [t[0] for t in cursor.fetchall() if not
                  os.path.exists(os.path.join(config.ICON_DIR, t[1]))]

    generate_icons(codepoints,
                   config.ICON_DIR,
                   logfile=wf.logfile,
                   font=wf.settings.get('font'),
                   size=wf.settings.get('size'))

    stop = time.time()
    log.debug('Finished icon generation at {}'.format(
              datetime.datetime.fromtimestamp(stop)
              .strftime('%Y-%m-%d %H:%M:%S')))
    log.info('{} icons generated in {}'.format(
             len(codepoints),
             util.readable_time(stop - start)))


if __name__ == '__main__':
    import sys
    from workflow import Workflow
    wf = Workflow()
    log = wf.logger
    sys.exit(wf.run(main))
