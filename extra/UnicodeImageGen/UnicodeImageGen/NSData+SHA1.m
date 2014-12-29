//
//  NSData+SHA1.m
//  UnicodeImageGen
//
//  Created by Dean Jackson on 10/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import <CommonCrypto/CommonDigest.h>
#import "NSData+SHA1.h"

@implementation NSData (sha1)

- (NSString*)sha1
{
    uint8_t digest[CC_SHA1_DIGEST_LENGTH];

    CC_SHA1(self.bytes, (unsigned int)self.length, digest);

    NSMutableString *hash = [NSMutableString
                             stringWithCapacity:CC_SHA1_DIGEST_LENGTH * 2];

    for (int i=0; i < CC_SHA1_DIGEST_LENGTH; i++) {
        [hash appendFormat:@"%02x", digest[i]];
    }
    return hash;
}

@end
