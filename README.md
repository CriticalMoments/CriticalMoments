# CriticalMoments iOS

[![Release Tests](https://github.com/CriticalMoments/CriticalMoments/actions/workflows/test_release.yml/badge.svg)](https://github.com/CriticalMoments/CriticalMoments/actions/workflows/test_release.yml)
[![Carthage compatible](https://img.shields.io/badge/Carthage-compatible-4BC51D.svg?style=flat)](https://github.com/Carthage/Carthage)
[![CI Status](https://img.shields.io/travis/scosman/CriticalMoments.svg?style=flat)](https://travis-ci.org/scosman/CriticalMoments)
[![Version](https://img.shields.io/cocoapods/v/CriticalMoments.svg?style=flat)](https://cocoapods.org/pods/CriticalMoments)
[![License](https://img.shields.io/cocoapods/l/CriticalMoments.svg?style=flat)](https://cocoapods.org/pods/CriticalMoments)
[![Platform](https://img.shields.io/cocoapods/p/CriticalMoments.svg?style=flat)](https://cocoapods.org/pods/CriticalMoments)

# Work in Progress

This project is a work in progress, and is not ready for any usage.

# Requirements

Currenlty CriticalMoments supports iOS and iPad OS, for iOS 11+.

The API supports both Objective-C and Swift.

# Installation Options

CriticalMoments can be installed several ways:

 - Swift Package Manager
 - CocoaPods
 - Direct framework download

Carthage is not supported, but direct xcframework download is similar, and suggested for those using Carthage.

## Swift Package Manager Installation

Follow Apple's instructions for [Adding package dependencies to your app using swift package manager](https://developer.apple.com/documentation/xcode/adding-package-dependencies-to-your-app).

The git URL to enter is simply: `https://github.com/CriticalMoments/CriticalMoments`

We suggest specifing "Up to next major version" as your dependancy rule. We don't recommend the dependancy rule `branch=main`, as main may contain pre-release code.

## CocoaPods Installation

CriticalMoments is available through [CocoaPods](https://cocoapods.org). 

To install it, follow the usual Cocoapods steps: 

 - Add the pod to your Podfile. A line like `pod 'CriticalMoments'`, optionally locking to a version or major release
 - Run `pod install` and confirm the output indicates the CM installation was successful
 - Clean and build your project
 - Restart Xcode (yes, this is usually needed...)
 - Link Critical Moments in the "Link Binary with Libraries" section, inside "Build Phases" tab of your target
 - Import critical moments where needed
   - Objective C: `#import "CriticalMoments.h"` 
   - Swift: `import CriticalMoments` 

## Direct Framework Download Installation

 - Download `CriticalMoments.xcframework` (link coming soon)
 - Add the framework to your project by dragging into the "Frameworks, Libraries, and Embedded Content" section of your project in xcode
 - Import and use the framework where needed
   - Objective C: `@import CriticalMoments;`
   - Swift: `import CriticalMoments`

## Carthage Installation

CriticalMoments is available through [Carthage](https://github.com/Carthage/Carthage)

There are a few extra steps to install via Carthage so please follow steps below carefully. These are needed because this project uses `Package.swift`, and Carthage doesn't yet support it. You must run the [XcodeGen](https://github.com/yonaskolb/XcodeGen) tool to build the project files Carthage needs. If you prefer to not install other tools, the direct framework download approach is very similar to Carthage.

  - Add CriticalMoments to you `Cartfile` with a line like `github https://github.com/CriticalMoments/CriticalMoments`, optionally including a version requirement
  - Run `carthage update --use-xcframeworks`. This will fail because of the missing xcodeproj, but is needed to populate your /Carthage/Checkouts cache
  - Run `xcodegen generate --spec ./Carthage/Checkouts/criticalmoments/ios/project.yml` (if you don't have xcodegen installed already, install with `brew install xcodegen`)
  - Run `carthage build --no-skip-current --use-xcframeworks`, this time the build should succeed as the step above created needed project files 
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
