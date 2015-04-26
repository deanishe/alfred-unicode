//
//  IGLogFileManager.h
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import <CocoaLumberjack/CocoaLumberjack.h>
#import "DDFileLogger.h"

@interface IGLogFileManager : DDLogFileManagerDefault
{
    NSString *logFile;
}

@property NSString *logFile;

- (instancetype)initWithLogFile:(NSString *)logFile;
- (BOOL)isLogFile:(NSString *)fileName;
- (NSString *)newLogFileName;
- (NSString *)createNewLogFile;

@end
