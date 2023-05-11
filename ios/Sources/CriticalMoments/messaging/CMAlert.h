//
//  CMAlert.h
//
//
//  Created by Steve Cosman on 2023-05-11.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

@interface CMAlert : NSObject

// No public APIs. There is really no reason to use our alert
// over UIAlertController from your own code.

/// :nodoc:
- (instancetype)init NS_UNAVAILABLE;

@end

NS_ASSUME_NONNULL_END
