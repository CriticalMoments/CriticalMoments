//
//  CMAudioPropertyProvider.h
//
//
//  Created by Steve Cosman on 2023-05-26.
//

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMAudioPlayingPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

@interface CMAudioPortPropertyProvider : NSObject <CMDynamicPropertyProvider>

+ (CMAudioPortPropertyProvider *)hasHeadphones;
+ (CMAudioPortPropertyProvider *)hasWiredHeadset;
+ (CMAudioPortPropertyProvider *)hasBtHeadphones;
+ (CMAudioPortPropertyProvider *)hasBtHeadset;
+ (CMAudioPortPropertyProvider *)hasCarAudio;

@end

NS_ASSUME_NONNULL_END
