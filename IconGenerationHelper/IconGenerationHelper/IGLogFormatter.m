//
//  IGLogFormatter.m
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import "IGLogFormatter.h"

@implementation IGLogFormatter

- (NSString *)formatLogMessage:(DDLogMessage *)logMessage
{
    NSDateFormatter *dateFormatter = [[NSDateFormatter alloc]init];
    [dateFormatter setDateFormat:@"HH:mm:ss"];
    NSString *date = [dateFormatter stringFromDate:(logMessage->_timestamp)];
    NSString *fileName = [NSString stringWithFormat:@"%@:%lu",
                          [logMessage->_file lastPathComponent],
                          logMessage->_line];
    NSString *logLevel;
    switch (logMessage->_flag) {
        case DDLogFlagError    : logLevel = @"ERROR  "; break;
        case DDLogFlagWarning  : logLevel = @"WARNING"; break;
        case DDLogFlagInfo     : logLevel = @"INFO   "; break;
        case DDLogFlagDebug    : logLevel = @"DEBUG  "; break;
        default                : logLevel = @""; break;
    }

    return [NSString stringWithFormat:@"%@ %@ %@ %@\n", date, fileName, logLevel, logMessage->_message];
}

@end
