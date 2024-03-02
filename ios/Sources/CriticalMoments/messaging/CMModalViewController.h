//
//  CMSheetViewController.h
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import <UIKit/UIKit.h>

#import "../themes/CMTheme.h"

@import Appcore;

NS_ASSUME_NONNULL_BEGIN

@interface CMModalViewController : UIViewController

- (instancetype)initWithDatamodel:(DatamodelModalAction *)model;

- (instancetype)init NS_UNAVAILABLE;

@end

NS_ASSUME_NONNULL_END
