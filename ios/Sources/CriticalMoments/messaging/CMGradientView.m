//
//  CMGradientView.m
//
//
//  Created by Steve Cosman on 2023-06-16.
//

#import "CMGradientView.h"

@interface CMGradientView ()

@property(nonatomic, readwrite) CAGradientLayer *gradient;

@end

@implementation CMGradientView

- (instancetype)init {
    self = [super init];
    if (self) {
        self.gradient = [CAGradientLayer layer];
        [self setGradientColorsForTraits];
        [self.layer insertSublayer:self.gradient atIndex:0];
    }
    return self;
}

- (CMTheme *)theme {
    if (self.customTheme) {
        return self.customTheme;
    }
    return CMTheme.current;
}

- (void)setCustomTheme:(CMTheme *)customTheme {
    _customTheme = customTheme;
    [self setGradientColorsForTraits];
}

- (void)layoutSubviews {
    [super layoutSubviews];
    self.gradient.frame = self.bounds;
}

- (void)traitCollectionDidChange:(UITraitCollection *)previousTraitCollection {
    [super traitCollectionDidChange:previousTraitCollection];
    [self setGradientColorsForTraits];
}

- (void)setGradientColorsForTraits {
    UIColor *clearBgColor = [self.theme.backgroundColor colorWithAlphaComponent:0.0];
    self.gradient.colors = @[ (id)clearBgColor.CGColor, (id)self.theme.backgroundColor.CGColor ];
}

@end
