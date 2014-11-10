//
//  PreviewApp.m
//  UnicodeImageGen
//
//  Created by Dean Jackson on 08/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import "PreviewApp.h"
#import "PreviewIcon.h"
#import "CharacterDictionary.h"

@implementation PreviewApp

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

    fileManager = [NSFileManager defaultManager];

    CharacterDictionary *myDict = [[CharacterDictionary alloc] init];
    _characters = [myDict characters];
    NSArray *characterNames = [_characters allKeys];

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

        PreviewIcon *icon = [[PreviewIcon alloc] init];
        icon.fontName = [settings fontName];
        icon.iconSize = [settings iconSize];
        icon.character = string;
        icon.filePath = filePath;
        if ([settings verbose]) {
            NSLog(@"%@ -> %@", string, filePath);
        }
        [icon savePreview];
        i++;
        pc = ((float)i / (float)total) * 100.0;
        int mod = i % 100;
        if (mod == 0) {
            CFTimeInterval elapsed = CFAbsoluteTimeGetCurrent() - startTime;
            NSLog(@"[%0.2f%%] Wrote %5d icons in %5.2f seconds", pc, i, elapsed);
        }

    }

    CFTimeInterval elapsed = CFAbsoluteTimeGetCurrent() - startTime;
    NSLog(@"Wrote %lu icons in %0.2f seconds", [characterNames count], elapsed);

    return 0;
}

@end
