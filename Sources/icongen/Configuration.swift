// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-07

import Foundation
import AppKit

extension NSColor {

    convenience init?(hex: String) {
        var hexCleaned = hex.trimmingCharacters(in: .whitespacesAndNewlines)
        hexCleaned = hexCleaned.replacingOccurrences(of: "#", with: "")
        var rgb: UInt64 = 0
        var r: CGFloat = 0.0
        var g: CGFloat = 0.0
        var b: CGFloat = 0.0
        var a: CGFloat = 1.0

        let length = hexCleaned.count
        guard Scanner(string: hexCleaned).scanHexInt64(&rgb) else { return nil }

        if length == 6 {
            r = CGFloat((rgb & 0xFF0000) >> 16) / 255.0
            g = CGFloat((rgb & 0x00FF00) >> 8) / 255.0
            b = CGFloat(rgb & 0x0000FF) / 255.0
        } else if length == 8 {
            r = CGFloat((rgb & 0xFF000000) >> 24) / 255.0
            g = CGFloat((rgb & 0x00FF0000) >> 16) / 255.0
            b = CGFloat((rgb & 0x0000FF00) >> 8) / 255.0
            a = CGFloat(rgb & 0x000000FF) / 255.0
        } else {
            return nil
        }
        self.init(red: r, green: g, blue: b, alpha: a)
    }

    func hex() -> String {
        let r = Int(self.redComponent * 255)
        let g = Int(self.greenComponent * 255)
        let b = Int(self.blueComponent * 255)
        let a = Int(self.alphaComponent * 255)
        if self.alphaComponent == 1.0 {
            return String(format: "#%02X%02X%02X", r, g, b)
        }
        return String(format: "#%02X%02X%02X%02X", r, g, b, a)
    }
}


fileprivate struct RawConfiguration: Decodable {
    var text: String
    var filepath: String
    var icon_size: Int?
    var font_name: String?
    var text_colour: String?
    var background_colour: String?
    var text_size: Float?
    var resize_factor: Float?
}

struct Configuration: Decodable {

    // let text : NSAttributedString
    let text : String
    let font : NSFont
    let url : URL
    var size : Int = 128
    var textColour = NSColor(red: 1, green: 0.184, blue: 0.663, alpha: 1.0)
    var bgColour : NSColor?

    init(from decoder: Decoder) throws {
        let raw = try RawConfiguration(from: decoder)
        self.text = raw.text
        self.url = URL(fileURLWithPath: NSString(string: raw.filepath).expandingTildeInPath)
        if let n = raw.icon_size {
            self.size = n
        }

        // font
        var fontSize : Float
        if let n = raw.text_size {
            fontSize = n
        } else {
            guard let screen = NSScreen.main else {
                throw IconGeneratorError.noScreenInfo
            }
            guard let resolution = screen.deviceDescription[NSDeviceDescriptionKey.resolution] as! NSSize? else {
                throw IconGeneratorError.noScreenInfo
            }
            let dpi = Float(resolution.width)
            fontSize = Float(size) * (raw.resize_factor ?? 1.0) / (dpi / 72)
        }
        let fontName = raw.font_name ?? "Arial Unicode MS"
        guard let font = NSFont(name: fontName, size: CGFloat(fontSize)) else {
            throw IconGeneratorError.invalidFont
        }

        self.font = font
        // log("fontName=\(fontName), fontSize=\(fontSize), descender=\(font.descender)")

        // colours
        if let hex = raw.text_colour {
            if let colour = NSColor.init(hex: hex) {
                self.textColour = colour
            } else {
                throw IconGeneratorError.invalidColour
            }
        }

        if let hex = raw.background_colour {
            if let colour = NSColor.init(hex: hex) {
                self.bgColour = colour
            } else {
                throw IconGeneratorError.invalidColour
            }
        }
    }
}
