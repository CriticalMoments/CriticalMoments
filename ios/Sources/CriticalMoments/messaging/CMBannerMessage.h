//
//  CMBannerMessage.h
//  
//
//  Created by Steve Cosman on 2023-04-22.
//

#import <Foundation/Foundation.h>

@import UIKit;

NS_ASSUME_NONNULL_BEGIN

@class CMBannerMessage;

/**
 A protocol for the action to call when a banner message is tapped.
 @param message the CMBannerMessage which was tapped
 @see CMBannerMessage.actionDelegate
 */
@protocol CMBannerActionDelegate
-(void) messageAction:(CMBannerMessage*)message;
@end

/**
 This class represents individual banner message which can be presented across the top or bottom of your app.
 
 It will take the default theme, which controls the font and colors.
 
 This class can be subclassed to implement completely custom views. Before subclassing, try using themes to achieve desired look. If you choose to subclass:
  - You should implement buildViewForMessage
  - Your view should have dismiss buttons, which targets the dismissTapped method
  - Your view should have a "next message" button, which targets the nextMessageButtonTapped method
  - Your view should call the actionDelegate on tap
  - Text should be loaded from the body property. If you hardcode, server driven messaging won't be available
  - You should respect showDismissButton and maxLineCount properties in your view
 */
@interface CMBannerMessage : NSObject

/**
 The body text to be rendered in the banner
 */
@property (nonatomic, readonly) NSString* body;

/**
 Controls if a dismiss button ("X") is rendered, allowing user to close the banner. Defaults to YES.
 */
@property (nonatomic, readwrite) bool showDismissButton;

/**
 Controlls the max number of lines allowed in the banner message text. Defaults to 4. Set to 0 for unlimited, although banners are still capped at 20% of screen height and will eventually tuncate.
 */
@property (nonatomic, readwrite) NSNumber* maxLineCount;

/**
 This action delegate will be called when the banner message is tapped. If not set, tapping will not perform any action.
 */
@property (nonatomic, readwrite) id<CMBannerActionDelegate> actionDelegate;

-(instancetype)init NS_UNAVAILABLE;

/**
 @param body The body text to be rendered in the banner
 */
-(instancetype)initWithBody:(NSString*)body;

/**
 This method is only exposed for subclassing and should not be called in normal usage. See CMBannerMessage class documentation for notes on how to subclass properly.
 */
-(UIView*) buildViewForMessage;

/**
 This method is only exposed for subclassing and should not be called in normal usage. See CMBannerMessage class documentation for notes on how to subclass properly.
 */
- (void)dismissTapped:(UIButton*)sender;

/**
 This method is only exposed for subclassing and should not be called in normal usage. See CMBannerMessage class documentation for notes on how to subclass properly.
 */
- (void)nextMessageButtonTapped:(UIButton*)sender;

@end

NS_ASSUME_NONNULL_END
