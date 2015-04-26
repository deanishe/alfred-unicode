//
//  main.m
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import <Foundation/Foundation.h>
#import <CoreFoundation/CoreFoundation.h>
#import <GBCli/GBCli.h>
#import <CocoaLumberjack/CocoaLumberjack.h>
#import "IGSettings.h"
#import "IGApp.h"

int main(int argc, char * argv[]) {
    int retcode;
    @autoreleasepool {

        // Default settings
        GBSettings *defaults = [GBSettings settingsWithName:@"Default" parent:nil];
        GBSettings *settings = [GBSettings settingsWithName:@"CLI" parent:defaults];
        [defaults applyFactoryDefaults];

        GBOptionsHelper *options = [[GBOptionsHelper alloc] init];
        options.applicationName = ^{ return [defaults appName]; };
        options.applicationVersion = ^{ return [defaults versionNumber]; };
        options.applicationBuild = ^{ return [defaults buildNumber]; };
        options.printHelpHeader = ^{
            return @"%APPNAME v%APPVERSION (build %APPBUILD)\n\n"
                   @"Generate icons for specified codepoints in output directory.\n"
                   @"\n"
                   @"Usage:\n"
                   @"\n"
                   @"    IconGen [-v] [-F] [-O] [-s <size>] [-l logfile] [-o <dirpath>] [-f <font>] [-l <limit>] [<codepoint> [<codepoint> ...]]\n"
                   @"\n"
                   @"If no codepoints are specified, they will be read from STDIN, one per line.\n"
                   @"If no output directory (-o) is specified, the current working directory will be used.\n"
                   @"If no logfile is specified, none will be written.\n"
                   @"\n"
                   @"Options:\n";

        };

        // Create parser and register all options.

        // Alternate actions
        [options registerOption:'F'
                           long:@"fontlist"
                    description:@"Print list of available fonts and exit."
                          flags:(GBOptionFlags)GBValueNone];
        [options registerOption:'h'
                           long:@"help"
                    description:@"Show this message and exit."
                          flags:(GBOptionFlags)GBValueNone];
        [options registerOption:'V'
                           long:@"version"
                    description:@"Show version number and exit."
                          flags:(GBOptionFlags)GBValueNone];

        // General options
        [options registerOption:'v'
                           long:@"verbose"
                    description:@"Log debugging information."
                          flags:(GBOptionFlags)GBValueNone];
        [options registerOption:'L'
                           long:@"logfile"
                    description:@"Log output to this file."
                          flags:(GBOptionFlags)GBValueRequired];

        // Icon-generation options
        [options registerOption:'l'
                           long:@"limit"
                    description:@"Only generate this many icons (for testing)."
                          flags:(GBOptionFlags)GBValueRequired];
        [options registerOption:'s'
                           long:@"size"
                    description:@"Generate icons of this size in pixels (default: 256)."
                          flags:(GBOptionFlags)GBValueRequired];
        [options registerOption:'f'
                           long:@"font"
                    description:@"Font to use to generate icons (default: 'ArialUnicodeMS')."
                          flags:(GBOptionFlags)GBValueRequired];
        [options registerOption:'o'
                           long:@"outputdir"
                    description:@"Directory to save icons in (default: current working directory)."
                          flags:(GBOptionFlags)GBValueRequired];
        [options registerOption:'O'
                           long:@"overwrite"
                    description:@"Overwrite existing icons (default: no). Useful if you've changed the colour."
                          flags:(GBOptionFlags)GBValueNone];

        // Register settings and parse CLI args
        GBCommandLineParser *parser = [[GBCommandLineParser alloc] init];
        [parser registerSettings:settings];
        [parser registerOptions:options];
        [parser parseOptionsUsingDefaultArguments];

        if (settings.verbose) {
            [options printValuesFromSettings:settings];
            [DDLog addLogger:[DDASLLogger sharedInstance] withLevel:DDLogLevelDebug];
            [DDLog addLogger:[DDTTYLogger sharedInstance] withLevel:DDLogLevelDebug];
        } else {
            [DDLog addLogger:[DDASLLogger sharedInstance] withLevel:DDLogLevelInfo];
            [DDLog addLogger:[DDTTYLogger sharedInstance] withLevel:DDLogLevelInfo];
        }

        DDLogDebug(@"codepoints on CLI : %@", [settings arguments]);

        if (settings.printHelp) {
            [options printHelp];
            return EXIT_SUCCESS;
        }

        if (settings.printVersion) {
            [options printVersion];
            return EXIT_SUCCESS;
        }

        IGApp *app = [[IGApp alloc] init];
        retcode = [app runWithSettings:settings];
    }
    return retcode;
}
