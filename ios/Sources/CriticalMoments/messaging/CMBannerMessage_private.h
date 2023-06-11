//
//  CMBannerMessage_private.h
//
//
//  Created by Steve Cosman on 2023-04-26.
//

NS_ASSUME_NONNULL_BEGIN

@import Appcore;

@protocol CMBannerMessageManagerDelegate
- (void)dismissedMessage:(CMBannerMessage *)message;
- (void)messageDidLayout:(CMBannerMessage *)message;
@end

@protocol CMBannerNextMessageDelegate
- (void)nextMessage;
@end

@interface CMBannerMessage ()

// _private header prevents exposing these to public SDK.

/**
 :nodoc:
 @param bannerData The appcore datamodel for this banner
 */
- (instancetype)initWithAppcoreDataModel:(DatamodelBannerAction *)bannerData;

// We want people to overrider CMBannerMessage buildViewForMessage without
// breaking/overriding our delegation system for dismiss/next
@property(nonatomic, readwrite, nullable) id<CMBannerNextMessageDelegate> nextMessageDelegate;
@property(nonatomic, readwrite) id<CMBannerMessageManagerDelegate> messageManagerDelegate;

@end

NS_ASSUME_NONNULL_END
