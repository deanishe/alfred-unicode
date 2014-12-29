//
//  NSDictionary+TSV.m
//  UnicodeImageGen
//
//  Created by Dean Jackson on 10/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import "NSDictionary+TSV.h"

@implementation NSDictionary (tsv)

- (void)saveToFilePath: (NSString*)filePath error:(NSError*)error
{
    NSMutableString *output = [NSMutableString string];

    for (NSString *key in self) {
        NSString *value = [self objectForKey: key];
        [output appendFormat: @"%@\t%@\n", key, value];
    }
    [output writeToFile:filePath
             atomically:YES
               encoding:NSUTF8StringEncoding
                  error:&error];
}

@end
