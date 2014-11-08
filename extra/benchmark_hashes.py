#!/usr/bin/env python
# encoding: utf-8
#
# Copyright © 2014 deanishe@deanishe.net
#
# MIT Licence. See http://opensource.org/licenses/MIT
#
# Created on 2014-11-06
#

"""benchmark_hashes.py file1 [file2 [...]]

Benchmark various hashes available in Python.
"""

from __future__ import print_function, unicode_literals, absolute_import

import logging
import sys
import timeit

LIMIT = 200  # Max no. of files to hash

HASH_FUNCS = {
    # name: (func, hexdigest?)
    'sha1': ('hashlib.sha1', 'True'),
    'md5': ('hashlib.md5', 'True'),
    'hash': ('hash', 'False'),
}

TIMER_TPL = """\
digest = {digest}
for f in files:
    with open(f) as fp:
        h = {func}(fp.read())
        if digest:
            h.hexdigest()
"""

logging.basicConfig(
    format='%(asctime)s %(levelname)8s %(filename)s:%(lineno)d %(message)s',
    datefmt='%H:%M:%S',
    level=logging.DEBUG)

log = logging.getLogger('')


def main():
    if '-h' in sys.argv or '--help' in sys.argv:
        print(__doc__)
        return 0
    if len(sys.argv) < 2:
        print(__doc__)
        return 1
    files = sys.argv[1:]

    if len(files) > LIMIT:
        files = files[:LIMIT]

    for name in HASH_FUNCS:
        func, digest = HASH_FUNCS[name]
        code = TIMER_TPL.format(func=func, digest=digest)
        t = timeit.Timer(
            code,
            'import hashlib; import sys; files=sys.argv[1:]')
        print('{:5s} : {}'.format(name, t.timeit(1)))


if __name__ == '__main__':
    sys.exit(main())
