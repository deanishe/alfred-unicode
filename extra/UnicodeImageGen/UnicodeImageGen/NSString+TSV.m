//
//  NSStringTSV.m
//  UnicodeImageGen
//
//  Created by Dean Jackson on 09/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import "NSString+TSV.h"

@implementation NSString (tsv)

- (NSArray *) arrayFromTSV
{
    NSMutableArray *result = [NSMutableArray array];
    NSArray *lines = [self componentsSeparatedByString:@"\n"];

    for (NSString *line in lines) {
        if ([line isEqualToString:@""]) {
            continue;
        }
        NSArray *values = [line componentsSeparatedByString:@"\t"];
        [result addObject:values];

    }
    return result;
}

@end
 