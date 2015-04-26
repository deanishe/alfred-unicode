//
//  UGApp.h
//  UnicodeImageGen
//
//  Created by Dean Jackson on 08/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "GBSettingsAdditions.h"
#import "UGCharacterDictionary.h"

@interface UGApp : NSObject
{
    GBSettings* settings;
    NSDictionary* characters;
    NSFileManager* fileManager;

}

@property GBSettings* settings;
@property (readonly) NSDictionary* characters;
@property NSFileManager* fileManager;

- (int)run;

//- (NSArray *)readCharactersFile;

- (NSString *)filePathForCodePoint:(NSString *)codePoint;

- (NSString *)stringWithCodePoint:(uint32_t)codePoint;

- (NSString *)resolvePath:(NSString *)path;

- (void)printFontList;

@end
