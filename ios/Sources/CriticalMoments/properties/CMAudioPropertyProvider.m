//
//  CMAudioPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-05-26.
//

#import "CMAudioPropertyProvider.h"

@import AVFoundation;

@implementation CMAudioPlayingPropertyProvider

- (BOOL)boolValue {
    return AVAudioSession.sharedInstance.isOtherAudioPlaying;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
}

@end
