//
//  CriticalMoments.h
//  CriticalMoments
//
//  Created by Steve Cosman on 2023-04-17.
//

#import <Foundation/Foundation.h>

#import "../messaging/CMBannerManager.h"
#import "../messaging/CMBannerMessage.h"
#import "../themes/CMTheme.h"

NS_ASSUME_NONNULL_BEGIN

/**
 The primary interface into Critical Moments. See out getting started docs for usage instructions:
 https://docs.criticalmoments.io/get-started
 */
@interface CriticalMoments : NSObject

// Simple "ping" method for testing end to end integrations
/// :nodoc:
+ (NSString *)objcPing;

// Golang "ping" method for testing end to end integrations
/// :nodoc:
+ (NSString *)goPing;

/**
 Start should be called once you've performed all needed initialization for
 critical moments. Critical moments won't perform actions until it is started.
 This is typically called in AppDelegate didfinishlaunchingwithoptions, but can
 be anywhere you like, as long as the primary root view controler is already
 rendering when you call start.

 Initializtion that should be performed before calling start:

 - Set critical moments API key (required)
 - Set critical moments config URLs (highly recomended)
 - Setup a default theme from code (optional). Can also be done through config
 or not at all.
 */
+ (void)start;

/**
 Set the API Key for critical moments.

 You can get a valid API key from criticalmoments.io

 API Keys are not transferable; each app requires it's own key.

 @param apiKey the API Key. Create one on criticalmoments.io
 @param error optional, any error created when validating the API key
 */
+ (void)setApiKey:(NSString *)apiKey error:(NSError **)error;

/**
 Set the config URL for critical moments.

 We highly recommend https/web URLs, as Critical Moments is particularly useful
 for unexpected/unplanned customer messaging. With a remote URL you can update
 the config to handle these situations. Loading from a file in the bundle is
 supported, but mostly for testing.

 @param urlString the URL string of the json config file. Can be a local
 `file://` URL or a `https://` URL.
 @warning Be sure to secure who can upload files to this URL path. This config
 file can present messages directly to your users, and you should treat security
 seriously, as you would your app update release process or webpage.
 */
+ (void)setConfigUrl:(NSString *)urlString;

// TODO: improve docs
// TODO: enforce naming limits (ascii, no spaces)?
/**
 Use SendEvent to sent named events to Critical Moments (example:
 `user_updated_profile_photo`). These events may trigger actions, or may be used
 in conditions.

 @param eventName a string describing the event. Example:
 `user_updated_profile_photo`
 */
+ (void)sendEvent:(NSString *)eventName;

/**
 Check a  condition string, returning the result of evaluating it.

 A name is provided so that you can remotely override the condition string using a cloud based config file.

 @param name A name for this condition. Must be provided and can not be an empty string.
 The name allows you to override the hardcoded condition string remotely from the cloud-hosted
 CM config file later, if your business needs change
 @param condition The condition string, for example: "interface_orientation == 'landscape'". See documentation on
options here: https://docs.criticalmoments.io/conditional-targeting/intro-to-conditions
 @param error Any errors returned from evaluating the condition.
 @return The result of evaluating the condition. Always false for an error.

 */
+ (bool)checkNamedCondition:(NSString *)name condition:(NSString *)condition error:(NSError **)error;

@end

NS_ASSUME_NONNULL_END
