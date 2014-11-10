//
//  GBSettingsAdditions.h
//  UnicodeImageGen
//
//  Created by Dean Jackson on 09/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "GBSettings.h"

@interface GBSettings (GBSettingsAdditions)

@property (nonatomic, assign) NSString* fontName;
@property (nonatomic, assign) NSString* outputDirectory;
@property (nonatomic, assign) NSInteger limit;
@property (nonatomic, assign) BOOL verbose;
@property (nonatomic, assign) CGFloat iconSize;
@property (nonatomic, assign) BOOL printHelp;
@property (nonatomic, assign) BOOL printFontList;

//+ (id)mySettingsWithName:(NSString *)name parent:(GBSettings *)parent;

- (void)applyFactoryDefaults;

@end
