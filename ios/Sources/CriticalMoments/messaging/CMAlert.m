//
//  CMAlert.m
//
//
//  Created by Steve Cosman on 2023-05-11.
//

#import "CMAlert.h"
#import "CMAlert_private.h"

@import UIKit;

#import "../CriticalMoments_private.h"
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

    if (dataModel.showCancelButton) {
        NSString *cancelString = [CMUtils uiKitLocalizedStringForKey:@"Cancel"];
        UIAlertAction *cancelAction = [UIAlertAction actionWithTitle:cancelString
                                                               style:UIAlertActionStyleCancel
                                                             handler:^(UIAlertAction *_Nonnull action) {
                                                               // Action just for events here.
                                                               // "Cancel" not localized because event names should not
                                                               // be localized -2 constant for Cancel (documented)
                                                               [self buttonTappedWithAction:nil
                                                                                 buttonName:@"Cancel"
                                                                                buttonIndex:-2];
                                                             }];
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
        UIAlertAction *okAction = [UIAlertAction
            actionWithTitle:okString
                      style:UIAlertActionStyleDefault
                    handler:^(UIAlertAction *_Nonnull action) {
                      // Use "OK" since we don't want events localized
                      // use -1 a constant for OK (documented)
                      [self buttonTappedWithAction:self.dataModel.okButtonActionName buttonName:@"OK" buttonIndex:-1];
                    }];
        [alert addAction:okAction];

        // Only highlight ok as primary if there's other buttons.
        // This is an iOS UI standard.
        if (dataModel.showCancelButton || customButtonActions.count > 0) {
            [alert setPreferredAction:okAction];
        }
    }

    dispatch_async(dispatch_get_main_queue(), ^{
      UIViewController *topController = CMUtils.topViewController;
      if (!topController) {
          NSLog(@"CriticalMoments: can't find root vc for presenting alert");
      }

      // if popoverPresentationController is present, we need to set these or it
      // will crash. However, this shouldn't be present in any case we know of,
      // since we don't use action sheet look on iPad
      if (alert.popoverPresentationController) {
          alert.popoverPresentationController.permittedArrowDirections = 0;
          alert.popoverPresentationController.sourceRect =
              CGRectMake(topController.view.center.x, topController.view.center.y, 0, 0);
          alert.popoverPresentationController.sourceView = topController.view;
      }

      [topController presentViewController:alert animated:YES completion:nil];
    });
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
                                                                     [self buttonTappedWithAction:buttonModel.actionName
                                                                                       buttonName:buttonModel.label
                                                                                      buttonIndex:i];
                                                                   }];
        action.isPrimaryAction = [DatamodelAlertActionButtonStyleEnumPrimary isEqualToString:buttonModel.style];

        [customActions addObject:action];
    }

    return customActions;
}

- (void)buttonTappedWithAction:(NSString *)actionName buttonName:(NSString *)btnName buttonIndex:(int)btnIdx {
    if (self.completionEventSender && self.alertName.length > 0) {
        if (btnName) {
            NSString *tapEventName = [NSString stringWithFormat:@"sub-action:%@:button:%@", self.alertName, btnName];
            [self.completionEventSender sendEvent:tapEventName];
        }
        NSString *tapEventIndexName =
            [NSString stringWithFormat:@"sub-action:%@:button_index:%d", self.alertName, btnIdx];
        [self.completionEventSender sendEvent:tapEventIndexName];
    }

    if (actionName.length > 0) {
        [CriticalMoments.sharedInstance performNamedAction:actionName
                                                   handler:^(NSError *_Nullable error) {
                                                     if (error) {
                                                         NSLog(@"CriticalMoments: Alert tap unknown issue: %@", error);
                                                     }
                                                   }];
    }
}

@end
