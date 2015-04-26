//
//  IGIconGenerator.h
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h>

@interface IGIconGenerator : NSObject

{
    NSString* fontName;
    CGFloat iconSize;
    BOOL verbose;
}

@property NSString* fontName;
@property CGFloat iconSize;
@property BOOL verbose;

- (BOOL)saveIcon:(NSString *)filePath withString:(NSString *)string;

- (float)screenDPI;

@end
