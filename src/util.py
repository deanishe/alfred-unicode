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
Utility functions
"""

from __future__ import print_function, unicode_literals, absolute_import

import functools


def readable_time(seconds):
    """Convert seconds to hours, minutes, seconds"""

    def _mod(m, num):
        return (num / m, num % m)

    mod60 = functools.partial(_mod, 60)
    mod24 = functools.partial(_mod, 24)
    mod7 = functools.partial(_mod, 7)

    def _pluralise(num, unit):
        if num == 0:
            return ''
        suffix = unit + ('', 's')[num > 1]
        return '{:d} {:s}'.format(num, suffix)

    if seconds < 0.1:
        return '%0.3f secs' % seconds
    if seconds < 1:
        return '%0.2f secs' % seconds
    if seconds == 1:
        return '1 sec'
    if seconds < 10:
        return '%0.1f secs' % seconds
    if seconds < 60:
        return '%d secs' % seconds

    weeks = days = hours = minutes = 0
    minutes, seconds = mod60(seconds)
    hours, minutes = mod60(minutes)
    days, hours = mod24(hours)
    weeks, days = mod7(days)

    unit_map = [
        [weeks, 'week'],
        [days, 'day'],
        [hours, 'hr'],
        [minutes, 'min'],
        [seconds, 'sec'],
    ]
    # Drop leading entries that are zero
    while len(unit_map) and unit_map[0][0] < 1:
        unit_map.pop(0)
    # unit_map = [[i, u] for i, u in unit_map if i > 0]

    # Drop all but the largest two
    if len(unit_map) > 2:
        unit_map = unit_map[:2]

    # Round up if first unit if there are more than 2 units and
    # > 3 of the first
    if len(unit_map) > 1 and unit_map[0][0] > 3:
        # Round up to next day and drop hours if hour > 11
        if unit_map[0][1] == 'day' and unit_map[1][0] > 11:
            unit_map[0][0] += 1
            unit_map.pop()
        # Round up to next hour/minute if minute/second is > 29
        elif unit_map[1][0] > 29:
            unit_map[0][0] += 1
            unit_map.pop()

    return ', '.join([_pluralise(i, u) for i, u in unit_map if i > 0])


if __name__ == '__main__':
    for i in (120, 240, 3600, 84600, 12312,
              435, 57656, 2342, 2343425, 34534435,
              23432, 368997654, 25897, 343434, 7886907632,
              0.03, 0.7, 5.6, 2, 1.6, 90):
        print('{:10s}  ->  {}'.format(str(i), readable_time(i)))
