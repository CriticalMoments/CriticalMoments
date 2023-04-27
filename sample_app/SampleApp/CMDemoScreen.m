//
//  CMDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-26.
//

#import "CMDemoScreen.h"
#import "DemoViewContoller.h"
#import "Utils.h"

@import UIKit;

@interface CMDemoAction ()

@property (nonatomic) id actionTarget;
@property (nonatomic) SEL actionSelector;

@end

@implementation CMDemoAction

- (void)addTarget:(nullable id)target action:(SEL)action {
    self.actionTarget = target;
    self.actionSelector = action;
}

-(void)performAction {
    if (self.actionDelegate) {
        [self.actionDelegate performAction];
    } else if (self.actionNextScreen) {
        [self pushNextScreen];
    } else if (self.actionBlock) {
        self.actionBlock();
    } else if (self.actionTarget && self.actionSelector) {
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Warc-performSelector-leaks"
        [self.actionTarget performSelector:self.actionSelector];
#pragma clang diagnostic pop
    }
}

-(void) pushNextScreen {
    DemoViewContoller* demoVc = [[DemoViewContoller alloc] initWithDemoScreen:self.actionNextScreen];
    UINavigationController* navController;
    UIViewController *rootVC  = Utils.keyWindow.rootViewController;
    if ([rootVC isKindOfClass:[UITabBarController class]]) {
        UITabBarController* tab = (UITabBarController*)rootVC;
        rootVC = tab.selectedViewController;
    }
    if ([rootVC isKindOfClass:[UINavigationController class]]) {
        navController = (UINavigationController*)rootVC;
    } else {
        navController = rootVC.navigationController;
    }
    [navController pushViewController:demoVc animated:YES];
}

@end

@implementation CMDemoSection

-(instancetype)init {
    self = [super init];
    if (self) {
        self.actions = [[NSMutableArray alloc] init];
    }
    return self;
}

@end

@implementation CMDemoScreen

-(instancetype)init {
    self = [super init];
    if (self) {
        self.sections = [[NSMutableArray alloc] init];
    }
    return self;
}

-(void)addActionToRootSection:(CMDemoAction *)action {
    if (_sections.count == 0) {
        CMDemoSection* rootSection = [[CMDemoSection alloc] init];
        [_sections addObject:rootSection];
    }
    CMDemoSection* rootSection = _sections.firstObject;
    [rootSection.actions addObject:action];
}

-(void)addSection:(NSString *)title withActions:(NSArray<CMDemoAction *> *)actions {
    CMDemoSection* section = [[CMDemoSection alloc] init];
    section.title = title;
    NSMutableArray* mutableActions;
    if ([actions isKindOfClass:[NSMutableArray class]]) {
        mutableActions = (NSMutableArray*)actions;
    } else {
        mutableActions = [actions mutableCopy];
    }
    
    section.actions = mutableActions;
    [_sections addObject:section];
}

@end
