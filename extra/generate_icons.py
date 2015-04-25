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
"""

from __future__ import print_function, unicode_literals, absolute_import

from codecs import unicode_escape_decode
from collections import defaultdict
import csv
import hashlib
import json
import logging
from multiprocessing import Pool
import os
import subprocess
import sys
from time import time

from PIL import Image, ImageFont, ImageDraw

LOG_LEVEL = logging.INFO

logging.basicConfig(
    format='%(asctime)s %(levelname)8s %(filename)s:%(lineno)d %(message)s',
    datefmt='%H:%M:%S',
    level=LOG_LEVEL)

log = logging.getLogger('')


PARENT_DIR = os.path.abspath(os.path.dirname(__file__))

TSV_FILE = os.path.join(PARENT_DIR, 'characters.tsv')
ICON_TSV_FILE = os.path.join(PARENT_DIR, 'icons.tsv')
ICON_JSON_FILE = os.path.join(PARENT_DIR, 'icons.json')

FONT_DIR = os.path.join(PARENT_DIR, 'fonts')

# FONT_FILE = os.path.join(PARENT_DIR, 'fonts', 'Universalia+.ttf')
# FONT_FILE = os.path.join(PARENT_DIR, 'fonts', 'u0000.ttf')
# FONT_FILE = os.path.join(PARENT_DIR, 'fonts', 'NotoSans-Regular.ttf')
FONT_FILE = '/Library/Fonts/Arial Unicode.ttf'


ICON_DIR = os.path.join(os.path.dirname(PARENT_DIR),
                        'src', 'icons')

UNKNOWN_CHAR = 'fffd'
UNKNOWN_CHAR_ICON_PATH = os.path.join(ICON_DIR, 'unknown.png')

ICON_SIZE = 256
FONT_SIZE = int(ICON_SIZE * 0.7)
BG_COLOUR = (0, 0, 0, 0)
COLOUR = '#D100A6'
# For debugging
LIMIT = 0
GRAPHICS_MAGICK = '/usr/local/bin/gm'
PNGCRUSH = '/usr/local/bin/pngcrush'

FONTS = [f for f in os.listdir(FONT_DIR) if f.startswith('u')]
FONTS = [(int(f[1:-4], 16), f) for f in FONTS]
FONTS = [(i, os.path.join(FONT_DIR, f)) for i, f in FONTS]
FONTS.sort()


def get_font(h):
    return FONT_FILE
    n = int(h, 16)
    for i, (j, path) in enumerate(FONTS):
        if n < j:
            i -= 1
            break

    return FONTS[i][1]


def get_image_path(h):
    """Generate an image path from the Unicode codepoint"""
    filename = '{:0>8}.png'.format(h.upper())
    dirpath = os.path.join(
        ICON_DIR,
        filename[:2],
        filename[2:4],
        filename[4:6],
    )
    filepath = os.path.join(dirpath, filename)
    return filepath


# def save_image_gm(char, imagefile, fontpath):
#     cmd = [
#         GRAPHICS_MAGICK,
#         'convert',
#         '-pointsize', unicode(ICON_SIZE),
#         '-font', fontpath,
#         'label:{}'.format(char),
#         imagefile,
#     ]
#     cmd = [u.encode('utf-8') for u in cmd]
#     log.debug(cmd)
#     subprocess.call(cmd)
#     return pngcrush(imagefile)


# def save_image_old(char, imagefile, fontpath):
#     image = Image.new('RGBA', (ICON_SIZE, ICON_SIZE),
#                       color=BG_COLOUR)
#     draw = ImageDraw.Draw(image)
#     font = ImageFont.truetype(fontpath, int(ICON_SIZE * 0.8))
#     draw.text((0, 0), char, font=font, fill=COLOUR)
#     draw = ImageDraw.Draw(image)
#     image.save(imagefile)
#     log.info('Saved `{}`'.format(imagefile))


def save_image(char, imagefile, fontpath):
    """Save icon `char` from `fontpath` to `imagefile`. Return SHA1 digest"""
    image = Image.new('RGBA', (ICON_SIZE, ICON_SIZE),
                      color=BG_COLOUR)
    draw = ImageDraw.Draw(image)

    font = ImageFont.truetype(fontpath, FONT_SIZE)

    size = draw.textsize(char, font=font)
    offset = font.getoffset(char)

    width, height = map(sum, zip(size, offset))

    draw.text(((ICON_SIZE - width) / 2, (ICON_SIZE - height) / 2),
              char,
              font=font,
              fill=COLOUR)

    # Get bounding box
    bbox = image.getbbox()

    # Create alpha mask
    imagemask = Image.new("L", (ICON_SIZE, ICON_SIZE), 0)
    drawmask = ImageDraw.Draw(imagemask)

    # Draw icon on mask
    drawmask.text(((ICON_SIZE - width) / 2,
                  (ICON_SIZE - height) / 2),
                  char,
                  font=font,
                  fill=255)

    iconimage = Image.new('RGBA', (ICON_SIZE, ICON_SIZE), COLOUR)
    iconimage.putalpha(imagemask)

    if bbox:
        iconimage = iconimage.crop(bbox)

        borderw = int((ICON_SIZE - (bbox[2] - bbox[0])) / 2)
        borderh = int((ICON_SIZE - (bbox[3] - bbox[1])) / 2)

    else:
        borderw = borderh = 0

    # Create output image
    outimage = Image.new('RGBA', (ICON_SIZE, ICON_SIZE), BG_COLOUR)
    outimage.paste(iconimage, (borderw, borderh))

    # Save file
    # name, ext = os.path.splitext(imagefile)
    # imagefile = '{}-pil{}'.format(name, ext)
    outimage.save(imagefile)
    return pngcrush(imagefile)


def pngcrush(filepath):
    """Run file through `pngcrush` and return SHA1 hash"""
    name, ext = os.path.splitext(filepath)
    temppath = '{}.{}{}'.format(name, os.getpid(), ext)
    size_in = os.stat(filepath).st_size
    os.rename(filepath, temppath)
    cmd = [
        PNGCRUSH,
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
    log.debug('Optimised [{}b / {:0.1f}%] `{}`'.format(size_out, pc, filepath))
    with open(filepath) as fp:
        return hashlib.sha1(fp.read()).hexdigest()


def handle_duplicates(hash_file_map, hash_hex_map):
    """Delete duplicates and link them to the unknown icon instead"""
    hex_file_map = {}
    dupes = 0
    count = 0
    empty_dirs = 0
    for hsh in hash_file_map:
        files = hash_file_map[hsh]
        count += len(files)
        if len(files) == 1:
            hashes = hash_hex_map[hsh]
            assert len(hashes) == 1
            hex_file_map[hashes[0]] = files[0]
            continue
        log.debug('{} files for hash `{}`'.format(len(files), hsh))
        dupes += len(files)

        for h in hash_hex_map.get(hsh, []):
            hex_file_map[h] = UNKNOWN_CHAR_ICON_PATH

        for path in files:
            if os.path.exists(path):
                os.unlink(path)

    # Delete empty directories
    for root, dirnames, filenames in os.walk(ICON_DIR, topdown=False):
        for dirname in dirnames:
            path = os.path.join(root, dirname)
            if len(os.listdir(path)) == 0:
                log.debug('Deleting empty directory `{}`'.format(path))
                os.rmdir(path)
                empty_dirs += 1

    pc = (float(dupes) / count) * 100
    log.info('Deleted {} duplicates ({:0.1f}%) and {} empty directories'.format(
             dupes, pc, empty_dirs))

    # Make paths relative
    for h in hex_file_map:
        hex_file_map[h] = os.path.relpath(hex_file_map[h], ICON_DIR)

    return hex_file_map


def pool_wrapper(h, imagefile, fontpath):
    """multiprocessing wrapper for `save_image()`"""
    s = '\\U{:0>8s}'.format(h)
    char = unicode_escape_decode(s)[0]
    try:
        hsh = save_image(char, imagefile, fontpath)
    except Exception as err:
        log.error(err)
        hsh = None
    return (h, hsh, imagefile, fontpath)


def main():
    hash_file_map = defaultdict(list)
    hash_hex_map = defaultdict(list)
    fontless = []
    start = time()

    assert os.path.exists(FONT_FILE), 'FONT_FILE does not exist'

    # Create unknown file first
    s = '\\u{}'.format(UNKNOWN_CHAR)
    char = unicode_escape_decode(s)[0]
    save_image(char, UNKNOWN_CHAR_ICON_PATH, FONT_FILE)

    with open(TSV_FILE, 'rb') as fp:
        reader = csv.reader(fp, delimiter=b'\t')
        characters = []
        for row in reader:
            name, h, entity = [v.decode('utf-8') for v in row]
            filepath = get_image_path(h)
            characters.append((name, h, filepath))

        log.debug('{} characters in `{}`'.format(len(characters), TSV_FILE))

    total = len(characters)
    pool = Pool()
    results = []
    done = 0
    for i, (name, h, filepath) in enumerate(characters):
        if LIMIT and i == LIMIT:
            total = LIMIT
            break

        fontpath = get_font(h)
        if fontpath is None:
            log.error('No font found for U+{}'.format(h))
            fontless.append(h)
            continue

        log.debug('{} -> {}'.format(h, fontpath))

        # filepath = '{}.png'.format(h.upper())
        # s = '\\u{}'.format('26C5')
        # s = '\\u{}'.format(h)
        # char = unicode_escape_decode(s)[0]
        # char = '⛅'
        # log.debug(char)
        imagedir = os.path.dirname(filepath)
        if not os.path.exists(imagedir):
            os.makedirs(imagedir)
        # hsh = save_image(char, filepath, FONT_FILE)
        r = pool.apply_async(pool_wrapper, (h, filepath, fontpath))
        results.append(r)

    # Handle results
    while done < total:
        if not len(results):
            continue
        r = results.pop(0)
        if not r.ready():
            results.append(r)
            continue

        h, hsh, filepath, fontpath = r.get()
        done += 1

        pc = (float(done) / total) * 100
        log.info('[{:5d}/{:d}] {:6.2f}% : Saved `{}` ({})'.format(
                 done, total, pc, filepath.replace(ICON_DIR, '').lstrip('/'),
                 os.path.basename(fontpath)))

        if hsh is None:  # Failed
            log.error('Failed : {}'.format(h))
            continue

        hash_file_map[hsh].append(filepath)
        hash_hex_map[hsh].append(h)

    duration = time() - start
    av = duration / total
    log.info('{:d} icons generated in {:0.2f} seconds ({:0.2f}s/icon)'.format(
             done, duration, av))

    # Remove duplicates and link them to the unknown icon
    hex_files_map = handle_duplicates(hash_file_map, hash_hex_map)

    if len(fontless):
        log.warning('{} characters with no font'.format(len(fontless)))

    for h in fontless:
        hex_files_map[h] = UNKNOWN_CHAR_ICON_PATH

    with open(ICON_TSV_FILE, 'wb') as fp:
        for h in sorted(hex_files_map.keys()):
            fp.write('{}\t{}\n'.format(h, hex_files_map[h]))

    with open(ICON_JSON_FILE, 'wb') as fp:
        json.dump(hex_files_map, fp, sort_keys=True, indent=2)


if __name__ == '__main__':
    sys.exit(main())
