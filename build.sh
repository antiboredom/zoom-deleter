echo "generating icon"
cd icon
convert logo.png -resize 64x64 -background transparent icon.png
convert logo.png -resize 64x64 -flatten -colors 256 -background transparent icon.ico
bash make_icon.sh
bash make_icns.sh logo.png

echo "building"
cd ..
env GO111MODULE=on go build
env GO111MODULE=on GOOS=windows go build -ldflags "-H=windowsgui"

echo "copying to mac app"
cp zoom_deleter Zoom\ Deleter.app/Contents/MacOS/
cp icon/icon.icns Zoom\ Deleter.app/Contents/Resources/

# echo "signing mac app"
# codesign -s "Developer ID Application: Sam Lavigne (8VFY589HXQ)" --timestamp --options runtime -f  --deep Zoom\ Deleter.app


# echo "zipping"
# zip -r zoom_deleter_win.zip zoom_deleter.exe
# zip -r zoom_deleter_mac.zip Zoom\ Deleter.app
