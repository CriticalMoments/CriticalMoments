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

@protocol CMBannerActionDelegate
-(void) messageAction:(CMBannerMessage*)message;
@end

@protocol CMBannerDismissDelegate
-(void) dismissedMessage:(CMBannerMessage*)message;
@end

@protocol CMBannerNextMessageDelegate
-(void) nextMessage;
@end

@interface CMBannerMessage : NSObject

@property (nonatomic, readonly) NSString* body;
@property (nonatomic, readwrite) bool showDismissButton;
@property (nonatomic, readwrite) NSNumber* maxLineCount;
@property (nonatomic, readwrite) id<CMBannerActionDelegate> actionDelegate;

-(instancetype)init NS_UNAVAILABLE;

-(instancetype)initWithBody:(NSString*)body;

-(UIView*) buildViewForMessage;

@end

NS_ASSUME_NONNULL_END
