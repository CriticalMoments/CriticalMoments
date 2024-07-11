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
@import CriticalMoments;

@interface CMDemoAction ()

@property(nonatomic) id actionTarget, resetTestTarget;
@property(nonatomic) SEL actionSelector, resetTestSelector;

@end

@implementation CMDemoAction

- (void)addTarget:(nullable id)target action:(SEL)action {
    self.actionTarget = target;
    self.actionSelector = action;
}

- (void)addResetTestTarget:(nullable id)target action:(SEL)action {
    self.resetTestTarget = target;
    self.resetTestSelector = action;
}

- (void)performAction {
    if (self.actionDelegate) {
        [self.actionDelegate performAction];
    } else if (self.actionNextScreen) {
        [self pushNextScreen];
    } else if (self.actionCMEventName) {
        [CriticalMoments.sharedInstance sendEvent:self.actionCMEventName];
    } else if (self.actionCMActionName) {
        [CriticalMoments.sharedInstance performNamedAction:self.actionCMActionName
                                                   handler:^(NSError *_Nullable error) {
                                                     if (error) {
                                                         NSLog(@"SampleApp: Menu tap action unknown issue: %@", error);
                                                     }
                                                   }];
    } else if (self.actionBlock) {
        self.actionBlock();
    } else if (self.actionTarget && self.actionSelector) {
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Warc-performSelector-leaks"
        [self.actionTarget performSelector:self.actionSelector];
#pragma clang diagnostic pop
    }
}

- (void)resetForTests {
    if (self.resetTestTarget && self.resetTestSelector) {
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Warc-performSelector-leaks"
        [self.resetTestTarget performSelector:self.resetTestSelector];
#pragma clang diagnostic pop
    }
}

- (void)pushNextScreen {
    DemoViewContoller *demoVc = [[DemoViewContoller alloc] initWithDemoScreen:self.actionNextScreen];
    UINavigationController *navController = [Utils appNavControl];
    [navController pushViewController:demoVc animated:YES];
}

@end

@interface CMDemoSection ()

@property(nonatomic, readwrite) NSString *title;
@property(nonatomic, readwrite) NSMutableArray<CMDemoAction *> *actionsList;

@end

@implementation CMDemoSection

- (instancetype)init {
    self = [super init];
    if (self) {
        self.actionsList = [[NSMutableArray alloc] init];
    }
    return self;
}

- (NSArray<CMDemoAction *> *)actions {
    return self.actionsList;
}

@end

@interface CMDemoScreen ()

@property(nonatomic, readwrite) NSMutableArray<CMDemoSection *> *sections;
@property(nonatomic, readwrite) NSMutableArray<CMDemoAction *> *actions;

@end

@implementation CMDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.sections = [[NSMutableArray alloc] init];
    }
    return self;
}

- (void)addActionToRootSection:(CMDemoAction *)action {
    if (_sections.count == 0) {
        CMDemoSection *rootSection = [[CMDemoSection alloc] init];
        [_sections addObject:rootSection];
    }
    CMDemoSection *rootSection = _sections.firstObject;
    [rootSection.actionsList addObject:action];
}

- (void)addSection:(NSString *)title withActions:(NSArray<CMDemoAction *> *)actions {
    CMDemoSection *section = [[CMDemoSection alloc] init];
    section.title = title;
    NSMutableArray *mutableActions;
    if ([actions isKindOfClass:[NSMutableArray class]]) {
        mutableActions = (NSMutableArray *)actions;
    } else {
        mutableActions = [actions mutableCopy];
    }

    section.actionsList = mutableActions;
    [_sections addObject:section];
}

- (void)didAppear:(UIViewController *)vc {
}

@end
