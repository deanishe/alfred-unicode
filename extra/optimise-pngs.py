#!/usr/bin/env python
# encoding: utf-8
#
# Copyright © 2014 deanishe@deanishe.net
#
# MIT Licence. See http://opensource.org/licenses/MIT
#
# Created on 2014-11-10
#

"""
"""

from __future__ import print_function, unicode_literals, absolute_import

import os
import subprocess
import sys


def pngcrush(filepath):
    """Run file through `pngcrush` and return SHA1 hash"""
    name, ext = os.path.splitext(filepath)
    temppath = '{}.{}{}'.format(name, os.getpid(), ext)
    size_in = os.stat(filepath).st_size
    os.rename(filepath, temppath)
    cmd = [
        'pngcrush',
        '-rem', 'allb',
        '-m', '10',
        '-q',
        '-reduce',
        temppath,
        filepath,
    ]
    subprocess.call(cmd)
    if os.path.exists(filepath) and os.path.exists(temppath):
        os.unlink(temppath)
    size_out = os.stat(filepath).st_size
    pc = (float(size_out) / size_in) * 100
    print('Optimised [{:4d}b / {:0.1f}%] `{}`'.format(
          size_out, pc, filepath), file=sys.stderr)


def main():
    if not len(sys.argv) == 2:
        print('Usage: optimise-pngs.py <rootdir>')
        return 1

    rootdir = sys.argv[1]

    for root, dirnames, filenames in os.walk(rootdir):
        for filename in filenames:
            if not filename.lower().endswith('.png'):
                continue
            filepath = os.path.join(root, filename)
            pngcrush(filepath)


if __name__ == '__main__':
    sys.exit(main())
