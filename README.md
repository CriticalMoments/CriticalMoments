# CriticalMoments iOS

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

We suggest specifing "Up to next major version" as your dependancy rule. We don't recommend installing branch=main, as main may contain pre-release code.

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
