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

- (instancetype)init NS_UNAVAILABLE;

- (instancetype)initWithWithDataModel:(DatamodelButton *)model andTheme:(CMTheme *_Nullable)theme;

/// the "default" action, which won't be called if the model has preventDefault=true
@property(nonatomic, copy, nullable) void (^defaultAction)();
/// button tapped action, always called
@property(nonatomic, copy, nullable) void (^buttonTappedAction)();

@end

NS_ASSUME_NONNULL_END
