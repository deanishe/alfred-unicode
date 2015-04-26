//
//  IGSettings.h
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <GBCli/GBSettings.h>

@interface GBSettings (IGSettings)

@property (nonatomic, assign) NSString* fontName;
@property (nonatomic, assign) NSString* outputDirectory;
@property (nonatomic, assign) NSArray* codePoints;
@property (nonatomic, assign) NSInteger limit;
@property (nonatomic, assign) CGFloat iconSize;
@property (nonatomic, assign) BOOL verbose;
@property (nonatomic, assign) NSString* versionNumber;
@property (nonatomic, assign) NSString* buildNumber;
@property (nonatomic, assign) NSString* appName;
@property (nonatomic, assign) BOOL overwrite;
@property (nonatomic, assign) BOOL printHelp;
@property (nonatomic, assign) BOOL printFontList;
@property (nonatomic, assign) BOOL printVersion;

- (void)applyFactoryDefaults;

@end
