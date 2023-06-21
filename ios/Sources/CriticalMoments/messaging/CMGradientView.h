//
//  CMGradientView.h
//
//
//  Created by Steve Cosman on 2023-06-16.
//

#import <UIKit/UIKit.h>

#import "../themes/CMTheme.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMGradientView : UIView

/// The custom theme to use for style. Defaults to the default theme.
@property(nonatomic, readwrite) CMTheme *customTheme;

@end

NS_ASSUME_NONNULL_END
