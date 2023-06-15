//
//  CMImageView.m
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import "CMImageView.h"

@implementation CMImageView

- (instancetype)init {
    self = [super init];
    if (self) {
        [self buildSubviews];
    }
    return self;
}

- (CMTheme *)theme {
    if (self.customTheme) {
        return self.customTheme;
    }
    return CMTheme.current;
}

- (UIImage *)getImageFromDatamodel {
    // TODO case for options in data model. V1: symbol, built in.
    UIImage *image = [self imageForSymbolImage];

    if (!image) {
        // TODO -- get fallback image
    }
    return image;
}

- (void)buildSubviews {
    UIImage *image = [self getImageFromDatamodel];

    UIImageView *iv = [[UIImageView alloc] initWithImage:image];
    iv.tintColor = [self.theme primaryColorForView:self];
    iv.contentMode = UIViewContentModeScaleAspectFit;
    iv.translatesAutoresizingMaskIntoConstraints = NO;
    [self addSubview:iv];

    // Layout

    NSArray<NSLayoutConstraint *> *constraints = @[
        [iv.topAnchor constraintEqualToAnchor:self.topAnchor],
        [iv.leftAnchor constraintEqualToAnchor:self.leftAnchor],
        [iv.bottomAnchor constraintEqualToAnchor:self.bottomAnchor],
        [iv.rightAnchor constraintEqualToAnchor:self.rightAnchor],
    ];
    [NSLayoutConstraint activateConstraints:constraints];
}

/*
  - weight
  - color mode
    - default to primary color mono
    - Hierarchical with optional color override (single color with tints)
    - Palette: must provide colors
    - Multicolor
  - 1-3 colors (default to use primary)

 */

- (UIImage *)imageForSymbolImage {
    UIImage *image;
    if (@available(iOS 13.0, *)) {
        UIImageSymbolConfiguration *c = [UIImageSymbolConfiguration unspecifiedConfiguration];

        // TODO -- check propertes mentioned above
        // [UIImageSymbolConfiguration configurationWithPointSize:<#(CGFloat)#> weight:<#(UIImageSymbolWeight)#>
        // scale:<#(UIImageSymbolScale)#>];

        // TODO hardcode
        image = [UIImage systemImageNamed:@"square.and.pencil" withConfiguration:c];
    }
    return image;
}

@end
