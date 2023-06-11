//
//  CMBannerManagger.m
//
//
//  Created by Steve Cosman on 2023-04-23.
//

#import "CMBannerManager.h"
#import "../utils/CMUtils.h"
#import "CMBannerMessage_private.h"

#define MAX_BANNER_HEIGHT_PERCENTAGE 0.20

@interface CMBannerManager () <CMBannerMessageManagerDelegate, CMBannerNextMessageDelegate>

// Access should be @synchronized(self)
@property(nonatomic, strong) NSMutableArray<CMBannerMessage *> *appWideMessages;
@property(nonatomic, strong) CMBannerMessage *currentMessage;

// access syncronized by main queue
@property(nonatomic, strong) UIView *appWideContainerView;

// Track our changes, so they can be reverted, access syncronized by main queue
@property(nonatomic, strong) UIViewController *injectedInsetsRootVc;
@property(nonatomic) UIEdgeInsets insetAddedForBanner;

@end

@implementation CMBannerManager

static CMBannerManager *sharedInstance = nil;

+ (CMBannerManager *)shared {
    // avoid lock if we can
    if (sharedInstance) {
        return sharedInstance;
    }

    @synchronized(CMBannerManager.class) {
        if (!sharedInstance) {
            sharedInstance = [[self alloc] init];
        }

        return sharedInstance;
    }
}

- (instancetype)init {
    self = [super init];
    if (self) {
        _appWideMessages = [[NSMutableArray alloc] init];
        // Default to bottom -- less likely to conflict with hard-coded app
        // frame content
        self.appWideBannerPosition = CMBannerPositionBottom;
    }
    return self;
}

- (void)showAppWideMessage:(CMBannerMessage *)message {
    @synchronized(self) {
        if (![_appWideMessages containsObject:message]) {
            message.messageManagerDelegate = self;
            [_appWideMessages addObject:message];
        }

        _currentMessage = message;

        // update the app wide banner position if message has a preference on
        // location
        bool rendered = NO;
        if (message.preferredPosition != CMBannerPositionNoPreference) {
            rendered = [self setAppWideBannerPositionReturningRendered:message.preferredPosition];
        }

        if (!rendered) {
            [self renderForCurrentState];
        }
    }
}

- (void)removeAppWideMessage:(CMBannerMessage *)message {
    @synchronized(self) {
        [_appWideMessages removeObject:message];

        if (_currentMessage == message) {
            _currentMessage = _appWideMessages.lastObject;
        }

        [self renderForCurrentState];
    }
}

- (void)removeAllAppWideMessages {
    @synchronized(self) {
        [_appWideMessages removeAllObjects];
        _currentMessage = nil;
        [self renderForCurrentState];
    }
}

- (void)setAppWideBannerPosition:(CMBannerPosition)appWideBannerPosition {
    [self setAppWideBannerPositionReturningRendered:appWideBannerPosition];
}

- (bool)setAppWideBannerPositionReturningRendered:(CMBannerPosition)appWideBannerPosition {
    // don't do work work if no change
    if (appWideBannerPosition == _appWideBannerPosition) {
        return NO;
    }
    // No preference isn't valid for manager. Keep prior value.
    if (appWideBannerPosition == CMBannerPositionNoPreference) {
        return NO;
    }
    _appWideBannerPosition = appWideBannerPosition;
    [self removeAppWideBannerContainer];
    [self renderForCurrentState];
    return YES;
}

- (void)renderForCurrentState {
    // Always dispatch async. The caller might call us before the window/rootVC
    // relations are setup.
    dispatch_async(dispatch_get_main_queue(), ^{
      [self renderForCurrentStateSyncMain];
    });
}

