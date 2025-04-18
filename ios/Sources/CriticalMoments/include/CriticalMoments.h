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

/// :nodoc:
+ (CriticalMoments *)sharedInstance;

/**
 The default instance of Critical Moments. You should always use this instance.
 */
+ (CriticalMoments *)shared;

#pragma mark Setup

/**
 Start should be called once you've performed all needed initialization for
 critical moments. Critical moments won't perform actions until it is started.
 This is typically called in AppDelegate didfinishlaunchingwithoptions, but can
 be anywhere you like, as long as the primary root view controler is already
 rendering when you call start.

 Initializtion that should be performed before calling start:

 - Set critical moments API key (required)
 - Set critical moments config URLs (required). See setDevelopmentConfigName: and setReleaseConfigUrl:
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

/**
 Set a local development config file for Critical Moments by name. Path will be looked up in your main bundle.

 For local development you may use a local and unsigned JSON config file built into the app binary. See the Quick Start
 guide for how to create this file: https://docs.criticalmoments.io/quick-start

 This local config file will not be used on release builds / app store builds. You must also set a production config URL
 with setProductionConfigUrl for those builds.

 @param configFileName the name of the config file (e.g. `cmConfig.json`). The full path will be looked up in your main
 bundle.
 */
- (void)setDevelopmentConfigName:(NSString *)configFileName;

/**
 Set a local development config URL for Critical Moments.

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
 seriously, as you would your app update release process or webpage security.
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

/**
 Set to true to log events and conditions to console, as they occur. Primarily uses in debugging, and should be disabled
 in release builds. Disabled by default.

 @param devMode  if true, the CM SDK will log events and conditions as they occur
 */
- (void)setDeveloperMode:(bool)devMode;

/// :nodoc: Replaced by setDeveloperMode which logs events, and more.
- (void)setLogEvents:(bool)logEvents __attribute__((deprecated("Use setDeveloperMode: instead")));

#pragma mark Feature Flags / Named Conditions

/// :nodoc: Our old format. Renamed the selectors for swift async bindings.
- (void)checkNamedCondition:(NSString *_Nonnull)name
                    handler:(void (^_Nonnull)(bool result, NSError *_Nullable error))handler;

/**
 Checks a named condition string, returning the result of evaluating it. The provided name is used to lookup a condition
in the config file's namedConditions section.

 Configuration documentation: https://docs.criticalmoments.io/conditional-targeting/named-conditions

 The result is returned through the provided handler asynchronously. The result is asynchronous because some conditions
can use properties which are asynchronous (checking network state, battery state, and many others).  It is not called on
the main thread, so be sure to dispatch to the main thread if calling into UI libraries.

 @param name The name of this condition. Must be provided and can not be an empty string.
 The name is used as a lookup the condition-string of a namedCondition in the config file.
 @param result Returns the boolean result of the condition evaluation. The boolean value is false for any error,
including if the condition is not found in the config. Also returns/throws any errors occurred evaluating the condition.
 @warning Be sure to provide a unique name to each use case. Reusing names (even if the current conditional logic is
currently equivalent) will make it impossible to override each usage independently from remote configuration.
 */
- (void)checkNamedCondition:(NSString *_Nonnull)name
          completionHandler:(void (^_Nonnull)(bool result, NSError *_Nullable error))result;

#ifdef IS_CRITICAL_MOMENTS_INTERNAL
// Private, only for internal use (demo app/testing).
// Evaluate a condition from a conditional-string
// Do not use this in any other apps. It's against the TOS, and will always return an error.
/// :nodoc:
- (void)checkInternalTestCondition:(NSString *_Nonnull)conditionString
                           handler:(void (^_Nonnull)(bool result, NSError *_Nullable error))handler;
#endif

#ifdef IS_CRITICAL_MOMENTS_INTERNAL
/// :nodoc: This API is private, and should not be used. Use events + triggers
- (void)performNamedAction:(NSString *)name handler:(void (^_Nullable)(NSError *_Nullable error))handler;
#endif

#pragma mark Themes

#ifdef IS_CRITICAL_MOMENTS_INTERNAL
// Fetch the current theme for this CM instance
// Private, only for internal use (demo app).
/// :nodoc:
- (CMTheme *)currentTheme;
// Set the current theme for this CM instance.
// Private, only for internal use (demo app).
/// :nodoc:
- (void)setTheme:(CMTheme *)theme;
// Set the current theme for this CM instance to a built in theme
// Private, only for internal use (demo app).
/// :nodoc:
- (void)setBuiltInTheme:(NSString *)themeName;
// Private, only for internal use (demo app).
/// :nodoc:
- (int)builtInBaseThemeCount;
#endif

#pragma mark Properties

/// :nodoc: Old format. New format below to take advantage of swift bindings.
- (void)registerStringProperty:(NSString *)value forKey:(NSString *)name error:(NSError *_Nullable *)error;

