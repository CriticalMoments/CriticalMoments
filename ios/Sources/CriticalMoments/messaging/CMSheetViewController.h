//
//  CMSheetViewController.h
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import <UIKit/UIKit.h>

#import "../themes/CMTheme.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMSheetViewController : UIViewController

// TODO pretty much all private

/// If set, allows the sheet to be dismissed via a close button or swiping it away. If false, the user must select one
/// of the buttons to dismiss the sheet. Defaults to true.
@property(nonatomic, readwrite) BOOL showCloseButton;

/// A custom theme, which will be used instead of the default theme if set
@property(nonatomic, readwrite) CMTheme *customTheme;

@end

NS_ASSUME_NONNULL_END
