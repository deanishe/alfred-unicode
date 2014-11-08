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
Generate an HTML preview of Unicode characters.

Output file is `preview.html` in the same directory.
"""

from __future__ import print_function, unicode_literals, absolute_import

from codecs import unicode_escape_decode
import csv
import logging
import os
import sys
from time import time

logging.basicConfig(
    format='%(asctime)s %(levelname)8s %(filename)s:%(lineno)d %(message)s',
    datefmt='%H:%M:%S',
    level=logging.DEBUG)

log = logging.getLogger('')


TSV_FILE = os.path.join(os.path.dirname(__file__), 'characters.tsv')
HTML_FILE = os.path.join(os.path.dirname(__file__), 'preview.html')
# For debugging
LIMIT = 1000

STYLE_TPL = """\
@font-face {
    font-family: 'universalia';
    src: url("Universalia+.ttf") format("truetype");
    font-weight: normal;
    font-style: normal;
}
ul#icons {
    list-style: none;
}
ul#icons li {
    float: left;
    height: 120px;
    width: 120px;
    border: 1px solid #ccc;
    text-align: center;
    padding: 20px;
}
ul#icons li div.char {
    font-family: universalia;
    font-size: 48px;
    height: 80px;
}
ul#icons li div.info {
    font-family: 'Helvetica Neue';
    font-size: 18px;
    color: #444;
}
"""

PAGE_TPL = """\
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <title>Preview</title>
    <style>
{style}
    </style>
</head>
<body>
    <ul id="icons">
        {content}
    </ul>
</body>
</html>
"""

ICON_TPL = """\
<li>
    <div class="char">{char}</div>
    <div class="info">{info}</div>
</li>
"""


def main():
    start = time()
    count = 0
    items = []
    with open(TSV_FILE, 'rb') as fp:
        reader = csv.reader(fp, delimiter=b'\t')
        for row in reader:
            name, h, entity = [v.decode('utf-8') for v in row]
            s = '\\u{}'.format(h)
            u = unicode_escape_decode(s)[0]
            # char = '❤'
            log.info(u)
            count += 1
            items.append(ICON_TPL.format(char=u, info='U+{}'.format(h)))
            if LIMIT and count == LIMIT:
                break
    html = PAGE_TPL.format(content='\n'.join(items), style=STYLE_TPL)

    with open(HTML_FILE, 'wb') as fp:
        fp.write(html.encode('utf-8'))

    log.info('{:d} icons generated in {:0.2f} seconds'.format(
             count, time() - start))


if __name__ == '__main__':
    sys.exit(main())
