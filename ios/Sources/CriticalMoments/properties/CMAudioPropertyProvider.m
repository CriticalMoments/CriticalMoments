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

@interface CMAudioPortPropertyProvider ()
@property(nonatomic, strong) NSArray<AVAudioSessionPort> *matchingPorts;
@end

@implementation CMAudioPortPropertyProvider

- (instancetype)initWithPorts:(NSArray<AVAudioSessionPort> *)ports {
    self = [super init];
    if (self) {
        self.matchingPorts = ports;
    }
    return self;
}

+ (nonnull CMAudioPortPropertyProvider *)hasHeadphones {
    return [[CMAudioPortPropertyProvider alloc]
        initWithPorts:@[ AVAudioSessionPortHeadphones, AVAudioSessionPortBluetoothA2DP ]];
}

+ (nonnull CMAudioPortPropertyProvider *)hasBtHeadphones {
    return [[CMAudioPortPropertyProvider alloc] initWithPorts:@[ AVAudioSessionPortBluetoothA2DP ]];
}

+ (nonnull CMAudioPortPropertyProvider *)hasWiredHeadset {
    return [[CMAudioPortPropertyProvider alloc]
        initWithPorts:@[ AVAudioSessionPortHeadsetMic, AVAudioSessionPortHeadphones ]];
}

+ (nonnull CMAudioPortPropertyProvider *)hasBtHeadset {
    return [[CMAudioPortPropertyProvider alloc] initWithPorts:@[ AVAudioSessionPortBluetoothHFP ]];
}

+ (CMAudioPortPropertyProvider *)hasCarAudio {
    return [[CMAudioPortPropertyProvider alloc] initWithPorts:@[ AVAudioSessionPortCarAudio ]];
}

- (BOOL)boolValue {
    AVAudioSessionRouteDescription *cr = AVAudioSession.sharedInstance.currentRoute;
    NSMutableArray<AVAudioSessionPortDescription *> *allPorts =
        [[NSMutableArray<AVAudioSessionPortDescription *> alloc] init];
    [allPorts addObjectsFromArray:cr.outputs];
    [allPorts addObjectsFromArray:cr.inputs];

    for (AVAudioSessionPortDescription *portDesc in allPorts) {
        if ([self.matchingPorts containsObject:portDesc.portType]) {
            return true;
        }
    }
    return false;
}

- (CMPropertyProviderType)type {
    return CMPropertyProviderTypeBool;
}

@end
