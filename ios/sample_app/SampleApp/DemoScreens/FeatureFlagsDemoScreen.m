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

- (NSString *)sectionSubtitleForResult:(bool)result withDescription:(NSString *)description {
    return result ? [NSString stringWithFormat:@"✅ Feature Flag True\n\n%@", description]
                  : [NSString stringWithFormat:@"❌ Feature Flag False\n\n%@", description];
}

- (void)buildSections {
    [CriticalMoments.sharedInstance
        checkNamedCondition:@"is_iphone_with_recent_os"
                    handler:^(_Bool result, NSError *_Nullable error) {
                      if (error != nil) {
                          NSLog(@"Issue with feature flag: %@", error);
                          return;
                      }
                      CMDemoAction *action = [[CMDemoAction alloc] init];
                      action.title = @"iPhone with iOS Version >= 17.0?";
                      action.subtitle =
                          [self sectionSubtitleForResult:result
                                         withDescription:
                                             @"This feature flag checks device type (phone vs tablet), and OS version."
                                             @"\n\nCondition: device_model_class == 'iPhone' "
                                             @"&&\n!versionLessThan(os_version, '17.0')"];
                      [self addSection:@"Target Device Properties" withActions:@[ action ]];
                    }];

    [CriticalMoments.sharedInstance
        checkNamedCondition:@"ab_test_group_for_experiment_five"
                    handler:^(_Bool result, NSError *_Nullable error) {
                      if (error != nil) {
                          NSLog(@"Issue with feature flag: %@", error);
                          return;
                      }
                      CMDemoAction *action = [[CMDemoAction alloc] init];
                      action.title = @"AB Test Group Assignment";
                      NSString *subtitle = result ? @"🅰️ Assigned to Group A" : @"🅱️ Assigned to Group B";
                      subtitle =
                          [subtitle stringByAppendingString:
                                        @"\n\nSplit users into AB tests including: 1) random asignment, 2) filtering "
                                        @"by user properties "
                                        @"(is_pro_user), 3) filtering by built-in properties (is_ipad, "
                                        @"app_install_date), and much "
                                        @"more. These can be remotely updated to rollout or rollback.\n\nExample: "
                                        @"randForKey('experiment5', stableRand()) % 100 < 25 && "
                                        @"!(is_pro_user ?? false)"];
                      action.subtitle = subtitle;
                      [self addSection:@"AB Testing" withActions:@[ action ]];
                    }];

    [CriticalMoments.sharedInstance
        checkNamedCondition:@"weather_warm"
                    handler:^(_Bool result, NSError *_Nullable error) {
                      if (error != nil) {
                          NSLog(@"Issue with feature flag: %@", error);
                          return;
                      }
                      CMDemoAction *action = [[CMDemoAction alloc] init];
                      action.title = @"Weather: Temperature > 20";
                      action.subtitle =
                          [self sectionSubtitleForResult:result
                                         withDescription:
                                             @"This feature flag checks ourdoor weather, using GeoIP "
                                             @"location.\n\nCondition: (weather_approx_location_temperature > 20)"];
                      [self addSection:@"Weather Example" withActions:@[ action ]];
                    }];

    [CriticalMoments.sharedInstance
        checkNamedCondition:@"app_not_recently_installed"
                    handler:^(_Bool result, NSError *_Nullable error) {
                      if (error != nil) {
                          NSLog(@"Issue with feature flag: %@", error);
                          return;
                      }
                      CMDemoAction *action = [[CMDemoAction alloc] init];
                      action.title = @"App Installed Over 10 Mins Ago?";
                      action.subtitle = [self sectionSubtitleForResult:result
                                                       withDescription:@"This flag looks at user engagement history to "
                                                                       @"determine it's value. In this case, how long "
                                                                       @"ago the app was installed.\n\nCondition: "
                                                                       @"app_install_date < now() - duration('10m')"];
                      [self addSection:@"User Engagement History" withActions:@[ action ]];
                    }];
}

@end
