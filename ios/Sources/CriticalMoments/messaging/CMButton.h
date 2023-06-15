//
//  CMButtonStack.h
//
//
//  Created by Steve Cosman on 2023-06-15.
//

#import <UIKit/UIKit.h>

#import "../themes/CMTheme.h"

NS_ASSUME_NONNULL_BEGIN

@interface CMButton : UIButton

// TODO data model
+ (UIButton *)buttonWithWithDataModel:(NSObject *)o andTheme:(CMTheme *_Nullable)theme;

@end

NS_ASSUME_NONNULL_END
