//
//  CMBannerMessage_private.h
//  
//
//  Created by Steve Cosman on 2023-04-26.
//

NS_ASSUME_NONNULL_BEGIN

@protocol CMBannerDismissDelegate
-(void) dismissedMessage:(CMBannerMessage*)message;
@end

@protocol CMBannerNextMessageDelegate
-(void) nextMessage;
@end

@interface CMBannerMessage ()

// _private header prevents exposing these to public SDK.

// We want people to overrider CMBannerMessage buildViewForMessage without
// breaking/overriding our delegation system for dismiss/next
@property (nonatomic, readwrite, nullable) id<CMBannerNextMessageDelegate> nextMessageDelegate;
@property (nonatomic, readwrite) id<CMBannerDismissDelegate> dismissDelegate;

@end

NS_ASSUME_NONNULL_END
