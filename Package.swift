// swift-tools-version: 5.7
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "CriticalMoments",
    platforms: [.iOS(.v11)],
    products: [
        // Products define the executables and libraries a package produces, making them visible to other packages.
        .library(
            name: "CriticalMoments",
            targets: ["CriticalMoments", "CriticalMomentsSwift"]),
    ],
    targets: [
        // Targets are the basic building blocks of a package, defining a module or a test suite.
        // Targets can depend on other targets in this package and products from dependencies.
        .target(
            name: "CriticalMomentsSwift",
            path: "ios/Sources/CriticalMomentsSwift"),
        .target(
            name: "CriticalMoments",
            dependencies: ["Appcore"],
            path: "ios/Sources/CriticalMoments",
            publicHeadersPath:"include"),
        .binaryTarget(
            name: "Appcore",
            path: "go/appcore/build/Appcore.xcframework"),
        .testTarget(
            name: "CriticalMomentsSwiftTests",
            dependencies: ["CriticalMomentsSwift"],
            path: "ios/Tests/CriticalMomentsSwiftTests"),
        .testTarget(
            name: "CriticalMomentsTests",
            dependencies: ["CriticalMoments"],
            path: "ios/Tests/CriticalMomentsTests",
            cSettings: [
                .headerSearchPath("../../Sources/CriticalMoments"),
            ]
        ),
    ],
    swiftLanguageVersions: [.v5]
)
