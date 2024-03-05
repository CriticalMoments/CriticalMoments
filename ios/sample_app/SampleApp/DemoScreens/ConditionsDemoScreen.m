//
//  ConditionsDemoScreen.m
//  SampleApp
//
//  Created by Steve Cosman on 2023-05-26.
//

#import "ConditionsDemoScreen.h"

#import "Utils.h"

@implementation ConditionsDemoScreen

- (instancetype)init {
    self = [super init];
    if (self) {
        self.title = @"Conditions Demos";
        self.infoText = @"Powerful conditions for targeting, with over 100 properties";
        self.buttonLink = @"https://docs.criticalmoments.io/conditional-targeting/intro-to-conditions";

        [self buildSections];
    }
    return self;
}

- (void)buildSections {

    // Simple conditions

    CMDemoAction *horizontalCondition = [[CMDemoAction alloc] init];
    horizontalCondition.title = @"Interface Orientation";
    horizontalCondition.subtitle = @"Condition is true of device is in landscape "
                                   @"orientation.\n\n(interface_orientation == 'landscape')";
    horizontalCondition.actionCMActionName = @"conditional_landscape";
    [horizontalCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    // Boring, but add to UI tests only, skip in demo UI
    CMDemoAction *isIpad = [[CMDemoAction alloc] init];
    isIpad.title = @"Is iPad";
    isIpad.subtitle = @"";
    isIpad.actionCMActionName = @"conditional_ipad";
    isIpad.skipInUI = true;
    [isIpad addResetTestTarget:self action:@selector(dismissAlerts)];

    CMDemoAction *warmCondition = [[CMDemoAction alloc] init];
    warmCondition.title = @"Current Temperature";
    warmCondition.subtitle =
        @"Condition true if it over 20 degrees celcius outside.\n\n(weather_approx_location_temperature > 20)";
    warmCondition.actionCMActionName = @"conditional_warm";
    warmCondition.skipInUiTesting = true;
    [warmCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    CMDemoAction *cloudyCondition = [[CMDemoAction alloc] init];
    cloudyCondition.title = @"Cloud Cover";
    cloudyCondition.subtitle = @"Condition true if it cloudy outside.\n\n(weather_approx_location_cloud_cover > 0.5)";
    cloudyCondition.actionCMActionName = @"conditional_cloudy";
    cloudyCondition.skipInUiTesting = true;
    [cloudyCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    CMDemoAction *flatCondition = [[CMDemoAction alloc] init];
    flatCondition.title = @"Device position";
    flatCondition.subtitle = @"Condition it true if device is laying flat on a "
                             @"table\n\n(device_orientation == 'face_up' || "
                             @"device_orientation == 'face_down')";
    flatCondition.actionCMActionName = @"conditional_flat";
    flatCondition.skipInUiTesting = true;
    [flatCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    CMDemoAction *wifiConditon = [[CMDemoAction alloc] init];
    wifiConditon.title = @"Wifi connection";
    wifiConditon.subtitle = @"Condition true if the device's primary network is wifi (not "
                            @"cellular).\n\n(network_connection_type == "
                            @"'wifi')";
    wifiConditon.actionCMActionName = @"conditional_wifi";
    wifiConditon.skipInUiTesting = true;
    [wifiConditon addResetTestTarget:self action:@selector(dismissAlerts)];

    CMDemoAction *chargingCondition = [[CMDemoAction alloc] init];
    chargingCondition.title = @"Charging battery";
    chargingCondition.subtitle = @"Condition is true if device is charging.\n\n(device_battery_state == "
                                 @"'charging' || "
                                 @"device_battery_state == 'full')";
    chargingCondition.actionCMActionName = @"conditional_charging";
    chargingCondition.skipInUiTesting = true;
    [chargingCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    CMDemoAction *geoCondition = [[CMDemoAction alloc] init];
    geoCondition.title = @"Permissionless Location";
    geoCondition.subtitle =
        @"Condition true if this device is in Canada currently. Checked using IP address, without needing location/GPS "
        @"permissions.\n\n(location_approx_country == 'CA')";
    geoCondition.actionCMActionName = @"conditional_canada";
    geoCondition.skipInUiTesting = true;
    [geoCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    CMDemoAction *dateCondition = [[CMDemoAction alloc] init];
    dateCondition.title = @"Installed app in last hour";
    dateCondition.subtitle =
        @"Condition is true if this app was installed in last hour.\n\napp_install_date > now() - hours(1)";
    dateCondition.actionCMActionName = @"conditional_installed_recently";
    dateCondition.skipInUiTesting = true;
    [dateCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    [self addSection:@"Simple conditions"
         withActions:@[
             horizontalCondition, warmCondition, cloudyCondition, flatCondition, wifiConditon, chargingCondition,
             geoCondition, dateCondition, isIpad
         ]];

    CMDemoAction *compoundCondition = [[CMDemoAction alloc] init];
    compoundCondition.title = @"Combining conditions";
    compoundCondition.subtitle = @"Condition true if device is landscape and connected to "
                                 @"wifi.\n\n(interface_orientation == 'landscape') && "
                                 @"(network_connection_type == "
                                 @"'wifi')";
    compoundCondition.actionCMActionName = @"conditional_compound";
    compoundCondition.skipInUiTesting = true;
    [compoundCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    [self addSection:@"Compound conditions" withActions:@[ compoundCondition ]];

    CMDemoAction *osVersion = [[CMDemoAction alloc] init];
    osVersion.title = @"Function Example";
    osVersion.subtitle = @"Checks if running iOS 16 or newer, using a function to extract the "
                         @"major release number "
                         @"from the longer release string, eg "
                         @"'16.4.1'.\n\n(versionNumberComponent(os_version,0) >= "
                         @"16)";
    osVersion.actionCMActionName = @"conditional_os_version";
    [osVersion addResetTestTarget:self action:@selector(dismissAlerts)];

    [self addSection:@"Conditions with functions" withActions:@[ osVersion ]];

    CMDemoAction *complex = [[CMDemoAction alloc] init];
    complex.title = @"Complex condition";
    complex.subtitle = @"Condition true if 2 or more subconditions are met: "
                       @"has wifi, is landscape, and is charging.\n\n"
                       @"(interface_orientation == 'landscape' ? 1 : 0) +\n "
                       @"(network_connection_type == 'wifi' ? 1 : 0) +\n "
                       @"((device_battery_state == 'charging' || "
                       @"device_battery_state == 'full') ? 1 : 0) >= 2";
    complex.actionCMActionName = @"conditional_complex";
    complex.skipInUiTesting = true;
    [complex addResetTestTarget:self action:@selector(dismissAlerts)];

    [self addSection:@"Complex condition" withActions:@[ complex ]];
}

- (void)dismissAlerts {
    [Utils.keyWindow.rootViewController dismissViewControllerAnimated:NO completion:nil];
}

@end
