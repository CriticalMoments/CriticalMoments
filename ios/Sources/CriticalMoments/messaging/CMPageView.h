//
//  CMPageView.h
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import <UIKit/UIKit.h>

#import "../themes/CMTheme.h"

@import Appcore;

NS_ASSUME_NONNULL_BEGIN

@interface CMPageView : UIView

- (instancetype)initWithDatamodel:(DatamodelPage *)model andTheme:(CMTheme *)theme;

- (instancetype)init NS_UNAVAILABLE;

/// the "default" action, which will be called after any button tap
@property(nonatomic, copy, nullable) void (^anyButtonDefaultAction)();
/// Called after any button tap, for events
@property(nonatomic, copy, nullable) void (^buttonCallback)(NSString *_Nullable buttonName, int buttonIndex);

@end

NS_ASSUME_NONNULL_END
