//
//  CMButtonStack.h
//
//
//  Created by Steve Cosman on 2023-06-15.
//

#import <UIKit/UIKit.h>

#import "../themes/CMTheme.h"

@import Appcore;

NS_ASSUME_NONNULL_BEGIN

@interface CMButton : UIView

/// :nodoc:
- (instancetype)init NS_UNAVAILABLE;

/// :nodoc:
- (instancetype)initWithWithDataModel:(DatamodelButton *)model andTheme:(CMTheme *_Nullable)theme;

/// :nodoc:
// the "default" action, which won't be called if the model has preventDefault=true
@property(nonatomic, copy, nullable) void (^defaultAction)();

@end

NS_ASSUME_NONNULL_END
