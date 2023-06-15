//
//  CMPageView.h
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import <UIKit/UIKit.h>

#import "../themes/CMTheme.h"

NS_ASSUME_NONNULL_BEGIN

// TODO: confirm not importable in apps / is internal

@interface CMPageView : UIView

/// The custom theme to use for the style of this page. Defaults to the default theme.
@property(nonatomic, readwrite) CMTheme *customTheme;

@end

NS_ASSUME_NONNULL_END
