//
//  NSDictionary+TSV.h
//  UnicodeImageGen
//
//  Created by Dean Jackson on 10/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>

@interface NSDictionary (tsv)

- (void)saveToFilePath: (NSString*)filePath error:(NSError*)error;

@end
