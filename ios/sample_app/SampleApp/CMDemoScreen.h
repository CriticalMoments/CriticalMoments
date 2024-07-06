//
//  CMDemoScreen.h
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-26.
//

#import <Foundation/Foundation.h>

#import <UIKit/UIKit.h>

NS_ASSUME_NONNULL_BEGIN

@class CMDemoScreen;

@protocol DemoActionDelegate

- (void)performAction;

@end

@interface CMDemoAction : NSObject

@property(nonatomic, readwrite) NSString *title, *snapshotTitle;
@property(nonatomic, readwrite) NSString *subtitle;
@property(nonatomic, readwrite) bool skipInUiTesting, skipInUI;

// Only should use one of these
@property(nonatomic, readwrite) id<DemoActionDelegate> actionDelegate;
@property(nonatomic, readwrite) CMDemoScreen *actionNextScreen;
@property(nonatomic, readwrite) NSString *actionCMEventName;
@property(nonatomic, readwrite) NSString *actionCMActionName;
@property(nonatomic, copy) void (^actionBlock)(void);
- (void)addTarget:(nullable id)target action:(SEL)action;
- (void)addResetTestTarget:(nullable id)target action:(SEL)action;

- (void)performAction;
- (void)resetForTests;

@end

@interface CMDemoSection : NSObject

@property(nonatomic, readonly) NSString *title;

- (NSArray<CMDemoAction *> *)actions;

@end

@interface CMDemoScreen : NSObject

@property(nonatomic, readwrite) NSString *title;
@property(nonatomic, readwrite) NSString *infoText, *buttonLink, *buttonTitle;

//@property (nonatomic, copy, nullable) void (^willAppear)(void);

- (NSArray<CMDemoSection *> *)sections;

- (void)addSection:(NSString *)section withActions:(NSArray<CMDemoAction *> *)actions;
- (void)addActionToRootSection:(CMDemoAction *)action;
- (void)didAppear:(UIViewController *)vc;

@end

NS_ASSUME_NONNULL_END
