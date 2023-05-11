//
//  CMAlert.m
//
//
//  Created by Steve Cosman on 2023-05-11.
//

#import "CMAlert.h"
#import "CMAlert_private.h"

@import UIKit;

#import "../utils/CMUtils.h"

@interface CMAlert ()

@property(nonatomic, readwrite) DatamodelAlertAction *dataModel;

//@property (nonatomic, readwrite) NSString* title, *message,
//*okButtonActionName, *style;
//@property (nonatomic, readwrite) bool showCancelButton, showOkButton;

/*

 type AlertAction struct {
     ShowCancelButton   bool
     ShowOkButton       bool
     OkButtonActionName string
     Style              string // AlertActionStyleEnum
     CustomButtons      []*AlertActionCustomButton
 }

 type AlertActionCustomButton struct {
     Label      string
     ActionName string
     Style      string // AlertActionButtonStyleEnum
 }
 */

@end

@implementation CMAlert

- (nonnull instancetype)initWithAppcoreDataModel:
    (DatamodelAlertAction *)alertDataModel {
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
    if ([DatamodelAlertActionStyleEnumLarge isEqualToString:dataModel.style]) {
        style = UIAlertControllerStyleActionSheet;
    }

    // TODO: setPreferredAction maybe disable?

    UIAlertController *alert =
        [UIAlertController alertControllerWithTitle:title
                                            message:message
                                     preferredStyle:style];

    NSArray<UIAlertAction *> *customButtonActions = [self customButtonActions];
    for (UIAlertAction *customButtonAction in customButtonActions) {
        [alert addAction:customButtonAction];
        [alert setPreferredAction:customButtonAction];
    }

    if (dataModel.showCancelButton) {
        UIAlertAction *cancelAction =
            [UIAlertAction actionWithTitle:@"Cancel"
                                     style:UIAlertActionStyleCancel
                                   handler:nil];
        [alert addAction:cancelAction];
        [alert setPreferredAction:cancelAction];
    }

    if (dataModel.showOkButton) {
        UIAlertAction *okAction = [UIAlertAction
            actionWithTitle:@"OK"
                      style:UIAlertActionStyleDefault
                    handler:^(UIAlertAction *_Nonnull action) {
                      if (self.dataModel.okButtonActionName.length > 0) {
                          [self
                              performAction:self.dataModel.okButtonActionName];
                      }
                    }];
        [alert addAction:okAction];
        [alert setPreferredAction:okAction];
    }

    UIWindow *keyWindow = [CMUtils keyWindow];
    UIViewController *rootVc = keyWindow.rootViewController;
    if (!rootVc) {
        NSLog(@"CriticalMoments: can't find root vc for presenting alert");
    } else {
        [keyWindow.rootViewController presentViewController:alert
                                                   animated:YES
                                                 completion:nil];
    }
}

- (NSArray<UIAlertAction *> *)customButtonActions {
    NSMutableArray<UIAlertAction *> *customActions = [[NSMutableArray alloc]
        initWithCapacity:self.dataModel.customButtonsCount];
    for (int i = 0; i < self.dataModel.customButtonsCount; i++) {
        DatamodelAlertActionCustomButton *buttonModel =
            [self.dataModel customButtonAtIndex:i];
        if (!buttonModel) {
            continue;
        }

        UIAlertActionStyle style = UIAlertActionStyleDefault;
        if ([DatamodelAlertActionButtonStyleEnumDestructive
                isEqualToString:buttonModel.style]) {
            style = UIAlertActionStyleDestructive;
        }

        UIAlertAction *action = [UIAlertAction
            actionWithTitle:buttonModel.label
                      style:style
                    handler:^(UIAlertAction *action) {
                      if (buttonModel.actionName.length > 0) {
                          [self performAction:buttonModel.actionName];
                      }
                    }];

        [customActions addObject:action];
    }

    return customActions;
}

- (void)performAction:(NSString *)actionName {
    NSError *error;
    [AppcoreSharedAppcore() performNamedAction:actionName error:&error];
    if (error) {
        NSLog(@"CriticalMoments: Alert tap unknown issue: %@", error);
    }
}

@end
