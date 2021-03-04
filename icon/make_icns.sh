#!/bin/bash

# source: https://github.com/GitableRoy/MakeICNS/blob/master/makeicns.sh
# This script accepts a [ PNG | JPG | TIF | SVG ] file and returns an .icns

# if user doesn't have imagemagick, exit
checkMagick=$(which magick)
if  [ -z $checkMagick ]; then
  echo "You must have imagemagick installed to use this script!"
  exit 0
fi

# check for arguments
if [[ $# -eq 0 ]]; then
  echo "No argument entered. Please supply this script with a PNG!"
  exit 0
fi


CreateImageFile() {
 # if file placed in script is not a PNG, JPG, TIF, or SVG then exit
 if ! [[ "$1" == *"png"* || "$1" == *"jpg"* || "$1" == *"jpeg"* \
       || "$1" == *"tif"* || "$1" == *"svg"* ]]; then
   echo "The file "$1" is a "$(echo $1|cut -d "." -f2)" file. \
         Please enter a .png, .jpg, .tif, or .svg file"
   continue
 fi

 filename=$(echo $1 | cut -d "." -f1)

 # resize file to same dimensions
 size=$(echo magick identify -format "%wx%h\n" $1)
 width=$($size | cut -d "x" -f1 | bc)
 height=$($size | cut -d "x" -f2 | bc)

 # create a resized png if dimenions were not the same
 if ! [[ $width == $height ]]; then
   if [[ $width > $height ]]; then
     lgrVal=$width
   else
     lgrVal=$height
   fi
   magick convert $1 -resize $lgrVal"x"$lgrVal\! -quality 100 $filename"_resize.png"
   ofilename=$filename
   filename=$(echo $filename"_resize")
 fi


 # create necessary icns files in /tmp
 mkdir /tmp/iconbuilder.iconset

 magick convert $1 -resize 512x512! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_512x512.png"
 magick convert $1 -resize 256x256! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_256x256.png"
 magick convert $1 -resize 128x128! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_128x128.png"
 magick convert $1 -resize 64x64! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_64x64.png"
 magick convert $1 -resize 32x32! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_32x32.png"
 magick convert $1 -resize 16x16! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_16x16.png"

 magick convert $1 -resize 512x512! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_512x512@2x.png"
 magick convert $1 -resize 256x256! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_256x256@2x.png"
 magick convert $1 -resize 128x128! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_128x128@2x.png"
 magick convert $1 -resize 64x64! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_64x64@2x.png"
 magick convert $1 -resize 32x32! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_32x32@2x.png"
 magick convert $1 -resize 16x16! -fuzz 01% -transparent white \
     -quality 100 /tmp/iconbuilder.iconset/"icon_16x16@2x.png"

 iconutil --convert icns --output "icon.icns" /tmp/iconbuilder.iconset


 # clean up by removing iconbuilder.folder
 rm -rf /tmp/iconbuilder.iconset

 # remove extra file if resize was necessary and change .icns to original filename
 if [[ $filename == *"resize"* ]]; then
   mv $filename".icns" $ofilename".icns"
   rm $(echo $filename".png")
 fi
}


for arg in "$@"
  do
    echo "Creating icon for $arg"

    # if argument is a file then call func 
    if [[ -f $arg ]]; then
      CreateImageFile $arg

    # if argument is a dir then loop through all files in dir and call func
    elif [[ -d $arg ]]; then
      for sub in $arg/*; do
        if [[ -f $sub ]]; then
          CreateImageFile $sub
        fi
      done
    fi

  done
exit 0

