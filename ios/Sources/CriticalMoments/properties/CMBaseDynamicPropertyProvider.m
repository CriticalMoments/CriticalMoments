//
//  CMDynamicPropertyProvider.m
//
//
//  Created by Steve Cosman on 2023-05-22.
//

#import "CMBaseDynamicPropertyProvider.h"

@import Appcore;

@implementation CMBaseDynamicPropertyProvider

- (BOOL)boolValue {
    return NO;
}

- (double)floatValue {
    return 0;
}

- (int64_t)intValue {
    return 0;
}

- (NSString *_Nonnull)stringValue {
    return nil;
}

- (long)type {
    return AppcoreLibPropertyProviderTypeString;
}

@end
