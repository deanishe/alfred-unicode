//
//  IGApp.m
//  IconGenerationHelper
//
//  Created by Dean Jackson on 26/04/2015.
//  Copyright (c) 2015 Dean Jackson. All rights reserved.
//

#import "IGApp.h"

@implementation IGApp

@synthesize codePoints;
@synthesize fileManager;
@synthesize iconGenerator;
@synthesize outputDirectory;
@synthesize settings;



- (void)printFontList
{
    NSFontManager *fontManager = [NSFontManager sharedFontManager];
    NSArray *allFonts = [fontManager availableFonts];

    printf("The following fonts are available on your system:\n");

    for (NSString *fontName in allFonts) {
        printf("\t%s\n", [fontName UTF8String]);
    }
}


- (void)printVersion
{
    printf("IconGen v%s\n", [[settings versionNumber] UTF8String]);
}


// Return absolute filepath for path
- (NSString *)resolvePath:(NSString *)path
{
    NSString *expandedPath = [path stringByStandardizingPath];
//    NSLog(@"path : %@, expandedPath : %@", path, expandedPath);

    if ([expandedPath hasPrefix:@"/"]) {
        // expandedPath is absolute
        return expandedPath;
    }

    NSString *absolutePath = [[[[[NSFileManager alloc] init]
                                currentDirectoryPath]
                               stringByAppendingPathComponent:expandedPath]
                              stringByStandardizingPath];

//    NSLog(@"absolutePath : %@", absolutePath);

    return absolutePath;
}

// Return the filepath for the specified codePoint
- (NSString *)filePathWithCodePoint:(NSString *)codePoint
{
    NSString *dirPath = [outputDirectory stringByAppendingPathComponent:
               [codePoint substringToIndex:2]];
    dirPath = [dirPath stringByAppendingPathComponent:
               [codePoint substringWithRange:NSMakeRange(2, 2)]];
    dirPath = [dirPath stringByAppendingPathComponent:
               [codePoint substringWithRange:NSMakeRange(4, 2)]];
    NSString *fileName = [NSString stringWithFormat: @"%@.png", codePoint];
    NSString *filePath = [dirPath stringByAppendingPathComponent: fileName];
    return filePath;
}

// Return int for codePoint
- (uint32_t)integerWithCodePoint:(NSString *)codePoint
{
    uint32_t value;
    [[NSScanner scannerWithString:codePoint] scanHexInt: &value];
    return value;
}

// Return Unicode string for codePoint
- (NSString *)stringWithCodePoint:(NSString *)codePoint
{
    uint32_t codePointInteger = [self integerWithCodePoint:codePoint];
    NSString *string = [[NSString alloc] initWithBytes:&codePointInteger
                                                length:4
                                              encoding:NSUTF32LittleEndianStringEncoding];
    return string;
}

- (NSString *)validateCodePoint:(NSString *)codePoint
{
    NSString *uc = [codePoint uppercaseString];
    NSString *regex = @"[A-F0-9]{8}";
    NSPredicate *regexText = [NSPredicate predicateWithFormat:@"SELF MATCHES %@", regex];
    if ([regexText evaluateWithObject:uc]) {
        return uc;
    } else {
        return nil;
    }
}

