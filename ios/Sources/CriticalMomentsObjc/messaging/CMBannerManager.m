//
//  CMBannerManagger.m
//  
//
//  Created by Steve Cosman on 2023-04-23.
//

#import "CMBannerManager.h"

#define BANNER_HEIGHT 60.0

@interface CMBannerManager ()

@property (nonatomic, strong) NSMutableArray<CMBannerMessage*>* appWideMessages;
@property (nonatomic, strong) UIView* appWideRootView;

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
    }
    return self;
}

-(void) showAppWideMessage:(CMBannerMessage*)message {
    @synchronized (self) {
        if ([_appWideMessages containsObject:message]) {
            return;
        }
        
        [_appWideMessages addObject:message];
    }
}

-(void) removeAppWideMessage:(CMBannerMessage*)message {
    @synchronized (self) {
        
    }
}

-(void) removeAllAppWideMessages {
    @synchronized (self) {
        
    }
}

-(void) createAppWideBannerContainer {
    /*UIView* v = [[CMBannerMessage alloc] initWithBody:@"Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world Helllooo world "];
    //v.backgroundColor = [UIColor greenColor];
    //v.frame = CGRectMake(50,50,50,50);
    v.translatesAutoresizingMaskIntoConstraints = NO;
    v.accessibilityIdentifier = @"banner";
    [w addSubview:v];
    UIView* rootView = rootVc.view;
    
    // Bottom || Top
    if (false) {
        // Bottom
        UIEdgeInsets additionalInset = rootVc.additionalSafeAreaInsets;
        additionalInset.bottom = additionalInset.bottom + BANNER_HEIGHT;
        rootVc.additionalSafeAreaInsets = additionalInset;
        
        NSArray<NSLayoutConstraint*>* constraints = @[
            
            // position below the window and to the edges
            [v.topAnchor constraintEqualToAnchor:rootView.layoutMarginsGuide.bottomAnchor],
            [v.leftAnchor constraintEqualToAnchor:w.leftAnchor],
            [v.rightAnchor constraintEqualToAnchor:w.rightAnchor],
            [v.bottomAnchor constraintEqualToAnchor:w.bottomAnchor],
        ];
        
        [NSLayoutConstraint activateConstraints:constraints];
    } else {
        // Top
        UIEdgeInsets additionalInset = rootVc.additionalSafeAreaInsets;
        additionalInset.top = additionalInset.top + BANNER_HEIGHT;
        rootVc.additionalSafeAreaInsets = additionalInset;
        
        NSArray<NSLayoutConstraint*>* constraints = @[
            
            // position below the window and to the edges
            [v.topAnchor constraintEqualToAnchor:w.topAnchor],
            [v.leftAnchor constraintEqualToAnchor:w.leftAnchor],
            [v.rightAnchor constraintEqualToAnchor:w.rightAnchor],
            [v.bottomAnchor constraintEqualToAnchor:rootView.layoutMarginsGuide.topAnchor],
        ];
        
        [NSLayoutConstraint activateConstraints:constraints];
    }
    
    [rootView setNeedsLayout];*/
}

-(void) removeAppWideBannerContainer {
    
}

@end
