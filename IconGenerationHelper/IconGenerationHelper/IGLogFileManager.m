//
//  IGLogFileManager.m
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import "IGLogFileManager.h"

@implementation IGLogFileManager

@synthesize logFile;

- (instancetype)initWithLogFile:(NSString *)filePath {
    logFile = [filePath lastPathComponent];
    NSString *dirPath = [filePath stringByDeletingLastPathComponent];
    return [self initWithLogsDirectory:dirPath];
}

- (BOOL)isLogFile:(NSString *)fileName
{
    return [fileName hasSuffix:@".log"];
}

- (NSString *)newLogFileName
{
    return logFile;
}

- (NSString *)createNewLogFile
{
    NSString *fileName = [self newLogFileName];
    NSString *logDirectory = [self logsDirectory];
    NSString *filePath = [logDirectory stringByAppendingPathComponent:fileName];
    NSFileManager *fileManager = [NSFileManager defaultManager];
    if (! [fileManager fileExistsAtPath:filePath]) {
        [fileManager createFileAtPath:filePath contents:nil attributes:nil];
    }
    return filePath;
}


@end
