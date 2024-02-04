//
//  ConversionDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2024-02-02.
//

#import "ConversionDemoScreen.h"

@import CriticalMoments;
#import "../Utils.h"

@implementation ConversionDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Conversions & Journey";
        self.infoText = @"Asking users to subscribe or buy before "
                        @"they have experienced the core value of your app is a sure fire way to get them to decline, "
                        @"or even worse, leave for good.\n\nCritical Moments helps you create a user journey:\n1) "
                        @"Ensure they discover your core features\n2) After seeing value, ask them to subscribe or "
                        @"buy\n3) Once they are loyal, ask them to review.\n\nThe example "
                        @"below walk through the user journey of a fictional “todo "
                        @"list” app, including UI nudges for users to progress if they get stuck.\n\nTo restart the "
                        @"demo user journey, kill and relaunch this demo app, which will clear user progress.";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {
    CMDemoAction *aha1 = [[CMDemoAction alloc] init];
    aha1.title = @"Evaluate Goal 1: User to create a list";
    aha1.subtitle =
        @"Check if the user has created their first list, a core feature of the app. If they haven't created "
        @"a list yet, nudge them to try it now with a banner. See the guide in our docs for best practices giving them "
        @"time to "
        @"explore before introducing nudges.\n\neventCount('created_list') == 0";
    aha1.actionCMActionName = @"journey_goal_1";
    [aha1 addResetTestTarget:self action:@selector(resetAll)];

    CMDemoAction *aha1complete = [[CMDemoAction alloc] init];
    aha1complete.title = @"Simulate Achieveing Goal 1: User Creates List";
    aha1complete.subtitle =
        @"Tap this to simulate the user creating a list, completing the first goal of the user journey.";
    aha1complete.actionCMEventName = @"created_list";
    [aha1complete addResetTestTarget:self action:@selector(resetAll)];

    // This action is hidden from UI, but visable in unit tests, creating and e2e test with database state
    CMDemoAction *aha1UiTestAfter = [[CMDemoAction alloc] init];
    aha1UiTestAfter.title = @"Goal 1: Hidden UI Test Goal";
    aha1UiTestAfter.skipInUI = YES;
    aha1UiTestAfter.subtitle = @"";
    aha1UiTestAfter.actionCMActionName = @"journey_goal_1";
    [aha1UiTestAfter addResetTestTarget:self action:@selector(resetAll)];

    [self addSection:@"Goal 1: User Creates a Todo List" withActions:@[ aha1, aha1complete, aha1UiTestAfter ]];

    CMDemoAction *aha2 = [[CMDemoAction alloc] init];
    aha2.title = @"Evaluate Goal 2: Add 2 tasks to lists";
    aha2.subtitle =
        @"Check if the user has added 2 tasks to any list, an important feature of the app. If they "
        @"haven't added any yet, nudge them to try it now with a modal UI.\n\neventCount('created_list') > 0 "
        @"&&\neventCount('created_list_item') < 2";
    aha2.actionCMActionName = @"journey_goal_2";
    [aha2 addResetTestTarget:self action:@selector(resetAll)];

    CMDemoAction *aha2complete = [[CMDemoAction alloc] init];
    aha2complete.title = @"Simulate Goal 2: Add task to list";
    aha2complete.subtitle = @"Tap this to simulate the user adding an task to a list. You'll need to tap this twice to "
                            @"compelte the second goal of the user journey.";
    aha2complete.actionCMEventName = @"created_list_item";
    [aha2complete addResetTestTarget:self action:@selector(resetAll)];

    // These two actions are hidden from UI, but visable in unit tests, creating and e2e test with database state
    CMDemoAction *aha2complete2 = [[CMDemoAction alloc] init];
    aha2complete2.title = @"Goal 2: Simulate User Action Hidden";
    aha2complete2.skipInUI = YES;
    aha2complete2.actionCMEventName = @"created_list_item";
    [aha2complete2 addResetTestTarget:self action:@selector(resetAll)];

    CMDemoAction *aha2UiTestAfter = [[CMDemoAction alloc] init];
    aha2UiTestAfter.title = @"Goal 2: Hidden UI Test Goal";
    aha2UiTestAfter.skipInUI = YES;
    aha2UiTestAfter.actionCMActionName = @"journey_goal_2";
    [aha2UiTestAfter addResetTestTarget:self action:@selector(resetAll)];

    [self addSection:@"Goal 2: User Adds Two Tasks to Lists"
         withActions:@[ aha2, aha2complete, aha2complete2, aha2UiTestAfter ]];

    CMDemoAction *buyAction = [[CMDemoAction alloc] init];
    buyAction.title = @"Evaluate Goal 3: Buy Pro Subscription";
    buyAction.subtitle = @"If the user completed past goals and then completes a task on a list, offer them a trial of "
                         @"the Pro subscription at that momoment. Tapping this button will simulate completing a task "
                         @"on a list.\n\neventCount('created_list') > 0 "
                         @"&&\neventCount('created_list_item') >= 2 &&\neventCount('completed_task') > 0";
    buyAction.actionCMEventName = @"completed_task";
    [buyAction addResetTestTarget:self action:@selector(resetAll)];

    [self addSection:@"Goal 3: Subscribe to Pro Plan" withActions:@[ buyAction ]];

    // TODO:
    // - add limit to how often we show subscribe prompt based on last time shown
    // - add limit to how often we show / total count of shows for buy
    // - Add app review after we hit max "try upgrade", with delay and
}

- (void)resetAll {
    [Utils.keyWindow.rootViewController dismissViewControllerAnimated:NO completion:nil];

    [CriticalMoments.sharedInstance removeAllBanners];
}

@end
