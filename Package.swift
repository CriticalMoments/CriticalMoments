// swift-tools-version: 5.7
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription
import Foundation

// Production release binary by default
var appcoreTarget = Target.binaryTarget(
    name: "Appcore",
    url: "https://github.com/CriticalMoments/CriticalMoments/releases/download/0.1.8-beta/Appcore.xcframework.zip",
    checksum: "5d96757dbe1103c98fc3dbcdcfa8a8a5d1bc0d99599cd3012fddc0cea13c83a0")

// If this device has built the appcore framework locally, use that. This is primarily for development.
// We highly recommend end users use the production binary.
// If you don't trust the precompiled binaries, you can verify the checksums/source from Github release action logs which built it https://github.com/CriticalMoments/CriticalMoments/actions.
// Building yourself should work, but requires additional tooling (golang) and we don't offer support for this flow.
let filePath = #filePath
let endOfPath = filePath.count - "Package.swift".count - 1
let dirPath = String(filePath[...String.Index.init(utf16Offset: endOfPath, in: filePath)])
let infoPath = dirPath + "go/appcore/build/Appcore.xcframework/Info.plist"
if (FileManager.default.fileExists(atPath: infoPath))
{
    print("Using Local Appcore Build From: " + infoPath);
    appcoreTarget = Target.binaryTarget(
        name: "Appcore",
        path: "go/appcore/build/Appcore.xcframework")
}

let package = Package(
    name: "CriticalMoments",
    platforms: [.iOS(.v12)],
    products: [
        // Products define the executables and libraries a package produces, making them visible to other packages.
        .library(
            name: "CriticalMoments",
            targets: ["CriticalMoments"]),
    ],
    targets: [
        // Targets are the basic building blocks of a package, defining a module or a test suite.
        // Targets can depend on other targets in this package and products from dependencies.
        .target(
            name: "CriticalMoments",
            dependencies: ["Appcore"],
            path: "ios/Sources/CriticalMoments",
            publicHeadersPath:"include"),
        appcoreTarget,
        .testTarget(
            name: "CriticalMomentsTests",
            dependencies: ["CriticalMoments"],
            path: "ios/Tests/CriticalMomentsTests",
            resources: [
                .copy("TestResources")
            ],
            cSettings: [
                .headerSearchPath("../../Sources/CriticalMoments"),
            ]
        ),
    ],
    swiftLanguageVersions: [.v5]
)
