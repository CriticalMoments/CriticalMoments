//
//  CMSheetViewController.h
//
//
//  Created by Steve Cosman on 2023-06-14.
//

#import <UIKit/UIKit.h>

#import "../themes/CMTheme.h"
#import "../utils/CMEventSender.h"

@import Appcore;

NS_ASSUME_NONNULL_BEGIN

@interface CMModalViewController : UIViewController

- (instancetype)initWithDatamodel:(DatamodelModalAction *)model;

- (instancetype)init NS_UNAVAILABLE;

/// For sending events when actions are performed in the modal
@property(nonatomic, weak, readwrite) id<CMEventSender> completionEventSender;
@property(nonatomic, strong, readwrite) NSString *modalName;

@end

NS_ASSUME_NONNULL_END
