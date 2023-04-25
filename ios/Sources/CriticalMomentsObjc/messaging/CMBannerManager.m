//
//  CMBannerManagger.m
//  
//
//  Created by Steve Cosman on 2023-04-23.
//

#import "CMBannerManager.h"

// TODO Dynamic height.
// https://stackoverflow.com/questions/19628568/uilabel-sizethatfits-not-working
// best just fully dnamic, with max height and max-lines? I think so
#define BANNER_HEIGHT 60.0

@interface CMBannerManager () <CMBannerDismissDelegate>

// Access should be @synchronized(self)
@property (nonatomic, strong) NSMutableArray<CMBannerMessage*>* appWideMessages;

// currentMessage managed by renderForCurrentState -- don't modify directly
@property (nonatomic, strong) CMBannerMessage* currentMessage;

// access syncronized by main queue
@property (nonatomic, strong) UIView* appWideContainerView;
@property (nonatomic) UIEdgeInsets addedInsetsForContainerView;

@end

@implementation CMBannerManager

static CMBannerManager *sharedInstance = nil;

+ (CMBannerManager*)sharedInstance
{
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
        self.addedInsetsForContainerView = UIEdgeInsetsZero;
    }
    return self;
}

-(void) showAppWideMessage:(CMBannerMessage*)message {
    @synchronized (self) {
        if ([_appWideMessages containsObject:message]) {
            return;
        }
        
        message.dismissDelegate = self;
        [_appWideMessages addObject:message];

        [self renderForCurrentState];
    }
}

-(void) removeAppWideMessage:(CMBannerMessage*)message {
    @synchronized (self) {
        [_appWideMessages removeObject:message];
        
        [self renderForCurrentState];
    }
}

-(void) removeAllAppWideMessages {
    @synchronized (self) {
        [_appWideMessages removeAllObjects];
        
        [self renderForCurrentState];
    }
}

-(void) renderForCurrentState {
    if (![NSThread isMainThread]) {
        dispatch_sync(dispatch_get_main_queue(), ^{
            [self renderForCurrentState];
        });
        return;
    }

    // Pick a valid new current message, preferring current, then the last added still active message
    CMBannerMessage* priorCurrentMessage = _currentMessage;
    if (![_appWideMessages containsObject:_currentMessage]) {
        _currentMessage = nil;
    }
    _currentMessage = _appWideMessages.lastObject;
    
    // if no messages left to render clear container view
    if (!_currentMessage) {
        [self removeAppWideBannerContainer];
        return;
    }
    
    if (priorCurrentMessage == _currentMessage) {
        // we are already rendering this message, no-op
        return;
    }
    
    // remove prior message from container
    [priorCurrentMessage removeFromSuperview];
    
    if (!_appWideContainerView) {
        [self createAppWideBannerContainer];
    }
    
    _currentMessage.translatesAutoresizingMaskIntoConstraints = NO;
    [_appWideContainerView addSubview:_currentMessage];
    NSArray<NSLayoutConstraint*>* constraints = @[
        [_currentMessage.topAnchor constraintEqualToAnchor:_appWideContainerView.topAnchor],
        [_currentMessage.leftAnchor constraintEqualToAnchor:_appWideContainerView.leftAnchor],
        [_currentMessage.rightAnchor constraintEqualToAnchor:_appWideContainerView.rightAnchor],
        [_currentMessage.bottomAnchor constraintEqualToAnchor:_appWideContainerView.bottomAnchor],
    ];
    
    [NSLayoutConstraint activateConstraints:constraints];
}

-(void) createAppWideBannerContainer {
    //if (dispatch_queue_get_label(dispatch_get_main_queue()) != dispatch_queue_get_label(DISPATCH_CURRENT_QUEUE_LABEL)) {
    // Dispatch UI work to main
    if (![NSThread isMainThread]) {
        dispatch_sync(dispatch_get_main_queue(), ^{
            [self createAppWideBannerContainer];
        });
        return;
    }
    
    if (_appWideContainerView) {
        return;
    }
    
    // TODO better primary -- UIApplication.shared.windows.first(where: { $0.isKeyWindow })?.addSubview(myView)
    // TODO Warning
    UIWindow* keyWindow = [[[UIApplication sharedApplication] windows] firstObject];
    UIViewController* appRootViewController = keyWindow.rootViewController;
    
    _appWideContainerView = [[UIView alloc] init];
    _appWideContainerView.translatesAutoresizingMaskIntoConstraints = NO;
    [keyWindow addSubview:_appWideContainerView];
    
    // Bottom || Top
    if (self.appWideBannerPosition == CMAppWideBannerPositionBottom) {
        // Container at bottom of app
        UIEdgeInsets rootVcInsets = appRootViewController.additionalSafeAreaInsets;
        rootVcInsets.bottom = rootVcInsets.bottom + BANNER_HEIGHT;
        _addedInsetsForContainerView.bottom = _addedInsetsForContainerView.bottom  + BANNER_HEIGHT;
        appRootViewController.additionalSafeAreaInsets = rootVcInsets;
        
        // TODO dynamic height here for many rows possible?
        NSArray<NSLayoutConstraint*>* constraints = @[
            // position below the window and to the edges
            [_appWideContainerView.topAnchor constraintEqualToAnchor:appRootViewController.view.layoutMarginsGuide.bottomAnchor],
            [_appWideContainerView.leftAnchor constraintEqualToAnchor:keyWindow.leftAnchor],
            [_appWideContainerView.rightAnchor constraintEqualToAnchor:keyWindow.rightAnchor],
            [_appWideContainerView.bottomAnchor constraintEqualToAnchor:keyWindow.bottomAnchor],
        ];
        
        [NSLayoutConstraint activateConstraints:constraints];
    } else {
        // Container at top of app
        UIEdgeInsets rootVcInsets = appRootViewController.additionalSafeAreaInsets;
        rootVcInsets.top = rootVcInsets.top + BANNER_HEIGHT;
        _addedInsetsForContainerView.top = _addedInsetsForContainerView.top  + BANNER_HEIGHT;
        appRootViewController.additionalSafeAreaInsets = rootVcInsets;
        
        NSArray<NSLayoutConstraint*>* constraints = @[
            // position above the window and to the edges
            [_appWideContainerView.topAnchor constraintEqualToAnchor:keyWindow.topAnchor],
            [_appWideContainerView.leftAnchor constraintEqualToAnchor:keyWindow.leftAnchor],
            [_appWideContainerView.rightAnchor constraintEqualToAnchor:keyWindow.rightAnchor],
            [_appWideContainerView.bottomAnchor constraintEqualToAnchor:appRootViewController.view.layoutMarginsGuide.topAnchor],
        ];
        
        [NSLayoutConstraint activateConstraints:constraints];
    }
}

// TODO main method dispatch
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
    
    UIWindow* keyWindow = _appWideContainerView.window;
    UIViewController* appRootViewController = keyWindow.rootViewController;
    UIEdgeInsets rootVcInsets = appRootViewController.additionalSafeAreaInsets;
    rootVcInsets.bottom = rootVcInsets.bottom - _addedInsetsForContainerView.bottom;
    _addedInsetsForContainerView.bottom = 0.0;
    rootVcInsets.top = rootVcInsets.top - _addedInsetsForContainerView.top;
    _addedInsetsForContainerView.top = 0.0;
    appRootViewController.additionalSafeAreaInsets = rootVcInsets;
    [_appWideContainerView removeFromSuperview];
    _appWideContainerView = nil;
}

#pragma mark

-(void) dismissedMessage:(CMBannerMessage*)message{
    [self removeAppWideMessage:message];
}

@end
