//
//  CMBannerMessage.h
//  
//
//  Created by Steve Cosman on 2023-04-22.
//

#import <Foundation/Foundation.h>

@import UIKit;

NS_ASSUME_NONNULL_BEGIN

// Properties
// DefaultView: title, body, theme
// CustomView: new class with UIView, and delegate to dismiss
// Tap action: delegate


/*
 manager: high level helper for showing/hidding/adding message to app. API supports multi message, even if not done in
 Message: individual message. Class with default implementation, and you can override buildSubviews or something like that.
    - provide insets set for default dismiss button and "<>".
    - Draw dismiss and "<>" by default, but overridding customDismissButton or customNextButton or customPrevButton will prevent these from getting defaults
    - has a delegate for what happens on tap/dismiss
 Delegate: tap/dismiss
 */

@class CMBannerMessage;

@protocol CMBannerActionDelegate
-(void) messageAction:(CMBannerMessage*)message;
@end

@protocol CMBannerDismissDelegate
-(void) dismissedMessage:(CMBannerMessage*)message;
@end

@interface CMBannerMessage : NSObject

@property (nonatomic, readonly) NSString* body;
@property (nonatomic, readwrite) bool showDismissButton;
@property (nonatomic, readwrite) id<CMBannerActionDelegate> actionDelegate;
@property (nonatomic, readwrite) id<CMBannerDismissDelegate> dismissDelegate;

-(instancetype)init NS_UNAVAILABLE;

-(instancetype)initWithBody:(NSString*)body;

-(UIView*) buildViewForMessage;

// TODO: preferred size for height method
// TODO: next/prev delegate

@end

NS_ASSUME_NONNULL_END
