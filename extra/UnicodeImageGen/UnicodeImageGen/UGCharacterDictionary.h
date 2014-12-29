//
//  UGCharacterDictionary.h
//  UnicodeImageGen
//
//  Created by Dean Jackson on 09/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface UGCharacterDictionary : NSObject
{
    @protected
    NSMutableDictionary *_characters;
}

@ property (readonly) NSDictionary * characters;

- (id) init;

@end
