//
//  IGSettings.m
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import "IGSettings.h"

@implementation GBSettings (IGSettings)

GB_SYNTHESIZE_OBJECT(NSString*, fontName, setFontName, @"fontname");
GB_SYNTHESIZE_OBJECT(NSString*, outputDirectory, setOutputDirectory, @"outputdir");
GB_SYNTHESIZE_OBJECT(NSArray*, codePoints, setCodePoints, @"codepoints");
GB_SYNTHESIZE_INT(limit, setLimit, @"limit");
GB_SYNTHESIZE_BOOL(verbose, setVerbose, @"verbose");
GB_SYNTHESIZE_OBJECT(NSString*, versionNumber, setVersionNumber, @"versionNumber");
GB_SYNTHESIZE_BOOL(overwrite, setOverwrite, @"overwrite");
GB_SYNTHESIZE_FLOAT(iconSize, setIconSize, @"iconsize");
GB_SYNTHESIZE_BOOL(printHelp, setPrintHelp, @"help");
GB_SYNTHESIZE_BOOL(printFontList, setPrintFontList, @"fontlist");
GB_SYNTHESIZE_BOOL(printVersion, setPrintVersion, @"version");

- (void)applyFactoryDefaults
{
    self.iconSize = 256.0f;
    self.fontName = @"ArialUnicodeMS";
    self.outputDirectory = [[[NSFileManager alloc] init] currentDirectoryPath];
    self.codePoints = [NSArray array];
    self.verbose = NO;
    self.overwrite = NO;
    self.limit = 0;
    self.printFontList = NO;
    self.printVersion = NO;
    self.versionNumber = [[NSBundle mainBundle] objectForInfoDictionaryKey: (NSString *) kCFBundleVersionKey];
}

@end
