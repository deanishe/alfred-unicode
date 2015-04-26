//
//  IGApp.h
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

@interface IGApp : NSObject
{
    NSFileManager *fileManager;
    NSUserDefaults *defaults;
    NSArray *codePoints;
}

@property NSFileManager *fileManager;
@property NSUserDefaults *defaults;
@property NSArray *codePoints;

- (int)run;

- (void)applyFactoryDefaults;

- (void)parseCommandLineArguments;

- (NSString *)stringWithCodePoint:(uint32_t)codePoint;

- (NSString *)resolvePath:(NSString *)path;

- (void)printHelp;

- (void)printFontList;

@end
