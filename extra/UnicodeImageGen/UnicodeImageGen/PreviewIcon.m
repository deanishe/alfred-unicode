//
//  PreviewIcon.m
//  UnicodeImageGen
//
//  Created by Dean Jackson on 08/11/2014.
//  Copyright (c) 2014 Dean Jackson. All rights reserved.
//

#import "PreviewIcon.h"

@implementation PreviewIcon

@synthesize fontName, iconSize, character, filePath, verbose;

- (id)init
{
    if (self = [super init]) {
        verbose = NO;
    }
    return self;
}


- (float)screenDPI
{
    NSScreen *screen = [NSScreen mainScreen];
    NSDictionary *description = [screen deviceDescription];
    NSSize displayPixelSize = [[description objectForKey:NSDeviceSize] sizeValue];
    CGSize displayPhysicalSize = CGDisplayScreenSize(
        [[description objectForKey:@"NSScreenNumber"] unsignedIntValue]);
    float dpi = (displayPixelSize.width / displayPhysicalSize.width) * 25.4f;
//    NSLog(@"DPI is %0.2f", dpi);
    return dpi;
}

/*

 Adapted from:
 http://stackoverflow.com/questions/11442993/how-to-convert-text-to-image-in-cocoa-objective-c
*/

- (void)savePreview
{

    // Create an attributed string with string and font information
    float factor = [self screenDPI] / 72;
    float fontSize = iconSize / factor;
    if (verbose) {
        NSLog(@"Adjusted fontsize : %f for font `%@`", fontSize, fontName);
    }
    CTFontRef font = CTFontCreateWithName((CFStringRef)fontName, fontSize, nil);
    NSDictionary* attributes = [NSDictionary dictionaryWithObjectsAndKeys:
                                (__bridge id)font,
                                kCTFontAttributeName,
                                kCFBooleanTrue,
                                kCTForegroundColorFromContextAttributeName,
                                nil];
    NSAttributedString* as = [[NSAttributedString alloc] initWithString:character
                                                             attributes:attributes];
    CFRelease(font);

    // Figure out how big an image we need
    CTLineRef line = CTLineCreateWithAttributedString((CFAttributedStringRef)as);
    CGFloat ascent, descent, leading;
    double fWidth = CTLineGetTypographicBounds(line, &ascent, &descent, &leading);

    // On iOS 4.0 and Mac OS X v10.6 you can pass null for data
    size_t charWidth = (size_t)ceilf(fWidth);
    size_t charHeight = (size_t)ceilf(ascent + descent);
    void* imageData = malloc(iconSize * iconSize * 4);
//    void* maskData = malloc(iconSize * iconSize * 4);

    // Create the context and fill it with white background
    CGColorSpaceRef space = CGColorSpaceCreateDeviceRGB();
    CGBitmapInfo bitmapInfo = (CGBitmapInfo)kCGImageAlphaPremultipliedLast;

    // Image
    CGContextRef ctx = CGBitmapContextCreate(imageData, iconSize, iconSize, 8,
                                             iconSize * 4, space, bitmapInfo);

    // Mask
//    CGContextRef maskCtx = CGBitmapContextCreate(maskData, iconSize, iconSize, 8,
//                                                 iconSize * 4, space, bitmapInfo);

    CGColorSpaceRelease(space);

    // Image
    CGContextSetRGBFillColor(ctx, 1.0, 1.0, 1.0, 1.0); // White background
    CGContextFillRect(ctx, CGRectMake(0.0, 0.0, iconSize, iconSize));

    // Mask
//    CGContextSetRGBFillColor(maskCtx, 1.0, 1.0, 1.0, 1.0); // white background
//    CGContextFillRect(maskCtx, CGRectMake(0.0, 0.0, iconSize, iconSize));

    // Draw the text
    CGFloat x = (iconSize - charWidth) / 2;
    CGFloat y = ((iconSize - charHeight) / 2) + descent;
    CGContextSetTextDrawingMode(ctx, kCGTextFill);
    CGContextSetRGBFillColor(ctx, 0.83, 0.0, 0.65, 1.0); // Black
    CGContextSetTextPosition(ctx, x, y);
    CTLineDraw(line, ctx);
    CGImageRef imageRef = CGBitmapContextCreateImage(ctx);

    // Mask
//    CGContextSetTextDrawingMode(maskCtx, kCGTextFill);
//    CGContextSetRGBFillColor(maskCtx, 0.0, 0.0, 0.0, 1.0); // black
//    CGContextSetTextPosition(maskCtx, x, y);
//    CTLineDraw(line, maskCtx);
//    CGImageRef maskRef = CGBitmapContextCreateImage(maskCtx);

    CFRelease(line);

    // Mask image
//    CGImageRef mask = CGImageMaskCreate(CGImageGetWidth(maskRef),
//                                        CGImageGetHeight(maskRef),
//                                        CGImageGetBitsPerComponent(maskRef),
//                                        CGImageGetBitsPerPixel(maskRef),
//                                        CGImageGetBytesPerRow(maskRef),
//                                        CGImageGetDataProvider(maskRef),
//                                        NULL,
//                                        false);
//
//    CGImageRef maskedImage = CGImageCreateWithMask(imageRef, mask);

    // Save as JPEG

    NSBitmapImageRep* imageRep = [[NSBitmapImageRep alloc]
                                  initWithCGImage:imageRef];
    NSAssert(imageRep, @"imageRep must not be nil");
    NSDictionary *properties = [NSDictionary dictionaryWithObjectsAndKeys:
                                [NSNumber numberWithBool:NO],
                                NSImageInterlaced, nil];
    NSData* imageBinaryData = [imageRep representationUsingType:NSPNGFileType
                                               properties:properties];
    [imageBinaryData writeToFile:filePath atomically:YES];
    if (verbose) {
        NSLog(@"Wrote %@", filePath);
    }

    // Clean up
    CGContextRelease(ctx);
    CGImageRelease(imageRef);
    free(imageData);
    // Mask
//    CGImageRelease(maskRef);
//    CGImageRelease(mask);
//    CGImageRelease(maskedImage);
//    free(maskData);

}

@end