- (void)renderForCurrentStateSyncMain {
    if (![NSThread isMainThread]) {
        dispatch_sync(dispatch_get_main_queue(), ^{
          [self renderForCurrentState];
        });
        return;
    }

    // Ensure current message is valid
    // Prefer current, then the last added still active message
    if (![_appWideMessages containsObject:_currentMessage]) {
        _currentMessage = nil;
    }
    if (!_currentMessage) {
        _currentMessage = _appWideMessages.lastObject;
    }

    // if no messages left to render clear container view
    if (!_currentMessage) {
        [self removeAppWideBannerContainer];
        return;
    }

    // remove prior message from container
    for (UIView *priorMessage in _appWideContainerView.subviews) {
        [priorMessage removeFromSuperview];
    }

    [self createAppWideBannerContainerIfMissing];
    if (!_appWideContainerView) {
        return;
    }

    if (_appWideMessages.count > 1) {
        _currentMessage.nextMessageDelegate = self;
    } else {
        _currentMessage.nextMessageDelegate = nil;
    }

    _currentMessage.translatesAutoresizingMaskIntoConstraints = NO;
    [_appWideContainerView addSubview:_currentMessage];
    NSArray<NSLayoutConstraint *> *constraints = @[
        [_currentMessage.topAnchor constraintEqualToAnchor:_appWideContainerView.topAnchor],
        [_currentMessage.leftAnchor constraintEqualToAnchor:_appWideContainerView.leftAnchor],
        [_currentMessage.rightAnchor constraintEqualToAnchor:_appWideContainerView.rightAnchor],
        [_currentMessage.bottomAnchor constraintEqualToAnchor:_appWideContainerView.bottomAnchor],
    ];

    [NSLayoutConstraint activateConstraints:constraints];

    [self setInsetsForCurrentBanner:_currentMessage];
}

- (void)createAppWideBannerContainerIfMissing {
    if (![NSThread isMainThread]) {
        dispatch_sync(dispatch_get_main_queue(), ^{
          [self createAppWideBannerContainerIfMissing];
        });
        return;
    }

    if (_appWideContainerView) {
        return;
    }

    // Find key window, falling back to first window
    UIWindow *keyWindow = [CMUtils keyWindow];
    if (!keyWindow) {
        // no window to render in
        NSLog(@"CriticalMoments: CMBannerManager could not find a key window, "
              @"aborting.");
        return;
    }

    // Check we're really setup properly.
    // TODO: do this check on CM.start. If called too soon, we should just
    // dispatch later, not fail
    UIViewController *rootVc = keyWindow.rootViewController;
    if (rootVc.view.window != keyWindow) {
        NSLog(@"CriticalMoments: tried to show a banner before root VC was "
              @"setup, "
              @"aborting.");
        return;
    }

    // Add the banner container view to the key window
    _appWideContainerView = [[UIView alloc] init];
    [keyWindow addSubview:_appWideContainerView];

    //
    // Layout
    //

    _appWideContainerView.translatesAutoresizingMaskIntoConstraints = NO;

    NSArray<NSLayoutConstraint *> *constraints = @[
        // position banner to the side edges
        [_appWideContainerView.leftAnchor constraintEqualToAnchor:keyWindow.leftAnchor],
        [_appWideContainerView.rightAnchor constraintEqualToAnchor:keyWindow.rightAnchor],

        // Make the banner at most 20% window height. Backstop for way too much
        // text.
        [_appWideContainerView.heightAnchor constraintLessThanOrEqualToAnchor:keyWindow.heightAnchor
                                                                   multiplier:MAX_BANNER_HEIGHT_PERCENTAGE],
    ];

    // Top vs Bottom layout for banner
    if (self.appWideBannerPosition == CMBannerPositionBottom) {
        // Banner container at bottom of app
        constraints = [constraints
            arrayByAddingObject:[_appWideContainerView.bottomAnchor constraintEqualToAnchor:keyWindow.bottomAnchor]];
    } else {
        // Banner container at top of app
        constraints = [constraints
            arrayByAddingObject:[_appWideContainerView.topAnchor constraintEqualToAnchor:keyWindow.topAnchor]];
    }

    [NSLayoutConstraint activateConstraints:constraints];
}

