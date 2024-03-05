//
//  Header.h
//
//
//  Created by Steve Cosman on 2024-03-04.
//

#ifndef CMEventSenderHeader_h
#define CMEventSenderHeader_h

NS_ASSUME_NONNULL_BEGIN

@protocol CMEventSender <NSObject>

- (void)sendEvent:(NSString *)eventName;

@end

NS_ASSUME_NONNULL_END

#endif /* CMEventSenderHeader_h */
