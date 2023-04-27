//
//  CMBannerManagger.h
//  
//
//  Created by Steve Cosman on 2023-04-23.
//

#import <Foundation/Foundation.h>

#import "CMBannerMessage.h"

NS_ASSUME_NONNULL_BEGIN

typedef NS_ENUM(NSUInteger, CMAppWideBannerPosition) {
    CMAppWideBannerPositionBottom,
    CMAppWideBannerPositionTop
};

@interface CMBannerManager : NSObject

@property (nonatomic) CMAppWideBannerPosition appWideBannerPosition;

+(CMBannerManager*) sharedInstance;

-(void) showAppWideMessage:(CMBannerMessage*)message;
-(void) removeAppWideMessage:(CMBannerMessage*)message;
-(void) removeAllAppWideMessages;

@end

NS_ASSUME_NONNULL_END
