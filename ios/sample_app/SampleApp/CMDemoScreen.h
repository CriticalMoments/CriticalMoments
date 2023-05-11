//
//  CMDemoScreen.h
//  SampleApp
//
//  Created by Steve Cosman on 2023-04-26.
//

#import <Foundation/Foundation.h>

NS_ASSUME_NONNULL_BEGIN

@class CMDemoScreen;

@protocol DemoActionDelegate

- (void)performAction;

@end

@interface CMDemoAction : NSObject

@property(nonatomic, readwrite) NSString *title;
@property(nonatomic, readwrite) NSString *subtitle;

// Only should use one of these
@property(nonatomic, readwrite) id<DemoActionDelegate> actionDelegate;
@property(nonatomic, readwrite) CMDemoScreen *actionNextScreen;
@property(nonatomic, readwrite) NSString *actionCMEventName;
@property(nonatomic, readwrite) NSString *actionCMActionName;
@property(nonatomic, copy) void (^actionBlock)(void);
- (void)addTarget:(nullable id)target action:(SEL)action;

- (void)performAction;

@end

@interface CMDemoSection : NSObject

@property(nonatomic, readonly) NSString *title;

- (NSArray<CMDemoAction *> *)actions;

@end

@interface CMDemoScreen : NSObject

@property(nonatomic, readwrite) NSString *title;

- (NSArray<CMDemoSection *> *)sections;

- (void)addSection:(NSString *)section
       withActions:(NSArray<CMDemoAction *> *)actions;
- (void)addActionToRootSection:(CMDemoAction *)action;

@end

NS_ASSUME_NONNULL_END
