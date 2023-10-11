//
//  CMCallPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-07-07.
//

#import "CMCallPropertyProvider.h"

@import CallKit;

@implementation CMCallPropertyProvider

- (long)type {
    return AppcoreLibPropertyProviderTypeBool;
}

- (BOOL)boolValue {
    CXCallObserver *ck = [[CXCallObserver alloc] init];
    return ck.calls.count > 0;
}

@end
