//
//  main.m
//  UnicodeImageGen
//
//  Created by Dean Jackson on 08/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <CoreFoundation/CoreFoundation.h>
#import <stdlib.h>
#import <stdio.h>
#import "GBCli.h"
#import "UGApp.h"
#import "GBSettingsAdditions.h"

// From http://www.cocoabuilder.com/archive/cocoa/193451-finding-out-executable-location-from-c-program.html
char *GetExecutableLocation() {
    CFBundleRef bundle = CFBundleGetMainBundle();
    CFURLRef executableURL = CFBundleCopyExecutableURL(bundle);
    CFStringRef executablePath = CFURLCopyFileSystemPath(executableURL, kCFURLPOSIXPathStyle);
    CFIndex maxLength = CFStringGetMaximumSizeOfFileSystemRepresentation(executablePath);
    char *result = malloc(maxLength);

    if(result) {
        if(!CFStringGetFileSystemRepresentation(executablePath, result, maxLength)) {
            free(result);
            result = NULL;
        }
    }

    CFRelease(executablePath);
    CFRelease(executableURL);
    
    return result;
}

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        // Command-line arguments
        GBSettings *factoryDefaults = [GBSettings settingsWithName:@"Factory" parent:nil];
        GBSettings *settings = [GBSettings settingsWithName:@"CmdLine" parent:factoryDefaults];
        [factoryDefaults applyFactoryDefaults];

        // Create options helper and register options
        GBOptionsHelper *options = [[GBOptionsHelper alloc] init];
        options.printHelpHeader = ^{
            return @"UnicodeImageGen [-h] [-v] [-f <fontname>] [-s <fontsize>] [-o <outputdir>] [-i <iconpaths.tsv>]\n\nGenerate preview icons for Unicode characters\n"; };
        [options registerOption:'f'
                           long:@"fontname"
                    description:@"Name of font"
                          flags:(GBOptionFlags)GBValueOptional];
        [options registerOption:'s'
                           long:@"iconsize"
                    description:@"Size of icon"
                          flags:(GBOptionFlags)GBValueOptional];
        [options registerOption:'o'
                           long:@"outputdir"
                    description:@"Directory to save icons to"
                          flags:(GBOptionFlags)GBValueOptional];
        [options registerOption:'i'
                           long:@"iconlist"
                    description:@"TSV file to save codepoint-iconpaths map to"
                          flags:(GBOptionFlags)GBValueOptional];
        [options registerOption:'l'
                           long:@"limit"
                    description:@"Max. number of icons to save (for debugging)"
                          flags:(GBOptionFlags)GBValueOptional];
        [options registerOption:'v'
                           long:@"verbose"
                    description:@"Emit much verbosity"
                          flags:(GBOptionFlags)GBValueOptional];
        [options registerOption:'F'
                           long:@"fontlist"
                    description:@"Show list of available fonts"
                          flags:(GBOptionFlags)GBValueNone];
        [options registerOption:'h'
                           long:@"help"
                    description:@"Display this message"
                          flags:(GBOptionFlags)GBValueNone];

        // Options parser
        GBCommandLineParser *parser = [[GBCommandLineParser alloc] init];
        [parser registerSettings:settings];
        [parser registerOptions:options];

        [parser parseOptionsUsingDefaultArguments];

        if (settings.verbose) {
            [options printValuesFromSettings:settings];
        }

        if (argc == 1 || settings.printHelp) {
            [options printHelp];
            return EXIT_SUCCESS;
        }

        if ([settings verbose]) {
            char *path = GetExecutableLocation();
            NSString *appDirectory = [[NSString stringWithUTF8String:path] stringByDeletingLastPathComponent];
            free(path);
            NSLog(@"Application directory : %@", appDirectory);
        }

        UGApp *app = [[UGApp alloc] init];
        app.settings = settings;
        [app run];
    }
    return EXIT_SUCCESS;
}
