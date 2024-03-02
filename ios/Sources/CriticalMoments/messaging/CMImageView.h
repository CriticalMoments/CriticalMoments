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

- (instancetype)initWithDatamodel:(DatamodelImage *)model andTheme:(CMTheme *)theme;

- (instancetype)init NS_UNAVAILABLE;

@end

NS_ASSUME_NONNULL_END
