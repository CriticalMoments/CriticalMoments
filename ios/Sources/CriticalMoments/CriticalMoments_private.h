//
//  CriticalMoments_private.h
//
//
//  Created by Steve Cosman on 2023-09-29.
//

#import "../CriticalMoments_private.h"
#import "../background/CMBackgroundHandler.h"
#import "../utils/CMEventSender.h"
#import "include/CriticalMoments.h"

NS_ASSUME_NONNULL_BEGIN

@import Appcore;

@interface CriticalMoments () <CMEventSender>

// _private header prevents exposing these to public SDK.

// Internal only -- use sharedSingleton in product
- (id)initInternal;

// Internal only -- use start in product
- (NSError *)startReturningError;

// Internal only -- for testing and built in event
- (void)sendEvent:(NSString *)eventName
          builtIn:(bool)builtIn
          handler:(void (^_Nullable)(NSError *_Nullable error))handler;

/// Access named themes
- (DatamodelTheme *)themeFromConfigByName:(NSString *)name;

/// Set the current theme for this CM instance.
/// Private, only for internal use (demo app).
- (void)setTheme:(CMTheme *)theme;

/// Fetch the current theme for this CM instance
/// Private, only for internal use (demo app).
- (CMTheme *)currentTheme;

/// Get API Key
- (nonnull NSString *)getApiKey;

/// This API is private, and should not be used externally. Use events + triggers to fire named events.
- (void)performNamedAction:(NSString *)name handler:(void (^_Nullable)(NSError *_Nullable error))handler;

/// Private API to process a CM notification when it's tapped
- (void)actionForNotification:(NSString *)identifier;

/// Private API to disable notification, as NSUserNotificationCenter isn't available in SPM tests
- (void)disableUserNotifications;
- (BOOL)userNotificationsDisabled;

/// Private API to perform appcore work in background
- (void)runAppcoreBackgroundWork;

/// Access the background handler, internal only, private API
@property(nonatomic, strong) CMBackgroundHandler *backgroundHandler;

// Private API to get the current notification plan
- (AppcoreNotificationPlan *)currentNotificationPlan:(NSError *_Nullable *_Nullable)error;

// Private API to update the current notification plan
- (void)updateNotificationPlan:(AppcoreNotificationPlan *_Nullable)notifPlan;

@end

NS_ASSUME_NONNULL_END