- (void)removeAppWideBannerContainer {
    if (![NSThread isMainThread]) {
        dispatch_sync(dispatch_get_main_queue(), ^{
          [self removeAppWideBannerContainer];
        });
        return;
    }

    if (!_appWideContainerView) {
        return;
    }

    UIViewController *rootVc = _appWideContainerView.window.rootViewController;
    [_appWideContainerView removeFromSuperview];
    _appWideContainerView = nil;

    // reset insets to zero added for banner
    if (rootVc == _injectedInsetsRootVc) {
        UIEdgeInsets newInsets = rootVc.additionalSafeAreaInsets;
        newInsets.bottom = MAX(0, rootVc.additionalSafeAreaInsets.bottom - _insetAddedForBanner.bottom);
        newInsets.top = MAX(0, rootVc.additionalSafeAreaInsets.top - _insetAddedForBanner.top);
        rootVc.additionalSafeAreaInsets = newInsets;
        _insetAddedForBanner = UIEdgeInsetsZero;
        _injectedInsetsRootVc = nil;
    }
}

#pragma mark CMBannerMessageManagerDelegate

- (void)dismissedMessage:(CMBannerMessage *)message {
    [self removeAppWideMessage:message];
}

- (void)messageDidLayout:(CMBannerMessage *)message {
    [self setInsetsForCurrentBanner:message];
}

- (void)setInsetsForCurrentBanner:(CMBannerMessage *)message {
    if (![NSThread isMainThread]) {
        dispatch_sync(dispatch_get_main_queue(), ^{
          [self messageDidLayout:message];
        });
        return;
    }

    if (message != _currentMessage || !_appWideContainerView) {
        return;
    }

    // get current root VC. Some apps could change this, so always get current
    UIViewController *rootVc = _appWideContainerView.window.rootViewController;

    // if this is a new VC we haven't seen before, reset our _addedInsets to
    // zero
    if (rootVc != _injectedInsetsRootVc) {
        _insetAddedForBanner = UIEdgeInsetsZero;
        _injectedInsetsRootVc = nil;
    }

    // We want to have minimal impact to apps when setting the
    // additionalSafeAreaInsets New inset is:
    // - The height of our banner
    // - subtract the height of the safeAreaInset on that side, since we're
    // laying out over that and the app shouldn't double up the safe area
    // - Add back in any insets the app set for their own UI reasons.

    UIEdgeInsets newInsets = rootVc.additionalSafeAreaInsets;
    UIEdgeInsets newInsetsAddedForBanner = UIEdgeInsetsZero;

    if (self.appWideBannerPosition == CMBannerPositionBottom) {
        // reset top to zero added for banner
        newInsets.top = MAX(0, rootVc.additionalSafeAreaInsets.top - _insetAddedForBanner.top);
        newInsetsAddedForBanner.top = 0;

        // Calculate bottom inset
        CGFloat appAddedBottomInset = MAX(0, rootVc.additionalSafeAreaInsets.bottom - _insetAddedForBanner.bottom);
        newInsetsAddedForBanner.bottom = MAX(0, message.frame.size.height - message.safeAreaInsets.bottom);
        newInsets.bottom = MAX(0, newInsetsAddedForBanner.bottom + appAddedBottomInset);
    } else {
        // reset bottom to zero added for banner
        newInsets.bottom = MAX(0, rootVc.additionalSafeAreaInsets.bottom - _insetAddedForBanner.bottom);
        newInsetsAddedForBanner.bottom = 0;

        // Calculate bottom inset
        CGFloat appAddedTopInset = MAX(0, rootVc.additionalSafeAreaInsets.top - _insetAddedForBanner.top);
        newInsetsAddedForBanner.top = MAX(0, message.frame.size.height - message.safeAreaInsets.top);
        newInsets.top = MAX(0, newInsetsAddedForBanner.top + appAddedTopInset);
    }
    if (UIEdgeInsetsEqualToEdgeInsets(newInsets, rootVc.additionalSafeAreaInsets)) {
        // save a layout, no change
        return;
    }

    rootVc.additionalSafeAreaInsets = newInsets;
    _injectedInsetsRootVc = rootVc;
    _insetAddedForBanner = newInsetsAddedForBanner;
}

#pragma mark CMBannerNextMessageDelegate

- (void)nextMessage {
    @synchronized(self) {
        if (_appWideMessages.count == 0) {
            return;
        }
        NSUInteger nextIndex = [_appWideMessages indexOfObject:_currentMessage] + 1;
        nextIndex = nextIndex % _appWideMessages.count;
        _currentMessage = [_appWideMessages objectAtIndex:nextIndex];
        [self renderForCurrentState];
    }
}

@end
