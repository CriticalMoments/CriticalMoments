//
//  CMBannerManagger.m
//  
//
//  Created by Steve Cosman on 2023-04-23.
//

#import "CMBannerManager.h"
#import "CMBannerMessage_private.h"

#define MAX_BANNER_HEIGHT_PERCENTAGE 0.20

@interface CMBannerManager () <CMBannerDismissDelegate, CMBannerNextMessageDelegate>

// Access should be @synchronized(self)
@property (nonatomic, strong) NSMutableArray<CMBannerMessage*>* appWideMessages;
@property (nonatomic, strong) CMBannerMessage* currentMessage;

// currentMessageView managed by renderForCurrentState
@property (nonatomic, weak) UIView* currentMessageView;

// access syncronized by main queue
@property (nonatomic, strong) UIView* appWideContainerView;

@end

@implementation CMBannerManager

static CMBannerManager *sharedInstance = nil;

+ (CMBannerManager*)sharedInstance
{
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

-(instancetype)init {
    self = [super init];
    if (self) {
        _appWideMessages = [[NSMutableArray alloc] init];
        // Default to bottom -- less likely to conflict with hard-coded app frame content
        self.appWideBannerPosition = CMAppWideBannerPositionBottom;
    }
    return self;
}

-(void) showAppWideMessage:(CMBannerMessage*)message {
    @synchronized (self) {
        if (![_appWideMessages containsObject:message]) {
            message.dismissDelegate = self;
            [_appWideMessages addObject:message];
        }
        
        _currentMessage = message;
        [self renderForCurrentState];
    }
}

-(void) removeAppWideMessage:(CMBannerMessage*)message {
    @synchronized (self) {
        [_appWideMessages removeObject:message];
        
        if (_currentMessage == message) {
            _currentMessage = _appWideMessages.lastObject;
        }
        
        [self renderForCurrentState];
    }
}

-(void) removeAllAppWideMessages {
    @synchronized (self) {
        [_appWideMessages removeAllObjects];
        _currentMessage = nil;
        [self renderForCurrentState];
    }
}

-(void)setAppWideBannerPosition:(CMAppWideBannerPosition)appWideBannerPosition {
    if (appWideBannerPosition == _appWideBannerPosition) {
        return;
    }
    _appWideBannerPosition = appWideBannerPosition;
    [self removeAppWideBannerContainer];
    [self renderForCurrentState];
}

-(void) renderForCurrentState {
    // Always dispatch async. The caller might call us before the window/rootVC relations are setup.
    dispatch_async(dispatch_get_main_queue(), ^{
        [self renderForCurrentStateSyncMain];
    });
}

-(void) renderForCurrentStateSyncMain {
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
    [_currentMessageView removeFromSuperview];
    
    [self createAppWideBannerContainerIfMissing];
    if (!_appWideContainerView) {
        return;
    }
    
    if (_appWideMessages.count > 1) {
        _currentMessage.nextMessageDelegate = self;
    } else {
        _currentMessage.nextMessageDelegate = nil;
    }
    UIView* messageView = [_currentMessage buildViewForMessage];
    _currentMessageView = messageView;
    messageView.translatesAutoresizingMaskIntoConstraints = NO;
    [_appWideContainerView addSubview:messageView];
    NSArray<NSLayoutConstraint*>* constraints = @[
        [messageView.topAnchor constraintEqualToAnchor:_appWideContainerView.topAnchor],
        [messageView.leftAnchor constraintEqualToAnchor:_appWideContainerView.leftAnchor],
        [messageView.rightAnchor constraintEqualToAnchor:_appWideContainerView.rightAnchor],
        [messageView.bottomAnchor constraintEqualToAnchor:_appWideContainerView.bottomAnchor],
    ];
    
    [NSLayoutConstraint activateConstraints:constraints];
}

-(void) createAppWideBannerContainerIfMissing {
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
    UIWindow* keyWindow = [[[UIApplication sharedApplication] windows] firstObject];
    for (UIWindow* w in [[UIApplication sharedApplication] windows]) {
        if (w.isKeyWindow) {
            keyWindow = w;
            break;
        }
    }
    if (!keyWindow) {
        // no window to render in
        NSLog(@"CriticalMoments: CMBannerManager could not find a key window, aborting.");
        return;
    }
    
    // Add the container view to the key window root VC
    UIViewController* appRootViewController = keyWindow.rootViewController;
    if (appRootViewController.view.window != keyWindow) {
        NSLog(@"CriticalMoments: tried to show a banner before root VC was setup, aborting.");
        return;
    }
    _appWideContainerView = [[UIView alloc] init];
    [keyWindow addSubview:_appWideContainerView];
    
    //
    // Layout
    //
    
    appRootViewController.view.translatesAutoresizingMaskIntoConstraints = NO;
    _appWideContainerView.translatesAutoresizingMaskIntoConstraints = NO;
    
    // These two low priority constraints aligns rootVC to window top/bottom,  but are overridden by high pri banner constraints if present
    NSLayoutConstraint* appAlignBottomWindowLowPriorityConstraint = [appRootViewController.view.bottomAnchor constraintEqualToAnchor:keyWindow.bottomAnchor];
    appAlignBottomWindowLowPriorityConstraint.priority = UILayoutPriorityDefaultLow;
    NSLayoutConstraint* appAlignTopWindowLowPriorityConstraint = [appRootViewController.view.topAnchor constraintEqualToAnchor:keyWindow.topAnchor];
    appAlignTopWindowLowPriorityConstraint.priority = UILayoutPriorityDefaultLow;
    
    NSArray<NSLayoutConstraint*>* constraints = @[
        // position banner to the side edges
        [_appWideContainerView.leftAnchor constraintEqualToAnchor:keyWindow.leftAnchor],
        [_appWideContainerView.rightAnchor constraintEqualToAnchor:keyWindow.rightAnchor],
        
        // Make the banner at most 20% window height. Backstop for way too much text.
        [_appWideContainerView.heightAnchor constraintLessThanOrEqualToAnchor:keyWindow.heightAnchor multiplier:MAX_BANNER_HEIGHT_PERCENTAGE],
        
        // Align root VC to the edges of the window
        appAlignBottomWindowLowPriorityConstraint,
        appAlignTopWindowLowPriorityConstraint,
        [appRootViewController.view.leftAnchor constraintEqualToAnchor:keyWindow.leftAnchor],
        [appRootViewController.view.rightAnchor constraintEqualToAnchor:keyWindow.rightAnchor],
    ];
    
    // Top vs Bottom layout for banner
    if (self.appWideBannerPosition == CMAppWideBannerPositionBottom) {
        // Banner container at bottom of app
        constraints = [constraints arrayByAddingObjectsFromArray:@[
            [_appWideContainerView.bottomAnchor constraintEqualToAnchor:keyWindow.bottomAnchor],
            [appRootViewController.view.bottomAnchor constraintEqualToAnchor:_appWideContainerView.topAnchor],
        ]];
    } else {
        // Banner container at top of app
        constraints = [constraints arrayByAddingObjectsFromArray:@[
            [_appWideContainerView.topAnchor constraintEqualToAnchor:keyWindow.topAnchor],
            [appRootViewController.view.topAnchor constraintEqualToAnchor:_appWideContainerView.bottomAnchor],
        ]];
    }
    
    [NSLayoutConstraint activateConstraints:constraints];
}

-(void) removeAppWideBannerContainer {
    if (![NSThread isMainThread]) {
        dispatch_sync(dispatch_get_main_queue(), ^{
            [self removeAppWideBannerContainer];
        });
        return;
    }
    
    if (!_appWideContainerView) {
        return;
    }
    
    [_appWideContainerView removeFromSuperview];
    _currentMessageView = nil;
    _appWideContainerView = nil;
}

#pragma mark CMBannerDismissDelegate

-(void) dismissedMessage:(CMBannerMessage*)message{
    [self removeAppWideMessage:message];
}

#pragma mark CMBannerNextMessageDelegate

-(void)nextMessage {
    @synchronized (self) {
        NSUInteger nextIndex = [_appWideMessages indexOfObject:_currentMessage] + 1;
        nextIndex = nextIndex % _appWideMessages.count;
        _currentMessage = [_appWideMessages objectAtIndex:nextIndex];
        [self renderForCurrentState];
    }
}

@end
