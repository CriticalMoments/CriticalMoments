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
typedef NS_ENUM(NSUInteger, CMAppWideBannerPosition) {
    /// Position banners at the bottom of the key window
    CMAppWideBannerPositionBottom,
    /// Position banners at the top of the key window
    CMAppWideBannerPositionTop
};

/**
 Manages the presentation of banner messages across your app.
 
 Example usage Swift:
 ```swift
 let message = CMBannerMessage.init(body: "Important msg")
 CMBannerManager.sharedInstance().showAppWideMessage(message)
 ```
 
 Example usage Objective-C:
 ```objc
 CMBannerMessage* bannerMessage = [[CMBannerMessage alloc] initWithBody:@"Important msg"];
 [CMBannerManager.sharedInstance showAppWideMessage:bannerMessage];
 ```
 */
@interface CMBannerManager : NSObject

/**
 Set this to configure if the banner messages are presented at the top or bottom of your app.
 @warning Be sure to test your app renders well with the chosen position. The banner manager will resize your root view controller to make room for the banner; if you've hard coded offsets, for example the notch or dyamic island, then your app layout may be a bit strange. If you encounter issues, you should adapt use apple layout guides for any offsets, which will solve most of these issues and help on future hardware.
 @see CMAppWideBannerPosition enum
 */
@property (nonatomic) CMAppWideBannerPosition appWideBannerPosition;

/**
 A shared instance reference. You should only use a single instance of CMBannerManager per app. This sharedInstance is suggested for convenience, but you can also create and maintain your own instance if you prefer.
 @return a shared instance of CMBannerManager
 */
+(CMBannerManager*) sharedInstance;

/**
 Shows a banner across your entire application, shifting the root viewcontroller of your key window. If called multiple times, will include UI to iterate though multiple banners.
 @param message the CMBannerMessage to present
 */
-(void) showAppWideMessage:(CMBannerMessage*)message;

/**
Removes a previously presented banner message.
 @param message the CMBannerMessage to remove
 */
-(void) removeAppWideMessage:(CMBannerMessage*)message;

/**
 Removes all  previously presented banner messages.
 */
-(void) removeAllAppWideMessages;

@end

NS_ASSUME_NONNULL_END
