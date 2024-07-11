//
//  CMBackgroundHandler.h
//
//
//  Created by Steve Cosman on 2024-07-11.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

@interface CMBackgroundHandler : NSObject

+ (void)registerBackgroundTasks;
+ (void)scheduleBackgroundTask;

@end

NS_ASSUME_NONNULL_END
