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
Generate list of all Unicode characters based on *flat* XML list
from http://www.unicode.org/Public/UCD/latest/ucdxml/ucd.nounihan.flat.zip

Saves results to `characters.tsv` and `characters.json` in the same
directory.
"""

from __future__ import print_function, unicode_literals, absolute_import

# from codecs import unicode_escape_decode
import json
import logging
import os
from time import time
import sys
import xml.sax

logging.basicConfig(
    format='%(asctime)s %(levelname)8s %(filename)s:%(lineno)d %(message)s',
    datefmt='%H:%M:%S',
    level=logging.DEBUG)

log = logging.getLogger('')

XML_PATH = os.path.join(os.path.dirname(__file__), 'ucd.nounihan.flat.xml')
ENTITY_LIST = os.path.join(os.path.dirname(__file__), 'entities.json')


def ignore(name):
    """Return True if character should be ignored"""
    bad_names = ('CJK COMPATIBILITY IDEOGRAPH-#',
                 'CJK UNIFIED IDEOGRAPH-#')
    if name in bad_names:
        return True
    return False


class UCDHandler(xml.sax.ContentHandler):

    def __init__(self):
        with open(ENTITY_LIST) as fp:
            self.entities = json.load(fp)
        self._stack = []
        self._aliases = []
        self._control = False
        self.count = 0
        self._characters = {}
        self.tsv_path = 'characters.tsv'
        self.json_path = 'characters.json'
        xml.sax.ContentHandler.__init__(self)

    def startDocument(self):
        self.start = time()

    def endDocument(self):
        duration = time() - self.start
        log.info('{:d} characters parsed in {:0.3f} seconds'.format(
                 len(self._characters), duration))

        with open(self.tsv_path, 'wb') as fp:
            for name, value in self._characters.items():
                num, entity = value
                msg = '{}\t{}\t'.format(name, num)
                if entity:
                    msg += entity
                msg += '\n'
                fp.write(msg.encode('utf-8'))

        with open(self.json_path, 'wb') as fp:
            json.dump(self._characters, fp, indent=2, sort_keys=True)

    def startElement(self, name, attrs):
        self._stack.append((name, attrs))

    def endElement(self, name):
        _, attrs = self._stack.pop()

        if name == 'name-alias':
            self._aliases.append(attrs['alias'])
            if attrs['type'] in ('control', 'figment'):
                self._control = True

        elif name == 'char':
            if self._control:
                log.debug('Ignoring control character : {}'.format(
                          self._aliases[0]))
            elif 'cp' not in attrs:
                log.warning('No codepoint : {}'.format(attrs['na']))
            else:
                if len(self._characters) % 1000 == 0:
                    log.debug('{:6d} characters in {:0.3f} seconds'.format(
                              len(self._characters), time() - self.start))
                num = attrs['cp']

                # Ignore control characters (<= 20)
                i = int(num, 16)
                if i > int('20', 16):
                    # s = '\\U{:0>8s}'.format(num)
                    # char = unicode_escape_decode(s)[0]
                    entity = self.entities.get(num)
                    name = attrs['na']
                    names = self._aliases[:]
                    if name:
                        names.append(name)
                    for name in names:
                        # if name == 'CJK UNIFIED IDEOGRAPH-#':
                        #     name = 'CJK UNIFIED IDEOGRAPH-{}'.format(num)

                        if not ignore(name):
                            if name in self._characters:
                                log.warning('Duplicate character : {}'.format(
                                            name))
                            self._characters[name] = (num, entity)
                        # print('{}\t{}'.format(name, num))

            # reset
            self._aliases = []
            self._control = False


def main():
    handler = UCDHandler()
    xml.sax.parse(XML_PATH, handler)


if __name__ == '__main__':
    sys.exit(main())
