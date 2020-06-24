rm -rf output/*
mkdir -p output/ios/
mkdir -p output/android/
echo "Building for iOS..."
CGO_CFLAGS_ALLOW='-fmodules|-fblocks' gomobile bind -target=ios -o=output/ios/atoken.framework github.com/tsfdsong/atoken-app-sdk/blockchain github.com/tsfdsong/atoken-app-sdk/vexchain github.com/tsfdsong/neo-utils

echo "Building for Android..."
ANDROID_HOME=/Users/$USER/Library/Android/sdk/android-ndk-r21 gomobile bind -target=android -o=output/android/atoken.aar github.com/tsfdsong/atoken-app-sdk/blockchain github.com/tsfdsong/atoken-app-sdk/vexchain github.com/tsfdsong/neo-utils
