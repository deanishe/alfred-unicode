// Copyright (c) 2021 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2021-02-08

import Foundation
import AppKit

struct FontInfo: Encodable {
	let family: String
	let face: String

	init(family: String, face: String) {
		self.family = family
		self.face = face
	}
}


func availableFonts() -> [FontInfo] {
    let fm = NSFontManager.shared
    var fonts = [FontInfo]()
    for family in fm.availableFontFamilies {
        if let faces = fm.availableMembers(ofFontFamily: family) {
            for arr in faces {
                let name = arr[0] as! String
                let face = arr[1] as! String
                if fm.fontNamed(name, hasTraits: NSFontTraitMask.italicFontMask) {
                    continue
                }
                fonts.append(FontInfo(family: family, face: face))
            }
        }
    }
    return fonts
}