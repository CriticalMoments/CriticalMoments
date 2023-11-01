//
//  CMDynamicPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-05-22.
//

#import "CMBaseDynamicPropertyProvider.h"

@import Appcore;

@interface CMDynamicPropertyProviderWrapper ()
@property(nonnull, strong) id<CMDynamicPropertyProvider> pp;
@end

@implementation CMDynamicPropertyProviderWrapper

- (instancetype)initWithPP:(id<CMDynamicPropertyProvider>)pp {
    self = [super init];
    if (self) {
        self.pp = pp;
    }
    return self;
}

- (BOOL)boolValue {
    if ([_pp respondsToSelector:@selector(boolValue)]) {
        return self.pp.boolValue;
    }
    return false;
}

- (double)floatValue {
    if ([_pp respondsToSelector:@selector(floatValue)]) {
        return self.pp.floatValue;
    } else if ([_pp respondsToSelector:@selector(nillableFloatValue)]) {
        NSNumber *v = self.pp.nillableFloatValue;
        if (v) {
            return v.doubleValue;
        }
    }
    return AppcoreLibPropertyProviderNilFloatValue;
}

- (int64_t)intValue {
    if ([_pp respondsToSelector:@selector(intValue)]) {
        return self.pp.intValue;
    } else if ([_pp respondsToSelector:@selector(nillableIntValue)]) {
        NSNumber *v = self.pp.nillableIntValue;
        if (v) {
            return v.longLongValue;
        }
    }
    return AppcoreLibPropertyProviderNilIntValue;
}

- (int64_t)timeEpochMilliseconds {
    NSDate *date = nil;
    if ([_pp respondsToSelector:@selector(dateValue)]) {
        date = self.pp.dateValue;
    }
    if (!date) {
        return AppcoreLibPropertyProviderNilIntValue;
    }
    int64_t epochMilliseconds = [@(floor([date timeIntervalSince1970] * 1000)) longLongValue];
    return epochMilliseconds;
}

- (NSString *_Nonnull)stringValue {
    if ([_pp respondsToSelector:@selector(stringValue)]) {
        NSString *v = self.pp.stringValue;
        if (v) {
            return v;
        }
    }
    return AppcoreLibPropertyProviderNilStringValue;
}

- (long)type {
    switch (self.pp.type) {
    case CMPropertyProviderTypeBool:
        return AppcoreLibPropertyProviderTypeBool;
    case CMPropertyProviderTypeString:
        return AppcoreLibPropertyProviderTypeString;
    case CMPropertyProviderTypeInt:
        return AppcoreLibPropertyProviderTypeInt;
    case CMPropertyProviderTypeFloat:
        return AppcoreLibPropertyProviderTypeFloat;
    case CMPropertyProviderTypeTime:
        return AppcoreLibPropertyProviderTypeTime;
    }

    return AppcoreLibPropertyProviderTypeBool;
}

@end
