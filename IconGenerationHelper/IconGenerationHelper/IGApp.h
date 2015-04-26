//
//  IGApp.h
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>
#import "IGSettings.h"
#import "IGIconGenerator.h"
#import "IGLogFileManager.h"
#import "IGLogFormatter.h"

@interface IGApp : NSObject
{
    NSArray *codePoints;
    NSFileManager *fileManager;
    IGIconGenerator *iconGenerator;
    NSString *outputDirectory;
    NSString *logFile;
    GBSettings *settings;
}

// NSString hexadecimal values of Unicode codepoints
@property NSArray *codePoints;
@property NSFileManager *fileManager;
@property IGIconGenerator *iconGenerator;
// Absolute path to outputDirectory
@property NSString *outputDirectory;
@property NSString *logFile;
@property GBSettings *settings;

- (int)runWithSettings:(GBSettings *)_settings;

// Return codePoint as 32-bit integer
- (uint32_t)integerWithCodePoint:(NSString *)codePoint;

// Return Unicode string corresponding to codePoint
- (NSString *)stringWithCodePoint:(NSString *)codePoint;

// Return valid codePoint or nil
- (NSString *)validateCodePoint:(NSString *)codePoint;

// Return absolute path for path. Also resolve symlinks
- (NSString *)resolvePath:(NSString *)path;

// Return appropriate filePath (under outputDirectory) for icon for codePoint
- (NSString *)filePathWithCodePoint:(NSString *)codePoint;

// Print a list of available fonts to STDOUT
- (void)printFontList;


@end
