//
//  UGApp.m
//  UnicodeImageGen
//
//  Created by Dean Jackson on 08/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import "NSDictionary+TSV.h"
#import "UGApp.h"
#import "UGIconGenerator.h"
#import "UGCharacterDictionary.h"

@implementation UGApp

//@synthesize iconSize, fontName, outputDirectory, charactersFile, limit, settings;
@synthesize settings, fileManager;
@synthesize characters = _characters;


// Return the filepath for the specified codePoint
- (NSString *)filePathForCodePoint:(NSString *)codePoint
{
    NSString *dirPath = [settings outputDirectory];
    dirPath = [dirPath stringByAppendingPathComponent:
               [codePoint substringToIndex:2]];
    dirPath = [dirPath stringByAppendingPathComponent:
               [codePoint substringWithRange:NSMakeRange(2, 2)]];
    dirPath = [dirPath stringByAppendingPathComponent:
               [codePoint substringWithRange:NSMakeRange(4, 2)]];
    NSString *fileName = [NSString stringWithFormat: @"%@.png", codePoint];
    NSString *filePath = [dirPath stringByAppendingPathComponent: fileName];
    return filePath;
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

- (int)run
{

    if ([settings printFontList]) {
        if ([settings verbose]) {
            NSLog(@"Showing available fonts");
        }
        [self printFontList];
        exit(EXIT_SUCCESS);
    }

    // Store hashes and filepaths
    NSMutableDictionary *hashFilePathMap = [NSMutableDictionary dictionary];
    NSMutableDictionary *codePointFilePathMap = [NSMutableDictionary dictionary];
    NSMutableArray *duplicateIconPaths = [NSMutableArray array];

    fileManager = [NSFileManager defaultManager];

    UGCharacterDictionary *myDict = [[UGCharacterDictionary alloc] init];
    _characters = [myDict characters];
    NSArray *characterNames = [_characters allKeys];

    // Ensure outputDirectory is absolute path
//    settings.outputDirectory = [[[[[NSFileManager alloc] init] currentDirectoryPath]
//                                 stringByAppendingPathComponent:[settings outputDirectory]]
//                                stringByStandardizingPath];
    settings.outputDirectory = [self resolvePath:[settings outputDirectory]];


    if ([settings verbose]) {
        NSLog(@"Saving preview icons to `%@`", [settings outputDirectory]);
    }

    // Truncate character list if --limit was specified
    if ([settings limit] > 0) {
        characterNames = [characterNames subarrayWithRange:NSMakeRange(0, MIN([settings limit], [characterNames count]))];
        if ([settings verbose]) {
            NSLog(@"Truncated characters to length %lu", [characterNames count]);
        }
    }

    if ([settings verbose]) {
        for (NSString *codePoint in characterNames) {
            NSString *string = [_characters valueForKey:codePoint];
            NSLog(@"%@ : %@", codePoint, string);
        }

    }

    // Check output directory isn't a file
    BOOL isDir;
    if ([fileManager fileExistsAtPath:[settings outputDirectory] isDirectory:&isDir] && !isDir) {
        NSLog(@"ERROR: Is not a directory : %@", [settings outputDirectory]);
        exit(EXIT_FAILURE);
    }

    NSLog(@"Using font : %@", [settings fontName]);

    // Generate icons
    CFAbsoluteTime startTime = CFAbsoluteTimeGetCurrent();
    int i = 0;
    float pc = 0.0;
    unsigned long total = [characterNames count];

    NSLog(@"%lu icons to generate and save to %@",
          total, [settings outputDirectory]);

    for (NSString *codePoint in characterNames) {
        NSString *string = [_characters valueForKey:codePoint];
        NSString *filePath = [self filePathForCodePoint:codePoint];

        // Ensure destination directory exists
        NSString *dirPath = [filePath stringByDeletingLastPathComponent];
        if (![fileManager fileExistsAtPath: dirPath]) {
            NSURL *url = [NSURL URLWithString:[NSString stringWithFormat:@"file://%@", dirPath]];
            if ([settings verbose]) {
                NSLog(@"Creating directory : %@", url);
            }
            [fileManager createDirectoryAtURL:url
                  withIntermediateDirectories:YES
                                   attributes:nil
                                        error:nil];
        }

        UGIconGenerator *iconGenerator = [[UGIconGenerator alloc] init];
        iconGenerator.fontName = [settings fontName];
        iconGenerator.iconSize = [settings iconSize];
        iconGenerator.character = string;
        iconGenerator.filePath = filePath;
        if ([settings verbose]) {
            NSLog(@"%@ -> %@", string, filePath);
        }

        NSString *shaHash = [iconGenerator savePreview];
        NSString *existingFilePath = [hashFilePathMap objectForKey:shaHash];

        if (existingFilePath != nil) {
            // Icon with this hash already exists
            if ([settings verbose]) {
                NSLog(@"`%@` is a duplicate of `%@`", filePath, existingFilePath);
            }
            [codePointFilePathMap setObject:existingFilePath forKey:codePoint];
            [duplicateIconPaths addObject:filePath];

        } else {
            // Save the hash and this filepath
            [hashFilePathMap setObject:filePath forKey:shaHash];
            [codePointFilePathMap setObject:filePath forKey:codePoint];
        }

        i++;
        // Log progress
        int mod = i % 10;
        if (mod == 0) {
            pc = ((float)i / (float)total) * 100.0;
            CFTimeInterval elapsed = CFAbsoluteTimeGetCurrent() - startTime;
//            float perIcon = elapsed / (float)i;
            float perSecond = (float)i / elapsed;
            NSLog(@"[%6.2f%%] Wrote %5d icons in %6.2fs (%1.1f icons/second)", pc, i, elapsed, perSecond);
        }
    }

    unsigned long dupeCount = [duplicateIconPaths count];
    pc = ((float)dupeCount / (float)total) * 100.0;
    NSLog(@"%lu duplicate icons to delete (%0.1f%% of all images)", dupeCount, pc);

    int j = 1;
    for (NSString* dupePath in duplicateIconPaths) {
        if ([fileManager fileExistsAtPath:dupePath]) {
            [fileManager removeItemAtPath:dupePath error:nil];
            if ([settings verbose]) {
                NSLog(@"[%4d/%4lu] Deleted : %@", j, dupeCount, dupePath);
            }
        }
        j++;
    }

    CFTimeInterval elapsed = CFAbsoluteTimeGetCurrent() - startTime;
    NSLog(@"Wrote %lu icons in %0.2f seconds", [characterNames count], elapsed);

    NSLog(@"Writing list of codepoints and icon paths to : %@",
          [settings iconListPath]);

    // Data we'll save to the TSV file
    NSMutableDictionary *tsvDictionary = [NSMutableDictionary dictionary];
    for (NSString *codePoint in codePointFilePathMap) {
        NSString *iconPath = [codePointFilePathMap objectForKey:codePoint];
        NSString *relativeIconPath = [iconPath stringByReplacingOccurrencesOfString: [settings outputDirectory]
                                                                         withString: @""];
        relativeIconPath = [relativeIconPath stringByTrimmingCharactersInSet:
                            [NSCharacterSet characterSetWithCharactersInString:@"/"]];

        [tsvDictionary setObject:relativeIconPath forKey:codePoint];

        if ([settings verbose]) {
            NSLog(@"Relative icon path : %@  ->  %@", iconPath, relativeIconPath);
        }
    }

    [tsvDictionary saveToFilePath:[settings iconListPath] error:nil];

    return 0;
}

@end