- (int)runWithSettings:(GBSettings *)_settings
{
    fileManager = [NSFileManager defaultManager];
    settings = _settings;
    NSMutableArray *allCodePoints = [NSMutableArray array];
    NSMutableArray *validCodePoints = [NSMutableArray array];
    outputDirectory = [self resolvePath:[settings outputDirectory]];

    // Alternative actions
    if ([settings printFontList]) {
        [self printFontList];
        return EXIT_SUCCESS;
    }

    if ([settings printVersion]) {
        [self printVersion];
        return EXIT_SUCCESS;
    }

    // Check output directory isn't a file
    BOOL isDir;
    if ([fileManager fileExistsAtPath:outputDirectory
                          isDirectory:&isDir] && !isDir) {
        printf("ERROR: Is not a directory : %s", [outputDirectory UTF8String]);
        return EXIT_FAILURE;
    }

    // Load codepoints from ARGV or STDIN
    if ([[settings arguments] count] > 0) {
        // Take rest of ARGV
        for (NSString *codePoint in [settings arguments]) {
            [allCodePoints addObject:codePoint];
        }
    } else {
        // Read from STDIN
        NSLog(@"Reading codepoints from STDIN...");
        NSFileHandle *fh = [NSFileHandle fileHandleWithStandardInput];
        NSData *inputData = [NSData dataWithData:[fh readDataToEndOfFile]];
        NSString *inputString = [[NSString alloc] initWithData:inputData encoding:NSUTF8StringEncoding];
        NSArray *inputLines = [inputString componentsSeparatedByString:@"\n"];
        for (NSString *line in inputLines) {
            if ([[line stringByTrimmingCharactersInSet:[NSCharacterSet whitespaceAndNewlineCharacterSet]] length] > 0) {
                [allCodePoints addObject:line];
            }
        }

    }

    // Validate codepoints
    for (NSString *codePoint in allCodePoints) {
        NSString *sanitisedCodePoint = [self validateCodePoint:codePoint];
        if (sanitisedCodePoint == nil) {
            printf("ERROR: Invalid code point : %s\n", [codePoint UTF8String]);
        } else {
            [validCodePoints addObject:sanitisedCodePoint];
        }
    }

    if ([validCodePoints count] == 0) {
        printf("ERROR: No codepoints to generate icons for. Try --help.\n");
        return EXIT_FAILURE;
    }

    // Truncate codepoints if limit is set
    if ([settings limit] > 0 && [validCodePoints count] > [settings limit]) {
        if ([settings verbose]) {
            NSLog(@"Truncating codepoints to length %lu", [settings limit]);
        }
        codePoints = [validCodePoints subarrayWithRange:NSMakeRange(0, MIN([settings limit], [validCodePoints count]))];
    } else {
        codePoints = [NSArray arrayWithArray:validCodePoints];
    }

    if ([settings verbose]) {
        NSLog(@"%lu icon(s) to generate in `%@`...", [codePoints count], outputDirectory);
    }

    NSLog(@"Settings:\n");
    NSLog(@"           Font name : %@", [settings fontName]);
    NSLog(@"           Icon size : %0.0f", [settings iconSize]);
    NSLog(@"          Icon count : %lu", [codePoints count]);
    NSLog(@"    Output directory : %@\n", outputDirectory);
    NSLog(@"  Overwrite existing : %d\n", [settings overwrite]);

    iconGenerator = [[IGIconGenerator alloc] init];
    iconGenerator.fontName = [settings fontName];
    iconGenerator.iconSize = [settings iconSize];
    iconGenerator.verbose = [settings verbose];

    // Generate icons
    CFAbsoluteTime startTime = CFAbsoluteTimeGetCurrent();
    int i = 0;
    float pc = 0.0;
    unsigned long total = [codePoints count];

    for (NSString *codePoint in codePoints) {
        NSString *string = [self stringWithCodePoint:codePoint];
        NSString *filePath = [self filePathWithCodePoint:codePoint];

        if ([fileManager fileExistsAtPath:filePath] && ! [settings overwrite]) {
            if ([settings verbose]) {
                NSLog(@"Skipping (already exists) : %@", codePoint);
            }
            i++;
            continue;
        }

        if ([settings verbose]) {
            NSLog(@"%@\t%@\t%@", codePoint, string, filePath);
        }

        // Ensure destination directory exists
        NSString *dirPath = [filePath stringByDeletingLastPathComponent];
        if (! [fileManager fileExistsAtPath:dirPath]) {
            NSURL *url = [NSURL URLWithString:[NSString stringWithFormat:@"file://%@", dirPath]];
            if ([settings verbose]) {
                NSLog(@"Creating directory : %@", url);
            }
            [fileManager createDirectoryAtURL:url
                  withIntermediateDirectories:YES
                                   attributes:nil
                                        error:nil];
        }

        BOOL worked = [iconGenerator saveIcon:filePath withString:string];
        if (! worked) {
            NSLog(@"ERROR: Couldn't create icon for `%@` at `%@`", codePoint, filePath);
        }

        i++;

        // Log progress
        int mod = i % 50;
        if (mod == 0) {
            pc = ((float)i / (float)total) * 100.0;
            CFTimeInterval elapsed = CFAbsoluteTimeGetCurrent() - startTime;
            float perSecond = (float)i / elapsed;
            NSLog(@"[%5d/%5lu] (%6.2f%%) Wrote %5d icons in %6.2fs (%1.1f icons/sec)",
                  i+1, total, pc, i, elapsed, perSecond);
        }
    }
    CFTimeInterval elapsed = CFAbsoluteTimeGetCurrent() - startTime;
    float perSecond = (float)total / elapsed;
    NSLog(@"[%5lu/%5lu] (100.00%%) Wrote %5lu icons in %6.2fs (%1.1f icons/sec)",
          total, total, total, elapsed, perSecond);

    return EXIT_SUCCESS;
}

@end
