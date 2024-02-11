//
//  CriticalMoments.h
//  CriticalMoments
//
//  Created by Steve Cosman on 2023-04-17.
//

#import <Foundation/Foundation.h>

#import "../themes/CMTheme.h"

NS_ASSUME_NONNULL_BEGIN

/**
 The primary interface into Critical Moments. See out getting started docs for usage instructions:
 https://docs.criticalmoments.io/get-started
 */
@interface CriticalMoments : NSObject

/// init is not available. Use sharedInstance for all use cases.
- (instancetype)init NS_UNAVAILABLE;

/**
 The default instance of critical moments. You should always use this instance
 */
+ (CriticalMoments *)sharedInstance;

#pragma mark Setup

/**
 Start should be called once you've performed all needed initialization for
 critical moments. Critical moments won't perform actions until it is started.
 This is typically called in AppDelegate didfinishlaunchingwithoptions, but can
 be anywhere you like, as long as the primary root view controler is already
 rendering when you call start.

 Initializtion that should be performed before calling start:

 - Set critical moments API key (required)
 - Set critical moments config URLs (required). See setDevelopmentConfigUrl: and setReleaseConfigUrl:
 - Setup a default theme from code (optional). Can also be done through config.
 or not at all.
 */
- (void)start;

/**
 Set the API Key for critical moments.

 You can get a valid API key from criticalmoments.io

 API Keys are not transferable; each app requires it's own key.

 @param apiKey the API Key. Create one on criticalmoments.io
 @param error optional, any error created when validating the API key
 */
- (void)setApiKey:(NSString *)apiKey error:(NSError **)error;

/// :nodoc:
- (nonnull NSString *)getApiKey;

/**
 Set a local development Config URL for critical moments.

 For local development you may use a local and unsigned JSON config file built into the app binary. See the Quick Start
 guide for how to create this file: https://docs.criticalmoments.io/quick-start

 This local config file will not be used on release builds / app store builds. You must also set a production config URL
 with setProductionConfigUrl for those builds.

 @param urlString the URL string of the json config file. Should begin with `file://`
 */
- (void)setDevelopmentConfigUrl:(NSString *)urlString;

/**
 Set the config URL for Critical Moments to be used on release builds / app store builds

 This url should begin with `https://`, and should link to a signed Critical Moments configuration file. See the docs
 for details: https://docs.criticalmoments.io/config-file

 @param urlString the URL string of the json config file. Should begin with `https://`
 @warning Be sure to secure who can upload files to this URL path. This config
 file can present messages directly to your users, and you should treat security
 seriously, as you would your app update release process or webpage secuirty.
 */
- (void)setReleaseConfigUrl:(NSString *)urlString;

#pragma mark Events

/**
 Use SendEvent to sent a named events to Critical Moments (example:
 `user_updated_profile_photo`). These events may trigger actions, or may be used
 in conditions.

 @param eventName a string describing the event. Example:
 `user_updated_profile_photo`
 */
- (void)sendEvent:(NSString *)eventName;

#pragma mark Feature Flags / Named Conditions

/**
 Checks a condition string, returning the result of evaluating it.

 A name is provided so that you can remotely override the condition string using a cloud based config file.

 The result is returned through the provided handler asynchronously. The result is asynchronous because some conditions
can use properties which are asyncronous (checking network state, battery state, and many others).  It is not called on
the main thread, so be sure to dispatch to the main thread if calling into UI libraries.

 @param name A name for this condition. Must be provided and can not be an empty string.
 The name allows you to override the hardcoded condition string remotely from the cloud-hosted
 CM config file later if needed.
 @param condition The condition string, for example: "interface_orientation == 'landscape'". See documentation on
options here: https://docs.criticalmoments.io/conditional-targeting/intro-to-conditions
 @param handler A callback block which will be called async with the boolean result of the condition evaluation. It also
returns any errors occured evaluating the condition. The boolean value is false for any error.
 @warning Be sure to provide a unique name to each condition you use. Reusing names will make it impossible to override
each usage independently from remote configuration. Reused names will log warnings in the debug console.
 */
- (void)checkNamedCondition:(NSString *_Nonnull)name
                  condition:(NSString *_Nonnull)condition
                    handler:(void (^_Nonnull)(bool result, NSError *_Nullable error))handler;

/// :nodoc: This API is private, and should not be used
- (void)performNamedAction:(NSString *)name handler:(void (^_Nullable)(NSError *_Nullable error))handler;

#pragma mark Themes

// Fetch the current theme for this CM instance
// Private, only for internal use (demo app).
/// :nodoc:
- (CMTheme *)currentTheme;
// Set the current theme for this CM instance.
// Private, only for internal use (demo app).
/// :nodoc:
- (void)setTheme:(CMTheme *)theme;

#pragma mark Properties

/**
 Register a custom or well-known string property for use in the CM condition engine.

 @param value The property value
 @param name The property key/name.  Can be used in conditions as "name" or "custom_name"
 @param error Any errors encountered setting the property
 */
- (void)registerStringProperty:(NSString *)value forKey:(NSString *)name error:(NSError *_Nullable *)error;
/**
 Register a custom or well-known integer (int64) property for use in the CM condition engine.

 @param value The property value
 @param name The property key/name.  Can be used in conditions as "name" or "custom_name"
 @param error Any errors encountered setting the property
 */
- (void)registerIntegerProperty:(long long)value forKey:(NSString *)name error:(NSError *_Nullable *)error;

/**
 Register a custom or well-known boolean property for use in the CM condition engine.

 @param value The property value
 @param name The property key/name.  Can be used in conditions as "name" or "custom_name"
 @param error Any errors encountered setting the property
 */
- (void)registerBoolProperty:(BOOL)value forKey:(NSString *)name error:(NSError *_Nullable *)error;

/**
 Register a custom or well-known floating point (double) property for use in the CM condition engine.

 @param value The property value
 @param name The property key/name.  Can be used in conditions as "name" or "custom_name"
 @param error Any errors encountered setting the property
 */
- (void)registerFloatProperty:(double)value forKey:(NSString *)name error:(NSError *_Nullable *)error;

/**
 Register a custom or well-known timestamp property (NSDate) for use in the CM condition engine.

 @param value The property value
 @param name The property key/name.  Can be used in conditions as "name" or "custom_name"
 @param error Any errors encountered setting the property
 */
- (void)registerTimeProperty:(NSDate *)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error;

/**
 Register a set of custom or well-known properties from JSON formatted data.

 The JSON object should be a single level JSON object, with string keys and bool, string or number values.

 On an issue, it will skip they problematic key/value pair, but continue on and parse as many supported key/value pairs
 as possible.

 All JSON number values are parsed into float64 values (including integers).

 @param jsonData The json data, in the format described above
 @param error Any errors encountered setting these properties. An error does not necessarily indicate that all fields
 failed, just that some field(s) failed.
 */
- (void)registerPropertiesFromJson:(NSData *)jsonData error:(NSError *_Nullable __autoreleasing *)error;

// Simple "ping" method for testing end to end integrations
/// :nodoc:
- (NSString *)objcPing;

// Golang "ping" method for testing end to end integrations
/// :nodoc:
- (NSString *)goPing;

/// :nodoc: Private api for sample app.
- (void)removeAllBanners;

@end

NS_ASSUME_NONNULL_END
