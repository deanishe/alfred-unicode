//
//  GBSettingsAdditions.m
//  UnicodeImageGen
//
//  Created by Dean Jackson on 09/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import "GBSettingsAdditions.h"

@implementation GBSettings (GBSettingsAdditions)

GB_SYNTHESIZE_OBJECT(NSString*, fontName, setFontName, @"fontname");
GB_SYNTHESIZE_OBJECT(NSString*, outputDirectory, setOutputDirectory, @"outputdir");
GB_SYNTHESIZE_INT(limit, setLimit, @"limit");
GB_SYNTHESIZE_BOOL(verbose, setVerbose, @"verbose");
GB_SYNTHESIZE_FLOAT(iconSize, setIconSize, @"iconsize");
GB_SYNTHESIZE_BOOL(printHelp, setPrintHelp, @"help");
GB_SYNTHESIZE_BOOL(printFontList, setPrintFontList, @"fontlist");

//+ (id)mySettingsWithName:(NSString *)name parent:(GBSettings *)parent {
//    id result = [self settingsWithName:name parent:parent];
//    if (result) {
//        [result registerArrayForKey:@"output"];
//    }
//    return result;
//}

- (void)applyFactoryDefaults
{
    self.iconSize = 256.0f;
    self.fontName = @"ArialUnicodeMS";
    self.outputDirectory = [[[NSFileManager alloc] init] currentDirectoryPath];
    self.verbose = NO;
    self.limit = 0;
    self.printFontList = NO;
}

@end
