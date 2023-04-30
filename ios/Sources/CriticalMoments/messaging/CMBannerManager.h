//
//  CMBannerManagger.h
//
//
//  Created by Steve Cosman on 2023-04-23.
//

#import <Foundation/Foundation.h>

#import "CMBannerMessage.h"

NS_ASSUME_NONNULL_BEGIN

/**
 Enumeration of CMBannerMannager positions
 @see CMBannerManager.appWideBannerPosition
 */
typedef NS_ENUM(NSUInteger, CMBannerPosition) {
    /// Position banners at the bottom of the key window
    CMBannerPositionBottom,
    /// Position banners at the top of the key window
    CMBannerPositionTop
};

/**
 Manages the presentation of banner messages across your app.

 Example usage Swift:
 ```swift
 let message = CMBannerMessage.init(body: "Important msg")
 CMBannerManager.shared().showAppWideMessage(message)
 ```

 Example usage Objective-C:
 ```objc
 CMBannerMessage* bannerMessage = [[CMBannerMessage alloc]
 initWithBody:@"Important msg"]; [CMBannerManager.shared
 showAppWideMessage:bannerMessage];
 ```
 */
@interface CMBannerManager : NSObject

#pragma mark Shared Instance

/**
 A shared instance reference. You should only use a single instance of
 CMBannerManager per app. This `shared` instance is available for convenience,
 but you can also create and maintain your own instance if you prefer.
 @return a shared instance of CMBannerManager
 */
+ (CMBannerManager *)shared;

#pragma mark Banner Position

/**
 Set this to configure if the banner messages are presented at the top or bottom
 of your app.
 @warning Be sure to test your app renders well with the chosen position. The
 banner manager will resize your root view controller to make room for the
 banner; if you've hard coded offsets, for example the notch or dyamic island,
 then your app layout may be a bit strange. If you encounter issues, you should
 adapt use apple layout guides for any offsets, which will solve most of these
 issues and help on future hardware.
 @see CMBannerPosition enumeration
 */
@property(nonatomic) CMBannerPosition appWideBannerPosition;

#pragma mark Displaying and removing banners

/**
 Shows a banner across your entire application, shifting the root viewcontroller
 of your key window. If called multiple times, will include UI to iterate though
 multiple banners.
 @param message the CMBannerMessage to present
 */
- (void)showAppWideMessage:(CMBannerMessage *)message;

/**
Removes a previously presented banner message.
 @param message the CMBannerMessage to remove
 */
- (void)removeAppWideMessage:(CMBannerMessage *)message;

/**
 Removes all  previously presented banner messages.
 */
- (void)removeAllAppWideMessages;

@end

NS_ASSUME_NONNULL_END
