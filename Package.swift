// swift-tools-version: 5.8
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "CriticalMoments",
    platforms: [.iOS(.v11)],
    products: [
        // Products define the executables and libraries a package produces, making them visible to other packages.
        .library(
            name: "CriticalMoments",
            targets: ["CriticalMomentsSwift", "CriticalMomentsObjc"]),
    ],
    targets: [
        // Targets are the basic building blocks of a package, defining a module or a test suite.
        // Targets can depend on other targets in this package and products from dependencies.
        .target(
            name: "CriticalMomentsSwift",
            dependencies: ["CriticalMomentsObjc"],
            path: "ios/Sources/CriticalMomentsSwift"),
        .target(
            name: "CriticalMomentsObjc",
            path: "ios/Sources/CriticalMomentsObjc",
            publicHeadersPath:"include"),
        .testTarget(
            name: "CriticalMomentsSwiftTests",
            dependencies: ["CriticalMomentsSwift"],
            path: "ios/Tests/CriticalMomentsSwiftTests"),
        .testTarget(
            name: "CriticalMomentsObjcTests",
            dependencies: ["CriticalMomentsObjc"],
            path: "ios/Tests/CriticalMomentsObjcTests"),
    ],
    swiftLanguageVersions: [.v5]
)
