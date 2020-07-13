rm -rf output/*
mkdir -p output/ios/
mkdir -p output/android/
echo "Building for iOS..."

gomobile bind -target=ios -o=output/ios/atoken.framework github.com/tsfdsong/atoken-app-sdk/blockchain github.com/tsfdsong/atoken-app-sdk/defi github.com/tsfdsong/atoken-app-sdk/vexchain github.com/tsfdsong/neo-utils

echo "Building for Android..."
gomobile bind -target=android -o=output/android/atoken.aar github.com/tsfdsong/atoken-app-sdk/blockchain github.com/tsfdsong/atoken-app-sdk/defi github.com/tsfdsong/atoken-app-sdk/vexchain github.com/tsfdsong/neo-utils

echo "Building for zip..."
mkdir -p build/
cd build
zip -r "atoken-$(date +"%Y%m%d-%H:%M:%S").zip" ../output/ios/atoken.framework ../output/android/*
