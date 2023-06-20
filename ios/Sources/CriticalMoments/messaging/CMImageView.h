//
//  CMImageView.h
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import <UIKit/UIKit.h>

#import "../themes/CMTheme.h"

@import Appcore;

NS_ASSUME_NONNULL_BEGIN

@interface CMImageView : UIView

/// :nodoc:
- (instancetype)initWithDatamodel:(DatamodelImage *)model;

/// :nodoc:
- (instancetype)init NS_UNAVAILABLE;

// TODO: does this work if set after init?
/// The custom theme to use for style. Defaults to the default theme.
@property(nonatomic, readwrite) CMTheme *customTheme;

@end

NS_ASSUME_NONNULL_END
