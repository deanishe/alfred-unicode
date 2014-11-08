<?php

// From https://github.com/sperelson/Awesome2PNG/blob/master/awesome2png.php
class TTF2PNG
{
	/*
	 * Generate PNG from Font-Awesome TTF
	 * ----------------------------------
	 * The constructor takes the Wanted PNG height in pixels and brute-forces
	 * the correct (or best effort attempt) height in points to render the font.
	 * This may be a slow process, but the result is expected to be cached in a
	 * file. So, pretty pictures that are correctly sized, instead of a speedy
	 * result.
	 *
	 * The contructor also takes in the unicode hex value of the wanted
	 * character and a hex colour code. Also, a padding array with pixel counts:
	 *     array(top, bottom, left, right)
	 *
	 * The result is an awesome transparent PNG glyph with the character of
	 * choice rendered in the colour expected with the padding of choice
	 * applied. Padding can be positive or negative. A negative value will cause
	 * the glyph to crop.
	 *
	 * Requires GD version 2, FreeType, and a FontAwesome TTF file
	 *
	 * Copyright Stephen Perelson
	 */
	public function __construct($unicodeChar, $pixelshigh=30, $color='000000', $alpha=0, $padding=array(0, 0, 3, 3)) {
		// Variables for brute-forcing the correct point height
		$ratio = 96 / 72;
		$ratioadd = 0.0001;
		$heightalright = false;
		$count = 0;
		$maxcount = 20000;
		// Set the enviroment variable for GD
		putenv('GDFONTPATH=' . realpath('.'));
		$font = 'Arial Unicode';

		// $text = json_decode('"'.$unicodeChar.'"');
		$text = html_entity_decode($unicodeChar, ENT_NOQUOTES, 'UTF-8');
		// Brute-force point height
		while (!$heightalright && $count < $maxcount) {
			$x = $pixelshigh / $ratio;
			$count++;
			$bounds = imagettfbbox($x, 0, $font, $text);
			$height = abs($bounds[7] - abs($bounds[1]));
			if ($height == $pixelshigh) {
				$heightalright = true;
			} else {
				if ($height < $pixelshigh) {
					$ratio -= $ratioadd;
				} else {
					$ratio += $ratioadd;
				}
			}
		}
		$width = abs($bounds[4]) + abs($bounds[0]);
		// Create the image
		$im = imagecreatetruecolor($width + $padding[2] + $padding[3], $pixelshigh + $padding[0] + $padding[1]);
		imagesavealpha($im, true);
		$trans = imagecolorallocatealpha($im, 0, 0, 0, 127);
		imagefill($im, 0, 0, $trans);
		imagealphablending($im, true);
		$fontcolor = self::makecolor($color, $alpha);
		// Add the text
		imagettftext($im, $x, 0, 1 + $padding[2], $height-abs($bounds[1])-1 + $padding[0], $fontcolor, $font, $text);
		imagesavealpha($im, true);
		imagepng($im); // Outputs the image
		imagedestroy($im);
	}
	static private function makecolor($hexcolor, $alpha) {
		return $alpha << 24 | hexdec($hexcolor);
	}
}


function main() {
	$hex = strtoupper($_GET['h']);
	$colour = $_GET['c'];
	$char = '&#x' . $hex . ';';
	header('Content-type: image/png');
	new TTF2PNG($char, 256, $colour);
	// header('Content-type: text/plain; charset=utf-8');
	// $text = html_entity_decode($char, ENT_NOQUOTES, 'UTF-8');
	// echo($text);

}

main();

?>