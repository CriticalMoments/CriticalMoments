//
//  FeatureFlagsDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2024-02-04.
//

#import "FeatureFlagsDemoScreen.h"

@import CriticalMoments;

@implementation FeatureFlagsDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Feature Flags";
        self.infoText = @"This page shows demos of features flags. The titles contain the value of each flag.";
        self.buttonLink = @"https://docs.criticalmoments.io/feature-flags/conditional-feature-flags";
        [self buildSections];
    }
    return self;
}

- (void)buildSections {
    [CriticalMoments.sharedInstance
        checkNamedCondition:@"featureForPhonesNewOs"
                  condition:@"device_model_class == 'iPhone' && !versionLessThan(os_version, '17.0')"
                    handler:^(_Bool result, NSError *_Nullable error) {
                      if (error != nil) {
                          NSLog(@"Issue with feature flag: %@", error);
                          return;
                      }
                      CMDemoAction *action = [[CMDemoAction alloc] init];
                      action.title = result ? @"State: Is Phone with iOS 17+" : @"State: Is NOT Phone with iOS 17+";
                      action.subtitle =
                          @"This feature flag will change state, depending on which device and OS it's run "
                          @"on.\n\nCondition: device_model_class == 'iPhone' &&\n!versionLessThan(os_version, '17.0')";
                      [self addSection:@"Target Device Properties" withActions:@[ action ]];
                    }];

    [CriticalMoments.sharedInstance
        checkNamedCondition:@"abTestGroupForExperimentFive"
                  condition:@"(randForKey('experiment5', stableRand()) % 100) < 25 && !(is_pro_user ?? false)"
                    handler:^(_Bool result, NSError *_Nullable error) {
                      if (error != nil) {
                          NSLog(@"Issue with feature flag: %@", error);
                          return;
                      }
                      CMDemoAction *action = [[CMDemoAction alloc] init];
                      action.title = result ? @"AB Test Group: A" : @"AB Test Group: B";
                      action.subtitle =
                          @"Split users into AB tests including: 1) random asignment, 2) filtering by user properties "
                          @"(is_pro_user), 3) filtering by built-in properties (is_ipad, app_install_date), and much "
                          @"more. These can be remotely updated to rollout or rollback.\n\nExample: "
                          @"randForKey('experiment5', stableRand()) % 100 < 25 && "
                          @"!(is_pro_user ?? false)";
                      [self addSection:@"AB Testing" withActions:@[ action ]];
                    }];

    [CriticalMoments.sharedInstance
        checkNamedCondition:@"weatherExample"
                  condition:@"(weather_approx_location_temperature > 20)"
                    handler:^(_Bool result, NSError *_Nullable error) {
                      if (error != nil) {
                          NSLog(@"Issue with feature flag: %@", error);
                          return;
                      }
                      CMDemoAction *action = [[CMDemoAction alloc] init];
                      action.title = result ? @"Offer: Explore European Vacation Deals"
                                            : @"Offer: Explore Caribbean Vacation Deals";
                      action.subtitle = @"Show different offers to "
                                        @"different users based on local weather. Caribbean when cold, and Europe when "
                                        @"warm.\n\nCondition: (weather_approx_location_temperature > 10)";
                      [self addSection:@"Live Weather Example" withActions:@[ action ]];
                    }];

    [CriticalMoments.sharedInstance
        checkNamedCondition:@"custom_feature_1"
                  condition:@"true"
                    handler:^(_Bool result, NSError *_Nullable error) {
                      if (error != nil) {
                          NSLog(@"Issue with feature flag: %@", error);
                          return;
                      }
                      CMDemoAction *action = [[CMDemoAction alloc] init];
                      action.title =
                          result ? @"Feature Flag: Enabled in Code" : @"Feature Flag: Disabled with Remote Override";
                      action.subtitle =
                          @"This feature was enabled in code, but should be remotely disabled via cloud update. This "
                          @"can be useful for fixing unexpected issues, or rolling out successful experiments.";
                      [self addSection:@"Remote Update" withActions:@[ action ]];
                    }];

    [CriticalMoments.sharedInstance
        checkNamedCondition:@"app_launced_several_times"
                  condition:@"app_install_date < now() - duration('10m')"
                    handler:^(_Bool result, NSError *_Nullable error) {
                      if (error != nil) {
                          NSLog(@"Issue with feature flag: %@", error);
                          return;
                      }
                      CMDemoAction *action = [[CMDemoAction alloc] init];
                      action.title = result ? @"Feature Flag: App Installed Over 10 Mins Ago"
                                            : @"Feature Flag: App Installed in Last 10 Mins";
                      action.subtitle =
                          @"This flag looks at user engagement history to determine it's value. In this case, how long "
                          @"ago the app was installed.\n\nCondition: app_install_date < now() - duration('10m')";
                      [self addSection:@"User Engagment History" withActions:@[ action ]];
                    }];
}

@end
