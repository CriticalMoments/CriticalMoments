//
//  CMAlert.m
//
//
//  Created by Steve Cosman on 2023-05-11.
//

#import "CMAlert.h"
#import "CMAlert_private.h"

@import UIKit;

#import "../include/CriticalMoments.h"
#import "../utils/CMUtils.h"

// Wrapper for our custom data
@interface CMCustomAlertButton : UIAlertAction

@property(nonatomic) bool isPrimaryAction;

@end

@implementation CMCustomAlertButton
@end

@interface CMAlert ()

@property(nonatomic, readwrite) DatamodelAlertAction *dataModel;

@end

@implementation CMAlert

- (nonnull instancetype)initWithAppcoreDataModel:(DatamodelAlertAction *)alertDataModel {
    self = [super init];
    if (self) {
        self.dataModel = alertDataModel;
    }
    return self;
}

- (void)showAlert {
    // TODO -- main thread
    UIViewController *topController = CMUtils.topViewController;
    if (!topController) {
        NSLog(@"CriticalMoments: can't find top vc for presenting alert");
    }

    DatamodelAlertAction *dataModel = self.dataModel;
    NSString *title = dataModel.title.length > 0 ? dataModel.title : nil;
    NSString *message = dataModel.message.length > 0 ? dataModel.message : nil;

    UIAlertControllerStyle style = UIAlertControllerStyleAlert;
    BOOL isPhone = UIDevice.currentDevice.userInterfaceIdiom == UIUserInterfaceIdiomPhone;
    // Only use action sheet style on iPhone. On iPad it's visually the same as
    // alert, but without cancel button. The "large" format actually ends up
    // smaller on iPad if we don't do this. Also adds other usability issues
    // (tapping away to dismiss not easy to discoverable without an
    // permittedArrowDirections and position skewed on rotate)
    if (isPhone && [DatamodelAlertActionStyleEnumLarge isEqualToString:dataModel.style]) {
        style = UIAlertControllerStyleActionSheet;
    }

    UIAlertController *alert = [UIAlertController alertControllerWithTitle:title message:message preferredStyle:style];

    // if popoverPresentationController is present, we need to set these or it
    // will crash. However, this shouldn't be present in any case we know of,
    // since we don't use action sheet look on iPad
    if (alert.popoverPresentationController) {
        alert.popoverPresentationController.permittedArrowDirections = 0;
        alert.popoverPresentationController.sourceRect =
            CGRectMake(topController.view.center.x, topController.view.center.y, 0, 0);
        alert.popoverPresentationController.sourceView = topController.view;
    }

    if (dataModel.showCancelButton) {
        NSString *cancelString = [CMUtils uiKitLocalizedStringForKey:@"Cancel"];
        UIAlertAction *cancelAction = [UIAlertAction actionWithTitle:cancelString
                                                               style:UIAlertActionStyleCancel
                                                             handler:nil];
        [alert addAction:cancelAction];
    }

    NSArray<CMCustomAlertButton *> *customButtonActions = [self customButtonActions];
    for (CMCustomAlertButton *customButtonAction in customButtonActions) {
        [alert addAction:customButtonAction];
        if (customButtonAction.isPrimaryAction) {
            [alert setPreferredAction:customButtonAction];
        }
    }

    if (dataModel.showOkButton) {
        NSString *okString = [CMUtils uiKitLocalizedStringForKey:@"OK"];
        UIAlertAction *okAction = [UIAlertAction actionWithTitle:okString
                                                           style:UIAlertActionStyleDefault
                                                         handler:^(UIAlertAction *_Nonnull action) {
                                                           if (self.dataModel.okButtonActionName.length > 0) {
                                                               [self performAction:self.dataModel.okButtonActionName];
                                                           }
                                                         }];
        [alert addAction:okAction];

        // Only highlight ok as primary if there's other buttons.
        // This is an iOS UI standard.
        if (dataModel.showCancelButton || customButtonActions.count > 0) {
            [alert setPreferredAction:okAction];
        }
    }

    [topController presentViewController:alert animated:YES completion:nil];
}

- (NSArray<CMCustomAlertButton *> *)customButtonActions {

    NSMutableArray<CMCustomAlertButton *> *customActions =
        [[NSMutableArray alloc] initWithCapacity:self.dataModel.customButtonsCount];

    for (int i = 0; i < self.dataModel.customButtonsCount; i++) {
        DatamodelAlertActionCustomButton *buttonModel = [self.dataModel customButtonAtIndex:i];
        if (!buttonModel) {
            continue;
        }

        UIAlertActionStyle style = UIAlertActionStyleDefault;
        if ([DatamodelAlertActionButtonStyleEnumDestructive isEqualToString:buttonModel.style]) {
            style = UIAlertActionStyleDestructive;
        }

        CMCustomAlertButton *action = [CMCustomAlertButton actionWithTitle:buttonModel.label
                                                                     style:style
                                                                   handler:^(UIAlertAction *action) {
                                                                     if (buttonModel.actionName.length > 0) {
                                                                         [self performAction:buttonModel.actionName];
                                                                     }
                                                                   }];
        action.isPrimaryAction = [DatamodelAlertActionButtonStyleEnumPrimary isEqualToString:buttonModel.style];

        [customActions addObject:action];
    }

    return customActions;
}

- (void)performAction:(NSString *)actionName {
    NSError *error;
    [CriticalMoments.sharedInstance performNamedAction:actionName error:&error];
    if (error) {
        NSLog(@"CriticalMoments: Alert tap unknown issue: %@", error);
    }
}

@end
