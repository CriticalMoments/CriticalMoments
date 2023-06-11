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

    CMDemoAction *flatCondition = [[CMDemoAction alloc] init];
    flatCondition.title = @"Device position";
    flatCondition.subtitle = @"Condition it true if device is laying flat on a "
                             @"table\n\n(device_orientation == 'face_up' || "
                             @"device_orientation == 'face_down')";
    flatCondition.actionCMActionName = @"conditional_flat";
    [flatCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    CMDemoAction *wifiConditon = [[CMDemoAction alloc] init];
    wifiConditon.title = @"Wifi connection";
    wifiConditon.subtitle = @"Condition true if the device's primary network is wifi (not "
                            @"cellular).\n\n(network_connection_type == "
                            @"'wifi')";
    wifiConditon.actionCMActionName = @"conditional_wifi";
    [wifiConditon addResetTestTarget:self action:@selector(dismissAlerts)];

    CMDemoAction *chargingCondition = [[CMDemoAction alloc] init];
    chargingCondition.title = @"Charging battery";
    chargingCondition.subtitle = @"Condition is true if device is charging.\n\n(device_battery_state == "
                                 @"'charging' || "
                                 @"device_battery_state == 'full')";
    chargingCondition.actionCMActionName = @"conditional_charging";
    [chargingCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    [self addSection:@"Simple conditions"
         withActions:@[ horizontalCondition, flatCondition, wifiConditon, chargingCondition ]];

    CMDemoAction *compoundCondition = [[CMDemoAction alloc] init];
    compoundCondition.title = @"Combining conditions";
    compoundCondition.subtitle = @"Condition true if device is landscape and connected to "
                                 @"wifi.\n\n(interface_orientation == 'landscape') && "
                                 @"(network_connection_type == "
                                 @"'wifi')";
    compoundCondition.actionCMActionName = @"conditional_compound";
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
    [complex addResetTestTarget:self action:@selector(dismissAlerts)];

    [self addSection:@"Complex condition" withActions:@[ complex ]];
}

- (void)dismissAlerts {
    [Utils.keyWindow.rootViewController dismissViewControllerAnimated:NO completion:nil];
}

@end
