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

@interface CMButton : UIButton

// TODO Private? Or confirm not exported?

/// :nodoc:
- (instancetype)init NS_UNAVAILABLE;

// TODO: actual CMButton that contains a UIButton? init method, callbacks, ... yeah.

/// :nodoc:
+ (UIButton *)buttonWithWithDataModel:(DatamodelButton *)model andTheme:(CMTheme *_Nullable)theme;

@end

NS_ASSUME_NONNULL_END
