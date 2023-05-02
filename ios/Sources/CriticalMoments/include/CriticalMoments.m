//
//  CriticalMoments.m
//  CriticalMoments
//
//  Created by Steve Cosman on 2023-04-17.
//

#import "CriticalMoments.h"

@import Appcore;

@implementation CriticalMoments

+ (NSString *)objcPing {
    return @"objcPong";
}

+ (NSString *)goPing {
    return AppcoreGoPing();
}

@end
