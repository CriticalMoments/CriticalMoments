// swift-tools-version: 5.7
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

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
        .binaryTarget(
            name: "Appcore",
            url: "https://github.com/CriticalMoments/CriticalMoments/releases/download/0.1.7-beta/Appcore.xcframework.zip",
            checksum: "d3281ac6f8592830f6adb41524777e67f95fead539b09afa42bfc392fb964737"),
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
