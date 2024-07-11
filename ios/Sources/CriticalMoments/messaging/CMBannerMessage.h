//
//  CMBannerMessage.h
//
//
//  Created by Steve Cosman on 2023-04-22.
//

#import <Foundation/Foundation.h>

@import UIKit;

#import "../themes/CMTheme.h"
#import "../utils/CMEventSender.h"

NS_ASSUME_NONNULL_BEGIN

/**
 Enumeration of CMBannerMessage positions
 @see CMBannerMessage.bannerPosition
 @see CMBannerManager.appWideBannerPosition
 */
typedef NS_ENUM(NSUInteger, CMBannerPosition) {
    /// Position banners at the bottom of the key window (default)
    CMBannerPositionBottom,
    /// Position banners at the top of the key window
    CMBannerPositionTop,
    /// Positions banner at last used position, or default bottom
    CMBannerPositionNoPreference,
};

@class CMBannerMessage;

/**
 A protocol for the action which occurs when a banner message is tapped.
 @see CMBannerMessage.actionDelegate
 */
@protocol CMBannerActionDelegate
/**
 The delegate method which will be called when a banner is tapped
 @param message The CMBannerMessage which was tapped
 */
- (void)messageAction:(CMBannerMessage *)message;
@end

/**
 This class represents individual banner message which can be presented across
 the top or bottom of your app.
 */
@interface CMBannerMessage : UIView

#pragma mark Initializers

- (instancetype)init NS_UNAVAILABLE;

/**
 @param body The body text to be rendered in the banner
 */
- (instancetype)initWithBody:(NSString *)body;

#pragma mark Position

/**
 The preferred position to show this banner: the top or bottom of the screen.

 By default set to "no preference" and message will use the last banner position
 used, or system default.

 If multiple banners are presented, the position of the last banner shown is
 used for all banners. You can't have both bottom and top at the same time.

 @see CMBannerManagerappWideBannerPosition for changing the app wide banner
 position any time, independent of any single banner message's preference.
 @warning Be sure to test your app renders well with the chosen position. The
 banner manager will add to your root view controller's
 `additionalSafeAreaInsets` to make room for the banner; if you've hard coded
 offsets, for example the notch or dyamic island, then your app layout may be a
 bit strange. If you encounter issues, you should adapt use apple layout guides
 / safe-area-insets for any offsets, which will solve most of these issues and
 help on future hardware.
 */
@property(nonatomic) CMBannerPosition preferredPosition;

#pragma mark Theme

/// A custom theme, which will be used instead of the default theme if set
@property(nonatomic, readwrite) CMTheme *customTheme;

#pragma mark Content/Text

/**
 The body text to be rendered in the banner
 */
@property(nonatomic, readonly) NSString *body;

#pragma mark Configuration Options

/**
 Controls if a dismiss button ("X") is rendered, allowing user to close the
 banner. Defaults to YES.
 */
@property(nonatomic, readwrite) bool showDismissButton;

/**
 Controls the max number of lines allowed in the banner message text. Defaults
 to 4. Set to 0 for unlimited, although banners are still capped at 20% of
 screen height and will eventually tuncate.
 */
@property(nonatomic, readwrite) NSNumber *maxLineCount;

#pragma mark Action delegate

/**
 This action delegate will be called when the banner message is tapped. If not
 set, tapping will not perform any action.
 */
@property(nonatomic, readwrite) id<CMBannerActionDelegate> actionDelegate;

/// For sending events when the banner is tapped or closed
@property(nonatomic, weak, readwrite) id<CMEventSender> completionEventSender;
@property(nonatomic, strong, readwrite) NSString *bannerName;

@end

NS_ASSUME_NONNULL_END
