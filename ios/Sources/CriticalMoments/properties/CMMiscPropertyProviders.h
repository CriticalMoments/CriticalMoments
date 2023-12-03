//
//  CMMiscPropertyProviders.h
//
//
//  Created by Steve Cosman on 2023-07-07.
//

#import "CMBaseDynamicPropertyProvider.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMAppInstallDatePropertyProviders : NSObject <CMDynamicPropertyProvider>
@end

@interface CMHasWatchPropertyProviders : NSObject <CMDynamicPropertyProvider>
@end

@interface CMLanguageDirectionPropertyProvider : NSObject <CMDynamicPropertyProvider>
@end

NS_ASSUME_NONNULL_END
