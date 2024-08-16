//
//  CMDynamicPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-05-22.
//

#import <Foundation/Foundation.h>

// This class hides the idiosyncrasies of gomobile.
// Expose a clean protocol for objective C - CMDynamicPropertyProvider
// And the wrapper implements AppcoreLibPropertyProvider

typedef NS_ENUM(NSUInteger, CMPropertyProviderType) {
    CMPropertyProviderTypeBool,
    CMPropertyProviderTypeString,
    CMPropertyProviderTypeInt,
    CMPropertyProviderTypeFloat,
    CMPropertyProviderTypeTime,
};

@protocol CMDynamicPropertyProvider <NSObject>
- (CMPropertyProviderType)type;
@optional
- (BOOL)boolValue;
- (double)floatValue;
- (NSNumber *_Nullable)nillableFloatValue;
- (int64_t)intValue;
- (NSNumber *_Nullable)nillableIntValue;
- (NSString *_Nullable)stringValue;
- (NSDate *_Nullable)dateValue;
@end

@import Appcore;

NS_ASSUME_NONNULL_BEGIN

// intended to be subclassed
@interface CMDynamicPropertyProviderWrapper : NSObject <AppcoreLibPropertyProvider>

- (instancetype)init NS_UNAVAILABLE;
- (instancetype)initWithPP:(id<CMDynamicPropertyProvider>)pp;

@end

NS_ASSUME_NONNULL_END
