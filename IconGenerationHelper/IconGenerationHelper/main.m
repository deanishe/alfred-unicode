//
//  main.m
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <CoreFoundation/CoreFoundation.h>
#import "IGApp.h"

int main(int argc, const char * argv[]) {
    int retcode;
    @autoreleasepool {
        IGApp *app = [[IGApp alloc] init];
        retcode = [app run];
    }
    return retcode;
}
