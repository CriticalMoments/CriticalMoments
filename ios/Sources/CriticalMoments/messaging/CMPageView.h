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

/// A custom theme, which will be used instead of the default theme if set
@property(nonatomic, readwrite) CMTheme *customTheme;

@end

NS_ASSUME_NONNULL_END
