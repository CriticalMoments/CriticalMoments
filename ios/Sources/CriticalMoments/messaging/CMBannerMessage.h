//
//  CMBannerMessage.h
//
//
//  Created by Steve Cosman on 2023-04-22.
//

#import <Foundation/Foundation.h>

@import UIKit;

#import "../themes/CMTheme.h"

NS_ASSUME_NONNULL_BEGIN

/**
 Enumeration of CMBannerMannager positions
 @see CMBannerManager.appWideBannerPosition
 @see CMBannerMessage.bannerPosition
 */
typedef NS_ENUM(NSUInteger, CMBannerPosition) {
    /// Position banners at the bottom of the key window
    CMBannerPositionBottom,
    /// Position banners at the top of the key window
    CMBannerPositionTop
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

 This class can be subclassed to implement completely custom views. Before
 subclassing, try using themes to achieve desired look. If you choose to
 subclass:
  - You should implement buildViewForMessage
  - Your view should have dismiss buttons, which targets the dismissTapped
 method
  - Your view should have a "next message" button, which targets the
 nextMessageButtonTapped method
  - Your view should call the actionDelegate on tap
  - Text should be loaded from the body property. If you hardcode, server driven
 messaging won't be available
  - You should respect showDismissButton and maxLineCount properties in your
 view
 */
@interface CMBannerMessage : UIView

#pragma mark Initializers

/// :nodoc:
- (instancetype)init NS_UNAVAILABLE;

/**
 @param body The body text to be rendered in the banner
 */
- (instancetype)initWithBody:(NSString *)body;

#pragma mark Position

/**
 The position to show this banner: the top or bottom of the screen.

 If multiple banners are presented, the position of the last banner shown is
 used for all banners. You can't have both bottom and top at the same time.

 @see CMBannerManagerappWideBannerPosition for changing the app wide banner
 position any time, independent of any single banner message.
 @warning Be sure to test your app renders well with the chosen position. The
 banner manager will resize your root view controller to make room for the
 banner; if you've hard coded offsets, for example the notch or dyamic island,
 then your app layout may be a bit strange. If you encounter issues, you should
 adapt use apple layout guides for any offsets, which will solve most of these
 issues and help on future hardware.
 */
@property(nonatomic) CMBannerPosition bannerPosition;

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
 Controlls the max number of lines allowed in the banner message text. Defaults
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

#pragma mark Subclassing

/**
 This method is only exposed for subclassing and should not be called in normal
 usage. See CMBannerMessage class documentation for notes on how to subclass
 properly.
 */
//- (UIView *)buildViewForMessage;

/**
 This method is only exposed for subclassing and should not be called in normal
 usage. See CMBannerMessage class documentation for notes on how to subclass
 properly.
 */
- (void)dismissTapped:(UIButton *)sender;

/**
 This method is only exposed for subclassing and should not be called in normal
 usage. See CMBannerMessage class documentation for notes on how to subclass
 properly.
 */
- (void)nextMessageButtonTapped:(UIButton *)sender;

@end

NS_ASSUME_NONNULL_END
