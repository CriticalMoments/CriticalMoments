//
//  CMActionDispatcher.h
//
//
//  Created by Steve Cosman on 2023-05-05.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

@interface CMLibBindings : NSObject

#pragma mark Shared Instance

/**
 :nodoc:
 A shared instance reference.
 @return a shared instance of CMActionDispatcher
 */
+ (CMLibBindings *)shared;

/**
 :nodoc:
 Register the shared instance with appcore
 */
+ (void)registerWithAppcore;

@end

NS_ASSUME_NONNULL_END
