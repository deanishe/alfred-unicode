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
Save JSON and TSV lists of HTML named entities.
Source is file `HTML named entities.html`, which is a fragment
of the Wikipedia XML and HTML character entity references page:
http://en.wikipedia.org/wiki/List_of_XML_and_HTML_character_entity_references
"""

from __future__ import print_function, unicode_literals, absolute_import

import logging
import json
import os
import re
import sys


logging.basicConfig(
    format='%(asctime)s %(levelname)8s %(filename)s:%(lineno)d %(message)s',
    datefmt='%H:%M:%S',
    level=logging.DEBUG)

log = logging.getLogger('')


HTML_FILE = os.path.join(os.path.dirname(__file__), 'HTML named entities.html')

get_rows = re.compile(r'<tr.*?>(.+?)</tr>', re.DOTALL).findall
get_cells = re.compile(r'<td.*?>(.*)</td>').findall
parse_code = re.compile(r'U\+([0123456789ABCDEF]+) \((\d+)\)').match


def main():
    entities = {}
    with open(HTML_FILE) as fp:
        html = fp.read()
    rows = get_rows(html)
    log.debug('{:d} rows'.format(len(rows)))
    for row in rows:
        cells = get_cells(row)
        if not cells:
            log.debug('No cells : {}'.format(row))
            continue
        if len(cells) != 7:
            log.debug('Wrong cell count {}'.format(cells))
            continue
        name, char, code = cells[:3]
        h, d = parse_code(code).groups()
        # print('{}\t{}\t{}'.format(name, h, d))
        entities[h] = name

    with open('entities.tsv', 'wb') as fp:
        for h in sorted(entities.keys()):
            fp.write('{}\t{}\n'.format(h, entities[h]))

    with open('entities.json', 'wb') as fp:
        json.dump(entities, fp, indent=2, sort_keys=True)


if __name__ == '__main__':
    sys.exit(main())
