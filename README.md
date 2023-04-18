# CriticalMoments iOS

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
 - Carthage
 - CocoaPods
 - Direct framework download

## Swift Package Manager Installation

Instructions coming soon.

## Carthage Installation

CriticalMoments is available through [Carthage](https://github.com/Carthage/Carthage)

To install, follow the usual Carthage steps:

 - Add CriticalMoments to you `Cartfile` with a line like `github https://github.com/CriticalMoments/CriticalMoments`
 - Run `carthage update --use-xcframeworks`
 - Drag the built .xcframework bundles from Carthage/Build into the "Frameworks and Libraries" section of your application’s Xcode project.
 - Build the project
 - Import CriticalMoments where needed:
   - Objective C: `@import CriticalMoments;`
   - Swift: `import CriticalMoments`

For details see [Carthage Docs](https://github.com/Carthage/Carthage#adding-frameworks-to-an-application).

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

Instructions coming soon.

# Author

Steve Cosman: https://scosman.net

# License

Copyright (c) 2023 Stephen Cosman.

All rights reserved. 
