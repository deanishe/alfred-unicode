//
//  IGApp.m
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import "IGApp.h"

@implementation IGApp

@synthesize fileManager;
@synthesize defaults;
@synthesize codePoints;


// Apply application default settings
- (void)applyFactoryDefaults
{
    [defaults setObject:0 forKey:@"limit"];
    [defaults setObject:[NSNumber numberWithFloat:256.0f] forKey:@"iconsize"];
    [defaults setObject:@"ArialUnicodeMS" forKey:@"font"];
    [defaults setObject:[fileManager currentDirectoryPath] forKey:@"outputdir"];
//    [defaults setObject:[[fileManager currentDirectoryPath]
//                         stringByAppendingPathComponent:@"characters.tsv"]
//                 forKey:@"charlist"];
    [defaults setBool:NO forKey:@"help"];
    [defaults setBool:NO forKey:@"list"];
    [defaults setBool:NO forKey:@"verbose"];
}

- (void)parseCommandLineArguments
{
    NSArray *args = [[NSProcessInfo processInfo] arguments];
    for (NSString *arg in args) {
        if ([arg isEqualToString:@"-verbose"] ||  [arg isEqualToString:@"-v"]) {
            [defaults setBool:YES forKey:@"verbose"];
        }
        if ([arg isEqualToString:@"-help"] ||  [arg isEqualToString:@"-h"]) {
            [defaults setBool:YES forKey:@"help"];
        }
        if ([arg isEqualToString:@"-list"] ||  [arg isEqualToString:@"-l"]) {
            [defaults setBool:YES forKey:@"list"];
        }
    }
    NSString *codePoint = [defaults stringForKey:@"codepoint"];
    if (codePoint == NULL) {
        codePoints = [NSArray array];
        return;
    }
    if ([codePoint isEqualToString:@"-"]) {
        // Load codepoints from STDIN
    } else {
        codePoints = [NSArray arrayWithObject:codePoint];
    }
}

// Return Unicode string for codePoint
- (NSString *)stringWithCodePoint:(uint32_t)codePoint
{
    NSString *string = [[NSString alloc] initWithBytes:&codePoint length:4 encoding:NSUTF32LittleEndianStringEncoding];
    return string;
}

// Return absolute filepath for path
- (NSString *)resolvePath:(NSString *)path
{
    NSString *expandedPath = [path stringByStandardizingPath];
    NSLog(@"path : %@, expandedPath : %@", path, expandedPath);

    if ([expandedPath hasPrefix:@"/"]) {
        // expandedPath is absolute
        return expandedPath;
    }

    NSString *absolutePath = [[[[[NSFileManager alloc] init]
                                currentDirectoryPath]
                               stringByAppendingPathComponent:expandedPath]
                              stringByStandardizingPath];

    NSLog(@"absolutePath : %@", absolutePath);

    return absolutePath;
}


- (void)printFontList
{
    NSFontManager *fontManager = [NSFontManager sharedFontManager];
    NSArray *allFonts = [fontManager availableFonts];

    printf("The following fonts are available on your system:\n");

    for (NSString *fontName in allFonts) {
        printf("\t%s\n", [fontName UTF8String]);
    }
}

- (void)printHelp
{
    NSString *help = @"IconGen [-v] [-font <fontname>] [-outputdir <dirpath>] -codepoint <codepoint>\n\n"
                     @"Generate preview icon for `codepoint` in `outputdir`.\n\n"
                     @"Usage:\n"
                     @"    IconGen [-v] -codepoint 0001F4A9\n"
                     @"\n"
                     @"Options:\n"
                     @"    -limit <NUM>     Only generate NUM icons (for testing).\n"
                     @"    -codepoint <CP>  Generate icon for this codepoint.\n"
                     @"                     Codepoint must be 8-digit hexadecimal or '-'.\n"
                     @"                     If '-', codepoints will be read, one per line, from STDIN.\n"
                     @"    -h, -help        Show this messsage and exit.\n"
                     @"    -v, -verbose     Show more status output.\n"
                     @"\n";
    printf("%s", [help UTF8String]);
}


- (int)run
{
    fileManager = [NSFileManager defaultManager];
    defaults = [NSUserDefaults standardUserDefaults];
    [self applyFactoryDefaults];
    [self parseCommandLineArguments];



    if ([defaults boolForKey:@"verbose"]) {
        NSLog(@"Settings :");
        NSLog(@"  limit      : %ld", (long)[defaults integerForKey:@"limit"]);
        NSLog(@"  font       : %@", [defaults stringForKey:@"font"]);
        NSLog(@"  outputdir  : %@", [defaults stringForKey:@"outputdir"]);
        NSLog(@"  codepoint  : %@", [defaults stringForKey:@"codepoint"]);
        NSLog(@"  codepoints : %@", codePoints);
        NSLog(@"  iconsize   : %f", [defaults floatForKey:@"iconsize"]);
        NSLog(@"  help       : %d", [defaults boolForKey:@"help"]);
        NSLog(@"  list       : %d", [defaults boolForKey:@"list"]);
    }

    if ([defaults boolForKey:@"help"] == YES) {
        [self printHelp];
        return EXIT_SUCCESS;
    }

    if ([defaults stringForKey:@"codepoint"] == NULL) {
        printf("ERROR: No codepoint specified.\n\n");
        [self printHelp];
        return EXIT_FAILURE;
    }

    return EXIT_SUCCESS;
}

@end
