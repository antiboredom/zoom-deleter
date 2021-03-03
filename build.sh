# make the tray icons
echo "generating icon"
cd icon
convert icon.png -flatten -colors 256 -background transparent icon.ico
bash make_icon.sh

echo "building"
# build it
cd ..
env GO111MODULE=on go build
env GO111MODULE=on GOOS=windows go build -ldflags "-H=windowsgui"

echo "copying to mac app"
# copy to the mac application
cp zoom_deleter Zoom\ Deleter.app/Contents/MacOS
