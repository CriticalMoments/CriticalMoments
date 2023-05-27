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
    horizontalCondition.subtitle =
        @"Evaluate condition: `(interface_orientation == 'landscape')`.\n\nTry "
        @"turning the device landscape/portrait.";
    horizontalCondition.actionCMActionName = @"conditional_landscape";
    [horizontalCondition addResetTestTarget:self
                                     action:@selector(dismissAlerts)];

    CMDemoAction *flatCondition = [[CMDemoAction alloc] init];
    flatCondition.title = @"Device position";
    flatCondition.subtitle =
        @"Evaluate condition: `(device_orientation == 'face_up' || "
        @"device_orientation == 'face_down')`.\n\nTry laying the device flat "
        @"on the table.";
    flatCondition.actionCMActionName = @"conditional_flat";
    [flatCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    CMDemoAction *wifiConditon = [[CMDemoAction alloc] init];
    wifiConditon.title = @"Wifi connection";
    wifiConditon.subtitle = @"Evaluate condition `(network_connection_type == "
                            @"'wifi')`.\n\nTry connecting/disconnecting wifi.";
    wifiConditon.actionCMActionName = @"conditional_wifi";
    [wifiConditon addResetTestTarget:self action:@selector(dismissAlerts)];

    CMDemoAction *chargingCondition = [[CMDemoAction alloc] init];
    chargingCondition.title = @"Charging battery";
    chargingCondition.subtitle =
        @"Evaluate condition: `(device_battery_state == 'charging' || "
        @"device_battery_state == 'full')`.\n\nTry plugging in/removing "
        @"charging cable.";
    chargingCondition.actionCMActionName = @"conditional_charging";
    [chargingCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    [self addSection:@"Simple conditions"
         withActions:@[
             horizontalCondition, flatCondition, wifiConditon, chargingCondition
         ]];

    CMDemoAction *compoundCondition = [[CMDemoAction alloc] init];
    compoundCondition.title = @"Compound condition";
    compoundCondition.subtitle =
        @"Evaluate condition: `(interface_orientation == 'landscape') && "
        @"(network_connection_type == "
        @"'wifi')`.\n\nOnly true when device is landscape and connected to "
        @"wifi.";
    compoundCondition.actionCMActionName = @"conditional_compound";
    [compoundCondition addResetTestTarget:self action:@selector(dismissAlerts)];

    [self addSection:@"Compound conditions" withActions:@[ compoundCondition ]];

    CMDemoAction *osVersion = [[CMDemoAction alloc] init];
    osVersion.title = @"Version Number Check";
    osVersion.subtitle =
        @"Evaluate condition: `(versionNumberComponent(os_version,0) >= "
        @"16)`.\n\nChecks ios version, extracting the major release number "
        @"from the longer release string, eg '16.4.1'";
    osVersion.actionCMActionName = @"conditional_os_version";
    [osVersion addResetTestTarget:self action:@selector(dismissAlerts)];

    [self addSection:@"Conditions with functions" withActions:@[ osVersion ]];
}

- (void)dismissAlerts {
    [Utils.keyWindow.rootViewController dismissViewControllerAnimated:NO
                                                           completion:nil];
}

@end
