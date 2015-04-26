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


def human_readable_time(seconds):
    """Convert seconds to hours, minutes, seconds"""

    def mod60(number):
        result = number / 60
        remainder = number % 60
        return result, remainder

    def mod24(number):
        result = number / 24
        remainder = number % 24
        return result, remainder

    def format_units(number, unit, plural=None):
        if number == 0:
            return ''
        if not plural:
            plural = unit + "s"
        if number == 1:
            return '%d %s' % (number, unit)
        else:
            return '%d %s' % (number, plural)

    if seconds < 1:
        return '%0.2f secs' % seconds
    if seconds == 1:
        return '1 sec'
    if seconds < 5:
        return '%0.1f secs' % seconds
    if seconds < 60:
        return '%d secs' % seconds

    days, hours, minutes = 0, 0, 0
    minutes, seconds = mod60(seconds)
    hours, minutes = mod60(minutes)
    days, hours = mod24(hours)
    unit_map = [(days, 'day'), (hours, 'hr'),
                (minutes, 'min'), (seconds, 'sec')]
    unit_map = [[i, u] for i, u in unit_map if i > 0]

    while len(unit_map) > 2:
        unit_map.pop()

    if len(unit_map) > 1 and unit_map[0][0] > 3:
        if unit_map[0][1] == 'day' and unit_map[1][0] > 11:
            unit_map[0][0] += 1
            unit_map.pop()
        elif unit_map[1][0] > 29:
            unit_map[0][0] += 1
            unit_map.pop()

    return ' '.join([format_units(i, u) for i, u in unit_map])
