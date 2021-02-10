// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-07

import Foundation
import AppKit

let iconSize = 128

enum OutputType {
    case error
    case standard
}


func log(_ message: String, to: OutputType = .standard) {
    switch to {
    case .standard:
        print("\(message)")
    case .error:
        fputs("\(message)\n", stderr)
    }
}

func listFonts() throws {
    let fonts = availableFonts()
    log("\(fonts.count) fonts", to: .error)
    let encoder = JSONEncoder()
    let options : JSONEncoder.OutputFormatting = [.prettyPrinted, .sortedKeys]
    encoder.outputFormatting = options
    let data = try encoder.encode(fonts)
    print(String(data: data, encoding: .utf8)!)
}


func main() throws {
    if CommandLine.arguments.count == 2 {
        let s = CommandLine.arguments[1]
        if s == "-h" || s == "-help" || s == "--help" {
            print("icongen [-h] [--fonts] <file>...")
            print("")
            print("Generate images with text")
            print("")
            print("Usage:")
            print("\ticongen <file>...")
            print("\ticongen --fonts")
            print("\ticongen -h")
            print("")
            print("Options:")
            print("\t--fonts\tlist available fonts")
            print("\t-h, --help\tshow this message and exit")
            return
        }
        if s == "--fonts" {
            try listFonts()
            return
        }
    }
    for i in 1..<CommandLine.arguments.count {
        let fileManager = FileManager.default
        let jsonPath = CommandLine.arguments[i]
        let jsonURL = URL(fileURLWithPath: NSString(string: jsonPath).expandingTildeInPath)

        log("reading JSON from \(jsonPath) ...")
        let data = try Data(contentsOf: jsonURL)
        let decoder = JSONDecoder()
        let configs = try decoder.decode([Configuration].self, from: data)

        let generator = IconGenerator()
        for (index, conf) in configs.enumerated() {
            let status = "[\(index+1)/\(configs.count)]"
            log("\(status) colour=\(conf.textColour.hex()), bg=\(conf.bgColour?.hex() ?? "transparent"), size=\(conf.size), path=\(conf.url.path)")
            try fileManager.createDirectory(at: conf.url.deletingLastPathComponent(),
                   withIntermediateDirectories: true,
                                     attributes: nil)
            try generator.saveIcon(size: conf.size, text: conf.text, font: conf.font, to: conf.url, textColour: conf.textColour, bgColour: conf.bgColour)
        }
    }
}

do {
    try main()
} catch {
    log("error: \(error)", to: .error)
    exit(1)
}