/**
 Set a custom or well-known string property for use in the CM condition engine.

 @param value The property value
 @param name The property key/name.  Can be used in conditions as "name" or "custom_name"
 @param error Any errors encountered setting the property
 @return True if call was successful. An error will be returned/thrown if false.
 */
- (BOOL)setStringProperty:(NSString *)value forKey:(NSString *)name error:(NSError *_Nullable *)error;

/// :nodoc: Old format. New format below to take advantage of swift bindings.
- (void)registerIntegerProperty:(long long)value forKey:(NSString *)name error:(NSError *_Nullable *)error;

/**
 Set a custom or well-known integer (int64) property for use in the CM condition engine.

 @param value The property value
 @param name The property key/name.  Can be used in conditions as "name" or "custom_name"
 @param error Any errors encountered setting the property
 @return True if call was successful. An error will be returned/thrown if false.
 */
- (BOOL)setIntegerProperty:(long long)value forKey:(NSString *)name error:(NSError *_Nullable *)error;

/// :nodoc: Old format. New format below to take advantage of swift bindings.
- (void)registerBoolProperty:(BOOL)value forKey:(NSString *)name error:(NSError *_Nullable *)error;

/**
 Set a custom or well-known boolean property for use in the CM condition engine.

 @param value The property value
 @param name The property key/name.  Can be used in conditions as "name" or "custom_name"
 @param error Any errors encountered setting the property
 @return True if call was successful. An error will be returned/thrown if false.
 */
- (BOOL)setBoolProperty:(BOOL)value forKey:(NSString *)name error:(NSError *_Nullable *)error;

/// :nodoc: Old format. New format below to take advantage of swift bindings.
- (void)registerFloatProperty:(double)value forKey:(NSString *)name error:(NSError *_Nullable *)error;

/**
 Set a custom or well-known floating point (double) property for use in the CM condition engine.

 @param value The property value
 @param name The property key/name.  Can be used in conditions as "name" or "custom_name"
 @param error Any errors encountered setting the property
 @return True if call was successful. An error will be returned/thrown if false.
 */
- (BOOL)setFloatProperty:(double)value forKey:(NSString *)name error:(NSError *_Nullable *)error;

/// :nodoc: Old format. New format below to take advantage of swift bindings.
- (void)registerTimeProperty:(NSDate *)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error;

/**
 Register a custom or well-known timestamp property (NSDate) for use in the CM condition engine.

 @param value The property value
 @param name The property key/name.  Can be used in conditions as "name" or "custom_name"
 @param error Any errors encountered setting the property
 @return True if call was successful. An error will be returned/thrown if false.
 */
- (BOOL)setTimeProperty:(NSDate *)value forKey:(NSString *)name error:(NSError *_Nullable __autoreleasing *)error;

/// :nodoc: Old format. New format below to take advantage of swift bindings.
- (void)registerPropertiesFromJson:(NSData *)jsonData error:(NSError *_Nullable __autoreleasing *)error;

/**
 Set a set of custom or well-known properties from JSON formatted data.

 The JSON object should be a single level JSON object, with string keys and bool, string or number values.

 On an issue, it will skip they problematic key/value pair, but continue on and parse as many supported key/value pairs
 as possible.

 All JSON number values are parsed into float64 values (including integers).

 @param jsonData The json data, in the format described above
 @param error Any errors encountered setting these properties. An error does not necessarily indicate that all fields
 failed, just that some field(s) failed.
 @return True if call was successful. An error will be returned/thrown if false.
 */
- (BOOL)setPropertiesFromJson:(NSData *)jsonData error:(NSError *_Nullable __autoreleasing *)error;

#pragma mark Notifications

/**
 Requests a person’s authorization to allow local and remote notifications for your app, using the standard system
 prompt to the user.

 This API calls the system's requestAuthorizationWithOptions. If the user approves authorization, Critical Moments will
 schedule any queued notifications.

 @param completionHandler Optional. In swift, return 2 values from an async call to this function.
 This is returned after permissions are granted or denied in the prompt.
 Returns two values. The first BOOL indicates if a permission prompt was shown. The second BOOL indicates if the user
 allowed notifications.
 */
- (void)requestNotificationPermissionWithCompletionHandler:
    (void (^_Nullable)(BOOL prompted, BOOL granted, NSError *__nullable error))completionHandler;

#ifdef IS_CRITICAL_MOMENTS_INTERNAL
// Simple "ping" method for testing end to end integrations
/// :nodoc:
- (NSString *)objcPing;

// Golang "ping" method for testing end to end integrations
/// :nodoc:
- (NSString *)goPing;

/// :nodoc: Private api for sample app.
- (void)removeAllBanners;
#endif

@end

NS_ASSUME_NONNULL_END
