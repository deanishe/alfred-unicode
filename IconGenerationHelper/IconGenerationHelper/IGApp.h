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

@interface IGApp : NSObject
{
    NSArray *codePoints;
    NSFileManager *fileManager;
    IGIconGenerator *iconGenerator;
    NSString *outputDirectory;
    GBSettings *settings;
}

// NSString hexadecimal values of Unicode codepoints
@property NSArray *codePoints;
@property NSFileManager *fileManager;
@property IGIconGenerator *iconGenerator;
// Absolute path to outputDirectory
@property NSString *outputDirectory;
@property GBSettings *settings;

- (int)runWithSettings:(GBSettings *)_settings;

// Return codePoint as 32-bit integer
- (uint32_t)integerWithCodePoint:(NSString *)codePoint;

// Return Unicode string corresponding to codePoint
- (NSString *)stringWithCodePoint:(NSString *)codePoint;

- (NSString *)validateCodePoint:(NSString *)codePoint;

- (NSString *)resolvePath:(NSString *)path;

- (NSString *)filePathWithCodePoint:(NSString *)codePoint;

- (void)printFontList;

- (void)printVersion;

@end
