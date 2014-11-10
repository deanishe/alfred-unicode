//
//  PreviewApp.h
//  UnicodeImageGen
//
//  Created by Dean Jackson on 08/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import "GBSettingsAdditions.h"
#import "CharacterDictionary.h"

@interface PreviewApp : NSObject
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

- (void)printFontList;

@end
