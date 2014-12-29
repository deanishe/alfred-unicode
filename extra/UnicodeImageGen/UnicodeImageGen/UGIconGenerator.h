//
//  UGIconGenerator.h
//  UnicodeImageGen
//
//  Created by Dean Jackson on 08/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

@interface UGIconGenerator : NSObject
{
    NSString* fontName;
    NSString* character;
    NSString* filePath;
    CGFloat iconSize;
    BOOL verbose;
}

@property NSString* fontName;
@property NSString* character;
@property NSString* filePath;
@property CGFloat iconSize;
@property BOOL verbose;

- (NSString*)savePreview;

- (float)screenDPI;

@end
