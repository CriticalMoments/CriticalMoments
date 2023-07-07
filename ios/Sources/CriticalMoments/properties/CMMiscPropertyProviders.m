//
//  CMMiscPropertyProviders.m
//
//
//  Created by Steve Cosman on 2023-07-07.
//

#import "CMMiscPropertyProviders.h"

#import "../utils/CMUtils.h"

@implementation CMAppInstallDatePropertyProviders

- (int64_t)intValue {
    NSURL *docsFolderUrl = [[[NSFileManager defaultManager] URLsForDirectory:NSDocumentDirectory
                                                                   inDomains:NSUserDomainMask] lastObject];
    if (!docsFolderUrl) {
        return 0;
    }

    NSError *err;
    NSDate *appInstallDate = [[[NSFileManager defaultManager] attributesOfItemAtPath:docsFolderUrl.path error:&err]
        objectForKey:NSFileCreationDate];

    if (err != nil || appInstallDate == nil) {
        return 0;
    }

    return [CMUtils cmTimestampFromDate:appInstallDate];
}

- (long)type {
    return AppcoreLibPropertyProviderTypeInt;
}

@end
