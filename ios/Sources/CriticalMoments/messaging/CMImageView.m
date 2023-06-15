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

- (void)buildSubviews {
    // TODO case for options. V1: symbol, built in.
    [self buildSubviewsForSymbolImage];
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

- (void)buildSubviewsForSymbolImage {
    if (@available(iOS 13.0, *)) {
        UIImageSymbolConfiguration *c = [UIImageSymbolConfiguration unspecifiedConfiguration];

        // TODO -- check propertes mentioned above
        // [UIImageSymbolConfiguration configurationWithPointSize:<#(CGFloat)#> weight:<#(UIImageSymbolWeight)#>
        // scale:<#(UIImageSymbolScale)#>];

        UIImage *image = [UIImage systemImageNamed:@"square.and.pencil" withConfiguration:c];
        UIImageView *iv = [[UIImageView alloc] initWithImage:image];
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
    } else {
        // TODO -- get fallback image
        return;
    }
}

@end
