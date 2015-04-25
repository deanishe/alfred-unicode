#!/usr/bin/env python
# encoding: utf-8
#
# Copyright © 2014 deanishe@deanishe.net
#
# MIT Licence. See http://opensource.org/licenses/MIT
#
# Created on 2014-11-09
#

"""
"""

from __future__ import print_function, unicode_literals, absolute_import

import csv
import logging
import os
import sys

LOG_LEVEL = logging.INFO

logging.basicConfig(
    format='%(asctime)s %(levelname)8s %(filename)s:%(lineno)d %(message)s',
    datefmt='%H:%M:%S',
    level=LOG_LEVEL)

log = logging.getLogger('')


PARENT_DIR = os.path.abspath(os.path.dirname(__file__))

TSV_FILE = os.path.join(PARENT_DIR, 'characters.tsv')
OBJC_FILE = os.path.join(PARENT_DIR,
                         'UnicodeImageGen',
                         'UnicodeImageGen',
                         'UGCharacterDictionary.m')

COCOA_TPL = b"""

//
//  CharacterArray.m
//  UnicodeImageGen
//
//  Created by Dean Jackson on 09/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import "UGCharacterDictionary.h"

@implementation UGCharacterDictionary

@synthesize characters = _characters;

- (id) init
{
    if (self = [super init]) {
        _characters = [NSMutableDictionary dictionary];
%(content)s

    }
    return self;
}

@end



"""

ENTRY_TPL = b"""\
        [_characters setValue:@"{char}" forKey:@"{codepoint}"];"""

ESCAPED = [
    '2028',
    '2029',
]


def main():
    entries = []
    with open(TSV_FILE, 'rb') as fp:
        reader = csv.reader(fp, delimiter=b'\t')
        for row in reader:
            name, h, entity = [v for v in row]
            codepoint = b'{:0>8s}'.format(h)
            i = int(codepoint, 16)
            if i < 128:
                char = unichr(i).encode('utf-8')
                if char == b'\\':
                    char = b'\\\\'
                elif char == b'"':
                    char = b'\\"'
            else:
                char = b'\U' + codepoint
            # s = '\\U{:0>8s}'.format(h)
            # char = unicode_escape_decode(s)[0]
            # if char == '\\':
            #     char = '\\\\'
            # elif char == '"':
            #     char = '\\"'
            # elif codepoint in ESCAPED:
            #     char = '\\\\\\\\U{:0>8s}'.format(h)
            entries.append(ENTRY_TPL.format(
                           codepoint=codepoint,
                           char=char))
    content = '\n'.join(entries)
    text = COCOA_TPL % {'content': content}
    with open(OBJC_FILE, 'wb') as fp:
        fp.write(text)


if __name__ == '__main__':
    sys.exit(main())
