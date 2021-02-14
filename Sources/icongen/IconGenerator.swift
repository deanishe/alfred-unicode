// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-07

import Foundation
import AppKit

enum IconGeneratorError: Error {
    case noScreenInfo
    case invalidFont
    case invalidColour
    case initialiseGenerator
    case jsonEncoding

    case noFontForText
    case invalidImageData
}

struct IconGenerator {

    func saveIcon(size: Int, text: String, font: NSFont, to url: URL, textColour: NSColor, bgColour: NSColor?) throws {
        let imageSize = NSSize(width: size, height: size)
        // generate image
        let image = NSImage(size: imageSize)
        let imageRect = NSRect(x: 0, y: 0, width: size, height: size)
        image.lockFocus()

        if let bg = bgColour {
            bg.setFill()
            imageRect.fill()
        }

        let textAttributes : [NSAttributedString.Key : Any] = [
            .font: font,
            .foregroundColor: textColour,
        ]
        let attrStr = NSAttributedString.init(string: text, attributes: textAttributes)

        let textSize = attrStr.size()
        let x = (size - Int(textSize.width)) / 2
        let y = (size - Int(textSize.height) + Int(font.descender)) / 2
        log("size=\(size), textSize=\(textSize), origin=(\(x), \(y))")
        attrStr.draw(in: NSRect(origin: CGPoint(x: x, y: -y), size: imageSize))
        image.unlockFocus()

        // save image
        guard let tiffData = image.tiffRepresentation else {
            throw IconGeneratorError.invalidImageData
        }

        guard let imgRep = NSBitmapImageRep(data: tiffData) else {
            throw IconGeneratorError.invalidImageData
        }

        let imgProps : [NSBitmapImageRep.PropertyKey : Any] = [.interlaced: 1]

        guard let pngData = imgRep.representation(using: NSBitmapImageRep.FileType.png, properties: imgProps) else {
            throw IconGeneratorError.invalidImageData
        }

        do {
            try pngData.write(to: url, options: Data.WritingOptions.atomic)
        } catch {
            throw error
        }
    }
}
