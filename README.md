# CriticalMoments iOS

[![Release Build](https://github.com/CriticalMoments/CriticalMoments/actions/workflows/publish_xcframework.yml/badge.svg)](https://github.com/CriticalMoments/CriticalMoments/actions/workflows/publish_xcframework.yml)
[![Release Tests](https://github.com/CriticalMoments/CriticalMoments/actions/workflows/test_release.yml/badge.svg)](https://github.com/CriticalMoments/CriticalMoments/actions/workflows/test_release.yml)
[![Carthage compatible](https://img.shields.io/badge/Carthage-compatible-4BC51D.svg?style=flat)](https://github.com/Carthage/Carthage)

# Work in Progress

This project is a work in progress, and is not ready for any usage.

# Requirements

Currenlty CriticalMoments supports iOS and iPad OS, for iOS 11+.

The API supports both Objective-C and Swift.

# Installation Options

CriticalMoments can be installed several ways:

 - [Swift Package Manager (recommended)](#swift-package-manager-installation)
 - [CocoaPods](#cocoapods-installation)
 - [Direct framework download](#direct-framework-download-objective-c)
 - [Carthage](#carthage-installation)

## Swift Package Manager Installation

Follow Apple's instructions for [Adding package dependencies to your app using swift package manager](https://developer.apple.com/documentation/xcode/adding-package-dependencies-to-your-app).

The git URL to enter is simply: `https://github.com/CriticalMoments/CriticalMoments`

We suggest specifing "Up to next major version" as your dependancy rule. We don't recommend the dependancy rule `branch=main`, as main may contain pre-release code.

Then import and use where needed:
 - Objective C: `@import CriticalMoments`
 - Swift: `import CriticalMoments`

## CocoaPods Installation

CriticalMoments is available through [CocoaPods](https://cocoapods.org). 

To install it, follow the usual Cocoapods steps: 

 - Add the pod to your Podfile. A line like `pod 'CriticalMoments', '>= 0.1.4-beta'`, optionally modifing the version requirement
 - Run `pod install` and confirm the output indicates the CM installation was successful
 - Clean and build your project
 - Restart Xcode (yes, this is usually needed...)
 - Link Critical Moments in the "Link Binary with Libraries" section, inside "Build Phases" tab of your target
 - Import critical moments where needed
   - Objective C: `#import "CriticalMoments.h"` 
   - Swift: `import CriticalMoments` 

## Direct Framework Download (Objective C)

You can download an XCFramework for use in your Objective C apps from our [Github Releases](https://github.com/CriticalMoments/CriticalMoments/releases/latest). This approach is not recommended for Swift, as the framework must be built with the exact same toolchain as you app, which is unlikely to match. For Swift, use [Swift Package Manager](#swift-package-manager-installation).

Note: if you choose this route, you should manually update for bug fixes and enhancements.

Process: 

 - Download `xcframework.zip` from the latest [Github Release](https://github.com/CriticalMoments/CriticalMoments/releases/latest)
 - Add the framework to your project by dragging into the "Frameworks, Libraries, and Embedded Content" section of your project in xcode
 - Import and use the framework where needed
   - Objective C: `@import CriticalMoments;`
   - Swift: `import CriticalMoments`

## Carthage Installation

There are a few extra steps to install via Carthage so please follow steps below carefully. These are needed because this project uses `Package.swift`, and Carthage doesn't yet support it. You must run the [XcodeGen](https://github.com/yonaskolb/XcodeGen) tool to build the project files Carthage needs. If you prefer to not install other tools, we suggest using Swift Package Manager.

  - Add CriticalMoments to you `Cartfile` with a line like `github "https://github.com/CriticalMoments/CriticalMoments" >= 0.1.4-beta`, optionally modifing the version requirement
  - Run `carthage update --no-build`
  - Run `xcodegen generate --spec ./Carthage/Checkouts/criticalmoments/ios/project.yml`. If you don't already have xcodegen installed already, you'll need to install with `brew install xcodegen`.
  - Run `carthage build --use-xcframeworks`, this time the build should succeed as the step above created needed project files 
  - Drag the built .xcframework bundles from Carthage/Build into the "Frameworks and Libraries" section of your application’s Xcode project.
  - Build the project
  - Import CriticalMoments where needed:
    - Objective C: `@import CriticalMoments;`
    - Swift: `import CriticalMoments`

# Author

Steve Cosman: https://scosman.net

# License

Copyright (c) 2023 Stephen Cosman.

All rights reserved. 
