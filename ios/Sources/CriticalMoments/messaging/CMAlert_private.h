//
//  CMAlert_private.h
//
//
//  Created by Steve Cosman on 2023-05-11.
//

NS_ASSUME_NONNULL_BEGIN

@import Appcore;

@interface CMAlert ()

// _private header prevents exposing these to public SDK.

/**
 :nodoc:
 @param alertDataModel The appcore datamodel for this alert
 */
- (instancetype)initWithAppcoreDataModel:(DatamodelAlertAction *)alertDataModel;

/// :nodoc:
- (void)showAlert;

@end

NS_ASSUME_NONNULL_END
