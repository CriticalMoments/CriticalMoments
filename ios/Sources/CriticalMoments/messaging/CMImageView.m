//
//  CMImageView.m
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import "CMImageView.h"

#import "../themes/CMTheme_private.h"
#import "../utils/CMUtils.h"

#define CM_DEFAULT_IMAGE_SIZE 50

@interface CMImageView ()

@property(nonatomic, strong) DatamodelImage *model;
@property(nonatomic) CGFloat height;
@property(nonatomic, readwrite) CMTheme *customTheme;

@end

@implementation CMImageView

- (nonnull instancetype)initWithDatamodel:(nonnull DatamodelImage *)model andTheme:(CMTheme *)theme {
    self = [super init];
    if (self) {
        _model = model;
        _customTheme = theme;
        _height = 0.0; // zero until valid image loaded
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

- (UIImage *)getImageFromDatamodel:(DatamodelImage *)model {
    UIImage *image;

    if ([DatamodelImageTypeEnumSFSymbol isEqualToString:model.imageType]) {
        image = [self buildSymbolImage:model.symbolImageData];
    } else if ([DatamodelImageTypeEnumLocal isEqualToString:model.imageType]) {
        image = [self buildLocalImage:model.localImageData];
    }

    if (image) {
        self.height = model.height > 0 ? model.height : CM_DEFAULT_IMAGE_SIZE;
    }

    if (!image && model.fallback) {
        image = [self getImageFromDatamodel:model.fallback];
    }

    return image;
}

- (UIImage *)buildLocalImage:(DatamodelLocalImage *)model {
    return [UIImage imageNamed:model.path];
}

- (CGSize)intrinsicContentSize {
    // square
    return CGSizeMake(self.height, self.height);
}

- (void)buildSubviews {
    UIImage *image = [self getImageFromDatamodel:self.model];

    UIImageView *iv = [[UIImageView alloc] initWithImage:image];
    iv.contentMode = UIViewContentModeScaleAspectFit;
    iv.translatesAutoresizingMaskIntoConstraints = NO;
    [self addSubview:iv];

    // Layout

    NSArray<NSLayoutConstraint *> *constraints = @[
        [iv.topAnchor constraintEqualToAnchor:self.topAnchor],
        [iv.leftAnchor constraintEqualToAnchor:self.leftAnchor],
        [iv.bottomAnchor constraintEqualToAnchor:self.bottomAnchor],
        [iv.rightAnchor constraintEqualToAnchor:self.rightAnchor],
        [iv.heightAnchor constraintLessThanOrEqualToConstant:self.height],
    ];
    [NSLayoutConstraint activateConstraints:constraints];
}

- (UIImage *)buildSymbolImage:(DatamodelSymbolImage *)model {
    UIImage *image;
    if (@available(iOS 13.0, *)) {
        UIImageSymbolConfiguration *c = [UIImageSymbolConfiguration unspecifiedConfiguration];

        if (model.weight.length > 0) {
            UIImageSymbolWeight w = [self weightForConfigString:model.weight];
            c = [c configurationByApplyingConfiguration:[UIImageSymbolConfiguration configurationWithWeight:w]];
        }

        // Color priority: image data, custom theme, global theme, system tint, system default
        UIColor *primaryColor = [CMUtils colorFromHexString:model.primaryColor];
        if (!primaryColor) {
            primaryColor = [self.theme primaryColorForView:self];
        }
        UIColor *secondaryColor = [CMUtils colorFromHexString:model.secondaryColor];

        // Modes and color
        // This nighmare of a code block brought to you by the UIImageSymbolConfiguration API.
        bool hasiOS15 = NO;
        bool isMono = NO;
        if (@available(iOS 15.0, *)) {
            hasiOS15 = YES;
        }
        if (hasiOS15 && [DatamodelSystemSymbolModeEnumHierarchical isEqualToString:model.mode]) {
            // Hierarchical color
            if (@available(iOS 15.0, *)) {
                c = [c configurationByApplyingConfiguration:[UIImageSymbolConfiguration
                                                                configurationWithHierarchicalColor:primaryColor]];
            }
        } else if (hasiOS15 && secondaryColor && [DatamodelSystemSymbolModeEnumPalette isEqualToString:model.mode]) {
            // palette colors
            if (@available(iOS 15.0, *)) {
                c = [c
                    configurationByApplyingConfiguration:[UIImageSymbolConfiguration configurationWithPaletteColors:@[
                        primaryColor, secondaryColor
                    ]]];
            }
        } else {
            // Mono is default and fallback
            isMono = YES;
        }
        image = [UIImage systemImageNamed:model.symbolName withConfiguration:c];
        if (isMono) {
            // apply primary as tint (for mono only)
            image = [image imageWithTintColor:primaryColor renderingMode:UIImageRenderingModeAlwaysOriginal];
        }
    }

    return image;
}

- (UIImageSymbolWeight)weightForConfigString:(NSString *)s API_AVAILABLE(ios(13.0)) {
    if ([DatamodelSystemSymbolWeightEnumUltraLight isEqualToString:s]) {
        return UIImageSymbolWeightUltraLight;
    }

    if ([DatamodelSystemSymbolWeightEnumThin isEqualToString:s]) {
        return UIImageSymbolWeightThin;
    }

    if ([DatamodelSystemSymbolWeightEnumLight isEqualToString:s]) {
        return UIImageSymbolWeightLight;
    }

    if ([DatamodelSystemSymbolWeightEnumRegular isEqualToString:s]) {
        return UIImageSymbolWeightRegular;
    }

    if ([DatamodelSystemSymbolWeightEnumMedium isEqualToString:s]) {
        return UIImageSymbolWeightMedium;
    }

    if ([DatamodelSystemSymbolWeightEnumSemiBold isEqualToString:s]) {
        return UIImageSymbolWeightSemibold;
    }

    if ([DatamodelSystemSymbolWeightEnumBold isEqualToString:s]) {
        return UIImageSymbolWeightBold;
    }

    if ([DatamodelSystemSymbolWeightEnumHeavy isEqualToString:s]) {
        return UIImageSymbolWeightHeavy;
    }

    if ([DatamodelSystemSymbolWeightEnumBlack isEqualToString:s]) {
        return UIImageSymbolWeightBlack;
    }

    return UIImageSymbolWeightRegular;
}

@end
